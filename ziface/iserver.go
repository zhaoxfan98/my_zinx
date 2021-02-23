/**
  @Time    : 2021/2/20 16:06
  @Author  : zhaoxfan
*/
package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId uint32, router IRouter)
	//得到连接管理
	GetConnMgr() IConnManager
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func (connection IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func (connection IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(connection IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(connection IConnection)
}
