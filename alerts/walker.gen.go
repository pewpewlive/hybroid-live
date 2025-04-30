// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedParenthesiss struct {
  Specifier Snippet
  Symbol string
}

func (ep *ExpectedParenthesiss) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedParenthesiss) GetSpecifier() Snippet {
  return ep.Specifier
}

func (ep *ExpectedParenthesiss) GetNote() string {
  return ""
}

func (ep *ExpectedParenthesiss) GetID() string {
  return "hyb001W"
}

func (ep *ExpectedParenthesiss) GetAlertType() Type {
  return Error
}

