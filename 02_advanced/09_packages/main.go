/*
是什么：包是组织和复用代码的单位；模块是由 go.mod 声明路径和版本规则的一组包。
作用：import 按路径引入包；首字母大写的名称可以从其他包访问，小写名称用于包内部实现。
使用场景：把请求处理、业务逻辑和辅助函数放到合适的包中，避免所有代码堆在 main，同时控制公开能力。

目标：模块路径、包导入、导出名字、同一模块内的代码复用。
对照：go.mod 类似项目身份与依赖声明；导入的是包目录，不是 JS 的某个源文件。
运行：go run ./02_advanced/09_packages（从仓库根目录运行整个包）。
观察：Hello, Ada! 与 Hello, Go learner!；name 包中的小写 fallback 不能从这里访问。
练习：在 name 包新增导出函数，再从 main 调用；执行 go doc 查看包的公开 API。
*/
package main

import (
	"fmt"

	"example.com/golearn/02_advanced/09_packages/name"
)

func main() {
	fmt.Println(name.Greeting(" Ada "))
	fmt.Println(name.Greeting("  "))
}
