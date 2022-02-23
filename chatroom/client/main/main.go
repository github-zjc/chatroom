package main
import (
	"fmt"
	"go_code/chatroom/client/process"
)
var UserID int
var Userpwd string
var UserName string
func main(){
	var key int
	var loop bool = true
	for loop {
		fmt.Println("----------------欢迎登陆多人聊天系统：-----------------")
		fmt.Println("                   1 登录聊天系统")
		fmt.Println("                   2 注册用户")
		fmt.Println("                   3 退出系统")
		fmt.Println("请选择(1-3)")
		fmt.Println("-----------------------")
		fmt.Scanf("%d\n",&key)
		switch key {
		case 1:
			fmt.Println("登录...")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n",&UserID)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n",&Userpwd)
			//定义一个userProcess结构体，完成登录的请求
			up := &process.UserProcess{

			}
			err := up.Login(UserID,Userpwd)
			if err != nil {
				return
			}
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n",&UserID)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n",&Userpwd)
			fmt.Println("请输入用户的名字")
			fmt.Scanf("%s\n",&UserName)
			////定义一个userProcess结构体,完成注册的请求
			up := &process.UserProcess{}
			err := up.Register(UserID,Userpwd,UserName)
			if err != nil {
				return
			}

		case 3:
			fmt.Println("退出系统")
			loop = false
		default:
			fmt.Println("您的输入有误，请重新选择(1-3)")
		}

	}





}