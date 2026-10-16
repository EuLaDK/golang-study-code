// 是什么：I/O 指数据的输入与输出；io.Reader 和 io.Writer 是读取、写入字节的接口，bufio 在其外层增加缓冲。
// 作用：让文件、网络和内存使用相同的流处理方式；Scanner 可以逐行读，缓冲 Writer 需要 Flush 才把暂存数据写出。
// 使用场景：按行处理日志、大文件或流式请求，减少零碎读写；同时检查读取错误并设置合理的数据大小上限。
//
// 目标：io.Reader/Writer 是通用字节流接口；bufio 在流外增加缓冲。
// 运行（仓库根目录）：go run ./04_stdlib/11_io_bufio
// 观察：Scanner 按行读取，Writer 必须 Flush；例子不需要真实文件。
// 练习：使用 strings.ToUpper 转换每一行，并统计输出字节数。
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(strings.NewReader("first\nsecond\n"))
	// Scanner 默认 token 限制约 64 KiB；按业务上限显式调整，而非无限增长。
	scanner.Buffer(make([]byte, 1024), 128*1024)
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	for scanner.Scan() {
		if _, err := fmt.Fprintln(writer, "line:", scanner.Text()); err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}
	fmt.Print(output.String())
	var limited bytes.Buffer
	n, err := io.Copy(&limited, io.LimitReader(strings.NewReader("abcdef"), 3))
	if err != nil {
		panic(err)
	}
	fmt.Println("copied:", n, limited.String())
	// LimitReader 到上限就返回 EOF，不会报告原输入超限；拒绝超限要读 limit+1。
}
