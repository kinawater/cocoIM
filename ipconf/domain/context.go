package domain

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type IpConfContext struct {
	// 请求上下文
	Ctx *context.Context
	// 用户请求上下文
	UserRequestContext *app.RequestContext
	// 客户端请求上下文
	// 其实当前用不上这个上下文，只是预留一下
	ClientCtx *ClientCtx
}
type ClientCtx struct {
	// 从用户请求的上下文里其实已经能拿到ip地址了，暂时先放着，用到再说
	IP string `json:"ip"`
}

func BuildIpConfContext(ctx *context.Context, requestContext *app.RequestContext) *IpConfContext {
	return &IpConfContext{
		Ctx:                ctx,
		UserRequestContext: requestContext,
		ClientCtx:          &ClientCtx{},
	}
}
