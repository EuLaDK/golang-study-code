// 是什么：embed 把资源文件在构建时打包进程序；io/fs 用统一接口表示文件系统，embed.FS 是其中一种实现。
// 作用：程序可以直接读取自带的模板或文字资源，而不用依赖运行时某个磁盘目录存在。
// 使用场景：分发包含模板、静态资源或默认配置的单个可执行文件；资源变化后需要重新构建，密钥不应嵌入程序。
//
// 目标：embed 把资源编入程序；io/fs 提供统一的只读文件系统接口。
// 运行（仓库根目录）：go run ./04_stdlib/12_fs_embed
// 观察：无需依赖当前工作目录；修改资源后需要重新编译（go run 会重新构建）。
// 练习：增加一个 assets/hello.txt，用 fs.ReadDir 列出两个文件。
package main

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed assets/*.txt
var assets embed.FS

func main() {
	data, err := fs.ReadFile(assets, "assets/message.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("embedded: %s", data)
	entries, err := fs.ReadDir(assets, "assets")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		fmt.Println("entry:", entry.Name(), "directory:", entry.IsDir())
	}
	// embed 模式相对源文件目录；fs 路径使用 /，不能用 ../ 越过根目录。
	// 不要把密钥嵌入二进制：能拿到程序就能提取资源。
}
