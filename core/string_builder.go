package core

import (
	"fmt"
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

func (sb *StringBuilder) Writef(format string, args ...any) {
	sb.WriteString(fmt.Sprintf(format, args...))
}
