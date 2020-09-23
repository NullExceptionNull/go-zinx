package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {

	listen, err := net.Listen("tcp", "127.0.0.1:7777")

	if err != nil {
		fmt.Println("server listen err : ", err)
		return
	}

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("server Accept err : ", err)
				return
			}

			//处理客户端的请求
			go func(conn2 net.Conn) {
				//定义一个拆包的对象
				dp := NewDataPack()

				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					head, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("unpack head error")
						return
					}
					if head.GetMsgLen() > 0 {
						msg := head.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("unpack data error")
							return
						}
						fmt.Println("----> Recv MsgId : ", msg.Id, ", DataLen : ", msg.GetMsgLen())
					}
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:7777")

	if err != nil {
		fmt.Println("client dial error : ", err)
		return
	}

	dp := NewDataPack()

	//
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}

	pack1, err := dp.Pack(msg1)

	if err != nil {
		fmt.Println("client pack msg1 error", err)
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', '0', '!', '.'},
	}

	pack2, err := dp.Pack(msg2)

	if err != nil {
		fmt.Println("client pack msg2 error", err)
	}

	pack1 = append(pack1, pack2...)

	conn.Write(pack1)

	select {}
}
