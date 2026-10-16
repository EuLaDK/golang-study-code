// 是什么：模板是带占位符的文本，程序把数据填入后生成最终内容；text/template 与 html/template 面向不同输出。
// 作用：text/template 生成普通文本；html/template 按 HTML 上下文自动转义，避免普通用户输入成为可执行标签。
// 使用场景：生成邮件、报表文本或服务端 HTML；模板源码应可信，用户输入应作为数据传入。
//
// 目标：对比 text/template 文本生成与 html/template 上下文自动转义。
// 运行（仓库根目录）：go run ./04_stdlib/20_template
// 观察：HTML 模板输出 &lt; 等实体，用户输入不会变为可执行的 script 标签。
// 练习：给数据增加 []string 字段，用 {{range .Tags}} 遍历。
package main

import (
	"fmt"
	htmltemplate "html/template"
	"os"
	texttemplate "text/template"
)

func main() {
	data := struct{ Name string }{Name: `<script>alert("hello")</script>`}
	// 模板源码必须可信；不可信内容作为数据传入，而非拼进模板源码。
	plain := texttemplate.Must(texttemplate.New("plain").Parse("text: {{.Name}}\n"))
	if err := plain.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
	html := htmltemplate.Must(htmltemplate.New("html").Parse("<p>Hello, {{.Name}}</p>\n"))
	if err := html.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
	fmt.Println("HTML output is escaped; text output is raw.")
	// 不要把用户输入转换为 template.HTML：这会主动绕过转义，带来 XSS 风险。
	// 模板适合内容生成；构造 SQL 使用参数占位符，不能用模板拼接。
}
