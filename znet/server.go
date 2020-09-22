package znet

import (
	"fmt"
	"go-zinx/ziface"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Version   string
	Port      int
}

//定义当前连接的回调方法 目前demo 就是把收到的写出去
func callBack(conn *net.TCPConn, bytes []byte, cnt int) error {
	fmt.Println("Connection handle ....")
	if _, err := conn.Write(bytes[:cnt]); err != nil {
		fmt.Println("write back error", err)
		return err
	}

	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listening at ip %s , port %d \n", s.IP, s.Port)

	go func() {
		//1:获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			fmt.Println("Resolve error", err)
			return
		}
		//2:监听地址
		tcpListen, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP error", err)
			return
		}
		fmt.Println("Start Zinx Server Succ ", s.Name, " Listening...")
		//3:阻塞获取链接
		var cid uint32
		cid = 0
		for {
			//在这里死循环 不停的获取新的连接
			conn, err := tcpListen.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error")
				continue
			}
			NewConnection(conn, cid, callBack).Start()
			cid++
		}
	}()
}

func (s *Server) Stop() {
	panic("implement me")
}

func (s *Server) Serve() {
	s.Start()
	//这里应该阻塞
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
