import { LanguageClient } from "vscode-languageclient/node";
import * as vscode from "vscode";

export function registerCommands(client: LanguageClient) {
  const disposables: vscode.Disposable[] = [];

  disposables.push(
    vscode.commands.registerCommand("vocab.collectFromThisFile", async () => {
      const response = await client.sendRequest("vocab/collectFromThisFile");
      debugger;
    })
  );

  disposables.push(
    vscode.commands.registerCommand("vocab.collectAll", async () => {
      const response = await client.sendRequest("vocab/collectAll");
      debugger;
    })
  );

  return disposables;
}
