# 06 测试与 Go 工具

每个目录都有可运行的 `main.go` 和实际测试文件 `main_test.go`。测试使用标准库 `testing`，不安装测试框架。测试命令会编译 `main.go`，但不会执行它的 `main()`。

| 顺序 | 目录 | 学习重点 |
|---|---|---|
| 01 | `01_unit_table` | 测试命名、表驱动、t.Run 子测试、错误链、Example 输出断言 |
| 02 | `02_httptest` | ResponseRecorder 测 handler、NewServer 测真实 HTTP 往返 |
| 03 | `03_benchmarks` | b.N、准备数据、ns/op、分配统计 |
| 04 | `04_fuzz` | 种子、随机输入、解析及序列化往返性质 |

## 常用命令

以下命令都在项目根目录运行。PowerShell 中正则使用单引号，防止 `$` 等符号被 shell 解释。

```powershell
# 格式化当前学习章节；gofmt 接受目录并递归处理其中的 Go 文件。
gofmt -w ./06_testing

# 编译并测试所有包。-count=1 禁用成功结果缓存。
go test ./...
go test -count=1 ./06_testing/...

# 指定测试与子测试；-v 展示名字和日志。
go test -v ./06_testing/01_unit_table -run '^TestParsePort$/maximum$'

# 静态检查会发现可疑用法，但不能替代测试。
go vet ./...

# 覆盖率衡量哪些语句被执行，不证明业务断言足够。
go test ./06_testing/... -coverprofile coverage.out
go tool cover -func coverage.out
go tool cover -html coverage.out -o coverage.html

# -run '^$' 跳过普通测试；重复测量时关注波动，不比较一次运行的微小差异。
go test ./06_testing/03_benchmarks -run '^$' -bench . -benchmem -count=3

# 原生 fuzz 一次指定一个包及一个 fuzz 目标；普通 go test 只运行种子。
go test ./06_testing/04_fuzz -run '^$' -fuzz '^FuzzParseID$' -fuzztime 3s
```

`Fatal/Fatalf` 终止当前测试；`Error/Errorf` 记录失败后继续。不要在被测函数启动的工作 goroutine 中调用 `t.Fatal`，应先把结果传回测试 goroutine。`Example` 的 `Output` 注释必须与标准输出一致，示例才会通过。

Fuzz 找到失败输入后会写入该包的 `testdata/fuzz/`，保留这些文件作为回归用例。Fuzz 是有时限的探索，不能保证检查所有输入。

## 并发检测和性能分析

```powershell
# -race 需要受支持的平台、CGO 和适用的 C 工具链；Windows 通常需要 GCC。
go test -race ./03_concurrency/... ./06_testing/...

# 生成 CPU profile。-cpuprofile 也会在当前目录保留测试二进制。
go test ./06_testing/03_benchmarks -run '^$' -bench . -benchtime 2s -cpuprofile cpu.out
go tool pprof -top cpu.out
```

`go test -race` 只检测运行路径中的数据竞争；它不是“没有竞争”的证明。若 Windows 环境未安装 C 编译器，先完成普通测试，再配置工具链运行它。

`pprof` 的热点分析用于决定优化哪里；采样结果会受运行时间影响。学习本章先理解测试边界和错误分支，再根据真实服务的性能数据做优化。覆盖率、profile 和测试二进制是生成物，不需要当作学习源码提交。
