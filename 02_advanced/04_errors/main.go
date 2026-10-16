/*
是什么：error 是 Go 表达可预期失败的接口；哨兵错误是预先定义、可作为固定错误标识的值。
作用：用 %w 给错误补充上下文并保留原始原因；errors.Is 判断错误链中的标识，errors.As 提取其中某种错误类型。
使用场景：解析参数、读文件、访问数据库或上游接口失败时返回 error，由调用方记录日志、重试或转换为 HTTP 响应。

目标：错误返回、哨兵错误、使用 %w 包装、errors.Is/As 解包。
对照：正常失败通过返回值传递，不用 JS/Python 的 try/catch；错误也是有类型的值。
运行：go run ./02_advanced/04_errors
观察："42" 得到 ID 42；"0" 匹配 ErrInvalidID；"oops" 可提取 *strconv.NumError。
练习：为 ID 增加上限，返回可用 errors.Is 判断的新哨兵错误，并增加测试用例。
*/
package main

import (
	"errors"
	"fmt"
	"strconv"
)

// 1. 定义原有的错误和新增加的“上限错误”
var (
	ErrInvalidID  = errors.New("id must be positive")
	ErrIDTooLarge = errors.New("id exceeds maximum limit") // <--- 1. 新哨兵错误：上限错误
)

const MaxID = 1000 // 定义 ID 的上限为 1000

func parseID(text string) (int, error) {
	id, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("parse id %q: %w", text, err)
	}

	// 下限检查
	if id <= 0 {
		return 0, fmt.Errorf("validate id %d: %w", id, ErrInvalidID)
	}

	// 2. 上限检查
	if id > MaxID { // <--- 2. 检查是否超过上限
		return 0, fmt.Errorf("validate id %d: %w", id, ErrIDTooLarge)
	}

	return id, nil
}

func main() {
	// 4. 增加测试用例 "9999"（用来触发超过上限的错误）
	for _, text := range []string{"42", "0", "oops", "9999"} { // <--- 4. 增加了 "9999"
		id, err := parseID(text)
		if err == nil {
			fmt.Println("id:", id)
			continue
		}

		fmt.Println("error:", err)

		// 判断是否是下限错误
		if errors.Is(err, ErrInvalidID) {
			fmt.Println("invalid id: ask caller for a positive number")
		}

		// 3. 用 errors.Is 判断是否是上限错误
		if errors.Is(err, ErrIDTooLarge) { // <--- 3. 捕获并处理上限错误
			fmt.Println("invalid id: id is too large, must be <= 1000")
		}

		var numberError *strconv.NumError
		if errors.As(err, &numberError) {
			fmt.Println("invalid number:", numberError.Num)
		}
	}
}
