// MCP 原生 Go HTTP 示例服务，不使用 MCP SDK 或 Web 框架。
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mcp_server_go/internal/nativehttp"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8081", "HTTP 监听地址")
	flag.Parse()
	srv := &http.Server{Addr: *addr, Handler: nativehttp.NewHandler(), ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		log.Printf("native MCP HTTP server listening on http://%s/mcp", *addr)
		errCh <- srv.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			log.Fatal(err)
		}
	}
}
