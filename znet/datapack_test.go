/**
  @Time    : 2021/2/21 15:20
  @Author  : zhaoxfan
*/
package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只是负责测试datapack拆包、封包的单元测试
func TestDataPack(t *testing.T)  {
	//模拟的服务器
	//1.创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ",err)
		return
	}

	//创建一个go承载 负责从客户端读取黏包的数据，然后进行解析
	go func() {
		//2.从客户端读取数据，拆包处理
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}
		//处理客户端请求
		go func(conn net.Conn) {
			//创建封包拆包对象dp
			dp := NewDataPack()
			for {
				// 1.第一次从conn读，把包的head读出来
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error")
					return
				}
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack err: ", err)
					return
				}
				if msgHead.GetDataLen() > 0{
					//msg是有数据的，需要进行第二次读取
					//2.第二次从conn读，根据head中的datalen再读取data内容
					//类型断言，接口转回具体的数据类型
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetDataLen())
					//根据datalen的长度再次从io流中读取
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}
					//完整的一个消息已经读取完毕
					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}
			}
		}(conn)
	}()

	//模拟客户端
	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	//创建一个封包对象 dp
	dp := NewDataPack()
	//模拟黏包过程，封装两个msg一同发送
	//封装一个msg1包
	msg1 := &Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil{
		fmt.Println("client pack msg1 error", err)
		return
	}
	//封装第二个msg2包
	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil{
		fmt.Println("client pack msg1 error", err)
		return
	}
	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)
	//向服务器端写数据
	_, _ = conn.Write(sendData1)
	//客户端阻塞
	select {}
}