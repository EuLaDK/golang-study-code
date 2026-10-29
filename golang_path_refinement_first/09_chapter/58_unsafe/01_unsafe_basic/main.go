package main

import (
	"fmt"
	"unsafe"
)

type Foo struct {
	a int
	b string
	c [10]byte
	d float64
}

func main() {
	// sizeof用于获取一个表达式值的大小
	var i int = 5
	var a = [100]int{}
	var sl = a[:]
	var f Foo
	var b float64 = 5.00

	fmt.Println(unsafe.Sizeof(i))   // 8
	fmt.Println(unsafe.Sizeof(a))   // 800
	fmt.Println(unsafe.Sizeof(sl))  // 24 (注：返回的是切片描述符的大小)
	fmt.Println(unsafe.Sizeof(f))   // 48
	fmt.Println(unsafe.Sizeof(f.c)) // 10
	fmt.Println(unsafe.Sizeof((*int)(nil)))

	// Alignof用于获取一个表达式的内存地址对齐系数
	// 对齐系数是一个计算机体系架构层面的术语。在不同的计算机体系结构下，处理器对变量地址都有着对齐要求，即变量的地址必须可被改变量的对齐系数整
	fmt.Println("Alignof的示例")
	fmt.Println(unsafe.Alignof(i))          // 8
	fmt.Println(unsafe.Alignof(f.a))        // 8
	fmt.Println(unsafe.Alignof(a))          // 8
	fmt.Println(unsafe.Alignof(sl))         // 8
	fmt.Println(unsafe.Alignof(f))          // 8
	fmt.Println(unsafe.Alignof(f.c))        // 1
	fmt.Println(unsafe.Alignof(struct{}{})) // 1 (注：空结构体的对齐系数为1)
	fmt.Println(unsafe.Alignof([0]int{}))   // 8 (注：长度为0的数组，其对齐系数依然

	// Offsetof 用于获取结构体中某字段的地址偏移量（相对于结构体变量的地址）。Offsetof函数的应用面较窄，近用于求结构体中某字段的偏移值

	// unsafe.Pointer类型。unsafe.Pointer可用于表示任意类型的指针
	// 1、任意类型的置身都可以被转换为unsafe.Pointer
	p1 := (unsafe.Pointer)(&i) // // *int -> unsafe.Pointer
	p2 := (unsafe.Pointer)(&b) // *float64 -> unsafe.Pointer
	p3 := (unsafe.Pointer)(&a) // *[10]string -> unsafe.Pointer
	p4 := (unsafe.Pointer)(&f) // *Foo -> unsafe.Pointer

	// 2、unsafe.Pointer也可以被转换为任意类型的指针值
	var pa = (*int)(p1)          // unsafe.Pointer -> *int
	var pb = (*float64)(p2)      // unsafe.Pointer -> *float64
	var parr = (*[10]string)(p3) // unsafe.Pointer -> *[10]string
	var pf = (*Foo)(p4)          // unsafe.Pointer -> *Foo

	// 3、uintptr类型值可以被转换为一个unsafe.Pointer
	var pi uintptr = 0x80010203
	pi1 := unsafe.Pointer(pi)

	// 4、ubsafe.Pointer也可以被转换为一个unitptr类型值
	pp := unsafe.Pointer(&a)
	var ii = uintptr(pp)
	fmt.Println(pa, pb, parr, pf, pi1, ii)
}
