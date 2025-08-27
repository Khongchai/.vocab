Notes for how go lsp for vscode is implemented

# Extension install gopl when?

# How does node extension process starts the go server

1.
https://github.com/golang/vscode-go/blob/85d7f0ca21fc18762ba9f7981de0f7c9a197d572/extension/src/goMain.ts#L140

2. 
https://github.com/golang/vscode-go/blob/master/extension/src/language/goLanguageServer.ts#L383

3.
https://github.com/golang/vscode-go/blob/master/extension/src/commands/startLanguageServer.ts#L88

# How do they communicate?

vscode-languageserver-node does (GoLanguageClient extends that one.)



TODO

- [ ] Try connecting go vocab server and vscode 