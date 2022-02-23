package process

import (
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/client/model"
)
//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User,10)
var CurUser model.CurUser //我们在用户登录成功后，完成对CurUser初始化

//在客户端显示当前在线用户
func outputOnlineUser() {
	fmt.Println("当前用户在线列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:\t",id)
	}
}


//编写一个方法处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	user,ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {//次userid没有对应用户

		user = &message.User {
			UserId : notifyUserStatusMes.UserId,
		}

	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	
	outputOnlineUser()
}