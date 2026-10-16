/*
是什么：切片是对底层数组某一段的描述，保存数据位置、长度 len 和容量 cap；[]int 表示整数切片。
作用：append 增加元素，copy 复制元素；多个切片可能共享底层数组，因此修改元素可能影响另一个切片。
      len 是当前可访问的元素数，cap 是从切片起点到底层数组末尾可容纳的元素数；扩容可能换一个底层数组。
使用场景：表示查询结果、请求列表和批量任务；需要互不影响的数据时，显式复制元素。

目标：切片的长度/容量、共享底层数组、append、独立复制。
对照：切片是底层数组的一段视图；传值复制切片描述信息，不会复制元素。
运行：go run ./01_basics/08_slices
观察：修改 view 会修改 original；修改 clone 不会。满容量 append 后新元素不影响原数组。
练习：移除 full slice 表达式中的第三个下标，观察 append 如何覆盖原数组中的 30。
*/
package main

import "fmt"

func clone(values []int) []int {
	result := make([]int, len(values))
	copy(result, values)
	return result
}

func main() {
	original := []int{10, 20, 30}
	view := original[:2]
	view[0] = 99
	independent := clone(original)
	independent[0] = 1
	fmt.Println("shared:", original, view, "clone:", independent)
	fmt.Println("view len/cap:", len(view), cap(view))

	limited := original[:2:2] // 限制容量为 2，让 append 必须分配新的底层数组。
	appended := append(limited, 40)
	fmt.Println("after append:", original, appended)
	var empty []int
	empty = append(empty, 7) // nil 切片可以直接 append；索引访问仍需要检查长度。
	fmt.Println("append to nil:", empty)
}
