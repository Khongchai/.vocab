package lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type InputReader struct {
	reader *bufio.Reader
}

// ref:
// https://github.com/microsoft/typescript-go/blob/main/internal/lsp/lsproto/baseproto.go#L31

var (
	ErrInvalidHeader        = errors.New("lsp: invalid header")
	ErrInvalidContentLength = errors.New("lsp: invalid content length")
	ErrNoContentLength      = errors.New("lsp: no content length")
)

func NewInputReader(reader io.Reader) *InputReader {
	wrapped := bufio.NewReader(reader)
	return &InputReader{
		reader: wrapped,
	}
}

// Read the content of a json-rpc formatted message.
// Returns the byte array containing the content of that json rpc request
func (inputReader *InputReader) Read() ([]byte, error) {
	var contentLength int64

	// parses content length
	for {
		line, err := inputReader.reader.ReadBytes('\n')
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
	if _, err := io.ReadFull(inputReader.reader, data); err != nil {
		return nil, fmt.Errorf("lsp: read content: %w", err)
	}

	return data, nil
}
