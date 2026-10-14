# Go MCP 示例服务

这是一个用于学习 Model Context Protocol（MCP）的 Go 示例服务集合。

## 功能概览

所有示例均提供以下两个工具：

- `greet`：根据姓名生成问候语；
- `calculate`：执行两个数字的加减乘除运算。

## 示例列表

项目 `cmd/` 目录下包含 4 个独立示例，分别演示不同的实现方式和传输方式：

| 示例 | 传输方式 | 依赖 | 说明 |
|------|---------|------|------|
| [cmd/mcp-sse](file:///d:/my_github/mcp-go-example/cmd/mcp-sse) | stdio | MCP SDK | 本地 Agent 以子进程启动，最常用的集成方式 |
| [cmd/mcp-http](file:///d:/my_github/mcp-go-example/cmd/mcp-http) | Streamable HTTP | MCP SDK | 远程 HTTP 调用，由 SDK 封装协议细节 |
| [cmd/mcp-native-http](file:///d:/my_github/mcp-go-example/cmd/mcp-native-http) | Streamable HTTP | 仅标准库 | 纯 Go 标准库实现，便于学习协议底层 |
| [cmd/mcp-go-sdk](file:///d:/my_github/mcp-go-example/cmd/mcp-go-sdk) | stdio | MCP SDK | 官方 SDK 最小示例，含服务端与客户端 |

其中 `cmd/mcp-native-http` 是独立的标准库实现，不依赖 MCP SDK 或 Web 框架；其余示例均使用 MCP Go SDK。

## 环境要求

- Go 1.27 或更高版本；
- Node.js 和 `npx`（仅使用 MCP Inspector 调试时需要）。

## 安装依赖

```bash
go mod tidy
```

## 运行 stdio 服务（mcp-sse）

stdio 服务启动后会等待 MCP 客户端连接，不会主动输出交互界面：

```bash
go run ./cmd/mcp-sse
```

> 注意：stdout 是 MCP 的 JSON-RPC 通道，服务日志已重定向到 stderr，避免污染协议数据。

## 运行 HTTP 服务（mcp-http）

基于 SDK 的 Streamable HTTP 服务，默认监听 `127.0.0.1:8080`：

```bash
go run ./cmd/mcp-http
```

| 路径 | 说明 |
|------|------|
| `POST /mcp` | MCP Streamable HTTP 端点 |
| `GET /healthz` | 健康检查，返回 `{"status":"ok"}` |

修改监听地址：

```bash
go run ./cmd/mcp-http -addr 0.0.0.0:8080
```

`0.0.0.0` 会监听所有网卡，仅建议在已经配置访问控制、TLS 和认证后使用。

## 运行原生 HTTP 服务（mcp-native-http）

纯标准库实现，默认监听 `127.0.0.1:8081`，采用无会话 JSON 响应模式：

```bash
go run ./cmd/mcp-native-http
```

MCP 端点：`http://127.0.0.1:8081/mcp`

修改监听地址：

```bash
go run ./cmd/mcp-native-http -addr 127.0.0.1:9090
```

## 运行 SDK 最小示例（mcp-go-sdk）

该示例包含服务端和客户端，客户端会以子进程方式启动服务端并调用 `greet` 工具：

```bash
# 先构建服务端可执行文件
go build -o ./cmd/mcp-go-sdk/server/server.exe ./cmd/mcp-go-sdk/server

# 运行客户端（会自动启动服务端子进程）
go run ./cmd/mcp-go-sdk/client
```

> 客户端代码中硬编码了服务端可执行文件的路径，请根据实际路径修改 [client.go](file:///d:/my_github/mcp-go-example/cmd/mcp-go-sdk/client/client.go) 中的 `server.exe` 路径。

## 使用 MCP Inspector 调试

[MCP Inspector](https://github.com/modelcontextprotocol/inspector) 是官方提供的可视化调试工具，可以查看 MCP 通信的完整过程。

### 调试 stdio 服务

Inspector 会以子进程方式启动服务，适合调试 `mcp-sse` 和 `mcp-go-sdk/server`：

```bash
# 调试 mcp-sse
npx @modelcontextprotocol/inspector go run ./cmd/mcp-sse
```

打开后在 Tools 面板中即可看到 `greet` 和 `calculate` 两个工具，选择工具并填写参数即可调用。

### 调试 HTTP 服务

对于已经在运行的 HTTP 服务，直接将端点地址交给 Inspector：

```bash
# 调试 mcp-http（需先启动服务）
npx @modelcontextprotocol/inspector http://127.0.0.1:8080/mcp

# 调试 mcp-native-http（需先启动服务）
npx @modelcontextprotocol/inspector http://127.0.0.1:8081/mcp
```

### 指定协议版本

Inspector 默认使用较旧的协议版本。若要使用最新的 MCP 协议版本（`2026-07-28`），添加 `--protocol-era modern` 参数（注意**不要**加 `--cli`，否则会切换到脚本化 CLI 模式，需要额外指定 `--method`）：

```bash
# 以最新协议版本（modern era, 2026-07-28）调试 mcp-http
npx @modelcontextprotocol/inspector --protocol-era modern http://127.0.0.1:8080/mcp
```

> ⚠️ 对于 `mcp-http` 示例，仅添加 `--protocol-era modern` 还不够：Streamable HTTP 传输需要在服务端设置 `StreamableHTTPOptions.Stateless = true` 才会接受 `2026-07-28`，否则会降级到 `2025-11-25`。详见 [mcp-http README](file:///d:/my_github/mcp-go-example/cmd/mcp-http/README.MD)。
>
> ⚠️ `mcp-native-http` 示例为原生实现，最高仅支持到 `2025-11-25`，不支持 `2026-07-28`，因此**不要**对其使用 `--protocol-era modern`，否则会连接失败。

### Inspector 使用要点

1. **初始化握手**：连接后 Inspector 会自动发送 `initialize` 请求，服务端返回协议版本、能力和服务端信息；
2. **工具列表**：Inspector 自动发送 `tools/list`，展示所有工具的名称、描述和输入参数 Schema；
3. **调用工具**：在工具面板中填写参数后点击 Call，可查看完整的 JSON-RPC 请求和响应；
4. **消息日志**：底部的 Messages 面板会记录所有 JSON-RPC 通信（`initialize`、`tools/list`、`tools/call` 等），方便学习协议细节。

## 构建

```bash
go build -o bin/mcp-sse ./cmd/mcp-sse
go build -o bin/mcp-http ./cmd/mcp-http
go build -o bin/mcp-native-http ./cmd/mcp-native-http
go build -o bin/mcp-go-sdk-server ./cmd/mcp-go-sdk/server
```

## 测试

```bash
go test ./...
go vet ./...
```

## 客户端配置示例

### stdio 服务（mcp-sse）

客户端需要配置可执行命令，使用生成二进制文件的绝对路径：

```json
{
  "mcpServers": {
    "mcp-example-go": {
      "command": "/absolute/path/to/bin/mcp-sse"
    }
  }
}
```

### HTTP 服务（mcp-http / mcp-native-http）

远程 HTTP 客户端使用 MCP 地址：

```text
http://127.0.0.1:8080/mcp
```

> ⚠️ 当前 HTTP 示例均未加入身份认证，禁止直接暴露到公网；正式部署前应增加 HTTPS、Bearer Token、来源限制和限流。
