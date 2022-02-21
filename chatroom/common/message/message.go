package message

const (
	LoginMesType = "LoginMes"
	LoginResMesType = "loginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
)

type Message struct {
	Type string	`json:"type"`//消息类型
	Data string `json:"data"`//消息的内容
}

type LoginMes struct {
	UserID int	`json:"userid"`//用户id
	Userpwd string `json:"userpwd"`//用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code int	`json:"code"`//返回码500表示用户未注册，200表示登录成功
	Error string `json:"error"`//返回错误信息
}

type RegisterMes struct {
	User User `json:"user"` //类型就是User结构体
}

type RegisterResMes struct {
	Code int	`json:"code"`//返回码400表示用户以占有，200表示注册成功
	Error string `json:"error"`//返回错误信息
}