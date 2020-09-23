package main

import (
	"fmt"
	"go-zinx/ziface"
	"go-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call Handle...")

	err := request.GetConnection().Send([]byte(" ping ...\n"), 1)
	if err != nil {
		fmt.Println(" ping error")
	}
}

func main() {
	//创建服务端
	server := znet.NewServer()

	server.AddRouter(new(PingRouter))

	server.Serve()
}
