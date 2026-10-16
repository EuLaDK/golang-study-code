package main

import (
	"errors"
	"fmt"
	"time"
)

func spawn(f func() error) <-chan error {
	c := make(chan error)
	go func() {
		c <- f()
	}()
	return c
}

func main() {
	// 创建一个goroutine 使用go关键字创建
	go fmt.Println("I am a goroutine")
	// 基于匿名函数 / 闭包创建 goroutine，并通过channel 传递结果
	var c = make(chan int)
	go func(a, b int) {
		c <- a + b
	}(3, 4)
	// c:= srv.newConn(rw)
	// go c.serve(connCtx)
	// 当 goroutine 执行的函数或方法返回时，该 goroutine 自动退出
	// 即使 goroutine

	ch := spawn(func() error {
		time.Sleep(2 * time.Second)
		return errors.New("timeout")
	})
	fmt.Println(<-ch)
}
