package mytest

import (
	"cocoIM/common/sdk"
	"net"
)

var TcpConnNum int64

func TestConnRunMain() {
	for i := 0; i < int(TcpConnNum); i++ {
		sdk.MakeNewChat(net.ParseIP("127.0.0.1"), 8900, "logic", "1115858", "1212")
	}
}
