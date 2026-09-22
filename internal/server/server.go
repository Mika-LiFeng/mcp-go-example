// Package server 负责组装 MCP Server 及其对外暴露的能力。
package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp_server_go/internal/tools"
)

// New 创建并配置一个 MCP Server。
func New() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-example-go",
			Version: "v0.1.0",
		},
		nil,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "greet",
			Description: "根据姓名生成问候语",
		},
		tools.Greet,
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "calculate",
			Description: "执行两个数字的加减乘除运算",
		},
		tools.Calculate,
	)

	return server
}
