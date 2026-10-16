// 是什么：这是优雅关闭的实际 HTTP 测试，用 channel 控制一个请求何时开始、何时允许完成。
// 作用：触发取消后确认服务正在关闭但尚未返回，再放行请求，检查响应完整且服务正常结束。
// 使用场景：修改退出流程后，确认服务不会把仍在处理的请求直接打断；临时端口避免依赖固定端口。
// 运行：go test ./05_backend/07_graceful_shutdown；t.Cleanup 安排测试退出时清理资源，超时用于防止测试无限等待。
//
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestServeWaitsForActiveRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	deadline, stopDeadline := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(stopDeadline)
	started := make(chan struct{})
	release := make(chan struct{})
	releaseHandler := sync.OnceFunc(func() { close(release) })
	shutdownStarted := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		fmt.Fprint(w, "work finished")
	})}
	server.RegisterOnShutdown(func() { close(shutdownStarted) })
	served := make(chan error, 1)
	go func() {
		served <- serve(ctx, server, listener)
		close(served)
	}()
	client := &http.Client{Timeout: 5 * time.Second}
	t.Cleanup(func() {
		releaseHandler()
		cancel()
		server.Close()
		listener.Close()
		client.CloseIdleConnections()
		select {
		case <-served:
		case <-deadline.Done():
			t.Error("server did not stop during cleanup")
		}
	})

	responseResult := make(chan error, 1)
	go func() {
		response, err := client.Get("http://" + listener.Addr().String())
		if err != nil {
			responseResult <- err
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err == nil && (response.StatusCode != http.StatusOK || string(body) != "work finished") {
			err = fmt.Errorf("status=%d body=%q", response.StatusCode, body)
		}
		responseResult <- err
	}()
	select {
	case <-started:
	case <-deadline.Done():
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case <-shutdownStarted:
	case <-deadline.Done():
		t.Fatal("shutdown did not start")
	}
	// 确认 Shutdown 已开始时，请求仍被阻塞，服务不能提前返回。
	select {
	case err := <-served:
		t.Fatalf("server returned before active request finished: %v", err)
	default:
	}
	releaseHandler()
	select {
	case err := <-responseResult:
		if err != nil {
			t.Fatalf("active request failed during shutdown: %v", err)
		}
	case <-deadline.Done():
		t.Fatal("request did not finish")
	}
	select {
	case err := <-served:
		if err != nil {
			t.Fatalf("graceful shutdown failed: %v", err)
		}
	case <-deadline.Done():
		t.Fatal("server did not finish")
	}
}
