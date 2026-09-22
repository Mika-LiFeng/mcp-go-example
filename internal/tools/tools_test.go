package tools

import (
	"context"
	"testing"
)

func TestGreet(t *testing.T) {
	_, output, err := Greet(context.Background(), nil, GreetInput{Name: "小明"})
	if err != nil {
		t.Fatalf("Greet() error = %v", err)
	}

	const want = "你好，小明！欢迎使用 MCP 示例服务。"
	if output.Greeting != want {
		t.Fatalf("Greet() greeting = %q, want %q", output.Greeting, want)
	}
}

func TestGreetRejectsEmptyName(t *testing.T) {
	_, _, err := Greet(context.Background(), nil, GreetInput{Name: "  "})
	if err == nil {
		t.Fatal("Greet() error = nil, want an error")
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		input     CalculateInput
		want      float64
		wantError bool
	}{
		{
			name:  "add",
			input: CalculateInput{A: 2, B: 3, Operation: "add"},
			want:  5,
		},
		{
			name:  "divide",
			input: CalculateInput{A: 9, B: 3, Operation: "divide"},
			want:  3,
		},
		{
			name:      "divide by zero",
			input:     CalculateInput{A: 9, B: 0, Operation: "divide"},
			wantError: true,
		},
		{
			name:      "unsupported operation",
			input:     CalculateInput{A: 1, B: 2, Operation: "power"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, output, err := Calculate(context.Background(), nil, tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatal("Calculate() error = nil, want an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}
			if output.Result != tt.want {
				t.Fatalf("Calculate() result = %v, want %v", output.Result, tt.want)
			}
		})
	}
}
