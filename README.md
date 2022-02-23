# chatroom
海量数据通讯系统

#本项目在vsCode上编译

1:使用前需要连接redis数据库，并开启redis-cli

2:分别对server包和client包下的main进行go build
例如：
go bulid -o server.exe ./go_code/chatroom/server/main
go bulid -o client.exe ./go_code/chatroom/client/main

3：go build 后，到项目的src目录上输入cmd打开终端，终端数量根据自己需求来

4：打开终端后，输入server.exe 就开启了服务器，输入client.exe就开启了客户端

5：服务器和客户端都开启后，根据自己需求操作
