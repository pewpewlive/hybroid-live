// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultilineString struct {
  Specifier SnippetSpecifier
}

func (ms *MultilineString) GetMessage() string {
  return "multiline strings are not allowed"
}

func (ms *MultilineString) GetSpecifier() SnippetSpecifier {
  return ms.Specifier
}

func (ms *MultilineString) GetNote() string {
  return ""
}

func (ms *MultilineString) GetID() string {
  return "hyb001"
}

func (ms *MultilineString) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnterminatedString struct {
  Specifier SnippetSpecifier
}

func (us *UnterminatedString) GetMessage() string {
  return "unterminated string"
}

func (us *UnterminatedString) GetSpecifier() SnippetSpecifier {
  return us.Specifier
}

func (us *UnterminatedString) GetNote() string {
  return ""
}

func (us *UnterminatedString) GetID() string {
  return "hyb002"
}

func (us *UnterminatedString) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MalformedNumber struct {
  Specifier SnippetSpecifier
  Number string
}

func (mn *MalformedNumber) GetMessage() string {
  return fmt.Sprintf("malformed number: '%s'", mn.Number)
}

func (mn *MalformedNumber) GetSpecifier() SnippetSpecifier {
  return mn.Specifier
}

func (mn *MalformedNumber) GetNote() string {
  return ""
}

func (mn *MalformedNumber) GetID() string {
  return "hyb003"
}

func (mn *MalformedNumber) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidNumberPostfix struct {
  Specifier SnippetSpecifier
  Postfix string
}

func (inp *InvalidNumberPostfix) GetMessage() string {
  return fmt.Sprintf("invalid number postfix: '%s'", inp.Postfix)
}

func (inp *InvalidNumberPostfix) GetSpecifier() SnippetSpecifier {
  return inp.Specifier
}

func (inp *InvalidNumberPostfix) GetNote() string {
  return "a valid postfix is either 'f', 'fx', 'r' or 'd'"
}

func (inp *InvalidNumberPostfix) GetID() string {
  return "hyb004"
}

func (inp *InvalidNumberPostfix) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnsupportedCharacter struct {
  Specifier SnippetSpecifier
  Character string
}

func (uc *UnsupportedCharacter) GetMessage() string {
  return fmt.Sprintf("unsupported character: '%s'", uc.Character)
}

func (uc *UnsupportedCharacter) GetSpecifier() SnippetSpecifier {
  return uc.Specifier
}

func (uc *UnsupportedCharacter) GetNote() string {
  return ""
}

func (uc *UnsupportedCharacter) GetID() string {
  return "hyb005"
}

func (uc *UnsupportedCharacter) GetAlertType() AlertType {
  return Error
}

