package main

import "fmt"

func main() {
	// Go语言规定：Go工具链将忽略名为testdata的目录。这样开发者在
	// 编写测试时，就可以在名为testdata的目录下存放和管理测试代码依赖的
	// 数据文件。而go test命令在执行时会将被测试程序包源码所在目录设置
	// 为其工作目录，这样如果要使用testdata目录下的某个数据文件，我们无
	// 须再处理各种恼人的路径问题
	fmt.Println("使用testData管理测试依赖的外部数据文件")
}
