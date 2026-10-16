/*
是什么：控制流程决定“执行哪段代码、重复执行多少次”；if 做条件判断，switch 选择分支，for 表示循环。
作用：range 逐个取出集合中的元素；break 提前退出循环，continue 跳过当前轮剩余语句。
使用场景：检查请求是否合法、按状态选择响应、批量处理列表；循环还可用于持续读取任务或倒计时。

目标：条件、switch、多种 for 写法、range、循环控制。
对照：条件必须是 bool；只有 for，没有 while；switch 默认不会穿透到下个 case。
运行：go run ./01_basics/05_control_flow
观察：positive、success；过滤后的循环和为 9；range 同时返回索引和值。
练习：用 switch 给分数划分等级；再用 for condition 实现倒计时。
*/
package main

import "fmt"

func main() {
	if n := 3; n > 0 { // n 只存在于这个 if/else 内。
		fmt.Println("positive")
	} else {
		fmt.Println("non-positive")
	}

	status := 201
	switch {
	case status >= 200 && status < 300:
		fmt.Println("success")
	case status >= 400:
		fmt.Println("error")
	default:
		fmt.Println("other")
	}
	score := 75
	switch {
	case score >= 60 && score < 70:
		fmt.Println("及格")
	case score >= 70 && score < 80:
		println("良好")
	case score >= 80 && score < 90:
		println("优秀")
	case score >= 90 && score < 100:
		println("非常优秀")
	case score == 100:
		println("特别优秀")
	default:
		fmt.Println("不及格")
	}

	total := 0
	for i := 1; i <= 10; i++ {
		if i > 5 {
			break
		}
		if i%2 == 0 {
			continue
		}
		total += i // 1,4,9
	}
	fmt.Println("sum:", total)
	remaining := 2
	for remaining > 0 {
		remaining--
	}
	fmt.Println("remaining:", remaining)
	for index, value := range []string{"GET", "POST"} {
		fmt.Println(index, value)
	}
	time := 10
	for time > 0 {
		println("倒计时：", time)
		time--
	}
	fmt.Println("倒计时结束！")
}
