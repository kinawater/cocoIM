package cocoIM

import (
	"context"
	"net"
	"time"
)

type Server interface {
	// SetAcceptor 设置接收者，主要用于接收Start方法启动后监听到的连接，然后交给处理器去处理握手的相关问题
	SetAcceptor(Acceptor)
	// SetStateListener 设置状态监听器，主要是上报连接状态，让上层做出对应的处理
	SetStateListener(StateListener)
	// 设置读取超时
	SetReadWait(time.Duration)
	// 连接管理channel
	SetChannelMap(ChannleMap)
	// 设置消息监听
	SetMessageListener(MessageListener)

	//Push 发送指定消息到指定的channel里
	Push(channelID string, msg []byte) error

	Start() error
	Shutdown(ctx context.Context) error
}

type Client interface {
	ID() string
	Name() string
	Connect(string) error
	SetDialer(Dialer)
	Send([]byte) error
	Read() (Frame, error)
	Close()
}

// Acceptor 连接接收器，用于保存收到的连接
type Acceptor interface {
	//Accept 接收器
	//Conn 被接收的连接
	//time 超时时间
	Accept(Conn, time.Duration) (string, error)
}

type Conn interface {
	// Conn 网络连接
	net.Conn
	ReadFrame()
	WriteFrame()
	Flush() error
}

//Frame websocket 一帧的封装
type Frame interface {
	SetOpCode(code OpCode)
	GetOpCode() OpCode
	SetPayload([]byte)
	GetPayload() []byte
}

// StateListener 状态监听器
type StateListener interface {
	// Disconnect 断开连接
	Disconnect(string) error
}

type ChannleMap interface {
	//Add 添加一个连接到MAP
	//TODO:Add(channel Channel)
	Add()
	//Remove 移除
	Remove(id string)
	//Get 获取
	//TODO: Get(id string)(Channel,bool)
	Get(id string)
	//All 全部
	//TODO: All()[]Channel
	All()
}

//MessageListener 消息监听
type MessageListener interface {
	//Receive 接收
	Receive(Agent, []byte)
}

//Agent 发送方
type Agent interface {
	ID() string
	Push([]byte) error
}

//Dialer 拨号器
type Dialer interface {
	// DialAndHandshake 握手
	DialAndHandshake(DialerContext) (net.Conn, error)
}

//DialerContext 拨号连接上下文
type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}

// OpCode 定义一个ws的operation code  和 ws的官方文档保持一致
type OpCode byte

const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)
