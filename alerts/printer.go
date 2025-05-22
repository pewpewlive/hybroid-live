package alerts

import (
	"bufio"
	"fmt"
	"hybroid/tokens"
	"os"
	"sort"
	"strconv"
	"strings"

	color "github.com/mitchellh/colorstring"
)

type lineBuffer map[int][]byte

func lineRangeFromSpan(span []int) (int, int) {
	start := span[0]
	end := start
	if len(span) == 2 {
		end = span[1]
	}

	return start, end
}

func (lb lineBuffer) appendLine(line int, buffer []byte) {
	lb[line] = make([]byte, len(buffer))
	copy(lb[line], buffer)
}

func (lb lineBuffer) freeLines(lines ...int) {
	if len(lb) == 0 {
		return
	}

	start, end := lineRangeFromSpan(lines)

	for line := start; line < end+1; line++ {
		delete(lb, line)
	}
}

func (lb lineBuffer) retrieveLines(lines ...int) lineBuffer {
	buffer := make(lineBuffer)

	if len(lb) == 0 || len(lines) == 0 {
		return buffer
	}

	for _, line := range lines {
		if bufferLine, ok := lb[line]; ok {
			buffer[line] = bufferLine
		}
	}
	return buffer
}

func (lb lineBuffer) isAvailable(lines ...int) bool {
	if len(lb) == 0 {
		return false
	}

	for _, line := range lines {
		if _, ok := lb[line]; !ok {
			return false
		}
	}
	return true
}

type Printer struct {
	alertsByFile map[string][]Alert

	// The current scanner of a file
	source *bufio.Scanner
	// Line buffer object that keeps track of necessary lines
	buffer lineBuffer
	// Keep track of the current line
	line int
	// Keep track of lines that are no longer necessary
	freeStart int
}

func NewPrinter() Printer {
	return Printer{
		alertsByFile: make(map[string][]Alert),
		source:       nil,
		buffer:       nil,
		line:         -1,
		freeStart:    -1,
	}
}

func (p *Printer) StageAlerts(sourcePath string, alerts []Alert) {
	fileAlerts, existed := p.alertsByFile[sourcePath]
	if !existed {
		p.alertsByFile[sourcePath] = make([]Alert, 0)
		fileAlerts = p.alertsByFile[sourcePath]
	}
	fileAlerts = append(fileAlerts, alerts...)

	// Sort alerts by line
	sort.Slice(fileAlerts, func(i, j int) bool {
		return fileAlerts[i].GetSpecifier().GetTokens()[0].Line < fileAlerts[j].GetSpecifier().GetTokens()[0].Line
	})

	p.alertsByFile[sourcePath] = fileAlerts
}

func (p *Printer) PrintAlerts() error {
	for sourcePath, alerts := range p.alertsByFile {
		sourceFile, err := os.OpenFile(sourcePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// Initialize the Printer state
		p.source = bufio.NewScanner(sourceFile)
		if p.buffer != nil {
			for k := range p.buffer {
				delete(p.buffer, k)
			}
		}
		p.buffer = make(lineBuffer)
		p.line = 0

		alertMsg := strings.Builder{}
		for _, alert := range alerts {
			var msg string
			switch alert.GetAlertType() {
			case Error:
				msg = "[light_red][bold]error[%s]: [reset]"
			case Warning:
				msg = "[light_yellow][bold]warning[%s]: [default]"
			}

			alertMsg.WriteString(fmt.Sprintf(msg, alert.GetID()))
			p.writeMessage(&alertMsg, alert)
			p.writeLocation(&alertMsg, sourcePath, alert)
			err := p.writeCodeSnippet(&alertMsg, alert)
			p.writeNote(&alertMsg, alert)
			if err == nil {
				color.Printf(alertMsg.String() + "\n")
			} else {
				fmt.Printf("Fatal error: %s", err)
				break
			}

			alertMsg.Reset()
		}
	}

	return nil
}

func (p *Printer) writeLocation(alertMsg *strings.Builder, sourcePath string, alert Alert) {
	var largestLineNumber int
	tokens := alert.GetSpecifier().GetTokens()
	for _, token := range tokens {
		largestLineNumber = max(largestLineNumber, token.Line)
	}
	lineNumberSpaces := strings.Repeat(" ", len(strconv.Itoa(largestLineNumber)))
	locationStr := fmt.Sprintf("[light_gray]%s --- %s:%d:%d ---\n", lineNumberSpaces, sourcePath, tokens[0].Line, tokens[0].Column.Start)
	alertMsg.WriteString(locationStr)
}

func (p *Printer) writeMessage(alertMsg *strings.Builder, alert Alert) {
	messageStr := fmt.Sprintf("[bold]%s[reset]\n", alert.GetMessage())
	alertMsg.WriteString(messageStr)
}

func (p *Printer) writeNote(alertMsg *strings.Builder, alert Alert) {
	if alert.GetNote() != "" {
		location := alert.GetSpecifier().GetTokens()[0].Location
		lineNumberSpaces := strings.Repeat(" ", len(strconv.Itoa(location.Line)))
		noteStr := fmt.Sprintf("[cyan]%s = note:[default] %s\n", lineNumberSpaces, alert.GetNote())
		alertMsg.WriteString(noteStr)
		return
	}
}

func (p *Printer) writeCodeSnippet(alertMsg *strings.Builder, alert Alert) error {
	specifier := alert.GetSpecifier()
	location := mergeLocations(specifier.GetTokens())

	// The current line does not match the closest alert line
	// or is not already in the buffer, scanning is needed
	if p.line != location.Line && !p.buffer.isAvailable(location.Line) {
		// Scan the lines until we reach the start of the alert line location
		for p.line != location.Line {
			if !p.source.Scan() {
				err := p.source.Err()
				if err != nil {
					return err
				}
			}

			p.line++
		}

		// Failed to get the line (an error occurred)
		if p.line != location.Line {
			return fmt.Errorf("p.line != location.Line (%d, %d). Stopping", p.line, location.Line)
		}

		// Discard no longer necessary lines and append the necessary line
		if p.freeStart != -1 {
			p.buffer.freeLines(p.freeStart, p.line-1)
		}
		p.buffer.appendLine(p.line, p.source.Bytes())
	}

	switch specifier := specifier.(type) {
	case SingleLine:
		p.freeStart = p.line
		lines := p.buffer.retrieveLines(location.Line)
		alertMsg.WriteString(specifier.GetSnippet(lines))
	case MultiLine:
		tokens := specifier.GetTokens()
		startLocation, endLocation := tokens[0].Location, tokens[1].Location

		// The locations match on the same line, convert to a SingleLine snippet
		if startLocation.Line == endLocation.Line {
			newToken := tokens[0]
			newToken.Location = location
			lines := p.buffer.retrieveLines(location.Line)
			alertMsg.WriteString(NewSingle(newToken).GetSnippet(lines))
			break
		}

		// The current range of lines is already in the buffer, use it
		if p.buffer.isAvailable(startLocation.Line, endLocation.Line) {
			lines := p.buffer.retrieveLines(startLocation.Line, endLocation.Line)
			alertMsg.WriteString(specifier.GetSnippet(lines))
			break
		}

		// The range is missing some lines, scan them, append them and proceed as usual
		p.freeStart = p.line
		for p.line != endLocation.Line && !p.buffer.isAvailable(endLocation.Line) {
			if !p.source.Scan() {
				err := p.source.Err()
				if err != nil {
					return err
				}
			}

			p.line++
			p.buffer.appendLine(p.line, p.source.Bytes())
		}

		lines := p.buffer.retrieveLines(startLocation.Line, endLocation.Line)
		alertMsg.WriteString(specifier.GetSnippet(lines))
	}

	return nil
}

func mergeLocations(tokens []tokens.Token) (location tokens.Location) {
	location.Line = tokens[0].Line
	location.Column = tokens[0].Column

	for _, token := range tokens {
		location.Line = min(location.Line, token.Line)
		location.Column.Start = min(location.Column.Start, token.Column.Start)
		location.Column.End = max(location.Column.End, token.Column.End)
	}

	return
}
