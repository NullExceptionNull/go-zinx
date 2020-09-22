package ziface

//将连接与请求数据绑定
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}
