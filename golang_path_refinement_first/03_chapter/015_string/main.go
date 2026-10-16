package main

import (
	"unsafe"
)

// string在Go中的运行时
type stringStruct struct {
	str unsafe.Pointer
	len int
}

// 在能预估出最终字符串长度的情况下，使用预初始化的strings.Builder连接构建字符串效率最高
// strings.Join 连接构建字符串的平均性能最稳定，如果输出的多个字符串是以[]string承载的，那么strings.Join也是不错的选择
// 使用操作符连接的方式最直观、最自然，在编译器知晓欲连接的字符串个数的情况下，使用此种方式可以得到编译器的优化处理
// fmt.Sprintf虽然效率不高，但也不是一无是处，如果是由多种不同类型变量来构建特定格式的字符串，那么这种方式还是最适合的
func main() {

}
