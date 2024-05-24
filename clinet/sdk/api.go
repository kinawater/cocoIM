package sdk

// 定义 聊天 结构体
type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

// 定义 消息 结构体
type Message struct {
	Type       string
	Name       string
	FromUserID string
	ToUserID   string
	Content    string
	Session    string
}
