// MCP 示例服务的启动入口。
package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp_server_go/internal/server"
)

func main() {
	// stdout 是 MCP 的 JSON-RPC 通道，日志必须写到 stderr，避免污染协议数据。
	log.SetOutput(os.Stderr)

	mcpServer := server.New()
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("run MCP server: %v", err)
	}
}
