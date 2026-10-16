/*
是什么：泛型通过类型参数，让一份函数或类型定义支持多种具体类型，同时保留编译期检查。
作用：约束规定允许的类型和操作；comparable 表示可用 == 比较，~int 表示底层类型为 int 的类型也可使用。
使用场景：编写对多种数值类型通用的计算或集合操作；业务类型明确时直接用具体类型，避免增加阅读负担。

目标：类型参数、联合约束、底层类型 ~、comparable。
对照：类似 TS 泛型，但 Go 约束还控制可用的运算符；any 不代表可以随意相加。
运行：go run ./02_advanced/08_generics
观察：int sum=6、float sum=4、Score sum=30；字符串查找为 true。
练习：为 Number 增加 int32；尝试传入 []string，阅读编译错误后恢复。
*/
package main

import "fmt"

type Number interface {
	~int | ~int64 | ~float64 | ~int32
}

func sum[T Number](values []T) T {
	var total T
	for _, value := range values {
		total += value
	}
	return total
}

func contains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

type Score int

func main() {
	fmt.Println("int sum:", sum[int]([]int{1, 2, 3})) // 可以这么写但是不必要
	fmt.Println("float sum:", sum([]float64{1.5, 2.5}))
	fmt.Println("Score sum:", sum([]Score{10, 20}))
	fmt.Println("contains:", contains([]string{"GET", "POST"}, "POST"))
	fmt.Println("empty sum:", sum[int](nil)) //
	// fmt.Println("errors test", sum([]string{"tom", "jerry", "name"})) // 报错，不符合要求
	// contains 用于展示约束；实际代码可直接用标准库 slices.Contains。
}
