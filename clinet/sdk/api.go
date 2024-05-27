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

// 创建对话
func MakeNewChat(serverAddr, nick, userID, seessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: seessionID,
		conn:      newConnet(serverAddr),
	}
}
func (chat *Chat) Send(msg *Message) {
	chat.conn.send(msg)
}
func (chat *Chat) Recv() <-chan *Message {
	return chat.conn.recv()
}
func (chat *Chat) Close() {
	chat.conn.close()
}
