import * as vscode from "vscode";

export async function addNewWordSectionUseCase(
  document: vscode.TextDocument,
  words: {
    it: string[];
    de: string[];
  }
) {
  if (words.it.length + words.de.length === 0) {
    vscode.window.showInformationMessage("Nothing to review for now!");
    return;
  }

  let content = (() => {
    const contents: string[] = [];

    const today = new Date();
    contents.push(
      `${today.getDate().toString().padStart(2, "0")}/${(today.getMonth() + 1)
        .toString()
        .padStart(2, "0")}/${today.getFullYear()}`.trim()
    );

    if (words.it.length > 0) {
      contents.push(`>> (it) ${words.it.join(", ")}`);
    }

    if (words.de.length > 0) {
      contents.push(`>> (de) ${words.de.join(", ")}`);
    }

    const joined = ["\n", ...contents.join("\n")].join("");
    return joined;
  })();

  const lineCount = document.lineCount;
  const lastLine = document.lineAt(lineCount - 1);
  const endPosition = lastLine.range.end;
  const edit = new vscode.WorkspaceEdit();
  edit.insert(document.uri, endPosition, content);
  await vscode.workspace.applyEdit(edit);

  vscode.window.showInformationMessage(
    `Added ${words.it.length} it words and ${words.de.length} words`
  );
}
