// nice refs
// https://github.com/microsoft/typescript-go/blob/0a3c816da9be581f3b567df9f05b73533f5c9384/internal/lsp/lsproto/jsonrpc.go#L221
// https://www.jsonrpc.org/specification#request_object

package lsproto

import (
	"encoding/json"
	"fmt"
)

const jsonRPCVersion = `"2.0"`

type RequestMethod = string

const (
	Initialize RequestMethod = "initialize"
)

type MessageKind int

const (
	MessageKindNotification MessageKind = iota
	MessageKindRequest
	MessageKindResponse
)

type Message struct {
	Kind MessageKind
	msg  any
}

type Notification struct {
	Method RequestMethod `json:"method"`
	Params any           `json:"params"`
}

type RequestMessage struct {
	ID     int           `json:"id"`
	Method RequestMethod `json:"method"`
	Params any           `json:"params,omitempty"`
}

type ResponseMessage struct {
	ID     int `json:"id,omitempty"`
	Result any `json:"result,omitempty"`
	Error  any `json:"error,omitempty"`
}

func UnmarshalJson(raw []byte) (*Message, error) {
	var out map[string]interface{} = nil
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("%#v: %w", ErrInvalidRequest, err)
	}

	if out["id"] == nil {
		var notification *Notification
		json.Unmarshal(raw, notification)
		return &Message{
			Kind: MessageKindNotification,
			msg:  notification,
		}, nil
	}

	if out["method"] != nil {
		var request *RequestMessage
		json.Unmarshal(raw, request)
		return &Message{
			Kind: MessageKindRequest,
			msg:  request,
		}, nil
	}

	var response *ResponseMessage
	json.Unmarshal(raw, response)
	return &Message{
		Kind: MessageKindResponse,
		msg:  response,
	}, nil
}
