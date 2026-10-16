package main

import (
	"fmt"
	"unicode/utf8"
)

func encodeRune() {
	var r rune = 0x4E2D
	fmt.Printf("the unicode charactor is %c\n", r)
	buf := make([]byte, 3)
	_ = utf8.EncodeRune(buf, r) // 对rune进行UTF-8编码
	fmt.Printf("utf-8 representation is 0x%X\n", buf)
}

func decodeRune() {
	var buf = []byte{0xE4, 0xB8, 0xAD}
	r, _ := utf8.DecodeRune(buf) // 对buf进行UTF-8解码
	fmt.Printf("the unicode charactor after decoding %s\n", string(r))
}

func main() {
	a := 10        // 10进制
	b := 0700      // 8进制
	c1 := 0xaabbcc // 16进制

	// 字符串
	fmt.Println(a, b, c1)
	s1 := "hello"
	// 可以换行
	s2 := `hello
	world`
	fmt.Println(s1, s2)
	encodeRune()
	decodeRune()
}
