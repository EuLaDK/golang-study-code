package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	var a = []byte{'a', 'b', 'c'}
	var b = []byte{'a', 'b', 'c'}

	// if a == b {} // 报错不可以直接比较 可以使用Equal函数
	fmt.Println(bytes.Equal(a, b))
	if a != nil && b != nil {
		// 可以这么比较
		fmt.Println("slice a is not nil")
	}

	// strings和bytes包还共同提供了EqualFold函数，用于进行不区分大小写的Unicode字符的等值比较。字节切片在比较时，切片内的字节序列将被解释成字符的UTF-8编码表示后再进行比较
	fmt.Println(strings.EqualFold("Golang", "golang"))
	fmt.Println(bytes.EqualFold([]byte("Golang"), []byte("golang")))

	// string的比较直接用操作符就行 < > bytes的必须用Compare方法
	var a1 = []byte{'a', 'b', 'c'}
	var b1 = []byte{'a', 'b', 'd'}
	var c1 = []byte{}
	var d1 []byte

	fmt.Println(bytes.Compare(a1, b1))
	fmt.Println(bytes.Compare(b1, a1))
	fmt.Println(bytes.Compare(c1, d1))
	fmt.Println(bytes.Compare(c1, nil))
	fmt.Println(bytes.Compare(d1, nil))
	fmt.Println(bytes.Compare(nil, nil))
}
