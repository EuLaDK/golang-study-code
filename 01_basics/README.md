# 第一阶段：基础语法

适合已经写过 JS / TS / Python、第一次学习 Go 的开发者。按编号阅读每节 `main.go` 顶部的说明，再运行、修改、重新运行。命令都从项目根目录执行，例如 `go run ./01_basics/01_hello`。

| 顺序 | 目录 | 学习重点 | 易错点 |
| --- | --- | --- | --- |
| 01 | [hello](01_hello/main.go) | package、import、main、格式化输出 | 一个目录通常是一个包，不是每个文件各运行各的 |
| 02 | [variables](02_variables/main.go) | var、短变量声明、零值、显式转换 | `:=` 只能用于函数内部；不会隐式把字符串转数字 |
| 03 | [constants](03_constants/main.go) | const、无类型常量、iota | const 不能保存运行时计算结果 |
| 04 | [operators](04_operators/main.go) | 整数除法、比较、逻辑、位运算 | 整数相除会截断；`++` 不能作为表达式 |
| 05 | [control_flow](05_control_flow/main.go) | if、switch、for、range、break/continue | 条件必须是 bool；没有 JS 的 truthy/falsy |
| 06 | [functions](06_functions/main.go) | 多返回值、可变参数、函数值 | 调用可能失败的函数，要检查 error |
| 07 | [arrays](07_arrays/main.go) | 固定长度数组、复制、比较 | `[3]int` 与 `[4]int` 是不同类型 |
| 08 | [slices](08_slices/main.go) | len/cap、append、copy、共享底层数组 | 切片赋值不会复制元素；append 可能换底层数组 |
| 09 | [maps](09_maps/main.go) | make、读取存在性、删除、遍历 | nil map 不能写；map 遍历顺序无保证 |
| 10 | [strings](10_strings/main.go) | byte、rune、UTF-8、字符串不可变 | len(string) 是字节数，不是字符数 |
| 11 | [pointers](11_pointers/main.go) | 地址、解引用、nil、传参 | Go 仍然按值传参，传指针复制的是地址 |
| 12 | [scope_and_values](12_scope_and_values/main.go) | 作用域、遮蔽、值语义 | 内层 `:=` 可能创建同名新变量 |

前置条件：Go 1.22 或更新版本。先运行 `go version`；不需要安装任何第三方库。

每节顶部都有小练习，建议先猜输出，再执行。完成后应能解释：为什么数组赋值是独立副本，而切片和 map 赋值后却能观察到共享数据。

验证：`go test ./01_basics/...`。切片课的测试固定检查共享底层数组、独立复制和容量边界；其他课可以直接通过输出观察语义。
