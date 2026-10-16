package main

import "fmt"

// Go中，函数参数采用的是值传递方式。所谓值传递，指的是将实际参数在内存中的表示逐位赋值到形参中，对于整型、数组、结构体这类类型，它们的内存表示就是数据本身
// 因此作为实参时，值传递意味着直接复制这些数据，且传递开销也与数据大小成正比
// 像 string、切片、map这样的引用类型，其内存表示的是对数据内容的描述符，当这些类型作为实参时，值传递仅复制描述符(内存地址)而不复制实际数据内容的复制过程也被称为浅复制

// 在 Go 中，变长参数(elems ...int)实际上是通过切片来实现的
func myAppend(s1 []int, elems ...int) []int {
	fmt.Printf("%T\n", elems)
	if len(elems) == 0 {
		fmt.Println("no elems to append")
		return s1
	}
	s1 = append(s1, elems...)
	return s1
}

// 具名返回值，为函数返回值命名，可以在函数里面作为局部变量一样使用，如果具名返回值可以提高代码的可读性，则应优先考虑使用
// func parseNanoseconds(value string, nbytes int) (ns int, rangeErrString string, err error) {
// 	if !commaOrPeriod(value[0]) {
// 		err = errBad
// 		return
// 	}
// 	if ns, err =atoi(value[1:nbytes]); err != nil {
// 		return
// 	}
// }

// 在Go中，只有在函数（或方法）内部才能使用 defer
// defer 关键字后面只能接函数（或方法），这些函数被称为deferred（被延迟的）函数
// 这些函数会被注册到他们所在的goroutime中，并存放用于存储deferred函数的栈数据结构中
// 当执行defer的函数退出前，这些函数按后进先出的顺序被程序调度执行
// func doSomething() error {
// var mu sync.Mutex
// mu.Lock()
// defer mu.Unlock()
// r1, err := OpenResource1()
// if err != nil {
// return err
// }
// defer r1.Close()
// r2, err := OpenResource2()
// if err != nil {
// return err
// }
// defer r2.Close()
// r3, err := OpenResource3()
// if err != nil {
// return err
// }
// defer r3.Close()
// // 使用 r1、r2、r3
// return doWithResources()
// }

// 对于自定义的函数或方法，defer可以无条件的支持他们，但是对于有返回值的自定义函数或方法，返回值会在deferred函数调度执行时贝总丢弃
// 这些内置函数不可以被defer设置为deferred函数 append/cap/len/make/new/imag

func foo1() {
	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
}

func foo2() {
	for i := 0; i <= 3; i++ {
		defer func(n int) {
			fmt.Println(n)
		}(i)
	}
}

func foo3() {
	for i := 0; i <= 3; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

func main() {
	s1 := []int{1, 2, 3}
	s1 = myAppend(s1)
	fmt.Println((s1))
	s1 = myAppend(s1, 4, 5, 6)
	fmt.Println(s1)

	// defer 关键字的求值时机，defer关键字后面的表达式是在将deferred函数注册到deferred函数栈时进行求值的
	fmt.Println("fool result:")
	foo1()
	fmt.Println("\nfoo2 result:")
	foo2()
	fmt.Println("\nfoo3 result:")
	foo3()
}
