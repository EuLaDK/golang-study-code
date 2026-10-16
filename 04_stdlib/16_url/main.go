// 是什么：URL 是包含协议、主机、路径和查询参数的地址，net/url 把这些部分解析为可操作的结构。
// 作用：读取和修改 query 参数，并把空格、& 等字符按所在位置正确转义，避免手动拼接改变参数含义。
// 使用场景：构造搜索链接、读取分页参数或调用上游接口；能解析地址并不等于允许请求该地址。
//
// 目标：net/url 解析地址、读取/设置查询参数并正确转义。
// 运行（仓库根目录）：go run ./04_stdlib/16_url
// 观察：Query 是独立 Values；修改后需通过 Encode 写回 RawQuery。
// 练习：增加多值 tag 参数，用 Query()["tag"] 读取全部值。
package main

import (
	"fmt"
	"net/url"
)

func main() {
	u, err := url.Parse("https://example.com/search?page=1")
	if err != nil {
		panic(err)
	}
	query := u.Query()
	query.Set("q", "Go & backend")
	query.Add("tag", "go")
	query.Add("tag", "web")
	u.RawQuery = query.Encode()
	fmt.Println(u.String())
	fmt.Println("host:", u.Hostname(), "tags:", u.Query()["tag"])
	fmt.Println("path segment:", url.PathEscape("a/b"), "query:", url.QueryEscape("a b"))
	// Parse 成功不代表可安全请求；外部 URL 还需校验协议/目标主机，防范 SSRF。
}
