import { LanguageClient } from "vscode-languageclient/node";
import * as vscode from "vscode";
import { addNewWordSectionUseCase } from "./usecases/addNewWordSection.usecase";

export function registerCommands(client: LanguageClient) {
  const disposables: vscode.Disposable[] = [];

  type CollectResponse = {
    words: {
      it: string[];
      de: string[];
    };
  };

  disposables.push(
    vscode.commands.registerCommand("vocab.collectFromThisFile", async () => {
      const editor = vscode.window.activeTextEditor;
      if (!editor) {
        return;
      }

      const response = (await client.sendRequest("vocab/collectFromThisFile", {
        currentDocumentUri: editor.document.uri.toString(),
      })) as CollectResponse;

      await addNewWordSectionUseCase(editor.document, response.words);
    })
  );

  disposables.push(
    vscode.commands.registerCommand("vocab.collectAll", async () => {
      const editor = vscode.window.activeTextEditor;
      if (!editor) {
        return;
      }

      const response = (await client.sendRequest(
        "vocab/collectAll"
      )) as CollectResponse;

      await addNewWordSectionUseCase(editor.document, response.words);
    })
  );

  return disposables;
}
