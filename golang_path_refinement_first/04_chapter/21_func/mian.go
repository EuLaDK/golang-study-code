package main

import (
	"fmt"
	"net/http"
)

// type HandlerFunc func(ResponseWriter, *Request) // 这个是http.HandlerFunc的定义

func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome, Gopher!\n")
}

type BinaryAdder interface {
	Add(int, int) int
}

type MyAdderFunc func(int, int) int

func (f MyAdderFunc) Add(x, y int) int {
	return f(x, y)
}

func MyAdd(x, y int) int {
	return x + y
}

// Go编写科里化函数
func times(x, y, z int) int {
	return x * y * z
}

// 闭包是在函数内部定义的匿名函数，并且允许该匿名函数访问定义它的外部函数的作用域。本质上，闭包是将函数内部和函数外部连接起来的桥梁
func partialTimes(x int) func(int) func(int) int {
	return func(y int) func(int) int {
		return func(z int) int {
			return times(x, y, z)
		}
	}
}

// 函子
// 函子本身是一个容器类型，以Go一眼为例，这个容器可以是切片、map甚至channel
// 该容器类型需要实现一个方法，该方法接受一个函数类型参数，并在容器的每个元素上应用那个函数，得到一个新函子，原函子容器内部元素值不受影响
type IntSliceFunctor interface {
	Fmap(fn func(int) int) IntSliceFunctor
}

type intSliceFunctorImpl struct {
	ints []int
}

func (isf intSliceFunctorImpl) Fmap(fn func(int) int) IntSliceFunctor {
	newInts := make([]int, len(isf.ints))
	for i, v := range isf.ints {
		retInt := fn(v)
		newInts[i] = retInt
	}
	return intSliceFunctorImpl{ints: newInts}
}

func NewIntSliceFunctor(slice []int) IntSliceFunctor {
	return intSliceFunctorImpl{ints: slice}
}

func main() {
	// 像对整型变量那样对函数进行显式的类型转换
	var a int = 5
	var b int32 = 6
	// fmt.Println(a+b) // 这样会报错因为go认为这是两个类型
	fmt.Println(a + int(b))
	// http.ListenAndServe(":8080", http.HandlerFunc(greeting)) // http.HandlerFunc把我们自己的greeting函数转成http.HandlerFunc类型的，将函数greeting显式转换为HandlerFunc类型
	var i BinaryAdder = MyAdderFunc(MyAdd)
	fmt.Println(i.Add(5, 6))
	timesTwo := partialTimes(2)
	timesThree := partialTimes(3)
	timesFour := partialTimes(4)
	timesFive := timesTwo(3)
	timesSix := timesThree(4)
	timeSeven := timesFour(5)
	fmt.Println(timesFive(4))
	fmt.Println(timesSix(5))
	fmt.Println(timeSeven(6))

	// 函子
	// 原切片
	intSlice := []int{1, 2, 3, 4}
	fmt.Printf("init a function form int slice: %#v\n", intSlice)
	f := NewIntSliceFunctor(intSlice)
	fmt.Printf("original functor: %+v\n", f)

	mapperFunc1 := func(i int) int {
		return i + 10
	}
	mapped1 := f.Fmap(mapperFunc1)
	fmt.Printf("mapped functor1: %+v\n", mapped1)

	mapperFunc2 := func(i int) int {
		return i * 3
	}
	mapped2 := mapped1.Fmap(mapperFunc2)
	fmt.Printf("mapped functor2: %+v\n", mapped2)
	fmt.Printf("origin functor: %+v\n", f) // 原函子没有改变
	fmt.Printf("composite functor: %+v\n", f.Fmap(mapperFunc1).Fmap(mapperFunc2))
}
