package main

// conn.Read的行为特点 1. Socket中无数据 连接建立后，如果客户端未发送数据，服务端会阻塞在Socket的读操作上，这和前面提到的阻塞I/O模型的行为模式是一致的。执行该读操作的goroutine也会被挂起。Go运行时会监视该Socket，直到其有数据读事件才会重新调度该Socket对应的goroutine完成读操作
// 2. 如果Socket中有部分数据就绪，且数据量小于一次读操作所期望读出的数据长度，那么读操作将会成功读出这部分数据并返回，而不是等待期望长度数据全部读取后再返回
// 3. Socket中有足够多的数据 如果连接上有数据，且数据长度大于或等于一次Read操作所期望读出的数据长度，那么Read将会成功读出这部分数据并返回
// 4. socket 关闭 有数据关闭和无数据关闭  Read返回了错误EOF（代表连接断开）
// 5. 读操作超时
// 6. 成功写
// 7. 写阻塞 TCP通信连接两端的操作系统内核都会为该连接保留数据缓冲区，一段调用Write后，实际上数据是写入操作系统协议栈的数据缓冲区中的。TCP是全双工通信，每个方向都有独立数据缓冲区。当发送方将对方的接收缓冲区及自身的发送缓冲区都写满后，Write调用就会阻塞
// 8. 写入部分数据
// 9. 写入超时 客户端增加conn.SetWriteDeadline(time.Now().Add(time.Microsecond * 10)) 写入超时时，依旧存在数据部分写入的情况
// 10. goroutine安全的并发读写

import (
	"log"
	"net"
	"time"
)

// 测试2和3运行client的main.go传10字节以下和10字节以上
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		var buf = make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

// 测试7 write阻塞
func handleConn2(c net.Conn) {
	defer c.Close()
	time.Sleep(time.Second * 10)
	for {
		// 从连接上读取数据
		time.Sleep(5 * time.Second)
		var buf = make([]byte, 60000)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes, error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf))
	}
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// go handleConn(conn)
		go handleConn2(conn)
	}
}
