package vocabulary

import lsproto "vocab/lsp"

type Handlers struct {
	OnNotification func(notification lsproto.Notification)
	OnRequest      func(request lsproto.RequestMessage)
	OnResponse     func(response lsproto.ResponseMessage)
}
