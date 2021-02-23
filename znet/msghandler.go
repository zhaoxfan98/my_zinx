/**
  @Time    : 2021/2/22 0:40
  @Author  : zhaoxfan
*/
package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

//消息处理模块的实现
type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32] ziface.IRouter
	//业务工作worker池的数量
	WorkerPoolSize	uint32
	//负责Worker取任务的消息队列
	TaskQueue	[]chan ziface.IRequest
}
//初始化 创建MsgHandle方法
func NewMsgHandle() *MsgHandle{
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest){
	//1. 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok{
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
		return
	}
	//2. 根据MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

}
// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter){
	//1. 判断的当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok{
		panic("repeated api, msgId = " + strconv.Itoa(int(msgId)))
	}
	//2.添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}
//启动一个work工作池（开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池）
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启Worker 每个Worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++{
		//一个Worker被启动
		//1.给当前Worker对应的任务队列开辟空间 第0个Worker用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2.启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}
//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest){
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
			//有消息则取出队列的Request，并执行绑定的业务方法
			case request := <-taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}
//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle)SendMsgToTaskQueue(request ziface.IRequest){
	//1.根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//2.将消息发送给对应的Worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}