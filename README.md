# Go MCP 示例服务

这是一个用于学习 Model Context Protocol（MCP）的 Go 示例服务。

当前版本提供两种传输方式，并提供两个 Tool：

- `greet`：根据姓名生成问候语；
- `calculate`：执行加减乘除运算。

传输方式：

- `stdio`：适合本地 MCP 客户端启动子进程；
- `Streamable HTTP`：适合远程 MCP 客户端通过 HTTP 调用。

其中 `cmd/mcp-native-http.go` 是独立的标准库实现，不依赖 MCP SDK 或 Web 框架，提供 `greet` 和 `calculate` 工具。项目中的其他示例仍使用 MCP Go SDK。

## 环境要求

- Go 1.27 或更高版本；
- Node.js 和 `npx`（仅使用 MCP Inspector 调试时需要）。

## 安装依赖

```bash
go mod tidy
```

## 运行服务

stdio 服务启动后会等待 MCP 客户端连接，不会主动输出交互界面：

```bash
go run ./cmd/mcp-example
```

注意：stdout 是 MCP 协议通道，服务日志应写入 stderr。

## 运行 HTTP 服务

HTTP 服务默认监听 `127.0.0.1:8080`：

```bash
go run ./cmd/mcp-http
```

远程 MCP 端点：

```text
http://127.0.0.1:8080/mcp
```

健康检查：

```text
http://127.0.0.1:8080/healthz
```

修改监听地址：

```bash
go run ./cmd/mcp-http -addr 0.0.0.0:8080
```

`0.0.0.0` 会监听所有网卡，仅建议在已经配置访问控制、TLS 和认证后使用。

## 使用 MCP Inspector 调试

```bash
npx @modelcontextprotocol/inspector go run ./cmd/mcp-example
```

在 Inspector 中打开 Tools，选择 `greet` 或 `calculate` 即可调用。

HTTP 服务也可以直接交给 Inspector：

```bash
npx @modelcontextprotocol/inspector http://127.0.0.1:8080/mcp
```

## 运行原生 Go HTTP 示例

该实现只使用 Go 标准库，默认监听 `127.0.0.1:8081`，支持 MCP 初始化、工具列表和工具调用，采用无会话 JSON 响应模式：

```bash
go run ./cmd/mcp-native-http.go
```

将 Inspector 连接到 `http://127.0.0.1:8081/mcp` 。要修改监听地址：

```bash
go run ./cmd/mcp-native-http.go -addr 127.0.0.1:9090
```

## 构建

```bash
go build -o bin/mcp-example ./cmd/mcp-example
go build -o bin/mcp-http ./cmd/mcp-http
go build -o bin/mcp-native-http ./cmd/mcp-native-http.go
```

## 测试

```bash
go test ./...
go vet ./...
```

## 客户端配置示例

客户端需要配置可执行命令。请使用生成二进制文件的绝对路径：

```json
{
  "mcpServers": {
    "mcp-example-go": {
      "command": "/absolute/path/to/mcp-example"
    }
  }
}
```

远程 HTTP 客户端则使用下面的 MCP 地址：

```text
http://127.0.0.1:8080/mcp
```

当前 HTTP 示例未加入身份认证，禁止直接暴露到公网；正式部署前应增加 HTTPS、Bearer Token、来源限制和限流。
