/**
  @Time    : 2021/2/23 8:11
  @Author  : zhaoxfan
*/
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx/myDemo/protobufDemo/pb"
)

func main()  {
	//定义一个Person结构对象
	person := &pb.Person{
		Name: "zhao",
		Age: 22,
		Emails: []string{"zhaoxfan98@gmail.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "18810617098",
				Type: pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "13731998947",
				Type: pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "18822255664",
				Type: pb.PhoneType_WORK,
			},
		},
	}
	//编码
	//将person对象  就是将protobuf的message进行序列化，得到一个二进制文件
	data, err := proto.Marshal(person)
	//data就是需要进行网络传输的数据 对端需要按照Message Person格式进行解析
	if err != nil{
		fmt.Println("marshal err:", err)
	}
	//解码
	newPerson := &pb.Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil{
		fmt.Println("unmarshal err:", err)
	}
	fmt.Println("原数据：",person)
	fmt.Println("解码后：", newPerson)
}

