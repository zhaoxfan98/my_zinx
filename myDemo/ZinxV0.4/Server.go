/**
  @Time    : 2021/2/20 16:29
  @Author  : zhaoxfan
*/
package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于Zinx框架开发的服务端应用程序
 */

//ping test 自定义路由
type PingRouter struct {
	//先基础BaseRouter
	znet.BaseRouter
}
//Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("call router PreHandler")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back ping err")
	}
}
//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping ping err")
	}
}
//Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("call router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
	if err != nil {
		fmt.Println("call back ping ping err")
	}
}
func main()  {
	//1.创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.4]")
	//2. 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//3.启动Server
	s.Serve()
}
