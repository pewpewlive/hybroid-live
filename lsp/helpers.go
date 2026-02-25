package lsp

import (
	"hybroid/tokens"
	"path/filepath"
	"strings"
)

// isInCommentOrString checks if the given position is inside a comment or a string.
// This is a simplified version that doesn't use the full lexer for performance.
func isInCommentOrString(text string, line, col int) bool {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if line < 0 || line >= len(lines) {
		return false
	}

	// Check for single line comments and strings first
	l := lines[line]
	if col > len(l) {
		col = len(l)
	}

	// Simple state machine
	isComment := false
	isMultilineComment := false
	isString := false

	for i := 0; i <= line; i++ {
		segment := lines[i]
		if i == line {
			segment = segment[:col]
		}
		runes := []rune(segment)
		for j := 0; j < len(runes); j++ {
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

			// Not in any special state
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
	if character < 0 || character >= len(l) {
		return ""
	}

	start := character
	for start > 0 && IsWordChar(rune(l[start-1])) {
		start--
	}
	end := character
	for end < len(l) && IsWordChar(rune(l[end])) {
		end++
	}

	if start == end {
		return ""
	}

	return l[start:end]
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
