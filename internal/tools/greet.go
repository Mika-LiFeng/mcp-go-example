// Package tools 提供 MCP Tool 的输入模型和处理函数。
package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GreetInput 是 greet Tool 的输入参数。
type GreetInput struct {
	Name string `json:"name" jsonschema:"需要问候的姓名"`
}

// GreetOutput 是 greet Tool 的结构化输出。
type GreetOutput struct {
	Greeting string `json:"greeting" jsonschema:"生成的问候语"`
}

// Greet 根据姓名生成问候语。
func Greet(_ context.Context, _ *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, GreetOutput{}, fmt.Errorf("name must not be empty")
	}

	return nil, GreetOutput{Greeting: fmt.Sprintf("你好，%s！欢迎使用 MCP 示例服务。", name)}, nil
}
