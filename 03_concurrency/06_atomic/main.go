// 是什么：原子操作把针对单个状态的一次操作作为不可被拆开的整体完成，sync/atomic 提供原子计数等能力。
// 作用：Add 安全累加，Load 读取，CompareAndSwap 只有在当前值等于预期值时才更新。
// 使用场景：并发统计请求数或更新简单状态；多个字段必须一起变化时，单独做多个原子操作不能保证整体一致。
// 它做的锁也能做，在单个状态上比 Mutex 更轻量，但不能替代 Mutex。
// 目标：用 sync/atomic 对单个计数器做原子更新。
// 对照 JS：概念接近 Atomics；Go 的原子类型不能在开始使用后复制。
// 在项目根目录运行：go run ./03_concurrency/06_atomic
// 观察：输出 4000，然后 CompareAndSwap 成功，计数变成 0。
// 练习：解释“余额增加同时扣减库存”为什么不能靠两个独立原子操作保证整体一致。如果用两个原子操作，可能在增加余额后，扣减库存前，另一个 goroutine 也增加了余额和扣减了库存，导致整体不一致。需要用 Mutex 或其他同步机制来保证多个字段的一致性。
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter atomic.Int64
	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 1000; n++ {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
	if counter.Load() != 4000 {
		panic("计数错误")
	}
	fmt.Println(counter.Load())
	fmt.Println("交换成功：", counter.CompareAndSwap(4000, 0))
	fmt.Println(counter.Load())
	// 原子操作适合单个状态；多个字段之间的业务约束通常使用 Mutex。
}
