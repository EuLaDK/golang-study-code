# 第二阶段：进阶语法

前置：完成 `01_basics`，能解释值、指针、切片和 map。仍然只使用标准库，每节独立运行，例如 `go run ./02_advanced/01_structs_methods`。

| 顺序 | 目录 | 学习重点 | 易错点 |
| --- | --- | --- | --- |
| 01 | [structs_methods](01_structs_methods/main.go) | 结构体、方法、值接收者和指针接收者 | 修改接收者需要指针；字段大写才可导出 |
| 02 | [interfaces](02_interfaces/main.go) | 隐式实现、方法集、类型断言、类型 switch | 接口持有带类型的 nil 指针时，接口本身不等于 nil |
| 03 | [composition](03_composition/main.go) | 嵌入、方法提升、显式字段访问 | 嵌入是组合，不是类继承 |
| 04 | [errors](04_errors/main.go) | error、%w、errors.Is/As | 不要靠比较错误字符串决定业务逻辑 |
| 05 | [defer](05_defer/main.go) | 延迟调用、逆序执行、参数求值时间 | defer 在函数返回时执行，不在块结束时执行 |
| 06 | [panic_recover](06_panic_recover/main.go) | panic 栈展开、recover 的位置 | recover 只能在同一 goroutine 的延迟函数中直接生效 |
| 07 | [closures](07_closures/main.go) | 闭包状态、循环变量 | Go 1.22 的循环内新变量与循环外共享变量不同 |
| 08 | [generics](08_generics/main.go) | 类型参数、约束、comparable、~ | 泛型不是运行时的 any；约束决定可执行的操作 |
| 09 | [packages](09_packages/main.go) | module、package、导出与导入 | go.mod 定义模块路径；目录确定包的导入路径 |
| 10 | [reflection](10_reflection/main.go) | Type、Value、Kind、字段标签 | 能用具体类型或接口解决时，优先不用反射 |

运行模块课时用 `go run ./02_advanced/09_packages`，不要把辅助包当程序启动。它是本阶段唯一为教学目的拆成多个包的示例。

验证：`go test ./02_advanced/...`。错误课的表驱动测试覆盖成功值、哨兵错误和被包装的解析错误；泛型课会在输出中展示内建类型与自定义类型复用同一函数。

完成后再进入并发编程和后端实践。特别留意：正常失败返回 error，资源清理使用 defer，panic/recover 用于不可预期失败的边界处理，不能代替常规参数校验。
