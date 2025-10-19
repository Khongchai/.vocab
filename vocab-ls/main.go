package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"vocab/engine"
	"vocab/lib"
	lsproto "vocab/lsp"
	"vocab/vocabulary/compiler"
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
	compiler := compiler.NewCompiler(ctx, func(any) {})
	engine := engine.NewEngine(ctx, inputReader.Read, outputWriter.Write, logger, map[string]func(lsproto.Notification) any{
		"textDocument/didChange": func(rm lsproto.Notification) any {
			var params lsproto.DidChangeTextDocumentParams
			marshalled, err := json.Marshal(rm.Params)
			if err != nil {
				logger.Log("Error while unmarshalling params")
				return nil
			}
			json.Unmarshal(marshalled, &params)

			for i := range params.ContentChanges {
				change := params.ContentChanges[i]
				// for now, sequential. In the future we can make this parallel
				compiler.Accept(params.TextDocument.Uri, change.Text, change.Range)
			}

			diagnostics := compiler.Compile()
			response := lsproto.NewPublishDiagnosticsNotfication(
				lsproto.PublishDiagnosticsParams{
					Uri:         params.TextDocument.Uri,
					Version:     params.TextDocument.Version,
					Diagnostics: diagnostics,
				},
			)

			// uncomment to test error all
			// response := lsproto.NewPublishDiagnosticsNotfication(
			// 	lsproto.PublishDiagnosticsParams{
			// 		Uri:     params.TextDocument.Uri,
			// 		Version: params.TextDocument.Version,
			// 		Diagnostics: []lsproto.Diagnostic{
			// 			{
			// 				Severity: lsproto.DiagnosticsSeverityError,
			// 				Message:  "This is a test; no need to panick!",
			// 				Range: lsproto.Range{
			// 					Start: lsproto.Position{
			// 						Line:      0,
			// 						Character: 0,
			// 					},
			// 					End: lsproto.Position{
			// 						Line:      99999999,
			// 						Character: 99999999,
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// )

			return response
		},
	}, map[string]func(lsproto.RequestMessage) any{
		"textDocument/diagnostic": func(message lsproto.RequestMessage) any {
			response := lsproto.NewFullDocumentDiagnosticResponse(
				message.ID,
				[]lsproto.Diagnostic{},
				map[string][]lsproto.Diagnostic{},
			)

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
