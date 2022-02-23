//处理和用户相关的请求   2.登录   3.注册   4.注销   5.用户列表管理
package process
import (
	"fmt"
	"net"
	"go_code/chatroom/common/message"
	"encoding/json"
	"go_code/chatroom/server/utils"
	"go_code/chatroom/server/model"
)
//用户上线
type UserProcess struct {
	//字段？
	Conn net.Conn
	//添加一个字段，表示该Conn是哪个用户
	UserId int
}

func (this *UserProcess) NotifyOthersOnlineUser(userid int) {

	for id, up := range SuserMgr.onlineUsers {
		if id == userid {
			continue
		}
		up.NotifyMeOnline(userid)
	}

}

func (this *UserProcess) NotifyMeOnline(userid int) {

	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userid
	notifyUserStatusMes.Status = message.UserOnline

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("NotifyMeOnline notifyUserStatusMes json.Marshal err=",err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("NotifyMeOnline mes json.Marshal err=",err)
		return
	}

	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.Writepkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline writepkg err=",err)
		return
	}

}

//用户下线
func (this *UserProcess) NotifyOthersOfflineUser(userid int) {

	for id, up := range SuserMgr.onlineUsers {
		if id == userid {
			continue
		}
		up.NotifyMeOffline(userid)
	}

}

func (this *UserProcess) NotifyMeOffline(userid int) {

	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userid
	notifyUserStatusMes.Status = message.UserOffline

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("NotifyMeOffline notifyUserStatusMes json.Marshal err=",err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("NotifyMeOffline mes json.Marshal err=",err)
		return
	}

	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.Writepkg(data)
	if err != nil {
		fmt.Println("NotifyMeOffline writepkg err=",err)
		return
	}

}



//用户登陆处理
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//将mes.Data的值取出
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data),&loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err= ",err)
		return
	}
	//声明一个resMes
	var resMes message.Message
	//对resMes完成赋值
	resMes.Type = message.LoginResMesType
	//声明一个LoginResMes，并完成赋值
	var loginResMes message.LoginResMes
	//判断是否合法
	//现在我们需要到redis数据库上去完成验证
	//1.使用model》MyUserDao,到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserID,loginMes.Userpwd)
	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD  {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}

	} else {
		loginResMes.Code = 200
		//添加在线用户
		this.UserId = loginMes.UserID
		SuserMgr.AddOnlineUser(this)

		//通知其他用户我上线了
		this.NotifyOthersOnlineUser(loginMes.UserID)

		//将当前在线用户的id 放入到loginResMes。UserId
		//遍历userMgr.OnlineUsers
		for id,_ := range SuserMgr.onlineUsers {
			loginResMes.UserId = append(loginResMes.UserId,id)
		}		
		
		fmt.Println(user,"登陆成功")
	}

	//将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}

	//因为使用了分层模式（mvc）我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.Writepkg(data)

	return
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data),&registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err= ",err)
		return
	}
	//声明一个resMes
	var resMes message.Message
	//对resMes完成赋值
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes
	//我们需要到redis数据库去完成注册
	//1.使用model.MyUserDao到redis去验证
	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}

	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}

	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}

	//因为使用了分层模式（mvc）我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.Writepkg(data)
	return
}

//用户退出处理
func (this *UserProcess) ServerProcessOff(mes *message.Message) (err error) {
	//将mes.Data的值取出
	var notifyUserStatusMes message.NotifyUserStatusMes
	err = json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err= ",err)
		return
	}
	//声明一个resMes
	var resMes message.Message
	//对resMes完成赋值
	resMes.Type = message.LoginResMesType
	//声明一个LoginResMes，并完成赋值
	var loginResMes message.LoginResMes

	loginResMes.Code = 200
	//删除退出用户
	this.UserId = notifyUserStatusMes.UserId
	SuserMgr.DelOnlineUser(this.UserId)

	//通知其他用户我下线了
	this.NotifyOthersOfflineUser(notifyUserStatusMes.UserId)

	//将当前在线用户的id 放入loginResMes.UserId中
	//遍历userMgr.OnlineUsers
	for id,_ := range SuserMgr.onlineUsers {
		loginResMes.UserId = append(loginResMes.UserId,id)
	}		
		

	//将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=",err)
		return
	}

	//因为使用了分层模式（mvc）我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.Writepkg(data)

	return
}