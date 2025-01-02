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

type Singleline struct { // Alert(alerts.DoesNotExistException{}, Singleline{token}, "your params")
	Token tokens.Token
}

func (ss *Singleline) GetSnippet(src string, index, columnCount, lineCount int) string {
	snippet := strings.Builder{} // how are we getting the error line then
	line := src[index-columnCount+1 : index-1]
	location := ss.Token.Location
	if location.ColStart > 80 { // ok
		var content string
		short := false
		if location.ColEnd+20 > columnCount && columnCount < 120 {
			content = string(line[location.ColStart-20 : columnCount-1])
		} else {
			content = string(line[location.ColStart-20 : location.ColEnd+20])
			short = true
		}
		var shortt string
		if short {
			shortt = "..."
		}
		// TODO: move prints into separate functions
		snippet.WriteString(fmt.Sprintf("[cyan]%*s |\n", len(strconv.Itoa(location.LineEnd)), ""))
		snippet.WriteString(fmt.Sprintf("[cyan]%d |[default]   [dark_gray]...[default]%s[dark_gray]%s[reset]\n", lineCount, content, shortt))
		snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red]   %s%s\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat(" ", 22), strings.Repeat("^", location.ColEnd-location.ColStart+1)))
	} else {
		snippet.WriteString(fmt.Sprintf("[cyan]%*s |\n", len(strconv.Itoa(location.LineEnd)), ""))
		snippet.WriteString(fmt.Sprintf("[cyan]%d |[default]   %s\n", lineCount, line))
		snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red]   %s%s\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat(" ", location.ColStart-1), strings.Repeat("^", location.ColEnd-location.ColStart+1)))
	}

	return snippet.String()
}

func (ss *Singleline) GetTokens() []tokens.Token {
	return []tokens.Token{ss.Token}
}

type Multiline struct {
}

func (ml *Multiline) GetSnippet(src string, index, columnCount, lineCount int) string {
	return ""
}

func (ml *Multiline) GetTokens() []tokens.Token {
	return []tokens.Token{}
}
