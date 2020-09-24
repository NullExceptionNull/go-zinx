package znet

import (
	"fmt"
	"go-zinx/utils"
	"go-zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	TaskQueue      []chan ziface.IRequest
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	m := &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerSize),
		WorkerPoolSize: utils.GlobalObject.WorkerSize,
	}
	m.StartWorkerPool()
	return m
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	router, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("None Router ,msgId = ", request.GetMsgId())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {

	if _, ok := m.Apis[msgId]; ok {
		fmt.Println("repeat api ,msgId = ", msgId)
	}

	m.Apis[msgId] = router

	fmt.Println("Add api MsgId = ", msgId, " ok !")
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerSize)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}

}

func (m *MsgHandle) StartOneWorker(i int, req chan ziface.IRequest) {
	fmt.Println("worker id = ", i, "Starting")
	for {
		select {
		case request := <-req:
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandle) SendMsgToQueue(req ziface.IRequest) {

	workerId := req.GetConnection().GetConnID() % m.WorkerPoolSize

	fmt.Println("add connection id = ", req.GetConnection().GetConnID(), "to ", workerId)

	m.TaskQueue[workerId] <- req
}
