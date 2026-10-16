// 是什么：配置是随部署环境变化的参数，例如端口和请求超时；环境变量是启动进程时传入的一种配置来源。
// 作用：把“运行在哪里、等待多久”从业务代码中独立出来，设置默认值，并在启动时拒绝无效设置。
// 使用场景：同一份程序在开发、测试和正式环境使用不同参数，避免每次修改代码后重新构建。
//
// 学习目标：读取环境变量，并在启动时发现无效配置。
// 运行：go run ./05_backend/01_configuration
// 默认输出：port=8080 request_timeout=3s；APP_PORT / APP_TIMEOUT 可覆盖默认值。
// 对照 JS：os.LookupEnv 类似读取 process.env，但还会报告变量是否存在。
// 练习：增加 APP_DEBUG，使用 strconv.ParseBool 校验它。
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type config struct {
	Port           int
	RequestTimeout time.Duration
}

func loadConfig() (config, error) {
	cfg := config{Port: 8080, RequestTimeout: 3 * time.Second}
	if value, ok := os.LookupEnv("APP_PORT"); ok {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return config{}, fmt.Errorf("APP_PORT must be an integer from 1 to 65535")
		}
		cfg.Port = port
	}
	if value, ok := os.LookupEnv("APP_TIMEOUT"); ok {
		timeout, err := time.ParseDuration(value)
		if err != nil || timeout <= 0 {
			return config{}, fmt.Errorf("APP_TIMEOUT must be a positive duration, such as 3s")
		}
		cfg.RequestTimeout = timeout
	}
	return cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("port=%d request_timeout=%s\n", cfg.Port, cfg.RequestTimeout)
}
