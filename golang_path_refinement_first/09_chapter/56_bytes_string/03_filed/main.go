package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// Fields
	fmt.Printf("%q\n", strings.Fields("go java python"))
	fmt.Printf("%q\n", bytes.Fields([]byte("go java python")))

	// ==========================================
	// 1. 使用 strings 相关函数分割字符串
	// ==========================================

	// 【strings.Split】：全量分割，完全丢弃分隔符
	// 结果：["a" "b" "c"]（按逗号完全切开，去掉逗号）
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))

	// 结果：["a," ",c"]（以 "b" 为刀切开，"b" 被移除，两侧逗号保留）
	fmt.Printf("%q\n", strings.Split("a,b,c", "b"))

	// 结果：["G" "o" "社" "区" "欢" "迎" "你"]（分隔符为空字符串时，按 UTF-8 字符逐字拆分）
	fmt.Printf("%q\n", strings.Split("Go社区欢迎你", ""))

	// 结果：["abc"]（找不到分隔符 "de"，不做任何切分，原样返回一个元素的切片）
	fmt.Printf("%q\n", strings.Split("abc", "de"))

	// 【strings.SplitN】：限量分割，参数 n 控制最多切出的片段数
	// 结果：["a" "b,c,d"]（n=2 表示最多拆成 2 份，切一刀后，剩下的不切直接作为整体）
	fmt.Printf("%q\n", strings.SplitN("a,b,c,d", ",", 2))

	// 结果：["a" "b" "c,d"]（n=3 表示最多拆成 3 份，切两刀，最后一份是 "c,d"）
	fmt.Printf("%q\n", strings.SplitN("a,b,c,d", ",", 3))

	// 【strings.SplitAfter】：全量分割，但保留分隔符（贴在每个片段末尾）
	// 结果：["a," "b," "c," "d"]（逗号依然存在，紧跟在被切开的内容后面）
	fmt.Printf("%q\n", strings.SplitAfter("a,b,c,d", ","))

	// 【strings.SplitAfterN】：限量分割，且保留分隔符
	// 结果：["a," "b,c,d"]（只切成 2 份，并且把第一个分隔符留在第一份末尾）
	fmt.Printf("%q\n", strings.SplitAfterN("a,b,c,d", ",", 2))

	fmt.Println(strings.Repeat("-", 50))

	// ==========================================
	// 2. 使用 bytes 相关函数分割字节切片（逻辑完全相同）
	// ==========================================

	// 【bytes.Split】全量切分，去掉分隔符
	// 结果：["a" "b" "c"]
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte(",")))

	// 结果：["a," ",c"]
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte("b")))

	// 结果：逐字节拆分（注意：第二个参数传 nil 相当于按单个 byte 字节逐个拆碎）
	fmt.Printf("%q\n", bytes.Split([]byte("Go社区欢迎你"), nil))

	// 结果：["abc"]（找不到分隔符，原样返回）
	fmt.Printf("%q\n", bytes.Split([]byte("abc"), []byte("de")))

	// 【bytes.SplitN】限量切分
	// 结果：["a" "b,c,d"]（最多分成 2 份）
	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c,d"), []byte(","), 2))

	// 结果：["a" "b" "c,d"]（最多分成 3 份）
	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c,d"), []byte(","), 3))

	// 【bytes.SplitAfter】全量切分，保留分隔符
	// 结果：["a," "b," "c," "d"]
	fmt.Printf("%q\n", bytes.SplitAfter([]byte("a,b,c,d"), []byte(",")))

	// 【bytes.SplitAfterN】限量切分，保留分隔符
	// 结果：["a," "b,c,d"]
	fmt.Printf("%q\n", bytes.SplitAfterN([]byte("a,b,c,d"), []byte(","), 2))
}
