package gateway

import (
	"cocoIM/common/config"
	"fmt"
	"github.com/panjf2000/ants"
)

var wPool *ants.Pool

// 初始化一个pool
func initWorkPoll() {
	var err error
	if wPool, err = ants.NewPool(config.GetGatewayWorkerPoolNum()); err != nil {
		fmt.Printf("InitWorkPoll.err:%s num:%d \n", err.Error(), config.GetGatewayEpollChanNum())
	}
}
