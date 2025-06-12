package main

func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}

/*

➜  IM go build -o server main.go server.go user.go
➜  IM ./server

➜  IM ./server
链接建立成功

➜  IM nc 127.0.0.1 8888
[127.0.0.1:54767]127.0.0.1:54767:已上线

*/
