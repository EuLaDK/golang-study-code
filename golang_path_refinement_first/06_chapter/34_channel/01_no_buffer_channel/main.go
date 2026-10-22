package main

import (
	"fmt"
	"sync"
	"time"
)

// 创建无缓冲channel T是类型 无缓冲channel的接收和发送操作是同步的
// c := make(chan T)
// 即对于同一个无缓冲channel，只有在对其进行接收操作的goroutine和对其进行发送操作的goroutine都存在的情况下，通信才能进行，否则单方面的操作会让对应的goroutine陷入阻塞状态
// 无缓冲 channel（make(chan T)）没有缓冲区，发送和接收必须同时发生才行：
// 发送方执行 ch <- v 时，如果没有接收方在等，它会阻塞，直到有人来接收。
// 接收方执行 <-ch 时，如果没有发送方在等，它也会阻塞，直到有人来发送。
// 两边必须“碰面”，值才能从发送方转移到接收方。2. 为什么有两条“一定发生在……之前”Go 的 happens-before 规则对无缓冲 channel 是这样规定的（官方 memory model）：
// 发送动作 happens-before 对应的接收动作的完成
// （send happens before the completion of the corresponding receive）
// 接收动作 happens-before 对应的发送动作的完成
// （receive happens before the completion of the corresponding send）
// 用更直白的话说：
// 发送方把值交给 channel 之后，接收方才能把这个值取走并继续往下执行。
// → 所以“发送”一定发生在“接收完成”之前。
// 接收方把值取走之后，发送方才能认为发送成功并继续往下执行。
// → 所以“接收”一定发生在“发送完成”之前。
// 两边都在等待对方，谁也不能先“完成”自己的操作。

// 发送方
// ch <- 42          // ① 到达这里后阻塞，等待接收方
// fmt.Println("send done")  // ③ 只有接收方取走值后，才执行到这里

// // 接收方
// x := <-ch         // ② 到达这里后阻塞，等待发送方
// fmt.Println("recv done")  // ④ 只有发送方把值交过来后，才执行到这里

var cc = make(chan int)
var a string

func f() {
	a = "hello world"
	<-cc
}

// 一对一通知信号 无缓冲channel常被用于在两个goroutine之间一对一地传递通知信号
type signal struct{}

func workerOneToOne() {
	println("worker is working...")
	time.Sleep(3 * time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		println("工人开始工作...")
		f()
		c <- signal(struct{}{})
	}()
	return c
}

// 一对多信号通知
func workerOneToMany(i int) {
	fmt.Printf("工人 %d: 正在工作...\n", i)
	time.Sleep(3 * time.Second)
	fmt.Printf("工人 %d: 工作完毕\n", i)
}

func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1) // 没启动一个工人，计数器加1
		go func(i int) {
			<-groupSignal // 等开始信号
			fmt.Printf("工人 %d: 开始工作...\n", i)
			f(i)      // 开始干活
			wg.Done() // 干完后计数器 -1
		}(i + 1)
	}
	go func() {
		wg.Wait()               // 阻塞，直到计数器归零
		c <- signal(struct{}{}) // 所有工人干完了，才发完成信号
	}()
	return c
}

// 无缓冲channel用于替代锁机制
type counter struct {
	c chan int
	i int
}

var cter counter

func InitCounter() {
	cter = counter{
		c: make(chan int),
	}
	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	fmt.Println("计数器初始化完成")
}

func Increase() int {
	return <-cter.c
}

func init() {
	InitCounter()
}

func main() {
	x1 := 1
	x2 := 2
	c1 := make(chan int)    // 创建一个无缓冲(unbuffered)的int类型的channel
	c2 := make(chan int, 5) // 创建一个带缓冲的int类型的channel
	go func() {             // 无缓冲的必须要有接受不然会报错
		c1 <- x1
		c1 <- x2
		close(c1)
	}()
	c2 <- x2
	fmt.Println(<-c1) // 从channel c中接收一个值
	y1 := <-c1        // 从channel c接收一个值并赋值给变量y1
	y2, ok := <-c2    // 从 channel c中接收一个值，若channel关闭了，ok将被置为false
	fmt.Println(y1, y2, ok)
	// for i := range c1 {...} 将 for range 与channel结合使用
	close(c2) // 关闭channel c

	// 	c := make(chan chan int) // 创建一个无缓冲的chan int类型的channel
	// func stream(ctx context.Context, out chan<- Value) error // 将只发送(sendonly) channel作为函数参数
	// func spawn(...) <-chan T // 将只接收(receive-only)类型channel作为返回值

	// 	select {
	// case x := <-c1: // 从channel c1接收数据
	// ...
	// case y, ok := <-c2: // 从channel c2接收数据，并根据ok值判断c2是否已
	// 经关闭
	// ...
	// case c3 <- z: // 将z值发送到channel c3中
	// ...
	// default: // 当上面case中的channel通信均无法实施时，执行该默认分支
	// }

	go f()
	cc <- 5
	println(a)
	fmt.Println("一对一通信")
	println("开始工作...")
	one := spawn(workerOneToOne)
	<-one
	fmt.Println("工人工作完毕")

	fmt.Println("一对多通信")
	fmt.Println("开始多工人...")
	groupSignal := make(chan signal)
	many := spawnGroup(workerOneToMany, 5, groupSignal)
	time.Sleep(5 * time.Second)
	fmt.Println("很多工人开始工作...")
	close(groupSignal)
	<-many
	fmt.Println("这个组的工人工作完毕！")

	fmt.Println("这是用channel替代锁实现的计数器")
	for i := 0; i < 10; i++ {
		go func(i int) {
			v := Increase()
			fmt.Printf("goroutine-%d:current counter value is %d\n", i, v)
		}(i)
	}
	time.Sleep(5 * time.Second)
}
