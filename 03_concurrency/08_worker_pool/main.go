// 是什么：worker pool 是固定数量的工作 goroutine 从任务队列领取工作的一种模式。
// 作用：限制并发量，避免为大量任务一次创建过多 goroutine；给任务附上索引可以把结果恢复为输入顺序。
// 使用场景：批量处理文件或调用上游 API；真实 I/O 任务还要结合 context、超时和错误处理。
//
// 目标：限制并发数量，用 channel 分配任务、收集结果，并保持输入顺序。
// 对照 JS：Promise.all(items.map(...)) 会一次启动所有任务；固定 worker 能限制资源占用。
// 在项目根目录运行：go run ./03_concurrency/08_worker_pool
// 观察：输出 [1 4 9 16 25]；计算顺序不固定，但结果按输入位置排列。
// 练习：把平方改成字符串长度；接着加入上一课的 context 取消机制。
package main

import (
	"errors"
	"fmt"
	"sync"
)

type job struct {
	index int
	value int
}

func squareAll(values []int, workers int) ([]int, error) {
	if workers < 1 {
		return nil, errors.New("workers 必须大于 0")
	}
	jobs := make(chan job)
	results := make(chan job)
	var wg sync.WaitGroup
	for n := 0; n < workers; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				results <- job{index: task.index, value: task.value * task.value}
			}
		}()
	}
	go func() {
		for index, value := range values {
			jobs <- job{index: index, value: value}
		}
		close(jobs) // 唯一生产者负责关闭 jobs。
		wg.Wait()
		close(results) // 所有 worker 停止发送后才能关闭 results。
	}()

	output := make([]int, len(values))
	for result := range results {
		output[result.index] = result.value // 仅由当前 goroutine 写切片，不需要锁。
	}
	return output, nil
}

func main() {
	output, err := squareAll([]int{1, 2, 3, 4, 5}, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
	// 本例任务必定结束；网络/磁盘任务要设置超时，并让收发分支监听 context。
}
