// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedPaarenthesissss struct {
  Specifier Snippet
  Symbol string
}

func (ep *ExpectedPaarenthesissss) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedPaarenthesissss) GetSpecifier() Snippet {
  return ep.Specifier
}

func (ep *ExpectedPaarenthesissss) GetNote() string {
  return ""
}

func (ep *ExpectedPaarenthesissss) GetID() string {
  return "hyb001G"
}

func (ep *ExpectedPaarenthesissss) GetAlertType() Type {
  return Error
}

