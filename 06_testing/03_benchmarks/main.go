// 是什么：基准测试 Benchmark 重复运行同一项操作，测量每次操作的耗时和内存分配。
// 作用：testing 自动调整运行次数，报告 ns/op、B/op、allocs/op，即每次耗时、分配字节数和分配次数。
// 使用场景：比较实现或发现优化方向；先确认结果正确，再测量，避免把机器波动误认为性能改进。
//
// 目标：用 Benchmark 测量执行时间和内存分配，不凭感觉判断快慢。
// 根目录运行：go run ./06_testing/03_benchmarks
// PowerShell 基准命令：go test ./06_testing/03_benchmarks -run '^$' -bench . -benchmem
// 观察：普通运行输出 600；基准报告 ns/op、B/op 和 allocs/op，数值因机器而异。
// 练习：把输入大小从 1000 改到 10000，观察时间变化；不要把时间写成测试断言。
package main

import "fmt"

func sumSizes(sizes []int64) int64 {
	var total int64
	for _, size := range sizes {
		total += size
	}
	return total
}

func main() {
	fmt.Println(sumSizes([]int64{100, 200, 300}))
	// 教学输入是小整数；处理不可信的大数时，需要额外定义负数和溢出的规则。
}
