const path = require('path');
const { workspace, languages, ExtensionContext } = require('vscode');
const { LanguageClient, TransportKind } = require('vscode-languageclient/node');

function activate(context) {
    let serverModule = path.join(__dirname, '..', 'server', 'kuneiform-lsp');

    // Get the log level from the kuneiform extension configuration
    const config = workspace.getConfiguration('kuneiform');
    const logLevel = config.get('logLevel', 'info');

    let serverOptions = {
        run: {
            command: serverModule,
            args: ['-loglevel', logLevel],
            transport: TransportKind.stdio
        },
        debug: {
            command: serverModule,
            args: ['-loglevel', logLevel],
            transport: TransportKind.stdio
        }
    };

    let clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'kuneiform' }],
        middleware: {
            provideDocumentSemanticTokens: (document, token) => {
                return client.sendRequest('textDocument/semanticTokens/full', {
                    textDocument: { uri: document.uri.toString() }
                }).then(response => {
                    return response;
                });
            }
        }
    };

    let client = new LanguageClient(
        'kuneiformLanguageServer',
        'Kuneiform Language Server',
        serverOptions,
        clientOptions
    );

    client.start();
}

exports.activate = activate;
