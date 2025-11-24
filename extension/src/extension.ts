import * as path from "path";
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";
import { registerCommands } from "./registerCommands";

let client: LanguageClient | undefined;

export async function activate(context: vscode.ExtensionContext) {
  const serverOptions = (() => {
    const lsPath = path.resolve(context.extensionPath, "binaries");
    const exePath = path.join(
      lsPath,
      `vocab-ls${process.platform === "win32" ? ".exe" : ""}`
    );
    const option: ServerOptions = {
      // Note: if we can't find the package during build, take a look at this
      // https://github.com/microsoft/typescript-go/blob/main/_packages/native-preview/lib/getExePath.js
      command: exePath,
      transport: TransportKind.stdio,
    };
    return option;
  })();

  const clientOptions: LanguageClientOptions = (() => {
    const outputChannel = vscode.window.createOutputChannel("vocab");
    const traceOutputChannel =
      vscode.window.createOutputChannel("vocab (trace)");
    return {
      outputChannel,
      traceOutputChannel,
      documentSelector: [
        {
          scheme: "file",
          language: "vocab",
        },
        {
          scheme: "untitled",
          language: "vocab",
        },
      ],
    };
  })();

  client = new LanguageClient(
    "vocab",
    "vocab language server",
    serverOptions,
    clientOptions
  );

  await client.start();
  await client.sendNotification("workspace/didChangeConfiguration", {
    settings: {
      "vscode-languageclient": {
        trace: { server: "verbose" },
      },
    },
  });

  context.subscriptions.push(...registerCommands(client));
}

export async function deactivate() {
  await client?.stop();
}
