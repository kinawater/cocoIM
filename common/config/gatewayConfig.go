package config

import "github.com/spf13/viper"

func GatewayConfigInit(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// GetGatewayMaxTcpNum 当前服务能承受的最大的tcp连接数量
func GetGatewayMaxTcpNum() int32 {
	return viper.GetInt32("gateway.tcp_max_num")
}

// GetGatewayEpollChanNum 单个服务器繁忙时刻，能接受积压不处理的连接最大数量
// 超过该数量的连接，需要等待
func GetGatewayEpollChanNum() int32 { return viper.GetInt32("gateway.epoll_channel_max_size") }

// GetGatewaySingleEpollNum 作为二级reactor的epoll的数量限制
// 和cpu内核相关，比如cpu是英特尔i7-7700，是4核8线程，那就设置成8
func GetGatewaySingleEpollNum() int { return viper.GetInt("gateway.single_epoll_max_num") }

// GetGatewayWaitQueueSize epoll在获取事件的时候，准备的事件列表的长度
func GetGatewayWaitQueueSize() int { return viper.GetInt("gateway.epoll_wait_queue_size") }

// GetGatewayServerPort gateway服务器通信端口
func GetGatewayServerPort() int { return viper.GetInt("gateway.server_port") }

// GetGatewayWorkerPoolNum 线程池配置
func GetGatewayWorkerPoolNum() int { return viper.GetInt("gateway.worker_pool_num") }
