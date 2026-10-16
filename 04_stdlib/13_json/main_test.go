// 是什么：这是 JSON 解析边界的测试，包含合法输入和故意构造的非法输入。
// 作用：检查 int64 大整数精度，并确认未知字段、多份 JSON、尾随垃圾、无效业务值和过大输入被拒绝。
// 使用场景：修改 API 请求结构或解析逻辑后，确保外部输入仍按约定处理。
// 运行：go test ./04_stdlib/13_json；预期非法输入返回 error，只有错误地接受它才算测试失败。
//
package main

import (
	"strings"
	"testing"
)

func TestDecodeUser(t *testing.T) {
	user, err := decodeUser(`{"id":9007199254740993,"name":"Ada"}`)
	if err != nil || user.ID != 9007199254740993 {
		t.Fatalf("int64 precision lost or valid input rejected: user=%+v err=%v", user, err)
	}
	for _, input := range []string{
		`{"id":1,"name":"Ada","admin":true}`,
		`{"id":1,"name":"Ada"} {"id":2,"name":"Bo"}`,
		`{"id":1,"name":"Ada"} trailing`,
		`{"id":0,"name":"Ada"}`,
		`{"id":1,"name":" "}`,
		`null`,
		strings.Repeat(" ", 4097),
	} {
		if _, err := decodeUser(input); err == nil {
			t.Errorf("accepted invalid JSON request: %.80q", input)
		}
	}
}
