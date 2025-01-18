package alerts

import (
	"fmt"
	"hybroid/tokens"
	"strconv"
	"strings"
)

type SnippetSpecifier interface {
	GetSnippet(src string, index, columnCount, lineCount int) string
	GetTokens() []tokens.Token
}

func NewSingle(token tokens.Token) *Singleline {
	return &Singleline{
		Token: token,
	}
}

type Singleline struct {
	Token tokens.Token
}

func (ss *Singleline) GetSnippet(src string, index, columnCount, lineCount int) string {
	snippet := strings.Builder{} // how are we getting the error line then
	line := src[index-columnCount+1 : index]
	loc := ss.Token.TokenLocation
	if loc.Column.Start > 80 { // ok
		var content string
		short := false
		if loc.Column.End+20 > columnCount && columnCount < 120 {
			content = string(line[loc.Column.Start-20 : columnCount])
		} else {
			content = string(line[loc.Column.Start-20 : loc.Column.End+20])
			short = true
		}
		var shortt string
		if short {
			shortt = "..."
		}
		// TODO: move prints into separate functions
		snippet.WriteString(fmt.Sprintf("[cyan]%*s |\n", len(strconv.Itoa(loc.Line.End)), ""))
		snippet.WriteString(fmt.Sprintf("[cyan]%d |[default]   [dark_gray]...[default]%s[dark_gray]%s[reset]\n", lineCount, content, shortt))
		snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red]   %s%s\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat(" ", 22), strings.Repeat("^", loc.Column.End-loc.Column.Start+1)))
	} else {
		snippet.WriteString(fmt.Sprintf("[cyan]%*s |\n", len(strconv.Itoa(loc.Line.End)), ""))
		snippet.WriteString(fmt.Sprintf("[cyan]%d |[default]   %s\n", lineCount, line))
		snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red]   %s%s\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat(" ", loc.Column.Start-1), strings.Repeat("^", loc.Column.End-loc.Column.Start+1)))
	}

	return snippet.String()
}

func (ss *Singleline) GetTokens() []tokens.Token {
	return []tokens.Token{ss.Token}
}

func NewMulti(startToken, endToken tokens.Token) *Multiline {
	return &Multiline{
		StartToken: startToken,
		EndToken:   endToken,
	}
}

type Multiline struct {
	StartToken tokens.Token
	EndToken   tokens.Token
}

func (ml *Multiline) GetSnippet(src string, index, columnCount, lineCount int) string {
	snippet := strings.Builder{}
	first_line := src[index-columnCount+1 : index-1]
	startLoc := ml.StartToken.TokenLocation
	endLoc := ml.EndToken.TokenLocation

	diff := startLoc.Line.End - endLoc.Line.Start

	if diff == 0 {
		singleline := NewSingle(tokens.Token{TokenLocation: tokens.NewLocation(
			startLoc.Line.Start,
			startLoc.Line.End,
			startLoc.Column.Start,
			endLoc.Column.End, 0, 0)})

		return singleline.GetSnippet(src, index, columnCount, lineCount)
	}

	snippet.WriteString(fmt.Sprintf("[cyan]%*s |\n", len(strconv.Itoa(startLoc.Line.End)), ""))
	snippet.WriteString(fmt.Sprintf("[cyan]%*d |[default]   %s\n", len(strconv.Itoa(startLoc.Line.End)), lineCount, string(first_line)))
	snippet.WriteString(fmt.Sprintf("[cyan]%*s |[light_red]  %s^\n", len(strconv.Itoa(startLoc.Line.End)), "", strings.Repeat("_", startLoc.Column.Start)))

	for i := index; i <= len(src)-1; i++ {
		columnCount++
		index = i
		//print(string(src[i]))

		if endLoc.Line.Start == lineCount && src[i] == '\n' {
			index--
			columnCount--
			//println("\nhappened")
			break
		}

		if src[i] == '\n' {
			//println("\nhappened2")
			columnCount = 0
			lineCount++
		}
	}

	last_line := src[index-columnCount+1 : index+1]

	//color.Printf("[cyan]%*s |\n", len(strconv.Itoa(location.LineEnd)), "")
	if diff != 1 {
		snippet.WriteString(fmt.Sprintf("[dark_gray]...[default]%s[light_red]|\n", strings.Repeat(" ", len(strconv.Itoa(lineCount)))))
	}
	snippet.WriteString(fmt.Sprintf("[cyan]%d |[light_red] | [default]%s\n", lineCount, string(last_line)))
	snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red] |%s^\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat("_", endLoc.Column.End)))

	return snippet.String()
}

func (ml *Multiline) GetTokens() []tokens.Token {
	return []tokens.Token{ml.StartToken, ml.EndToken}
}
