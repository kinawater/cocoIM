package gateway

import (
	"cocoIM/common/config"
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

}

func runProc(c *connection, ep *epollSingle) {

}
