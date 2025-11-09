package harvester

import (
	"context"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"vocab/lib"
	lsproto "vocab/lsp"
	"vocab/vocabulary/forest"
	"vocab/vocabulary/parser"
)

type Harvester struct {
	logger lib.Logger
	forest *forest.Forest
	engine *Engine
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

func NewHarvester(
	ctx context.Context,
	forest *forest.Forest,
	readCallback func() ([]byte, error),
	writeCallback func(msg any),
	logger lib.Logger,
) *Harvester {
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

	engine := NewEngine(ctx, readCallback, writeCallback, logger,
		map[string]func(lsproto.Notification) (any, error){
			"workspace/didDeleteFiles": func(request lsproto.Notification) (any, error) {
				params, err := lib.UnmarshalInto(request.Params, &lsproto.DeleteFilesParms{})
				if err != nil {
					return nil, err
				}

				for _, file := range params.Files {
					forest.Remove(file.Uri)
				}

				return nil, nil
			},
			"textDocument/didOpen": func(rm lsproto.Notification) (any, error) {
				params, err := lib.UnmarshalInto(rm.Params, &lsproto.DidOpenDocumentParams{})
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
				params, err := lib.UnmarshalInto(rm.Params, &lsproto.DidChangeTextDocumentParams{})
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
			"vocab/collectFromThisFile": func(rm lsproto.RequestMessage) (any, error) {
				params, err := lib.UnmarshalInto(rm.Params, &lsproto.CollectParams{})
				if err != nil {
					return nil, err
				}

				harvested := forest.Harvest()
				thisDocInfo := harvested[params.CurrentDocumentUri]

				itWordSet := make(map[string]struct{})
				deWordSet := make(map[string]struct{})
				for _, harvested := range thisDocInfo {
					if harvested.Diagnostic.Severity != lsproto.DiagnosticsSeverityError {
						continue
					}
					if harvested.Lang == parser.Deutsch {
						_, found := deWordSet[harvested.Word]
						if found {
							continue
						}
						deWordSet[harvested.Word] = struct{}{}
					} else {
						_, found := itWordSet[harvested.Word]
						if found {
							continue
						}
						itWordSet[harvested.Word] = struct{}{}
					}
				}

				return lsproto.NewCollectResponse(rm.ID,
					slices.Collect(maps.Keys(itWordSet)),
					slices.Collect(maps.Keys(deWordSet)),
				), nil
			},
			"vocab/collectAll": func(rm lsproto.RequestMessage) (any, error) {

				harvesteds := forest.Harvest()
				deWordSet := make(map[string]struct{})
				itWordSet := make(map[string]struct{})

				for _, diagnostics := range harvesteds {
					for _, harvested := range diagnostics {
						if harvested.Diagnostic.Severity != lsproto.DiagnosticsSeverityError {
							continue
						}
						if harvested.Lang == parser.Deutsch {
							_, found := deWordSet[harvested.Word]
							if found {
								continue
							}
							deWordSet[harvested.Word] = struct{}{}
						} else {
							_, found := itWordSet[harvested.Word]
							if found {
								continue
							}
							itWordSet[harvested.Word] = struct{}{}
						}
					}
				}

				return lsproto.NewCollectResponse(rm.ID,
					slices.Collect(maps.Keys(itWordSet)),
					slices.Collect(maps.Keys(deWordSet)),
				), nil
			},
			"textDocument/hover": func(rm lsproto.RequestMessage) (any, error) {
				// not used yet
				_, err := lib.UnmarshalInto(rm.Params, &lsproto.HoverParams{})
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
				request, err := lib.UnmarshalInto(message.Params, &lsproto.DocumentDiagnosticsParams{})
				if err != nil {
					return nil, err
				}

				diagnostics := forest.Harvest()
				var thisDocDiags []lsproto.Diagnostic
				for _, d := range diagnostics[request.TextDocument.Uri] {
					thisDocDiags = append(thisDocDiags, d.Diagnostic)
				}

				delete(diagnostics, request.TextDocument.Uri)
				restDiags := make(map[string][]lsproto.Diagnostic)
				for key := range diagnostics {
					for _, d := range diagnostics[key] {
						restDiags[key] = append(restDiags[key], d.Diagnostic)
					}
				}

				response := lsproto.NewFullDocumentDiagnosticResponse(
					message.ID,
					thisDocDiags,
					restDiags,
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
	return &Harvester{
		logger: logger,
		engine: engine,
		forest: forest,
	}
}

func (h *Harvester) Start() {
	h.engine.Start()
}
