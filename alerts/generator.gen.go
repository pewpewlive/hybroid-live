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

func (ep *ExpectedPaarenthesissss) Message() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedPaarenthesissss) SnippetSpecifier() Snippet {
  return ep.Specifier
}

func (ep *ExpectedPaarenthesissss) Note() string {
  return ""
}

func (ep *ExpectedPaarenthesissss) ID() string {
  return "hyb001G"
}

func (ep *ExpectedPaarenthesissss) AlertType() Type {
  return Error
}

