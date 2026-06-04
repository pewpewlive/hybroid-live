# Hybroid Live - Project Context

## Project Overview

**Hybroid Live** is a statically-typed programming language designed specifically for creating content for the game **PewPew Live**. It transpiles to Lua, which is the scripting language used by the game.

**Key Goals:**
*   Provide a better developer experience than raw Lua.
*   Add features missing in Lua (classes, enums, strict typing).
*   Optimize for PewPew Live specifics (tick loops, fixed-point math).
*   Provide robust error messages.

**Status:** Alpha (expect breaking changes).

## Architecture

The project is written in **Go** and follows a standard compiler architecture:
1.  **Lexer (`lexer/`)**: Tokenizes source code.
2.  **Parser (`parser/`)**: Constructs the AST (`ast/`).
3.  **Walker (`walker/`)**: Performs semantic analysis, type checking, and scope resolution.
4.  **Generator (`generator/`)**: Transpiles the AST into Lua code.
5.  **LSP (`lsp/`)**: Provides Language Server Protocol support for editors (VS Code).

### LSP single-file fallback

When the editor opens a `.hyb` file without a containing folder (no `rootUri` in the `initialize` request), the LSP walks the parent directories looking for `hybconfig.toml` — the same pattern that `tsserver` uses to discover a `tsconfig.json` for loose files. If a project root is found, full workspace analysis runs as if the folder had been opened. If no marker is found, the file is analyzed in isolation and a single Information diagnostic is published at the top of the buffer: *"This file is open without its Hybroid project. Open the folder containing `hybconfig.toml` to resolve all `use` references."* See `lsp/find_project_root.go` and `lsp/handle_text_document_did_change.go`.

## Development Environments

Hybroid supports distinct "environments" that dictate available standard libraries and compilation behavior:
*   `Level`: For game levels (access to `fmath`, restricted Lua stdlib).
*   `Mesh`: For generating meshes (full math support).
*   `Sound`: For generating sounds.

## Build & Usage

### Prerequisites
*   Go 1.23+
*   Python 3 (for build scripts)

### Building the CLI
To build the native CLI executable:
```bash
go build -o hybroid main_native.go
```

### Running the CLI
```bash
./hybroid <command> [arguments]
```
Common commands:
*   `init`: Initialize a new Hybroid project.
*   `build`: Compile Hybroid code to Lua.
*   `watch`: Watch for changes and recompile.

### Building for Release
Use the Python helper script:
```bash
python utils/build_hybroid.py
```

### Testing
Run standard Go tests:
```bash
go test ./...
```

### Releasing the LSP alpha

1.  Tag the merge commit on `master` as `v0.1.0-lsp-alpha`.
2.  Run `python utils/build_hybroid.py` to produce binaries for all platforms in `./build/`. Upload each to the GitHub Release with the name `hybroid-<platform><.exe>`.
3.  The `install.sh` / `install.ps1` scripts at hybroid.pewpew.live fetch the matching asset and copy it to `~/.hybroid/hybroid`.
4.  The VS Code extension is packaged and published separately in the `hybroid-vscode` repo. Bump its `version` in `package.json` to `0.1.0` and run `vsce package`.
5.  No new CI workflow is added in this repo — the existing `.github/workflows/go.yml` validates builds and tests on push to master, which is sufficient for the alpha.

### Install location and logs

The Hybroid CLI/LSP binary and the LSP debug log both live under `~/.hybroid/`:

*   Binary: `~/.hybroid/hybroid` (or `~/.hybroid/hybroid.exe` on Windows)
*   LSP log (debug mode only): `~/.hybroid/logs/lsp.log`

The `HYBROID_LS_LOG` environment variable overrides the log path. On macOS, VS Code has a sanitized PATH that does not include `~/.hybroid`, so the extension searches there explicitly. The `hybroid.languageServerPath` user setting is the override if both mechanisms fail. See `lsp/logpath.go` for the resolution contract and `vscode-ext/src/path-resolver.ts` (in the submodule) for the binary search chain.

## Directory Structure

*   `alerts/`: Error reporting system (diagnostics, pretty printing).
*   `ast/`: Abstract Syntax Tree node definitions.
*   `cli/`: Implementation of CLI commands.
*   `core/`: Core data structures (Queue, Stack, Span).
*   `docs/`: Documentation website (Astro).
*   `evaluator/`: Constant evaluation and testing logic.
*   `examples/`: Sample Hybroid projects and code snippets.
*   `generator/`: Lua code generation logic.
*   `lexer/`: Source code tokenization.
*   `lsp/`: Language Server Protocol implementation.
*   `parser/`: Recursive descent parser.
*   `tokens/`: Token type definitions.
*   `utils/`: Python scripts for build, API generation, and maintenance.
*   `walker/`: Semantic analysis and type checking.
*   `wasm/`: WASM bindings for the web playground.

## Key Files

*   `main_native.go`: Entry point for the CLI.
*   `spec.md`: The Hybroid language specification.
*   `hybconfig.toml`: Project configuration file (seen in examples).
*   `go.mod`: Go dependencies.

## Conventions

*   **Language:** Go for the toolchain, Python for build scripts.
*   **Style:** Follows standard Go formatting (`gofmt`).
*   **Testing:** Uses Go's built-in testing framework.
*   **Documentation:** Maintained in `docs/` and `spec.md`.
