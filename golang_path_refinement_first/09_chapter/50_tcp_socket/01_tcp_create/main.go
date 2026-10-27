package main

import (
	"log"
	"net"
)

// TCP Socket网络编程模型 分为阻塞和非阻塞
// 阻塞和非阻塞是以内核是否等数据全部就绪才返回（给发起系统调用的应用线程）来区分的。如果内核一直等到全部数据就绪才返回，则这种行为模式称为阻塞；如果内核查看数据就绪状态后，即便没有就绪也立即返回错误（给发起系统调用的应用线程），则这种行为模式称为非阻塞
// 1、阻塞 I/O 模型
// 2、非阻塞 I/O 模型
// 3、I/O 多路复用模型

// tcp的建立 客户端Go语言使用Dial或DialTimeout函数发起连接请求
// Dial在调用后将一直阻塞，直到连接建立成功或失败
// conn, err := net.Dial("tcp","taobao.com:80")
// if err != nil {...} 连接建立成功，可以进行读写操作

// DialTimeout 是带有超时机制的Dial
// conn,err:=net.DialTimeout("tcp","localhost:8080",2*time.Second) if err!=nil{...} 连接客户端，可以进行读写操作

func main() {
	// 对于客户端而言，建立连接时可以回遇到以下几种情况
	// 1、网络不可达或对方服务未启动 如果传给Dial的服务端地址是网络不可达的，或者服务地址中端口对应服务并没有启动，端口未被监听，则Dial几乎会返回错误
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	// 2、对方服务的listen backlog队列满了 瞬间有大量客户端尝试与服务端建立连接，服务端可能会出现listen backlog队列满了，接收连接accept不及时的情况，这将导致客户端的Dial调用阻塞
	// 3、若网络延迟较大，Dial将阻塞并超时
}
