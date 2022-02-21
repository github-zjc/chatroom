//1.把一些常用的工具的函数，结构体  2.提供常用方法和函数
package utils
import (
	"fmt"
	"net"
	"go_code/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
)

//这里我们将这些方法关联到结构体中
type Transfer struct {
	//分析应该有哪些字段
	Conn net.Conn
	Buf  [8096]byte//这是传输时，使用的缓冲
}

func (this *Transfer) Readpkg() (mes message.Message,err error) {
	// //接收客户端发送消息的长度
	// buf := make([]byte,8096)
	n, err := this.Conn.Read(this.Buf[:4])
	if n != 4 || err != nil {
		//fmt.Println("conn.Read err=",err)
		return
	}
	fmt.Println("读到buf=",this.Buf[:4])
	
	//读取客户端消息
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	n, err = this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read err=",err)
		return
	}
	err = json.Unmarshal(this.Buf[:pkgLen],&mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=",err)
		return
	}
	return
}

func (this *Transfer) Writepkg(data []byte) (err error) {
	//7.现在data就是我们要发送的数据
	//7.1 首先把data的长度发送给服务器
	//现获取到data的长度->转成一个表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4],pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=",err)
		return
	}
	//发送data本身
	n, err = this.Conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write err=",err)
		return
	}
	return
}
