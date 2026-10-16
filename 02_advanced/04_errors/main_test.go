// 是什么：这是错误处理的表驱动测试，把输入、期望结果和期望错误放在同一张表中。
// 作用：验证正常 ID、非法范围和数字格式错误，并用 errors.Is/As 确认包装后的错误仍保留原始原因。
// 使用场景：给解析函数补充逻辑或改变错误提示后，确认调用方仍能正确识别错误。
// 运行：go test ./02_advanced/04_errors；t.Run 为每组输入建立独立的子测试。
//
package main

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseID(t *testing.T) {
	for _, test := range []struct {
		text    string
		want    int
		wantErr error
	}{
		{text: "42", want: 42},
		{text: "0", wantErr: ErrInvalidID},
		{text: "-1", wantErr: ErrInvalidID},
		{text: "oops", wantErr: strconv.ErrSyntax},
	} {
		t.Run(test.text, func(t *testing.T) {
			got, err := parseID(test.text)
			if got != test.want || !errors.Is(err, test.wantErr) {
				t.Fatalf("parseID(%q) = %d, %v; want %d, %v", test.text, got, err, test.want, test.wantErr)
			}
			if test.wantErr == strconv.ErrSyntax {
				var parseError *strconv.NumError
				if !errors.As(err, &parseError) || parseError.Num != test.text {
					t.Fatal("被包装的错误应保留原始解析类型与输入")
				}
			}
		})
	}
}
