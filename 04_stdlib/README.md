# 04 标准库：后端开发的日常工具

每个编号目录是一个独立的 `package main`。从仓库根目录执行表中的命令；先读 `main.go` 顶部的目标、观察点和练习，再运行并修改。示例只使用标准库，不需要安装依赖。

这里是**面向初学后端的常用标准库子集**，不是把 Go 每一个包罗列一遍。文件名和路径使用英文，中文解释写在代码旁边，便于边看边运行。多数示例会直接输出结果；遇到写死的示例数据失败时使用 `panic`，真实外部请求的失败应返回错误，由调用方决定如何响应。

| 顺序 | 主题与包 | 运行命令 | 重点 |
| --- | --- | --- | --- |
| 01 | `fmt` | `go run ./04_stdlib/01_fmt` | 格式化占位符、字符串生成 |
| 02 | `strings`、`bytes` | `go run ./04_stdlib/02_strings_bytes` | 清洗、拆分、拼接、字节与字符串 |
| 03 | `strconv` | `go run ./04_stdlib/03_strconv` | 严格解析、整数范围、错误 |
| 04 | `unicode`、`unicode/utf8` | `go run ./04_stdlib/04_unicode_utf8` | 字节、码点、UTF-8 |
| 05 | `time` | `go run ./04_stdlib/05_time` | 时间布局、时区、间隔、计时 |
| 06 | `math`、`math/rand/v2` | `go run ./04_stdlib/06_math_rand` | 数学、浮点数、可复现的随机模拟 |
| 07 | `sort`、`slices`、`cmp` | `go run ./04_stdlib/07_sort_slices_cmp` | 原地排序、浅拷贝、二分查找 |
| 08 | `maps` | `go run ./04_stdlib/08_maps` | 浅拷贝、比较、稳定遍历 |
| 09 | `regexp` | `go run ./04_stdlib/09_regexp` | 编译、匹配、分组、替换 |
| 10 | `os`、`path/filepath` | `go run ./04_stdlib/10_os_filepath` | 临时目录、文件读写、自动清理 |
| 11 | `io`、`bufio` | `go run ./04_stdlib/11_io_bufio` | Reader/Writer、扫描、缓冲与上限 |
| 12 | `io/fs`、`embed` | `go run ./04_stdlib/12_fs_embed` | 内嵌资源、统一文件系统接口 |
| 13 | `encoding/json` | `go run ./04_stdlib/13_json` | 标签、严格字段、单值和大小校验 |
| 14 | `encoding/csv` | `go run ./04_stdlib/14_csv` | 带分隔符的字段、Flush 与错误 |
| 15 | `encoding/base64`、`encoding/hex` | `go run ./04_stdlib/15_base64_hex` | 二进制文本编码，不是加密 |
| 16 | `net/url` | `go run ./04_stdlib/16_url` | 查询参数、多值、转义 |
| 17 | `net/http`、`net/http/httptest` | `go run ./04_stdlib/17_http` | 临时服务、路由、客户端超时与关闭 |
| 18 | `log/slog` | `go run ./04_stdlib/18_slog` | 结构化日志、级别、上下文字段 |
| 19 | `flag`、`os` | `go run ./04_stdlib/19_flag -name Ada -port 9000 file.txt` | 命令行参数、环境变量 |
| 20 | `text/template`、`html/template` | `go run ./04_stdlib/20_template` | 文本生成和 HTML 自动转义 |
| 21 | `crypto/sha256`、`crypto/hmac`、`crypto/rand` | `go run ./04_stdlib/21_crypto` | 摘要、消息认证、安全随机数 |
| 22 | `archive/zip`、`compress/gzip` | `go run ./04_stdlib/22_archive_compress` | 打包、压缩、完成写入、解压大小 |

## 与 JS / TS / Python 的联系

- 不依赖隐式类型转换：解析字符串先看 `strconv` 的 `error`，再做业务范围校验。
- 字符串字节数、Unicode 码点数和用户看到的字符数不同；处理中文时不要直接按字节截断。
- `slices.Sort` 和 `maps.Clone` 分别是原地排序与浅拷贝；有嵌套引用时要理解共享数据。
- `io.Reader` / `io.Writer` 类似流的最小能力约定；文件、网络连接、内存缓冲能复用同一套处理函数。
- `json.Unmarshal` 不是任意对象校验器；字段类型、是否必填、业务范围需要明确设计。
- 标准库已能编写 HTTP API，熟悉它后更容易理解框架帮你做了什么。

## 自检与边界

```powershell
go test ./04_stdlib/...
go vet ./04_stdlib/...
```

`13_json/main_test.go` 验证大整数精度，以及未知字段、多个 JSON 值、尾随垃圾、无效业务字段和超大输入的拒绝行为。`21_crypto` 的运行示例会检查原消息通过、篡改消息失败。文件示例使用临时目录后清理，HTTP 示例只访问本机临时服务并自动退出；随机值、耗时和日志时间可能每次不同。

其他章节还会讲解 `errors`、`context`、`sync`、`sync/atomic`、`reflect`、`testing` 等与语法、并发和后端实践紧密相关的包。`database/sql` 只提供数据库访问接口，真实数据库还需要驱动，本轮不接外部驱动。`net` 底层套接字、`crypto/tls` 证书管理、`encoding/xml`、`unsafe` 和 `os/exec` 可以在出现具体需求后再学；不要为了“用到标准库”把它们加入项目。
