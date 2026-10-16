// 是什么：这个文件同时包含正确性测试和性能基准；testing.B 用于基准测量，b.N 是工具决定的运行次数。
// 作用：先检查求和结果，再重复测量；ResetTimer 排除准备数据的时间，ReportAllocs 显示内存分配统计。
// 使用场景：在功能正确的前提下比较实现成本；保存最后结果是为了减少“结果没被使用”带来的编译优化干扰。
// 运行：go test ./06_testing/03_benchmarks -run '^$' -bench . -benchmem；-run '^$' 跳过普通测试。
//
package main

import "testing"

// 保存结果，让编译器不能因为结果完全没用而删除被测计算。
var benchmarkResult int64

func TestSumSizes(t *testing.T) {
	if got := sumSizes([]int64{100, 200, 300}); got != 600 {
		t.Fatalf("got %d, want 600", got)
	}
	if got := sumSizes(nil); got != 0 {
		t.Fatalf("empty input: got %d, want 0", got)
	}
}

func BenchmarkSumSizes(b *testing.B) {
	sizes := make([]int64, 1000)
	for i := range sizes {
		sizes[i] = int64(i)
	}
	b.ReportAllocs()
	b.ResetTimer() // 排除准备数据的时间；b.N 由 Go 自动调整。
	var total int64
	for i := 0; i < b.N; i++ {
		total = sumSizes(sizes)
	}
	benchmarkResult = total
}
