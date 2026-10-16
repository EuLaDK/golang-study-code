package main

import (
	"fmt"
	"os"
)

// error 类型
// type error interface {
// 	Error() string
// }

// err := errors.New("your first demo error")
// errWithCtx = fmt.Errorf("index %d is out of bounds", i)

func readFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func processFile(filename string) error {
	data, err := readFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %s, err: %s", filename, err)
	}
	fmt.Println(string(data))
	return nil
}

// 异常处理 panic
// panic主要有两类来源：一类是Go运行时触发的，另一类则是开发者通过调用panic函数主动触发的。无论是哪一类，一旦panic函数被触发，后续程序的执行过程是相同的，这个过程在Go中被称为panicking
// Go官方文档以手动调用panic函数为例解释panicking过程，当函数F调用panic函数后，函数F的执行将立即停止。不过，函数F中已进行求值的deferred函数仍然会正常执行，这些deferred函数执行完毕后，控制权才会返回其调用者

// 以下代码直观地展示了 panicking
func foo() {
	println("call foo")
	bar()
	println("exit foo") // 这行不会被执行
}

func bar() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover the panic:", e)
		}
	}()
	println("call bar")
	panic("panic occurs in bar")
	zoo()               // 这行不会被执行
	println("exit bar") // 这行不会被执行
}

func zoo() {
	println("call zoo")
	println("exit zoo")
}

func main() {
	err := processFile("demo.txt")
	if err != nil {
		fmt.Println(err)
	}

	// 异常处理
	println("call main")
	foo()
	println("exit main") // 这行不会被执行
}
