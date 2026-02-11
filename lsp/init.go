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
		f, err := os.OpenFile("D:\\testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
		// We can't defer f.Close() here because Init returns while server is running.
		// However, for a CLI tool it's usually fine as OS will close it.
	}

	log.Println("Starting Integrated Language Server for Hybroid")
	log.Println("WARNING: THIS SERVER IS IN PRE-ALPHA STATE!!! USE WITH CAUTION!")

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
