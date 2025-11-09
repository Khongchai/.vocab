package harvester

import (
	"vocab/lib"
	lsproto "vocab/lsp"
	"vocab/vocabulary/forest"
)

type NotificationWorker struct {
	forest *forest.Forest
}

func diagnosticsToNotificationResponse(uri string, version float64, diags []lsproto.Diagnostic) *lsproto.PublishDiagnosticsNotification {
	return lsproto.NewPublishDiagnosticsNotfication(
		lsproto.PublishDiagnosticsParams{
			Uri:         uri,
			Version:     version,
			Diagnostics: diags,
		},
	)
}

func NewNotificationWorker(f *forest.Forest) *NotificationWorker {
	return &NotificationWorker{
		forest: f,
	}
}

func (n *NotificationWorker) DeleteFileWorker(request lsproto.Notification) (any, error) {
	params, err := lib.UnmarshalInto(request.Params, &lsproto.DeleteFilesParms{})
	if err != nil {
		return nil, err
	}

	for _, file := range params.Files {
		n.forest.Remove(file.Uri)
	}

	return nil, nil
}

func (n *NotificationWorker) DidOpenWorker(request lsproto.Notification) (any, error) {
	params, err := lib.UnmarshalInto(request.Params, &lsproto.DidOpenDocumentParams{})
	if err != nil {
		return nil, err
	}

	n.forest.Plant(params.TextDocument.Uri, params.TextDocument.Text, nil)

	return diagnosticsToNotificationResponse(
		params.TextDocument.Uri,
		params.TextDocument.Version,
		nil,
	), nil
}

func (n *NotificationWorker) DidChangeWorker(request lsproto.Notification) (any, error) {
	params, err := lib.UnmarshalInto(request.Params, &lsproto.DidChangeTextDocumentParams{})
	if err != nil {
		return nil, err
	}

	for i := range params.ContentChanges {
		change := params.ContentChanges[i]
		// for now, sequential. In the future we can make this parallel
		n.forest.Plant(params.TextDocument.Uri, change.Text, change.Range)
	}

	return nil, nil
}
