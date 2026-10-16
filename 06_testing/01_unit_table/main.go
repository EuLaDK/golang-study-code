// 是什么：单元测试用预设输入和预期结果检查一个小功能；表驱动测试把多组测试情况放进一张数据表。
// 作用：testing 按 Test 开头的函数运行检查，t.Run 为每一行创建子测试，让失败时容易看出是哪种情况。
// 使用场景：检查解析、校验或计算函数的正常值和边界值；修改实现后再次运行，及时发现行为变化。
//
// 目标：把边界明确的业务函数从 main 分离，使用表驱动测试检查输入。
// 对照 JS/Python：Go 的 testing 是标准库，不需要安装 Jest 或 pytest。
// 根目录运行：go run ./06_testing/01_unit_table
// 根目录测试：go test -v ./06_testing/01_unit_table
// 观察：示例输出 8080，测试中的每个输入拥有独立的子测试名称。
// 练习：添加空字符串和带空格输入，先写清楚“是否允许”的约定，再添加测试。
package main

import (
	"fmt"
	"strconv"
)

func parsePort(raw string) (int, error) {
	port, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("解析端口 %q: %w", raw, err)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("端口必须在 1..65535 之间，收到 %d", port)
	}
	return port, nil
}

func main() {
	port, err := parsePort("8080")
	if err != nil {
		panic(err)
	}
	fmt.Println(port)
}
