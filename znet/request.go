/**
  @Time    : 2021/2/21 9:33
  @Author  : zhaoxfan
*/
package znet

import "zinx/ziface"

type Request struct {
	//已经和客户端建立好的链接
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}
//获取请求的连接信息
func (r *Request) GetConnection() ziface.IConnection  {
	return r.conn
}
//获取请求信息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
