package main

import (
	"log"
	"net"
	"time"
)

func main() {
	// 测试server下面的2,3用这段代码
	// if len(os.Args) <= 1 {
	// 	fmt.Println("usage: go run client2.go YOUR_CONTENT")
	// 	return
	// }
	// log.Println("begin dial...")
	// conn, err := net.Dial("tcp", ":8888")
	// if err != nil {
	// 	log.Println("dial error:", err)
	// 	return
	// }
	// defer conn.Close()
	// log.Println("dial ok")
	// time.Sleep(time.Second * 2)
	// data := os.Args[1]
	// conn.Write([]byte(data))
	// log.Println("send data ok")
	// time.Sleep(time.Second * 10000)
	// 测试7 write阻塞
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
	data := make([]byte, 65536)
	var total int
	for {
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}
	log.Printf("write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
}
