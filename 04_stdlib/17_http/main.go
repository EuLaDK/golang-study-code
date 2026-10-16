// 是什么：HTTP 是客户端与服务端交换请求和响应的协议；net/http 同时提供客户端和服务端能力。
// 作用：ServeMux 按路由分派请求，handler 处理请求并写响应，Client 发请求；httptest.NewServer 创建临时本机测试服务。
// 使用场景：编写 API 或调用第三方服务；区分网络错误和 HTTP 错误状态，并设置超时、限制读取大小和关闭响应体。
//
// 目标：用 net/http 创建路由与客户端，处理超时、状态码和响应体。
// 运行（仓库根目录）：go run ./04_stdlib/17_http
// 观察：本机随机端口临时服务返回 JSON，演示完成后自动退出，不访问外网。
// 练习：请求不存在的路径，观察 Do 不因 HTTP 404 而返回 error。
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

func run() error {
	mux := http.NewServeMux()
	// Go 1.22 起支持方法和路径参数；PathValue 取得路径中的 name。
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "hello " + r.PathValue("name")}); err != nil {
			fmt.Fprintln(os.Stderr, "write response:", err)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := &http.Client{Timeout: 2 * time.Second}
	defer client.CloseIdleConnections()
	response, err := client.Get(server.URL + "/hello/Ada")
	if err != nil {
		return err
	}
	defer response.Body.Close() // 收到响应就安排关闭；读取结束后可复用连接。
	const limit = 4096
	body, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return err
	}
	if len(body) > limit {
		return fmt.Errorf("response exceeds %d bytes", limit)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", response.Status)
	}
	fmt.Println("status:", response.Status)
	fmt.Println("content type:", response.Header.Get("Content-Type"))
	fmt.Printf("body: %s", body)
	// 默认客户端没有整体超时；正式服务也应配置 http.Server 的超时。
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
