const path = require('path');
const os = require('os');
// const { workspace, languages, ExtensionContext } = require('vscode');
const { LanguageClient, TransportKind } = require('vscode-languageclient/node');

function activate(context) {
    // Figure out the server binary based on the platform

    const serverModule = getServerPath();
 
    // Get the log level from the kuneiform extension configuration
    // const config = workspace.getConfiguration('kuneiform');
    // const logLevel = config.get('logLevel', 'info');

    let serverOptions = {
        run: {
            command: serverModule,
            // args: ['-loglevel', logLevel],
            transport: TransportKind.stdio
        },
        debug: {
            command: serverModule,
            // args: ['-loglevel', logLevel],
            transport: TransportKind.stdio
        }
    };

    let clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'kuneiform' }]
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

function getServerPath() {
    const platform = os.platform();
    const arch = os.arch();
    let binaryName;

    switch (platform) {
        case 'darwin':
            binaryName = arch === 'arm64' ? 'kuneiform-lsp-darwin-arm64' : 'kuneiform-lsp-darwin-amd64';
            break;
        case 'linux':
            binaryName = arch === 'arm64' ? 'kuneiform-lsp-linux-arm64' : 'kuneiform-lsp-linux-amd64';
            break;
        case 'win32':
            binaryName = arch === 'arm64' ? 'kuneiform-lsp-windows-arm64.exe' : 'kuneiform-lsp-windows-amd64.exe';
            break;
        default:
            throw new Error(`Unsupported platform: ${platform}-${arch}`);
    }

    return path.join(__dirname, '..', 'server', '.build', binaryName);
}