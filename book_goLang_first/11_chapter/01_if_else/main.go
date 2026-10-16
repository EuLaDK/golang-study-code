package main

import "fmt"

// if boolean_expression {
// 	// if 分支
// } else if boolean_expression {} else {}
// // 原分支

func var_foo() (int, int) {
	a := 10
	b := 20
	return a, b
}

func main() {
	a := 60
	if a >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}
	b := 60
	if (a >= 60) && (b >= 60) {
		fmt.Println("期末考试全部及格")
	} else {
		fmt.Println("期末考试没有及格")
	}

	if c, d := var_foo(); c > d {
		fmt.Println(c)
	} else {
		fmt.Println(d)
	}
}
