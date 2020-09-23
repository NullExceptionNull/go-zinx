package znet

import (
	"errors"
	"fmt"
	"go-zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	conn     *net.TCPConn
	ConnID   uint32
	IsClosed bool
	ExitChan chan bool
	Router   ziface.IRouter
}

//初始化链接模块
func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		conn:     conn,
		ConnID:   connId,
		Router:   router,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) Start() {
	fmt.Println("Conn Start () ... ConnID = ", c.ConnID)
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop () ... ConnID = ", c.ConnID)

	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
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
	_, err = c.GetTCPConnection().Write(binary)
	if err != nil {
		fmt.Println("write error msg id =", msgId)
		return err
	}
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
		go func(re ziface.IRequest) {
			c.Router.PreHandle(re)
			c.Router.Handle(re)
			c.Router.PostHandle(re)
		}(r)
	}
}
