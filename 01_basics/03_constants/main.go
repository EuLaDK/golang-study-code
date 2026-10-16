/*
是什么：const 声明编译时就能确定、运行时不能修改的常量；iota 在一组常量声明中按行递增编号。
作用：给固定值起名字；用 1 << iota 为不同能力分配不同二进制位，再通过 | 组合、& 检查某一位。
使用场景：定义状态码、默认上限和一组权限标志；iota 更适合内部编号，外部协议中固定编号应明确指定。

目标：编译期常量、iota、位标志。
对照：JS 的 const 约束绑定不被重新赋值；Go 的 const 要求值在编译期可确定。
运行：go run ./01_basics/03_constants
观察：状态编号为 0 1 2，读写权限为 3，包含写权限为 true。
练习：增加 CanDelete 权限，计算同时具备三种权限的值。
*/
package main

import "fmt"

const application = "learning-api"

const (
	Pending = iota
	Running
	Done
)

const (
	CanRead = 1 << iota // 左移运算符，0000 -> 0001就是1
	CanWrite
	CanDelete
)

func main() {
	const limit = 10 // 无类型常量可以在允许的范围内赋给不同数字类型。
	var small int8 = limit
	var large float64 = limit
	permissions := CanRead | CanWrite // 或运算 0001 | 0010 = 0011(3)
	fmt.Println(application, Pending, Running, Done)
	fmt.Println("limits:", small, large)
	fmt.Println("permissions:", permissions, "can write:", permissions&CanWrite != 0) // 0011 & 0010 = 0010(2) 不为0
	if CanDelete == 4 {
		fmt.Println(large+1, CanDelete, "可以删除")
	}
}
