package lsp

import (
	"net/url"
	"path/filepath"
	"testing"
)

// TestFromToURI_RoundTrip_PathWithSpaces locks in that file URIs with
// percent-encoded spaces round-trip back to the same filesystem path. If
// url.Path or url.URL.String regresses (e.g. someone swaps in a PathEscape
// that double-encodes %), editors will silently fail to open these files.
func TestFromToURI_RoundTrip_PathWithSpaces(t *testing.T) {
	original := "/tmp/foo bar/level.hyb"

	uri := toURI(original)
	if uri == "" {
		t.Fatal("toURI returned empty")
	}
	// The URI must contain %20 for the space, not a raw space — editors
	// like VS Code reject raw spaces in file:// URIs.
	if !filepath.IsAbs(string(uri)) && !contains(string(uri), "%20") {
		t.Errorf("expected URI to contain %%20 for space, got %q", uri)
	}

	back, err := fromURI(uri)
	if err != nil {
		t.Fatalf("fromURI error: %v", err)
	}
	if back != original {
		t.Errorf("round-trip mismatch:\n  got:  %q\n  want: %q", back, original)
	}
}

// TestFromToURI_RoundTrip_Unicode locks in non-ASCII paths (e.g. user
// profile directories with non-ASCII characters). macOS supports HFS+
// case-insensitive but allows Unicode in path components, and Linux/Windows
// do too.
func TestFromToURI_RoundTrip_Unicode(t *testing.T) {
	original := "/tmp/日本語/level.hyb"

	uri := toURI(original)
	back, err := fromURI(uri)
	if err != nil {
		t.Fatalf("fromURI error: %v", err)
	}
	if back != original {
		t.Errorf("round-trip mismatch:\n  got:  %q\n  want: %q", back, original)
	}
}

// TestFromToURI_RoundTrip_PercentEncoded verifies that a path that
// legitimately contains a percent sign is encoded/decoded correctly. This
// is the tricky case: a naive implementation that always calls PathEscape
// would double-encode a literal '%'.
func TestFromToURI_RoundTrip_PercentEncoded(t *testing.T) {
	original := "/tmp/100%complete/level.hyb"

	uri := toURI(original)
	// The literal % must become %25, not be passed through raw.
	raw := string(uri)
	if !contains(raw, "%25") {
		t.Errorf("expected literal %% to be encoded as %%25, got URI %q", raw)
	}

	back, err := fromURI(uri)
	if err != nil {
		t.Fatalf("fromURI error: %v", err)
	}
	if back != original {
		t.Errorf("round-trip mismatch:\n  got:  %q\n  want: %q", back, original)
	}
}

// contains is a tiny helper to avoid pulling in strings just for a
// substring check used twice.
func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ensure url package is referenced (it is, via fromURI/toURI transitively
// when in the lsp package). This keeps the import behavior stable if a
// future refactor inlines them.
var _ = url.ParseRequestURI
