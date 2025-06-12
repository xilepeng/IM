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

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的 channel
	Message chan string
}

// 创建一个 Server 的接口
func NewServer(ip string, port int) *Server {

	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听 Message 广播消息的 goroutine, 一旦有消息就发送给全部在线的 User
func (this *Server) LestenMessager() {
	for {
		msg := <-this.Message
		// 将 msg 发送给全部的在线 User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Hander(conn net.Conn) {
	fmt.Println("链接建立成功")

	// 用户上线，将用户加入到 onlinemap 中
	user := NewUser(conn, this)
	user.Online()
	// this.mapLock.Lock()
	// this.OnlineMap[user.Name] = user
	// this.mapLock.Unlock()

	// 广播当前用户上线消息
	// this.BroadCast(user, "已上线")

	// 监听用户是否活跃的 channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 100000)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				// this.BroadCast(user, "下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户的消息，去除 '\n'
			msg := string(buf[:n-1])

			// 将得到的消息广播
			// this.BroadCast(user, msg)
			// 用户针对 message 进行消息处理
			user.DoMessage(msg)
			// 用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()

	// 当前 handler 阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，应该重置定时器
			// 不做任何事情，为了激活 select,更新下面的定时器
		case <-time.After(time.Second * 100000):
			// 已经超时，
			// 将当前的 user 强制关闭
			user.sendMsg("你已超时，连接已关闭")
			// 销毁用户的资源
			close(user.C)
			// 关闭连接
			conn.Close()
			// 退出当前 handler
			return
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	//	socket listen
	//net.Listen("tcp", "127.0.0.1:8888")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Lesten is err", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// 启动监听 Message 的 grotinue
	go this.LestenMessager()

	for {
		//	accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// do hander

		go this.Hander(conn)
	}

}
