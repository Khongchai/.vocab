package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
				print("Received notification ", n.Method)
				// use n
				// text document will be sent here.
				if n.Method == "textDocument/didChange" {
					params := n.Params.(map[string]interface{})
					document := params["textDocument"].(map[string]interface{})
					documentUri := document["uri"].(string)
					documentVersion := document["version"].(float64)

					diagnosticsRange := map[string]any{
						"start": map[string]any{
							"line":      0,
							"character": 0,
						},
						"end": map[string]any{
							"line":      99999999,
							"character": 99999999,
						},
					}

					response := map[string]any{
						"jsonrpc": "2.0",
						"method":  "textDocument/publishDiagnostics",
						"params": map[string]any{
							"uri":     documentUri,
							"version": documentVersion,
							"diagnostics": [1]map[string]any{
								// diagnostic object
								{
									"severity": lsproto.DiagnosticsSeverityError,
									"range":    diagnosticsRange,
									"message":  "This is a test; no need to panick!",
								},
							},
						},
					}

					print("Sending all hell loose notification")

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
			}
		case lsproto.MessageKindRequest:
			if r, ok := data.Msg.(lsproto.RequestMessage); ok {
				var response map[string]interface{} = nil
				if r.Method == "initialize" {
					response = map[string]any{
						"jsonrpc": "2.0",
						"id":      r.ID, // echo the request id
						"result": map[string]any{
							"capabilities": map[string]any{
								// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#serverCapabilities
								"textDocumentSync": map[string]any{
									"openClose": true,
									"change":    lsproto.TextDocumentSyncKindFull,
								},
								// TODO pull mode diagnostic
								// "diagnosticProvider": map[string]any{
								// 	// a change of date in one vocab can affect another (spaced repetition)
								// 	"interFileDependencies": true,
								// },
							},
							// optional, helps debugging in client logs
							"serverInfo": map[string]any{
								"name":    "vocab-ls",
								"version": "0.0.1",
							},
						},
					}
				}

				if response == nil { // not handled.
					continue
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
}
