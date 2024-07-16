package ipconf

import (
	"cocoIM/ipconf/domain"
)

func TopEndports(eds []*domain.Endport, top int) []*domain.Endport {
	if top <= 0 {
		return nil
	}
	if len(eds) < top {
		return eds
	}
	return eds[:top]
}

// 格式化返回值
func PackRes(eds []*domain.Endport) Response {
	return Response{
		Message: "ok",
		Code:    0,
		Data:    eds,
	}

}
