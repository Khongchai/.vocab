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
	forest "vocab/vocabulary/forest"
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
	forest := forest.NewForest(ctx, func(any) {})
	engine := engine.NewEngine(ctx, inputReader.Read, outputWriter.Write, logger, map[string]func(lsproto.Notification) (any, error){
		"textDocument/didOpen": func(rm lsproto.Notification) (any, error) {
			params, err := unmarshalInto(rm.Params, &lsproto.DidOpenDocumentParams{})
			if err != nil {
				return nil, err
			}
			forest.Plant(params.TextDocument.Uri, params.TextDocument.Text, nil)
			diagnostics := forest.Harvest()
			response := lsproto.NewPublishDiagnosticsNotfication(
				lsproto.PublishDiagnosticsParams{
					Uri:         params.TextDocument.Uri,
					Version:     params.TextDocument.Version,
					Diagnostics: diagnostics,
				},
			)

			return response, nil
		},
		"textDocument/didChange": func(rm lsproto.Notification) (any, error) {
			params, err := unmarshalInto(rm.Params, &lsproto.DidChangeTextDocumentParams{})
			if err != nil {
				return nil, err
			}

			for i := range params.ContentChanges {
				change := params.ContentChanges[i]
				// for now, sequential. In the future we can make this parallel
				forest.Plant(params.TextDocument.Uri, change.Text, change.Range)
			}

			diagnostics := forest.Harvest()
			response := lsproto.NewPublishDiagnosticsNotfication(
				lsproto.PublishDiagnosticsParams{
					Uri:         params.TextDocument.Uri,
					Version:     params.TextDocument.Version,
					Diagnostics: diagnostics,
				},
			)

			return response, nil
		},
	}, map[string]func(lsproto.RequestMessage) (any, error){
		// Not actually doing anything right now.
		"textDocument/diagnostic": func(message lsproto.RequestMessage) (any, error) {
			_, err := unmarshalInto(message.Params, &lsproto.DocumentDiagnosticsParams{})
			if err != nil {
				return nil, err
			}
			response := lsproto.NewFullDocumentDiagnosticResponse(
				message.ID,
				[]lsproto.Diagnostic{},
				map[string][]lsproto.Diagnostic{},
			)

			return response, nil
		},
		"initialize": func(message lsproto.RequestMessage) (any, error) {
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
			return response, nil
		},
	})

	engine.Start()
}

func unmarshalInto[T any](unmarshalled any, params *T) (*T, error) {
	marshalled, err := json.Marshal(unmarshalled)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(marshalled, &params)
	return params, nil
}
