// 是什么：中间件是包在 HTTP handler 外的一层处理函数，可以在业务处理前后执行通用逻辑。
// 作用：收到请求后决定是否调用 next，或提前返回；多个中间件组合后形成从外到内进入、从内到外返回的调用链。
// 使用场景：记录请求日志、设置响应头、检查身份或传递超时；本例设置的 context deadline 需要下游协作处理。
//
// 学习目标：使用 func(http.Handler) http.Handler 组合中间件。
// 运行：go run ./05_backend/05_middleware
// 观察：日志顺序为 before -> after；响应包含安全头，handler 能看到 deadline。
// 对照 Express：next.ServeHTTP(w, r) 类似 next()，但下一层返回后还能继续执行代码。
// 练习：增加一个要求 X-API-Key 的教学中间件，失败时直接 return。
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("before: %s %s", r.Method, r.URL.Path)
		defer log.Println("after: handler returned")
		next.ServeHTTP(w, r)
	})
}

func withTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		// 取消是协作式的；handler / 数据库 / HTTP 客户端需要使用 ctx。
		// 单独设置 deadline 不会强行终止 goroutine，也不会自动写出超时响应。
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.SetFlags(0)
	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, exists := r.Context().Deadline()
		fmt.Fprintf(w, "has_deadline=%t\n", exists)
	})
	// 请求从外向内进入，返回时从内向外离开。
	handler := logging(securityHeaders(withTimeout(endpoint)))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	fmt.Println("header:", response.Header().Get("X-Content-Type-Options"))
	fmt.Print(response.Body.String())
}
