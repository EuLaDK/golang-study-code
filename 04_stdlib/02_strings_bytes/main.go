// 是什么：strings 处理 string 文本，bytes 处理 []byte 字节数据；Builder 用于逐步拼出字符串。
// 作用：去除空白、拆分、连接、查找和替换；选择工具时要区分文字内容和原始字节。
// 使用场景：清洗标签、解析简单文本配置或处理请求数据；处理中文时不要把一个字节当成一个完整字符。
//
// 目标：strings 操作不可变字符串，bytes 操作可变的 []byte。
// 运行（仓库根目录）：go run ./04_stdlib/02_strings_bytes
// 观察：Fields 合并连续空白；Split 会保留空字段；Builder 适合循环拼接。
// 练习：把逗号分隔的标签清洗为小写，再用 Join 连接。
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	raw := "  Go   Backend \n"
	fmt.Println(strings.TrimSpace(raw))
	fmt.Println(strings.Fields(raw))
	fmt.Printf("split=%q\n", strings.Split("a,,b", ","))
	fmt.Println(strings.Contains(strings.ToLower(raw), "go"))

	var builder strings.Builder
	for _, part := range []string{"learn", " ", "Go"} {
		builder.WriteString(part)
	}
	fmt.Println(builder.String())

	data := []byte("hello") // 转换后修改字节不改变原始字符串。
	data[0] = 'H'
	fmt.Println(string(bytes.ReplaceAll(data, []byte("llo"), []byte("y"))))
	// []byte 的索引是字节；处理中文字符应先理解 unicode/utf8 一课。
}
