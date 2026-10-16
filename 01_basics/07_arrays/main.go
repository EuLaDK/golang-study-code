/*
是什么：数组是长度固定、元素类型相同的一组数据；[3]int 表示恰好包含三个整数。
作用：用索引读取和修改元素；数组赋值会复制元素，长度也是类型的一部分，因此 [3]int 与 [4]int 是不同类型。
使用场景：表示固定长度的摘要或固定大小的数据；长度经常变化的列表通常使用下一节的切片。

目标：固定长度数组、值复制、数组比较。
对照：Go 数组不是 JS Array 或 Python list；日常动态集合通常用下一节的切片。
运行：go run ./01_basics/07_arrays
观察：original 为 [10 20 30]，copied 为 [99 20 30]，数组赋值后元素独立。
练习：新增 view := original[:] 并修改 view[0]，猜测 original 的输出如何变化。
*/
package main

import "fmt"

func main() {
	original := [3]int{10, 20, 30}
	copied := original
	view := original[:]
	view[0] = 100
	copied[0] = 99
	fmt.Println("original:", original, "copied:", copied, "view", view)
	fmt.Println("equal:", original == copied, "length:", len(original))
	weekend := [...]string{"Saturday", "Sunday"}
	fmt.Printf("inferred array type: %T\n", weekend) // `%T` 是格式化动词：**print the type of the value**，输出变量是什么类型
	grid := [2][2]int{{1, 2}, {3, 4}}
	fmt.Println("grid[1][0]:", grid[1][0])
}
