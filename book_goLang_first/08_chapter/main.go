package main

import (
	"fmt"
)

// iota
// << 左移一位的意思
// 0000 -> 0001 = 1
// 0001 -> 0010 = 2
// 0010 -> 0100 = 4
const (
	mutexLocked = 1 << iota
	mutexWoken
	mutexStarving
	mutexWaiterShift    = iota // 这里iota是3，所以直接赋值就是3
	starvingThresholdNs = 1e6
)

const (
	Apple, Banana  = 1, iota + 10 // 这里iota是0
	Cherry, Durian                // 1
	Orange, Papaya                // 2 可以理解为TS的枚举默认第一个就是0
)

func main() {
	fmt.Println(mutexLocked, mutexWoken, mutexStarving, mutexWaiterShift, starvingThresholdNs)
	fmt.Println(Apple, Banana, Cherry, Durian, Orange, Papaya)
}
