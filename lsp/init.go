package lsp

import (
	"context"
	"hybroid/core"
	"hybroid/walker"
	"io"
	"log"
	"os"
	"path/filepath"

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
		// Pick a log destination that is always writable. The cwd of the
		// spawned server process is not under our control (e.g. an
		// extension host may launch us with a read-only cwd), so prefer
		// an explicit path from HYBROID_LS_LOG, then the user's home
		// directory, and only fall back to a relative path if both fail.
		logPath := os.Getenv("HYBROID_LS_LOG")
		if logPath == "" {
			if home, err := os.UserHomeDir(); err == nil {
				logPath = filepath.Join(home, "hybroid_ls.log")
			} else {
				logPath = "hybroid_ls.log"
			}
		}
		if f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			log.SetOutput(f)
			log.Println("Debug mode enabled, logging to", logPath)
			// Note: we intentionally do not defer f.Close() — the file is closed
			// by the OS on process exit. Closing earlier would prevent any
			// post-disconnect logging from being flushed.
		} else {
			// Could not open the log file (e.g. read-only home dir).
			// Fall back to discarding log output rather than crashing the
			// server with log.Fatalf — the JSON-RPC stream is more important.
			log.SetOutput(io.Discard)
		}
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
