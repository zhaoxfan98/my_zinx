/**
  @Time    : 2021/2/20 17:08
  @Author  : zhaoxfan
*/
package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

//模拟客户端
func main() {
	fmt.Println("client0 start...")
	time.Sleep(1 * time.Second)
	//1. 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		//发送封包的Msg消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinx client Test0 Message")))
		if err != nil {
			fmt.Println("Pack error, ",err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error, ",err)
			return
		}
		//服务器应该给我们回复一个message数据，MsgID:1 ping..ping..ping

		//1.先读取流中的head部分，得到ID和dataLen
		binaryHaed := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHaed); err != nil{
			fmt.Println("reaf head error, ", err)
			break
		}
		//将二进制head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHaed)
		if err != nil{
			fmt.Println("client unpack  MsgHead err:", err)
			break
		}
		if msgHead.GetDataLen() > 0{
			//2. 再根据DataLen进行第二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data error, ", err)
				return
			}
			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		//cpu阻塞
		time.Sleep(1*time.Second)
	}

}
