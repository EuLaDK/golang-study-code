package main

import (
	"fmt"
)

// switch initStmt; expr {
// case expr1:
// // 执行分支 1
// case expr2:
// // 执行分支 2
// case expr3_1, expr3_2, expr3_3:
// // 执行分支 3
// case expr4:
// // 执行分支 4 …
// case exprN:
// // 执行分支 N
// default:
// // 执行默认分支
// }

type person struct {
	name string
	age  int
}

// 多个值符合可以这样
func checkWorkday(a int) {
	switch a {
	case 1, 2, 3, 4, 5, 6:
		fmt.Println("It is a work day")
	case 7, 8:
		fmt.Println("It is a Weekend day")
	default:
		fmt.Println("哈哈哈哈哈")
	}
}

func case1() int {
	fmt.Println("eval case1 expr")
	return 1
}

func case2() int {
	fmt.Println("eval case2 expr")
	return 2
}

func switchexpr() int {
	fmt.Println("eval switch expr")
	return 1
}

func main() {
	// 一般不会用来做结构体的这种，只是可以这么做
	p := person{"tom", 12}
	switch p {
	case person{"tony", 33}:
		fmt.Println("match tony")
	case person{"tom", 12}:
		fmt.Println("match tom")
	default:
		fmt.Println("no match")
	}
	checkWorkday(10)

	switch switchexpr() {
	case case1():
		fmt.Println("exex case1")
		fallthrough // 一般来说如果符合要求就会退出整个case，但是加上了fallthrough关键字就会继续执行后面的case语句
		// 需要注意的是，由于 fallthrough 的存在，case2 语句的表达式不会被重新求值，而是直接执行其分支代码
	case case2():
		fmt.Println("exec case2")
		fallthrough
	default:
		fmt.Println("exec default")
	}

	// type switch
	var x interface{} = 13
	switch v := x.(type) { // v 存储的是变量 x 的动态类型对应的值信息
	case nil:
		println("v is nil")
	case int:
		println("the type of v is int, v ={", v)
	case string:
		println("the type of v is string, v =", v)
	case bool:
		println("the type of v is bool, v =", v)
	default:
		println("don't support the type")
	}
}
