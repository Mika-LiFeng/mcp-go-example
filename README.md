# Go MCP 示例服务

这是一个用于学习 Model Context Protocol（MCP）的 Go 示例服务。

当前版本使用 `stdio` 传输，提供两个 Tool：

- `greet`：根据姓名生成问候语；
- `calculate`：执行加减乘除运算。

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

## 使用 MCP Inspector 调试

```bash
npx @modelcontextprotocol/inspector go run ./cmd/mcp-example
```

在 Inspector 中打开 Tools，选择 `greet` 或 `calculate` 即可调用。

## 构建

```bash
go build -o bin/mcp-example ./cmd/mcp-example
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
