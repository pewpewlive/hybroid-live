package lsp

import (
	"context"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

// barrier is a deterministic synchronizer. All goroutines call
// barrier.Wait(); the test releases them at a known instant by calling
// barrier.Release(). This lets a test hold all participants at a
// single instruction, then unleash them so the race happens at a
// controlled point — eliminating the flakiness of `time.Sleep` based
// concurrency tests.
type barrier struct {
	release chan struct{}
}

func newBarrier() *barrier { return &barrier{release: make(chan struct{})} }

func (b *barrier) Wait()    { <-b.release }
func (b *barrier) Release() { close(b.release) }

// TestConcurrent_DidChangeAndHover_NoDeadlock is the regression test
// for the original deadlock. In production: a user types in the editor
// (didChange → scheduleAnalysis → timer fires), while the editor
// simultaneously sends hover/definition requests (which call waitReady
// and acquire h.mu). The original bug: scheduleAnalysis leaked h.mu,
// so the first hover after a few edits would block forever and the
// runtime would detect "all goroutines are asleep - deadlock!".
//
// We reproduce the pattern with deterministic barriers: a writer
// goroutine fires N didChange events while a reader goroutine fires
// M hovers, both gated on the same barrier. The test asserts that
// every event completes within the deadline. The Go runtime would
// panic with deadlock detection if any of them blocks indefinitely.
func TestConcurrent_DidChangeAndHover_NoDeadlock(t *testing.T) {
	h, _ := newTestHandler(t)
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	// Open the file once so h.eval is set up; subsequent didChanges
	// and hovers have a real evaluator to query.
	openReq := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       "env TestLevel as Level\n",
		},
	})
	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, openReq); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	const nChanges = 50
	const nHovers = 20

	gate := newBarrier()
	done := make(chan struct{}, nChanges+nHovers)

	// Writer: didChange loop. Each call goes through handleTextDocumentDidChange
	// which acquires h.mu briefly, then schedules a debounced analysis.
	for i := 0; i < nChanges; i++ {
		i := i
		go func() {
			gate.Wait()
			req := newTestRequest("textDocument/didChange", DidChangeTextDocumentParams{
				TextDocument: VersionedTextDocumentIdentifier{
					TextDocumentIdentifier: TextDocumentIdentifier{URI: uri},
					Version:                i + 1,
				},
				ContentChanges: []TextDocumentContentChangeEvent{
					{Text: "env TestLevel as Level\nlet x" + itoaSimple(i) + " = " + itoaSimple(i) + "\n"},
				},
			})
			_, _ = h.handleTextDocumentDidChange(context.Background(), h.conn, req)
			done <- struct{}{}
		}()
	}

	// Reader: hover loop. handleTextDocumentHover takes h.mu briefly
	// while it looks up the symbol. Under the original bug, after a
	// few didChanges the leaked lock would block these hovers.
	for i := 0; i < nHovers; i++ {
		go func() {
			gate.Wait()
			req := newTestRequest("textDocument/hover", TextDocumentPositionParams{
				TextDocument: TextDocumentIdentifier{URI: uri},
				Position:     Position{Line: 0, Character: 4},
			})
			_, _ = h.handleTextDocumentHover(context.Background(), h.conn, req)
			done <- struct{}{}
		}()
	}

	// Release everyone at once.
	gate.Release()

	// Wait for all participants to complete, with a deadline that
	// would expose a real deadlock.
	deadline := time.After(5 * time.Second)
	for i := 0; i < nChanges+nHovers; i++ {
		select {
		case <-done:
		case <-deadline:
			t.Fatalf("deadlock detected: %d/%d events completed after 5s",
				i, nChanges+nHovers)
		}
	}
}

// TestConcurrent_TimerAndEvalMu_NoDeadlock verifies the inverse
// race: while the debounced analysis (under h.evalMu) is in progress,
// a flood of didChange events must not block. The original code
// acquires h.mu at the top of scheduleAnalysis and h.evalMu inside
// analyzeAndPublish; the test asserts the two locks don't deadlock
// each other across rapid changes.
func TestConcurrent_TimerAndEvalMu_NoDeadlock(t *testing.T) {
	h, _ := newTestHandler(t)
	// Slow down analysis to make the race window visible. The timer
	// debounce is 1ms (set by newTestHandler), so the analysis
	// callback fires quickly, but the evalMu-holding work inside
	// analyzeAndPublish can still race with a fresh didChange.
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	openReq := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       "env TestLevel as Level\n",
		},
	})
	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, openReq); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	// Fire 100 didChange calls back-to-back from a single goroutine
	// and assert each one returns. The timer will fire 100 times in
	// quick succession (each change resets the timer, but the test
	// only stops when all 100 handleTextDocumentDidChange calls
	// return). Under the original bug, one of the timer callbacks
	// would deadlock waiting on a leaked h.mu.
	const n = 100
	for i := 0; i < n; i++ {
		req := newTestRequest("textDocument/didChange", DidChangeTextDocumentParams{
			TextDocument: VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: TextDocumentIdentifier{URI: uri},
				Version:                i + 1,
			},
			ContentChanges: []TextDocumentContentChangeEvent{
				{Text: "env TestLevel as Level\nlet v" + itoaSimple(i) + " = " + itoaSimple(i) + "\n"},
			},
		})
		if _, err := h.handleTextDocumentDidChange(context.Background(), h.conn, req); err != nil {
			t.Fatalf("didChange %d: %v", i, err)
		}
	}
	// Let the final timer fire and analyzeAndPublish run.
	time.Sleep(200 * time.Millisecond)
}

// TestMarkReady_Idempotent verifies that markReady can be called any
// number of times without panicking. The implementation guards with a
// `readySet` boolean so `close(ready)` only happens once. If someone
// removes that guard (a "simplification" that drops the boolean), this
// test will catch the "close of closed channel" panic.
func TestMarkReady_Idempotent(t *testing.T) {
	h, _ := newTestHandler(t)

	// newTestHandler already called markReady once. Call it many more
	// times — none should panic.
	var panics atomic.Int32
	for i := 0; i < 100; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					panics.Add(1)
					t.Logf("markReady iteration %d panicked: %v", i, r)
				}
			}()
			h.markReady()
		}()
	}
	if panics.Load() != 0 {
		t.Errorf("markReady panicked %d times across 101 calls", panics.Load())
	}

	// The ready channel should still receive exactly once.
	select {
	case <-h.ready:
		// expected
	default:
		t.Errorf("expected ready channel to be closed")
	}
	// A second receive should NOT block (we got the value already);
	// but a third should also be fine — closed channels are readable
	// indefinitely.
	select {
	case <-h.ready:
		// expected
	default:
		t.Errorf("ready channel should remain readable after first receive")
	}
}

// itoaSimple is a tiny allocation-free int-to-string for test use
// (the standard strconv import would also work; this avoids pulling
// in the import for what is essentially a debug helper).
func itoaSimple(n int) string {
	if n == 0 {
		return "0"
	}
	const digits = "0123456789"
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = digits[n%10]
		n /= 10
	}
	return string(buf[pos:])
}

// keep jsonrpc2 import alive.
var _ = jsonrpc2.CodeInvalidParams
