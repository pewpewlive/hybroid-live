// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedPaarenthesissss struct {
  Specifier SnippetSpecifier
  Symbol string
}

func (ep *ExpectedPaarenthesissss) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedPaarenthesissss) GetSpecifier() SnippetSpecifier {
  return ep.Specifier
}

func (ep *ExpectedPaarenthesissss) GetNote() string {
  return ""
}

func (ep *ExpectedPaarenthesissss) GetID() string {
  return "hyb001c"
}

func (ep *ExpectedPaarenthesissss) GetAlertType() AlertType {
  return Error
}

