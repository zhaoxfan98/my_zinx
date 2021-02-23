/**
  @Time    : 2021/2/21 9:33
  @Author  : zhaoxfan
*/
package ziface

//实际上是把客户端请求的链接信息和请求的数据包装到一个Request中
type IRequest interface {
	//得到当前链接
	GetConnection() IConnection
	//得到请求的消息数据
	GetData() []byte
	//得到请求的消息ID
	GetMsgID() uint32
}
