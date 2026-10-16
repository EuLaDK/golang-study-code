// 是什么：channel 是带元素类型的通信通道，用于在 goroutine 之间发送和接收值。
// 作用：无缓冲 channel 的发送与接收需要双方接手才能完成；chan<- 只允许发送，<-chan 只允许接收。
//
//	close 表示发送方不会再发数据，接收方可以用 range 读到通道耗尽；关闭不会替你取消生产任务。
//
// 使用场景：把结果交给另一项任务、连接生产者和消费者，或通知任务结束；关闭通常由明确的发送方负责。
//
// 目标：用无缓冲 channel 传递数据，了解发送、接收、关闭和单向类型。
// 对照 JS：channel 更像带同步能力的队列，并不是事件广播；一次发送只被接收一次。
// 在项目根目录运行：go run ./03_concurrency/02_channels
// 观察：依次输出 1、4、9；关闭后接收得到零值和 ok=false。
// 练习：生产 5 个平方；解释为什么接收者不负责关闭这个 channel。
// 向已关闭的通道发送数据会导致程序直接崩溃 在Go的底层机制中，向已经关闭的通道发送数据(ch <- data), 会直接引发 panic
// 只有发送者才知道什么时候发完了，比方我去送快递，只有我才知道快递送完没有，收快递的人是不知道的
// 如果有多个发送者，会触发“重复关闭”的崩溃
package main

import "fmt"

func produce(output chan<- int) {
	defer close(output) // 发送方知道何时不再有数据，因此由它关闭。
	for n := 1; n <= 3; n++ {
		output <- n * n // 无缓冲：等到接收方接手数据才继续。
	}
}

func consume(input <-chan int) {
	for value := range input { // 直到 channel 关闭并且数据耗尽。
		fmt.Println(value)
	}
}

func main() {
	values := make(chan int)
	go produce(values)
	consume(values)
	value, ok := <-values
	fmt.Printf("关闭后：value=%d, ok=%t\n", value, ok)
	// 向已关闭的 channel 发送或再次关闭它会 panic；nil channel 的收发永远阻塞。
}
