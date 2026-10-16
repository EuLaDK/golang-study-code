// 是什么：CRUD 指创建 Create、读取 Read、更新 Update、删除 Delete；API 是客户端调用这些操作的约定。
// 作用：用 GET、POST、PUT、DELETE 等 HTTP 方法管理任务资源，组合路由、JSON、输入校验、锁和关闭流程。
// 使用场景：理解一个任务或订单管理后端的完整请求过程；本例把任务保存在进程内存，重启会丢失，真实持久化需要数据库。
//
// 学习目标：把路由、JSON、校验、错误处理、锁和优雅关闭组合成任务 API。
// 运行：go run ./05_backend/08_crud_api -addr 127.0.0.1:8081
// 测试：go test ./05_backend/08_crud_api -v
// 观察：POST /tasks 返回 201；GET /tasks 返回数组；具体 PowerShell 命令见上级 README。
// 练习：增加截止日期，同时修改输入、存储、响应和测试。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
)

type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// 请求类型不含 ID，防止客户端修改服务端分配的标识。
type taskInput struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type api struct {
	// ponytail: 单进程内存存储，重启丢数据；需要持久化或多实例时换数据库。
	mu     sync.RWMutex
	tasks  map[int]task
	nextID int
}

func newHandler() http.Handler {
	a := &api{tasks: make(map[int]task), nextID: 1}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /tasks", a.list)
	mux.HandleFunc("POST /tasks", a.create)
	mux.HandleFunc("GET /tasks/{id}", a.get)
	mux.HandleFunc("PUT /tasks/{id}", a.update)
	mux.HandleFunc("DELETE /tasks/{id}", a.remove)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("encode response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if _, err := w.Write(append(data, '\n')); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func readTaskInput(w http.ResponseWriter, r *http.Request) (taskInput, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return taskInput{}, false
	}
	body := http.MaxBytesReader(w, r.Body, 1<<20)
	defer body.Close()
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	var input taskInput
	err = decoder.Decode(&input)
	if err == nil {
		// 第二次 Decode 必须读到 EOF，拒绝 {..}{..} 等多个 JSON 值。
		var extra any
		if nextErr := decoder.Decode(&extra); nextErr != io.EOF {
			err = nextErr
			if err == nil {
				err = errors.New("body must contain exactly one JSON value")
			}
		}
	}
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "body must not exceed 1 MiB")
		} else {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		}
		return taskInput{}, false
	}
	input.Title = strings.TrimSpace(input.Title)
	if n := utf8.RuneCountInString(input.Title); n < 1 || n > 100 {
		writeError(w, http.StatusBadRequest, "title must contain 1 to 100 Unicode code points")
		return taskInput{}, false
	}
	return input, true
}

func readID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func (a *api) list(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	// 非 nil 的空切片编码为 []；nil 切片编码为 null。
	items := make([]task, 0, len(a.tasks))
	for _, item := range a.tasks {
		items = append(items, item)
	}
	a.mu.RUnlock()
	// map 遍历无固定顺序；复制后再排序、编码，缩短持锁时间。
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	writeJSON(w, http.StatusOK, items)
}

func (a *api) create(w http.ResponseWriter, r *http.Request) {
	input, ok := readTaskInput(w, r)
	if !ok {
		return
	}
	a.mu.Lock()
	item := task{ID: a.nextID, Title: input.Title, Completed: input.Completed}
	a.tasks[item.ID] = item
	a.nextID++
	a.mu.Unlock()
	w.Header().Set("Location", fmt.Sprintf("/tasks/%d", item.ID))
	writeJSON(w, http.StatusCreated, item)
}

func (a *api) get(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	a.mu.RLock()
	item, exists := a.tasks[id]
	a.mu.RUnlock()
	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *api) update(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	input, ok := readTaskInput(w, r)
	if !ok {
		return
	}
	item := task{ID: id, Title: input.Title, Completed: input.Completed}
	a.mu.Lock()
	_, exists := a.tasks[id]
	if exists {
		a.tasks[id] = item
	}
	a.mu.Unlock()
	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *api) remove(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	a.mu.Lock()
	_, exists := a.tasks[id]
	delete(a.tasks, id)
	a.mu.Unlock()
	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func run(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server := &http.Server{
		Handler:           newHandler(),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- server.Serve(listener) }()
	log.Printf("task API: http://%s (Ctrl+C to stop)", listener.Addr())
	select {
	case err := <-serveErrors:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return errors.Join(err, server.Close())
	}
	if err := <-serveErrors; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.Println("shutdown complete")
	return nil
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8081", "listen address")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *addr); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
