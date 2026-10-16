// 是什么：HTTP 客户端是向其他服务发请求的工具；被当前后端调用的另一个服务通常称为上游服务。
// 作用：Client.Do 完成请求，context 和 Timeout 控制等待时间，再检查状态码、读取并解析响应。
// 使用场景：调用支付、地图或内部服务；网络请求成功不保证业务成功，响应体还需限制大小并及时关闭。
//
// 学习目标：调用上游 API 时限制时间、检查状态码、限制响应大小并关闭 Body。
// 运行：go run ./05_backend/06_http_client
// 输出：task={ID:1 Title:Learn Go}；upstream status: 503。
// 对照 fetch：Do 返回 nil error 不代表 HTTP 状态码是 2xx。
// 练习：让测试服务返回错误 JSON，观察解析错误如何向调用者传播。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

type task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func fetchTask(ctx context.Context, client *http.Client, url string) (task, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return task{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return task{}, fmt.Errorf("fetch task: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return task{}, fmt.Errorf("upstream status: %d", response.StatusCode)
	}
	const maxBody = 1 << 20
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBody+1))
	if err != nil {
		return task{}, err
	}
	if len(body) > maxBody {
		return task{}, fmt.Errorf("upstream response exceeds 1 MiB")
	}
	var value task
	if err := json.Unmarshal(body, &value); err != nil {
		return task{}, fmt.Errorf("decode task: %w", err)
	}
	return value, nil
}

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			http.Error(w, "temporarily unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(task{ID: 1, Title: "Learn Go"}); err != nil {
			log.Printf("encode: %v", err)
		}
	}))
	defer server.Close()
	// Client / Transport 支持并发复用；不必为每次请求创建一个。
	client := &http.Client{Timeout: 2 * time.Second}
	defer client.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	value, err := fetchTask(ctx, client, server.URL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("task=%+v\n", value)
	if _, err := fetchTask(ctx, client, server.URL+"/fail"); err != nil {
		fmt.Println(err)
	}
}
