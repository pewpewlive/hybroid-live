package generator

import (
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
