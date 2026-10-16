/*
是什么：panic 是 Go 的“当前执行流程无法正常继续”机制；recover 是接住正在发生的 panic 的内建函数。
作用：panic 停止当前函数的正常执行，并沿调用链向外退出，途中执行已登记的 defer。

	defer 表示“函数退出时执行清理”；recover 必须由同一 goroutine 的延迟函数直接调用。
	接住后不会从 panic 原位置恢复；恢复边界函数返回后，它的调用者可以继续执行。未接住的 panic 会终止进程。

使用场景：数组越界、解引用 nil 指针等编程错误会触发 panic；可在独立任务或请求边界记录异常并隔离失败。

	参数不合法、文件不存在等可预期失败，通常返回 error，让调用者正常处理。

目标：panic 时执行 defer，理解 recover 的边界。
对照：不要把 panic 当作常规业务异常；参数错误、查询失败等应返回 error。
运行：go run ./02_advanced/06_panic_recover
观察：inner cleanup、recovered: demonstration failure、program continues；panic 后的原语句不会继续执行。
练习：注释 recover 的 defer，观察栈追踪；不要在公共服务中故意触发这个练习。
*/
package main

import "fmt"

func fail() {
	defer fmt.Println("inner cleanup")
	panic("demonstration failure")
}

func boundary() {
	defer func() {
		if cause := recover(); cause != nil {
			fmt.Println("recovered:", cause)
		}
	}()
	fail()
	fmt.Println("unreachable after panic")
}

func main() {
	boundary()
	fmt.Println("program continues")
	// recover 必须在同一 goroutine 的 defer 中直接调用；不能捕获别的 goroutine 的 panic。
}
