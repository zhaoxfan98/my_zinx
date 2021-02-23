/**
  @Time    : 2021/2/20 16:29
  @Author  : zhaoxfan
*/
package main

import "zinx/znet"

/*
	基于Zinx框架开发的服务端应用程序
 */

func main()  {
	//1.创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	//2.启动Server
	s.Serve()
}
