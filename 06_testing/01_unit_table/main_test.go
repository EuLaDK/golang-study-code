// 是什么：这是 testing 的入门示例；表驱动测试让不同输入共用同一套检查，Example 可以核对示例输出。
// 作用：用 t.Run 区分测试情况，Fatalf 记录失败并终止当前测试，Errorf 记录失败后继续检查。
// 使用场景：测试端口解析等有明确输入输出约定的函数，并用 errors.As 确认底层错误类型没有丢失。
// 运行：go test -v ./06_testing/01_unit_table；Example 中的 Output 注释是 Go 测试工具识别的输出断言。
//
package main

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestParsePort(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"common", "8080", 8080, false},
		{"minimum", "1", 1, false},
		{"maximum", "65535", 65535, false},
		{"zero", "0", 0, true},
		{"negative", "-1", 0, true},
		{"too_large", "65536", 0, true},
		{"not_a_number", "abc", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePort(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("parsePort(%q): err=%v, wantErr=%t", tc.input, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("parsePort(%q)=%d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestParsePortPreservesCause(t *testing.T) {
	_, err := parsePort("abc")
	var numberError *strconv.NumError
	if !errors.As(err, &numberError) {
		t.Fatalf("want wrapped *strconv.NumError, got %v", err)
	}
}

// Example 的 Output 注释也是断言；go test 会执行并核对标准输出。
func Example_parsePort() {
	port, err := parsePort("8080")
	fmt.Println(port, err)
	// Output: 8080 <nil>
}
