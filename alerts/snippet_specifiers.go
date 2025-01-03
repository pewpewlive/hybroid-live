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
	line := src[index-columnCount+1 : index+1]
	location := ss.Token.Location
	if location.ColStart > 80 { // ok
		var content string
		short := false
		if location.ColEnd+20 > columnCount && columnCount < 120 {
			content = string(line[location.ColStart-20 : columnCount])
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
	StartToken tokens.Token
	EndToken   tokens.Token
}

func (ml *Multiline) GetSnippet(src string, index, columnCount, lineCount int) string {
	snippet := strings.Builder{}
	first_line := src[index-columnCount+1 : index-1]
	startLocation := ml.StartToken.Location
	endLocation := ml.EndToken.Location

	diff := startLocation.LineEnd - endLocation.LineStart

	if diff == 0 {
		singleline := &Singleline{Token: tokens.Token{Location: tokens.TokenLocation{
			LineStart: startLocation.LineStart,
			LineEnd:   startLocation.LineEnd,
			ColStart:  startLocation.ColStart,
			ColEnd:    endLocation.ColEnd,
		}}}

		return singleline.GetSnippet(src, index, columnCount, lineCount)
	}

	snippet.WriteString(fmt.Sprintf("[cyan]%*s |\n", len(strconv.Itoa(startLocation.LineEnd)), ""))
	snippet.WriteString(fmt.Sprintf("[cyan]%*d |[default]   %s\n", len(strconv.Itoa(startLocation.LineEnd)), lineCount, string(first_line)))
	snippet.WriteString(fmt.Sprintf("[cyan]%*s |[light_red]  %s^\n", len(strconv.Itoa(startLocation.LineEnd)), "", strings.Repeat("_", startLocation.ColStart)))

	for i := index; i <= len(src)-1; i++ {
		columnCount++
		index = i
		//print(string(src[i]))

		if endLocation.LineStart == lineCount && src[i] == '\n' {
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
	snippet.WriteString(fmt.Sprintf("[cyan]%s |[light_red] |%s^\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat("_", endLocation.ColEnd)))

	return snippet.String()
}

func (ml *Multiline) GetTokens() []tokens.Token {
	return []tokens.Token{ml.StartToken, ml.EndToken}
}
