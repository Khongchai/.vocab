package harvester

import (
	"context"
	"errors"
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
	notificationHandlers map[string]func(lsproto.Notification) (any, error)
	requestHandlers      map[string]func(lsproto.RequestMessage) (any, error)
}

func NewEngine(
	ctx context.Context,
	read ReadCallback,
	write WriteCallback,
	logger lib.Logger,
	notificationHandlers map[string]func(lsproto.Notification) (any, error),
	requestHandlers map[string]func(lsproto.RequestMessage) (any, error),
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

				response, err := engine.onNotification(n)
				if err != nil {
					engine.logger.Logf("Got error while handling message %+v", err)
				}
				if response == nil {
					continue
				}
				engine.write(response)
			}
		case lsproto.MessageKindRequest:
			if r, ok := data.Msg.(lsproto.RequestMessage); ok {
				engine.logger.Log("Received request ", r.Method)

				response, err := engine.onRequest(r)
				if err != nil {
					engine.logger.Logf("Got error while handling request %+v", err)
				}
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

func (engine *Engine) onRequest(message lsproto.RequestMessage) (any, error) {
	handler := engine.requestHandlers[message.Method]
	if handler == nil {
		return nil, errors.New("request handler nil")
	}
	result, err := handler(message)
	return result, err
}

func (engine *Engine) onNotification(message lsproto.Notification) (any, error) {
	handler := engine.notificationHandlers[message.Method]
	if handler == nil {
		return nil, errors.New("notification handlers nil")
	}
	result, err := handler(message)
	return result, err
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
