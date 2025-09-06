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
)

type DocumentDiagnosticReportKind string

const (
	DocumentDiagnosticReportKindFull      DocumentDiagnosticReportKind = "full"
	DocumentDiagnosticReportKindUnchanged DocumentDiagnosticReportKind = "unchanged"
)

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
