package gateway

import (
	"cocoIM/common/config"
	"cocoIM/common/tcp"
	"errors"
	"fmt"
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
	initWorkPoll()
	initEpoll(listener, runProc)
	fmt.Println("---------------im gateway stated----------------")
	// 这里阻塞掉
	select {}
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
	err = wPool.Submit(func() {
		packetData := tcp.DataPacket{
			Len:  uint32(len(data)),
			Data: data,
		}
		marshalPacketData, _ := packetData.Marshal()
		_ = tcp.WriteData(marshalPacketData, c.conn)
	})
	if err != nil {
		fmt.Errorf("runProc:err:%v\n", err.Error())
	}
	return
}
