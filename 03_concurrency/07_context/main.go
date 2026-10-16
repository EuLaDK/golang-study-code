// 是什么：context 是沿函数调用链传递取消信号、截止时间和少量请求元数据的对象。
// 作用：WithCancel 创建可取消的子 context，WithTimeout 设置最长等待时间，Done 通道通知任务应该停止。
//
//	cancel 只是发通知，任务需检查信号或调用支持 context 的 API；需要确认退出时，还要等待任务完成。
//
// 使用场景：客户端断开后停止无用工作，给上游请求或数据库操作设置超时，并传递请求标识。
//
// 目标：传播取消信号、截止时间和请求级元数据。
// 对照 JS：类似 AbortController/AbortSignal；取消是协作式的，不会强制杀死 goroutine。
// 在项目根目录运行：go run ./03_concurrency/07_context
// 观察：先输出 request-42 和 context canceled，再输出 context deadline exceeded。
// 练习：把父 context 传给两个子任务，验证一次 cancel 能通知两者。
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type requestIDKey struct{} // 私有类型避免不同包使用同一个字符串 key 发生冲突。

func waitForCancellation(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func main() {
	parent := context.WithValue(context.Background(), requestIDKey{}, "request-42")
	ctx, cancel := context.WithCancel(parent)
	defer cancel() // 即使提前返回，也释放 context 关联的资源。
	finished := make(chan error, 1)
	go func() {
		finished <- waitForCancellation(ctx)
	}()
	fmt.Println(ctx.Value(requestIDKey{}), "requestIDKey")
	cancel()
	err := <-finished // 明确等待子任务退出，而不是调用 cancel 后直接认为它已结束。
	if !errors.Is(err, context.Canceled) {
		panic("没有收到取消通知")
	}
	fmt.Println(err, "err")

	deadlineCtx, stop := context.WithTimeout(parent, 10*time.Millisecond)
	defer stop()
	err = waitForCancellation(deadlineCtx)
	if !errors.Is(err, context.DeadlineExceeded) {
		panic("没有收到超时通知")
	}
	fmt.Println(err, "deadlineCtx")
	// context 通常作为函数第一个参数；Value 放请求标识等元数据，业务参数显式传递。
}
