/*
是什么：闭包是能够访问其定义位置外层变量的函数；即使外层函数已经返回，所需变量也能继续存在。
作用：让函数携带状态，例如记住上一次计数；捕获的是变量，因此多个闭包也可能共享同一个变量。
使用场景：生成回调、计数器或中间件；循环中创建闭包时，要分清每轮新声明的变量与外层共享变量。

目标：闭包捕获变量、独立状态、循环变量的版本语义。
对照：与 JS 闭包类似；Go 1.22 起 for/range 中用 := 声明的循环变量每轮独立。
运行：go run ./02_advanced/07_closures
观察：两个计数器分别为 1 2 和 1；循环内新变量输出 Ada、Lin，循环外共享变量输出 Lin、Lin。
练习：让计数器接收步长；说明多个 goroutine 同时调用同一计数器为何需要同步。
*/
package main

import "fmt"

func newCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func counterWithStep(step int) func() int {
	count := 0
	return func() int {
		count += step
		return count
	}
}

func main() {
	first, second := newCounter(), newCounter()
	fmt.Println("first:", first(), first(), "second:", second())

	names := []string{"Ada", "Lin"}
	var independent []func() string
	for _, name := range names {
		independent = append(independent, func() string { return name })
	}
	fmt.Println("per iteration:", independent[0](), independent[1]())
	var shared []func() string
	var name string             // 这里应该是函数拿的是这个相当于一个外层变量的指针，所以值变了函数去取值的时候都是拿到的一样的值
	for _, name = range names { // 赋值给外层变量，不是每轮声明一个新变量。
		shared = append(shared, func() string { return name })
	}
	fmt.Println("shared variable:", shared[0](), shared[1]())
	firstStep, secondStep := counterWithStep(2), counterWithStep(3)
	fmt.Println("first step:", firstStep(), firstStep(), "second step:", secondStep())
}
