package lsp

import (
	"context"
	"hybroid/core"
	"hybroid/walker"
	"io"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (c stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (c stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}

func Init(debug bool) {
	core.IsDebug = debug
	if !core.IsDebug {
		log.SetOutput(io.Discard)
	}

	if core.IsDebug {
		// Resolve where the debug log should go. The full
		// precedence and fallback chain lives in logpath.go;
		// see the resolveLogPath docstring for the contract.
		home, _ := os.UserHomeDir()
		cfg := resolveLogPath(os.Getenv("HYBROID_LS_LOG"), home)
		configureLog(cfg)
	}

	log.Println("Starting HybroidLS")
	log.Println("Warning: HybroidLS is experimental. Expect bugs or missing features!")

	walker.SetupLibraryEnvironments()

	log.Println("Preparing to communicate via stdio")

	var connOpt []jsonrpc2.ConnOpt
	if core.IsDebug {
		connOpt = append(connOpt, jsonrpc2.LogMessages(log.Default()))
	}

	handler := NewHandler()
	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}),
		handler, connOpt...).DisconnectNotify()

	log.Println("All Connections Closed")
}
