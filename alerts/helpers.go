package alerts

import (
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
		token := actualAlert.SnippetSpecifier().GetTokens()[0]
		loc := token.Location
		name := reflect.ValueOf(alert).Elem().Type().Name()
		msg := strings.TrimSpace(actualAlert.Message())
		t.Logf("%d. %s (%s) at line %d, column %d on token '%s'", i+1, name, msg, loc.Line, loc.Column.Start, token.Lexeme)
	}
}
