package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp_server_go/internal/server"
)

func TestHealthz(t *testing.T) {
	testServer := httptest.NewServer(NewHandler(server.New()))
	defer testServer.Close()

	response, err := testServer.Client().Get(testServer.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read health response: %v", err)
	}
	if string(body) != `{"status":"ok"}` {
		t.Fatalf("body = %q, want %q", body, `{"status":"ok"}`)
	}
}
