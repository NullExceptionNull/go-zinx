package ziface

import "net"

type IConnection interface {
	//启动链接
	Start()
	//停止连接
	Stop()
	//获取当前连接的socket conn
	GetTCPConnection() *net.TCPConn
	//获取连接ID
	GetConnID() uint32
	//获取远程客户端的状态
	RemoteAddr() net.Addr
	//发送数据
	Send(data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error
