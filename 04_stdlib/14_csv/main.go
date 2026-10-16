// 是什么：CSV 是用分隔符表示表格行列的文本格式，字段中也可以包含逗号、引号或换行。
// 作用：encoding/csv 的 Reader/Writer 按 CSV 规则解析和转义，比简单按逗号 Split 更可靠。
// 使用场景：导入导出用户或订单报表；写入完成后 Flush 并检查错误，供电子表格打开的内容还需考虑公式注入。
//
// 目标：encoding/csv 正确处理逗号、引号和换行，不要手动 Split CSV。
// 运行（仓库根目录）：go run ./04_stdlib/14_csv
// 观察：包含逗号的字段被自动加引号；writer.Flush 后还要检查 Error。
// 练习：把名字改成含双引号的文本，观察转义后的输出。
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

func main() {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	for _, row := range [][]string{{"name", "role"}, {"Ada, Jr.", "developer"}} {
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}
	fmt.Print(buffer.String())
	reader := csv.NewReader(&buffer)
	reader.FieldsPerRecord = 2
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("row: %q\n", row)
	}
	// 导出不可信内容供电子表格打开时，还要防范 =、+ 等开头的公式注入。
}
