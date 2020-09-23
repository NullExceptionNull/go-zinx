package main

import (
	"fmt"
	"go-zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	//直接连接远程连接 得到一个conn
	fmt.Println("client start ...")
	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error")
		return
	}
	for {
		dp := znet.NewDataPack()

		binary, err := dp.Pack(znet.NewMessage(0, []byte("msg from client ...")))

		if err != nil {
			fmt.Println("client pack error", err)
			continue
		}

		if _, err := conn.Write(binary); err != nil {
			fmt.Println("client write error", err)
			continue
		}
		binaryHeadData := make([]byte, dp.GetHeadLen())

		_, err = io.ReadFull(conn, binaryHeadData)

		if err != nil {
			fmt.Println("client read buf error", err)
			break
		}
		msg, err := dp.Unpack(binaryHeadData)

		if err != nil {
			fmt.Println("client Unpack error", err)
			break
		}
		//var data []byte
		if msg.GetMsgLen() > 0 {

			msgData := msg.(*znet.Message)

			msgData.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msgData.Data); err != nil {
				fmt.Println("read msg data error ,", err)
			}
		}
		fmt.Println("recv msg callback ", string(msg.GetMsgData()))
		time.Sleep(time.Second)
	}

}
