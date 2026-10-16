// 是什么：这是 HTTP handler 的表驱动测试，ResponseRecorder 充当内存中的响应记录器，不建立真实 TCP 连接。
// 作用：为不同方法和路径构造请求，检查状态码、Content-Type 和 JSON 内容；JSON 解码检查响应格式可正常读取。
// 使用场景：覆盖默认参数、指定参数、未知路径和错误请求方法，快速确认接口处理逻辑。
// 运行：go test -v ./06_testing/02_httptest；需要测试网络交互时，可使用配套 main 示例里的 NewServer。
//
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		path    string
		status  int
		message string
	}{
		{"named", "GET", "/hello?name=Alice", 200, "你好，Alice"},
		{"default", "GET", "/hello", 200, "你好，Go"},
		{"method", "POST", "/hello", 405, ""},
		{"missing", "GET", "/missing", 404, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.path, nil)
			recorder := httptest.NewRecorder()
			newHandler().ServeHTTP(recorder, request) // 直接调用，不建立真实 TCP 连接。
			response := recorder.Result()
			defer response.Body.Close()
			if response.StatusCode != tc.status {
				t.Fatalf("status=%d, want %d", response.StatusCode, tc.status)
			}
			if tc.status != http.StatusOK {
				return
			}
			if got := response.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type=%q", got)
			}
			var body map[string]string
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["message"] != tc.message {
				t.Errorf("message=%q, want %q", body["message"], tc.message)
			}
		})
	}
}
