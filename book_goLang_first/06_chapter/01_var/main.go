// 第六章变量
package main

import (
	"fmt"
)

// 变量声明 var 变量名(a) 类型(int) = 10
var a int = 10

// 声明多个
var a1, a2, a3 int = 1, 2, 3

// 集中声明
var (
	a4 int    = 4
	a5 bool   = true
	a6 string = "hello"
)

// 省略类型声明 自动推导类型
var a7 = 7

func main() {
	// 简短声明 := 只能在函数体内使用
	a8 := 8
	// 显示转换类型
	a9 := int64(a8)
	var (
		a10 int
		a11 string
	)
	fmt.Println(a10, a11, a9)
}
