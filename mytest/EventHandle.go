package mytest

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// 先定义一下处理函数接口
type EventHandle interface {
	HandleFun(conn net.Conn) error
}

// 写一个Echo处理，就是发的啥，返回啥
type EchoHandle struct{}

func (e *EchoHandle) HandleFun(conn net.Conn) error {
	var buffer = make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			var ErrNet net.Error
			errors.As(err, &ErrNet)
			if ErrNet.Timeout() || ErrNet.Temporary() {
				fmt.Println("临时错误，不用在意")
				time.Sleep(time.Second)
				continue
			}
			fmt.Println("读取错误")
			return err
		}
		if n > 0 {
			data := buffer[:n]
			fmt.Println("输入:", data)
			_, err = conn.Write(data)
			if err != nil {
				var ErrNet net.Error
				errors.As(err, &ErrNet)
				if ErrNet.Timeout() || ErrNet.Temporary() {
					fmt.Println("临时错误，不用在意")
					time.Sleep(time.Second)
					continue
				}
				fmt.Println("写入错误")
				return err
			}
		}
	}
}

type Reactor struct {
	listener net.Listener
	handle   EventHandle
}

func NewReactor(address string, handle EventHandle) (*Reactor, error) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Reactor{
		listener: listen,
		handle:   handle,
	}, nil
}

func (r *Reactor) Start() {
	for {
		accept, err := r.listener.Accept()
		if err != nil {
			fmt.Println("错误")
			return
		}
		go r.handle.HandleFun(accept)
	}
}

func (r *Reactor) Stop() {
	r.listener.Close()
}
