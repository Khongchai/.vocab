Notes for how go lsp for vscode is implemented

# Extension install gopl when?


# How does node extension process starts the go server

1.
https://github.com/golang/vscode-go/blob/85d7f0ca21fc18762ba9f7981de0f7c9a197d572/extension/src/goMain.ts#L140

2. 
https://github.com/golang/vscode-go/blob/master/extension/src/language/goLanguageServer.ts#L383

3.
https://github.com/golang/vscode-go/blob/master/extension/src/commands/startLanguageServer.ts#L88

4. 
https://github.com/golang/vscode-go/blob/85d7f0ca21fc18762ba9f7981de0f7c9a197d572/extension/src/language/goLanguageServer.ts#L414

5.
https://github.com/microsoft/vscode-languageserver-node/blob/3412a17149850f445bf35b4ad71148cfe5f8411e/client/src/node/main.ts#L486

# How do they communicate?

https://github.com/microsoft/vscode-languageserver-node/tree/main/jsonrpc

# Notes
- gopl for how go server is implemented.
- vscode-go for glue and extension (vscode) side.
- vscode-lanagueserver-node for the internal of vscode-go.

TODO

- [ ] Try connecting go vocab server and vscode 