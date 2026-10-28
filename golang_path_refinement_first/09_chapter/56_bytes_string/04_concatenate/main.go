package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// join函数实现字符串字节切片的拼接
	s := []string{"I", "love", "Go"}
	fmt.Println(strings.Join(s, " ")) // I love Go
	b := [][]byte{[]byte("I"), []byte("love"), []byte("Go")}
	fmt.Printf("%q\n", bytes.Join(b, []byte(" ")))

	// strings 提供了Builder类型和相关方法搞笑构建字符串 bytes对应的是buffer类型及相关方法
	// 在 Go 中，string 是不可变类型。如果你用 + 在循环里拼接字符串，底层会不断分配新内存并发生内存拷贝，导致性能极差、GC 压力巨大

	// 	strings.Builder 就是为了解决这个问题而在 Go 1.10 引入的。
	// 为什么大家都爱用它？
	// 零拷贝转换（Zero-Copy）： 它的 String() 方法利用了 unsafe.Pointer，把底层的 []byte 直接当成 string 返回，不需要重新开辟内存复制一遍。
	// 支持预分配： 可以用 Grow(n) 提前开辟好内存，做到全程 0 次额外内存分配。

	ss := []string{"I", "love", "Go"}
	var builder strings.Builder
	for i, w := range ss {
		builder.WriteString(w)
		if i != len(ss)-1 {
			builder.WriteString(" ")
		}
	}
	fmt.Printf("%s\n", builder.String())

	bb := [][]byte{[]byte("I"), []byte("love"), []byte("Go")}
	var buf bytes.Buffer
	for i, w := range bb {
		buf.Write(w)
		if i != len(b)-1 {
			buf.WriteString(" ")
		}
	}
	fmt.Printf("%s\n", buf.String())
}
