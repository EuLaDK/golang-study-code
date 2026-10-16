/*
是什么：defer 把一次函数调用登记到当前函数退出时执行，类似“离开前记得做这件事”。
作用：无论正常返回还是因为 panic 向外退出，已登记的 defer 都会执行；多个 defer 按后登记先执行的顺序运行。
      普通调用参数在登记时求值；闭包则可在执行时读取外层变量，所以两种写法可能看到不同值。
使用场景：成功打开文件后安排 Close、加锁后安排 Unlock、创建取消函数后安排 cancel，减少漏掉清理的机会。

目标：defer 的执行时间、参数立即求值、后进先出。可以想成类似于栈的特性
对照：类似在函数退出时运行 finally；defer 不会在 if/for 的块结束时执行。
运行：go run ./02_advanced/05_defer
观察：body 先输出，随后 closure=2、argument=1、last cleanup；命名返回值最终为 2。
练习：在 for 循环里 defer fmt.Println(i)，预测顺序；解释为何大量文件操作不应都 defer 到大函数末尾。
*/
package main

import "fmt"

func demonstrate() {
	defer fmt.Println("last cleanup")
	value := 1
	defer fmt.Println("argument:", value)             // 调用参数现在求值，保存的是 1。
	defer func() { fmt.Println("closure:", value) }() // 闭包执行时读取变量。
	value = 2
	fmt.Println("body")
}

func namedReturn() (result int) { // 命名返回值，这里虽然是return 1，但是实际上返回的是result = 1 ,return result
	defer func() { result++ }()
	return 1 // 先把 result 设为 1，再执行 defer，最终返回 2。
}

func main() {
	demonstrate()
	fmt.Println("named return:", namedReturn())
	for i := 0; i < 5; i++ {
		defer fmt.Println(i)
	}
	// 实际后端中，资源成功打开后立刻安排 defer Close；文件章节会演示关闭错误处理。
}
