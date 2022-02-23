package process

import (
	"fmt"
)

//因为UserMgr 实例在服务器端有且只有一个
//因为在很多的地方，都会使用到，所以，我们定义将其定义为全局变量

var (
	SuserMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化
func init() {
	SuserMgr = &UserMgr{
		onlineUsers : make(map[int]*UserProcess,1024),
	}
}

//完成对onlineUsers的添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {

	this.onlineUsers[up.UserId] = up
}

//del
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers,userId)
}

//返回当前所以在线用户
func (this *UserMgr) GetOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

//根据id返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess,err error) {

	up, ok := this.onlineUsers[userId]
	if !ok {
		//说明你当前查找的用户不在线
		err = fmt.Errorf("用户%d 不存在",userId)
		return
	}
	return
}