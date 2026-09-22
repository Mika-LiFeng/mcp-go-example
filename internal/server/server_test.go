package server

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewServerExposesTools(t *testing.T) {
	ctx := context.Background()
	mcpServer := New()

	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := mcpServer.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("connect server: %v", err)
	}

	client := mcp.NewClient(
		&mcp.Implementation{Name: "mcp-example-test", Version: "v0.1.0"},
		nil,
	)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	defer func() {
		if err := clientSession.Close(); err != nil {
			t.Errorf("close client session: %v", err)
		}
		if err := serverSession.Wait(); err != nil {
			t.Errorf("wait server session: %v", err)
		}
	}()

	listResult, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(listResult.Tools) != 2 {
		t.Fatalf("tool count = %d, want 2", len(listResult.Tools))
	}

	greetResult, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "测试用户"},
	})
	if err != nil {
		t.Fatalf("call greet: %v", err)
	}
	if greetResult.IsError {
		t.Fatal("call greet returned an MCP error")
	}

	calculateResult, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "calculate",
		Arguments: map[string]any{
			"a":         2,
			"b":         3,
			"operation": "add",
		},
	})
	if err != nil {
		t.Fatalf("call calculate: %v", err)
	}
	if calculateResult.IsError {
		t.Fatal("call calculate returned an MCP error")
	}
}
