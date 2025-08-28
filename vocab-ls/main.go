package main

import (
	"context"
	"io"
	"os"
)

func main() {
	stream := protocol.NewJSONStream(os.Stdin, os.Stdout)
	conn := protocol.ServerDispatcher(stream)

	// Register handlers
	conn.Go(ctxHandler)
	select {}
}

func ctxHandler(ctx context.Context, conn protocol.Conn) error {
	for {
		msg, err := conn.Recv(ctx)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch req := msg.(type) {
		case *protocol.DidOpenTextDocumentParams:
			// Highlight the entire document red
			diagnostic := protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 9999, Character: 0}, // big range
				},
				Severity: protocol.SeverityError, // red squiggle
				Message:  "Everything is wrong :)",
			}
			_ = conn.Notify(ctx, "textDocument/publishDiagnostics", &protocol.PublishDiagnosticsParams{
				URI:         req.TextDocument.URI,
				Diagnostics: []protocol.Diagnostic{diagnostic},
			})
		}
	}
}
