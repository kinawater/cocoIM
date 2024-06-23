package source

import (
	"cocoIM/common/config"
	"cocoIM/common/discovery"
	"context"
	"github.com/bytedance/gopkg/util/logger"
)

// todo: init

func Init() {
	eventChan = make(chan *Event)
	ctx := context.Background()
	go DataHandler(&ctx)

}

// 服务发现处理
func DataHandler(ctx *context.Context) {
	serviceDiscovery := discovery.NewServiceDiscovery(ctx)
	defer serviceDiscovery.Close()
	setFunc := func(key, value string) {
		if endpoint, err := discovery.UnMarshal([]byte(value)); err == nil {
			// 由于语法限制，endpoint的判空得放在后面，所以NewEvent方法里要多判断一次
			if event := NewEvent(endpoint); endpoint != nil {
				event.Type = AddNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "datahandler 数据处理出错，err:%s", err.Error())
		}
	}
	delFunc := func(key, value string) {
		if endpoint, err := discovery.UnMarshal([]byte(value)); err == nil {

			if event := NewEvent(endpoint); endpoint != nil {
				event.Type = DelNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "datahandler 数据处理出错，err:%s", err.Error())
		}
	}
	err := serviceDiscovery.WatchService(config.GetServicePathForIPConf(), setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}
