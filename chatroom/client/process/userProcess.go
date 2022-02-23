package process
import (
	"fmt"
	"net"
	"go_code/chatroom/common/message"
	"encoding/json"
	"encoding/binary"
	"go_code/chatroom/client/utils"
	"os"
)

type UserProcess struct {
	//暂时不需要字段
}

//写一个函数完成登录
func (this *UserProcess) Login(UserID int,Userpwd string) (err error) {
	
	// fmt.Println("你输入的userid=%d pwd=%s",UserID,Userpwd)
	// return nil
	//1.与服务器建立连接
	conn, err := net.Dial("tcp","localhost:8889")
	defer conn.Close()
	if err != nil {
		fmt.Println("net.Dial err=",err)
	}
	//2.通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserID = UserID
	loginMes.Userpwd = Userpwd
	//4.对login结构体进行序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}
	//5.把data赋值给mes.Data字段
	mes.Data = string(data)
	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}
	//7.现在data就是我们要发送的数据
	//7.1 首先把data的长度发送给服务器
	//现获取到data的长度->转成一个表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4],pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=",err)
		return
	}

	fmt.Printf("客户端发送消息的长度=%d 内容=%s\n",len(data),string(data))
	// return

	//发送客户端信息data
	n, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write err=",err)
		return
	}

	//接收服务器的消息并进行处理
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn : conn,
	}
	mes, err = tf.Readpkg()
	if err != nil {
		fmt.Println("readpkg err=",err)
		return
	}
	fmt.Println("接收到的消息mes=",mes)
	//将mes.Data部分反序列化验证用户是否登录成功
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data),&loginResMes)

	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = UserID
		CurUser.UserStatus = message.UserOnline
		//fmt.Println("登录成功")
		//显示在线用户
		fmt.Println("当前在线用户有:")
		for _, v := range loginResMes.UserId {
			
			if UserID == v {
				continue
			}
			fmt.Println("用户id:\t",v)

			//完成客户端的onlineUsers 完成初始化
			user := &message.User {
				UserId : v,
				UserStatus : message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		//起一个协程用来读取给客户端发来的消息
		go serverProcessMes(conn)

		//循环显示登陆成功后的界面
		for {
			ShowMenu(conn,UserID)
		}


	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

func (this *UserProcess) Register(UserID int,Userpwd string,UserName string) (err error) {
	//与服务器建立连接
	conn, err := net.Dial("tcp","localhost:8889")
	defer conn.Close()
	if err != nil {
		fmt.Println("net.Dial err=",err)
	}

	//2.通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3.创建一个LoginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = UserID
	registerMes.User.UserPwd = Userpwd
	registerMes.User.UserName = UserName

	//4.对registerMes结构体进行序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}

	//5.把data赋值给mes.Data字段
	mes.Data = string(data)
	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}
	
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn : conn,
	}

	//发送data给服务器端
	err = tf.Writepkg(data)
	if err != nil {
		fmt.Println("注册发送信息出错了 err=",err)
	}

	mes, err = tf.Readpkg()
	if err != nil {
		fmt.Println("readpkg err=",err)
		return
	}
	fmt.Println("接收到的消息mes=",mes)
	
	//将mes.Data部分反序列化验证用户是否注册成功
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data),&registerResMes)
	if registerResMes.Code == 200 {
		
		fmt.Println("注册成功，你可以重新登录")
		os.Exit(0)

	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}

// func exitMes(Conn net.Conn,userid int) {

// 	tf := &utils.Transfer{
// 		Conn : Conn,
// 	}

// 	var mes message.Message
// 	mes.Type = message.NotifyUserStatusMesType

// 	var notifyUserStatusMes message.NotifyUserStatusMes
// 	notifyUserStatusMes.UserId = userid
// 	notifyUserStatusMes.Status = message.UserOffline
		
// 	data, err := json.Marshal(notifyUserStatusMes)
// 	if err != nil {
// 		fmt.Println("ExitMes(Conn net.Conn,userid int) notifyUserStatusMes json.Marshal err=",err)
// 		return
// 	}

// 	mes.Data = string(data)

// 	data, err = json.Marshal(mes)
// 	if err != nil {
// 		fmt.Println("ExitMes(Conn net.Conn,userid int) mes json.Marshal err=",err)
// 		return
// 	}

// 	err = tf.Writepkg(data)
// 	if err != nil {
// 		if err != nil {
// 			fmt.Println("ExitMes(Conn net.Conn,userid int) Weitepkg err=",err)
// 			return
// 		}	
// 	}

// }

