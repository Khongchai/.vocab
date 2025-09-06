package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/go-json-experiment/json"
)

type OutputWriter struct {
	writer *bufio.Writer
}

func NewOutputWriter(writer io.Writer) *OutputWriter {
	return &OutputWriter{
		writer: bufio.NewWriter(writer),
	}
}

func (writer *OutputWriter) Write(message any) {
	out, err := json.Marshal(message)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	if _, err := fmt.Fprintf(writer.writer, "Content-Length: %d\r\n\r\n", len(out)); err != nil {
		panic("wtf")
	}
	if _, err := writer.writer.Write(out); err != nil {
		panic("wtf bro")
	}
	writer.writer.Flush()
}
