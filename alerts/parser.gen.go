// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedParenthesis struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Symbol string
}

func (ep *ExpectedParenthesis) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedParenthesis) GetTokens() []tokens.Token {
  return []tokens.Token{ep.Token}
}

func (ep *ExpectedParenthesis) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ep.Location}
}

func (ep *ExpectedParenthesis) GetNote() string {
  return ""
}

func (ep *ExpectedParenthesis) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (ee *ExpectedEnvironment) GetMessage() string {
  return "Expected environment statement"
}

func (ee *ExpectedEnvironment) GetTokens() []tokens.Token {
  return []tokens.Token{ee.Token}
}

func (ee *ExpectedEnvironment) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ee.Location}
}

func (ee *ExpectedEnvironment) GetNote() string {
  return "environment statement has to be the first statement in any hybroid file. example: env HellloWorld as Level"
}

func (ee *ExpectedEnvironment) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Symbol string
}

func (ei *ExpectedIdentifier) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ei.Symbol)
}

func (ei *ExpectedIdentifier) GetTokens() []tokens.Token {
  return []tokens.Token{ei.Token}
}

func (ei *ExpectedIdentifier) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ei.Location}
}

func (ei *ExpectedIdentifier) GetNote() string {
  return ""
}

func (ei *ExpectedIdentifier) GetAlertType() AlertType {
  return Error
}

