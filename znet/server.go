package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/ziface"
)

type Server struct {
	//服务器的名称
	Name string
	//tcp4
	IPVersion string
	//
	IP string
	//端口
	Port int
}

//开始网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP : %s , Port %d , is starting\n ", s.IP, s.Port)

	//开启一个 GO 协程去
	go func() {
		//解析出地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			fmt.Println("resolve tcp addr err :", err)
			return
		}

		//监听服务器地址 拿到监听器
		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		//监听成功后 就循环的等待客户端连接

		for {
			conn, err := listener.AcceptTCP()

			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//开启一个协程去拿数据
			go func() {
				for {
					buf := make([]byte, 512)

					i, err := conn.Read(buf)

					if err != nil {
						fmt.Println("recv buf err ", err)
						continue
					}
					//回显
					if _, err := conn.Write(buf[:i]); err != nil {
						fmt.Println("write back buf err ", err)
						continue
					}

				}

			}()

		}

	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
}

func (s *Server) Serve() {
	s.Start()

	//在此阻塞
	for {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
