// 是什么：正则表达式是用一段模式描述文本格式的规则，regexp 将模式编译成可重复使用的匹配工具。
// 作用：判断文本是否符合模式、提取括号分组中的内容、替换匹配片段；^ 和 $ 可以限定整个文本的开始和结束。
// 使用场景：解析固定格式的版本号或提取文本片段；简单包含判断用 strings 更直接，Go 正则也有语法支持边界。
//
// 目标：编译正则、匹配与提取子组；重复使用时只编译一次。
// 运行（仓库根目录）：go run ./04_stdlib/09_regexp
// 观察：提取结果第 0 项是整个匹配；Go 正则不支持反向引用/环视。
// 练习：让版本号同时接受 v12.34 和 v12.34.56。
package main

import (
	"fmt"
	"regexp"
)

func main() {
	// MustCompile 适合开发者写死的表达式；外部输入用 Compile 并处理错误。
	version := regexp.MustCompile(`^v([0-9]+)\.([0-9]+)$`)
	fmt.Println("valid:", version.MatchString("v1.22"))
	fmt.Println("groups:", version.FindStringSubmatch("v1.22"))
	spaces := regexp.MustCompile(`\s+`)
	fmt.Println(spaces.ReplaceAllString("Go   web\tserver", " "))
	// 简单前缀/包含判断用 strings；正则不适合完整校验 HTML、URL 或邮箱语义。
}
