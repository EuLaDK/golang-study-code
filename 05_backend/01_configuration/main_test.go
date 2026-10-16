// 是什么：这是启动配置测试，通过测试专用的环境变量控制程序读取到的配置。
// 作用：t.Setenv 临时设置变量并在测试结束后恢复，验证有效端口与超时能读取，无效值被拒绝。
// 使用场景：调整环境变量名称、默认值或校验规则后，确认错误配置不会被静默接受。
// 运行：go test ./05_backend/01_configuration；环境变量是进程级状态，这类测试不能随意并行执行。
//
package main

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	t.Setenv("APP_PORT", "9000")
	t.Setenv("APP_TIMEOUT", "2s")
	cfg, err := loadConfig()
	if err != nil || cfg.Port != 9000 || cfg.RequestTimeout != 2*time.Second {
		t.Fatalf("config=%+v err=%v", cfg, err)
	}
	for _, port := range []string{"", "abc", "0", "65536"} {
		t.Setenv("APP_PORT", port)
		if _, err := loadConfig(); err == nil {
			t.Fatalf("accepted invalid port %q", port)
		}
	}
	t.Setenv("APP_PORT", "9000")
	t.Setenv("APP_TIMEOUT", "-1s")
	if _, err := loadConfig(); err == nil {
		t.Fatal("accepted negative timeout")
	}
}
