// 是什么：httptest 是测试 HTTP 代码的标准库工具；handler 是接收请求并写出响应的处理代码。
// 作用：ResponseRecorder 在内存中记录响应，NewServer 启动临时本机服务，让客户端测试真实 HTTP 往返。
// 使用场景：检查路由、状态码、响应头和 JSON 内容；先用内存测试覆盖处理逻辑，再按需测试网络交互。
//
// 目标：使用 httptest 测试 HTTP handler，无需固定端口或手动启动后端。
// 对照前端：类似 mock server，但直接复用真实 Go handler。
// 根目录运行：go run ./06_testing/02_httptest
// 根目录测试：go test -v ./06_testing/02_httptest
// 观察：临时本地服务器返回 200 和 {"message":"你好，Go"}，随后自动关闭。
// 练习：添加 ?name= 的空值测试；区分 ResponseRecorder 与 NewServer 的用途。
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

func newHandler() http.Handler {
	mux := http.NewServeMux()
	// Go 1.22 的 ServeMux 支持“方法 路径”；GET 模式也匹配 HEAD。
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Go"
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "你好，" + name}); err != nil {
			log.Printf("写入响应失败: %v", err)
		}
	})
	return mux
}

func main() {
	server := httptest.NewServer(newHandler())
	defer server.Close()
	client := server.Client()
	client.Timeout = 2 * time.Second
	response, err := client.Get(server.URL + "/hello")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.StatusCode)
	fmt.Print(string(body))
}
