package httpserver

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp_server_go/internal/server"
)

func TestMCPOverStreamableHTTP(t *testing.T) {
	testServer := httptest.NewServer(NewHandler(server.New()))
	defer testServer.Close()

	client := mcp.NewClient(
		&mcp.Implementation{Name: "mcp-http-test", Version: "v0.1.0"},
		nil,
	)
	clientSession, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint: testServer.URL + "/mcp",
	}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer clientSession.Close()

	listResult, err := clientSession.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list tools over HTTP: %v", err)
	}
	if len(listResult.Tools) != 2 {
		t.Fatalf("tool count = %d, want 2", len(listResult.Tools))
	}

	callResult, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "HTTP 用户"},
	})
	if err != nil {
		t.Fatalf("call greet over HTTP: %v", err)
	}
	if callResult.IsError {
		t.Fatal("call greet over HTTP returned an MCP error")
	}
}
