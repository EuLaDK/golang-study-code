package main

import (
	"fmt"
	"time"
)

func demo1() {
	var m = []int{1, 2, 3, 4, 5}

	for i, v := range m {
		go func() {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}()
	}
	time.Sleep(time.Second * 10)
}

func demo2() {
	var m = []int{1, 2, 3, 4, 5}

	for i, v := range m {
		go func(i, v int) {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}(i, v)
	}
	time.Sleep(time.Second * 10)
}

// 参与迭代的是a的副本
func arrayRangeExpression() {
	var a = []int{1, 2, 3, 4, 5} // 这里用值和切片是两种结果因为用值是创建一个新的数组，用切片则是指针指向原数组的值
	var r [5]int
	fmt.Println("-------------------")
	fmt.Println("a=", a)

	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}

		r[i] = v
	}
	fmt.Println("r=", r)
	fmt.Println("a=", a)
}

func main() {
	// demo1()
	// demo2()
	arrayRangeExpression()
	s1 := "中国人"
	for i, v := range s1 {
		fmt.Printf("%d %s 0x%x\n", i, string(v), v)
	}
}
