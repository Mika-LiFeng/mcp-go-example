// Package nativehttp 使用 Go 标准库实现无会话 MCP Streamable HTTP JSON 模式。
package nativehttp

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

const protocolVersion = "2025-11-25"
const maxBodySize = 1 << 20

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

// NewHandler 创建 /mcp 和 /healthz 路由。
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", handleMCP)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}

func validVersion(v string) bool {
	return v == protocolVersion || v == "2025-06-18" || v == "2025-03-26"
}
func accepts(header, want string) bool {
	for _, part := range strings.Split(header, ",") {
		mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(part))
		if err == nil && mediaType == want && params["q"] != "0" {
			return true
		}
	}
	return false
}
func validID(id json.RawMessage) bool {
	if len(id) == 0 { return true }
	var v any
	if json.Unmarshal(id, &v) != nil { return false }
	switch v.(type) { case string, float64: return true; default: return false }
}
func writeRPC(w http.ResponseWriter, status int, id json.RawMessage, result any, rpcErr *rpcError) {
	if len(id) == 0 { id = json.RawMessage("null") }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(rpcResponse{JSONRPC: "2.0", ID: id, Result: result, Error: rpcErr})
}
func handleMCP(w http.ResponseWriter, r *http.Request) {
	// 拒绝跨源网页访问，降低本地服务遭 DNS rebinding 攻击的风险。
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		scheme := "http"
		if r.TLS != nil { scheme = "https" }
		if err != nil || u.Scheme != scheme || !strings.EqualFold(u.Host, r.Host) || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
			http.Error(w, "invalid Origin", http.StatusForbidden); return
		}
	}
	if v := r.Header.Get("MCP-Protocol-Version"); v != "" && !validVersion(v) {
		http.Error(w, "unsupported MCP-Protocol-Version", http.StatusBadRequest); return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType); return
	}
	if !accepts(r.Header.Get("Accept"), "application/json") || !accepts(r.Header.Get("Accept"), "text/event-stream") {
		http.Error(w, "Accept must include application/json and text/event-stream", http.StatusNotAcceptable); return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodySize))
	if err != nil {
		status := http.StatusBadRequest
		if _, ok := err.(*http.MaxBytesError); ok { status = http.StatusRequestEntityTooLarge }
		http.Error(w, "cannot read request body", status); return
	}
	if !json.Valid(body) { writeRPC(w, http.StatusBadRequest, nil, nil, &rpcError{-32700, "Parse error"}); return }
	var req rpcRequest
	if json.Unmarshal(body, &req) != nil || req.JSONRPC != "2.0" || req.Method == "" || !validID(req.ID) {
		writeRPC(w, http.StatusBadRequest, nil, nil, &rpcError{-32600, "Invalid Request"}); return
	}
	if len(req.Params) > 0 && bytes.TrimSpace(req.Params)[0] != '{' {
		if len(req.ID) == 0 { http.Error(w, "invalid notification params", http.StatusBadRequest) } else { writeRPC(w, http.StatusOK, req.ID, nil, &rpcError{-32602, "Invalid params"}) }
		return
	}
	// 无 ID 的 JSON-RPC 通知按 Streamable HTTP 约定返回空 202。
	if len(req.ID) == 0 { w.WriteHeader(http.StatusAccepted); return }
	result, rpcErr := dispatch(req)
	writeRPC(w, http.StatusOK, req.ID, result, rpcErr)
}

func dispatch(req rpcRequest) (any, *rpcError) {
	switch req.Method {
	case "initialize":
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
			Capabilities map[string]any `json:"capabilities"`
			ClientInfo *struct { Name string `json:"name"`; Version string `json:"version"` } `json:"clientInfo"`
		}
		if json.Unmarshal(req.Params, &p) != nil || p.ProtocolVersion == "" || p.Capabilities == nil || p.ClientInfo == nil || p.ClientInfo.Name == "" || p.ClientInfo.Version == "" {
			return nil, &rpcError{-32602, "Invalid initialize params"}
		}
		version := p.ProtocolVersion
		if !validVersion(version) { version = protocolVersion }
		return map[string]any{"protocolVersion": version, "capabilities": map[string]any{"tools": map[string]any{}}, "serverInfo": map[string]string{"name": "mcp-native-go", "version": "v0.1.0"}}, nil
	case "ping": return map[string]any{}, nil
	case "tools/list": return map[string]any{"tools": toolDefinitions}, nil
	case "tools/call": return callTool(req.Params)
	default: return nil, &rpcError{-32601, "Method not found"}
	}
}
