package main

import (
	"fmt"
	"unsafe"
)

// 测试二级指针用的
func foo(p **int) {
	fmt.Println(p, "二级指针")
	var b int = 40
	*p = &b
	fmt.Println(*p)
}

func main() {
	// 声明一个指针变量
	// var p *T
	// 如果有一个整型int 那么它的指针类型就是*int 这里int被称为*int的基类型
	// Go中 有一种特殊的指针类型不需要基类型 就是 unsafe.Pointer,任何指针类型都可以显示转换为unsafe.Pointer，反之亦然
	var p *int
	var p1 = unsafe.Pointer(p) // 将任意指针类型显示转换为 unsafe.Pointer
	p = (*int)(p1)             // 将unsafe.Pointer 显示转换为任意指针类型

	// 指针类型变量未被显式赋予初始值，那么它的值将是 nil
	var count *int
	fmt.Println(count)

	// 给指针赋值 必须要类型相同
	var a int = 6
	var p2 *int = &a
	fmt.Println(a, p2)

	// 通过指针变量读取或修改其指向的内存地址上的变量值的操作被称为指针的解引用
	*p2 += 2
	fmt.Println(a, *p2)

	// 指向指针类型变量的指针 一级指针通常用来改变普通变量的值，二级指针可以用来改变指针变量的值，也就是改变指针变量的指向
	var p3 **int = &p2
	fmt.Println("p3", **p3)
	**p3 += 2
	fmt.Println(a, *p2, **p3, "p2的地址", p2)
	foo(&p2)
	fmt.Println(a, *p2, **p3, "p2的地址", p2)

	// 如果需要不同类型的指针转换，需要使用unsafe.Pointer
	var p4 int = 0x12345678
	var p5 *int = &p4
	var p6 *byte = (*byte)(unsafe.Pointer(p5)) // 正确
	fmt.Printf("%x\n", *p6)                    // 当使用unsafe包中的类型或函数时，代码的安全性需要由开发者自行保证，这意味着开发者需要明确知道自己在做什么
}
