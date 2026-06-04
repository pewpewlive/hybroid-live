package lsp

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleShutdown(_ context.Context, conn notifier, _ *jsonrpc2.Request) (result any, err error) {
	return nil, nil
}
