package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建一个用户的api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前 user channel消息的goroutine
	go user.ListenMessage()
	return user
}

//用户上线业务
func (this *User) Online() {
	//用户上线，将用户加入到onlinemap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BoardCast(this, "已上线\r\n使用【rename|名称】可以修改名字\r\n使用【who】可以查看在线用户\r\n")
}

//用户下线业务
func (this *User) Offline() {
	//用户下线，将用户从onlinemap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户下线消息
	this.server.BoardCast(this, "已下线")
}

//给当前user对应客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

//处理用户消息功能
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前用户都有哪些
		this.server.mapLock.Lock()

		for _, user := range this.server.OnlineMap {
			OnlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线..."
			this.SendMsg(OnlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]
		//this.server.OnlineMap[this.Name]=newName
		//判断当前用户是否已存在
		if _, ok := this.server.OnlineMap[newName]; ok {
			//fmt.Printf("当前用户%s已存在",newName)
			this.SendMsg("当前用户" + newName + "已存在")
			return
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName

			this.SendMsg("您已修改当前用户名：" + newName + "\r\n")
		}

	} else {
		this.server.BoardCast(this, msg)

	}
}

//监听当前user channel的方法，一旦有消息就直接发送给对方客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\r\n"))
	}
}
