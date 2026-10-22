package main

import (
	"fmt"
	"unsafe"
)

// 接口类型变量的内部表示
// 接口类型变量在运行时的表示：
// 用于表示其余拥有方法的接口（interface）类型变量
type iface struct {
	tab  *itab
	data unsafe.Pointer // 指向当前赋值给该接口类型变量的动态类型变量的值
}

// 而iface除了要存储动态类型信息之外，还要存储接口本身的信息
// （接口的类型信息、方法列表信息等）以及动态类型所实现的方法的信
// 息，因此iface的第一个字段指向一个itab类型结构
type itab struct {
	inter *interfacetype // 存储着该接口类型自身的信息
	_type *_type
	hash  uint32
	_     [4]byte
	fun   [1]uintptr
}

type interfacetype struct {
	typ _type // 类型信息
	// pkgpath name      // 包路径名
	// mhdr    []imethod // 接口方法集合切片
}

// 表示没有方法的空接口，类型变量，即interface{}类型的变量
type eface struct {
	_type *_type         // eface所表示的空接口类型并无方法列表，因此其第一个指针字段指向一个_type类型结构，该结构为该接口类型变量的动态类型的信息
	data  unsafe.Pointer // 指向当前赋值给该接口类型变量的动态类型变量的值
}

type _type struct {
	size    uintptr
	ptrdata uintptr
	hash    uint32
	// tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	// alg        *typeAlg
	gcdata *byte
	// str        nameOff
	// ptrToThis  typeOff
}

// eface 表示空接口类型变量的例子
type T struct {
	n int
	s string
}

// iface表示非空接口类型变量的例子
type T1 struct {
	n int
	s string
}

func (T1) M1() {}
func (T1) M2() {}

type NonEmptyInterface interface {
	M1()
	M2()
}

func main() {
	var t = T{
		n: 17,
		s: "hello,interface",
	}
	var ei interface{} = t // Go运行时使用eface结构表示ei
	var t1 = T1{
		n: 18,
		s: "hello, interface",
	}
	var i NonEmptyInterface = t1
	fmt.Println(ei, i)
}

// 小结
// 接口类型变量在运行时表示为eface和iface，eface用于表示空接口类型变量，iface用于表示非空接口类型变量
// 当且仅当两个接口类型变量的类型信息(eface._type/iface.tab._type)相同，且数据指针(eface.data/iface.data)所指数据相同时，两个接口类型才是相等的
// 通过println可以输出接口类型变量的两部分指针变量的值
// 可通过复制runtime包eface和iface相关类型源码，自定义输出eface/iface详尽信息的函数
// 接口类型变量的装箱操作有Go编译器和运行时共同完成
