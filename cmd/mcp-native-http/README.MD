# MCP 原生 Go HTTP 示例

不依赖 MCP SDK 或任何 Web 框架，仅使用 Go 标准库 `net/http` 从零实现的 **MCP Streamable HTTP** 服务。适合想深入理解 MCP 协议底层细节的学习者。

## 功能

- 纯标准库实现，无第三方依赖
- 支持 MCP Streamable HTTP **JSON 响应模式**（无会话、请求-响应）
- 实现以下 JSON-RPC 方法：
  - `initialize`：协议握手，返回服务端能力
  - `ping`：心跳检测
  - `tools/list`：返回工具列表
  - `tools/call`：调用工具
- 提供两个工具：`greet` 和 `calculate`
- 安全措施：Origin 校验，防止 DNS rebinding 攻击

## 运行

默认监听 `127.0.0.1:8081`：

```bash
go run ./cmd/mcp-native-http
```

自定义监听地址：

```bash
go run ./cmd/mcp-native-http -addr 127.0.0.1:9090
```

## 端点

| 路径 | 说明 |
|------|------|
| `POST /mcp` | MCP Streamable HTTP 端点 |
| `GET /healthz` | 健康检查，返回 `{"status":"ok"}` |

MCP 端点地址：

```
http://127.0.0.1:8081/mcp
```

## 使用 MCP Inspector 调试

```bash
npx @modelcontextprotocol/inspector http://127.0.0.1:8081/mcp
```

> **协议版本说明**：本示例为原生实现，最高仅支持到 `2025-11-25`，**不支持**最新的 `2026-07-28` 版本。因此不要使用 `--cli --protocol-era modern` 参数（该参数会指定使用最新版本，导致连接失败）。Inspector 默认版本即可正常工作。

## 用 curl 手动调用

请求头需包含 `Content-Type: application/json` 和 `Accept: application/json, text/event-stream`。

> **关于协议版本**：`initialize` 请求中的 `protocolVersion` 字段指定客户端使用的 MCP 协议版本。本示例为原生实现，仅支持 `2025-11-25`、`2025-06-18`、`2025-03-26` 三个版本（见 [handler.go#L14-L46](file:///d:/my_github/mcp-go-example/internal/nativehttp/handler.go#L14-L46)），默认协商到 `2025-11-25`。若传入不支持的版本（如 `2026-07-28`），服务端会返回 `400 Bad Request`。

1. 初始化：

```bash
curl -X POST http://127.0.0.1:8081/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"curl","version":"1.0"}}}'
```

服务端响应中的 `protocolVersion` 字段即为协商后的实际协议版本。

2. 列出工具：

```bash
curl -X POST http://127.0.0.1:8081/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

3. 调用 `calculate` 工具：

```bash
curl -X POST http://127.0.0.1:8081/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"calculate","arguments":{"a":7,"b":5,"operation":"multiply"}}}'
```

预期返回：

```json
{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"35"}],"structuredContent":{"result":35},"isError":false}}
```

## 实现要点

核心代码在 [internal/nativehttp](file:///d:/my_github/mcp-go-example/internal/nativehttp)：

- [handler.go](file:///d:/my_github/mcp-go-example/internal/nativehttp/handler.go)：HTTP 入口与 JSON-RPC 分发逻辑
  - 校验 `MCP-Protocol-Version`、`Content-Type`、`Accept` 头
  - 解析 JSON-RPC 请求，根据 `method` 字段分发到对应处理函数
  - 无 `id` 的通知请求返回 `202 Accepted`
- [tools.go](file:///d:/my_github/mcp-go-example/internal/nativehttp/tools.go)：工具定义与调用逻辑
  - `toolDefinitions` 直接以 JSON 字符串定义工具的 name、description、inputSchema
  - `callTool` 根据工具名分发到具体处理逻辑

### 与 SDK 实现的对比

| 方面 | `cmd/mcp-http`（SDK） | `cmd/mcp-native-http`（原生） |
|------|----------------------|------------------------------|
| 依赖 | `go-sdk/mcp` | 仅标准库 |
| 工具注册 | `mcp.AddTool` 自动推导 Schema | 手写 JSON Schema |
| 会话管理 | SDK 内部维护 | 无状态，每次请求独立 |
| 协议细节 | 由 SDK 封装 | 需手动处理 JSON-RPC 编解码 |

> ⚠️ 本示例为学习用途，未实现 SSE 流式响应、工具变更通知、会话管理等高级特性，也未加入身份认证。
