package lsp

import (
	"hybroid/alerts"
	"hybroid/tokens"
	"testing"
)

// makeTokenAt builds a Token with a 1-based Location suitable for feeding
// into a SingleLine snippet. line, colStart, colEnd are all 1-based; the LSP
// conversion subtracts 1 when emitting positions.
func makeTokenAt(line, colStart, colEnd int) tokens.Token {
	return tokens.Token{
		Location: tokens.NewLocation(line, colStart, colEnd),
		Type:     tokens.String,
	}
}

// TestAlertsToDiagnostics_ErrorSeverity locks in that alerts of type
// alerts.Error map to LSP DiagnosticSeverity 1. If this regresses, editors
// will display lexer/parser errors as warnings (yellow squiggle instead of
// red) — a silent UX break.
func TestAlertsToDiagnostics_ErrorSeverity(t *testing.T) {
	tok := makeTokenAt(3, 5, 10)
	alert := &alerts.UnterminatedString{Specifier: alerts.NewSingle(tok)}

	diags := alertsToDiagnostics("file:///x.hyb", []alerts.Alert{alert})
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Severity != 1 {
		t.Errorf("Error alert should map to severity 1, got %d", diags[0].Severity)
	}
}

// TestAlertsToDiagnostics_WarningSeverity locks in that alerts of type
// alerts.Warning map to LSP DiagnosticSeverity 2. This is the case for the
// most common diagnostic in the wild: hyb073W "variable is not used".
func TestAlertsToDiagnostics_WarningSeverity(t *testing.T) {
	tok := makeTokenAt(7, 2, 4)
	alert := &alerts.UnusedElement{Specifier: alerts.NewSingle(tok)}

	diags := alertsToDiagnostics("file:///x.hyb", []alerts.Alert{alert})
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Severity != 2 {
		t.Errorf("Warning alert should map to severity 2, got %d", diags[0].Severity)
	}
}

// TestAlertsToDiagnostics_MalformedTokenLocation verifies that the
// conversion is defensive against pathological token locations. The
// production code path in alertsToDiagnostics falls back to {0,0} when the
// snippet has no tokens; here we exercise the equivalent path with a
// zero-width token. The contract is "no panic, no NaN, both ends equal".
func TestAlertsToDiagnostics_MalformedTokenLocation(t *testing.T) {
	// Token at line=0 col=0 col=0 — pathological but possible from a
	// hand-constructed alert in a test or a future refactor.
	tok := tokens.Token{
		Location: tokens.Location{Line: 0, Column: struct{ Start, End int }{0, 0}},
		Type:     tokens.String,
	}
	alert := &alerts.UnterminatedString{Specifier: alerts.NewSingle(tok)}

	diags := alertsToDiagnostics("file:///x.hyb", []alerts.Alert{alert})
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	d := diags[0]
	// 1-based line/col is converted to 0-based; both should be 0 (clamped).
	if d.Range.Start.Line != 0 || d.Range.Start.Character != 0 ||
		d.Range.End.Line != 0 || d.Range.End.Character != 0 {
		t.Errorf("expected range 0,0–0,0 for zero token, got %+v", d.Range)
	}
}

// TestAlertsToDiagnostics_NotePopulatesRelated verifies that an alert with
// a non-empty Note() produces RelatedInformation that points back at the
// originating URI. This is how editors render "see also" links in the
// gutter when hovering a diagnostic.
func TestAlertsToDiagnostics_NotePopulatesRelated(t *testing.T) {
	// Use a custom alert type to inject a Note(). Since Alert is an
	// interface, we implement it inline.
	tok := makeTokenAt(2, 1, 4)
	uri := DocumentURI("file:///example.hyb")
	alert := &notedAlert{
		id:      "hyb999T",
		typ:     alerts.Warning,
		msg:     "test message",
		note:    "see also here",
		snippet: alerts.NewSingle(tok),
	}

	diags := alertsToDiagnostics(uri, []alerts.Alert{alert})
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	d := diags[0]
	if len(d.RelatedInformation) != 1 {
		t.Fatalf("expected 1 related info, got %d", len(d.RelatedInformation))
	}
	if d.RelatedInformation[0].Location.URI != uri {
		t.Errorf("related URI = %q, want %q", d.RelatedInformation[0].Location.URI, uri)
	}
	if d.RelatedInformation[0].Message != "see also here" {
		t.Errorf("related message = %q, want %q",
			d.RelatedInformation[0].Message, "see also here")
	}
}

// TestAlertsToDiagnostics_MessageFormat locks in the "hybxxxX" prefix on
// every diagnostic message. The editor UI filters and de-duplicates
// diagnostics by this prefix; if someone refactors and drops the brackets,
// every existing user-configured warning filter breaks silently.
func TestAlertsToDiagnostics_MessageFormat(t *testing.T) {
	tok := makeTokenAt(1, 1, 1)
	alert := &alerts.UnusedElement{Specifier: alerts.NewSingle(tok)}

	diags := alertsToDiagnostics("file:///x.hyb", []alerts.Alert{alert})
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	want := "[hyb073W] " + alert.Message()
	if diags[0].Message != want {
		t.Errorf("message = %q, want %q", diags[0].Message, want)
	}
	// Also: the source must be "hybroid" so the editor can group by tool.
	if diags[0].Source == nil || *diags[0].Source != "hybroid" {
		t.Errorf("source = %v, want pointer to \"hybroid\"", diags[0].Source)
	}
}

// notedAlert is a test-only Alert implementation that lets us inject a
// non-empty Note() — none of the alerts in alerts/*.gen.go provide a way
// to do that without a generator round-trip.
type notedAlert struct {
	id      string
	typ     alerts.Type
	msg     string
	note    string
	snippet alerts.Snippet
}

func (n *notedAlert) ID() string                    { return n.id }
func (n *notedAlert) Message() string               { return n.msg }
func (n *notedAlert) Note() string                  { return n.note }
func (n *notedAlert) AlertType() alerts.Type        { return n.typ }
func (n *notedAlert) SnippetSpecifier() alerts.Snippet { return n.snippet }
