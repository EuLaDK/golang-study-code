package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

// unsafe.Pointer和其他常规类型指针一样，可以作为对象引用。如果一个对象仍然被某个unsafe.Pointer变量引用着，那么该对象是不会被垃圾回收的
// 但是uintptr并不是指针，它仅仅是一个整型值，即便它存储的是某个对象的内存地址，它也不会被算作对该对象的引用。
// 如果认为将对象地址存储在一个uintptr变量中，该对象就不会被垃圾回收器回收，那就是对uintptr的最大误解

type Foo struct {
	name string
}

func finalizer(p *Foo) {
	fmt.Printf("Foo: [%s]被垃圾回收\n", p.name)
}

func NewFoo(name string) *Foo {
	var f Foo = Foo{
		name: name,
	}
	runtime.SetFinalizer(&f, finalizer)
	return &f
}

func allocLargeObject() *[1000000]uint64 {
	a := [1000000]uint64{}
	return &a
}

// 递归压栈，往栈上塞局部变量，迫使栈扩容
func a(depth int) {
	if depth == 0 {
		return
	}
	// 局部数组，分配在栈上，增加栈占用
	var buf [1024]byte
	_ = buf
	a(depth - 1)
}

func main() {
	// 与操作系统的通信主要通过系统调用进行，这在之前已提过。而与非Go编写的代码的通信则主要通过cgo方式
	// func SetIon(iconBytes []byte) {
	// 	// 转成一个C char类型
	// 	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	// 	// 调用来自systray.h的函数
	// 	C.setIcon(cstr,(C.int)(len(iconBytes)))
	// }
	// 	（2）高效类型转换
	// 使用unsafe包，Gopher可以绕开Go类型系统的安全检查，因此可以
	// 通过unsafe包实现性能更好的类型转换。最常见的类型转换是string与
	// []byte类型间的相互转换：
	// func Bytes2String(b []byte) string {
	// return *(*string)(unsafe.Pointer(&b))
	// }
	// func String2Bytes(s string) []byte {
	// sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// bh := reflect.SliceHeader{
	// Data: sh.Data,
	// Len: sh.Len,
	// Cap: sh.Len,
	// }
	// return *(*[]byte)(unsafe.Pointer(&bh))
	// }

	// unsafe.Pointer和uintptr的区别
	var p1 = uintptr(unsafe.Pointer(NewFoo("FooRefByUintptr")))
	var p2 = unsafe.Pointer(NewFoo("FooRefByPointer"))

	for i := 0; i < 5; i++ {
		allocLargeObject()

		// 尝试输出p1和p2地址上的值
		q1 := (*Foo)(unsafe.Pointer(p1))
		fmt.Printf("object ref by uinptr: %+v\n", *q1)

		q2 := (*Foo)(p2)
		fmt.Printf("object ref by pointer: %+v\n", *q2)

		runtime.GC() // 运行垃圾回收
		time.Sleep(1 * time.Second)
	}

	var x = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmt.Printf("变量x的值=%d\n", x)
	println("变量x的地址=", &x)

	// ========== 坑开始 ==========
	// 第一步：把指针转成uintptr裸数字存下来
	// `uintptr`：只是一个普通整数，**不跟踪对象，不会自动更新地址**
	var p = uintptr(unsafe.Pointer(&x))
	// 这是合法的unsafe.Pointer，runtime能跟踪它
	// `unsafe.Pointer`：是 runtime 能识别的指针，**对象搬家时，runtime 自动修正它内部的地址**
	var q = unsafe.Pointer(&x)

	// 中间执行函数！！这里有可能触发栈扩容，x被拷贝到新栈
	a(500) // 递归深度大，栈上累积大量局部变量，更容易触发栈扩容

	// 修改x的值（修改的是新栈上那份拷贝）
	for i := 0; i < 10; i++ {
		x[i] = x[i] + 10
	}

	println("栈扩容后，变量x的地址=", &x)
	fmt.Printf("栈扩容后，变量x的值=%d\n", x)

	// ❌ 违规！隔了很多代码之后，再拿之前存的uintptr转回指针读取
	// p里面存的是【旧栈地址】，x已经搬走了
	valP := *(*[10]int)(unsafe.Pointer(p))
	fmt.Printf("变量p(uintptr)存储的地址上的值=%d\n", valP)

	// ✅ unsafe.Pointer会被runtime自动更新，永远指向x的最新位置
	valQ := *(*[10]int)(q)
	fmt.Printf("变量q(unsafe.Pointer)引用的地址上的值=%d\n", valQ)

	// 提问 但是原来的不会被回收吗？
	// 	栈扩容的时候：

	// 1. `x` 被拷贝到**新栈**，新栈上这份是有效的；
	// 2. **旧栈帧这块内存，不会立刻被 GC 回收！**
	// Go 的 GC 只管**堆**，栈不是 GC 负责回收的。

	// ---

	// ### 栈扩容的完整流程（重点）

	// 1. 栈不够用 → runtime 分配一块更大的新栈；
	// 2. 把当前栈帧里所有变量原样拷贝到新栈；
	// 3. **扫描这个栈帧里所有指针（包括 `unsafe.Pointer` 类型的 q）**，把指针的值全部改成新地址；
	// 4. 旧栈标记成「空闲栈内存」，放回 runtime 的栈内存池，**等待后面别的 goroutine 复用**，不是直接还给操作系统。

	// 👉 所以旧栈里的老数据还躺在那里，只是 runtime 认为这块栈现在不属于当前这个 goroutine 了。

	// - 短时间内没人复用这块内存：你用旧地址去读，还能读到原来的数据（就是书上截图效果）
	// - 一旦别的 goroutine 申请栈、复用了这块内存：里面的数据就被覆盖，你读出来就是乱码

	// >
	// > 注意：**栈根本不走 GC！GC 只扫堆上的对象。**
	// > GC 的 “回收” 这个概念，不能套在栈内存上。

	// ### 对比：如果是堆上对象，会怎么样？

	// 如果 `x` 是在堆上，并且**没有任何指针指向它** → GC 会标记它，可以被回收。
	// 但你这个例子里，旧地址只是一个存放在 `uintptr` 里的数字：

	// - GC**不把 uintptr 当成指针**！GC 扫描的时候直接忽略 uintptr，不会认为这个数字是一个引用；
	// - 就算你拿着这个数字去访问，GC 也不认它是指针引用。

	// ### 一句话区分

	// - `unsafe.Pointer`：GC / 栈迁移扫描器**认得它是指针** → 对象搬家时修正地址；
	// - `uintptr`：在扫描器眼里，它就是一个普通数字，**不是指针**，既不会阻止内存复用，也不会自动更新。

	// ### 再回答你最核心的疑问

	// >
	// > 原来那份旧数据什么时候消失？
	// > 不是 GC 回收，是**栈内存池复用**。什么时候别的 goroutine 需要栈，把这块内存拿去用，覆写内容，旧数据就没了。
	// > 在这之前，数据还在原地，属于一块已经 “失效、不属于你” 的内存，读它就是**未定义行为**。

	// ### 补充一个容易踩坑的点

	// 很多人会误以为：只要对象搬走了，原来的内存就无效、一碰就 panic。
	// 不是的。内存权限还在进程里，只是 runtime 不再管理它了，所以大概率还能读，只是内容不受保证。这也是为什么这个反面例子有时候能读出旧值。
}
