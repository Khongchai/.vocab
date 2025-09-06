package engine

import (
	"context"
	"fmt"
	"os"
	"vocab/lib"
	lsproto "vocab/lsp"
)

type ReadCallback = func() ([]byte, error)
type WriteCallback = func(any)

type Engine struct {
	ctx    context.Context
	read   ReadCallback
	write  WriteCallback
	logger lib.Logger
}

func NewEngine(ctx context.Context, read ReadCallback, write WriteCallback, logger lib.Logger) *Engine {
	engine := &Engine{
		ctx,
		read,
		write,
		logger,
	}
	return engine
}

// Blocks and read content of this json rpc message
func (engine *Engine) readNext() (*lsproto.Message, error) {
	bytes, err := engine.read()
	if err != nil {
		return nil, err
	}

	message, err := lsproto.UnmarshalJson(bytes)
	if err != nil {
		return nil, fmt.Errorf("%#v: %w", lsproto.ErrInvalidRequest, err)
	}

	return message, nil
}

// Start up main loop.
func (engine *Engine) Start() {
	for { // https://github.com/microsoft/typescript-go/blob/main/internal/lsp/server.go#L246
		data, err := engine.readNext()

		if err != nil {
			engine.logger.Log("Decode error: ", err)
			fmt.Fprintln(os.Stderr, "decode error:", err)
			continue
		}

		switch data.Kind {
		case lsproto.MessageKindNotification:
			if n, ok := data.Msg.(lsproto.Notification); ok {
				engine.logger.Log("Received notification ", n.Method)
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

					engine.write(response)
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

				engine.write(response)
			}
		case lsproto.MessageKindResponse:
			if r, ok := data.Msg.(lsproto.ResponseMessage); ok {
				engine.logger.Log(r.ID)
			}
		default:
			engine.logger.Log("No default message handler found.")
		}
	}

}

func (engine *Engine) handle() {

}
