//1.要显示登陆成功的界面  2.保持和服务器通讯（即启动协程）  3.当读取服务器发送的消息后，就会显示在界面
package process

import (
	"fmt"
	"os"
	"net"
	"go_code/chatroom/client/utils"
)

//显示登录成功后的界面
func ShowMenu() {

	fmt.Println("---------登录xxx登陆成功----------")
	fmt.Println("---------1. 显示在线用户列表-----------")
	fmt.Println("---------2. 发送消息----------")
	fmt.Println("---------3. 消息列表----------")
	fmt.Println("---------4. 退出系统----------")
	fmt.Println("请选择(1-4):")
	var key int
	fmt.Scanf("%d\n",&key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表-")
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
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
		fmt.Printf("mes %v\n",mes)
	}

}
