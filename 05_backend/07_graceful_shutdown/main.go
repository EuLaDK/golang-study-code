// 是什么：优雅关闭是服务退出时先停止接收新连接，再等待正在处理的请求完成，避免直接打断工作。
// 作用：signal.NotifyContext 把 Ctrl+C 等信号转为取消通知，Server.Shutdown 等待关闭，超时后可用 Close 强制结束。
// 使用场景：服务发布更新或进程停止；关闭等待要使用新的 context，因为触发退出的信号 context 已经取消。
//
// 学习目标：Ctrl+C 时停止接收新请求，并等待已有请求完成。
// 运行：go run ./05_backend/07_graceful_shutdown -addr 127.0.0.1:8080
// 观察：访问 /slow 需要约 2 秒；请求过程中 Ctrl+C，服务会等待请求结束。
// Shutdown 的 context 决定最多等待多久，不是正在处理的请求的 context。
// 练习：将慢请求改为 10 秒，观察 5 秒关闭超时后的强制关闭。
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "try GET /slow, then press Ctrl+C")
	})
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("slow request started")
		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()
		select {
		case <-timer.C:
			fmt.Fprintln(w, "work finished")
			log.Println("slow request finished")
		case <-r.Context().Done():
			log.Println("client disconnected or connection closed")
		}
	})
	return mux
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
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	return serve(ctx, server, listener)
}

func serve(ctx context.Context, server *http.Server, listener net.Listener) error {
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- server.Serve(listener) }()
	log.Printf("listening on http://%s", listener.Addr())
	select {
	case err := <-serveErrors:
		return err
	case <-ctx.Done():
		log.Println("shutting down; waiting up to 5s")
	}
	// 信号 ctx 已取消，关闭流程必须使用新的 context。
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		closeErr := server.Close()
		return errors.Join(err, closeErr)
	}
	if err := <-serveErrors; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.Println("shutdown complete")
	return nil
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *addr); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
