// 是什么：gzip 压缩一个字节流，ZIP 把带文件名的多个条目组成归档，并可分别压缩条目。
// 作用：缩小数据表示或打包文件；压缩 Writer.Close 会完成尾部数据，解压时还要限制还原后的大小。
// 使用场景：压缩传输内容、日志归档和报表打包；处理外部压缩包时还需检查条目数量、路径和总解压大小。
//
// 目标：gzip 压缩单个字节流；zip 把多个带名称的文件打包。
// 运行（仓库根目录）：go run ./04_stdlib/22_archive_compress
// 观察：Writer.Close 完成压缩尾部；本例只在内存读写，不解压到磁盘。
// 练习：在 ZIP 中再加入第二个文本文件，并输出两个条目名称。
package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func boundedRead(reader io.Reader) ([]byte, error) {
	const limit = 4096
	data, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if err != nil {
		return nil, err
	}
	if len(data) > limit {
		return nil, fmt.Errorf("uncompressed data exceeds %d bytes", limit)
	}
	return data, nil
}

func run() error {
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write([]byte("hello gzip\n")); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	reader, err := gzip.NewReader(&compressed)
	if err != nil {
		return err
	}
	defer reader.Close()
	plain, err := boundedRead(reader)
	if err != nil {
		return err
	}
	fmt.Printf("gzip restored: %s", plain)

	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	entry, err := zipWriter.Create("notes.txt")
	if err != nil {
		return err
	}
	if _, err := io.WriteString(entry, "learn archives\n"); err != nil {
		return err
	}
	if err := zipWriter.Close(); err != nil {
		return err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(archive.Bytes()), int64(archive.Len()))
	if err != nil {
		return err
	}
	for _, file := range zipReader.File {
		stream, err := file.Open()
		if err != nil {
			return err
		}
		data, readErr := boundedRead(stream)
		closeErr := stream.Close() // 在循环中及时关闭，避免 defer 积累资源。
		if readErr != nil {
			return readErr
		}
		if closeErr != nil {
			return closeErr
		}
		fmt.Printf("zip entry %s: %s", file.Name, data)
	}
	// 不可信压缩包需限制条目数、单文件和总解压大小；文件名还可能带 ../。
	// 本例只读取自己创建的一个条目；不应直接扩展为任意压缩包解压工具。
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
