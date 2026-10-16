# Go 后端学习代码

面向有 JavaScript / TypeScript / Python 基础、正在学习 Go 的开发者。
按编号逐个阅读、运行和修改，每个目录解决一个知识点。全部使用 Go 标准库，无需安装第三方依赖、数据库或容器。

## 从这里开始

在本目录打开终端（下面的命令也适用于 PowerShell）：

```powershell
go version
go run ./01_basics/01_hello
go run ./01_basics/02_variables
go test ./...
go vet ./...
```

需要 Go 1.22 或更新版本，本机编写环境为 Go 1.26.3。`go.mod` 声明语言最低版本 1.22；示例使用了 1.22 的 HTTP 路由、循环变量语义和 `math/rand/v2`。

使用 `go run ./目录` 运行整个示例包，不要把多个示例的 `main.go` 同时传给 `go run`。少数示例包含辅助包或嵌入资源，运行目录比运行单文件更可靠。

## 学习路线

| 顺序 | 内容 | 建议掌握的能力 |
| --- | --- | --- |
| 1 | [基础语法：12 节](01_basics/README.md) | 看懂 Go 类型、值传递、切片和 map，独立编写小程序 |
| 2 | [进阶语法：10 节](02_advanced/README.md) | 用结构体、接口、错误链和包组织代码；反射为选读 |
| 3 | [并发编程：8 节](03_concurrency/README.md) | 理解 goroutine、channel、共享状态和取消传播 |
| 4 | [标准库：22 节](04_stdlib/README.md) | 处理文本、文件、编码、时间、HTTP、模板和加密基础 |
| 5 | [后端实践：8 节](05_backend/README.md) | 处理真实 HTTP 请求，运行并测试一个任务 CRUD API |
| 6 | [测试与工具：4 节](06_testing/README.md) | 编写表驱动测试、HTTP 测试、基准测试和模糊测试 |

共 64 个主题。测试可以提前穿插学习：学完函数就能看单元测试，学完 HTTP 就能看 `httptest`。

每节代码头部包含学习目标、运行命令、观察重点和一个小练习。建议先猜输出，再执行，再做练习。练习会改变原来的输出或测试结果，修改时同步更新相应断言。

## 从前端经验迁移

| 熟悉的概念 | Go 中需要重新理解的地方 |
| --- | --- |
| JS 动态类型 / TS 类型注解 | Go 在编译期检查类型，类型转换通常需要显式编写 |
| JS Array | 数组的长度属于类型；日常动态列表使用切片，切片可能共享底层数组 |
| JS Object / Map | 结构体表达固定字段，map 表达键值集合；读 map 要区分零值与键不存在 |
| class / extends | 使用结构体、方法和组合；接口由方法集合隐式满足 |
| throw / try / catch | 可预期失败通常通过 `error` 返回；`panic` 不是日常业务错误机制 |
| Promise / async / await | goroutine 是并发执行单元，不等同于 Promise；需考虑同步、取消与退出 |
| JS 事件循环 | Go 的 goroutine 可以并行执行，共享内存需要锁或其他同步机制 |
| npm / pip | `go.mod` 声明模块和依赖；本仓库仅依赖标准库，因此没有 `go.sum` |
| JSON.stringify / JSON.parse | `encoding/json` 按导出字段和 struct tag 编解码，还需处理未知字段、数字与零值 |

## 运行后端示例

```powershell
go run ./05_backend/08_crud_api
```

默认监听 `127.0.0.1:8081`。在另一个 PowerShell 窗口中：

```powershell
Invoke-RestMethod http://127.0.0.1:8081/healthz
$body = @{ title = 'Learn Go'; completed = $false } | ConvertTo-Json
Invoke-RestMethod http://127.0.0.1:8081/tasks -Method Post -ContentType 'application/json' -Body $body
Invoke-RestMethod http://127.0.0.1:8081/tasks
```

按 Ctrl+C 停止。完整增删改查、错误示例和设计边界见 [后端说明](05_backend/README.md)。

## 检查代码

```powershell
go test ./...
go vet ./...
go test ./06_testing/01_unit_table -cover
go test ./06_testing/03_benchmarks -run '^$' -bench . -benchmem
```

多数示例运行后会自动退出。`05_backend/07_graceful_shutdown` 和 `05_backend/08_crud_api` 是持续运行的服务，需要 Ctrl+C 停止。HTTP 标准库和客户端示例只访问本机测试服务，无需互联网。

`go test -race ./...` 需要平台支持的 C 编译器/CGO 环境；普通 `go test` 不需要它。测试工具章节包含详细说明。

## 后续学习

这里覆盖常用后端标准库及核心语法，不是标准库每个包、每个函数的完整参考。遇到新 API，可以运行 `go doc net/http.Server` 或 `go doc strings.Cut` 阅读本机文档。

CRUD 使用进程内存存储，重启会丢数据。掌握本仓库之后，再学习 `database/sql` + 实际数据库驱动、SQL 与事务、身份认证、部署和 Web 框架。`database/sql` 自身不提供数据库驱动，本轮不引入外部依赖。
