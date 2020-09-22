package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	//直接连接远程连接 得到一个conn
	fmt.Println("client start ...")
	time.Sleep(time.Second * 3)
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error")
		return
	}
	for {
		_, err = conn.Write([]byte("hello"))

		if err != nil {
			fmt.Println("write conn error", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error", err)
			return
		}
		fmt.Printf("Server call back %s ,cnt = %d\n", buf, cnt)

		time.Sleep(time.Second)
	}

}
