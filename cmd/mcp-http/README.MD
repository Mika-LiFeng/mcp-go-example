# MCP Streamable HTTP 示例

基于官方 MCP Go SDK 的 **Streamable HTTP** 传输示例，将 MCP 服务暴露为 HTTP 端点，供远程 MCP 客户端（如 Agent、MCP Inspector）通过 HTTP 调用。

## 功能

- 复用 [internal/server](file:///d:/my_github/mcp-go-example/internal/server/server.go) 中注册的两个工具：
  - `greet`：根据姓名生成问候语
  - `calculate`：执行加减乘除运算
- 通过 [internal/httpserver](file:///d:/my_github/mcp-go-example/internal/httpserver/handler.go) 封装 SDK 的 `mcp.NewStreamableHTTPHandler`，提供 `/mcp` 端点
- 附带 `/healthz` 健康检查端点

## 运行

默认监听 `127.0.0.1:8080`：

```bash
go run ./cmd/mcp-http
```

自定义监听地址：

```bash
go run ./cmd/mcp-http -addr 0.0.0.0:8080
```

## 端点

| 路径 | 说明 |
|------|------|
| `POST /mcp` | MCP Streamable HTTP 端点，接收 JSON-RPC 请求 |
| `GET /healthz` | 健康检查，返回 `{"status":"ok"}` |

MCP 端点地址：

```
http://127.0.0.1:8080/mcp
```

## SDK 版本与协议版本

本示例使用 MCP Go SDK **v1.8.0**。SDK 版本与 MCP 协议版本的支持关系如下：

| SDK 版本 | 支持的最高协议版本 | 说明 |
|---------|------------------|------|
| v1.6.x 及更早 | `2025-11-25` | 不支持 `2026-07-28` |
| **v1.7.0** | `2026-07-28` | 首次完整支持 `2026-07-28` 协议（beta 在 v1.7.0-pre.1） |
| **v1.8.0**（本示例） | `2026-07-28` | 完整支持，向后兼容 |

**向后兼容性**：SDK 始终向下兼容，同时支持 `2026-07-28`、`2025-11-25`、`2025-06-18`、`2025-03-26`、`2024-11-05` 共 5 个版本。连接时客户端与服务端协商取双方都支持的最新版本。

**如何查看协议版本常量**：

1. **源码位置**：协议版本常量定义在 SDK 源码的 `mcp/shared.go` 文件中。Go module 缓存中的路径为：

   ```
   $GOMODCACHE/github.com/modelcontextprotocol/go-sdk@v1.8.0/mcp/shared.go
   ```

   关键代码片段：

   ```go
   const (
       latestProtocolVersion   = protocolVersion20260728
       protocolVersion20260728 = "2026-07-28"
       protocolVersion20251125 = "2025-11-25"
       protocolVersion20250618 = "2025-06-18"
       protocolVersion20250326 = "2025-03-26"
       protocolVersion20241105 = "2024-11-05"
   )

   var supportedProtocolVersions = []string{
       protocolVersion20260728,
       protocolVersion20251125,
       protocolVersion20250618,
       protocolVersion20250326,
       protocolVersion20241105,
   }
   ```

2. **通过代码获取**：SDK 提供了导出函数 `mcp.SupportedProtocolVersions()`，返回当前 SDK 支持的所有协议版本（从新到旧）：

   ```go
   import "github.com/modelcontextprotocol/go-sdk/mcp"

   versions := mcp.SupportedProtocolVersions()
   // ["2026-07-28", "2025-11-25", "2025-06-18", "2025-03-26", "2024-11-05"]
   ```

3. **运行时查看**：在 `initialize` 响应的 `protocolVersion` 字段中可以看到协商后的实际协议版本。

> **本示例的实际协商版本**：虽然 SDK v1.8.0 支持 `2026-07-28`，但 Streamable HTTP 传输**仅在** `StreamableHTTPOptions.Stateless = true` 时才接受 `2026-07-28`。本示例的 [handler.go](file:///d:/my_github/mcp-go-example/internal/httpserver/handler.go) 未设置 `Stateless`（默认为有状态模式），因此 HTTP 客户端实际会协商到 **`2025-11-25`**。如需使用 `2026-07-28`，需在 `StreamableHTTPOptions` 中添加 `Stateless: true`。

## 使用 MCP Inspector 调试

```bash
# 以默认协议版本（legacy era，协商到 2025-11-25）调试 mcp-http
npx @modelcontextprotocol/inspector http://127.0.0.1:8080/mcp

# 以最新协议版本（modern era，2026-07-28）调试 mcp-http
# 注意：需在服务端设置 StreamableHTTPOptions.Stateless = true，否则会降级到 2025-11-25
npx @modelcontextprotocol/inspector --protocol-era modern http://127.0.0.1:8080/mcp
```

在 Inspector 中打开 Tools 面板，即可看到 `greet` 和 `calculate` 两个工具并直接调用。

> **`--protocol-era` 说明**：该参数控制 Inspector 协商的协议版本，可选 `legacy`（默认，使用 `initialize` 握手）、`auto`（先探测 `server/discover`，失败则回退）、`modern`（强制使用 `2026-07-28`，不回退）。**不要**与 `--cli` 一起使用——`--cli` 会切换到脚本化 CLI 模式，要求额外指定 `--method` 参数（如 `--method tools/list`），不会启动 Web 图形界面。

> **指定最新协议版本**：本示例默认协商到 `2025-11-25`（见上文"SDK 版本与协议版本"）。若要使用最新的 `2026-07-28`，除了在 Inspector 中添加 `--protocol-era modern` 参数外，还需在 [handler.go](file:///d:/my_github/mcp-go-example/internal/httpserver/handler.go) 的 `StreamableHTTPOptions` 中设置 `Stateless: true`，否则连接会降级到 `2025-11-25`。

## 用 curl 手动调用

MCP Streamable HTTP 要求请求头包含 `Content-Type: application/json` 和 `Accept: application/json, text/event-stream`。

> **关于协议版本**：`initialize` 请求中的 `protocolVersion` 字段指定客户端使用的 MCP 协议版本。本示例基于 MCP Go SDK v1.8.0，服务端支持 `2026-07-28` ~ `2024-11-05` 共 5 个版本，但由于未启用 `Stateless` 模式，HTTP 传输实际协商到 `2025-11-25`。下面示例使用 `2025-11-25`。若启用 `Stateless: true`，可将 `protocolVersion` 改为 `2026-07-28`。

1. 初始化：

```bash
curl -X POST http://127.0.0.1:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"curl","version":"1.0"}}}'
```

服务端响应中的 `protocolVersion` 字段即为协商后的实际协议版本。

2. 列出工具：

```bash
curl -X POST http://127.0.0.1:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

3. 调用 `greet` 工具：

```bash
curl -X POST http://127.0.0.1:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"greet","arguments":{"name":"世界"}}}'
```

## 实现要点

[handler.go](file:///d:/my_github/mcp-go-example/internal/httpserver/handler.go) 中使用 `JSONResponse: true`，表示服务端以普通 JSON 响应（而非 SSE 流）返回结果，便于调试和学习。生产环境若需要服务端主动推送通知（如 `notifications/tools/list_changed`），可关闭此选项使用 SSE 流式响应。

> ⚠️ 当前示例未加入身份认证，禁止直接暴露到公网；正式部署前应增加 HTTPS、Bearer Token、来源限制和限流。
