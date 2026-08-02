package lsp

import (
	"context"
	"testing"
)

func TestHandleShutdownLeavesConnectionOpenForExitNotification(t *testing.T) {
	conn := newFakeConn()
	handler := &langHandler{}

	result, err := handler.handleShutdown(context.Background(), conn, nil)
	if err != nil {
		t.Fatalf("handleShutdown returned an error: %v", err)
	}
	if result != nil {
		t.Fatalf("handleShutdown returned %v, want nil", result)
	}

	conn.mu.Lock()
	closed := conn.closed
	conn.mu.Unlock()
	if closed {
		t.Fatal("handleShutdown closed the connection before the exit notification")
	}
}
