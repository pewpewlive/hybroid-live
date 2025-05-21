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

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidAccessValue struct {
  Specifier Snippet
  ValueType string
}

func (iav *InvalidAccessValue) GetMessage() string {
  return fmt.Sprintf("value is of type '%s', so it cannot be accessed from", iav.ValueType)
}

func (iav *InvalidAccessValue) GetSpecifier() Snippet {
  return iav.Specifier
}

func (iav *InvalidAccessValue) GetNote() string {
  return "only lists, maps, classes, entities, structs and enums can be used to access values from"
}

func (iav *InvalidAccessValue) GetID() string {
  return "hyb006W"
}

func (iav *InvalidAccessValue) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type FieldAccessOnListOrMap struct {
  Specifier Snippet
  Field string
  AccessType string
}

func (faolom *FieldAccessOnListOrMap) GetMessage() string {
  return fmt.Sprintf("cannot access field '%s' from the %s", faolom.Field, faolom.AccessType)
}

func (faolom *FieldAccessOnListOrMap) GetSpecifier() Snippet {
  return faolom.Specifier
}

func (faolom *FieldAccessOnListOrMap) GetNote() string {
  return fmt.Sprintf("to access a value from a %s you use brackets, e.g. example[%s]", faolom.AccessType, faolom.Field)
}

func (faolom *FieldAccessOnListOrMap) GetID() string {
  return "hyb007W"
}

func (faolom *FieldAccessOnListOrMap) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MemberAccessOnNonListOrMap struct {
  Specifier Snippet
  Member string
  AccessType string
}

func (maonlom *MemberAccessOnNonListOrMap) GetMessage() string {
  return fmt.Sprintf("cannot access member '[%s]' from the %s", maonlom.Member, maonlom.AccessType)
}

func (maonlom *MemberAccessOnNonListOrMap) GetSpecifier() Snippet {
  return maonlom.Specifier
}

func (maonlom *MemberAccessOnNonListOrMap) GetNote() string {
  return "to access a value you use a dot and then an identifier, e.g. example.identifier"
}

func (maonlom *MemberAccessOnNonListOrMap) GetID() string {
  return "hyb008W"
}

func (maonlom *MemberAccessOnNonListOrMap) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidMemberIndex struct {
  Specifier Snippet
  AccessType string
  Index string
}

func (imi *InvalidMemberIndex) GetMessage() string {
  return fmt.Sprintf("'%s' is not of type number to be an index for the %s", imi.Index, imi.AccessType)
}

func (imi *InvalidMemberIndex) GetSpecifier() Snippet {
  return imi.Specifier
}

func (imi *InvalidMemberIndex) GetNote() string {
  return "for lists, an index (number) is used to access values, for maps, a key (text) is used"
}

func (imi *InvalidMemberIndex) GetID() string {
  return "hyb009W"
}

func (imi *InvalidMemberIndex) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidField struct {
  Specifier Snippet
  AccessType string
  FieldName string
}

func (_if *InvalidField) GetMessage() string {
  return fmt.Sprintf("field '%s' does not belong to the %s", _if.FieldName, _if.AccessType)
}

func (_if *InvalidField) GetSpecifier() Snippet {
  return _if.Specifier
}

func (_if *InvalidField) GetNote() string {
  return ""
}

func (_if *InvalidField) GetID() string {
  return "hyb010W"
}

func (_if *InvalidField) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MixedMapOrListContents struct {
  Specifier Snippet
  ContainerType string
  Type1 string
  Type2 string
}

func (mmolc *MixedMapOrListContents) GetMessage() string {
  return fmt.Sprintf("a %s's members must be the same type (found types: %s and %s)", mmolc.ContainerType, mmolc.Type1, mmolc.Type2)
}

func (mmolc *MixedMapOrListContents) GetSpecifier() Snippet {
  return mmolc.Specifier
}

func (mmolc *MixedMapOrListContents) GetNote() string {
  return ""
}

func (mmolc *MixedMapOrListContents) GetID() string {
  return "hyb011W"
}

func (mmolc *MixedMapOrListContents) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateKeyInMap struct {
  Specifier Snippet
}

func (dkim *DuplicateKeyInMap) GetMessage() string {
  return "found duplicate key in map"
}

func (dkim *DuplicateKeyInMap) GetSpecifier() Snippet {
  return dkim.Specifier
}

func (dkim *DuplicateKeyInMap) GetNote() string {
  return ""
}

func (dkim *DuplicateKeyInMap) GetID() string {
  return "hyb012W"
}

func (dkim *DuplicateKeyInMap) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCallerType struct {
  Specifier Snippet
  Type string
}

func (ict *InvalidCallerType) GetMessage() string {
  return fmt.Sprintf("caller is not a function (type found: %s)", ict.Type)
}

func (ict *InvalidCallerType) GetSpecifier() Snippet {
  return ict.Specifier
}

func (ict *InvalidCallerType) GetNote() string {
  return ""
}

func (ict *InvalidCallerType) GetID() string {
  return "hyb013W"
}

func (ict *InvalidCallerType) GetAlertType() Type {
  return Error
}

