package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//就好像GIN 的中间件
	AddRouter(msgId uint32, router IRouter)
}
