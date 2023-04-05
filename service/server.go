package service

import (
	"cocoIM/pkg/logger"
	"cocoIM/pkg/wsCommProto"
	"cocoIM/service/handle"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"net/http"
	"sync"
	"time"
)

var TimeFormatLayOut = "2006-01-02 15:04:05"

type ServerOptions struct {
	writeWait time.Duration //write timeout line
	readWait  time.Duration //read timeout line
}

type Server struct {
	once    sync.Once
	options ServerOptions
	id      string //服务器ID
	address string //服务器地址
	users   handle.ConnectUserList
	sync.Mutex
}

// 初始化
func NewServer(id, address string) *Server {
	return &Server{
		id:      id,
		address: address,
		users:   handle.ConnectUserList{},
		options: ServerOptions{
			writeWait: time.Second * 10,
			readWait:  time.Minute * 2,
		},
	}
}

func (s *Server) Start() error {
	// 监听服务
	mux := http.NewServeMux()
	logger.Info(map[string]string{
		"module": "server",
		"listen": s.address,
		"id":     s.id,
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 升级ws
		upgradeConn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			_ = upgradeConn.Close()
			logger.Error("ws升级失败")
			return
		}
		//读取userId
		user := r.URL.Query().Get("user")
		if user == "" {
			_ = upgradeConn.Close()
			logger.Error("获取userid失败")
			return
		}
		// 添加到会话列表
		oldConnectUser, ok := s.users.AddUser(handle.ConnectUser{
			UserID:  user,
			Connect: upgradeConn,
		})
		if ok {
			// 断开旧链接，同一时间只能运行一个地方登录
			oldConnectUser.Connect.Close()
			logger.Info(fmt.Printf("user: %s ,oldConnet address: %v", oldConnectUser.Connect.RemoteAddr()))
		}
		// 独立协程处理
		go func(user string, conn net.Conn) error {
			for {
				// 心跳检测
				// 设置读取超时时间，readWait之后认为超时
				_ = conn.SetReadDeadline(time.Now().Add(s.options.readWait))
				// 读帧
				readFrame, err := ws.ReadFrame(conn)
				if err != nil {
					return err
				}
				switch readFrame.Header.OpCode {
				case ws.OpPing:
					// 返回一个pong
					err := wsutil.WriteServerMessage(conn, ws.OpPing, nil)
					if err != nil {
						return err
					}
					// 这里暂时用这个
					logger.Info(fmt.Printf("user : %s , wirte pong...", user))
					continue
				case ws.OpClose:
					// 连接已经关闭
					return errors.New(fmt.Sprintf("user:%s", user))
				}
				logger.Info(fmt.Printf("user:%s , %s , Header : %s", user, time.Now().Format(TimeFormatLayOut), readFrame.Header))
				// 正常数据
				// 是否采用的掩码
				if readFrame.Header.Masked {
					// 解码
					ws.Cipher(readFrame.Payload, readFrame.Header.Mask, 0)
				}

				switch readFrame.Header.OpCode {
				case ws.OpText:
					go s.handleText(user, string(readFrame.Payload))
				case ws.OpBinary:
					go s.handleBinary(user, string(readFrame.Payload))
				}
			}
		}(user, upgradeConn)

	})
	logger.Info("started")
	return http.ListenAndServe(s.address, mux)
}
func (s *Server) shutdown() {
	s.once.Do(func() {
		s.Lock()
		defer s.Unlock()
		s.users.Range(func(_, userConn any) bool {
			userConn.(handle.ConnectUser).Connect.Close()
			return true
		})
	})
}

// 普通文本处理
func (s *Server) handleText(user string, message string) {
	// 群发广播
	s.massSend(user, message)
	// TODO 单聊
}

const (
	CommandPing = 100
	CommandPong = 101
)

func (s *Server) handleBinary(user, message string) {
	proto := wsCommProto.CommProto{}
	err := proto.UnMarshal([]byte(message))
	if err != nil {
		logger.Error(fmt.Printf("user : %s message: %s ,解析失败,err: %v", user, message, err))
	}
	switch proto.Command {
	case CommandPing:
		u, ok := s.users.Load(user)
		if ok {
			nowUser := u.(handle.ConnectUser)
			returnProto := wsCommProto.CommProto{
				Command: CommandPong,
			}
			pongBytes, err := returnProto.Marshal()
			if err != nil {
				logger.Error(fmt.Printf("user: %s marshal pong error:%v", user, err))
			}
			err = wsutil.WriteServerBinary(nowUser.Connect, pongBytes)
			if err != nil {
				logger.Error(fmt.Printf("user: %s ,write pong error:%v", user, err))
			}
		}
	}
}
func (s *Server) massSend(user, message string) {
	s.users.Range(func(key, value any) bool {
		if key == user {
			// 不要发给自己
			return true
		}
		nowUser := value.(handle.ConnectUser)
		err := s.writeText(nowUser.Connect, message)
		if err != nil {
			logger.Error(fmt.Printf("FromUser:%s 广播 ToUser:%s 失败 :%v", user, key, err))
		}
		return true
	})
}

func (s *Server) writeText(conn net.Conn, message string) error {
	// 创建一个帧
	frame := ws.NewTextFrame([]byte(message))
	// 设置超时时间
	err := conn.SetWriteDeadline(time.Now().Add(s.options.writeWait))
	if err != nil {
		return err
	}
	return ws.WriteFrame(conn, frame)
}
