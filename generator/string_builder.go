package generator

import (
	"hybroid/core"
	"strings"
)

type StringBuilder struct {
	strings.Builder
}

func (sb *StringBuilder) Write(chunks ...string) {
	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
}

func (sb *StringBuilder) WriteTabbed(chunks ...string) {
	sb.WriteString(getTabs())

	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
}

func (sb *StringBuilder) ReplaceSpan(str string, span core.Span[int]) {
	buffer := sb.String()
	sb.Reset()
	sb.Write(buffer[:span.Start], str, buffer[span.End:])
}
