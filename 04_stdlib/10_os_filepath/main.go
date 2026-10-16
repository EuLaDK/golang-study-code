// 是什么：os 提供文件和操作系统相关操作，path/filepath 按当前平台的路径规则组合或拆分文件路径。
// 作用：创建临时目录、读写文件、查看文件大小；Join 可以避免手写 Windows 或 Linux 的分隔符。
// 使用场景：保存导出结果、处理临时文件和读取本地配置；路径来自外部输入时，组合路径还不能保证访问范围安全。
//
// 目标：临时目录、文件读写、元信息与跨平台路径。
// 运行（仓库根目录）：go run ./04_stdlib/10_os_filepath
// 观察：路径分隔符随平台变化；临时文件在程序结束前清理。
// 练习：增加一个子目录，使用 os.MkdirAll 后再写入文件。
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func run() error {
	dir, err := os.MkdirTemp("", "golearn-files-*")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Fprintln(os.Stderr, "cleanup:", err)
		}
	}()
	name := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(name, []byte("learn Go\n"), 0600); err != nil {
		return err
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	info, err := os.Stat(name)
	if err != nil {
		return err
	}
	fmt.Printf("file=%s bytes=%d content=%q\n", filepath.Base(name), info.Size(), data)
	// ReadFile 适合已知较小文件；大文件使用流式 io。权限位在 Windows 上语义不同。
	// filepath.Join/Clean 不构成目录沙箱，不要直接接收不可信路径作为下载路径。
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
