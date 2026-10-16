// 是什么：log/slog 是结构化日志工具，把一条事件记录为消息、级别和有名字的字段。
// 作用：Info、Warn、Error 表达严重程度，With 为后续记录附上公共字段；JSONHandler 输出便于机器检索的 JSON。
// 使用场景：按 request_id 查请求、统计响应状态或定位失败；日志级别只影响记录方式，Error 不会自动停止程序。
//
// 目标：使用 log/slog 输出带级别和结构化字段的日志。
// 运行（仓库根目录）：go run ./04_stdlib/18_slog
// 观察：JSON 日志保留字段类型；默认 Info 级别不输出 Debug；时间每次不同。
// 练习：设置 HandlerOptions.Level = slog.LevelDebug，再运行比较输出。
package main

import (
	"errors"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	requestLog := logger.With("service", "learning-api", "request_id", "demo-001")
	requestLog.Debug("this line is filtered")
	requestLog.Info("request completed", "method", "GET", "status", 200, "duration_ms", 12)
	requestLog.Warn("retry scheduled", slog.Group("retry", "attempt", 2, "max", 3))
	requestLog.Error("operation failed", "error", errors.New("demo failure"))
	// log.Print 适合简单日志；slog 让后端平台可以按字段过滤，而非解析拼接字符串。
	// 密码、访问令牌和敏感请求体不应出现在日志中；Error 级别不会让进程退出。
}
