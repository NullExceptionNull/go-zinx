package main

import (
	"fmt"
	"go-zinx/ziface"
	"go-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (b *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("call preHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Println("before ping error")
	}
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte(" ping ...\n"))
	if err != nil {
		fmt.Println(" ping error")
	}
}

func (b *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("call PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte(" Post Ping ...\n"))
	if err != nil {
		fmt.Println(" PostPing error")
	}
}

func main() {
	//创建服务端
	server := znet.NewServer()

	server.AddRouter(new(PingRouter))

	server.Serve()
}
