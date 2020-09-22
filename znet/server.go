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
	Router    ziface.IRouter
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("add router success")
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
			NewConnection(conn, cid, s.Router).Start()
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
		Router:    nil,
	}
	return s
}
