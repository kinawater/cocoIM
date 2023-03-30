package errorInfo

// GetMsg 获取错误代码对应的消息
func GetMsg(code int) string {
	msg, ok := ErrorMsgCN[code]
	if ok {
		return msg
	}
	return ErrorMsgCN[ERROR]
}
