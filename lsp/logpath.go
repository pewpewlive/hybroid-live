package lsp

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// logConfig is the result of resolving where the LSP debug log
// should go. A zero-value path means "discard logs" (e.g. when
// the home directory is unavailable).
type logConfig struct {
	// path is the absolute path to the log file. Empty means
	// "use io.Discard" (the log output is silenced).
	path string
	// source records which resolution branch produced this
	// config: "env" (HYBROID_LS_LOG override), "home"
	// (~/.hybroid/logs/lsp.log), or "discard" (no path). It's
	// surfaced in the startup log line so users can see why
	// logs went where they did.
	source string
}

// resolveLogPath picks the destination for the LSP debug log
// without touching the filesystem. The order of precedence is:
//
//  1. envOverride (HYBROID_LS_LOG) — explicit override always wins.
//  2. homeDir + "/.hybroid/logs/lsp.log" — the documented install
//     location, available on every platform via os.UserHomeDir().
//  3. "" + "discard" — when homeDir is empty (rare, but
//     os.UserHomeDir() can fail on misconfigured CI runners).
//
// The function does NOT create directories and does NOT check
// writability — that's configureLog's job. Splitting the two
// phases makes resolveLogPath trivially testable: no filesystem,
// no I/O, no goroutines.
func resolveLogPath(envOverride, homeDir string) logConfig {
	if envOverride != "" {
		return logConfig{path: envOverride, source: "env"}
	}
	if homeDir == "" {
		return logConfig{path: "", source: "discard"}
	}
	return logConfig{
		path:   filepath.Join(homeDir, ".hybroid", "logs", "lsp.log"),
		source: "home",
	}
}

// configureLog opens the log file (creating the parent directory
// if needed) and points the standard logger at it. On any
// failure — missing home dir, unwritable path, permission
// denied — it falls back to io.Discard so the JSON-RPC server
// keeps running. The fallback is logged to stderr once at
// startup so the user can see why their logs aren't going where
// they expected.
//
// When the path came from the env override, configureLog does
// NOT create the parent directory: the caller who set the env
// var is responsible for ensuring the path is usable. This
// preserves the "explicit override" contract: an override that
// points at a nonexistent path should fail loudly, not be
// silently recreated somewhere else.
func configureLog(cfg logConfig) {
	if cfg.path == "" {
		log.SetOutput(io.Discard)
		return
	}

	dir := filepath.Dir(cfg.path)
	// For env-override paths, do NOT mkdir — the override is
	// supposed to be an existing path. For home-resolved paths,
	// mkdir the ~/.hybroid/logs/ dir if missing.
	if cfg.source != "env" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			// Couldn't create the directory. Don't crash —
			// the JSON-RPC stream is more important than
			// the debug log. Fall back to discard and tell
			// the user on stderr.
			log.SetOutput(io.Discard)
			os.Stderr.WriteString("hybroid-ls: could not create log directory " + dir + ": " + err.Error() + "; debug logging disabled\n")
			return
		}
	}

	f, err := os.OpenFile(cfg.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.SetOutput(io.Discard)
		os.Stderr.WriteString("hybroid-ls: could not open log file " + cfg.path + ": " + err.Error() + "; debug logging disabled\n")
		return
	}
	log.SetOutput(f)
	log.Println("Debug mode enabled, logging to", cfg.path, "(source:", cfg.source+")")
	// Note: we intentionally do not defer f.Close() — the file
	// is closed by the OS on process exit. Closing earlier would
	// prevent any post-disconnect logging from being flushed.
}
