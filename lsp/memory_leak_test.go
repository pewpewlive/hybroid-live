package lsp

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// closeForTest wraps the boilerplate of issuing a didClose for a single
// file. It does not check the response — this helper exists purely so
// the memory-leak tests can iterate open/close cycles without
// duplicating the request-construction code from cross_file_diagnostics_test.go.
func closeForTest(t *testing.T, h *langHandler, uri DocumentURI) {
	t.Helper()
	req := newTestRequest("textDocument/didClose", DidCloseTextDocumentParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
	})
	if _, err := h.handleTextDocumentDidClose(context.Background(), h.conn, req); err != nil {
		t.Fatalf("didClose: %v", err)
	}
}

// TestFilesMapCleanedUpOnClose pins the basic invariant: after a
// didOpen followed by a didClose, the handler's internal files map is
// empty. This is the foundation everything else (heap, goroutine
// counts) builds on — if entries aren't removed, the heap will grow
// linearly with the number of opens over a server's lifetime.
func TestFilesMapCleanedUpOnClose(t *testing.T) {
	h, conn := newTestHandler(t)
	uri := DocumentURI("file:///leak-test.hyb")

	openForTest(t, h, conn, uri, minimalLevelSource)
	if got := len(h.files); got != 1 {
		t.Fatalf("after open: got %d files, want 1", got)
	}

	closeForTest(t, h, uri)
	if got := len(h.files); got != 0 {
		t.Errorf("after close: got %d files, want 0", got)
	}
}

// TestNoGoroutineLeakOnRapidOpenClose opens and closes 500 files in
// quick succession, then asserts that the goroutine count is back at
// the baseline (within a small tolerance for the test runtime itself
// — Go's test runner spawns per-test goroutines, and runtime.GC
// itself briefly elevates the count).
//
// What this catches: a didChange handler that spawns a goroutine per
// change and never lets it exit; a debounce timer that's reset but
// not stopped, leaving an AfterFunc goroutine hanging; a closure
// that captures the handler reference and blocks on a channel that
// never closes.
//
// What this does NOT catch: leaks that grow slowly (the test uses a
// fixed iteration count and a generous tolerance). Long-running
// memory growth is covered by TestNoHeapGrowthOnRepeatedCycles below.
func TestNoGoroutineLeakOnRapidOpenClose(t *testing.T) {
	h, conn := newTestHandler(t)

	// Settle: do one full open+close cycle to make sure all
	// initialization goroutines (if any) have completed before we
	// record the baseline.
	settleURI := DocumentURI("file:///leak-settle.hyb")
	openForTest(t, h, conn, settleURI, minimalLevelSource)
	closeForTest(t, h, settleURI)
	runtime.GC()
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	const iterations = 500
	for i := 0; i < iterations; i++ {
		uri := DocumentURI(fmt.Sprintf("file:///leak/%d.hyb", i))
		openForTest(t, h, conn, uri, minimalLevelSource)
		closeForTest(t, h, uri)
	}

	// Let any pending debounce timers (1ms in test mode) fire and
	// their goroutines exit.
	time.Sleep(50 * time.Millisecond)
	runtime.GC()
	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	final := runtime.NumGoroutine()
	const tolerance = 5
	if final > baseline+tolerance {
		t.Errorf("goroutine count grew: baseline=%d, final=%d (delta=%d, tolerance=%d)",
			baseline, final, final-baseline, tolerance)
	}
}

// TestEvaluatorWalkerListStableAcrossOpenClose pins the contract that
// the LSP's per-file cleanup is reflected in the evaluator's
// walkerList. After N opens and N closes of distinct URIs, the
// walkerList length should match what was discovered via the
// first didOpen (which is what the evaluator was created with),
// not the cumulative count of opens.
//
// This catches a known issue: didClose currently only calls
// `delete(h.files, uri)` but does not tell the evaluator to drop
// the file from its walkers/programs/fileContents maps. Each new
// single-file open in the same session grows the evaluator's
// internal state unboundedly. A future fix should add
// `Evaluator.RemoveFile(path)` and call it from didClose.
//
// As of this writing, this test is expected to FAIL — it
// documents the leak. The first iteration's walkerList length is
// the reference point; subsequent iterations must not grow it.
func TestEvaluatorWalkerListStableAcrossOpenClose(t *testing.T) {
	h, conn := newTestHandler(t)

	// Establish reference: open one file, capture walkerList length.
	uri0 := DocumentURI("file:///eval-ref.hyb")
	openForTest(t, h, conn, uri0, minimalLevelSource)
	h.evalMu.Lock()
	baseline := len(h.eval.WalkerList())
	h.evalMu.Unlock()
	if baseline == 0 {
		t.Fatal("evaluator has no walkers after first open — test setup wrong")
	}

	// Open and close N more files. After this, the evaluator should
	// still have `baseline` walkers (the file0 walker), not baseline+N.
	const extra = 20
	for i := 0; i < extra; i++ {
		uri := DocumentURI(fmt.Sprintf("file:///eval-extra-%d.hyb", i))
		openForTest(t, h, conn, uri, minimalLevelSource)
	}

	// At this point the evaluator has baseline+1+extra walkers (we
	// haven't closed anything). Now close them all and check the
	// count returns to baseline.
	for i := 0; i < extra; i++ {
		uri := DocumentURI(fmt.Sprintf("file:///eval-extra-%d.hyb", i))
		closeForTest(t, h, uri)
	}
	closeForTest(t, h, uri0)

	h.evalMu.Lock()
	final := len(h.eval.WalkerList())
	h.evalMu.Unlock()

	if final > baseline {
		t.Errorf("evaluator walkerList grew from %d to %d after %d open+close cycles; didClose should drop the file from the evaluator",
			baseline, final, extra)
	}
}

// TestInfoNoticesMapBoundedByUniqueURIs is a documentation test: it
// records the current behavior that publishInfoOnce keeps a
// one-entry-per-URI in h.infoNoticesPublished, and asserts the
// invariant "no duplicate entries for the same URI". This isn't a
// leak per se (the map has at most one entry per unique URI), but a
// future refactor that accidentally added a slice or list would
// blow up the memory cost of long-lived servers opening many
// distinct single-file buffers.
func TestInfoNoticesMapBoundedByUniqueURIs(t *testing.T) {
	h, conn := newTestHandler(t)

	// Open the same URI multiple times. infoNoticesPublished should
	// have at most one entry for it.
	uri := DocumentURI("file:///info-bound.hyb")
	for i := 0; i < 5; i++ {
		openForTest(t, h, conn, uri, minimalLevelSource)
	}

	h.mu.Lock()
	count := len(h.infoNoticesPublished)
	h.mu.Unlock()
	if count != 1 {
		t.Errorf("infoNoticesPublished has %d entries for one URI, want 1", count)
	}
}
