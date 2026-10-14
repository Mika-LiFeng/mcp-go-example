# MCP stdio 传输示例

基于官方 MCP Go SDK 的 **stdio 传输**示例。服务端通过标准输入/输出与客户端通信，适合本地 MCP 客户端（如 Claude Desktop、Cursor 等 Agent）以子进程方式启动和连接。

> 注意：本示例实际使用的是 `mcp.StdioTransport`（stdio 传输），通过子进程的 stdin/stdout 交换 JSON-RPC 消息。如果需要远程 HTTP 调用，请参考 [cmd/mcp-http](../mcp-http/README.MD)。

## 功能

- 复用 [internal/server](file:///d:/my_github/mcp-go-example/internal/server/server.go) 中注册的两个工具：
  - `greet`：根据姓名生成问候语
  - `calculate`：执行加减乘除运算
- 使用 `mcp.StdioTransport` 启动服务

## 运行

服务启动后会等待 MCP 客户端连接，不会主动输出交互界面：

```bash
go run ./cmd/mcp-sse
```

> 注意：stdout 是 MCP 的 JSON-RPC 通道，服务日志已重定向到 stderr，避免污染协议数据。

## 作为 MCP 服务集成到客户端

构建二进制文件：

```bash
go build -o bin/mcp-sse ./cmd/mcp-sse
```

在 MCP 客户端配置中引用该可执行文件（使用绝对路径）：

```json
{
  "mcpServers": {
    "mcp-example-go": {
      "command": "/absolute/path/to/bin/mcp-sse"
    }
  }
}
```

客户端启动后会自动：
1. 发送 `initialize` 握手请求；
2. 发送 `tools/list` 获取 `greet` 和 `calculate` 两个工具及其参数 Schema；
3. 根据需要发送 `tools/call` 调用工具。

## 使用 MCP Inspector 调试

```bash
npx @modelcontextprotocol/inspector go run ./cmd/mcp-sse
```

Inspector 会自动以子进程方式启动该服务，在 Tools 面板中可看到 `greet` 和 `calculate` 两个工具并直接调用。

## 实现要点

[main.go](file:///d:/my_github/mcp-go-example/cmd/mcp-sse/main.go) 的核心逻辑：

```go
mcpServer := server.New()  // 创建并注册工具
mcpServer.Run(context.Background(), &mcp.StdioTransport{})
```

- `server.New()` 在 [internal/server/server.go](file:///d:/my_github/mcp-go-example/internal/server/server.go) 中通过 `mcp.AddTool` 注册 `greet` 和 `calculate` 两个工具，输入参数的 JSON Schema 由结构体的 `jsonschema` tag 自动推导。
- `mcp.StdioTransport` 从 stdin 逐行读取 JSON-RPC 请求，将响应写入 stdout。
- 日志通过 `log.SetOutput(os.Stderr)` 输出到 stderr，确保 stdout 仅承载协议数据。

## 与其他示例的对比

| 示例 | 传输方式 | 依赖 | 适用场景 |
|------|---------|------|---------|
| `cmd/mcp-sse`（本示例） | stdio | MCP SDK | 本地 Agent 以子进程启动 |
| `cmd/mcp-http` | Streamable HTTP | MCP SDK | 远程 HTTP 调用 |
| `cmd/mcp-native-http` | Streamable HTTP | 仅标准库 | 学习协议底层实现 |
| `cmd/mcp-go-sdk` | stdio | MCP SDK | 最小 SDK 用法演示 |
