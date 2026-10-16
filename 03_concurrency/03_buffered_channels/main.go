// 是什么：缓冲 channel 是带固定容量队列的通道，make(chan T, n) 可以暂存最多 n 个值。
// 作用：缓冲未满时发送方不必立即等接收方；缓冲满或队列空时，相应操作仍可能阻塞，即暂时等待。
// 使用场景：吸收短时间的生产与消费速度差异，或限制同时进行的任务数；缓冲只是暂存，不能替代退出和取消设计。
//
// 目标：认识缓冲区容量、队列顺序，以及关闭后仍然能读取剩余数据。
// 对照 Python：类似有容量上限的 queue.Queue；容量满时发送会阻塞。
// 在项目根目录运行：go run ./03_concurrency/03_buffered_channels
// 观察：len=2 cap=2，随后按提交顺序输出两个任务，最后 len=0。
// 练习：容量改成 1 后直接发送两个值会怎样？尝试加入并发接收者解决它。
package main

import "fmt"

func main() {
	jobs := make(chan string, 6)
	jobs <- "发"
	jobs <- "送"
	jobs <- "邮"
	jobs <- "件"
	fmt.Printf("len=%d cap=%d\n", len(jobs), cap(jobs)) // 先进先出，和队列有点类似
	close(jobs)
	for job := range jobs {
		fmt.Println(job)
	}
	fmt.Printf("len=%d\n", len(jobs))
	// len 只是当前快照，不能拿它判断并发环境下下一次发送是否一定不会阻塞。
}
