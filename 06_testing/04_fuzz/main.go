// 是什么：模糊测试 fuzz 从少量种子输入出发，自动生成和探索大量不同输入，寻找崩溃或违反约定的情况。
// 作用：测试“无论输入是什么都应满足”的性质，本例检查合法 ID 解析后再格式化、解析，结果仍应一致。
// 使用场景：测试解析器、编码器和外部输入边界，发现手写用例容易遗漏的组合；非法输入返回 error 本身是正常行为。
//
// 目标：用 fuzz 自动生成边界输入，验证解析器不会崩溃且支持规范化往返。
// 对照 JS/TS：类型只约束程序内部；HTTP 路径传来的字符串仍需要运行时校验。
// 根目录运行：go run ./06_testing/04_fuzz
// 普通测试（种子）：go test ./06_testing/04_fuzz
// PowerShell 模糊测试：go test ./06_testing/04_fuzz -run '^$' -fuzz '^FuzzParseID$' -fuzztime 3s
// 观察：示例输出 42；fuzz 自动尝试各种输入，不要求生成次数固定。
// 练习：增加空格、中文和 uint64 最大值加一作为种子，观察 strconv 的规则。
package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parseID(raw string) (uint64, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("解析 ID %q: %w", raw, err)
	}
	if id == 0 {
		return 0, errors.New("ID 必须大于 0")
	}
	return id, nil
}

func main() {
	id, err := parseID("42")
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}
