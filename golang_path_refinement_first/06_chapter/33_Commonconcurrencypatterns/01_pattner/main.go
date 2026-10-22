package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Go 针对 CSP模型提供了三种并发原语
// 1、goroutine：对应CSP模型中的P，封装了数据的处理逻辑，是Go运行时调度的基本执行单元
// 2、channel：对应CSP模型中的输入/输出原语，用于goroutine之间的通信和同步
// 3、select：用于应对多路输入/输出，可以让goroutine同时协调处理多个channel操作

// func spawn(f func()) chan T {
// 	c:=make(chan T)
// 	go func() {
// 		// 使用channel变量c(通过闭包方式)与调用spawn的goroutine通信
// 		...
// 		f()
// 		...
// 	}()
// 	return c
// }

func worker(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	interval, ok := args[0].(int)
	if !ok {
		return
	}
	time.Sleep(time.Second * (time.Duration(interval)))
}

func spawn(f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	go func() {
		f(args...)
		c <- struct{}{}
	}()
	return c
}

var OK = errors.New("ok")

func worker2(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	interval, ok := args[0].(int)
	if !ok {
		return errors.New("invalid interval arg")
	}
	time.Sleep(time.Second * (time.Duration(interval)))
	return OK
}

func spawn2(f func(args ...interface{}) error, args ...interface{}) chan error {
	c := make(chan error)
	go func() {
		c <- f(args...)
	}()
	return c
}

func spawnGroup(n int, f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			name := fmt.Sprintf("worker-%d:", i)
			f(args...)
			fmt.Println(name, "done")
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		c <- struct{}{}
	}()
	return c
}

func main() {
	// 创建模式
	// 1、使用go关键字+函数/方法创建goroutine
	// c := srv.newConn(rw)
	// go c.serve(connCtx)

	// 2、使用channel变量c(通过闭包方式)与调用spawn的goroutine通信
	// c := spawn(func(){})
	// 使用channel变量c与新创建的goroutine通信

	// 退出模式
	// 1、分离模式
	// 这类goroutine在启动后即与其创建者彻底分离，其生命周期与其执行的主函数相关，函数返回即goroutine退出。这类goroutine有创建两个用途
	// 1.1、一次性任务，新创建的goroutine用来执行一个简单的任务，执行后退出
	// 1.2、常驻后台执行一些特定任务，如监视(monitor)、观察(watch)等。其实现通常采用for{...}或for{select{...}}代码段形式，并多以定时器（timer）或事件（event）驱动执行
	// Go为每个goroutine调度模型中P内置的GC goroutine就是这种类型
	// $GOROOT/src/runtime/mgc.go
	// func gcBgMarkStartWorkers() {
	// // 每个P都有一个运行在后台的用于标记的G
	// for _, p := range allp {
	// if p.gcBgMarkWorker == 0 {
	// go gcBgMarkWorker(p) // 为每个P创建一个goroutine，以运行
	// gcBgMarkWorker
	// notetsleepg(&work.bgMarkReady, -1)
	// noteclear(&work.bgMarkReady)
	// }
	// }
	// }
	// func gcBgMarkWorker(_p_ *p) {
	// gp := getg()
	// ...
	// for { // 常驻后台处理GC事宜
	// ...
	// }
	// }

	// 2、join模式 goroutine的创建者需要等待新goroutine结束
	// 2.1、等待一个goroutine退出
	done := spawn(worker, 5)
	fmt.Println("spawn a worker goroutine")
	<-done
	fmt.Println("worker done")
	// 2.2、获取goroutine的退出状态
	done2 := spawn2(worker2, 5)
	fmt.Println("spawn worker1")
	err := <-done2
	fmt.Println("worker1 done:", err)
	done2 = spawn2(worker2)
	fmt.Println("spawn worker2")
	err = <-done2
	fmt.Println("worker2 done:", err)

	// 2.3、等待多个goroutine退出
	// 在有些场景中，goroutine的创建者可能会创建不止一个goroutine，并且需要等待全部新goroutine退出。可以通过Go语言提供的sync.WaitGroup实现等待多个goroutine退出的模式
	doneWait := spawnGroup(5, worker, 3)
	fmt.Println("spawn a group of workers")
	<-doneWait
	fmt.Println("group workers done")
	// 2.4 支持超时机制的等待
	doneWait2 := spawnGroup(5, worker, 30)
	fmt.Println("spawn a group of workers")
	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()
	select {
	case <-timer.C:
		fmt.Println("wait group workers exit timeout!")
	case <-doneWait2:
		fmt.Println("group workers done")
	}
}
