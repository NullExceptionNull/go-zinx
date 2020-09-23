package main

import (
	"fmt"
	"go-zinx/ziface"
	"go-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

type HelloRouter struct {
	znet.BaseRouter
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call Handle...")

	err := request.GetConnection().Send([]byte(" ping ...\n"), 1)
	if err != nil {
		fmt.Println(" ping error")
	}
}

func (b *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("call Handle...")

	err := request.GetConnection().Send([]byte(" hello ...\n"), 1)
	if err != nil {
		fmt.Println(" hello error")
	}
}

func main() {
	//创建服务端
	server := znet.NewServer()

	server.AddRouter(0, new(PingRouter))
	server.AddRouter(1, new(HelloRouter))

	server.Serve()
}
