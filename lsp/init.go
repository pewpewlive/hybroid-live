package lsp

import (
	"context"
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

func Init() {
	f, err := os.OpenFile("D:\\testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	log.Println("Starting Internal Hybroid Language Server")
	log.Println("WARNING: THIS SERVER IS IN PRE-ALPHA STATE!!! USE WITH CAUTION!")

	log.Println("Preparing to communicate via stdio")

	var connOpt []jsonrpc2.ConnOpt
	connOpt = append(connOpt, jsonrpc2.LogMessages(log.Default()))

	handler := NewHandler()
	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}),
		handler, connOpt...).DisconnectNotify()

	log.Println("All Connections Closed")
	f.Close()
}
