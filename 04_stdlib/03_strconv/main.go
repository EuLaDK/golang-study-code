// 是什么：strconv 在字符串与数字、布尔值之间转换；解析是把文本读成有类型的值，格式化是反过来生成文本。
// 作用：Atoi/Itoa 转换整数，ParseInt/ParseBool 解析指定类型，FormatFloat 格式化浮点数；解析失败会返回 error。
// 使用场景：读取 HTTP 查询参数、路径 ID、端口或环境变量；解析成功后仍需检查业务允许的范围。
//
// 目标：显式转换字符串与数字/布尔值，并检查解析失败。
// 运行（仓库根目录）：go run ./04_stdlib/03_strconv
// 观察：Atoi("42px") 报错；它不像 JS parseInt 那样接受数字前缀。
// 练习：把端口字符串解析后限制在 1..65535，而非只判断能否转为 int。
package main

import (
	"fmt"
	"strconv"
)

func main() {
	n, err := strconv.Atoi("42")
	if err != nil {
		panic(err)
	}
	fmt.Println(n+1, strconv.Itoa(n))
	if _, err := strconv.Atoi("42px"); err != nil {
		fmt.Println("invalid integer:", err)
	}
	// bitSize 限制数值范围；接口 ID 通常明确使用 int64。
	id, err := strconv.ParseInt("9223372036854775807", 10, 64)
	if err != nil {
		panic(err)
	}
	active, err := strconv.ParseBool("true")
	if err != nil {
		panic(err)
	}
	fmt.Println(id, active, strconv.FormatFloat(3.14159, 'f', 2, 64))
	// 本课固定字面量出错意味着示例写错，因此 panic；外部输入应返回 error。
}
