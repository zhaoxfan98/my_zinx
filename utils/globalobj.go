/**
  @Time    : 2021/2/21 11:15
  @Author  : zhaoxfan
*/
package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/**
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 zinx.json来配置
 */
type Globalobj struct {
	//Server
	TCPServer ziface.IServer	//当前Zinx的全局Server对象
	Host 	  string			//当前服务器主机IP
	TcpPort   int				//当前服务器主机监听端口号
	Name 	  string			//当前服务器名称

	//Zinx
	Version   string 			//当前Zinx版本号
	MaxConn   int				//当前服务器主机允许的最大链接数
	MaxPacketSize uint32		//当前zinx框架数据包的最大值

	WorkerPoolSize	uint32		//业务工作Worker池的数量
	MaxWorkerTaskLen uint32		//业务工作Worker对应负责的任务队列最大任务存储数量
}

//定义一个全局的对象
var GlobalObject *Globalobj
//读取用户的配置文件
func (g *Globalobj) Reload() {
	data, err := ioutil.ReadFile("D:\\goprojects\\src\\zinx\\mmo_game\\conf\\zinx.json")
	if err != nil{
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil{
		panic(err)
	}
}
//提供init方法，默认加载
func init(){
	//初始化GlobalObject变量，设置一些  默认值
	GlobalObject = &Globalobj{
		Name:    "ZinxServerApp",
		Version: "V0.10",
		TcpPort: 8999,
		Host:    "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize:4096,
		WorkerPoolSize: 10,			//Worker工作池的队列的个数
		MaxWorkerTaskLen: 1024,		//每个Worker对应的消息队列的任务的数量最大值
	}
	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}