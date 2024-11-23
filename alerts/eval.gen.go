// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedPaarenthesis struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Symbol string
}

func (ep *ExpectedPaarenthesis) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedPaarenthesis) GetTokens() []tokens.Token {
  return []tokens.Token{ep.Token}
}

func (ep *ExpectedPaarenthesis) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ep.Location}
}

func (ep *ExpectedPaarenthesis) GetNote() string {
  return ""
}

func (ep *ExpectedPaarenthesis) GetAlertType() AlertType {
  return Error
}

func (ep *ExpectedPaarenthesis) GetAlertStage() AlertStage {
  return Eval
}

