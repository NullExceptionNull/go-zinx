package znet

import (
	"errors"
	"fmt"
	"go-zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	conn       *net.TCPConn
	ConnID     uint32
	IsClosed   bool
	ExitChan   chan bool
	msgHandler ziface.IMsgHandle
	msgChan    chan []byte
}

//初始化链接模块
func NewConnection(conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		conn:       conn,
		ConnID:     connId,
		msgHandler: msgHandler,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}
	return c
}

//写消息 专门发送客户端消息的
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is Running...]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("Send data error ..", err)
			}
		case <-c.ExitChan:
			return
		}
	}

}

func (c *Connection) Start() {
	fmt.Println("Conn Start () ... ConnID = ", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop () ... ConnID = ", c.ConnID)

	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	c.ExitChan <- true
	defer c.conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Send(data []byte, msgId uint32) error {
	if c.IsClosed {
		return errors.New("Connection is closed when sending data...")
	}
	dp := NewDataPack()

	binary, err := dp.Pack(NewMessage(msgId, data))

	if err != nil {
		fmt.Println("Pack error msg id =", msgId)
		return err
	}
	c.msgChan <- binary
	return nil

}

//连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Read Goroutine is Running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())

		_, err := io.ReadFull(c.GetTCPConnection(), headData)

		if err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		message, err := dp.Unpack(headData)

		if err != nil {
			fmt.Println("unpack msg head error", err)
			break
		}
		var data []byte
		if message.GetMsgLen() > 0 {
			data = make([]byte, message.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error : ", err)
				break
			}
		}
		message.SetMsgData(data)
		r := &Request{
			conn: c,
			msg:  message,
		}

		c.msgHandler.SendMsgToQueue(r)
		//c.msgHandler.
	}
}
