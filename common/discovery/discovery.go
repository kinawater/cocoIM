package discovery

import (
	"cocoIM/common/config"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdClient "go.etcd.io/etcd/client/v3"
	"sync"
)

// 服务发现
type ServiceDiscovery struct {
	lock sync.Mutex
	// 客户端连接
	cli *etcdClient.Client
	ctx *context.Context
}

// NewServiceDiscovery 新建 “服务发现” 服务
func NewServiceDiscovery(ctx *context.Context) *ServiceDiscovery {
	// 创建一个etcd 连接
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		logger.Fatal(err)
	}
	return &ServiceDiscovery{
		cli: client,
		ctx: ctx,
	}
}

// 初始化服务发现列表和监视
func (s *ServiceDiscovery) WatchService(prefix string, set, del func(key, value string)) error {
	resp, err := s.cli.Get(*s.ctx, prefix, etcdClient.WithPrefix())
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		set(string(ev.Key), string(ev.Value))
	}
	// 监视前缀，修改变更的server
	s.watcher(prefix, resp.Header.Revision+1, set, del)
	return nil
}

// watcher 监视前缀
func (s *ServiceDiscovery) watcher(prefix string, rev int64, set, del func(key, value string)) {
	// 监听前缀和增加版本号
	rch := s.cli.Watch(*s.ctx, prefix, etcdClient.WithPrefix(), etcdClient.WithRev(rev))
	logger.CtxInfof(*s.ctx, "watching prefix:%s now...", prefix)
	for response := range rch {
		for _, event := range response.Events {
			switch event.Type {
			case mvccpb.PUT:
				set(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				del(string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}
}

// Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
