package sdk

// 定义连接结构体
type connect struct {
	//服务器监听地址
	serverAddr string
	//发送，回复消息
	sendChan, recvChan chan *Message
}

// 创建新连接
func newConnet(serverAddr string) *connect {
	return &connect{
		serverAddr: serverAddr,
		sendChan:   make(chan *Message),
		recvChan:   make(chan *Message),
	}
}

func (c *connect) send(data *Message) {
	data.Content = "服务器回复给你的：" + data.Content
	c.recvChan <- data
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

// 关闭连接，回收资源
func (c *connect) close() {
	// 暂时没有可处理的
}
