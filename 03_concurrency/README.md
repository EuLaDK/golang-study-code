# 03 并发编程

先完成基础语法和函数、错误处理，再按顺序阅读每个目录的 `main.go`。每课顶部都有运行命令、预期现象和练习，所有示例都会自行结束。

| 顺序 | 目录 | 学习重点 |
|---|---|---|
| 01 | `01_goroutines` | goroutine 生命周期、WaitGroup、非确定输出顺序 |
| 02 | `02_channels` | 无缓冲收发、方向类型、range、关闭的所有权 |
| 03 | `03_buffered_channels` | 缓冲容量、阻塞、关闭后排空队列 |
| 04 | `04_select` | 多路等待、default、超时、nil 禁用分支 |
| 05 | `05_mutex` | 数据竞争、临界区、共享计数器 |
| 06 | `06_atomic` | 原子计数、CompareAndSwap、适用边界 |
| 07 | `07_context` | 取消、超时、请求级元数据、等待任务退出 |
| 08 | `08_worker_pool` | 限制并发、集中关闭、结果恢复输入顺序 |

在项目根目录运行：

```powershell
go run ./03_concurrency/01_goroutines
go test ./03_concurrency/...
go vet ./03_concurrency/...
```

这些示例中的 goroutine 可以并行运行，不能依赖 JS 单线程的直觉。`time.Sleep` 只让时间过去，不证明另一项任务完成；使用 channel 或 WaitGroup 建立同步关系。`cancel()` 只是通知任务，任务需要检查 `ctx.Done()`，调用者仍要在必要时等待其退出。

数据竞争检测需要平台支持和 C 工具链。Windows 上通常需要安装适用的 GCC，并启用 CGO。环境满足后运行：

```powershell
go test -race ./03_concurrency/...
go run -race ./03_concurrency/05_mutex
go run -race ./03_concurrency/06_atomic
```

`go test -race` 只检查实际执行的路径；没有测试的 `main` 不会被它运行，因此计数示例另列 `go run -race`。锁是保护状态，channel 是传递数据，context 是传播取消；按需求选择，不必把每个函数都改成 goroutine。
