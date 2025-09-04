package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vocab/jsonrpc"
	lsproto "vocab/lsp"

	"github.com/go-json-experiment/json"
)

// two options:
// json rpc https://github.com/golang/tools/blob/e8ff82cb45564142dd895df0a1df546687d861e9/internal/jsonrpc2/stream.go#L26
// hand loop https://github.com/microsoft/typescript-go/blob/bcb8510f109a472fe8ce00ab4c6512dba31bedb7/internal/lsp/server.go#L246

func main() {
	print("Starting vocab...\n")

	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	reader := jsonrpc.NewJsonrpc(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	// to attach a debugger, give sometime before starting the main loop
	time.Sleep(10 * time.Second)

	// engine := NewEngine(reader, writer)

	// engine.start()

	for { // https://github.com/microsoft/typescript-go/blob/main/internal/lsp/server.go#L246
		data, err := reader.Read()

		if err != nil {
			fmt.Fprintln(os.Stderr, "decode error:", err)
			continue
		}

		switch data.Kind {
		case lsproto.MessageKindNotification:
			if n, ok := data.Msg.(lsproto.Notification); ok {
				// use n
				print(n.Method)
			}
		case lsproto.MessageKindRequest:
			if r, ok := data.Msg.(lsproto.RequestMessage); ok {
				print(r.Method)
				response := map[string]any{
					"jsonrpc": "2.0",
					"id":      r.ID, // echo the request id
					"result": map[string]any{
						"capabilities": map[string]any{
							"textDocumentSync": map[string]any{
								"openClose": true,
								// 1 = Full, 2 = Incremental (LSP 3.17); pick what you implement
								"change": 2,
							},
						},
						// optional, helps debugging in client logs
						"serverInfo": map[string]any{
							"name":    "vocab-ls",
							"version": "0.0.1",
						},
					},
				}
				out, err := json.Marshal(response)
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				}

				if _, err := fmt.Fprintf(writer, "Content-Length: %d\r\n\r\n", len(out)); err != nil {
					panic("wtf")
				}
				if _, err := writer.Write(out); err != nil {
					panic("wtf bro")
				}
				writer.Flush()
			}
		case lsproto.MessageKindResponse:
			if r, ok := data.Msg.(lsproto.ResponseMessage); ok {
				print(r.ID)
			}
		default:
			print("No default message handler found.")
		}
	}

	// for {
	// 	var req Request
	// 	if err := decoder.Decode(&req); err != nil {
	// 		fmt.Fprintln(os.Stderr, "decode error:", err)
	// 		return
	// 	}

	// 	switch req.Method {
	// 	case "initialize":
	// 		// Respond with empty capabilities
	// 		encoder.Encode(Response{
	// 			Jsonrpc: "2.0",
	// 			ID:      req.ID,
	// 			Result: map[string]any{
	// 				"capabilities": map[string]any{},
	// 			},
	// 		})

	// 	case "initialized":
	// 		// After initialization, send diagnostics for the whole file
	// 		diagnostic := Diagnostic{
	// 			Severity: 1, // Error
	// 			Source:   "demo-lsp",
	// 			Message:  "Everything is red 😈",
	// 		}
	// 		diagnostic.Range.Start = Position{Line: 0, Character: 0}
	// 		diagnostic.Range.End = Position{Line: 9999, Character: 0}

	// 		// Normally the URI is the open file, but we'll fake one
	// 		params := PublishDiagnosticsParams{
	// 			URI:         "file:///demo.go",
	// 			Diagnostics: []Diagnostic{diagnostic},
	// 		}

	// 		notification := map[string]any{
	// 			"jsonrpc": "2.0",
	// 			"method":  "textDocument/publishDiagnostics",
	// 			"params":  params,
	// 		}
	// 		encoder.Encode(notification)

	// 	default:
	// 		// Reply with empty result for unhandled requests
	// 		encoder.Encode(Response{
	// 			Jsonrpc: "2.0",
	// 			ID:      req.ID,
	// 			Result:  nil,
	// 		})
	// 	}

	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 	}
	// }
}
