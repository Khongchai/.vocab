package jsonrpc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	lsproto "vocab/lsp"

	"github.com/go-json-experiment/json"
)

// ref:
// https://github.com/microsoft/typescript-go/blob/main/internal/lsp/lsproto/baseproto.go#L31

var (
	ErrInvalidHeader        = errors.New("lsp: invalid header")
	ErrInvalidContentLength = errors.New("lsp: invalid content length")
	ErrNoContentLength      = errors.New("lsp: no content length")
)

type Jsonrpc struct {
	instance *bufio.Reader
}

func NewJsonrpc(reader io.Reader) *Jsonrpc {
	instance := bufio.NewReader(reader)
	return &Jsonrpc{
		instance,
	}
}

// Blocks and read content of this json rpc message
func (r *Jsonrpc) Read() (map[string]any, error) {
	bytes, err := r.ReadBody()
	if err != nil {
		return nil, err
	}

	req := &lsproto.Request{}
	if err := json.Unmarshal(bytes, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	return req, nil
}

// Read the content of a json-rpc formatted message.
// Returns the byte array containing the content of that json rpc request
func (r *Jsonrpc) ReadBody() ([]byte, error) {
	var contentLength int64

	// parses content length
	for {
		line, err := r.instance.ReadBytes('\n')
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
	if _, err := io.ReadFull(r.instance, data); err != nil {
		return nil, fmt.Errorf("lsp: read content: %w", err)
	}

	return data, nil
}
