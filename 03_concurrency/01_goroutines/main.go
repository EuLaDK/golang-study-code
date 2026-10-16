// 是什么：goroutine 是 Go 运行时调度的轻量执行单元，go f() 让函数与当前流程并发执行。
// 作用：让多个任务独立推进；WaitGroup 记录未完成的任务数，用 Add 登记、Done 完成、Wait 等待全部结束。
// 使用场景：同时请求多个上游服务或处理独立任务；必须明确如何等待结束和处理共享数据，不能假设输出顺序固定。
//
// 目标：用 goroutine 并发执行任务，用 WaitGroup 等待它们完成。
// 对照 JS：go f() 不是 Promise；main 返回时，其他 goroutine 也会结束。
// 在项目根目录运行：go run ./03_concurrency/01_goroutines
// 观察：三个任务的输出顺序不固定，但“全部完成”一定最后输出。
// 练习：把任务数改为 5；不要用 time.Sleep 猜测任务何时完成。
package main

import (
	"fmt"
	"sync"
	"time" // 1. 引入 time 包
)

func main() {
	var wg sync.WaitGroup
	startTime := time.Now() // 记录总开始时间

	for id := 1; id <= 5; id++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			// 2. 模拟每个任务不同的耗时
			// 比如：任务 1 睡 500ms，任务 5 睡 100ms
			sleepTime := time.Duration(6-taskID) * 100 * time.Millisecond
			time.Sleep(sleepTime)

			// 3. 打印任务完成时的时刻和经过的时间
			fmt.Printf("任务 %d 完成 | 耗时: %v | 完成时刻: %s\n",
				taskID,
				time.Since(startTime),             // 距程序启动过去了多久
				time.Now().Format("15:04:05.000"), // 精确到毫秒的时刻
			)
		}(id)
	}

	wg.Wait()
	// 4. 观察全部完成的总耗时
	fmt.Printf("全部完成，总耗时: %v\n", time.Since(startTime))
}
