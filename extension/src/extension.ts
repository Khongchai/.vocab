import * as path from "path";
import * as vscode from "vscode";
import { LanguageClient, TransportKind } from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
  const serverExe = {
    command: path.join(context.extensionPath, "red-ls"),
    transport: TransportKind.stdio,
  };

  client = new LanguageClient("redLS", "Red Language Server", serverExe, {
    documentSelector: [{ scheme: "file", language: "plaintext" }],
  });

  client.start();
}

export function deactivate() {
  if (!client) return undefined;
  return client.stop();
}
