package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 0:退出 1:公聊 2:私聊 3:更新用户名
}

// NewClient 创建一个新的客户端对象
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

func (client *Client) DealResponse() {
	// 处理服务器响应的消息
	// 这里可以实现接收服务器消息的逻辑
	buf := make([]byte, 4096)
	for {
		n, err := client.conn.Read(buf)
		if n == 0 {
			fmt.Println("服务器已关闭")
			return
		}
		if err != nil {
			fmt.Println("读取服务器消息失败:", err)
			return
		}
		// _ = string(buf[:n])
		msg := string(buf[:n])
		fmt.Println(msg)
	}
	// 将消息输出到标准输出,永久阻塞监听
	// io.Copy(os.Stdout, client.conn) // 将服务器的消息输出到标准输出
}

// UpdateName 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>> 更新用户名")
	fmt.Print("请输入用户名:")
	fmt.Scanln(&client.Name)

	// 发送给服务器
	// 这里可以实现发送更新用户名的逻辑
	// 比如发送一个特定格式的消息给服务器
	msg := fmt.Sprintf("rename|%s", client.Name)
	_, err := client.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("更新用户名失败:", err)
		return false
	}
	fmt.Println("用户名更新成功")
	return true
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	fmt.Println(">>>>> 公聊模式")
	fmt.Println("请输入消息内容,exit退出")
	var msg string
	for {
		fmt.Scanln(&msg)
		if msg == "exit" {
			fmt.Println("退出公聊模式")
			break
		}
		// 发送消息到服务器
		_, err := client.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("发送消息失败:", err)
			continue
		}
		fmt.Println("消息已发送:", msg)
		msg = "" // 清空消息内容，准备下一次输入
		fmt.Println("请输入消息内容,exit退出")
	}
}

// SelectUsers 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	// 发送查询在线用户的消息到服务器
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("查询在线用户失败:", err)
		return
	}
}

// PrivateChat 私聊模式
func (client *Client) PrivateChat() {
	// 首先查询在线用户
	client.SelectUsers()

	// 私聊模式的逻辑
	// fmt.Println(">>>>> 私聊模式")

	var targetUser string
	var msg string
	fmt.Println("请输入私聊对象的用户名:")
	fmt.Scanln(&targetUser)
	fmt.Println("请输入私聊内容:")
	for {
		// 读取用户输入的私聊内容
		fmt.Scanln(&msg)
		if msg == "exit" {
			fmt.Println("退出私聊模式")
			break
		}
		// 发送私聊消息到服务器
		// 这里可以实现发送私聊消息的逻辑
		privateMsg := fmt.Sprintf("to|%s|%s", targetUser, msg)
		_, err := client.conn.Write([]byte(privateMsg + "\n"))
		if err != nil {
			fmt.Println("发送私聊消息失败:", err)
			continue
		}
		// fmt.Println("私聊消息已发送给", targetUser, ":", msg)
		msg = "" // 清空消息内容，准备下一次输入
	}
}

// menu 显示菜单并获取用户选择
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入0~3")
		return false
	}
}

func (client *Client) Run() {
loop:
	for client.menu() {
		switch client.flag {
		case 0:
			fmt.Println(">>>>> 客户端退出")
			break loop //退出for循环
		case 1:
			// fmt.Println(">>>>> 公聊模式")
			client.PublicChat()
			break //退出switch,但不退出for循环
		case 2:
			// fmt.Println(">>>>> 私聊模式")
			client.PrivateChat()
			break //退出switch,但不退出for循环
		case 3:
			// fmt.Println(">>>>> 更新用户名")
			client.UpdateName()
			break //退出switch,但不退出for循环
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析
	flag.Parse()

	// client := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 链接服务器失败")
		return
	}
	fmt.Println(">>>>> 链接服务器成功")

	// 启动客户端处理服务器响应的 goroutine
	go client.DealResponse()

	// 启动客户端业务
	client.Run()

}

// go build -o client client.go
// ./client -ip 127.0.0.1 -port 8888
