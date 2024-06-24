package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kwilteam/kwil-db/parse"
	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

type lspHandler struct {
	docs     map[string]*kfDocs
	logger   *slog.Logger
	handlers map[string]func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request)
}

type Handler func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request)

func (l *lspHandler) registerHandlers() {
	l.handlers = map[string]func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request){
		"initialize":                  l.handleInitialize,
		"textDocument/didOpen":        l.handleDidOpen,
		"textDocument/didChange":      l.handleDidChange,
		"textDocument/didClose":       l.handleDidClose,
		"shutdown":                    l.handleShutdown,
		"$/cancelRequest":             l.handleCancelRequest,
		"textDocument/documentSymbol": l.handleDocumentSymbol,
		"textDocument/completion":     l.handleCompletion,
		// "textDocument/semanticTokens/full": l.handleSemanticTokens,
		// "completionItem/resolve":           l.handleCompletionItemResolve,
	}
}

func (l *lspHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	l.logger.Info("Received request: ", req.Method, req.ID)
	if handler, ok := l.handlers[req.Method]; ok {
		handler(ctx, conn, req)
	} else {
		l.logger.Info("Unknown request method: ", req.Method, req.ID)
	}
}

func (l *lspHandler) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	params := lsp.InitializeParams{}
	json.Unmarshal(*req.Params, &params)
	kind := lsp.TDSKFull
	res := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Kind: &kind,
			},
			// DocumentSymbolProvider: true,
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: triggerKeywords,
			},
			// HoverProvider: true,
		},
	}
	conn.Reply(ctx, req.ID, &res)
}

func (l *lspHandler) handleDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	params := lsp.DidOpenTextDocumentParams{}
	json.Unmarshal(*req.Params, &params)

	docID := string(params.TextDocument.URI)
	docText := params.TextDocument.Text

	if _, ok := l.docs[docID]; !ok {
		l.docs[docID] = &kfDocs{rawKf: docText}
	}

	_, diagnostics := l.validateKfDocument(docID, docText)
	conn.Notify(ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (l *lspHandler) handleDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	params := lsp.DidChangeTextDocumentParams{}
	json.Unmarshal(*req.Params, &params)

	if len(params.ContentChanges) != 1 {
		l.logger.Error("Expected exactly one change, got ", slog.Int("content changes", len(params.ContentChanges)))
		panic("Should be exactly one change")
	}

	docID := string(params.TextDocument.URI)
	docText := params.ContentChanges[0].Text

	doc, ok := l.docs[docID]
	if !ok {
		l.docs[docID] = &kfDocs{rawKf: docText}
	} else {
		doc.rawKf = docText // does this update the value in the map?
	}

	_, diagnostics := l.validateKfDocument(docID, docText)
	conn.Notify(ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (l *lspHandler) handleDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	params := lsp.DidCloseTextDocumentParams{}
	json.Unmarshal(*req.Params, &params)
	delete(l.docs, string(params.TextDocument.URI))
}

func (l *lspHandler) handleShutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	conn.Notify(ctx, "exit", req.ID, nil)
}

func (l *lspHandler) handleCancelRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// TODO: Should we do anything here?
}

func (l *lspHandler) handleDocumentSymbol(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// TODO: Extract action, procedures, tables, extension methods, foreign procedures, etc.
	// params := lsp.DocumentSymbolParams{}
	// json.Unmarshal(*req.Params, &params)
}

func (l *lspHandler) handleCompletion(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	params := lsp.CompletionParams{}
	err := json.Unmarshal(*req.Params, &params)
	if err != nil {
		l.logger.Error("error unmarshalling completion params: %v", err)
		return
	}

	docID := string(params.TextDocument.URI)
	doc, ok := l.docs[docID]
	if !ok {
		l.logger.Error("document not found: %s", slog.String("docID", docID))
		return
	}

	_, ds := l.validateKfDocument(docID, doc.rawKf)
	conn.Notify(ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: ds,
	})

	offset, err := l.getOffset(doc.rawKf, params.Position.Line, params.Position.Character)
	if err != nil {
		l.logger.Error("Error getting completionTriggerCharacter offset: ", err)
		return
	}

	items := l.getCompletionItems(doc.parsedSchema, offset)
	l.printSuggestions(items)
	conn.Reply(ctx, req.ID, items)
}

func (l *lspHandler) printSuggestions(items []lsp.CompletionItem) {
	// Optimization: Skip if log level is info, warn or error
	if logLevel.Level() >= slog.LevelInfo {
		return
	}

	var labels []string
	for _, item := range items {
		labels = append(labels, item.Label)
	}
	l.logger.Debug("Suggestions: ", slog.Any("", labels))
}

func (l *lspHandler) validateKfDocument(uri string, text string) (*parse.SchemaParseResult, []lsp.Diagnostic) {
	res, err := parse.ParseAndValidate([]byte(text))
	if err != nil {
		return nil, []lsp.Diagnostic{
			{
				Severity: lsp.Error,
				Message:  err.Error(),
			},
		}
	}

	if len(res.ParseErrs.Errors()) == 0 {
		doc, ok := l.docs[uri]
		if !ok {
			l.docs[uri] = &kfDocs{rawKf: text}
		}
		doc.parsedSchema = res
	}

	return res, getDiagnostics(res)
}

func (l *lspHandler) getOffset(text string, line, col int) (int, error) {
	lines := strings.Split(text, "\n")
	if line > len(lines) {
		return 0, fmt.Errorf("line %d is out of bounds,  max: %d , overall text: %s", line, len(lines), text)
	}

	offset := 0
	for i := 0; i < line-1; i++ {
		offset += len(lines[i]) + 1
	}
	return offset + col, nil
}
