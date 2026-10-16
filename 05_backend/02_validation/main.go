// 是什么：输入校验是判断外部数据是否满足格式和业务要求；清洗是按约定去掉首尾空白等无用内容。
// 作用：把任意原始输入转换成可放心使用的业务值，例如非空、长度在上限内的任务标题。
// 使用场景：处理 API 请求、文件导入或配置；服务端需要独立校验，不能假设前端已经拦截所有非法输入。
//
// 学习目标：将外部输入转换为经过校验的业务值。
// 运行：go run ./05_backend/02_validation
// 观察：空白标题被拒绝，首尾空白被移除，中文按 Unicode 码点计数。
// 对照 JS：len(string) 是字节数，不能直接替代字符数；前端校验不能替代后端校验。
// 练习：增加只允许未来日期的 deadline 字段。
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateTitle(raw string) (string, error) {
	if !utf8.ValidString(raw) {
		return "", fmt.Errorf("title must be valid UTF-8")
	}
	title := strings.TrimSpace(raw)
	// 码点不等同于用户看到的字形；某些 emoji 由多个码点组成。
	if count := utf8.RuneCountInString(title); count < 1 || count > 100 {
		return "", fmt.Errorf("title must contain 1 to 100 Unicode code points")
	}
	return title, nil
}

func main() {
	for _, raw := range []string{"  Learn Go  ", "学习 Go", "   ", strings.Repeat("a", 101)} {
		title, err := validateTitle(raw)
		if err != nil {
			fmt.Println("invalid:", err)
			continue
		}
		fmt.Printf("valid: %q, bytes=%d, code_points=%d\n", title, len(title), utf8.RuneCountInString(title))
	}
}
