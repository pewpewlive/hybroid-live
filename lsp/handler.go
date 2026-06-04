package lsp

import (
	"context"
	"fmt"
	"hybroid/core"
	"hybroid/evaluator"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode"

	"github.com/sourcegraph/jsonrpc2"
)

type eventType int

type lintRequest struct {
	URI       DocumentURI
	EventType eventType
}

// notifier is the minimal surface that langHandler uses to push messages to
// the LSP client. In production this is *jsonrpc2.Conn; in tests it is a fake
// that records calls. The variadic CallOption must be preserved verbatim —
// Go's method-set rules mean a fixed-arg interface cannot be satisfied by a
// variadic concrete method.
type notifier interface {
	Notify(ctx context.Context, method string, params any, opts ...jsonrpc2.CallOption) error
	Close() error
}

type File struct {
	LanguageID string
	Text       string
	Version    int
}

type langHandler struct {
	mu     sync.Mutex
	evalMu sync.Mutex
	logger *log.Logger
	// commands          []Command
	provideDefinition bool
	files             map[DocumentURI]*File
	eval              *evaluator.Evaluator
	lintDebounce      time.Duration
	request           chan lintRequest
	lintTimer         *time.Timer
	formatDebounce    time.Duration
	formatTimer       *time.Timer
	conn              notifier
	rootPath          string
	rootURI           DocumentURI
	filename          string
	folders           []string
	rootMarkers       []string
	triggerChars      []string

	// ready is closed once the initial pre-analysis has finished. Handlers
	// that depend on h.eval should wait on it to avoid racing with initialization.
	ready    chan struct{}
	readySet bool

	// pendingChange is the URI/text of the most recent didChange that has
	// not yet been analyzed because the lint timer hasn't fired.
	pendingChange struct {
		uri  DocumentURI
		text string
	}

	// lastPublishedURIs is mapping from LanguageID string to mapping of
	// whether diagnostics are published in a DocumentURI or not.
	lastPublishedURIs map[string]map[DocumentURI]struct{}

	// infoNoticesPublished tracks URIs that have already received a one-shot
	// "workspace context missing" Information diagnostic, so we don't republish
	// it on every didChange or didOpen of the same buffer.
	infoNoticesPublished map[DocumentURI]struct{}
}

func (h *langHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	core.DebugLog("Incoming request: %s (notification: %v)", req.Method, req.Notif)
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, conn, req)
	case "initialized":
		return
	case "$/setTrace":
		return
	case "$/cancelRequest":
		return
	case "shutdown":
		return h.handleShutdown(ctx, conn, req)
	case "exit":
		if h.conn != nil {
			_ = h.conn.Close()
		}
		os.Exit(0)
		return nil, nil
	case "textDocument/didOpen":
		return h.handleTextDocumentDidOpen(ctx, conn, req)
	case "textDocument/didChange":
		return h.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/didSave":
		return // h.handleTextDocumentDidSave(ctx, conn, req)
	case "textDocument/didClose":
		return h.handleTextDocumentDidClose(ctx, conn, req)
	case "textDocument/formatting":
		return // h.handleTextDocumentFormatting(ctx, conn, req)
	case "textDocument/rangeFormatting":
		return // h.handleTextDocumentRangeFormatting(ctx, conn, req)
	case "textDocument/documentSymbol":
		return // h.handleTextDocumentSymbol(ctx, conn, req)
	case "textDocument/completion":
		return h.handleTextDocumentCompletion(ctx, conn, req)
	case "textDocument/signatureHelp":
		return h.handleTextDocumentSignatureHelp(ctx, conn, req)
	case "completionItem/resolve":
		return h.HandleCompletionItemResolve(ctx, conn, req)
	case "textDocument/definition":
		return h.handleTextDocumentDefinition(ctx, conn, req)
	case "textDocument/references":
		return h.handleTextDocumentReferences(ctx, conn, req)
	case "textDocument/hover":
		return h.handleTextDocumentHover(ctx, conn, req)
	case "textDocument/rename":
		return h.handleTextDocumentRename(ctx, conn, req)
	case "textDocument/codeAction":
		return // h.handleTextDocumentCodeAction(ctx, conn, req)
	case "workspace/executeCommand":
		return // h.handleWorkspaceExecuteCommand(ctx, conn, req)
	case "workspace/didChangeConfiguration":
		return // h.handleWorkspaceDidChangeConfiguration(ctx, conn, req)
	case "workspace/didChangeWorkspaceFolders":
		return // h.handleDidChangeWorkspaceWorkspaceFolders(ctx, conn, req)
	case "workspace/workspaceFolders":
		return // h.handleWorkspaceWorkspaceFolders(ctx, conn, req)
	}

	return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
}

func NewHandler() jsonrpc2.Handler {
	// logger := log.New(os.Stderr, "", log.LstdFlags)

	handler := &langHandler{
		provideDefinition: true,
		files:             make(map[DocumentURI]*File),
		// evaluator will be initialized in handleInitialize
		request: make(chan lintRequest),
		conn:    nil,
		// filename:          config.Filename,
		rootMarkers:       []string{"hybconfig.toml"},
		// triggerChars:      config.TriggerChars,

		lintDebounce: 300 * time.Millisecond,
		ready:        make(chan struct{}),

		lastPublishedURIs:    make(map[string]map[DocumentURI]struct{}),
		infoNoticesPublished: make(map[DocumentURI]struct{}),
	}
	// handler
	return jsonrpc2.HandlerWithError(handler.handle)
}

func isWindowsDrivePath(path string) bool {
	if len(path) < 4 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

func isWindowsDriveURI(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}

func fromURI(uri DocumentURI) (string, error) {
	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("only file URIs are supported, got %v", u.Scheme)
	}
	if isWindowsDriveURI(u.Path) {
		u.Path = u.Path[1:]
	}
	return u.Path, nil
}

func toURI(path string) DocumentURI {
	if isWindowsDrivePath(path) {
		path = "/" + path
	}
	return DocumentURI((&url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(path),
	}).String())
}

func (h *langHandler) addFolder(folder string) {
	folder = filepath.Clean(folder)
	found := false
	for _, cur := range h.folders {
		if cur == folder {
			found = true
			break
		}
	}
	if !found {
		h.folders = append(h.folders, folder)
	}
}

// markReady closes the ready channel exactly once, signalling that h.eval is
// initialized and safe to use. Safe to call from any goroutine.
func (h *langHandler) markReady() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.readySet {
		close(h.ready)
		h.readySet = true
	}
}

// waitReady blocks until h.eval is ready or ctx is cancelled. Returns true
// if the evaluator became ready, false if the context expired first.
func (h *langHandler) waitReady(ctx context.Context) bool {
	h.mu.Lock()
	ch := h.ready
	alreadyReady := h.readySet
	h.mu.Unlock()
	if alreadyReady {
		return true
	}
	select {
	case <-ch:
		return true
	case <-ctx.Done():
		return false
	}
}

// scheduleAnalysis records the most recent change and (re)starts the lint
// debounce timer. When the timer fires, the pending change is analyzed
// and diagnostics are published. Multiple rapid changes coalesce into a
// single analysis run.
func (h *langHandler) scheduleAnalysis(uri DocumentURI, text string) {
	h.mu.Lock()
	h.pendingChange.uri = uri
	h.pendingChange.text = text
	if h.lintTimer != nil {
		h.lintTimer.Stop()
	}
	conn := h.conn
	debounce := h.lintDebounce
	h.lintTimer = time.AfterFunc(debounce, func() {
		h.mu.Lock()
		uri := h.pendingChange.uri
		text := h.pendingChange.text
		h.pendingChange.uri = ""
		h.pendingChange.text = ""
		h.mu.Unlock()
		if uri == "" {
			return
		}
		h.analyzeAndPublish(context.Background(), conn, uri, text)
	})
	h.mu.Unlock()
}

func (h *langHandler) preAnalyzeWorkspace() {
	if h.rootPath == "" {
		return
	}

	filesInfo, err := core.CollectFiles(h.rootPath)
	if err != nil {
		core.DebugLog("Workspace file discovery failed: %v", err)
		return
	}

	h.mu.Lock()
	h.eval = evaluator.NewEvaluator(filesInfo)
	eval := h.eval
	h.mu.Unlock()

	// 1. Parse all files from disk
	h.evalMu.Lock()
	err = eval.ParseAll(h.rootPath)
	if err != nil {
		core.DebugLog("Initial parse failed: %v", err)
	}

	// 2. Run analysis
	eval.RunAnalysis()

	// 3. Collect diagnostics for publish
	diagByPath := make(map[string][]Diagnostic, len(filesInfo))
	for _, info := range filesInfo {
		path := info.Path()
		diagByPath[path] = alertsToDiagnostics(toURI(filepath.Join(h.rootPath, path)), eval.GetAlerts(path))
	}
	h.evalMu.Unlock()

	// 4. Store file contents and publish diagnostics
	for _, info := range filesInfo {
		path := info.Path()
		uri := toURI(filepath.Join(h.rootPath, path))

		content, err := os.ReadFile(filepath.Join(h.rootPath, path))
		if err == nil {
			h.mu.Lock()
			h.files[uri] = &File{
				LanguageID: "hybroid",
				Text:       string(content),
				Version:    0,
			}
			h.mu.Unlock()
		}

		h.conn.Notify(context.Background(), "textDocument/publishDiagnostics", PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagByPath[path],
		})
	}

	h.markReady()
	core.DebugLog("Workspace pre-analysis complete. Analyzed %d files.", len(filesInfo))
}

// publishInfoOnce sends a single Information-severity diagnostic to the client
// for the given URI, the first time it is called for that URI. Subsequent
// calls for the same URI are a no-op. This is used to surface "your file is
// open without a project context" hints exactly once per buffer, so the user
// is informed without being re-pinged on every keystroke.
func (h *langHandler) publishInfoOnce(ctx context.Context, conn notifier, uri DocumentURI, message string) {
	h.mu.Lock()
	if _, ok := h.infoNoticesPublished[uri]; ok {
		h.mu.Unlock()
		return
	}
	h.infoNoticesPublished[uri] = struct{}{}
	connRef := h.conn
	h.mu.Unlock()

	if connRef == nil {
		return
	}

	severity := 3 // LSP DiagnosticSeverity.Information
	connRef.Notify(ctx, "textDocument/publishDiagnostics", PublishDiagnosticsParams{
		URI: uri,
		Diagnostics: []Diagnostic{
			{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 0},
				},
				Severity: severity,
				Message:  message,
				Source:   func() *string { s := "hybroid"; return &s }(),
			},
		},
	})
}
