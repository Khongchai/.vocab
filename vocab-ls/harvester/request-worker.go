package harvester

import (
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

type RequestWorker struct {
	forest *forest.Forest
	logger lib.Logger
}

func NewRequestWorker(f *forest.Forest, logger lib.Logger) *RequestWorker {
	return &RequestWorker{
		forest: f,
		logger: logger,
	}
}

func (n *RequestWorker) CollectFromThisFileWorker(rm lsproto.RequestMessage) (any, error) {
	params, err := lib.UnmarshalInto(rm.Params, &lsproto.CollectParams{})
	if err != nil {
		return nil, err
	}

	harvested := n.forest.Harvest()
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
}

func (n *RequestWorker) CollectFromAllFilesWorker(rm lsproto.RequestMessage) (any, error) {
	harvesteds := n.forest.Harvest()
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
}

func (n *RequestWorker) TextDocumentDiagnosticsWorker(message lsproto.RequestMessage) (any, error) {
	request, err := lib.UnmarshalInto(message.Params, &lsproto.DocumentDiagnosticsParams{})
	if err != nil {
		return nil, err
	}

	diagnostics := n.forest.Harvest()
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
}

func (n *RequestWorker) InitializeWorker(message lsproto.RequestMessage) (any, error) {
	root := message.Params["rootPath"].(string)
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
			n.logger.Logf("Can't read content at %s", root)
		}
		fileContent := string(bytes)
		n.forest.Plant(fmt.Sprintf("%s%s", "file://", path), fileContent, nil)

		return nil
	})

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

}
