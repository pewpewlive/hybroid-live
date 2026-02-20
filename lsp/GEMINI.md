# Hybroid LSP Implementation Details

This directory contains the implementation of the Language Server Protocol (LSP) for the Hybroid programming language.

## Core Logic (`lsp.go`, `analysis.go`, `handler.go`)

- **`NewHandler()`**: Initializes a new `langHandler` which manages the state of open files, analyzed walkers, and LSP capabilities.
- **`Analyze(uri DocumentURI, text string) AnalysisResult`**: The main entry point for semantic analysis. It runs the Lexer, Parser, and Walker on the provided text and collects diagnostics.
- **`alertsToDiagnostics(uri DocumentURI, alertsList []alerts.Alert) []Diagnostic`**: Converts internal compiler alerts into LSP-compliant diagnostics for error reporting in the editor.
- **`langHandler.handle(ctx, conn, req)`**: The main multiplexer for incoming JSON-RPC requests, dispatching them to specific handler functions.

## Lifecycle Handlers (`handle_initialize.go`, `handle_shutdown.go`, `init.go`)

- **`Init(debug bool)`**: Sets up the LSP server, configures logging (to `hybroid_ls.log.txt` if debug is enabled), and starts the JSON-RPC connection over stdio.
- **`handleInitialize(ctx, conn, req)`**: Negotiates capabilities with the client (e.g., VS Code), enabling features like Hover, Completion, and Signature Help.
- **`handleShutdown(ctx, conn, req)`**: Gracefully accepts the `shutdown` signal (keeping the socket active via `nil, nil`) allowing the client to subsequently emit the protocol-mandated `exit` event for a safe OS termination.

## Document Syncing (`handle_text_document_did_change.go`)

- **`handleTextDocumentDidOpen(ctx, conn, req)`**: triggered when a file is opened in the editor. It stores the file content and performs initial analysis.
- **`handleTextDocumentDidChange(ctx, conn, req)`**: Triggered as the user types. It updates the in-memory file content and re-runs analysis to provide live diagnostics.
- **`analyzeAndPublish(ctx, conn, uri, text)`**: Helper that runs `Analyze` and pushes the resulting diagnostics back to the client.

## Language Features

### Hover (`handle_hover.go`, `metadata.go`)

- **`handleTextDocumentHover(ctx, conn, req)`**: Provides tooltips for symbols. It identifies the word under the cursor and looks up metadata from both static documentation and the semantic walker's scope.
- **`getSymbolMetadata(label string)`**: Look up documentation/types for keywords, builtins, and namespace-qualified symbols (e.g., `Fmath:RandomFixed`).

### Completion (`handle_completion.go`, `handle_completion_item_resolve.go`)

- **`handleTextDocumentCompletion(ctx, conn, req)`**: Generates a list of suggestions. Includes keywords, native types, namespaces, builtin functions, and variables currently in scope.
- **`HandleCompletionItemResolve(ctx, conn, req)`**: provides additional details (like documentation) when a user selects a completion item from the list.

### Signature Help (`handle_signature_help.go`)

- **`handleTextDocumentSignatureHelp(ctx, conn, req)`**: Displays function parameters while typing a call. It parses the context to find the function name and the active parameter index.

### Rename (`handle_rename.go`)

- **`handleTextDocumentRename(ctx, conn, req)`**: Executes workspace-wide identifier refactoring. Uses `findReferences` dynamically passing the `rootDir` override to verify target `DocumentURI` pathways. Generates `WorkspaceEdit` objects containing specific text substitutions. Seamlessly captures structural symbols like `class Rectangle {}` alongside all of their derived expressions like `new Rectangle()` and `spawn ExampleEntity()`, injecting semantic trace references mapping perfectly to both absolute workspace and relative stray-file namespaces.

## Helpers (`helpers.go`)

- **`isInCommentOrString(text string, line, col int) bool`**: Scans the text to determine if a specific position is inside a comment or string literal, used to suppress features like Hover and Signature Help in those contexts.
- **`IsWordChar(r rune) bool`**: Defines what constitutes a "word" for symbol lookup (alphanumeric, underscores, and colons for namespaces).
- **`fromURI(uri DocumentURI)` / `toURI(path string)`**: Utilities for converting between LSP file URIs and local filesystem paths.

## Logging & Debugging

- **`core.DebugLog(format, v...)`**: Conditional logging that only outputs to a file if the server was started with the `--debug` flag.
