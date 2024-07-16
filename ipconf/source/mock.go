package source

import (
	"cocoIM/common/config"
	"cocoIM/common/discovery"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func testServiceRegister(ctx *context.Context, port, node string) {
	// 模拟服务发现
	go func() {
		discoverEnd := discovery.EndpointInfo{
			IP:   "127.0.0.1",
			Port: port,
			MetaData: map[string]any{
				"connect_num": float64(rand.Int63n(1000000000000000000)),
				"message_num": float64(rand.Int63n(1000000000000)),
			},
		}
		register, err := discovery.NewServiceRegister(ctx, fmt.Sprintf("%s/%s", config.GetServicePathForIPConf(), node), &discoverEnd, time.Now().Unix())
		if err != nil {
			panic(err)
		}
		go register.ListenLeaseRespChan()
		// 模拟已经注册的实例，每过一秒上报一次自己的服务器状态
		for {
			discoverEnd := discovery.EndpointInfo{
				IP:   "127.0.0.1",
				Port: port,
				MetaData: map[string]any{
					"connect_num": float64(rand.Int63n(1000000000000000000)),
					"message_num": float64(rand.Int63n(1000000000000)),
				},
			}
			register.UpdateValue(&discoverEnd)
			time.Sleep(1 * time.Second)
		}
	}()
}
