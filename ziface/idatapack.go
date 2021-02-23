/**
  @Time    : 2021/2/21 14:47
  @Author  : zhaoxfan
*/
package ziface

//封包数据和拆包数据  直接面向TCP连接中的数据流，为传输数据添加头部信息，用于处理TCP黏包问题
type IDataPack interface {
	GetHeadLen() uint32		//获取包头长度方法
	Pack(msg IMessage)([]byte, error)	//封装方法
	Unpack([]byte)(IMessage, error)		//拆包方法
}
