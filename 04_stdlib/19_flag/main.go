// 是什么：命令行选项是启动命令里的 -name 等参数，flag 负责解析它们；环境变量是操作系统传给进程的配置值。
// 作用：让使用者在启动时设置端口、名称和模式，无需修改代码；flag.Parse 把传入参数写入相应变量。
// 使用场景：编写可配置的命令行工具或服务；为选项提供默认值和帮助说明，再检查参数是否在允许范围内。
//
// 目标：标准 flag 解析命令行选项，os.LookupEnv 读取环境配置。
// 运行（仓库根目录）：go run ./04_stdlib/19_flag -name Ada -port 9000 file.txt
// 观察：flag.Parse 后指针才是用户传入的值；Args 获取剩余位置参数。
// 练习：增加 -debug 布尔开关，用 -h 查看自动生成的帮助。
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.String("name", "learner", "display name")
	port := flag.Int("port", 8080, "port between 1 and 65535")
	flag.Parse()
	if *port < 1 || *port > 65535 {
		fmt.Fprintln(os.Stderr, "port must be between 1 and 65535")
		os.Exit(2)
	}
	mode, exists := os.LookupEnv("GOLEARN_MODE")
	if !exists {
		mode = "development"
	}
	fmt.Printf("name=%s port=%d mode=%s args=%q\n", *name, *port, mode, flag.Args())
	// LookupEnv 区分“未设置”和“已设置为空”；Getenv 不能区分这两种情况。
	// 标准 flag 在第一个位置参数处停止解析，所以选项应放在 file.txt 之前。
}
