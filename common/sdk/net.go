package sdk

import (
	"cocoIM/common/tcp"
	"encoding/json"
	"fmt"
	"net"
)

// 定义连接结构体
type connect struct {
	//连接
	conn *net.TCPConn
	//发送，回复消息
	sendChan, recvChan chan *Message
}

// 创建新连接
func newConnet(ip net.IP, port int) *connect {
	clientConn := &connect{
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
	}
	// 创建一个TCP address
	addr := &net.TCPAddr{IP: ip, Port: port}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTCP.err= %+v", err)
	}

	clientConn.conn = conn
	// 循环读取连接数据
	go func() {
		for {
			data, err := tcp.ReadData(conn)
			if err != nil {
				fmt.Printf("client ReadData err:%+v \n", err)
			}
			msg := &Message{}
			_ = json.Unmarshal(data, msg)
			msg.Content = "服务器返回给你的:" + msg.Content
			clientConn.recvChan <- msg
		}
	}()
	return clientConn
}

func (c *connect) send(data *Message) {
	bytes, _ := json.Marshal(data)

	dataPacket := tcp.DataPacket{
		Data: bytes,
		Len:  uint32(len(bytes)),
	}
	marshalData, _ := dataPacket.Marshal()
	_, _ = c.conn.Write(marshalData)

}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

// 关闭连接，回收资源
func (c *connect) close() {
	// 暂时没有可处理的
	c.close()
}
