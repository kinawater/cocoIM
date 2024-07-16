package ipconf

import (
	"cocoIM/ipconf/domain"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 规定返回类型
type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

func GetIpInfoList(ctx context.Context, requestContext *app.RequestContext) {
	// 如果没有数据
	defer func() {
		if err := recover(); err != nil {
			requestContext.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	// 构建请求ipconflist
	ipRequestStruct := domain.BuildIpConfContext(&ctx, requestContext)
	// 请求dispatch得到已经排序的iplist
	candidateList := domain.Dispatch(ipRequestStruct)
	// 返回前5的数据,因为客户端可能还会根据一些客户本地的策略进行筛选，比如服务器的响应ping值
	ipRequestStruct.UserRequestContext.JSON(consts.StatusOK, PackRes(TopEndports(candidateList, 5)))
}
