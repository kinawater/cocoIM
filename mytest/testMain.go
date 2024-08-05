package mytest

import "fmt"

func RunMain() {
	handle := &EchoHandle{}
	reactor, err := NewReactor("localhost:12360", handle)
	if err != nil {
		fmt.Println("创建反应器出错", err)
		panic(err)
	}
	fmt.Println("反应器启动")
	reactor.Start()
}
