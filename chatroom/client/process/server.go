//1.要显示登陆成功的界面  2.保持和服务器通讯（即启动协程）  3.当读取服务器发送的消息后，就会显示在界面
package process

import (
	"fmt"
	"net"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"encoding/json"
	"os"
)

//显示登录成功后的界面
func ShowMenu(conn net.Conn,userid int) {

	fmt.Println("---------聊天系统登陆成功----------")
	fmt.Println("---------1. 显示在线用户列表-----------")
	fmt.Println("---------2. 发送消息----------")
	fmt.Println("---------3. 消息列表----------")
	fmt.Println("---------4. 退出系统----------")
	fmt.Println("请选择(1-4):")
	var key int
	fmt.Scanf("%d\n",&key)
	var content string

	//因为，我们总会使用到SmsProcess实例，因此我们将其定义在swtich外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key) 
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表-")
		outputOnlineUser()
	case 2:
		fmt.Println("发送消息")
		fmt.Println("你想对大家说的什么:)")
			fmt.Scanf("%s\n", &content)
			smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		exitMes(conn,userid)
		fmt.Println("你选择退出了系统...")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确")
	}
}

func serverProcessMes(Conn net.Conn) {

	tf := &utils.Transfer{
		Conn : Conn,
	}
	fmt.Println("客户端正在等待读取服务器发送的消息")
	for {
		mes, err := tf.Readpkg()
		if err != nil {
			fmt.Println("serverProcessMes err=",err)
			return
		}
		//如果读到了消息，又是下一步处理逻辑
		//fmt.Printf("mes %v\n",mes)
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			//1.取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes)
			if notifyUserStatusMes.Status == message.UserOnline {
			   //2.把这个用户的信息状态保存到客户段的map[int]User中
			   updateUserStatus(&notifyUserStatusMes)
			}
			if notifyUserStatusMes.Status == message.UserOffline {
				delete(onlineUsers,notifyUserStatusMes.UserId)
			}
		case message.SmsMesType:
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")

		}
	}

}

func exitMes(Conn net.Conn,userid int) {

	tf := &utils.Transfer{
		Conn : Conn,
	}

	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userid
	notifyUserStatusMes.Status = message.UserOffline
		
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("ExitMes(Conn net.Conn,userid int) notifyUserStatusMes json.Marshal err=",err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("ExitMes(Conn net.Conn,userid int) mes json.Marshal err=",err)
		return
	}

	err = tf.Writepkg(data)
	if err != nil {
		if err != nil {
			fmt.Println("ExitMes(Conn net.Conn,userid int) Weitepkg err=",err)
			return
		}	
	}

}




