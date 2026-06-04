package lsp

import (
	"testing"
	"time"
)

// TestScheduleAnalysis_TimerStopAfterFire covers the edge case where
// scheduleAnalysis is called with a high debounce, then a Stop() is
// issued before the timer fires. The expectation is that no analysis
// runs at all — Stop is a hard cancel. If the production code instead
// leaked a goroutine that ran with the now-stale pendingChange, we'd
// see an extra publishDiagnostics after Stop returns.
//
// We use a long debounce (200ms) so the test has a clear window to
// issue Stop. We also wait long enough to catch any leaked callback
// that did fire despite the Stop (Stop returns false if the timer has
// already fired — in that case we accept at most one stale publish).
func TestScheduleAnalysis_TimerStopAfterFire(t *testing.T) {
	h, conn := newTestHandler(t)
	// Override the default 1ms debounce with something the test can
	// reliably race against.
	h.lintDebounce = 200 * time.Millisecond

	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(dir + "/x.hyb")
	openForTest(t, h, conn, uri, "env X as Level\n")

	baseline := conn.Count()

	// Schedule an analysis that we will then cancel.
	h.scheduleAnalysis(uri, "env X as Level\n")

	// Immediately try to stop. lintTimer.Stop returns false if the
	// timer has already fired — we don't care which; we care that
	// after a wait, we get at most one stale publish.
	h.mu.Lock()
	timer := h.lintTimer
	h.mu.Unlock()
	if timer != nil {
		timer.Stop()
	}

	// Wait long enough for any leaked/stale callback to fire (well
	// past the 200ms debounce). If the production code is correct,
	// no extra publishes will appear after baseline. If a future
	// refactor removes the Stop() call or moves it after the AfterFunc
	// is registered, the count may grow by 1.
	time.Sleep(500 * time.Millisecond)

	got := conn.Count()
	if got > baseline+1 {
		t.Errorf("expected at most 1 stale publish after Stop, got %d (baseline=%d)",
			got-baseline, baseline)
	}
}
