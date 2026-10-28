package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// 修剪
	// TrimSpace函数去除输入字符串/字节切片首部和尾部的空白字符
	s := "     I love Go!!       "
	fmt.Println(s)
	fmt.Println(strings.TrimSpace(s))
	fmt.Printf("%q\n", bytes.TrimSpace([]byte(s)))

	// 场景 1：最常见的用法 —— 清理两端的特殊符号/空白/换行
	str := "!!! Hello, World !!!"
	// 清理两端的感叹号和空格（cutset 包含 '!' 和 ' '）
	fmt.Printf("%q\n", strings.Trim(str, "! "))      // 输出: "Hello, World"
	fmt.Printf("%q\n", strings.TrimLeft(str, "! "))  // 输出: "Hello, World !!!"
	fmt.Printf("%q\n", strings.TrimRight(str, "! ")) // 输出: "!!! Hello, World"

	// 场景 2：清理 URL 或路径两端多余的斜杠 '/'
	path := "///api/v1/users///"
	fmt.Println(strings.Trim(path, "/"))      // 输出: "api/v1/users"
	fmt.Println(strings.TrimLeft(path, "/"))  // 输出: "api/v1/users///"
	fmt.Println(strings.TrimRight(path, "/")) // 输出: "///api/v1/users"

	// 场景 3：对比 TrimLeft 与 TrimPrefix 的区别
	url := "https://example.com"

	// 错误做法（按字符集匹配，可能过度误删）：
	fmt.Println(strings.TrimLeft(url, "https://")) // 输出: "example.com" (运气好对了，但遇到 static 开头就翻车)

	// 正确做法（精确按整个前缀匹配）：
	fmt.Println(strings.TrimPrefix(url, "https://")) // 输出: "example.com"
	fmt.Println(strings.TrimSuffix(url, ".com"))

	fmt.Println("字符变换")
	// strings和bytes提供了ToUpper和ToLower函数
	fmt.Printf("%q\n", strings.ToUpper("i love golaNg"))
	fmt.Printf("%q\n", strings.ToLower("i LoVe gOlaNg"))

	fmt.Printf("%q\n", bytes.ToUpper([]byte("i LoVe gOlaNg!!")))
	fmt.Printf("%q\n", bytes.ToLower([]byte("i LoVe gOlaNg!!")))

	// Go标准库在strings和bytes包中提供了Map函数
	trans := func(r rune) rune {
		switch {
		case r == 'p':
			return 'g'
		case r == 'y':
			return 'o'
		case r == 't':
			return 'l'
		case r == 'h':
			return 'a'
		case r == 'o':
			return 'n'
		case r == 'n':
			return 'g'
		}
		return r
	}
	fmt.Printf("%q\n", strings.Map(trans, "I like python!!"))
	fmt.Printf("%q\n", bytes.Map(trans, []byte("I like python!!")))

}
