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
	ctx                  context.Context
	read                 ReadCallback
	write                WriteCallback
	logger               lib.Logger
	notificationHandlers map[string]func(lsproto.Notification) any
	requestHandlers      map[string]func(lsproto.RequestMessage) any
}

func NewEngine(
	ctx context.Context,
	read ReadCallback,
	write WriteCallback,
	logger lib.Logger,
	notificationHandlers map[string]func(lsproto.Notification) any,
	requestHandlers map[string]func(lsproto.RequestMessage) any,
) *Engine {
	engine := &Engine{
		ctx,
		read,
		write,
		logger,
		notificationHandlers,
		requestHandlers,
	}
	return engine
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

				response := engine.onNotification(n)
				if response == nil {
					continue
				}
				engine.write(response)
			}
		case lsproto.MessageKindRequest:
			if r, ok := data.Msg.(lsproto.RequestMessage); ok {
				response := engine.onRequest(r)
				if response == nil {
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

func (engine *Engine) onRequest(message lsproto.RequestMessage) any {
	handler := engine.requestHandlers[message.Method]
	if handler == nil {
		return nil
	}
	result := handler(message)
	return result
}

func (engine *Engine) onNotification(message lsproto.Notification) any {
	handler := engine.notificationHandlers[message.Method]
	if handler == nil {
		return nil
	}
	result := handler(message)
	return result
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
