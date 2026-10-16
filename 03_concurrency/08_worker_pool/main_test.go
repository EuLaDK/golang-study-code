// 是什么：这是工作池的自动化测试；工作池通过固定数量的 goroutine 从队列中领取任务。
// 作用：检查不同 worker 数量下的结果与顺序，验证空输入可以结束、非法 worker 数量会返回错误。
// 使用场景：调整批量任务的并发量或任务分发方式后，确认不会漏结果、乱序或因没有消费者而阻塞。
// 运行：go test ./03_concurrency/08_worker_pool；slices.Equal 比较两份结果的元素是否相同。
//
package main

import (
	"slices"
	"testing"
)

func TestSquareAll(t *testing.T) {
	for _, workers := range []int{1, 2, 8} {
		got, err := squareAll([]int{3, 0, -4, 2}, workers)
		if err != nil || !slices.Equal(got, []int{9, 0, 16, 4}) {
			t.Fatalf("workers=%d: got %v, err=%v", workers, got, err)
		}
	}
	if got, err := squareAll(nil, 2); err != nil || len(got) != 0 {
		t.Fatalf("empty input: got %v, err=%v", got, err)
	}
	if _, err := squareAll([]int{1}, 0); err == nil {
		t.Fatal("zero workers must be rejected instead of deadlocking")
	}
}
