/**
  @Time    : 2021/2/20 16:06
  @Author  : zhaoxfan
*/
package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP地址
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle
	//该Server的连接管理器
	ConnMgr ziface.IConnManager
	//该Server创建链接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server销毁链接之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)

}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router succ! ")
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP : %s, Port: %d is starting\n",
		s.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		//0. 开启消息队列及Worker工作池
		s.MsgHandle.StartWorkerPool()

		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//2.监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen:", s.IPVersion, " err:", err)
			return
		}
		fmt.Println("Start zinx server succ, ", s.Name, " succ Listenning...")
		var cid uint32
		cid = 0
		//3.阻塞的等待客户端连接 处理客户端连接业务
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}
			//设置最大链接个数的判断  如果超过了最大链接数量  则关闭此新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn{
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("=========>Too Many Connection, MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid ++
			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	//将一些服务的资源、状态、连接 进行停止或回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	//启动Server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

//得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	//先初始化全局配置文件
	utils.GlobalObject.Reload()
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
	return s
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func (ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func (ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
