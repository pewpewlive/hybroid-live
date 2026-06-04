package lsp

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureStderr redirects os.Stderr to a pipe for the duration of
// fn, returning whatever was written. Used by configureLog tests
// that need to assert on the fallback warning.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	_ = w.Close()
	return <-done
}

// withLogOutput swaps the default logger's output for the
// duration of fn, restoring the previous output on return. Every
// configureLog test uses this to keep its changes from leaking
// into other tests' log assertions.
func withLogOutput(t *testing.T, fn func()) {
	t.Helper()
	orig := log.Default().Writer()
	log.SetOutput(io.Discard) // baseline: also silence
	defer log.SetOutput(orig)
	fn()
}

// TestResolveLogPath_EnvOverrideWins pins precedence rule #1: an
// explicit HYBROID_LS_LOG value always wins, regardless of
// homeDir. This matches the historical behavior — the env var
// has been a documented escape hatch since v0.1.
func TestResolveLogPath_EnvOverrideWins(t *testing.T) {
	got := resolveLogPath("/tmp/custom.log", "/home/alice")
	if got.path != "/tmp/custom.log" {
		t.Errorf("path: got %q, want %q", got.path, "/tmp/custom.log")
	}
	if got.source != "env" {
		t.Errorf("source: got %q, want %q", got.source, "env")
	}
}

// TestResolveLogPath_HomeDir pins the standard case: env unset,
// home available → ~/.hybroid/logs/lsp.log.
func TestResolveLogPath_HomeDir(t *testing.T) {
	got := resolveLogPath("", "/home/alice")
	want := filepath.Join("/home/alice", ".hybroid", "logs", "lsp.log")
	if got.path != want {
		t.Errorf("path: got %q, want %q", got.path, want)
	}
	if got.source != "home" {
		t.Errorf("source: got %q, want %q", got.source, "home")
	}
}

// TestResolveLogPath_EmptyHomeFallsBack covers the rare case
// where os.UserHomeDir() returns "" (misconfigured CI, weird
// chroot, $HOME unset on Linux, $USERPROFILE unset on Windows).
// The function must return a discard config — the caller will
// then point the logger at io.Discard.
func TestResolveLogPath_EmptyHomeFallsBack(t *testing.T) {
	got := resolveLogPath("", "")
	if got.path != "" {
		t.Errorf("path: got %q, want \"\" (empty = discard)", got.path)
	}
	if got.source != "discard" {
		t.Errorf("source: got %q, want %q", got.source, "discard")
	}
}

// TestConfigureLog_CreatesMissingDir asserts that the first call
// to configureLog with a home-resolved path creates the
// ~/.hybroid/logs/ directory if it doesn't exist.
func TestConfigureLog_CreatesMissingDir(t *testing.T) {
	withLogOutput(t, func() {
		home := t.TempDir()
		// No .hybroid/logs/ exists yet. resolveLogPath picks the
		// home branch, configureLog must MkdirAll the parents.
		cfg := resolveLogPath("", home)
		if cfg.source != "home" {
			t.Fatalf("setup: expected home source, got %q", cfg.source)
		}

		configureLog(cfg)

		logDir := filepath.Join(home, ".hybroid", "logs")
		stat, err := os.Stat(logDir)
		if err != nil {
			t.Fatalf("expected log dir to exist after configureLog: %v", err)
		}
		if !stat.IsDir() {
			t.Errorf("expected dir, got file at %q", logDir)
		}
	})
}

// TestConfigureLog_OpensExistingFile asserts that a log.Println
// call lands in the resolved file. This is the happy path: the
// dir exists, the file doesn't, the call creates it and writes
// to it.
func TestConfigureLog_OpensExistingFile(t *testing.T) {
	withLogOutput(t, func() {
		home := t.TempDir()
		cfg := resolveLogPath("", home)

		configureLog(cfg)

		log.Println("hello from the test")
		// log.Println appends a newline, but the standard
		// logger also writes to its output (which is our
		// file). We need to flush. The standard logger has
		// no Flush, but the file is opened with O_APPEND so
		// the OS flushes on write. Read back and assert.
		data, err := os.ReadFile(cfg.path)
		if err != nil {
			t.Fatalf("read log file: %v", err)
		}
		if !strings.Contains(string(data), "hello from the test") {
			t.Errorf("log file %q did not contain expected message; contents: %q",
				cfg.path, string(data))
		}
	})
}

// TestConfigureLog_FallsBackOnUnwritableDir asserts the
// graceful fallback: if MkdirAll fails (parent is a file, not
// a dir), configureLog sets log output to io.Discard and
// prints a one-line warning to stderr.
func TestConfigureLog_FallsBackOnUnwritableDir(t *testing.T) {
	withLogOutput(t, func() {
		// Build a path where the immediate parent is a file.
		// MkdirAll(home/.hybroid/logs) will fail because
		// home/.hybroid is a file, not a directory.
		home := t.TempDir()
		collision := filepath.Join(home, ".hybroid")
		if err := os.WriteFile(collision, []byte("not a dir"), 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}
		cfg := resolveLogPath("", home)

		warning := captureStderr(t, func() {
			configureLog(cfg)
		})

		if !strings.Contains(warning, "could not create log directory") {
			t.Errorf("expected fallback warning on stderr, got %q", warning)
		}
		if !strings.Contains(warning, "debug logging disabled") {
			t.Errorf("expected warning to mention 'debug logging disabled', got %q", warning)
		}
	})
}

// TestConfigureLog_EnvOverrideDoesNotMkdir asserts the explicit
// override contract: HYBROID_LS_LOG points at a path whose
// parent directory doesn't exist, and configureLog must NOT
// create the parent. The override is supposed to be a path the
// user picked deliberately; if it points somewhere nonexistent
// we want the failure to be visible, not silently fixed up.
func TestConfigureLog_EnvOverrideDoesNotMkdir(t *testing.T) {
	withLogOutput(t, func() {
		home := t.TempDir()
		// /home/<tmp>/logs/missing.log — /home/<tmp>/logs/
		// does not exist. We expect configureLog to attempt
		// the OpenFile and fail (no parent dir), then fall
		// back to discard. Critically, the parent dir must
		// NOT have been created.
		override := filepath.Join(home, "logs", "missing.log")
		cfg := resolveLogPath(override, home)
		if cfg.source != "env" {
			t.Fatalf("setup: expected env source, got %q", cfg.source)
		}

		warning := captureStderr(t, func() {
			configureLog(cfg)
		})

		logsDir := filepath.Join(home, "logs")
		if _, err := os.Stat(logsDir); err == nil {
			t.Errorf("env override path: configureLog must not create %q", logsDir)
		}
		if !strings.Contains(warning, "could not open log file") {
			t.Errorf("expected fallback warning on stderr, got %q", warning)
		}
	})
}
