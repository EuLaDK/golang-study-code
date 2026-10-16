/*
是什么：运算符是对值进行计算或判断的符号，如 +、/、==、&&；位运算符处理整数的二进制位。
作用：算术运算计算数值，比较运算得到 bool，&& 和 || 组合条件并进行短路求值，即结果已确定时跳过右边。
使用场景：计算分页偏移、判断权限或检查输入；先判断切片非空再访问元素，避免索引越界。

目标：算术、比较、布尔与位运算。
对照：Go 没有 ===；不同数字类型通常不能直接运算；++ 是语句而非表达式。
运行：go run ./01_basics/04_operators
观察：7/2=3，余数=1，转为浮点后=3.5；短路求值不会访问空切片的第一个元素。
练习：把 a 改为 -7，观察整数除法向零截断，比较 Python 的 //。
*/
package main

import "fmt"

func main() {
	a, b := -7, 2
	fmt.Println("integer:", a/b, "remainder:", a%b, "float:", float64(a)/float64(b))
	a++
	fmt.Println("after increment:", a, "equal:", a == 8, "not equal:", a != b)

	items := []int{}
	fmt.Println("first is positive:", len(items) > 0 && items[0] > 0)
	fmt.Println("boolean:", !(a < b) || b == 0)
	fmt.Println("bits:", 6&3, 6|3, 6^3, 1<<3)
}
