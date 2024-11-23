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

func (ep *ExpectedParenthesis) GetAlertStage() AlertStage {
  return Parser
}

