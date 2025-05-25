package alerts

import (
	"fmt"
	"hybroid/tokens"
	"strconv"
	"strings"
)

type Snippet interface {
	GetSnippet(lines map[int][]byte) string
	GetTokens() []tokens.Token
}

// Tries to truncate the line if possible and writes it.
// Returns the truncated string's marker and alignment sizes
func writeTruncatedLine(snippet *strings.Builder, loc tokens.Location, line []byte) (int, int) {
	const startLimit = 80
	const midLimit = 60
	const endLimit = 50
	const segmentSize = 15
	const truncStr = "[dark_gray]...[default]"

	start, end := loc.Column.Start, loc.Column.End
	lineSize := len(line)

	leadingSpace := start - 1
	markerSize := end - start

	lineStart, lineMiddle, lineEnd := line[:start-1], line[start-1:end-1], line[end-1:]

	// Truncate the start:  ... + "Hello, Hybroid!"
	if start > startLimit {
		leadingSpace = 4 + segmentSize // `4` for "... " and `segmentSize` for the additional portion

		snippet.WriteString(truncStr + " ")
		snippet.Write(lineStart[len(lineStart)-segmentSize:])
	} else {
		snippet.Write(lineStart)
	}

	// Truncate the middle:  "Hello ... Hybroid!"
	if end-start > midLimit {
		markerSize = 5 + 2*segmentSize // `5` for " ... " and `2*segmentSize` for the additional portions

		snippet.Write(lineMiddle[:segmentSize])
		snippet.WriteString(" " + truncStr + " ")
		snippet.Write(lineMiddle[len(lineMiddle)-segmentSize:])
	} else {
		snippet.Write(lineMiddle)
	}

	// Truncate the end:  "Hello" + ...
	if lineSize-end > endLimit {
		snippet.Write(lineEnd[:segmentSize])
		snippet.WriteString(" " + truncStr)
	} else {
		snippet.Write(lineEnd)
	}

	snippet.WriteByte('\n')

	return leadingSpace, markerSize
}

type SingleLine struct {
	Token tokens.Token
}

func NewSingle(token tokens.Token) SingleLine {
	return SingleLine{
		Token: token,
	}
}

func (ss SingleLine) GetSnippet(lines map[int][]byte) string {
	snippet := strings.Builder{}
	loc := ss.Token.Location
	line := lines[loc.Line]

	lineNumberSpaces := strings.Repeat(" ", len(strconv.Itoa(loc.Line)))

	snippet.WriteString(fmt.Sprintf("[cyan]%s |\n", lineNumberSpaces))
	snippet.WriteString(fmt.Sprintf("[cyan]%d |[default]   ", loc.Line))
	spaceSize, markerSize := writeTruncatedLine(&snippet, loc, line)
	leadingSpace, marker := strings.Repeat(" ", spaceSize), strings.Repeat("^", markerSize)
	if ss.Token.Type == tokens.Eof {
		leadingSpace = strings.Repeat(" ", len(line))
		marker = "^- End Of File"
	}
	snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red]   %s%s\n", lineNumberSpaces, leadingSpace, marker))

	return snippet.String()
}

func (ss SingleLine) GetTokens() []tokens.Token {
	return []tokens.Token{ss.Token}
}

type MultiLine struct {
	StartToken tokens.Token
	EndToken   tokens.Token
}

func NewMulti(startToken, endToken tokens.Token) MultiLine {
	return MultiLine{
		StartToken: startToken,
		EndToken:   endToken,
	}
}

func (ml MultiLine) GetSnippet(lines map[int][]byte) string {
	snippet := strings.Builder{}
	startLoc, endLoc := ml.StartToken.Location, ml.EndToken.Location
	startLine, endLine := lines[startLoc.Line], lines[endLoc.Line]

	largestLineNumber := max(len(strconv.Itoa(startLoc.Line)), len(strconv.Itoa(endLoc.Line)))
	lineNumberSpaces := strings.Repeat(" ", largestLineNumber)

	snippet.WriteString(fmt.Sprintf("[cyan]%s |\n", lineNumberSpaces))
	snippet.WriteString(fmt.Sprintf("[cyan]%*d |[default]   ", largestLineNumber, startLoc.Line))
	spaceSize, markerSize := writeTruncatedLine(&snippet, startLoc, startLine)
	startHorizLine, marker := strings.Repeat("_", spaceSize+1), strings.Repeat("^", markerSize)
	snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red]  %s%s\n", lineNumberSpaces, startHorizLine, marker))

	if endLoc.Line-startLoc.Line != 1 {
		ellipsisAlignment := strings.Repeat(" ", largestLineNumber-1)
		snippet.WriteString(fmt.Sprintf("[dark_gray]%s...[light_red] |\n", ellipsisAlignment))
	}

	snippet.WriteString(fmt.Sprintf("[cyan]%*d |[light_red] | [default]", largestLineNumber, endLoc.Line))
	spaceSize, markerSize = writeTruncatedLine(&snippet, endLoc, endLine)
	endHorizLine, marker := strings.Repeat("_", spaceSize+1), strings.Repeat("^", markerSize)
	if ml.EndToken.Type == tokens.Eof {
		endHorizLine = strings.Repeat("_", len(endLine))
		marker = "_> End Of File"
	}
	snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red] |%s%s\n", lineNumberSpaces, endHorizLine, marker))

	return snippet.String()
}

func (ml MultiLine) GetTokens() []tokens.Token {
	return []tokens.Token{ml.StartToken, ml.EndToken}
}
