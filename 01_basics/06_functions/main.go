/*
是什么：函数是有名字或可保存到变量中的一段处理逻辑，接收参数并返回结果；Go 支持一次返回多个值。
作用：(结果, error) 同时传递成功数据和失败原因；error 是描述失败的值，nil error 通常表示成功。
      ...int 表示接收不定数量的整数；把函数当作值可以传给其他函数或稍后调用。
使用场景：把输入解析、校验、计算等步骤拆成可复用的小函数，并让调用者决定如何处理失败。

目标：参数类型、多返回值、可变参数、函数也是值。
对照：Go 常用 (value, error) 表达可预期失败，没有可选参数或函数重载。
运行：go run ./01_basics/06_functions
观察：divide: 4；error: divisor must not be zero；sum: 6；double: 8。
练习：增加接收两个整数并返回较大值的函数，然后把它保存到变量里调用。
*/
package main

import (
	"errors"
	"fmt"
)

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divisor must not be zero")
	}
	return a / b, nil
}

func sum(values ...int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func isMax(a, b int) int {
	if a > b {
		return a
	} else if a < b {
		return b
	} else {
		return 0
	}
}

func main() {
	result, err := divide(8, 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("divide:", result)
	if _, err := divide(8, 0); err != nil {
		fmt.Println("error:", err)
	}
	values := []int{1, 2, 3}
	fmt.Println("sum:", sum(values...))
	double := func(value int) int { return value * 2 }
	fmt.Println("double:", double(4))
	max := isMax(10, 20)
	println(max / 5)
}
