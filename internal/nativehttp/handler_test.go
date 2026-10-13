package nativehttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func post(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/event-stream")
	w := httptest.NewRecorder()
	NewHandler().ServeHTTP(w, r)
	return w
}

func TestInitializeAndListTools(t *testing.T) {
	init := post(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`)
	if init.Code != http.StatusOK { t.Fatalf("initialize status = %d, want 200: %s", init.Code, init.Body) }
	var got struct { Result struct { ProtocolVersion string `json:"protocolVersion"`; Capabilities map[string]any `json:"capabilities"` } `json:"result"` }
	if err := json.Unmarshal(init.Body.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Result.ProtocolVersion != protocolVersion { t.Fatalf("protocol = %q, want %q", got.Result.ProtocolVersion, protocolVersion) }

	listed := post(t, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	var tools struct { Result struct { Tools []struct { Name string `json:"name"` } `json:"tools"` } `json:"result"` }
	if err := json.Unmarshal(listed.Body.Bytes(), &tools); err != nil { t.Fatal(err) }
	if len(tools.Result.Tools) != 2 { t.Fatalf("tool count = %d, want 2", len(tools.Result.Tools)) }
}

func TestCallToolAndNotification(t *testing.T) {
	called := post(t, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"calculate","arguments":{"a":7,"b":5,"operation":"multiply"}}}`)
	if called.Code != http.StatusOK { t.Fatalf("call status = %d", called.Code) }
	var got struct { Result struct { Structured map[string]float64 `json:"structuredContent"`; IsError bool `json:"isError"` } `json:"result"` }
	if err := json.Unmarshal(called.Body.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Result.IsError || got.Result.Structured["result"] != 35 { t.Fatalf("call result = %+v, want result 35", got.Result) }

	notification := post(t, `{"jsonrpc":"2.0","method":"notifications/initialized"}`)
	if notification.Code != http.StatusAccepted || notification.Body.Len() != 0 { t.Fatalf("notification response = %d %q, want empty 202", notification.Code, notification.Body.String()) }
}

func TestRejectsCrossOriginAndInvalidProtocolVersion(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewBufferString(`{}`))
	r.Host = "localhost"
	r.Header.Set("Origin", "http://evil.example")
	w := httptest.NewRecorder()
	NewHandler().ServeHTTP(w, r)
	if w.Code != http.StatusForbidden { t.Fatalf("origin status = %d, want 403", w.Code) }

	r = httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewBufferString(`{}`))
	r.Header.Set("MCP-Protocol-Version", "unsupported")
	w = httptest.NewRecorder()
	NewHandler().ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest { t.Fatalf("version status = %d, want 400", w.Code) }
}
