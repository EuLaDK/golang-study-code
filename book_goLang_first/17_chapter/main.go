package main

import (
	"fmt"
	"sync"
)

// 泛型函数 [T ordered]泛型的类型参数列表 包含一个类型参数T，而ordered对类型参数T的类型约束
// func maxGenerics[T ordered] (sl []T) T {
// 	// 测试
// }
// Go 规范规定：函数的类型参数列表位于函数名于函数参数列表之间，由中括号括起的固定数量的、由逗号分隔的类型参数声明组成
// func genericsFunc[T1 constraint1, T2, constraint2, ...., Tn constraintN ] (ordinary parameters list) (return values list)

// 泛型类型 像这样在定义中包含类型参数的类型被称为泛型类型
// type maxableSlice[T ordered] struct {
// 	elems []T
// }

// 在泛型类型中，类型参数列表位于类型名字后面的中括号中。泛型类型可以包含多个类型参数，这些类型参数名通常为首字母大写，并且必须是具名的，命名唯一
// type TypeName[T1 constraint1, T2 constraint2, ..., Tn constraintN] TypeLiteral
// type Set[T comparable] map[T]struct{}
// type sliceFn[T any] struct {
// s []T
// cmp func(T, T) bool
// }
// type Map[K, V any] struct {
// root *node[K, V]
// compare func(K, K) int
// }
// type element[T any] struct {
// next *element[T]
// val T
// }
// type Numeric interface {
// ~int | ~int8 | ~int16 | ~int32 | ~int64 |
// ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
// ~float32 | ~float64 |
// ~complex64 | ~complex128
// }
// type NumericAbs[T Numeric] interface {
// Abs() T
// }

// 泛型类型与类型嵌入
type Lockable[T any] struct {
	t T
	sync.Mutex
}

func (l *Lockable[T]) Get() T {
	l.Lock()
	defer l.Unlock()
	return l.t
}

func (l *Lockable[T]) Set(v T) {
	l.Lock()
	defer l.Unlock()
	l.t = v
}

type Slice[T any] []T

func (s Slice[T]) String() string {
	if len(s) == 0 {
		return ""
	}
	var result = fmt.Sprintf("%v", s[0])
	for _, v := range s[1:] {
		result = fmt.Sprintf("%v, %v", result, v)
	}
	return result
}

type Lockable2[T any] struct {
	t T
	Slice[int]
	sync.Mutex
}

// 泛型方法 为泛型类型定义的方法称为泛型方法 maxableSlice这个是一个类型不用管
// func (sl *maxableSlice[T]) max() T {
// 	if len(sl.elems) == 0 {
// 		panic("slice is empty")
// 	}
// 	max := sl.elems[0]
// 	for _, v := range sl.elems[1:] {
// 		if v > max {
// 			max = v
// 		}
// 	}
// 	return max
// }

// 类型约束 使用any和空接口 interface{}
// 使用any约束的形参T可以进行以下操作
// 声明变量、同类型赋值、将变量传给其他函数或从函数返回、取变量地址、转换或赋值给interface{}类型变量、用在类型断言或type switch中、作为符合类型中的元素类型、传递给预定义函数，例如new
func doSomething[T1, T2 any](t1 T1, t2 T2) T1 {
	var a T1 // 声明变量
	var b T2
	a, b = t1, t2
	_ = b
	f := func(t T1) {

	}
	f(a)    // 传给其他函数
	p := &a // 取变量地址
	_ = p
	var i interface{} = a // 转换或赋值给interface{}类型变量
	_ = i
	c := new(T1) // 传递给预定义函数
	_ = c
	f(a)
	sl := make([]T1, 0, 10)
	_ = sl
	j, ok := i.(T1)
	_ = ok
	_ = j
	switch i.(type) {
	case T1:
	case T2:
	}
	return a
}

// 自定义约束 Go泛型最终决定使用interface语法来定义约束，所有接口类型均可作为类型参数的约束
func Stringify[T fmt.Stringer](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.String())
	}
	return ret
}

type MyString string

func (s MyString) String() string {
	return string(s)
}

func main() {
	// 调用泛型函数
	// 类型参数分为类型形参和类型实参
	// 泛型函数声明：T为类型形参
	// func maxGenerics[T ordered](sl []T)T
	// 调用泛型函数：int为类型实参
	// m := maxGenerics[int]([]int{1,2,-4,-6,7,0})
	// 不能依赖函数的返回值类型来推断类型实参
	n := Lockable2[string]{
		t:     "hello",
		Slice: []int{1, 2, 3},
	}
	println(n.String())
	sl := Stringify([]MyString{"I", "love", "golang"})
	fmt.Println(sl)
}
