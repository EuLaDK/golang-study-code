package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Go运行时和用户层有两个交互点，一个是上面说的承载信号交互channel，一个是运行时层引发的panic
// Go将信号分为两大类：一类是同步信号，另一类是异步信号

// 同步信号
// 同步信号是指那些由程序执行错误引发的信号，包括SIGBUS（总线错误/硬件异常）、SIGSEGV（段错误/无效内存引用）
// 一旦应用进程中的Go运行时收到这三个信号中的一个，意
// 味着应用极大可能出现了严重bug，无法继续执行下去，这时Go运行时
// 不会简单地将信号通过channel发送到用户层并等待用户层的异步处理，
// 而是直接将信号转换成一个运行时panic并抛出。如果用户层没有专门的
// panic恢复代码，那么Go应用将默认异常退出

// 异步信号
// 同步信号之外的信号都被Go划归为异步信号。异步信号不是由程
// 序执行错误引起的，而是由其他程序或操作系统内核发出的。异步信号
// 的默认处理行为因信号而异。SIGHUP、SIGINT和SIGTERM这三个信
// 号将导致程序直接退出；SIGQUIT、SIGILL、SIGTRAP、SIGABRT、
// SIGSTKFLT、SIGEMT和SIGSYS在导致程序退出的同时，还会将程序
// 退出时的栈状态打印出来；SIGPROF信号则是被Go运行时用于实现运
// 行时CPU性能剖析指标采集。其他信号不常用，均采用操作系统的默认
// 处理动作。对于用户层通过Notify函数捕获的信号，Go运行时则通过
// channel将信号发给用户层。

func catchAsyncSignal(c chan os.Signal) {
	for {
		s := <-c
		fmt.Println("没有收到异步信号", s)
	}
}

func triggerSyncSignal() {
	time.Sleep(3 * time.Second)
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("恢复panic：", e)
			return
		}
	}()
	var a, b = 1, 0
	fmt.Println(a / b)
}

func main() {
	// var wg sync.WaitGroup
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGFPE, syscall.SIGINT, syscall.SIGKILL)

	// wg.Add(2)
	// go func() {
	// 	catchAsyncSignal(c)
	// 	wg.Done()
	// }()

	// go func() {
	// 	triggerSyncSignal()
	// 	wg.Done()
	// }()

	// wg.Wait()

	// 如果多次调用Notify拦截某信号，但每次调用使用的channel不同，那么当应用进程收到异步信号时，Go运行时会给每个channel发送一份异步信号副本
	// c1 := make(chan os.Signal, 1)
	// c2 := make(chan os.Signal, 1)

	// signal.Notify(c1, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(c2, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	s := <-c1
	// 	fmt.Println("c1: 收到异步信号", s)
	// }()

	// s := <-c2
	// fmt.Println("c2: 收到异步信号", s)
	// time.Sleep(5 * time.Second)

	// 如何结合系统信号的使用来实现HTTP服务的优雅退出
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, Signal!\n")
	})

	srv := &http.Server{
		Addr: "localhost:8080",
	}

	// 1. 注册关闭时的清理钩子
	srv.RegisterOnShutdown(func() {
		fmt.Println("clean resources on shutdown...")
		time.Sleep(2 * time.Second)
		fmt.Println("clean resources ok")
	})

	// 2. 在独立的协程中启动服务
	serverErrors := make(chan error, 1)
	go func() {
		fmt.Println("server is running on :8080")
		// ListenAndServe 始终返回非 nil 的 error
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 3. 监听操作系统信号
	// 进阶：使用 signal.NotifyContext 可以极大简化 context 控制 (Go 1.16+)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 4. 等待中断信号 或 服务自身的致命错误（如端口被占用）
	select {
	case err := <-serverErrors:
		fmt.Printf("web server start failed: %v\n", err)
		return
	case sig := <-quit:
		fmt.Printf("\nreceived signal: %v, shutting down...\n", sig)
	}

	// 5. 设置优雅退出的最大超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 6. 调用 Shutdown，阻塞等待所有连接和 RegisterOnShutdown 结束
	// 如果超时，Shutdown 会立即返回 context.DeadlineExceeded
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("web server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("web server gracefully stopped")
	}

	fmt.Println("program exit ok")
}
