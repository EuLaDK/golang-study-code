package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 接着退出模式的第三种
// 3、notify-and-wait模式 goroutine创建者需要主动通知那些新goroutine退出
// 3.1 通知并等待一个goroutine退出
func worker(j int) {
	time.Sleep(time.Second * (time.Duration(j)))
}

func spawn(f func(int)) chan string {
	quit := make(chan string)
	go func() {
		var job chan int // 模拟job channel

		for {
			select {
			case j := <-job:
				f(j)
			case <-quit:
				quit <- "ok"
			}
		}
	}()
	return quit
}

// 3.2 通知并等待多个goroutine退出
func spawnGroup(n int, f func(int)) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done() // 保证wg.Done在goroutine退出前被执行
			name := fmt.Sprintf("worker-%d:", i)
			for {
				j, ok := <-job
				if !ok {
					println(name, "done")
					return
				}
				// 执行这个job
				worker(j)
			}
		}(i)
	}
	go func() {
		<-quit
		close(job) // 广播给所有新goroutine // Go语言的channel有一个特性是，当使用close函数关闭channel时，所有阻塞到该channel上的goroutine都会得到通知
		wg.Wait()
		quit <- struct{}{}
	}()
	return quit
}

// 4、退出模式的应用 优雅的解决goroutine退出问题
type GracefullyShutdowner interface {
	Shutdown(waitTimeout time.Duration) error
}

type ShutdownerFunc func(time.Duration) error

func (f ShutdownerFunc) Shutdown(waitTimeout time.Duration) error {
	return f(waitTimeout)
}

// 4.1、并发退出：各个goroutine的退出先后次序对数据处理无影响、因此各个goroutine可以并发执行退出逻辑
func ConcurrentShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	// 信号通道：所有shutdown执行完成后往这个通道发信号
	c := make(chan struct{})

	// 启动一个goroutine，批量并发执行所有关闭动作 + wg等待全部完成
	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			// 每个shutdowner单独起一个goroutine并发关闭
			go func(shutdowner GracefullyShutdowner) {
				defer wg.Done()
				shutdowner.Shutdown(waitTimeout)
			}(g) // 闭包捕获循环变量，这里传参是正确写法，没有经典坑
		}
		wg.Wait()       // 阻塞，等所有组件Shutdown执行完毕
		c <- struct{}{} // 全部跑完，给外部select发完成信号
	}()
	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()

	select {
	case <-c:
		// 全部shutdown在超时前执行完了，返回nil成功
		return nil
	case <-timer.C:
		// 超时先到，返回超时错误
		return errors.New("wait timeout")
	}
}

func shutdownMaker(processTm int) func(time.Duration) error {
	return func(time.Duration) error {
		time.Sleep(time.Second * time.Duration(processTm))
		return nil
	}
}

func TestConcurrentShutdown(t *testing.T) {
	f1 := shutdownMaker(2)
	f2 := shutdownMaker(6)

	err := ConcurrentShutdown(10*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err != nil {
		t.Errorf("want nil, actual: %s", err)
		return
	}
	err = ConcurrentShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	if err == nil {
		t.Error("want timeout, actual nil")
		return
	}
}

// 4.2 串行退出，各个goroutine之前的退出是按照一定程序逐个进行的，次序错了可能导致程序的状态混乱和错误
func SequentialShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	start := time.Now()
	var left time.Duration
	timer := time.NewTimer(waitTimeout)

	for _, g := range shutdowners {
		elapsed := time.Since(start)
		left = waitTimeout - elapsed

		c := make(chan struct{})
		go func(shutdowner GracefullyShutdowner) {
			shutdowner.Shutdown(left)
			c <- struct{}{}
		}(g)

		timer.Reset(left)
		select {
		case <-c:
			// 继续执行
		case <-timer.C:
			return errors.New("wait timeout")
		}
	}
	return nil
}

func main() {
	quit := spawn(worker)
	println("spawn a worker goroutine")

	time.Sleep(5 * time.Second)

	// 通知新创建的goroutine
	println("notify the worker to exit...")
	quit <- "exit"

	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()
	select {
	case status := <-quit:
		println("worker done:", status)
	case <-timer.C:
		println("wait worker exit timeout")
	}

	quit2 := spawnGroup(5, worker)
	println("spawn a group of workers")
	time.Sleep(5 * time.Second)
	// 通知 worker goroutine组退出
	println("notify the worker group to exit...")
	quit2 <- struct{}{}

	timer2 := time.NewTimer(time.Second * 5)
	defer timer2.Stop()
	select {
	case <-timer.C:
		println("wait group workers exit timeout!")
	case <-quit2:
		println("group workers done")
	}
	// TestConcurrentShutdown()
}
