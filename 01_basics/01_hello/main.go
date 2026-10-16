/*
是什么：package 用来声明代码属于哪个包，import 用来引入其他包；main 包里的 main 函数是可执行程序的入口。
作用：告诉 Go 从哪里开始运行；fmt 是标准库中的文本输出工具，Println 输出一行，Printf 按占位符输出。
使用场景：先用最小程序验证开发环境，再逐步把 main 扩展为命令行工具或 HTTP 服务的启动入口。

目标：认识 Go 程序入口、包与标准输出。
对照：不同于 Python 从文件第一行执行，Go 可执行程序从 main 包的 main 函数开始。
运行：go run ./01_basics/01_hello
观察：Hello, Go!；第二行显示 name=frontend learner, year=2, ready=true。
练习：增加一个字符串变量，使用 %q 输出带引号的内容。
*/
package main

import "fmt"

func main() {
	str := "daowuya"
	fmt.Println("Hello, Go!")
	fmt.Printf("name=%s, year=%d, ready=%t\n", "frontend learner", 2, true)
	fmt.Printf("type=%T, quoted=%q\n", 42, "Go\nbackend")
	fmt.Printf("name=%s", str)
}
