// 是什么：HTTP 响应由状态码、响应头和响应体组成；JSON 响应把业务数据作为 JSON 文本返回给客户端。
// 作用：状态码说明结果，Content-Type 说明数据格式，统一错误字段让客户端用稳定方式读取失败原因。
// 使用场景：向前端返回成功数据或参数错误；先完成编码再发出状态码，避免编码失败时已经承诺成功。
//
// 学习目标：先准备 JSON，再写响应头、状态码和响应体。
// 运行：go run ./05_backend/03_json_response
// 输出包含：200 {"name":"Go"}；400 {"error":"name is required"}。
// 对照前端：HTTP 400 仍然是成功到达服务端的响应，fetch 不会自动将它视为网络错误。
// 练习：把成功数据改为带 data 字段的结构，并更新输出观察。
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// 写出状态码后不能再改成 500；客户端断开等写入错误只能记录。
	if _, err := w.Write(append(data, '\n')); err != nil {
		log.Printf("write response: %v", err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": name})
}

func main() {
	for _, target := range []string{"/hello?name=Go", "/hello"} {
		response := httptest.NewRecorder()
		hello(response, httptest.NewRequest(http.MethodGet, target, nil))
		fmt.Printf("%d %s", response.Code, response.Body.String())
	}
}
