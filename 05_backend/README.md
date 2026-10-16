# 05 后端实践

前置：基础语法、结构体与接口、错误处理、锁、context，以及 JSON、HTTP 标准库示例。

| 顺序 | 目录 | 重点 |
| --- | --- | --- |
| 01 | [configuration](01_configuration/main.go) | 读取环境变量、默认值与启动校验 |
| 02 | [validation](02_validation/main.go) | 输入边界、去除空白、Unicode 字符数与业务约束 |
| 03 | [json_response](03_json_response/main.go) | JSON 响应、状态码和统一错误结构 |
| 04 | [routing](04_routing/main.go) | Go 1.22 方法路由、路径参数、404 与 405 |
| 05 | [middleware](05_middleware/main.go) | 用函数组合实现请求头、日志、请求超时 |
| 06 | [http_client](06_http_client/main.go) | 超时、context、状态码检查、响应体关闭与大小限制 |
| 07 | [graceful_shutdown](07_graceful_shutdown/main.go) | 信号、Server 配置、停止接收与等待进行中的请求 |
| 08 | [crud_api](08_crud_api/main.go) | 并发安全的完整任务增删改查 API |

前六节自动退出。最后两节启动本机 HTTP 服务，按 Ctrl+C 退出。

## 配置与关闭

```powershell
go run ./05_backend/01_configuration
$env:APP_PORT = '9000'
go run ./05_backend/01_configuration
Remove-Item Env:APP_PORT
go run ./05_backend/07_graceful_shutdown -addr 127.0.0.1:8080
```

配置章节只读取并验证配置，不会启动服务。关闭示例访问 `http://127.0.0.1:8080/slow` 会执行一个可取消的慢请求，期间 Ctrl+C 可以观察等待行为。

## 运行 CRUD API

```powershell
go run ./05_backend/08_crud_api -addr 127.0.0.1:8081
```

在另一个 PowerShell 窗口执行：

```powershell
$base = 'http://127.0.0.1:8081'
Invoke-RestMethod "$base/healthz"
$body = @{ title = 'Learn Go'; completed = $false } | ConvertTo-Json
$task = Invoke-RestMethod "$base/tasks" -Method Post -ContentType 'application/json' -Body $body
Invoke-RestMethod "$base/tasks"
Invoke-RestMethod "$base/tasks/$($task.id)"
$updated = @{ title = 'Read the tests'; completed = $true } | ConvertTo-Json
Invoke-RestMethod "$base/tasks/$($task.id)" -Method Put -ContentType 'application/json' -Body $updated
Invoke-RestMethod "$base/tasks/$($task.id)" -Method Delete
```

| 请求 | 成功状态 | 含义 |
| --- | --- | --- |
| `GET /healthz` | 200 | 进程健康检查 |
| `GET /tasks` | 200 | 按 ID 升序返回任务数组，空列表为 `[]` |
| `POST /tasks` | 201 | 创建任务；响应包含 `Location` |
| `GET /tasks/{id}` | 200 | 读取单个任务 |
| `PUT /tasks/{id}` | 200 | 全量替换 title 和 completed；completed 省略时为 false |
| `DELETE /tasks/{id}` | 204 | 删除任务，无响应体 |

写入请求必须使用 `Content-Type: application/json`，title 去除首尾空白后为 1～100 个 Unicode 码点。最大请求体为 1 MiB；拒绝未知字段和一个请求体中包含多个 JSON 值。客户端可通过错误响应 `{"error":"..."}` 查看原因。

常见失败：400 输入无效，404 资源不存在，405 方法不支持，413 请求体过大，415 内容类型不支持。Go 标准路由器产生的未知路径 404 和方法错误 405 使用它自带的文本格式，应用处理的错误使用 JSON。

```powershell
go test ./05_backend/...
go test ./05_backend/08_crud_api -v
```

## 阅读顺序与边界

先读 `newHandler` 的路由表，再读各个 handler、`readTaskInput` 输入边界，最后读启动与关闭流程；对照测试理解失败路径。

示例用 `sync.RWMutex` 保护进程内存 map，在锁内完成每次读写及 ID 分配。为便于阅读，数据层和 handler 放在一个文件中；测试放在相邻文件。服务仅默认监听本机，数据不持久化，不含身份认证。

这里的哈希、API 与中间件示例都是学习材料。真实对外服务还要结合需求增加认证授权、HTTPS、持久化、分页和监控。Go 的 CORS 不是服务端身份认证；前端跨域访问时再根据明确的允许源配置它。

练习：为任务添加截止时间，并同时更新 JSON 输入校验、存储、返回值和测试。注意多个修改是否都在正确的锁范围内。
