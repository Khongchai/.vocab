package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"vocab/engine"
	"vocab/lib"
	lsproto "vocab/lsp"
	"vocab/vocabulary/forest"
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

	collectVocabFilesInRootPath := func(rootPath string) {
		filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			isFile := d.Type().IsRegular()
			isVocab := func() bool {
				chunks := strings.Split(d.Name(), ".")
				extension := chunks[len(chunks)-1]
				return extension == "vocab"
			}()

			if !isFile || !isVocab {
				return nil
			}

			bytes, readErr := os.ReadFile(path)
			if readErr != nil {
				logger.Logf("Can't read content at %s", rootPath)
			}
			fileContent := string(bytes)
			forest.Plant(fmt.Sprintf("%s%s", "file://", path), fileContent, nil)

			return nil
		})
	}

	engine := engine.NewEngine(ctx, inputReader.Read, outputWriter.Write, logger,
		map[string]func(lsproto.Notification) (any, error){
			"workspace/didDeleteFiles": func(request lsproto.Notification) (any, error) {
				params, err := unmarshalInto(request.Params, &lsproto.DeleteFilesParms{})
				if err != nil {
					return nil, err
				}

				for _, file := range params.Files {
					forest.Remove(file.Uri)
				}

				return nil, nil
			},
			"textDocument/didOpen": func(rm lsproto.Notification) (any, error) {
				params, err := unmarshalInto(rm.Params, &lsproto.DidOpenDocumentParams{})
				if err != nil {
					return nil, err
				}

				forest.Plant(params.TextDocument.Uri, params.TextDocument.Text, nil)

				return diagnosticsToNotificationResponse(
					params.TextDocument.Uri,
					params.TextDocument.Version,
					nil,
				), nil
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

				return nil, nil
			},
		}, map[string]func(lsproto.RequestMessage) (any, error){
			"vocab/collectFromThisFile": func(lsproto.RequestMessage) (any, error) {
				return nil, nil
			},
			"vocab/collectAll": func(lsproto.RequestMessage) (any, error) {
				return nil, nil
			},
			"textDocument/hover": func(rm lsproto.RequestMessage) (any, error) {
				_, err := unmarshalInto(rm.Params, &lsproto.HoverParams{})
				if err != nil {
					return nil, err
				}

				return lsproto.NewTextDocumentHoverResponse(
					rm.ID,
					"lol",
					nil,
				), nil
			},
			"textDocument/diagnostic": func(message lsproto.RequestMessage) (any, error) {
				request, err := unmarshalInto(message.Params, &lsproto.DocumentDiagnosticsParams{})
				if err != nil {
					return nil, err
				}

				diagnostics := forest.Harvest()
				thisDocDiag := diagnostics[request.TextDocument.Uri]
				delete(diagnostics, request.TextDocument.Uri)

				response := lsproto.NewFullDocumentDiagnosticResponse(
					message.ID,
					thisDocDiag,
					diagnostics,
				)

				return response, nil
			},
			"initialize": func(message lsproto.RequestMessage) (any, error) {
				root := message.Params["rootPath"].(string)
				collectVocabFilesInRootPath(root)

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
							// "hoverProvider": true,
							"diagnosticProvider": map[string]any{
								// a change of date in one vocab can affect another (spaced repetition)
								"interFileDependencies": true,
							},
							// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#workspace_didChangeWatchedFiles
							"workspace": map[string]any{
								"fileOperations": map[string]any{
									"didDelete": map[string]any{
										"filters": []map[string]any{
											{
												"scheme":  "file",
												"pattern": map[string]any{"glob": "**/*.vocab"},
											},
										},
									},
								},
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

func diagnosticsToNotificationResponse(uri string, version float64, diags []lsproto.Diagnostic) *lsproto.PublishDiagnosticsNotification {
	return lsproto.NewPublishDiagnosticsNotfication(
		lsproto.PublishDiagnosticsParams{
			Uri:         uri,
			Version:     version,
			Diagnostics: diags,
		},
	)
}
