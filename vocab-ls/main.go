package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// two options:
// json rpc https://github.com/golang/tools/blob/e8ff82cb45564142dd895df0a1df546687d861e9/internal/jsonrpc2/stream.go#L26
// hand loop https://github.com/microsoft/typescript-go/blob/bcb8510f109a472fe8ce00ab4c6512dba31bedb7/internal/lsp/server.go#L246

type Request struct {
	Jsonrpc string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type Response struct {
	Jsonrpc string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

type Diagnostic struct {
	Range struct {
		Start Position `json:"start"`
		End   Position `json:"end"`
	} `json:"range"`
	Severity int    `json:"severity"` // 1 = Error
	Source   string `json:"source"`
	Message  string `json:"message"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// lkjsflk
func main() {
	reader := bufio.NewReader(os.Stdin)
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(os.Stdout)

	ctx := context.Background()

	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			fmt.Fprintln(os.Stderr, "decode error:", err)
			return
		}

		switch req.Method {
		case "initialize":
			// Respond with empty capabilities
			encoder.Encode(Response{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Result: map[string]any{
					"capabilities": map[string]any{},
				},
			})

		case "initialized":
			// After initialization, send diagnostics for the whole file
			diagnostic := Diagnostic{
				Severity: 1, // Error
				Source:   "demo-lsp",
				Message:  "Everything is red 😈",
			}
			diagnostic.Range.Start = Position{Line: 0, Character: 0}
			diagnostic.Range.End = Position{Line: 9999, Character: 0}

			// Normally the URI is the open file, but we'll fake one
			params := PublishDiagnosticsParams{
				URI:         "file:///demo.go",
				Diagnostics: []Diagnostic{diagnostic},
			}

			notification := map[string]any{
				"jsonrpc": "2.0",
				"method":  "textDocument/publishDiagnostics",
				"params":  params,
			}
			encoder.Encode(notification)

		default:
			// Reply with empty result for unhandled requests
			encoder.Encode(Response{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Result:  nil,
			})
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
