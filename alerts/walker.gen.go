// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedParenthesiss struct {
  Specifier Singleline
  Symbol string
}

func (ep *ExpectedParenthesiss) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedParenthesiss) GetSpecifier() SnippetSpecifier {
  return &ep.Specifier
}

func (ep *ExpectedParenthesiss) GetNote() string {
  return ""
}

func (ep *ExpectedParenthesiss) GetID() string {
  return "hyb001"
}

func (ep *ExpectedParenthesiss) GetAlertType() AlertType {
  return Error
}

