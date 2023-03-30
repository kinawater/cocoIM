package errorInfo

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

const (
	SUCCESS_CN        = "ok"
	ERROR_CN          = "fail"
	INVALID_PARAMS_CN = "请求参数错误"
)

var ErrorMsgCN = map[int]string{
	SUCCESS:        SUCCESS_CN,
	ERROR:          ERROR_CN,
	INVALID_PARAMS: INVALID_PARAMS_CN,
}
