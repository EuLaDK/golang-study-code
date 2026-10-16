package main

import (
	"fmt"
	"sync"
	"time"
)

// channel 是为 goroutine 间的通信设计的，因此，在大多数情况下，对 channel 的读写操作都会被分别放在不同的 goroutine 中执行
// 声明一个channel变量
var ch chan int // 默认值为nil

func produce(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i + 1
		time.Sleep(time.Second)
	}
	close(ch)
}

func consume(ch <-chan int) {
	for n := range ch {
		println(n)
	}
}

// func worker() {
// 	select {
// 	case <-c:
// 		// ....
// 	case <-time.After(30 * time.Second):
// 		return
// 	}
// }

// func worker2() {
// 	heartbeat := time.NewTicker(30 * time.Second)
// 	defer heartbeat.Stop()
// 	for {
// 		select {
// 		case <-c:
// 			// 执行某些操作
// 		}
// 	case <- heaheartbeat.C:
// 		// 执行某些心跳事件
// 	}
// }

func main() {
	// channel 类型变量赋初始值的唯一方法是使用 make 这个预定义函数
	// 无缓冲channel，元素类型为T
	ch1 := make(chan int)
	// 带有capacity 参数(就是5)就是带缓冲channel
	// 带缓存 channel 在实现中包含一个缓冲区。因此，在缓冲区未满时进行发送操作或在缓冲区非空时进行接收操作时异步的（即发送或接收不需要阻塞等待）
	// 	也就是说，对于一个带缓冲 channel，在缓冲区未满的情况下，执行发送操作的 goroutine 不会被阻塞；在缓冲区有数据时，执行接收操作的 goroutine 也不会被阻塞。
	// 但是，当缓冲区已满时，试图向其发送数据的 goroutine 将会被阻塞；当缓冲区为空时，尝试从中接收数据的 goroutine 也会被阻塞
	// 	ch2 := make(chan int, 1)
	// n := <-ch2 // 由于此时 ch2 的缓冲区中没有数据，因此执行接收操作将导致 goroutine 阻塞
	// ch3 := make(chan int, 1)
	// ch3 <- 17 // 向 ch3 发送整型型值 17
	// ch3 <- 27 // 由于此时 ch3 的缓冲区已满，因此再次发送数据将导致 goroutine 阻塞
	ch2 := make(chan int, 5)
	fmt.Println(ch1)
	// ch1 <- 13 // 将整型值 13 发送到无缓冲 channel 整型变量 ch1 中
	// n := <-ch1 // 从无缓冲 channel 类型变量 ch1 中接收一个整型值并存储到整型变量 n  这里只是做示范，因为接收和发送都在一个地方会触发死锁

	ch2 <- 17  // 将整型值 17 发送到带缓存 channel 类型变量 ch2 中
	m := <-ch2 // 从带缓冲 channel 类型变量 ch2
	fmt.Println(m)

	ch3 := make(chan int)
	// 对无缓冲 channel 类型的发送与接收操作，一定要放在两个不同的 goroutine 中进行，否则会导致程序进入死锁状态
	// 针对ch1的那个死锁问题可以这么改
	go func() {
		ch3 <- 13
	}()
	n3 := <-ch3
	fmt.Println(n3)

	// 使用操作符 <-,还可以声明只发送（chan<-）和只接收(<-chan)的channel类型
	ch4 := make(chan<- int, 1) // 只发送channel类型
	ch5 := make(<-chan int, 1) // 只接收 channel 类型
	fmt.Println(ch4, ch5, "--------------------")

	ch6 := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		produce(ch6)
		wg.Done()
	}()
	go func() {
		consume(ch6)
		wg.Done()
	}()
	wg.Wait()

	// 关闭channel后的返回值
	// n := <-ch // 对应类型的0值 假如这里的ch是int 那么 n就是0
	// m ,ok := <-ch  // 0，false
	// for v := range ch {...} // for range 语句结束

	// select 原语
	// 通过select，我们可以同时在多个channel上进行发送/接收操作
	// 	select {
	// case x := <-ch1: // 从 ch1 接收数据
	// …
	// case y, ok := <-ch2: // 从 ch2 接收数据，并根据 ok 值判断 ch2 是否已关闭
	// …
	// case ch3 <- z: // 将 z 值发送到 ch3 中
	// …
	// default: // 当上面 case 语句中的 chl 通信均无法实施时，执行该默认分支
	// }
	// 当select语句中没有default分支切所有case语句都处于阻塞式，select语句本身也会被阻塞，直到某个case语句对应的channel变成可发送或者可接收

	// 几种常见的channel和select结合使用的方法
	// 1、利用default分支避免阻塞 当其它case都阻塞时，select会执行default
	// 2、实现超时机制 上面的worker函数 在应用带有超时机制的select语句是，要特别注意timer资源的释放，尤其是在大量创建timer的情况下
	// 3、实现心跳机制 上面的worker2函数
}
