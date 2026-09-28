// MCP HTTP 示例服务的启动入口。
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

	"mcp_server_go/internal/httpserver"
	"mcp_server_go/internal/server"
)

func main() {
	log.SetOutput(os.Stderr)

	addr := flag.String("addr", "127.0.0.1:8080", "HTTP 监听地址")
	flag.Parse()

	mcpServer := server.New()
	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           httpserver.NewHandler(mcpServer),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownContext, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	serverError := make(chan error, 1)
	go func() {
		log.Printf("MCP HTTP server listening on http://%s", *addr)
		serverError <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run MCP HTTP server: %v", err)
		}
	case <-shutdownContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownContext); err != nil {
			log.Fatalf("shutdown MCP HTTP server: %v", err)
		}
	}
}
