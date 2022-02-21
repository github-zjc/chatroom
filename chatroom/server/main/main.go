//1.监听  2.等待客户的连接  3.初始化工作
package main 
import (
	"fmt"
	"net"
	"time"
	"go_code/chatroom/server/model"
)

func process1(conn net.Conn) {
	//设置延时关闭
	defer conn.Close()

	//这里调用总控，创建一个总控
	processor := &Processor{
		Conn : conn,
	}

	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误=err",err)
		return
	}
}

func init() {
	//当服务器启动时，我们就去初始化我们的redis的连接池
	initPool("localhost:6379", 16, 0, 300 * time.Second)
	initUserDao()
}

//这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool 本身就是一个全局的变量
	//这里需要注意一个初始化顺序问题
	//initPool, 在 initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

func main(){

	//提示信息
	fmt.Println("新的结构~~~服务器在8889端口进行监听")
	listen, err := net.Listen("tcp","0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=",err)
		return
	}
	//一旦监听成功，就等待客户来连接服务器
	fmt.Println("等待客户端来连接服务器")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=",err)
		}
		//一旦建立连接成功，就启动一个协程和客户端保持通讯
		go process1(conn)
	}
}