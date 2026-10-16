/*
是什么：string 是不可原地修改的字节序列；byte 表示一个字节，rune 是 int32 的别名，常用来表示 Unicode 码点。
作用：按字节处理编码数据，或按 rune 处理文字；UTF-8 用一个或多个字节编码一个码点，所以中文的字节数与码点数不同。
使用场景：处理姓名、标题和文本；按 rune 修改或计数可以避免把多字节文字从中间切开，但组合字符还需特别处理。

目标：字符串、byte、rune 与 UTF-8。
对照：JS 字符串通常按 UTF-16 代码单元计数；Go 字符串 len 按字节计数，range 解码 rune。
运行：go run ./01_basics/10_strings
观察："Go你好" 有 8 字节、4 个 rune；range 的索引为字节位置 0、1、2、5。
练习：换成含 emoji 的字符串；注意组合字符可能由多个 rune 组成，rune 数不等于视觉字符数。
*/
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	text := "Go你好"
	fmt.Println("bytes:", len(text), "runes:", utf8.RuneCountInString(text))
	for byteIndex, character := range text {
		fmt.Printf("index=%d rune=%c code=%U\n", byteIndex, character, character)
	}
	characters := []rune(text)
	fmt.Println(characters[2])
	characters[2] = '您' // 字符串不可原地修改；转换为 rune 切片后再构造新字符串。
	fmt.Println("original:", text, "changed:", string(characters))
	fmt.Printf("first byte=%d, first rune=%c\n", text[0], characters[0])
	fmt.Println(`raw string keeps \n literally`)
}
