package lsp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleInitialize(_ context.Context, conn notifier, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	h.conn = conn

	var params InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	// https://microsoft.github.io/language-server-protocol/specification#initialize
	// The rootUri of the workspace. Is null if no folder is open.
	if params.RootURI != "" {
		h.rootURI = params.RootURI
		workspacePath, err := fromURI(params.RootURI)
		if err != nil {
			return nil, err
		}
		workspacePath = filepath.Clean(workspacePath)
		h.addFolder(workspacePath)

		// A VS Code workspace can be broader than a Hybroid project (for
		// example, this repository contains projects under examples/). Only
		// pre-analyze here when the workspace itself is inside a project. If
		// it is a parent of one or more projects, didOpen will select the
		// nearest hybconfig.toml for the opened file instead of combining all
		// nested projects and test fixtures into one evaluator.
		projectRoot := findProjectRootFromDir(workspacePath, h.rootMarkers)
		if projectRoot != "" {
			h.rootPath = filepath.Clean(projectRoot)
			h.startPreAnalysis()
		}
	}

	var completion *CompletionProvider
	// var hasCompletionCommand bool
	var hasCodeActionCommand bool
	var hasSymbolCommand bool
	var hasFormatCommand bool
	var hasRangeFormatCommand bool

	if params.InitializationOptions != nil {
		//hasCompletionCommand = params.InitializationOptions.Completion
		hasCodeActionCommand = params.InitializationOptions.CodeAction
		hasSymbolCommand = params.InitializationOptions.DocumentSymbol
		hasFormatCommand = params.InitializationOptions.DocumentFormatting
		hasRangeFormatCommand = params.InitializationOptions.RangeFormatting
	}

	completion = &CompletionProvider{
		ResolveProvider:   true,
		TriggerCharacters: []string{":", "."},
	}
	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync:           TDSKFull,
			DocumentFormattingProvider: hasFormatCommand,
			RangeFormattingProvider:    hasRangeFormatCommand,
			DocumentSymbolProvider:     hasSymbolCommand,
			DefinitionProvider:         true,
			ReferencesProvider:         true,
			RenameProvider:             true,
			CompletionProvider:         completion,
			SignatureHelpProvider: &SignatureHelpProvider{
				TriggerCharacters: []string{"(", ","},
			},
			HoverProvider:      true,
			CodeActionProvider: hasCodeActionCommand,
			Workspace: &ServerCapabilitiesWorkspace{
				WorkspaceFolders: WorkspaceFoldersServerCapabilities{
					Supported:           true,
					ChangeNotifications: true,
				},
			},
		},
	}, nil
}
