package lsp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

// fakeNotify records every Notify call so tests can assert on the wire
// traffic the handler emits. It also serves as a synchronization point —
// tests that need to wait for a particular publish can poll Notifies until
// the expected count is reached (with a timeout).
type fakeNotify struct {
	mu       sync.Mutex
	notifies []capturedNotify
	closed   bool
}

type capturedNotify struct {
	Method string
	Params any
}

func newFakeConn() *fakeNotify {
	return &fakeNotify{notifies: make([]capturedNotify, 0, 16)}
}

func (f *fakeNotify) Notify(_ context.Context, method string, params any, _ ...jsonrpc2.CallOption) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.notifies = append(f.notifies, capturedNotify{Method: method, Params: params})
	return nil
}

func (f *fakeNotify) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
	return nil
}

func (f *fakeNotify) Notifies() []capturedNotify {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]capturedNotify, len(f.notifies))
	copy(out, f.notifies)
	return out
}

func (f *fakeNotify) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.notifies)
}

func (f *fakeNotify) CountByMethod(method string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.notifies {
		if c.Method == method {
			n++
		}
	}
	return n
}

// WaitFor polls until at least n notifies have been recorded or the deadline
// elapses. Returns the final count. Use this in tests that need to wait for
// an async publish (debounced analysis, goroutine-launched pre-analysis).
func (f *fakeNotify) WaitFor(n int, timeout time.Duration) int {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		c := f.Count()
		if c >= n {
			return c
		}
		time.Sleep(5 * time.Millisecond)
	}
	return f.Count()
}

// newTestHandler builds a langHandler with sensible test defaults:
//   - fakeNotify conn so tests can observe publishDiagnostics
//   - 1ms debounce so tests don't wait the production 300ms
//   - default rootMarkers ["hybconfig.toml"]
//   - markReady already called (no pre-analysis to wait for)
func newTestHandler(t *testing.T) (*langHandler, *fakeNotify) {
	t.Helper()
	conn := newFakeConn()
	h := &langHandler{
		provideDefinition:      true,
		files:                  make(map[DocumentURI]*File),
		request:                make(chan lintRequest),
		rootMarkers:            []string{"hybconfig.toml"},
		lintDebounce:           1 * time.Millisecond,
		ready:                  make(chan struct{}),
		lastPublishedURIs:      make(map[string]map[DocumentURI]struct{}),
		infoNoticesPublished:   make(map[DocumentURI]struct{}),
		conn:                   conn,
	}
	h.markReady()
	return h, conn
}

// newTestHandlerWithRoot is a convenience wrapper that pre-sets rootPath
// (mimicking handleInitialize) so tests can exercise the "workspace is open"
// branch.
func newTestHandlerWithRoot(t *testing.T, rootPath string) (*langHandler, *fakeNotify) {
	t.Helper()
	h, conn := newTestHandler(t)
	h.rootPath = rootPath
	h.addFolder(rootPath)
	return h, conn
}

// encodeParams marshals v to a *json.RawMessage the way jsonrpc2 expects.
func encodeParams(t *testing.T, v any) *json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("encodeParams: %v", err)
	}
	raw := json.RawMessage(b)
	return &raw
}

// newTestRequest builds a jsonrpc2.Request from a method and params struct.
func newTestRequest(method string, params any) *jsonrpc2.Request {
	var raw *json.RawMessage
	if params != nil {
		b, _ := json.Marshal(params)
		r := json.RawMessage(b)
		raw = &r
	}
	return &jsonrpc2.Request{
		Method: method,
		Params: raw,
	}
}

// writeProject writes files (relative path -> contents) into dir, creating
// parent directories as needed. It returns the absolute path of dir.
func writeProject(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("writeProject mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("writeProject write: %v", err)
		}
	}
	return dir
}

// minimalHybConfig is the smallest valid hybconfig.toml the LSP needs to
// recognize a directory as a Hybroid project. content is mostly irrelevant
// to the LSP today, but having the file present is what findProjectRoot
// checks for.
const minimalHybConfig = `[project]
name = "test"
output_directory = "out"

[level]
entry_point = "level.hyb"
`

// minimalLevelSource is a valid (if trivial) level source so the evaluator
// doesn't fail parsing.
const minimalLevelSource = `env TestLevel as Level

tick {
}
`
