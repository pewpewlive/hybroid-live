// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultilineString struct {
  Specifier Snippet
}

func (ms *MultilineString) Message() string {
  return "multiline strings are not allowed"
}

func (ms *MultilineString) SnippetSpecifier() Snippet {
  return ms.Specifier
}

func (ms *MultilineString) Note() string {
  return ""
}

func (ms *MultilineString) ID() string {
  return "hyb001L"
}

func (ms *MultilineString) AlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnterminatedString struct {
  Specifier Snippet
}

func (us *UnterminatedString) Message() string {
  return "unterminated string"
}

func (us *UnterminatedString) SnippetSpecifier() Snippet {
  return us.Specifier
}

func (us *UnterminatedString) Note() string {
  return ""
}

func (us *UnterminatedString) ID() string {
  return "hyb002L"
}

func (us *UnterminatedString) AlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MalformedNumber struct {
  Specifier Snippet
  Number string
}

func (mn *MalformedNumber) Message() string {
  return fmt.Sprintf("malformed number: '%s'", mn.Number)
}

func (mn *MalformedNumber) SnippetSpecifier() Snippet {
  return mn.Specifier
}

func (mn *MalformedNumber) Note() string {
  return ""
}

func (mn *MalformedNumber) ID() string {
  return "hyb003L"
}

func (mn *MalformedNumber) AlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidDigitInLiteral struct {
  Specifier Snippet
  Digit string
  Literal string
}

func (idil *InvalidDigitInLiteral) Message() string {
  return fmt.Sprintf("invalid digit '%s' in %s literal", idil.Digit, idil.Literal)
}

func (idil *InvalidDigitInLiteral) SnippetSpecifier() Snippet {
  return idil.Specifier
}

func (idil *InvalidDigitInLiteral) Note() string {
  return ""
}

func (idil *InvalidDigitInLiteral) ID() string {
  return "hyb004L"
}

func (idil *InvalidDigitInLiteral) AlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidNumberPostfix struct {
  Specifier Snippet
  Postfix string
}

func (inp *InvalidNumberPostfix) Message() string {
  return fmt.Sprintf("invalid number postfix: '%s'", inp.Postfix)
}

func (inp *InvalidNumberPostfix) SnippetSpecifier() Snippet {
  return inp.Specifier
}

func (inp *InvalidNumberPostfix) Note() string {
  return "a valid postfix is either 'f', 'fx', 'r' or 'd'"
}

func (inp *InvalidNumberPostfix) ID() string {
  return "hyb005L"
}

func (inp *InvalidNumberPostfix) AlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnsupportedCharacter struct {
  Specifier Snippet
  Character string
}

func (uc *UnsupportedCharacter) Message() string {
  return fmt.Sprintf("unsupported character: '%s'", uc.Character)
}

func (uc *UnsupportedCharacter) SnippetSpecifier() Snippet {
  return uc.Specifier
}

func (uc *UnsupportedCharacter) Note() string {
  return ""
}

func (uc *UnsupportedCharacter) ID() string {
  return "hyb006L"
}

func (uc *UnsupportedCharacter) AlertType() Type {
  return Error
}

