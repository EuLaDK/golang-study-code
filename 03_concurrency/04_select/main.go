// 是什么：select 从多个 channel 的收发操作中，选择一个当前可以执行的分支。
// 作用：把结果、取消和超时放在同一个等待点；default 表示当前无分支就绪时立即执行，nil channel 分支永远不会就绪。
// 使用场景：一边等任务结果、一边监听取消信号或超时；避免只等待一个通道导致程序无法及时退出。
//
// 目标：用 select 等待多个 channel，理解 default、超时和关闭后的 nil 分支。
// 对照 JS：类似 Promise.race，但 select 选择的是当前可以进行的 channel 操作。
// 在项目根目录运行：go run ./03_concurrency/04_select
// 观察：A/B 顺序不固定；随后输出“没有可接收的数据”和“等待超时”。
// 练习：删除关闭分支中的 = nil，观察为什么关闭的 channel 会持续被选中。
package main

import (
	"fmt"
	"time"
)

func main() {
	a, b := make(chan string, 1), make(chan string, 1)
	a <- "A"
	b <- "B"
	close(a)
	close(b)
	for a != nil || b != nil {
		select {
		case value, ok := <-a:
			if !ok {
				a = nil // nil channel 的分支永远不会就绪，相当于禁用这个分支。
				continue
			}
			fmt.Println(value)
		case value, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			fmt.Println(value)
		}
	}

	pending := make(chan string)
	select {
	case value := <-pending:
		fmt.Println(value)
	default:
		fmt.Println("没有可接收的数据") // default 不等待；不要无休止循环造成忙等。
	}

	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()
	select {
	case value := <-pending:
		fmt.Println(value)
	case <-timer.C:
		fmt.Println("等待超时") // timer 演示超时，不用它猜测另一个 goroutine 是否完成。
	}
}
