package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sourcegraph/jsonrpc2"
)

const (
	lsDir   = ".kwil-ls"
	logFile = "lsp.log"
)

var logLevel = &slog.LevelVar{}

type stdioRWC struct{}

func (s *stdioRWC) Close() error {
	return nil
}

func (s *stdioRWC) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (s *stdioRWC) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func main() {
	ctx := context.Background()
	logLevel.Set(slog.LevelDebug)
	logger := getLogger(logLevel)

	// Initialize the language server  and register the handlers
	lshandler := &lspHandler{
		docs:     make(map[string]*kfDocs),
		handlers: make(map[string]func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request)),
		logger:   logger,
	}
	lshandler.registerHandlers()

	conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(&stdioRWC{}, jsonrpc2.VSCodeObjectCodec{}), lshandler)

	logger.Info("Connected to the client...")
	defer conn.Close()

	<-conn.DisconnectNotify()
}

func getLogger(level *slog.LevelVar) *slog.Logger {
	// Initialize the logger
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory")
		os.Exit(1)
	}

	logDir := filepath.Join(homedir, lsDir)
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating log directory %s", logDir)
		os.Exit(1)
	}

	logFile := filepath.Join(logDir, logFile)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating log file %s", logFile)
		os.Exit(1)
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
