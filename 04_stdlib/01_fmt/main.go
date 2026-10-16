// 是什么：fmt 是标准库的文本格式化工具，把值变成便于阅读的文字，或组合成新的字符串。
// 作用：Println 输出一行，Printf 按 %d、%s 等占位符输出，Sprintf 只生成字符串，Errorf 生成错误值。
// 使用场景：查看变量、制作命令行输出、给错误补充信息；生产服务中需要按字段检索日志时可使用 slog。
//
// 目标：用 fmt 格式化输出；%v 看值，%T 看类型，%q 看带引号的字符串。
// 运行（仓库根目录）：go run ./04_stdlib/01_fmt
// 观察：Printf 直接输出；Sprintf 返回字符串，类似 JS 模板字符串的用途。
// 练习：增加一个浮点数字段，用 %.2f 保留两位小数。
package main

import "fmt"

func main() {
	name, age := "Go learner", 24
	fmt.Printf("name=%q age=%d type=%T\n", name, age, age)
	fmt.Printf("default=%v, Go syntax=%#v\n", []string{"Go", "HTTP"}, []string{"Go", "HTTP"})
	message := fmt.Sprintf("Hello, %s. Price: %.2f", name, 19.9)
	fmt.Println(message)
	// %w 专用于包装 error；普通文本格式化用 %v 或 %s。
	err := fmt.Errorf("invalid age: %d", -1)
	fmt.Println(err)
}
