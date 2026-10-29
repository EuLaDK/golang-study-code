package main

import (
	"fmt"
	"unsafe"
)

// unsafe.Pointer的六条安全使用模式
// 模式1、*T1 -> unsafe.Pointer -> *T2 本质上是内存块的重解释：将原本解释为T1类型的内存重新解释为T2类型
// $GOROOT/src/math/unsafe.go
// func Float64bits(f float64) uint64 {
// return *(*uint64)(unsafe.Pointer(&f))
// }
// &f                // 获取float64变量f的地址（*float64）
// unsafe.Pointer(&f)// 转成通用unsafe指针（万能桥接）
// (*uint64)(...)    // 把这个指针强行解释成 *uint64
// *()               // 解引用，读出这块内存，当成uint64
// func Float64frombits(b uint64) float64 {
// return *(*float64)(unsafe.Pointer(&b))
// }

// 模式2、unsafe.Pointer -> uintptr 模式2比较简单，就是将unsafe.Pointer显示转换为uintptr，并且转换后的uintptr类型变量不会再转换回unsafe.Pointer,只用于打印输出，并不参与其他操作

// 模式3：模拟指针运算 在一个表达式中，将unsafe.Pointer转换为uinptr类型，使用uinptr类型的值进行算术运算后，再转换回unsafe.Pointer
// 模式3经常用于访问结构体内字段或数组中的元素，也常用于实现对某内存对象的步进式检查：
// 模式三两点注意事项 1-不要越界 2-unsafe.Pointer -> uintptr -> unsafe.Pointer的转换要在一个表达式中

type Foo struct {
	s string
	b int
	c float64
	d [10]int
}

// 模式4：调用syscall.Syscall系列函数时指针类型到uintptr类型参数的转换
// var p *T // 待传给Syscall系列函数的指针变量
// syscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p))
// Go编译器会识别出这种特殊的使用模式，并保证在这个转换过程
// 中原内存对象（p）的有效性

// 模式5：将reflect.Value.Pointer或reflect.Value.UnsafeAddr转换为指针

// 模式6：reflect.SliceHeader和reflect.StringHeader必须通过模式1构建
func newSlice() *[]byte {
	var b = []byte("hello, gopher")
	return &b
}

func main() {
	// 模式一案例代码
	var a = [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var p = (*uint64)(unsafe.Pointer(&a[5]))

	fmt.Printf("a[5]的地址：0x%p\n", &a[5])
	fmt.Printf("a[5]的对齐系数：%d\n", unsafe.Alignof(a[5]))
	fmt.Printf("p的地址：0x%p\n", p)
	fmt.Printf("p的对齐系数：%d\n", unsafe.Alignof(p))
	fmt.Printf("*p = %d\n", *p)
	// 在使用模式一时要注意，转换后类型T2的对齐系数不能比转换前类型T1对齐系数更严格，即Alignof(T1) >= Alignof(T2)

	// 模式二案例代码
	var x2 = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var p2 = uintptr(unsafe.Pointer(&x2))
	fmt.Println(p2)

	// 模式三案例代码
	fmt.Println("模式三示例代码")
	var a3 = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var foo = Foo{
		s: "foo",
		b: 17,
		c: 3.1415,
		d: a3,
	}

	// ========== 1. 访问数组a的第4个元素（索引3，从0开始） ==========
	// a[0]偏移0，a[1]偏移 Sizeof(int)，a[3]偏移 3*unsafe.Sizeof(a[0])
	// 	1. `&a3`：拿到整个数组的起始地址，类型是 `*[10]int`
	// 2. `unsafe.Pointer(&a3)`：转为万能指针，打通类型壁垒
	// 3. `uintptr(...)`：把这个地址变成一个纯数字，现在可以做加法了
	// 4. `3*unsafe.Sizeof(a3[0])`：算出偏移多少字节：3 个 int 的总字节长度
	// 5. 加法：起始地址数字 + 偏移字节数 → 得到 a3 [3] 的内存地址数字
	// 6. `unsafe.Pointer(...)`：把计算好的地址数字，转回万能指针
	// 7. `(*int)(p3)`：告诉编译器：这块内存，按 int 来解释
	// 8. `*()`：解引用，读出这块内存的值，就是 4
	var p3 = unsafe.Pointer(uintptr(unsafe.Pointer(&a3)) + 3*unsafe.Sizeof(a3[0]))
	fmt.Println("a[3] =", *(*int)(p3)) // 4

	// ========== 2. 访问Foo结构体的字段c（float64） ==========
	// Offsetof 传字段：foo.c
	// &foo 起始地址：0x2000
	// 0x2000 ~ 0x2017  s string （24字节）
	// 0x2018 ~ 0x201F  b int     （8字节）
	// 0x2020 ~ 0x2027  c float64 （8字节） ← foo.c
	// 0x2028 ~ ...     d [10]int
	p3 = unsafe.Pointer(uintptr(unsafe.Pointer(&foo)) + unsafe.Offsetof(foo.c))
	fmt.Println("foo.c =", *(*float64)(p3)) // 3.1415

	// ========== 3. 对数组a的第一个元素 a[0] 逐字节步进检查 ==========
	// 取 a[0] 的地址，遍历这个int的每一个字节
	base := unsafe.Pointer(&a3[0])
	fmt.Println("\n--- a[0] 的内存字节（十六进制）---")
	for i := uintptr(0); i < unsafe.Sizeof(a3[0]); i++ {
		p3 = unsafe.Pointer(uintptr(base) + i)
		fmt.Printf("byte[%d] = 0x%x\n", i, *(*byte)(p3))
	}

	// ========== 额外再加一个：演示取 foo.b 做对比 ==========
	p3 = unsafe.Pointer(uintptr(unsafe.Pointer(&foo)) + unsafe.Offsetof(foo.b))
	fmt.Println("\nfoo.b =", *(*int)(p3)) // 17

	// ========== 额外：演示访问 foo.d[2] （结构体里的数组） ==========
	// 先拿到foo.d起始地址，再加偏移
	dBase := unsafe.Pointer(uintptr(unsafe.Pointer(&foo)) + unsafe.Offsetof(foo.d))
	p3 = unsafe.Pointer(uintptr(dBase) + 2*unsafe.Sizeof(a3[0]))
	fmt.Println("foo.d[2] =", *(*int)(p3)) // 3
}
