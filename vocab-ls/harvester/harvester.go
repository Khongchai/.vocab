package harvester

import (
	"context"
	"vocab/lib"
	lsproto "vocab/lsp"
	"vocab/vocabulary/forest"
)

type Harvester struct {
	engine             *Engine
	notificationWorker *NotificationWorker
	requestWorker      *RequestWorker
}

func NewHarvester(
	ctx context.Context,
	forest *forest.Forest,
	readCallback func() ([]byte, error),
	writeCallback func(msg any),
	logger lib.Logger,
) *Harvester {
	h := &Harvester{
		engine:             NewEngine(ctx, readCallback, writeCallback, logger),
		notificationWorker: NewNotificationWorker(forest),
		requestWorker:      NewRequestWorker(forest, logger),
	}

	h.engine.SetNotificationHandlers(map[string]func(lsproto.Notification) (any, error){
		"workspace/didDeleteFiles": h.notificationWorker.DeleteFileWorker,
		"textDocument/didOpen":     h.notificationWorker.DidOpenWorker,
		"textDocument/didChange":   h.notificationWorker.DidChangeWorker,
	}).SetRequestHandlers(map[string]func(lsproto.RequestMessage) (any, error){
		"vocab/collectFromThisFile": h.requestWorker.CollectFromThisFileWorker,
		"vocab/collectAll":          h.requestWorker.CollectFromAllFilesWorker,
		"textDocument/diagnostic":   h.requestWorker.TextDocumentDiagnosticsWorker,
		"initialize":                h.requestWorker.InitializeWorker,
	})

	return h
}

func (h *Harvester) Start() {
	h.engine.Start()
}
