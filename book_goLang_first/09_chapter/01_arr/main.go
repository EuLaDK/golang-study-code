package main

import (
	"fmt"
	"unsafe"
)

// 定义数组 var arr [N]T

// 如果两个数组类型的元素类型和长度都相同，那么这两个数组类型是等价的，如果有一个不同，那他们就是不同的数组类型
func foo(arr [3]int) {
	fmt.Println(arr)
}

func main() {
	var arr1 [3]int = [3]int{1, 2, 3}
	var arr2 [5]int
	fmt.Println(arr2) // 不初始化赋值就是全部会被赋值成0
	arr2 = [5]int{1, 2, 3, 4, 5}
	// var arr3 [3]string = [3]string{"a", "b", "c"}
	foo(arr1)
	fmt.Println(arr2)
	fmt.Println("数组长度", len(arr2))
	fmt.Println("数组大小", unsafe.Sizeof(arr2))
	// 用...可以代替长度
	var arr4 = [...]int{1, 2, 3, 4, 5}
	fmt.Println(arr4)
}
