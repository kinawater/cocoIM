package ipconf

import "cocoIM/common/config"

// 启动一个web容器，用来成为一个独立的服务，便于请求端请求

func RunMain(path string) {
	config.Init(path)
	// TODO: 初始化数据源
	// TODO: 初始化调度层
	// 启动一个server
}
