package nativehttp

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

var toolDefinitions = json.RawMessage(`[
 {"name":"greet","description":"根据姓名生成问候语","inputSchema":{"type":"object","properties":{"name":{"type":"string","description":"需要问候的姓名"}},"required":["name"]}},
 {"name":"calculate","description":"执行两个数字的加减乘除运算","inputSchema":{"type":"object","properties":{"a":{"type":"number"},"b":{"type":"number"},"operation":{"type":"string","enum":["add","subtract","multiply","divide"]}},"required":["a","b","operation"]}}
]`)

func toolResult(text string, structured any, isError bool) map[string]any {
	r := map[string]any{"content": []map[string]string{{"type":"text", "text":text}}, "isError":isError}
	if structured != nil { r["structuredContent"] = structured }
	return r
}
func callTool(raw json.RawMessage) (any, *rpcError) {
	var p struct { Name string `json:"name"`; Arguments json.RawMessage `json:"arguments"` }
	if json.Unmarshal(raw, &p) != nil || p.Name == "" { return nil, &rpcError{-32602, "Invalid tools/call params"} }
	switch p.Name {
	case "greet":
		var a struct { Name *string `json:"name"` }
		if json.Unmarshal(p.Arguments, &a) != nil || a.Name == nil { return toolResult("name must be a string", nil, true), nil }
		name := strings.TrimSpace(*a.Name)
		if name == "" { return toolResult("name must not be empty", nil, true), nil }
		greeting := fmt.Sprintf("你好，%s！欢迎使用 MCP 示例服务。", name)
		return toolResult(greeting, map[string]string{"greeting":greeting}, false), nil
	case "calculate":
		var a struct { A *float64 `json:"a"`; B *float64 `json:"b"`; Operation string `json:"operation"` }
		if json.Unmarshal(p.Arguments, &a) != nil || a.A == nil || a.B == nil || a.Operation == "" { return toolResult("a, b and operation are required with valid types", nil, true), nil }
		var value float64
		switch a.Operation {
		case "add": value = *a.A + *a.B
		case "subtract": value = *a.A - *a.B
		case "multiply": value = *a.A * *a.B
		case "divide": if *a.B == 0 { return toolResult("cannot divide by zero", nil, true), nil }; value = *a.A / *a.B
		default: return toolResult("unsupported operation", nil, true), nil
		}
		if math.IsInf(value, 0) || math.IsNaN(value) { return toolResult("result is not finite", nil, true), nil }
		return toolResult(fmt.Sprintf("%g", value), map[string]float64{"result":value}, false), nil
	default: return nil, &rpcError{-32602, "Unknown tool"}
	}
}
