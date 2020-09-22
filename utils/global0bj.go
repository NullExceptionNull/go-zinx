package utils

import (
	"encoding/json"
	"fmt"
	"go-zinx/ziface"
	"io/ioutil"
)

type GlobalObj struct {
	TcpServer      ziface.IServer
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Name           string `json:"name"`
	Version        string `json:"version"`
	MaxPackageSize int32  `json:"max_package_size"`
	MaxConnSize    int32  `json:"max_conn_size"`
}

var GlobalObject *GlobalObj

func (obj *GlobalObj) reload() {
	file, err := ioutil.ReadFile("./conf/zinx.json")
	if err != nil {
		fmt.Println("read file error")
	}
	err = json.Unmarshal(file, GlobalObject)
	if err != nil {
		fmt.Println("Unmarshal file error", err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Host:           "localhost",
		Port:           8999,
		Name:           "Zinx",
		MaxConnSize:    1000,
		MaxPackageSize: 4096,
		Version:        "0.3",
	}
	GlobalObject.reload()
}
