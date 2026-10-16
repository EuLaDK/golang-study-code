// 是什么：time 提供时间点 Time、时间间隔 Duration，以及定时器 Timer；时间点回答“什么时候”，间隔回答“多久”。
// 作用：解析和格式化日期、转换时区、计算时间差，或在一段时间后通知任务；Go 用参考时间表达格式布局。
// 使用场景：记录创建时间、判断过期、设置请求超时或测量耗时；同一时刻可以用不同的时区表示。
//
// 目标：时间解析、格式化、时区、时间间隔与可取消定时器。
// 运行（仓库根目录）：go run ./04_stdlib/05_time
// 观察：Go 布局使用参考时间 2006-01-02 15:04:05，不用 YYYY-MM-DD。
// 练习：把 deadline 提前 30 分钟，用 Before 判断是否早于原时间。
package main

import (
	"fmt"
	"time"
)

func main() {
	t, err := time.Parse(time.RFC3339, "2026-09-28T09:00:00+08:00")
	if err != nil {
		panic(err)
	}
	fmt.Println("UTC:", t.UTC().Format(time.RFC3339))
	fmt.Println("local layout:", t.Format("2006-01-02 15:04:05"))
	deadline := t.Add(90 * time.Minute)
	fmt.Println("difference:", deadline.Sub(t), "same instant:", t.Equal(t.UTC()))
	// 比较瞬间用 Equal，不用 ==（后者还比较内部时区/单调时钟信息）。
	started := time.Now()
	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()
	<-timer.C
	fmt.Println("elapsed >= 10ms:", time.Since(started) >= 10*time.Millisecond)
	// 测量耗时优先 time.Since；业务持久化时间通常统一为 UTC。
}
