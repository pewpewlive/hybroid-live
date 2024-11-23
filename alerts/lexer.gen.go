// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultilineString struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (ms *MultilineString) GetMessage() string {
  return "multiline strings are not allowed"
}

func (ms *MultilineString) GetTokens() []tokens.Token {
  return []tokens.Token{ms.Token}
}

func (ms *MultilineString) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ms.Location}
}

func (ms *MultilineString) GetNote() string {
  return ""
}

func (ms *MultilineString) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnterminatedString struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (us *UnterminatedString) GetMessage() string {
  return "unterminated string"
}

func (us *UnterminatedString) GetTokens() []tokens.Token {
  return []tokens.Token{us.Token}
}

func (us *UnterminatedString) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{us.Location}
}

func (us *UnterminatedString) GetNote() string {
  return ""
}

func (us *UnterminatedString) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MalformedNumber struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Number string
}

func (mn *MalformedNumber) GetMessage() string {
  return fmt.Sprintf("malformed number: '%s'", mn.Number)
}

func (mn *MalformedNumber) GetTokens() []tokens.Token {
  return []tokens.Token{mn.Token}
}

func (mn *MalformedNumber) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{mn.Location}
}

func (mn *MalformedNumber) GetNote() string {
  return ""
}

func (mn *MalformedNumber) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidNumberPostfix struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Postfix string
}

func (inp *InvalidNumberPostfix) GetMessage() string {
  return fmt.Sprintf("invalid number postfix: '%s'", inp.Postfix)
}

func (inp *InvalidNumberPostfix) GetTokens() []tokens.Token {
  return []tokens.Token{inp.Token}
}

func (inp *InvalidNumberPostfix) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{inp.Location}
}

func (inp *InvalidNumberPostfix) GetNote() string {
  return "valid number postfixes: 1f, 1fx, 1r, 1d"
}

func (inp *InvalidNumberPostfix) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnsupportedCharacter struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Character string
}

func (uc *UnsupportedCharacter) GetMessage() string {
  return fmt.Sprintf("unsupported character: '%s'", uc.Character)
}

func (uc *UnsupportedCharacter) GetTokens() []tokens.Token {
  return []tokens.Token{uc.Token}
}

func (uc *UnsupportedCharacter) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{uc.Location}
}

func (uc *UnsupportedCharacter) GetNote() string {
  return ""
}

func (uc *UnsupportedCharacter) GetAlertType() AlertType {
  return Error
}

