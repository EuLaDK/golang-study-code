package main

import (
	"fmt"
)

// go 切片在运行时其实是一个三元组结构，其定义如下
// type slice struct {
// 	array unsafe.Pointer // 指向底层数组的指针
// 	len   int            // 切片的长度
// 	cap   int            // 切片的容量，cap永远大于len
// }
// Go 编译器会自动为新创建的切片建立一个底层数组，默认底层数组的长度以切片初始元素数量相同

func main() {
	var num = []int{1, 2, 3, 4, 5}
	// 基于切片创建新的切片
	num2 := num[2:4]
	num2[0] = 100
	fmt.Println(num2)
	fmt.Println(num)
	// append函数可以向切片中添加元素，返回一个新的切片
	num = append(num, 6)
	fmt.Println(num)

	// 使用make函数创建切片
	s1 := make([]int, 6, 10) // 其中，10为cap值，底层的数组长度，6为切片的数组长度
	s2 := make([]byte, 6)    // 如果没有指定cap参数，那么底层数组的长度cap等于len
	copy(s1, []int{1, 2, 3, 4, 5, 6})
	s1 = append(s1, 7)
	s2 = append(s2, 6)
	fmt.Println(s1, s2)

	// 使用已有的数组切边创建一个切片 array[low:high:max] 因为是基于数组创建的，所以修改切片会影响数组的值
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s3 := arr[3:7:9]
	fmt.Println(s3)
}
