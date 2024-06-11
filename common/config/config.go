package config

import (
	"cocoIM/constants"
	"github.com/spf13/viper"
	"time"
)

// 使用viper封装一个简单的配置读取

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// 获取服务发现的地址
func GetEndpointsForDiscovery() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

// 获取连接服务器发现集群的超时时间 单位:秒
func GetTimeoutForDiscovery() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

func GetServicePathForIPConf() string {
	return viper.GetString("ip_conf.service_path")
}

// 判断debug参数，是否是debug环境
func IsDebug() bool {
	env := viper.GetString("global.env")
	return env == string(constants.Env.Debug)
}
