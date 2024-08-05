package gateway

import (
	"cocoIM/common/config"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

// 设定epoll对象
type ePoll struct {
	// 即将被分配的连接通道
	eChan chan *connection
	// 所有的fd->连接都存在一个大map里
	tables sync.Map
	// 根据cpu的核心数，限制同时能执行的携程数量
	eSize int
	// 优雅退出时候需要的信号通道
	done chan struct{}
	// 监听链接的listener
	listener *net.TCPListener
	// 处理accept连接的handleFun
	fc func(c *connection, e *epollSingle)
}

// 抽象出来的epoll所操作的单个对象
type epollSingle struct {
	fd int
}

// 当前服务器连接数(动态)
var tcpNum int32
var topEpoll *ePoll

func initEpoll(listener *net.TCPListener, f func(c *connection, ep *epollSingle)) {
	// 设置下linux的文件数量限制
	setLimit()
	topEpoll = newEPool(listener, f)
	topEpoll.createTopReactorForAcceptProcess()
	topEpoll.startEpoll()
}

func newEPool(listener *net.TCPListener, handleFunc func(c *connection, ep *epollSingle)) *ePoll {
	return &ePoll{
		eChan:    make(chan *connection, config.GetGatewayEpollChanNum()),
		tables:   sync.Map{},
		eSize:    config.GetGatewaySingleEpollNum(),
		done:     make(chan struct{}),
		listener: listener,
		fc:       handleFunc,
	}
}

// NewEpollSingle 创建一个子epoll，主要作用是作为subReactor
func NewEpollSingle() (*epollSingle, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epollSingle{fd: fd}, err
}

//总体实现应该是这样的
//1,第一个reactor开始接收连接，并保存到一个channel里，在此要做几个判断，一个是让tcp连接变成长连接，一个是不能超过设定的当前实例最大连接数
//2,按照cpu的核心数，启动数量相同的epoll，在此我们称之为子epoll，用来作为subReactor。
// ● 子epoll会从channel里接管这个fd，把fd注册到子epoll里，这样就能使用epoll来监听事件，和正常的监听一样用
// ● 同时这个子epoll还要监听epoll事件，调用处理器处理

//----------------第一层reactor创立（父级）---------------------

func (ep *ePoll) createTopReactorForAcceptProcess() {
	// 根据内核数量
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				// 拿出一个listener里面监听得到的accept连接
				acceptConn, err := ep.listener.AcceptTCP()
				// 判断是否需要触发熔断
				if !checkBreakPoint() {
					_ = acceptConn.Close()
					// TODO:这里其实应该触发熔断后主动上报到监测，但是这个上报频率和上报接口尚未实现
					// 即使实现也不影响当前的
				}
				// 给tcp连接增加配置
				setTcpConfig(acceptConn)
				if err != nil {
					var ErrNet net.Error
					errors.As(err, &ErrNet)
					if ErrNet.Timeout() || ErrNet.Temporary() {
						fmt.Println("临时错误，不用在意")
						continue
					}
					fmt.Errorf("accept error:%v", err)
				}
				// 拿到socket的fd并缓存
				c := connection{
					conn: acceptConn,
					fd:   getSocketFD(acceptConn),
				}
				ep.addTask(&c)
			}
		}()
	}
}

// 创立一组二级reactor，也就是epoll
func (e *ePoll) startEpoll() {
	for i := 0; i < e.eSize; i++ {
		go e.creatSingleSecondEPollAndHandleEvent()
	}
}

// 创立单个二级epoll，并监听事件，调用处理器
func (e *ePoll) creatSingleSecondEPollAndHandleEvent() {
	sEP, err := NewEpollSingle()
	if err != nil {
		// epoll都创建失败，还玩毛，停了吧
		panic(err)
	}
	// 创建监听事件
	go func() {
		for {
			select {
			case <-e.done:
				return
			case conn := <-e.eChan:
				// 计数器加1，方便统计和判断是否超过最大值
				addTcpNum()
				if err = sEP.add(conn); err != nil {
					fmt.Printf("增加连接到epoll失败（二级），err：%v \r\n", err)
					err := conn.Close()
					if err != nil {
						fmt.Printf("关闭连接失败了，err： %v \r\n", err)
						// 关闭连接都能失败，应该是出了更严重的问题
						panic(err)
					}
					continue
				}
				fmt.Printf("SingleEpollPool new connection[%v] tcpSize:%d \r\n", conn.RemoteAddr(), tcpNum)
			}
		}
	}()

	// 这里要使用epoll的wait来等待event，并调用对应的函数处理

	for {
		select {
		case <-e.done:
			return
		default:
			conns, err := sEP.wait(200)
			if err != nil && err != syscall.EINTR {
				fmt.Printf("failed to topEpoll wait %v \r\n", err)
				continue
			}
			for _, conn := range conns {
				// 调用函数处理
				if conn == nil {
					break
				}
				e.fc(conn, sEP)
			}
		}
	}
}

// 检查是否会熔断
func checkBreakPoint() bool {
	// 是否超过最大连接数
	if !checkTcpNumAsConfigMax() {
		return false
	}
	return true
}

// ************连接相关配置*******************
func setTcpConfig(conn *net.TCPConn) {
	var err error
	err = conn.SetKeepAlive(true)

	// 整体汇报到日志里,但是现在日志没有实现
	// TODO:汇报日志,临时写个输出
	fmt.Println(err)
}

// ************服务器tcp连接数相关操作************

// 获取当前连接数
func getTcpNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}

// 增加一个连接数
func addTcpNum() int32 {
	return atomic.AddInt32(&tcpNum, 1)
}

// 删除一个连接数
func subTcpNum() int32 {
	return atomic.AddInt32(&tcpNum, -1)
}

// 检查连接数是否超过配置的最大限制
func checkTcpNumAsConfigMax() bool {
	configTcpMaxNum := config.GetGatewayMaxTcpNum()
	nowTcpNum := getTcpNum()
	return nowTcpNum <= configTcpMaxNum
}

// ****************获取socket关键信息*****************

// 获取socket的句柄fd
func getSocketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	sysfd := pfdVal.FieldByName("Sysfd").Int()
	return int(sysfd)
}

// *************epoll相关操作***********************

// 把第一个reactor得到的连接放入channel里，等待二级reactor分配处理
func (e *ePoll) addTask(c *connection) {
	e.eChan <- c
}

// *************二级epoll相关操作

// 增加一个event到epoll的监听事件队列里
func (ep *epollSingle) add(conn *connection) error {
	fd := conn.fd
	err := unix.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	// 加入全局管理
	topEpoll.tables.Store(fd, conn)
	return nil
}

func (ep *epollSingle) wait(timeout int) (conns []*connection, err error) {
	// 准备一个events 事件数组
	events := make([]unix.EpollEvent, config.GetGatewayWaitQueueSize())
	n, err := unix.EpollWait(ep.fd, events, timeout)
	// 有事件的conn集合
	for i := 0; i < n; i++ {
		if conn, ok := topEpoll.tables.Load(int(events[i].Fd)); ok {
			conns = append(conns, conn.(*connection))
		}
	}
	return
}

// ******************** 系统相关设置 *********************
func setLimit() {
	// 定义一个 linuxRLimit 结构体变量，用于存储文件描述符限制
	var linuxRLimit syscall.Rlimit
	// 获取当前进程的文件描述符限制，并将其存储在 rLimit 变量中
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &linuxRLimit); err != nil {
		panic(err)
	}
	// 将当前限制设置为最大值
	linuxRLimit.Cur = linuxRLimit.Max
	// 把修改后的值重新设置回去
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &linuxRLimit); err != nil {
		panic(err)
	}
	log.Printf("set cur limit: %d", linuxRLimit)
}
