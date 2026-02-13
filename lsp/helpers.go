package lsp

import (
	"hybroid/tokens"
	"strings"
)

// isInCommentOrString checks if the given position is inside a comment or a string.
// This is a simplified version that doesn't use the full lexer for performance.
func isInCommentOrString(text string, line, col int) bool {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if line < 0 || line >= len(lines) {
		return false
	}

	// For comments/strings we might need to look at previous lines for multiline comments.
	// However, Hybroid multiline comments are /* */ and strings are single-line mostly (lexer alerts on multiline).
	
	// Check for single line comments and strings first
	l := lines[line]
	if col > len(l) {
		col = len(l)
	}

	inString := false
	inComment := false
	
	// We need to scan from the start of the file for multiline comments
	// or at least from a safe point. For now, let's scan the whole text up to the point.
	
	fullTextBefore := ""
	for i := 0; i < line; i++ {
		fullTextBefore += lines[i] + "\n"
	}
	fullTextBefore += l[:col]

	// Simple state machine
	isComment := false
	isMultilineComment := false
	isString := false
	
	runes := []rune(fullTextBefore)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		
		if isComment {
			if c == '\n' {
				isComment = false
			}
			continue
		}
		
		if isMultilineComment {
			if c == '*' && i+1 < len(runes) && runes[i+1] == '/' {
				isMultilineComment = false
				i++
			}
			continue
		}
		
		if isString {
			if c == '\\' && i+1 < len(runes) && runes[i+1] == '"' {
				i++
				continue
			}
			if c == '"' {
				isString = false
			}
			continue
		}
		
		// Not in any special state
		if c == '/' && i+1 < len(runes) {
			if runes[i+1] == '/' {
				isComment = true
				i++
				continue
			}
			if runes[i+1] == '*' {
				isMultilineComment = true
				i++
				continue
			}
		}
		
		if c == '"' {
			isString = true
			continue
		}
	}
	
	inComment = isComment || isMultilineComment
	inString = isString
	
	return inComment || inString
}

func IsWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == ':'
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
				Character: token.Column.End,
			},
		},
	}
}
