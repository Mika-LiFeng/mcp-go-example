package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CalculateInput 是 calculate Tool 的输入参数。
type CalculateInput struct {
	A         float64 `json:"a" jsonschema:"第一个数字"`
	B         float64 `json:"b" jsonschema:"第二个数字"`
	Operation string  `json:"operation" jsonschema:"运算符，只支持 add、subtract、multiply、divide"`
}

// CalculateOutput 是 calculate Tool 的结构化输出。
type CalculateOutput struct {
	Result float64 `json:"result" jsonschema:"计算结果"`
}

// Calculate 执行基础四则运算。
func Calculate(_ context.Context, _ *mcp.CallToolRequest, input CalculateInput) (*mcp.CallToolResult, CalculateOutput, error) {
	switch input.Operation {
	case "add":
		return nil, CalculateOutput{Result: input.A + input.B}, nil
	case "subtract":
		return nil, CalculateOutput{Result: input.A - input.B}, nil
	case "multiply":
		return nil, CalculateOutput{Result: input.A * input.B}, nil
	case "divide":
		if input.B == 0 {
			return nil, CalculateOutput{}, fmt.Errorf("cannot divide by zero")
		}
		return nil, CalculateOutput{Result: input.A / input.B}, nil
	default:
		return nil, CalculateOutput{}, fmt.Errorf("unsupported operation %q", input.Operation)
	}
}
