package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"vocab/engine"
	"vocab/lib"
	lsproto "vocab/lsp"
)

// two options:
// json rpc https://github.com/golang/tools/blob/e8ff82cb45564142dd895df0a1df546687d861e9/internal/jsonrpc2/stream.go#L26
// hand loop https://github.com/microsoft/typescript-go/blob/bcb8510f109a472fe8ce00ab4c6512dba31bedb7/internal/lsp/server.go#L246

func main() {
	print("Starting vocab-ls...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lib.NewInputReader(os.Stdin)
	outputWriter := lib.NewOutputWriter(os.Stdout)
	logger := lib.NewLogger(os.Stderr)
	engine := engine.NewEngine(ctx, inputReader.Read, outputWriter.Write, logger, map[string]func(lsproto.Notification) any{
		"textDocument/didChange": func(rm lsproto.Notification) any {
			params := rm.Params.(map[string]any)
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

			return response
		},
	}, map[string]func(lsproto.RequestMessage) any{
		// not working yet, circle back to this later.
		"textDocument/diagnostic": func(message lsproto.RequestMessage) any {
			params := message.Params.(map[string]any)
			document := params["textDocument"].(map[string]any)
			documentUri := document["uri"].(string)

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

			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      message.ID,
				"result": map[string]any{
					"relatedDocuments": map[string]any{
						documentUri: map[string]any{
							"kind": lsproto.DocumentDiagnosticReportKindFull,
							"items": []map[string]any{
								// diagnostic object
								{
									"severity": lsproto.DiagnosticsSeverityError,
									"range":    diagnosticsRange,
									"message":  "This is a test; no need to panick!",
								},
							},
						},
					},
				},
			}

			return response
		},
		"initialize": func(message lsproto.RequestMessage) any {
			response := map[string]any{
				"jsonrpc": "2.0",
				"id":      message.ID, // echo the request id
				"result": map[string]any{
					"capabilities": map[string]any{
						// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#serverCapabilities
						"textDocumentSync": map[string]any{
							"openClose": true,
							"change":    lsproto.TextDocumentSyncKindFull,
						},
						"diagnosticProvider": map[string]any{
							// a change of date in one vocab can affect another (spaced repetition)
							"interFileDependencies": true,
						},
					},
					// optional, helps debugging in client logs
					"serverInfo": map[string]any{
						"name":    "vocab-ls",
						"version": "0.0.1",
					},
				},
			}
			return response
		},
	})

	engine.Start()
}
