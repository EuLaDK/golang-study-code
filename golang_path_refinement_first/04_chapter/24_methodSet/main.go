package main

import (
	"fmt"
	"io"
	"reflect"
)

// 对于非接口类型的自定义类型T，其方法集合由所有receiver为T类型的方法组成；而类型*T的方法集合则包含所有receiver为T和*T类型的方法
// 一句话总结，只要是指针就包含非指针的，非指针没有指针的
func DumpMethodSet(i interface{}) {
	v := reflect.TypeOf(i)
	elemTyp := v.Elem()
	n := elemTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", elemTyp)
		return
	}
	fmt.Printf("%s's method set:\n", elemTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", elemTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}

type Interface interface {
	M1()
	M2()
}
type T struct{}

func (t T) M1()  {}
func (t *T) M2() {}

// Go 选择方法的次序
// 1、优先选择结构体自身实现的方法
// 2、如果结构体自身并未实现，那么将查找结构体中嵌入接口类型的方法集合中是否有该方法，如果有，则提升（promoted）为结构体的方法

type Interface1 interface {
	M11()
	M22()
}

type T1 struct {
	Interface1
}

func (T1) M11() {
	fmt.Println("T1 is M11")
}

type S struct{}

func (S) M11() {
	fmt.Println("S is M11")
}

func (S) M22() {
	fmt.Println("S is M22")
}

//  已有的类型I，T被称为underlying类型，而新类型被称为defined类型，就是定义了没赋值可以这么理解
// type MyInterface I
// type Mystruct T

type MT struct{}

func (MT) MT1()  {}
func (*MT) MT2() {}

type MyInterface interface {
	MT1()
	MT2()
}

type MT1 MT                   // Defined 类型不继承方法 使用 type T1 T 定义了一个全新的类型 T1，T1 是白手起家的，它绝对不会继承 T 身上的任何方法！
type MyInterface1 MyInterface // 接口 Defined 会继承方法 使用 type Interface1 Interface 定义新接口，因为接口只是“契约”，Interface1 拥有和 Interface 完全一样的方法列表

func main() {
	var t T
	var pt *T
	DumpMethodSet(&t)
	DumpMethodSet(&pt)
	DumpMethodSet((*Interface)(nil)) // 这是一个强制类型转换 把一个 nil（空指针），强行说成是 *Interface 类型的指针
	DumpMethodSet((*io.Writer)(nil))
	DumpMethodSet((*io.Reader)(nil))
	DumpMethodSet((*io.Closer)(nil))
	DumpMethodSet((*io.ReadWriter)(nil))
	DumpMethodSet((*io.ReadWriteCloser)(nil))

	var t1 = T1{
		Interface1: S{},
	}
	t1.M11()
	t1.M22()

	var mt MT
	var pmt *MT
	var mt1 MT1
	var pmt1 *MT1
	fmt.Println("------------------------")
	DumpMethodSet(&mt)
	DumpMethodSet(&mt1)
	DumpMethodSet(&pmt)
	DumpMethodSet(&pmt1)
	DumpMethodSet((*MyInterface)(nil))
	DumpMethodSet((*MyInterface1)(nil))
}
