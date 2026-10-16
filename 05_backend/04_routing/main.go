// 是什么：路由是根据请求方法和路径，把请求交给对应处理函数的分派规则，ServeMux 是标准库路由器。
// 作用：例如 GET /users/{id} 匹配读取用户的请求，{id} 是路径参数，handler 用 PathValue 读取并校验它。
// 使用场景：为用户、订单等资源组织接口；404 表示路径或资源未找到，405 表示目标路径不支持该请求方法。
//
// 学习目标：Go 1.22 ServeMux 支持请求方法与路径参数。
// 运行：go run ./05_backend/04_routing
// 观察：GET 返回 200；非法 ID 返回 400；POST 返回 405；未知路径返回 404。
// 对照 Express：PathValue("id") 类似 req.params.id，仍需要自己校验。
// 练习：增加 GET /users/{id}/posts，想一想静态路径与参数路径的优先级。
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "user=%d", id)
	})
	for _, request := range []struct{ method, path string }{
		{http.MethodGet, "/users/42"},
		{http.MethodGet, "/users/nope"},
		{http.MethodPost, "/users/42"},
		{http.MethodGet, "/missing"},
	} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(request.method, request.path, nil))
		fmt.Printf("%s %s -> %d %s\n", request.method, request.path, response.Code, strings.TrimSpace(response.Body.String()))
	}
}
