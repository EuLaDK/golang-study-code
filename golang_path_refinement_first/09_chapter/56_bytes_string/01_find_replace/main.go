package main

import (
	"bytes"
	"fmt"
	"strings"
)

// string包和bytes包的基本使用

func main() {
	// contains函数 查找有没有在这个字符里面 有true无false 第一个参数代表的字符串/字节切片中是否包含第二个参数代表的字符串/字节切片。
	fmt.Println(strings.Contains("Golang", "Go"))
	fmt.Println(strings.Contains("Golang", "go"))

	fmt.Println(bytes.Contains([]byte("Golang"), []byte("Go")))
	fmt.Println(bytes.Contains([]byte("Golang"), []byte("go")))

	// ContainsAny 函数 两个参数只要有任意一个字符相同就是true，比如golang和java都有字母a，输出true
	fmt.Println("ContainsAny方法")
	fmt.Println(strings.ContainsAny("Golang", "java"))
	fmt.Println(strings.ContainsAny("Golang", "Python"))
	fmt.Println(strings.ContainsAny("Golang", "C"))

	fmt.Println(bytes.ContainsAny([]byte("Golang"), "java"))
	fmt.Println(bytes.ContainsAny([]byte("Golang"), "c"))

	// ContainsRune函数 判断某一个Unicode字符（以码点形式即rune类型值传入）是否包含在第一个参数代表的字符串或字节切片中
	fmt.Println("ContainsRune函数")
	fmt.Println(strings.ContainsRune("Golang", 97))
	fmt.Println(bytes.ContainsRune([]byte("Golang"), 97))

	// HasPrefix和HasSuffix函数
	// HasPrefix函数用于判断第二个参数代表的字符串/字节切片是不是第一个参数的前缀
	// 同理，HasSuffix函数则用于判断第二个参数是不是一个参数的后缀。
	// 要注意的是在这两个函数的语义中，空字符串（""）是任何字符串的前缀和后缀，空字节切片（[]byte{}）和nil切片也是任何字节切片的前缀和后缀
	fmt.Println("HasPrefix和HasSuffix函数")
	fmt.Println(strings.HasPrefix("Golang", "Go"))                   // true
	fmt.Println(strings.HasPrefix("Golang", "Golang"))               // true
	fmt.Println(strings.HasPrefix("Golang", "lang"))                 // false
	fmt.Println(strings.HasPrefix("Golang", ""))                     // true
	fmt.Println(strings.HasPrefix("", ""))                           // true
	fmt.Println(strings.HasSuffix("Golang", "Go"))                   // false
	fmt.Println(strings.HasSuffix("Golang", "Golang"))               // true
	fmt.Println(strings.HasSuffix("Golang", "lang"))                 // true
	fmt.Println(strings.HasSuffix("Golang", ""))                     // true
	fmt.Println(strings.HasSuffix("", ""))                           // true
	fmt.Println(bytes.HasPrefix([]byte("Golang"), []byte("Go")))     // true
	fmt.Println(bytes.HasPrefix([]byte("Golang"), []byte("Golang"))) // true
	fmt.Println(bytes.HasPrefix([]byte("Golang"), []byte("lang")))   // false
	fmt.Println(bytes.HasPrefix([]byte("Golang"), []byte{}))         // true
	fmt.Println(bytes.HasPrefix([]byte("Golang"), nil))              // true
	fmt.Println(bytes.HasPrefix(nil, nil))                           // true
	fmt.Println(bytes.HasSuffix([]byte("Golang"), []byte("Go")))     // false
	fmt.Println(bytes.HasSuffix([]byte("Golang"), []byte("Golang"))) // true
	fmt.Println(bytes.HasSuffix([]byte("Golang"), []byte("lang")))   // true
	fmt.Println(bytes.HasSuffix([]byte("Golang"), []byte{}))         // true
	fmt.Println(bytes.HasSuffix([]byte("Golang"), nil))              // true
	fmt.Println(bytes.HasSuffix(nil, nil))

	// 定位查找
	// 定位相关查找函数会给出第二个参数代表的字符串/字节切片在第一个参数中第一次出现的位置（下标），如果没有找到，则返回-1。另外定位查找还有方向性，从左到右为正向定位查找（Index系列），反之为反向定位查找（LastIndex系列）
	// 	和ContainsAny只查看交集是否为空不同，IndexAny函数返回非空
	// 交集中第一个字符在第一个参数中的位置信息。另外要注意，反向查找
	// 空串或nil切片，返回的是第一个参数的长度，但作为位置（下标）信
	// 息，这个值已经越界了。
	fmt.Println(strings.Index("Learn Golang, Go!", "Go"))      // 6
	fmt.Println(strings.Index("Learn Golang, Go!", ""))        // 0
	fmt.Println(strings.Index("Learn Golang, Go!", "Java"))    // -1
	fmt.Println(strings.IndexAny("Learn Golang, Go!", "Java")) // 2
	fmt.Println(strings.IndexRune("Learn Golang, Go!", rune('a')))

	// 替换 Go标准库在strings包中提供了两种进行字符串替换的方法：Replace函数与Replacer类型。bytes包中则只提供了Replace函数用于字节切片的替换
	// 	Replace函数广泛用于简单的字符串/字节
	// 切片替换。它的最后一个参数是一个整型数，用于控制替换的次数。如
	// 果传入-1，则全部替换。而ReplaceAll函数本质上等价于最后一个参数
	// 传入-1的Replace函数
	fmt.Println(strings.Replace("I love java, java, java!!",
		"java", "go", 1)) // I love go, java, java!!
	fmt.Println(strings.Replace("I love java, java, java!!",
		"java", "go", 2)) // I love go, go, java!!
	fmt.Println(strings.Replace("I love java, java, java!!",
		"java", "go", -1)) // I love go, go, go!!
	fmt.Println(strings.Replace("math", "", "go", -1))
	// gomgoagotgohgo
	fmt.Println(strings.ReplaceAll("I love java, java, java!!", "java", "go")) // I love go, go, go!!
	replacer := strings.NewReplacer("java", "go", "python", "go")
	fmt.Println(replacer.Replace("I love java, python, go!!"))
	// I love go, go, go!!
	// 替换([]byte)
	fmt.Printf("%s\n", bytes.Replace([]byte("I love java, java, java!!"),
		[]byte("java"), []byte("go"), 1)) // I love go, java, java!!
	fmt.Printf("%s\n", bytes.Replace([]byte("I love java, java, java!!"),
		[]byte("java"), []byte("go"), 2)) // I love go, go, java!!
	fmt.Printf("%s\n", bytes.Replace([]byte("I love java, java, java!!"),
		[]byte("java"), []byte("go"), -1)) // I love go, go, go!!
	fmt.Printf("%s\n", bytes.Replace([]byte("math"), nil,
		[]byte("go"), -1)) // gomgoagotgohgo
	fmt.Printf("%s\n", bytes.ReplaceAll([]byte("I love java, java, java!!"),
		[]byte("java"), []byte("go")))
}
