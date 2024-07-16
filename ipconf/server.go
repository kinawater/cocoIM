package ipconf

import (
	"cocoIM/common/config"
	"cocoIM/ipconf/domain"
	"cocoIM/ipconf/source"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// 启动一个web容器，用来成为一个独立的服务，便于请求端请求

func RunMain(path string) {
	config.Init(path)
	// 初始化数据源
	source.Init()
	// 初始化调度层
	domain.Init()
	// 启动一个server
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
