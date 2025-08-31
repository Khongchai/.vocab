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
	Msg  any
}

type Notification struct {
	Method RequestMethod `json:"method"`
	Params any           `json:"params"`
}

type RequestMessage struct {
	ID     int           `json:"id"`
	Method RequestMethod `json:"method"`
	Params any           `json:"params,omitzero"`
}

type ResponseMessage struct {
	ID     int `json:"id,omitzero"`
	Result any `json:"result,omitzero"`
	Error  any `json:"error,omitzero"`
}

func UnmarshalJson(raw []byte) (*Message, error) {
	var out map[string]any = nil
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("%#v: %w", ErrInvalidRequest, err)
	}

	var kind MessageKind
	switch {
	case out["id"] == nil:
		kind = MessageKindNotification
	case out["method"] != nil:
		kind = MessageKindRequest
	default:
		kind = MessageKindResponse
	}

	message := &Message{
		Kind: kind,
		Msg:  out,
	}

	return message, nil
}
