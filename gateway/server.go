package gateway

import (
	"cocoIM/common/config"
	"cocoIM/common/tcp"
	"errors"
	"io"
	"log"
	"net"
)

func RunMain(path string) {
	config.GatewayConfigInit(path)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayServerPort()})
	if err != nil {
		log.Fatalf("StartTCPPollServer err:%s", err.Error())
		panic(err)
	}
	// TODO 一会删掉
	_ = listener

}

func runProc(c *connection, ep *epollSingle) {
	// 读取完整信息
	data, err := tcp.ReadData(c.conn)
	if err != nil {
		// 如果错误是EOF，那说明连接已经关闭，那么将会让连接剔除
		if errors.Is(err, io.EOF) {
			ep.remove(c)
		}
		return
	}
	// TODO 一会删掉
	_ = data

	return
}
