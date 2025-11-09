package lsproto

// https://github.com/microsoft/typescript-go/blob/0a3c816da9be581f3b567df9f05b73533f5c9384/internal/lsp/lsproto/baseproto.go#L106

type ErrorCode struct {
	Name string
	Code int32
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

type DiagnosticsSeverity int

const (
	DiagnosticsSeverityError DiagnosticsSeverity = iota + 1
	DiagnosticsSeverityWarning
	DiagnosticsSeverityInformation
	DiagnosticsSeverityHint
	DiagnosticsSeverityNone
)

type DocumentDiagnosticReportKind string

const (
	DocumentDiagnosticReportKindFull      DocumentDiagnosticReportKind = "full"
	DocumentDiagnosticReportKindUnchanged DocumentDiagnosticReportKind = "unchanged"
)

type OptionalVersionedTextDocumentIdentifier struct {
	Uri     string  `json:"uri"`
	Version float64 `json:"version,omitempty"`
}

type TextDocumentContentChangeEvent struct {
	// If the document is a full change, range is null
	Range *Range `json:"range,omitempty"`
	Text  string `json:"text"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   OptionalVersionedTextDocumentIdentifier `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent        `json:"contentChanges"`
}

type TextDocumentItem struct {
	Uri        string  `json:"uri"`
	LanguageId string  `json:"langaugeId"`
	Version    float64 `json:"version"`
	Text       string  `json:"text"`
}

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}

type FileDelete struct {
	Uri string `json:"uri"`
}

type DeleteFilesParms struct {
	Files []*FileDelete `json:"files"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier
}

type DidOpenDocumentParams struct {
	TextDocument *TextDocumentItem
}

func NewTextDocumentHoverResponse(requestId int, content string, r *Range) *map[string]any {
	return NewGenericResponse(
		requestId,
		map[string]any{
			"contents": map[string]any{
				"kind":  "plaintext",
				"value": content,
			},
			"range": r,
		},
	)
}

func NewGenericResponse(messageId int, result map[string]any) *map[string]any {
	return &map[string]any{
		"jsonrpc": JsonRPCVersion,
		"id":      messageId,
		"result":  result,
	}
}

type HoverParams struct {
	TextDocument string   `json:"textDocument"`
	Position     Position `json:"position"`
}

func NewFullDocumentDiagnosticResponse(id int, documentsDiagnostics []Diagnostic, relatedDocumentsDiagnostics map[string][]Diagnostic) *documentDiagnosticResponse {
	reports := map[string]FullDocumentDiagnosticReport{}

	for key := range relatedDocumentsDiagnostics {
		reports[key] = FullDocumentDiagnosticReport{
			Kind:  DocumentDiagnosticReportKindFull,
			Items: relatedDocumentsDiagnostics[key],
		}
	}

	return &documentDiagnosticResponse{
		Jsonrpc: JsonRPCVersion,
		ID:      id,
		Result: DocumentDiagnosticReport{
			Kind:             DocumentDiagnosticReportKindFull,
			Items:            documentsDiagnostics,
			RelatedDocuments: reports,
		},
	}
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_diagnostic
type documentDiagnosticResponse struct {
	Jsonrpc string                   `json:"jsonrpc"`
	ID      int                      `json:"id"`
	Result  DocumentDiagnosticReport `json:"result"`
}

type DocumentDiagnosticReport struct {
	Kind             DocumentDiagnosticReportKind            `json:"kind"`
	Items            []Diagnostic                            `json:"items"`
	RelatedDocuments map[string]FullDocumentDiagnosticReport `json:"relatedDocuments,omitempty"`
}

type FullDocumentDiagnosticReport struct {
	Kind  DocumentDiagnosticReportKind `json:"kind"`
	Items []Diagnostic                 `json:"items"`
}

func NewPublishDiagnosticsNotfication(params PublishDiagnosticsParams) *PublishDiagnosticsNotification {
	return &PublishDiagnosticsNotification{
		Jsonrpc: JsonRPCVersion,
		Method:  "textDocument/publishDiagnostics",
		Params:  params,
	}
}

type PublishDiagnosticsNotification struct {
	Jsonrpc string                   `json:"jsonrpc"`
	Method  string                   `json:"method"`
	Params  PublishDiagnosticsParams `json:"params"`
}

type PublishDiagnosticsParams struct {
	Uri         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Version     float64      `json:"version,omitempty"`
}

type TextDocument struct {
	Uri string `json:"uri"`
}

type DocumentDiagnosticsParams struct {
	TextDocument TextDocument `json:"textDocument"`
}

type CollectParams struct {
	CurrentDocumentUri string `json:"currentDocumentUri"`
}

func NewCollectResponse(requestId int, itWords []string, deWords []string) *map[string]any {
	return NewGenericResponse(
		requestId,
		map[string]any{
			"words": map[string]any{
				"it": itWords,
				"de": deWords,
			},
		},
	)
}

type Diagnostic struct {
	Range    Range               `json:"range"`
	Message  string              `json:"message,omitempty"`
	Severity DiagnosticsSeverity `json:"severity"`
}

func MakeDiagnostics(message string, line int, startPos int, endPos int, level DiagnosticsSeverity) *Diagnostic {
	return &Diagnostic{
		Message:  message,
		Severity: level,
		Range: Range{
			Start: Position{
				Line:      line,
				Character: startPos,
			},
			End: Position{
				Line:      line,
				Character: endPos,
			},
		},
	}
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

var (
	// Defined by JSON-RPC
	ErrParseError     = &ErrorCode{"ParseError", -32700}
	ErrInvalidRequest = &ErrorCode{"InvalidRequest", -32600}
	ErrMethodNotFound = &ErrorCode{"MethodNotFound", -32601}
	ErrInvalidParams  = &ErrorCode{"InvalidParams", -32602}
	ErrInternalError  = &ErrorCode{"InternalError", -32603}

	// Error code indicating that a server received a notification or
	// request before the server has received the `initialize` request.
	ErrServerNotInitialized = &ErrorCode{"ServerNotInitialized", -32002}
	ErrUnknownErrorCode     = &ErrorCode{"UnknownErrorCode", -32001}

	// A request failed but it was syntactically correct, e.g the
	// method name was known and the parameters were valid. The error
	// message should contain human readable information about why
	// the request failed.
	ErrRequestFailed = &ErrorCode{"RequestFailed", -32803}

	// The server cancelled the request. This error code should
	// only be used for requests that explicitly support being
	// server cancellable.
	ErrServerCancelled = &ErrorCode{"ServerCancelled", -32802}

	// The server detected that the content of a document got
	// modified outside normal conditions. A server should
	// NOT send this error code if it detects a content change
	// in it unprocessed messages. The result even computed
	// on an older state might still be useful for the client.
	//
	// If a client decides that a result is not of any use anymore
	// the client should cancel the request.
	ErrContentModified = &ErrorCode{"ContentModified", -32801}

	// The client has canceled a request and a server has detected
	// the cancel.
	ErrRequestCancelled = &ErrorCode{"RequestCancelled", -32800}
)
