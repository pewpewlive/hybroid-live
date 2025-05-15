// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "strings"
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

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForbiddenTypeInEnvironment struct {
  Specifier Snippet
  Type string
  Envs []string
}

func (ftie *ForbiddenTypeInEnvironment) GetMessage() string {
  return fmt.Sprintf("cannot have a %s in the following environments: %s", ftie.Type, strings.Join(ftie.Envs, ", "))
}

func (ftie *ForbiddenTypeInEnvironment) GetSpecifier() Snippet {
  return ftie.Specifier
}

func (ftie *ForbiddenTypeInEnvironment) GetNote() string {
  return ""
}

func (ftie *ForbiddenTypeInEnvironment) GetID() string {
  return "hyb002W"
}

func (ftie *ForbiddenTypeInEnvironment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentRedaclaration struct {
  Specifier Snippet
}

func (er *EnvironmentRedaclaration) GetMessage() string {
  return "cannot redeclare an environment"
}

func (er *EnvironmentRedaclaration) GetSpecifier() Snippet {
  return er.Specifier
}

func (er *EnvironmentRedaclaration) GetNote() string {
  return ""
}

func (er *EnvironmentRedaclaration) GetID() string {
  return "hyb003W"
}

func (er *EnvironmentRedaclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
  Specifier Snippet
}

func (ee *ExpectedEnvironment) GetMessage() string {
  return "expected environment declaration"
}

func (ee *ExpectedEnvironment) GetSpecifier() Snippet {
  return ee.Specifier
}

func (ee *ExpectedEnvironment) GetNote() string {
  return "the first declaration in any Hybroid file has to be an environment declaration"
}

func (ee *ExpectedEnvironment) GetID() string {
  return "hyb004W"
}

func (ee *ExpectedEnvironment) GetAlertType() Type {
  return Error
}

