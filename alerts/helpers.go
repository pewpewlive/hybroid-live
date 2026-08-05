package alerts

import (
	"hybroid/core"
	"reflect"
	"strings"
	"testing"
)

func PrintAlerts[A any](t *testing.T, kind string, alertsToPrint ...A) {
	if len(alertsToPrint) == 100 {
		t.Logf("%s 100+ alert(s):", kind)
	} else {
		t.Logf("%s %d alert(s):", kind, len(alertsToPrint))
	}

	for i, alert := range alertsToPrint {
		if alert, ok := any(alert).(reflect.Type); ok {
			t.Logf("%d. %s", i+1, alert.Name())
			continue
		}

		actualAlert := any(alert).(Alert)
		span := actualAlert.Span()
		name := reflect.ValueOf(alert).Elem().Type().Name()
		msg := strings.TrimSpace(actualAlert.Message())
		t.Logf("%d. %s (%s) at line %d, column %d", i+1, name, msg, span.Line, span.Column)
	}
}

type Snippet interface {
	//GetSnippet(lines map[int][]byte, alert Alert) string
	Span() core.Span
}

type SnippetProvider struct {
	span core.Span
}

func (sp *SnippetProvider) Span() core.Span { return sp.span }

type ContextProvider struct {
	context string
}
