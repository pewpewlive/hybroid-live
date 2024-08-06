package lsp

import (
	"context"
	"fmt"
	"log"
	"net/url"
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

type File struct {
	LanguageID string
	Text       string
	Version    int
}

type langHandler struct {
	mu     sync.Mutex
	logger *log.Logger
	// commands          []Command
	provideDefinition bool
	files             map[DocumentURI]*File
	lintDebounce      time.Duration
	request           chan lintRequest
	lintTimer         *time.Timer
	formatDebounce    time.Duration
	formatTimer       *time.Timer
	conn              *jsonrpc2.Conn
	rootPath          string
	filename          string
	folders           []string
	rootMarkers       []string
	triggerChars      []string

	// lastPublishedURIs is mapping from LanguageID string to mapping of
	// whether diagnostics are published in a DocumentURI or not.
	lastPublishedURIs map[string]map[DocumentURI]struct{}
}

func (h *langHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, conn, req)
	case "initialized":
		return
	case "shutdown":
		return h.handleShutdown(ctx, conn, req)
	case "textDocument/didOpen":
		return // h.handleTextDocumentDidOpen(ctx, conn, req)
	case "textDocument/didChange":
		return // h.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/didSave":
		return // h.handleTextDocumentDidSave(ctx, conn, req)
	case "textDocument/didClose":
		return // h.handleTextDocumentDidClose(ctx, conn, req)
	case "textDocument/formatting":
		return // h.handleTextDocumentFormatting(ctx, conn, req)
	case "textDocument/rangeFormatting":
		return // h.handleTextDocumentRangeFormatting(ctx, conn, req)
	case "textDocument/documentSymbol":
		return // h.handleTextDocumentSymbol(ctx, conn, req)
	case "textDocument/completion":
		return h.handleTextDocumentCompletion(ctx, conn, req)
	case "completionItem/resolve":
		return h.HandleCompletionItemResolve(ctx, conn, req)
	case "textDocument/definition":
		return // h.handleTextDocumentDefinition(ctx, conn, req)
	case "textDocument/hover":
		return // h.handleTextDocumentHover(ctx, conn, req)
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
		// provideDefinition: config.ProvideDefinition,
		files:   make(map[DocumentURI]*File),
		request: make(chan lintRequest),
		conn:    nil,
		// filename:          config.Filename,
		// rootMarkers:       *config.RootMarkers,
		// triggerChars:      config.TriggerChars,

		lastPublishedURIs: make(map[string]map[DocumentURI]struct{}),
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
