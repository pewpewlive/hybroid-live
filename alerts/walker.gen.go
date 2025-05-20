// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "strings"
)

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
  return "hyb001W"
}

func (ftie *ForbiddenTypeInEnvironment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
  Specifier Snippet
  GivenType string
}

func (iet *InvalidEnvironmentType) GetMessage() string {
  return fmt.Sprintf("'%s' is not a valid environment type", iet.GivenType)
}

func (iet *InvalidEnvironmentType) GetSpecifier() Snippet {
  return iet.Specifier
}

func (iet *InvalidEnvironmentType) GetNote() string {
  return "environment type can be 'Level', 'Mesh' or 'Sound'"
}

func (iet *InvalidEnvironmentType) GetID() string {
  return "hyb002W"
}

func (iet *InvalidEnvironmentType) GetAlertType() Type {
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

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateEnvironmentNames struct {
  Specifier Snippet
  Path1 string
  Path2 string
}

func (den *DuplicateEnvironmentNames) GetMessage() string {
  return fmt.Sprintf("duplicate environment names found between '%s' and '%s'", den.Path1, den.Path2)
}

func (den *DuplicateEnvironmentNames) GetSpecifier() Snippet {
  return den.Specifier
}

func (den *DuplicateEnvironmentNames) GetNote() string {
  return ""
}

func (den *DuplicateEnvironmentNames) GetID() string {
  return "hyb005W"
}

func (den *DuplicateEnvironmentNames) GetAlertType() Type {
  return Error
}

