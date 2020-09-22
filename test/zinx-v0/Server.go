package main

import "go-zinx/znet"

func main() {
	//创建服务端
	server := znet.NewServer("V0.1")

	server.Serve()
}
