/*
是什么：变量是有名字、有类型、可以重新赋值的数据；int、string、bool 分别表示整数、字符串和布尔值。
作用：var 声明变量，函数内的 := 可以声明变量并推断类型；零值是未手动赋值时的默认值，如 0、空字符串和 false。
使用场景：保存用户 ID、配置和处理结果；不同类型之间需要明确转换，例如把整数转换为浮点数后再做除法。

目标：变量声明、类型推断、零值、显式转换。
对照：类似 TS 的静态类型，但 Go 无隐式数字/字符串转换；未初始化变量有零值。
运行：go run ./01_basics/02_variables
观察：零值为 0、""、false；float64(3)/2 为 1.5，直接 3/2 为 1。
练习：将 age 改为 string，观察编译器如何拒绝原来的加法。
*/
package main

import "fmt"

func main() {
	var count int
	var title string
	var ready bool
	fmt.Printf("zero values: %d %q %t\n", count, title, ready)

	name := "Ada" // 函数内可用 :=，类型由右侧表达式推断。
	var age int = 28
	age = age + 1
	name, active := "Ada Lovelace", true // 同一作用域至少有一个新变量才可使用 :=。
	fmt.Println(name, age, active)
	var small int32 = 3
	large := int64(small)
	fmt.Printf("%T -> %T, division: %.1f / %d\n", small, large, float64(small)/2, small/2)
}
