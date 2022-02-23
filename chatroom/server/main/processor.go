//1.根据客户端的请求，调用对应的处理器，完成相应的任务
package main
import (
	"fmt"
	"net"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/utils"
	"go_code/chatroom/server/process"
	"io"
)

//先创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

func (this *Processor) serverProcessMes(mes *message.Message) (err error) {	

	fmt.Println("mes=",mes)
	switch mes.Type {
		case message.LoginMesType :
			//处理登录
			//创建一个UserProccss实例
			up := &process.UserProcess{
				Conn :this.Conn,
			}
			err = up.ServerProcessLogin(mes)
		case message.RegisterMesType :
			//处理注册
			up := &process.UserProcess{
				Conn :this.Conn,
			}
			err = up.ServerProcessRegister(mes)
		case message.NotifyUserStatusMesType :
			//处理用户退出
			up := &process.UserProcess{
				Conn :this.Conn,
			}
			up.ServerProcessOff(mes)
		case message.SmsMesType :
			//创建一个SmsProcess实例完成转发群聊消息.
			smsProcess := &process.SmsProcess{}
			smsProcess.SendGroupMes(mes)

		default :
			fmt.Println("mes.Type 类型有误")
	}
	return

}

func (this *Processor) process2() (err error) {
		//循环的读取客户端发送的信息
		fmt.Println("读取客户端发送的信息")
		for {
			//这里我们将读取数据包，直接封装成readpkg()，返回mes ，err
			//创建一个Transfer实例来完成读包任务
			tf := &utils.Transfer{
				Conn : this.Conn,
			}
			mes, err := tf.Readpkg()
			if err != nil {
				if err == io.EOF {
					fmt.Println("客户端退出，服务器也退出")
					return err
				} else {
					fmt.Println("readpkg err=",err)
					return err
				}
			}
			fmt.Println("接收到的消息mes=",mes)
			err = this.serverProcessMes(&mes)
			if err != nil {
				return err
			}
	
		}
}

