package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建一个server接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

//监听message广播消息channel的goroutine，一旦有消息就广播给所有在线用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		//将message发送给全部用户
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//广播当前消息
func (this *Server) BoardCast(user *User, msg string) {
	SendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- SendMsg
}
func (this *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	//fmt.Println("连接建立成功")
	//用户上线，将用户加入到onlinemap中
	user := NewUser(conn, this)
	user.Online()
	//监听用户是否是活跃的channel
	isLive := make(chan bool)

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户消息（取出结尾\r\n）
			msg := string(buf[:n-1])

			//用户针对msg消息进行处理
			user.DoMessage(msg)

			//用户任意消息，算是活跃用户
			isLive <- true
		}
	}()

	//当前handle阻塞
	for {
		select {
		case <-isLive:
		//当前用户是活跃的，应该重置定时器
		//不做任何处理，为了激活select，重置定时器

		case <-time.After(time.Second * 1800):
			//已经超时
			//将当前的user强制剔除
			user.SendMsg("你已经超时，被强制下线.")

			//销毁用户资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出当前handler
			return //runtime.Goexit()
		}

	}
}

//启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close listen scoket
	defer listener.Close()

	//启动监听message的goroutine
	go this.ListenMessager()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("连接失败:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}

}
