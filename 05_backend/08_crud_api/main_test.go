// 是什么：这是任务 API 的行为测试，通过 httptest 在内存中模拟请求并记录响应。
// 作用：覆盖增删改查、非法请求、并发创建与空服务关闭，并在更新后重新读取，确认存储确实改变。
// 使用场景：修改路由、校验或存储逻辑后，确认状态码、响应内容和共享数据行为仍符合接口约定。
// 运行：go test ./05_backend/08_crud_api；并发测试创建多个请求，检查 ID 唯一且列表结果完整。
//
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func request(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func TestCRUD(t *testing.T) {
	handler := newHandler()
	empty := request(handler, "GET", "/tasks", "")
	if empty.Code != 200 || strings.TrimSpace(empty.Body.String()) != "[]" {
		t.Fatalf("empty list: %d %s", empty.Code, empty.Body.String())
	}
	created := request(handler, "POST", "/tasks", `{"title":"  Learn Go  "}`)
	if created.Code != 201 || created.Header().Get("Location") != "/tasks/1" {
		t.Fatalf("create: %d %s", created.Code, created.Body.String())
	}
	var item task
	if err := json.Unmarshal(created.Body.Bytes(), &item); err != nil {
		t.Fatal(err)
	}
	if item.ID != 1 || item.Title != "Learn Go" || item.Completed {
		t.Fatalf("unexpected task: %+v", item)
	}
	if got := request(handler, "GET", "/tasks/1", ""); got.Code != 200 {
		t.Fatalf("get: %d", got.Code)
	}
	updated := request(handler, "PUT", "/tasks/1", `{"title":"Learn tests","completed":true}`)
	if updated.Code != 200 {
		t.Fatalf("update: %d %s", updated.Code, updated.Body.String())
	}
	if err := json.Unmarshal(updated.Body.Bytes(), &item); err != nil || !item.Completed || item.Title != "Learn tests" {
		t.Fatalf("updated=%+v err=%v", item, err)
	}
	stored := request(handler, "GET", "/tasks/1", "")
	var storedItem task
	if err := json.Unmarshal(stored.Body.Bytes(), &storedItem); err != nil || stored.Code != 200 || storedItem != item {
		t.Fatalf("stored=%+v want=%+v status=%d err=%v", storedItem, item, stored.Code, err)
	}
	deleted := request(handler, "DELETE", "/tasks/1", "")
	if deleted.Code != 204 || deleted.Body.Len() != 0 {
		t.Fatalf("delete: %d %s", deleted.Code, deleted.Body.String())
	}
	if got := request(handler, "GET", "/tasks/1", ""); got.Code != 404 {
		t.Fatalf("get deleted: %d", got.Code)
	}
}

func TestInvalidRequests(t *testing.T) {
	handler := newHandler()
	for _, tc := range []struct {
		name, method, path, body string
		status                   int
	}{
		{"empty title", "POST", "/tasks", `{"title":" "}`, 400},
		{"long title", "POST", "/tasks", `{"title":"` + strings.Repeat("a", 101) + `"}`, 400},
		{"unknown field", "POST", "/tasks", `{"title":"Go","id":999}`, 400},
		{"wrong type", "POST", "/tasks", `{"title":42}`, 400},
		{"broken JSON", "POST", "/tasks", `{`, 400},
		{"null", "POST", "/tasks", `null`, 400},
		{"two values", "POST", "/tasks", `{"title":"Go"} {}`, 400},
		{"trailing garbage", "POST", "/tasks", `{"title":"Go"} x`, 400},
		{"large body", "POST", "/tasks", `{"title":"` + strings.Repeat("a", 1<<20) + `"}`, 413},
		{"large trailing whitespace", "POST", "/tasks", `{"title":"Go"}` + strings.Repeat(" ", 1<<20), 413},
		{"bad ID", "GET", "/tasks/nope", "", 400},
		{"zero ID", "DELETE", "/tasks/0", "", 400},
		{"missing task", "GET", "/tasks/99", "", 404},
		{"update missing", "PUT", "/tasks/99", `{"title":"Go"}`, 404},
		{"delete missing", "DELETE", "/tasks/99", "", 404},
		{"unknown route", "GET", "/unknown", "", 404},
		{"wrong method", "PATCH", "/tasks/1", "", 405},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := request(handler, tc.method, tc.path, tc.body)
			if response.Code != tc.status {
				t.Fatalf("status=%d want=%d body=%s", response.Code, tc.status, response.Body.String())
			}
		})
	}
	r := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"Go"}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 415 {
		t.Fatalf("missing content type: %d", w.Code)
	}
	if list := request(handler, "GET", "/tasks", ""); strings.TrimSpace(list.Body.String()) != "[]" {
		t.Fatalf("invalid requests changed storage: %s", list.Body.String())
	}
}

func TestConcurrentCreates(t *testing.T) {
	handler := newHandler()
	const count = 32
	var workers sync.WaitGroup
	for i := 0; i < count; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			response := request(handler, "POST", "/tasks", `{"title":"Go"}`)
			if response.Code != 201 {
				t.Errorf("create: %d", response.Code)
			}
		}()
	}
	workers.Wait()
	response := request(handler, "GET", "/tasks", "")
	var items []task
	if err := json.Unmarshal(response.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != count {
		t.Fatalf("got %d tasks, want %d", len(items), count)
	}
	for i, item := range items {
		if item.ID != i+1 {
			t.Fatalf("ID=%d at position %d", item.ID, i)
		}
	}
}

func TestShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// 使用临时端口，并确认已取消的信号能让服务启动后自行关闭。
	if err := run(ctx, "127.0.0.1:0"); err != nil {
		t.Fatal(err)
	}
}
