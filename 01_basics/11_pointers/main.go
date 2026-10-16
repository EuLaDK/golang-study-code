/*
是什么：指针保存另一个值所在位置的地址；&value 取地址，*pointer 读取或修改指针指向的值。
作用：让函数能够修改调用者的数据；nil 指针表示没有指向任何值，解引用它会触发 panic。
使用场景：修改结构体字段、共享同一个对象，或区分“没有值”和“值为零”；Go 仍然按值传参，只是复制了地址。

目标：取地址 &、解引用 *、nil 指针与参数传递。
对照：Go 没有 JS 那种对象语法的隐式共享规则；参数始终按值复制，指针值复制后仍指向同一位置。
运行：go run ./01_basics/11_pointers
观察：by value 后仍为 10；by pointer 后为 11；未初始化指针等于 nil。
练习：写 swap(a, b *int) 交换两个整数，并先处理 nil 参数。
*/
package main

import "fmt"

func incrementValue(value int) {
	value++
}

func incrementPointer(value *int) {
	if value == nil {
		return
	}
	(*value)++
}

func swap(a, b *int) {
	if a == nil || b == nil {
		return
	}
	*a, *b = *b, *a
}

func main() {
	count := 10
	incrementValue(count)
	fmt.Println("by value:", count)
	incrementPointer(&count)
	fmt.Println("by pointer:", count)
	var missing *int
	fmt.Println("nil pointer:", missing == nil)
	incrementPointer(missing)
	pointer := new(int) // 得到指向零值 int 的指针；通常也可用局部变量配合 &。
	*pointer = 20
	fmt.Println("pointed value:", *pointer)
	var a int = 10
	var b int = 20
	fmt.Println(a, b)
	swap(&a, &b)
	fmt.Println(a == 20, b == 10)
}
