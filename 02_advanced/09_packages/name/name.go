// 是什么：这是供 main 包导入的辅助包，不是独立程序；导出函数指首字母大写、可从其他包访问的函数。
// 作用：把问候语生成逻辑集中到 Greeting，去除名字的首尾空白，并为缺少名字的情况提供默认值。
// 使用场景：多处代码需要相同处理时复用一个小包；fallback 是小写名称，只作为包内部实现使用。
// 运行：go run ./02_advanced/09_packages，通过导入并调用 name.Greeting 使用本文件。
//
// Package name 演示同一模块内的辅助包以及导出命名规则。
package name

import "strings"

const fallback = "Go learner"

// Greeting 返回去除首尾空白后的问候语；空名字使用默认名称。
func Greeting(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		value = fallback
	}
	return "Hello, " + value + "!"
}
