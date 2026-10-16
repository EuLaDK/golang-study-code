// 是什么：这是切片示例的自动化测试文件；断言是把实际结果与预期比较，不符合就让测试失败。
// 作用：检查视图共享、独立复制和满容量追加三种行为，避免把“复制切片变量”误认为“复制元素”。
// 使用场景：修改集合处理代码后确认共享关系没有意外变化；测试只运行明确写出的检查，不自动运行 main。
// 运行：go test ./01_basics/08_slices；t.Fatal 报告失败并停止当前测试。
//
package main

import "testing"

func TestSliceSharingAndCopy(t *testing.T) {
	original := []int{10, 20, 30}
	view := original[:2]
	view[0] = 99
	if original[0] != 99 {
		t.Fatal("切片视图应共享底层数组")
	}
	independent := clone(original)
	independent[0] = 1
	if original[0] != 99 {
		t.Fatal("clone 应复制元素")
	}
	appended := append(original[:2:2], 40)
	appended[0] = 2
	if original[0] != 99 || original[2] != 30 {
		t.Fatal("满容量追加应分配新的底层数组")
	}
}
