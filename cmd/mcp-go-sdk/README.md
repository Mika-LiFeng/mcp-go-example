# MCP Go SDK 最小示例

基于官方 [MCP Go SDK](https://go.sdk.modelcontextprotocol.io/) 编写的最小可运行示例，包含服务端和客户端两部分，演示如何通过 **stdio** 传输建立 MCP 连接并调用工具。

> 官方快速开始文档：<https://go.sdk.modelcontextprotocol.io/quick_start>

## 目录结构

```
cmd/mcp-go-sdk/
├── server/server.go   # MCP 服务端，注册一个 greet 工具
└── client/client.go   # MCP 客户端，启动服务端子进程并调用 greet
```

## 服务端

[server.go](file:///d:/my_github/mcp-go-example/cmd/mcp-go-sdk/server/server.go) 做了三件事：

1. 创建一个名为 `greeter` 的 MCP Server；
2. 通过 `mcp.AddTool` 注册 `greet` 工具，输入为 `{ "name": "..." }`；
3. 用 `mcp.StdioTransport` 启动服务，从 stdin 读取 JSON-RPC 请求，向 stdout 写回响应。

核心代码：

```go
server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
server.Run(context.Background(), &mcp.StdioTransport{})
```

## 客户端

[client.go](file:///d:/my_github/mcp-go-example/cmd/mcp-go-sdk/client/client.go) 通过 `mcp.CommandTransport` 启动服务端可执行文件作为子进程，建立 MCP 会话后调用 `greet` 工具并打印结果。

```go
transport := &mcp.CommandTransport{Command: exec.Command(`.../server.exe`)}
session, _ := client.Connect(ctx, transport, nil)
res, _ := session.CallTool(ctx, &mcp.CallToolParams{
    Name:      "greet",
    Arguments: map[string]any{"name": "you"},
})
```

## 运行步骤

> 注意：客户端代码中硬编码了服务端可执行文件的绝对路径，请根据你的实际路径修改 `client.go` 中的 `server.exe` 路径。

1. 构建服务端可执行文件：

```bash
go build -o ./cmd/mcp-go-sdk/server/server.exe ./cmd/mcp-go-sdk/server
```

2. 运行客户端（会自动启动服务端子进程）：

```bash
go run ./cmd/mcp-go-sdk/client
```

预期输出：

```
Hi you
```

## 关键点

- **stdio 传输**：服务端通过标准输入/输出与客户端通信，日志必须写入 stderr，避免污染协议通道。
- **CommandTransport**：客户端以子进程方式启动服务端，适合本地 Agent 集成。
- **工具发现**：客户端连接后会自动发送 `tools/list` 获取服务端提供的工具列表（此处 SDK 内部处理，未显式写出）。
