package discovery

import (
	"cocoIM/common/config"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

// 服务注册

type ServiceRegister struct {
	cli *clientv3.Client
	// 租约ID
	leaseID clientv3.LeaseID
	// 心跳检测的chan
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	value         string
	ctx           *context.Context
}

// NewServiceRegister 注册一个新的服务
// 请求上下文，主要是etcd要用  ctx *context.Context,
// 注册时候的服务名称，一般来说可以用实例名称 key string,
// 实例上提供的服务的具体信息 info *EndpointInfo,
// 租约超时时间 leaseTTL int64,
func NewServiceRegister(
	ctx *context.Context,
	key string,
	info *EndpointInfo,
	leaseTTL int64,
) (*ServiceRegister, error) {
	// 创建一个etcd连接
	client, err := clientv3.New(clientv3.Config{
		// 地址
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		log.Fatal(err)
	}
	register := &ServiceRegister{
		cli:   client,
		key:   key,
		value: info.Marshal(),
		ctx:   ctx,
	}

	// 申请租约
	if err := register.putKeyWithLease(leaseTTL); err != nil {
		return nil, err
	}
	return register, nil

}

// 创建一个带有lease的key
func (register *ServiceRegister) putKeyWithLease(leaseTTL int64) error {
	// 构造一个租约
	lease, err := register.cli.Grant(*register.ctx, leaseTTL)
	if err != nil {
		return err
	}
	// 设置一个带租约的key
	_, err = register.cli.Put(*register.ctx, register.key, register.value, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	// 设置租约的keepalive，自动续约，定期发送续约请求
	keepAliveChan, err := register.cli.KeepAlive(*register.ctx, lease.ID)
	if err != nil {
		return err
	}
	register.leaseID = lease.ID
	register.keepAliveChan = keepAliveChan
	return nil
}

// ListenLeaseRespChan 监听续租情况（监听心跳）
func (register *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepAliveResponse := range register.keepAliveChan {
		logger.CtxInfof(*register.ctx, "Lease renewal successful, the leaseID:%d,PUT key:%s,  value: %s,reps : +%v ",
			register.leaseID, register.key, register.value, leaseKeepAliveResponse)
	}
	logger.CtxInfof(*register.ctx, "ERROR! Lease renewal failed !!! the leaseID:%d,PUT key:%s,  value: %s")
}

// UpdateValue 更新
func (register *ServiceRegister) UpdateValue(val *EndpointInfo) error {
	value := val.Marshal()
	_, err := register.cli.Put(*register.ctx, register.key, value, clientv3.WithLease(register.leaseID))
	if err != nil {
		return err
	}
	register.value = value
	logger.CtxInfof(*register.ctx, "ServiceRegister.updateValue leaseID=%d Put key=%s,val=%s, success!", register.leaseID, register.key, register.value)
	return nil
}

// Close 注销服务
func (register *ServiceRegister) Close() error {
	if _, err := register.cli.Revoke(context.Background(), register.leaseID); err != nil {
		return err
	}
	logger.CtxInfof(*register.ctx, "lease close !!! leaseID:%d, Put key:%s,val:%s  success!", register.leaseID, register.key, register.value)
	return register.cli.Close()
}
