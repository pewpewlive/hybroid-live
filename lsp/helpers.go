package lsp

import (
	"hybroid/ast"
	"hybroid/tokens"
	"hybroid/walker"
	"path/filepath"
	"strings"
)

// isInCommentOrString checks if the given position is inside a comment or a string.
// This is a simplified version that doesn't use the full lexer for performance.
// `col` is treated as a rune index (matching how callers pass Position.Character).
func isInCommentOrString(text string, line, col int) bool {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if line < 0 || line >= len(lines) {
		return false
	}

	// Simple state machine
	isComment := false
	isMultilineComment := false
	isString := false

	for i := 0; i <= line; i++ {
		runes := []rune(lines[i])
		segmentLen := len(runes)
		if i == line {
			if col > segmentLen {
				col = segmentLen
			}
			segmentLen = col
		}
		for j := 0; j < segmentLen; j++ {
			c := runes[j]

			if isComment {
				continue
			}

			if isMultilineComment {
				if c == '*' && j+1 < len(runes) && runes[j+1] == '/' {
					isMultilineComment = false
					j++
				}
				continue
			}

			if isString {
				if c == '\\' && j+1 < len(runes) && runes[j+1] == '"' {
					j++
					continue
				}
				if c == '"' {
					isString = false
				}
				continue
			}

			if c == '/' && j+1 < len(runes) {
				if runes[j+1] == '/' {
					isComment = true
					break
				}
				if runes[j+1] == '*' {
					isMultilineComment = true
					j++
					continue
				}
			}

			if c == '"' {
				isString = true
				continue
			}
		}
		isComment = false
	}
	return isComment || isMultilineComment || isString
}

func IsWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == ':'
}

func getWordAt(text string, line, character int) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return ""
	}
	l := lines[line]
	if character < 0 {
		return ""
	}

	runes := []rune(l)
	if character >= len(runes) {
		return ""
	}

	start := character
	for start > 0 && IsWordChar(runes[start-1]) {
		start--
	}
	end := character
	for end < len(runes) && IsWordChar(runes[end]) {
		end++
	}

	if start == end {
		return ""
	}

	return string(runes[start:end])
}

func toLSPLocation(path string, token tokens.Token) Location {
	return Location{
		URI: toURI(path),
		Range: Range{
			Start: Position{
				Line:      token.Line - 1,
				Character: token.Column.Start - 1,
			},
			End: Position{
				Line:      token.Line - 1,
				Character: token.Column.End - 1,
			},
		},
	}
}

// getRelPath calculates the relative path from base to targ.
// If base is empty, or if an error occurs during Rel evaluation,
// the targ path is safely returned as a fallback to support single-file workspaces.
func getRelPath(base, targ string) string {
	if base == "" {
		return targ
	}
	rel, err := filepath.Rel(base, targ)
	if err != nil {
		return targ
	}
	return rel
}

// resolveBuiltinEnvByName returns the built-in library environment for the
// given user-facing namespace name (Pewpew, Fmath, Math, String, Table), or
// nil if the name does not match a built-in library. Centralises the switch
// that was previously duplicated across multiple handlers.
func resolveBuiltinEnvByName(name string) *walker.Environment {
	switch name {
	case "Pewpew":
		return walker.PewpewAPI
	case "Fmath":
		return walker.FmathAPI
	case "Math":
		return walker.MathAPI
	case "String":
		return walker.StringAPI
	case "Table":
		return walker.TableAPI
	}
	return nil
}

// resolveBuiltinEnv returns the built-in library environment for the given
// ast.Library value, or nil if the value is not recognised.
func resolveBuiltinEnv(lib ast.Library) *walker.Environment {
	switch lib {
	case ast.Pewpew:
		return walker.PewpewAPI
	case ast.Fmath:
		return walker.FmathAPI
	case ast.Math:
		return walker.MathAPI
	case ast.String:
		return walker.StringAPI
	case ast.Table:
		return walker.TableAPI
	}
	return nil
}
