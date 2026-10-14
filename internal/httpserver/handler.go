// Package httpserver 提供 MCP Streamable HTTP 传输及健康检查。
package httpserver

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewHandler 创建 HTTP 路由。
//
// /mcp 是远程 MCP 客户端使用的 Streamable HTTP 端点，
// /healthz 用于检查 HTTP 服务是否正常运行。
func NewHandler(mcpServer *mcp.Server) http.Handler {
	if mcpServer == nil {
		panic("mcpServer must not be nil")
	}

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server {
			return mcpServer
		},
		&mcp.StreamableHTTPOptions{
			// 使用 JSON 响应便于学习和使用 curl/Inspector 调试。
			JSONResponse: true,
			// 配置此项为 true 后，会使用 2026-07-28 协议，否则会降级到 2025-11-25 协议。
			//Stateless: true,
		},
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
