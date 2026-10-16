// 是什么：Unicode 给文字分配码点编号，UTF-8 是把码点编码成字节的规则；rune 常用来表示一个码点。
// 作用：unicode 判断字母等文字类别，unicode/utf8 校验编码、计数和解码；字节数、码点数与视觉字符数可能不同。
// 使用场景：校验多语言标题、处理中文输入和检查文本编码；rune 计数仍不等于完整的用户可见字符计数。
//
// 目标：区别 byte、rune（Unicode 码点）和用户眼中的字符。
// 运行（仓库根目录）：go run ./04_stdlib/04_unicode_utf8
// 观察："Go中文" 是 8 字节、4 个码点；range 的索引仍是字节偏移。
// 练习：尝试 "e\u0301"，解释为什么码点个数不等于可见字符个数。
package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	s := "Go中文"
	fmt.Println("bytes:", len(s), "runes:", utf8.RuneCountInString(s))
	for offset, r := range s {
		fmt.Printf("offset=%d rune=%c codepoint=%U letter=%t\n", offset, r, r, unicode.IsLetter(r))
	}
	fmt.Println("first Chinese rune:", string([]rune(s)[2]))
	fmt.Println("valid UTF-8:", utf8.ValidString(string([]byte{0xff})))
	// Go string 可装任意字节，不保证 UTF-8；JS length 则通常统计 UTF-16 码元。
	// 标准库不提供完整的字素簇分割；按 []rune 截断仍可能切开组合字符/emoji。
}
