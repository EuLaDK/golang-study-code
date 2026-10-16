package main

import (
	"fmt"
	"io"
)

// 接口类型的定义
type MyInterface interface {
	M1(a int) error
	M2(w io.Writer, strs ...string)
}

// go要求接口类型声明中的方法必须是具名的，并且方法名在其所属接口类型的方法集合中是唯一的
// type Interface1 interface {
// 	M1()
// }
// type Interface2 interface {
// 	M1(string)
// 	M2()
// }
// type Interface3 interface {
// 	Interface1
// 	Interface2 // 编译器错误：duplicate method M2
// 	M3()
// }

// 接口类型一旦被定义，就可以向其他Go类型一样用于声明变量
// 这些类型被称为接口类型变量，默认值为nil
// Go规定：如果一个类型的 T的方法集合是接口类型 I 的方法集合的等价集合或超集，就说类型T实现了接口类型 I，这意味着类型T的变量可以作为合法的右值赋值给接口类型I的变量
var err error   // err是一个error接口类型的实例变量
var r io.Reader // r是一个io.Reader 接口类型的实例变量

// 类型断言
// v, ok := i.(T)
// i是某个接口类型变量，如果T是一个非接口类型，并且是你想要还原的类型，那么这段代码的作用是断言存储咋接口类型变量i中的值的类型为T

type MyInterface2 interface {
	M1()
}

type TT int

func (TT) M1() {
	println("TT M1")
}

// 隐式契约，无需签署，自动生效。在 Go 中，接口类型与其实现者之间的关系是隐式的。与 Java 等其他编程语言不同，Go 不要求实现者显式声明 implements 关键字来表示实现了某个接口。只要一个类型实现了接口方法集合中的全部方法，它就被认为遵守了契约，并立即生效。
// 倾向于“小契约”
// 。如果契约过于繁杂，将会束缚手脚，减少灵活性并抑制表现力。所以，Go 提倡使用“小契约”
// ，这意味着尽量定义小接口，通常包含 1~3 个方法。正如 Go 之父罗伯·派克所说：
// “接口越大，抽象程度越弱。
// ”这句话也体现了 Go 社区倾向于定义小接口的原因

// 利用接口实现鸭子类型
type QuackableAnimal interface {
	Quack()
}

type Duck struct{}

func (Duck) Quack() {
	println("duck quack!")
}

type Dog struct{}

func (Dog) Quack() {
	println("dog quack!")
}

type Bird struct{}

func (Bird) Quack() {
	println("bird quack!")
}

func AnimalQuackInForest(a QuackableAnimal) {
	a.Quack()
}

// 在使用接口前，一定要搞清楚自己使用接口的目的，千万不能为了使用接口而使用接口。此外，还要牢记：尽量避免将空接口作为函数参数类型。
// 一旦将空接口作为函数参数类型，你便失去了编译器提供的类型安全保护屏障

func main() {
	// 任何类型都自动实现了空接口类型
	var i interface{} = 15
	i = "hello, golang"
	println(i)

	// 类型断言示例代码
	var a int64 = 13
	var t interface{} = a
	v1, ok := t.(int64)
	fmt.Printf("v1=%d, the type of v1 is %T, ok=%t\n", v1, v1, ok)
	v2, ok := t.(string)
	fmt.Printf("v2=%s, the type of v2 is %T, ok=%t\n", v2, v2, ok)
	v3 := t.(int64)
	fmt.Printf("v3=%d, the type of v3 is %T\n", v3, v3)
	// v4 := t.([]int) 抛出一个panic
	// fmt.Printf("the type of v4 is %T\n", v4)

	// T为接口类型的示例
	var tt TT
	var ti interface{} = tt
	vt, ok := ti.(MyInterface2)
	if !ok {
		panic("the value of i is not MyInterface2")
	}
	vt.M1()
	fmt.Printf("the type of vt is %T\n", vt)

	ti = int64(13)
	vt2, ok := ti.(MyInterface2)
	fmt.Printf("the type of v2 is %T\n", vt2)

	// 接口实现鸭子叫
	animals := []QuackableAnimal{new(Duck), new(Dog), new(Bird)}
	for _, animal := range animals {
		AnimalQuackInForest(animal)
	}
}
