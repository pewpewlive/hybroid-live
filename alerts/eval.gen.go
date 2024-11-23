// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedPaarenthesissss struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Symbol string
}

func (ep *ExpectedPaarenthesissss) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedPaarenthesissss) GetTokens() []tokens.Token {
  return []tokens.Token{ep.Token}
}

func (ep *ExpectedPaarenthesissss) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ep.Location}
}

func (ep *ExpectedPaarenthesissss) GetNote() string {
  return ""
}

func (ep *ExpectedPaarenthesissss) GetName() string {
  return "ExpectedPaarenthesissss"
}

func (ep *ExpectedPaarenthesissss) GetAlertType() AlertType {
  return Error
}

