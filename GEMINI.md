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
*   `parser/`: recursive descent parser.
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
