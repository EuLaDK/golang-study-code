package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println("当前次数", i)
	}

	a := 10
	for a < 20 { // 这个有点类似于别的语言的while循环
		a++
	}
	fmt.Println(a)

	var s1 []int = []int{1, 2, 3, 4, 5}
	for i, v := range s1 {
		fmt.Println(i, v)
	}

	// 迭代字符类型
	var s = " 中国人 "
	// 对于 string 类型，每次循环得到的 v 值是一个 Unicode 字符码点（即 rune 类型），而不是字节。返回的第一个值 i 为该 Unicode 字符码点（UTF-8 编码）的第一个字节在字符串内存序列中的位置
	for i, v := range s {
		fmt.Printf("%d %s 0x%x\n", i, string(v), v)
	}
	// 迭代map类型
	var m = map[string]string{
		"name": "tom",
		"addr": "重庆市开州区",
	}
	for k, v := range m {
		fmt.Println(k, v)
	}

	// continue 和 带label 的continue
	s2 := []int{1, 2, 3, 4, 5, 6, 7}
	value := 0
	for i, v := range s2 {
		if v%2 == 0 {
			continue
		}
		value += s2[i]
	}
	fmt.Println(value)

	// 带label的，这个label标记的是for循环的起始位置 可以使用 continue loop 的方式来实现循环体中断
	// 当处理复杂的嵌套循环结构时，带label的continue语句常用于跳转到外层循环
	var s3 = [][]int{
		{1, 34, 26, 35, 78},
		{3, 45, 13, 24, 99},
		{101, 13, 38, 7, 127},
		{54, 27, 40, 83, 81},
	}
outerloop:
	for i := 0; i < len(s3); i++ {
		for j := 0; j < len(s3[i]); j++ {
			if s3[i][j] == 13 {
				fmt.Printf("found 13 at [%d, %d]\n", i, j)
				continue outerloop
			}
		}
	}

	// break 语句
	var s4 = []int{5, 19, 6, 3, 8, 12}
	// 找出整型切片s4种的第一个偶数
	for i := 0; i < len(s4); i++ {
		if s4[i]%2 == 0 {
			fmt.Println("第一个偶数是", s4[i])
			break
		}
	}

	// break 也有和coninute类型的 label 也是跳出嵌套循环，不然break只是跳出当前循环

outerbreakloop:
	for i := 0; i < len(s3); i++ {
		for j := 0; j < len(s3[i]); j++ {
			if s3[i][j] == 38 {
				fmt.Printf("found gold at [%d, %d]\n", i, j)
				break outerbreakloop
			}
		}
	}

}
