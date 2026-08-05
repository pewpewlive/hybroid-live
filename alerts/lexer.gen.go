// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
	"hybroid/core"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultilineString struct {
	SnippetProvider
}

func NewMultilineString(span core.Span) *MultilineString {
	return &MultilineString{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ms *MultilineString) Message() string { return "multiline strings are not allowed" }
func (ms *MultilineString) Note() string    { return "" }
func (ms *MultilineString) Type() Type      { return Error }
func (ms *MultilineString) ID() string      { return "hyb001L" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnterminatedString struct {
	SnippetProvider
}

func NewUnterminatedString(span core.Span) *UnterminatedString {
	return &UnterminatedString{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (us *UnterminatedString) Message() string { return "unterminated string" }
func (us *UnterminatedString) Note() string    { return "" }
func (us *UnterminatedString) Type() Type      { return Error }
func (us *UnterminatedString) ID() string      { return "hyb002L" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MalformedNumber struct {
	SnippetProvider
	number string
}

func NewMalformedNumber(span core.Span, number string) *MalformedNumber {
	return &MalformedNumber{
		SnippetProvider: SnippetProvider{span: span},
		number:          number,
	}
}
func (mn *MalformedNumber) Message() string { return fmt.Sprintf("malformed number: '%s'", mn.number) }
func (mn *MalformedNumber) Note() string    { return "" }
func (mn *MalformedNumber) Type() Type      { return Error }
func (mn *MalformedNumber) ID() string      { return "hyb003L" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidDigitInLiteral struct {
	SnippetProvider
	digit   string
	literal string
}

func NewInvalidDigitInLiteral(span core.Span, digit string, literal string) *InvalidDigitInLiteral {
	return &InvalidDigitInLiteral{
		SnippetProvider: SnippetProvider{span: span},
		digit:           digit,
		literal:         literal,
	}
}
func (idil *InvalidDigitInLiteral) Message() string {
	return fmt.Sprintf("invalid digit '%s' in %s literal", idil.digit, idil.literal)
}
func (idil *InvalidDigitInLiteral) Note() string { return "" }
func (idil *InvalidDigitInLiteral) Type() Type   { return Error }
func (idil *InvalidDigitInLiteral) ID() string   { return "hyb004L" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidNumberPostfix struct {
	SnippetProvider
	postfix string
}

func NewInvalidNumberPostfix(span core.Span, postfix string) *InvalidNumberPostfix {
	return &InvalidNumberPostfix{
		SnippetProvider: SnippetProvider{span: span},
		postfix:         postfix,
	}
}
func (inp *InvalidNumberPostfix) Message() string {
	return fmt.Sprintf("invalid number postfix: '%s'", inp.postfix)
}
func (inp *InvalidNumberPostfix) Note() string {
	return "a valid postfix is either 'f', 'fx', 'r' or 'd'"
}
func (inp *InvalidNumberPostfix) Type() Type { return Error }
func (inp *InvalidNumberPostfix) ID() string { return "hyb005L" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnsupportedCharacter struct {
	SnippetProvider
	character string
}

func NewUnsupportedCharacter(span core.Span, character string) *UnsupportedCharacter {
	return &UnsupportedCharacter{
		SnippetProvider: SnippetProvider{span: span},
		character:       character,
	}
}
func (uc *UnsupportedCharacter) Message() string {
	return fmt.Sprintf("unsupported character: '%s'", uc.character)
}
func (uc *UnsupportedCharacter) Note() string { return "" }
func (uc *UnsupportedCharacter) Type() Type   { return Error }
func (uc *UnsupportedCharacter) ID() string   { return "hyb006L" }
