package main

import "net"

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
	this.server.BoardCast(this, "已上线")
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
	} else {
		this.server.BoardCast(this, msg)

	}
}

//监听当前user channel的方法，一旦有消息就直接发送给对方客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
