package znet

import (
	"fmt"
	"go-zinx/ziface"
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

func (c *Connection) Send(data []byte) error {
	panic("implement me")
}

//连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Read Goroutine is Running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		r := &Request{
			conn: c,
			data: buf,
		}
		go func(re ziface.IRequest) {
			c.Router.PreHandle(re)
			c.Router.Handle(re)
			c.Router.PostHandle(re)
		}(r)
	}
}
