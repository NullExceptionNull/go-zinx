package znet

import (
	"fmt"
	"go-zinx/ziface"
	"io"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Version   string
	Port      int
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
		for {
			//在这里死循环 不停的获取新的连接
			conn, err := tcpListen.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error")
				continue
			}
			go func() {
				//需要不停的从连接中获取数据
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil && err != io.EOF {
						fmt.Println("Recv buf err", err)
						continue
					}
					fmt.Printf("Recv client buf %s ,cnt %d\n", buf, cnt)
					bytes := append(buf[:cnt], []byte(" from server")...)
					if _, err := conn.Write(bytes); err != nil {
						fmt.Println("write buf err", err)
						continue
					}
				}
			}()
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
