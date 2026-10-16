/*
是什么：作用域是变量可以被访问的范围；遮蔽是内层声明同名变量，使该范围内暂时看不到外层变量。
作用：理解 := 与 = 的区别，以及参数复制后是否仍共享底层数据；数组复制元素，切片和 map 的复制仍可能共享数据。
使用场景：排查外层变量没被更新的问题，或判断函数修改参数后是否会影响调用者。

目标：块级作用域、变量遮蔽、数组/切片/map 的传值区别。
对照：所有参数都按值传递；共享行为取决于复制的值是否包含对底层数据的引用。
运行：go run ./01_basics/12_scope_and_values
观察：内层 message 不改变外层；数组保持 [1 2]，切片和 map 的元素变成 99。
练习：把内层的 := 改为 =，再次运行；解释两次 message 输出的差异。
*/
package main

import "fmt"

func change(array [2]int, slice []int, mapping map[string]int) {
	array[0] = 99
	slice[0] = 99
	mapping["value"] = 99
}

func main() {
	message := "outer"
	{
		message := "inner"
		fmt.Println("inside:", message)
	}
	fmt.Println("outside:", message)
	array := [2]int{1, 2}
	slice := []int{1, 2}
	mapping := map[string]int{"value": 1}
	change(array, slice, mapping)
	fmt.Println("array:", array, "slice:", slice, "map value:", mapping["value"])
}
