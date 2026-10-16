// 是什么：Mutex 是互斥锁，让同一时刻最多一个 goroutine 进入受保护的代码段，即临界区。
// 作用：Lock 取得锁，Unlock 释放锁，防止多个任务同时读改写共享数据造成数据竞争；WaitGroup 只负责等待，不负责保护数据。
// 使用场景：保护共享计数器、缓存或内存 map；有关联约束的多个字段通常一起放在同一锁保护下。
//
// 目标：用 Mutex 保护多个 goroutine 共享的状态。
// 对照 JS：goroutine 可能在多个 CPU 核心并行执行，不能套用单线程事件循环的直觉。
// 在项目根目录运行：go run ./03_concurrency/05_mutex
// 观察：最终计数始终为 4000；锁保障更新，WaitGroup 保障主函数等待。
// 练习：理解 counter++ 为什么是“读取、计算、写回”三个步骤；再了解 RWMutex。
package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex // Mutex 零值就能用；使用后不要复制它。
	var wg sync.WaitGroup
	counter := 0
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 1000; n++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	// 所有写入都已结束，此处读不需要再加锁。
	if counter != 4000 {
		panic("计数错误")
	}
	fmt.Println(counter)
}
