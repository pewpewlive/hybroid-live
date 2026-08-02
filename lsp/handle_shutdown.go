package lsp

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleShutdown(_ context.Context, _ notifier, _ *jsonrpc2.Request) (result any, err error) {
	// The shutdown request must receive a response before the client sends the
	// exit notification. Closing the connection here prevents that response
	// from reaching the client and makes a graceful shutdown look like a crash.
	return nil, nil
}
