package mcp

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Name            string                 `json:"name"`
		Arguments       map[string]interface{} `json:"arguments"`
		ProtocolVersion string                 `json:"protocolVersion"`
	} `json:"params"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// sessions tracks issued session IDs so we can validate them on
// follow-up requests. A sync.Map is safe for concurrent AI clients.
var sessions sync.Map

const defaultProtocolVersion = "2024-11-05"

func newSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Handler acts as the main entry point for the AI client
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request from AI: Method=%s | URL=%s\n", r.Method, r.URL.String())

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "MCP Server Valid", "ready": true}`))
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Use POST method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("-> AI requesting instruction: %s\n", req.Method)
	resp := Response{JSONRPC: "2.0", ID: req.ID}

	switch req.Method {
	case "initialize":
		version := req.Params.ProtocolVersion
		if version == "" {
			version = defaultProtocolVersion
		}

		sid := newSessionID()
		sessions.Store(sid, true)
		w.Header().Set("Mcp-Session-Id", sid)

		resp.Result = map[string]interface{}{
			"protocolVersion": version,
			"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
			"serverInfo":      map[string]interface{}{"name": "local-go-server", "version": "1.0.0"},
		}
	case "notifications/initialized":
		w.WriteHeader(http.StatusOK)
		return
	case "ping":
		resp.Result = map[string]interface{}{}
	case "tools/list":
		resp.Result = getToolsList()
	case "tools/call":
		handleToolCall(&req, &resp)
	default:
		resp.Error = map[string]interface{}{"code": -32601, "message": "MCP Method Not Supported"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
