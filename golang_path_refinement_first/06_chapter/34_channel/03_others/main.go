package main

import (
	"fmt"
	"time"
)

func main() {
	// nil channel 的妙用
	c1, c2 := make(chan int), make(chan int)
	go func() {
		time.Sleep(time.Second * 5)
		c1 <- 5
		close(c1)
	}()
	go func() {
		time.Sleep(time.Second * 7)
		c2 <- 7
		close(c2)
	}()

	for {
		select {
		case x, ok := <-c1:
			if !ok {
				c1 = nil // 这里显示的修改为nil可以防止一直输出0
			} else {
				fmt.Println(x)
			}
		case x, ok := <-c2:
			if !ok {
				c2 = nil
			} else {
				fmt.Println(x)
			}
		}
		if c1 == nil && c2 == nil {
			break
		}
	}
	fmt.Println("program end")

	// 与select结合使用的一些惯用法
	// 1、利用default分支避免阻塞 select语句的default分支的语义是在其他分支均因通信未就绪而无法被选择的时候执行，这就为default分支赋予了一种“避免阻塞”的特性
	// $GOROOT/src/time/sleep.go
	// func sendTime(c interface{}, seq uintptr) {
	// // 无阻塞地向c发送当前时间
	// // ...
	// select {
	// case c.(chan Time) <- Now():
	// default:
	// }
	// }

	// 2、实现超时机制
	// 	func worker() {
	// select {
	// case <-c:
	// // ...
	// case <-time.After(30 *time.Second): // 开启一个30s的定时器，如果这期间没有任何一个case进入就进入这个返回
	// return
	// }
	// }

	// 3、实现心跳机制
	// 	func worker() {
	// heartbeat := time.NewTicker(30 * time.Second)
	// defer heartbeat.Stop()
	// for {
	// select {
	// case <-c:
	// // ... 处理业务逻辑
	// case <- heartbeat.C:
	// //... 处理心跳
	// }
	// }
	// }
}
