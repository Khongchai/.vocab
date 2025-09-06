package engine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	lsproto "vocab/lsp"
)

type Engine struct {
	ctx    context.Context
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewEngine(ctx context.Context, reader io.Reader, writer io.Writer) *Engine {
	r := bufio.NewReader(reader)
	w := bufio.NewWriter(writer)
	engine := &Engine{
		ctx:    ctx,
		reader: r,
		writer: w,
	}
	return engine
}

// ref:
// https://github.com/microsoft/typescript-go/blob/main/internal/lsp/lsproto/baseproto.go#L31

var (
	ErrInvalidHeader        = errors.New("lsp: invalid header")
	ErrInvalidContentLength = errors.New("lsp: invalid content length")
	ErrNoContentLength      = errors.New("lsp: no content length")
)

// Read the content of a json-rpc formatted message.
// Returns the byte array containing the content of that json rpc request
func (engine *Engine) readInput() ([]byte, error) {
	var contentLength int64

	// parses content length
	for {
		line, err := engine.reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("lsp: read header: %w", err)
		}

		if bytes.Equal(line, []byte("\r\n")) {
			break
		}

		key, value, ok := bytes.Cut(line, []byte(":"))
		if !ok {
			return nil, fmt.Errorf("%w: %q", ErrInvalidHeader, line)
		}

		if bytes.Equal(key, []byte("Content-Length")) {
			contentLength, err = strconv.ParseInt(string(bytes.TrimSpace(value)), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: parse error: %w", ErrInvalidContentLength, err)
			}
			if contentLength < 0 {
				return nil, fmt.Errorf("%w: negative value %d", ErrInvalidContentLength, contentLength)
			}
		}
	}

	if contentLength <= 0 {
		return nil, ErrNoContentLength
	}

	// parses json body
	data := make([]byte, contentLength)
	if _, err := io.ReadFull(engine.reader, data); err != nil {
		return nil, fmt.Errorf("lsp: read content: %w", err)
	}

	return data, nil
}

// Blocks and read content of this json rpc message
func (engine *Engine) read() (*lsproto.Message, error) {
	bytes, err := engine.readInput()
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
		data, err := engine.read()

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
					if _, err := fmt.Fprintf(engine.writer, "Content-Length: %d\r\n\r\n", len(out)); err != nil {
						panic("wtf")
					}
					if _, err := engine.writer.Write(out); err != nil {
						panic("wtf bro")
					}
					engine.writer.Flush()
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

				if _, err := fmt.Fprintf(engine.writer, "Content-Length: %d\r\n\r\n", len(out)); err != nil {
					panic("wtf")
				}
				if _, err := engine.writer.Write(out); err != nil {
					panic("wtf bro")
				}
				engine.writer.Flush()
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

func (engine *Engine) handle() {

}
