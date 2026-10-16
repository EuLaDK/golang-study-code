package main

import (
	"fmt"
	"unsafe"
)

// 切片在Go运行时层面的内部表示 在运行时中，每个切片变量都是一个runtime.slice结构体类型的实例
type slice struct {
	array unsafe.Pointer // 指向下层数组某元素的指针，该元素也是切片的起始元素
	len   int            // 切片的长度，即切片中当前元素的个数
	cap   int            // 切片的最大容量，cap >= len
}

// 切片的高级特性 动态扩容 大致意思就是cap到达上限会扩充二倍
// 在当前底层数组容量无法满足
// 的情况下，动态分配新的数组，新数组长度会按一定算法扩展（参见
// $GOROOT/src/runtime/slice.go中的growslice函数）

func main() {
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s2 := s1[2:6]
	s3 := s1[3:6]
	fmt.Println(s1, s2, s3)
	s3[0] = 10
	fmt.Println(s1, s2, s3)
	u := []int{11, 12, 13, 14, 15}
	fmt.Println("array:", u)
	s := u[1:3]
	fmt.Printf("slice(len=%d,cap=%d):%v\n", len(s), cap(s), s)
	s = append(s, 24)
	fmt.Println("after append 24, array:", u)
	fmt.Printf("after append 24, slice(len=%d,cap=%d):%v\n", len(s), cap(s), s)
	s = append(s, 25)
	fmt.Println("after append 25, array:", u)
	fmt.Printf("after append 25, slice(len=%d,cap=%d):%v\n", len(s), cap(s), s)
	s = append(s, 26)
	fmt.Println("after append 26, array:", u)
	fmt.Printf("after append 26, slice(len=%d,cap=%d):%v\n", len(s), cap(s), s)
	s[0] = 22
	// 添加元素25后，切片的元素已经触碰到了底层数组u的边界，后面在添加26，append底层数组已经无法满足添加新元素的要求，于是新创建了一个底层数组，并将原切片的元素复制到新数组中
	// 这之后再修改切片中的元素值，原数组u的元素也没有发生任何改变，因为此时切片s与数组u已经解除了绑定关系，s已经不再是数组u的描述符了
	fmt.Println("after reassign 1st elem of slice, array:", u)
	fmt.Printf("after reassign 1st elem of slice, slice(len=%d,cap=%d):%v\n", len(s), cap(s), s)
}
