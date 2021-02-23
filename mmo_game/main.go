/**
  @Time    : 2021/2/23 9:04
  @Author  : zhaoxfan
*/
package main

import (
	"fmt"
	"zinx/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

//当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn ziface.IConnection){
	//创建一个玩家
	player := core.NewPlayer(conn)
	//同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SyncPid()
	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()
	//将当前新上线玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

func main()  {
	//创建服务器句柄
	s := znet.NewServer("MMO Game Zinx")
	//连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnecionAdd)
	//注册路由业务

	//启动服务
	s.Serve()

}
