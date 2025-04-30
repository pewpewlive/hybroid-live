// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultilineString struct {
  Specifier Snippet
}

func (ms *MultilineString) GetMessage() string {
  return "multiline strings are not allowed"
}

func (ms *MultilineString) GetSpecifier() Snippet {
  return ms.Specifier
}

func (ms *MultilineString) GetNote() string {
  return ""
}

func (ms *MultilineString) GetID() string {
  return "hyb001L"
}

func (ms *MultilineString) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnterminatedString struct {
  Specifier Snippet
}

func (us *UnterminatedString) GetMessage() string {
  return "unterminated string"
}

func (us *UnterminatedString) GetSpecifier() Snippet {
  return us.Specifier
}

func (us *UnterminatedString) GetNote() string {
  return ""
}

func (us *UnterminatedString) GetID() string {
  return "hyb002L"
}

func (us *UnterminatedString) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MalformedNumber struct {
  Specifier Snippet
  Number string
}

func (mn *MalformedNumber) GetMessage() string {
  return fmt.Sprintf("malformed number: '%s'", mn.Number)
}

func (mn *MalformedNumber) GetSpecifier() Snippet {
  return mn.Specifier
}

func (mn *MalformedNumber) GetNote() string {
  return ""
}

func (mn *MalformedNumber) GetID() string {
  return "hyb003L"
}

func (mn *MalformedNumber) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidNumberPostfix struct {
  Specifier Snippet
  Postfix string
}

func (inp *InvalidNumberPostfix) GetMessage() string {
  return fmt.Sprintf("invalid number postfix: '%s'", inp.Postfix)
}

func (inp *InvalidNumberPostfix) GetSpecifier() Snippet {
  return inp.Specifier
}

func (inp *InvalidNumberPostfix) GetNote() string {
  return "a valid postfix is either 'f', 'fx', 'r' or 'd'"
}

func (inp *InvalidNumberPostfix) GetID() string {
  return "hyb004L"
}

func (inp *InvalidNumberPostfix) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnsupportedCharacter struct {
  Specifier Snippet
  Character string
}

func (uc *UnsupportedCharacter) GetMessage() string {
  return fmt.Sprintf("unsupported character: '%s'", uc.Character)
}

func (uc *UnsupportedCharacter) GetSpecifier() Snippet {
  return uc.Specifier
}

func (uc *UnsupportedCharacter) GetNote() string {
  return ""
}

func (uc *UnsupportedCharacter) GetID() string {
  return "hyb005L"
}

func (uc *UnsupportedCharacter) GetAlertType() Type {
  return Error
}

