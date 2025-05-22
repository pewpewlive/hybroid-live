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
  return fmt.Sprintf("field '%s' does not belong to the '%s'", _if.FieldName, _if.AccessType)
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
  return fmt.Sprintf("a %s's members must be the same type (found types: '%s' and '%s')", mmolc.ContainerType, mmolc.Type1, mmolc.Type2)
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
  return fmt.Sprintf("caller is not a function (type found: '%s')", ict.Type)
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

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MethodOrFieldNotFound struct {
  Specifier Snippet
  Name string
}

func (mofnf *MethodOrFieldNotFound) GetMessage() string {
  return fmt.Sprintf("no method or field named '%s'", mofnf.Name)
}

func (mofnf *MethodOrFieldNotFound) GetSpecifier() Snippet {
  return mofnf.Specifier
}

func (mofnf *MethodOrFieldNotFound) GetNote() string {
  return ""
}

func (mofnf *MethodOrFieldNotFound) GetID() string {
  return "hyb014W"
}

func (mofnf *MethodOrFieldNotFound) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForeignLocalVariableAccess struct {
  Specifier Snippet
  Name string
}

func (flva *ForeignLocalVariableAccess) GetMessage() string {
  return fmt.Sprintf("cannot access local variable '%s' belonging to a different environment", flva.Name)
}

func (flva *ForeignLocalVariableAccess) GetSpecifier() Snippet {
  return flva.Specifier
}

func (flva *ForeignLocalVariableAccess) GetNote() string {
  return ""
}

func (flva *ForeignLocalVariableAccess) GetID() string {
  return "hyb015W"
}

func (flva *ForeignLocalVariableAccess) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidArgumentType struct {
  Specifier Snippet
  GivenType string
  ExpectedType string
}

func (iat *InvalidArgumentType) GetMessage() string {
  return fmt.Sprintf("argument was of type %s, but should be %s", iat.GivenType, iat.ExpectedType)
}

func (iat *InvalidArgumentType) GetSpecifier() Snippet {
  return iat.Specifier
}

func (iat *InvalidArgumentType) GetNote() string {
  return ""
}

func (iat *InvalidArgumentType) GetID() string {
  return "hyb016W"
}

func (iat *InvalidArgumentType) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type PublicDeclarationInLocalScope struct {
  Specifier Snippet
}

func (pdils *PublicDeclarationInLocalScope) GetMessage() string {
  return "cannot have a public declaration that is in a local scope"
}

func (pdils *PublicDeclarationInLocalScope) GetSpecifier() Snippet {
  return pdils.Specifier
}

func (pdils *PublicDeclarationInLocalScope) GetNote() string {
  return ""
}

func (pdils *PublicDeclarationInLocalScope) GetID() string {
  return "hyb017W"
}

func (pdils *PublicDeclarationInLocalScope) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type Redeclaration struct {
  Specifier Snippet
  VarName string
  DeclType string
}

func (r *Redeclaration) GetMessage() string {
  return fmt.Sprintf("a %s named '%s' already exists", r.DeclType, r.VarName)
}

func (r *Redeclaration) GetSpecifier() Snippet {
  return r.Specifier
}

func (r *Redeclaration) GetNote() string {
  return ""
}

func (r *Redeclaration) GetID() string {
  return "hyb018W"
}

func (r *Redeclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnnecessaryTypeInConstDeclaration struct {
  Specifier Snippet
}

func (uticd *UnnecessaryTypeInConstDeclaration) GetMessage() string {
  return "an explicit type is not necessary for a const declaration"
}

func (uticd *UnnecessaryTypeInConstDeclaration) GetSpecifier() Snippet {
  return uticd.Specifier
}

func (uticd *UnnecessaryTypeInConstDeclaration) GetNote() string {
  return ""
}

func (uticd *UnnecessaryTypeInConstDeclaration) GetID() string {
  return "hyb019W"
}

func (uticd *UnnecessaryTypeInConstDeclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type NoValueGivenForConstant struct {
  Specifier Snippet
}

func (nvgfc *NoValueGivenForConstant) GetMessage() string {
  return "constant must be declared with a value"
}

func (nvgfc *NoValueGivenForConstant) GetSpecifier() Snippet {
  return nvgfc.Specifier
}

func (nvgfc *NoValueGivenForConstant) GetNote() string {
  return ""
}

func (nvgfc *NoValueGivenForConstant) GetID() string {
  return "hyb020W"
}

func (nvgfc *NoValueGivenForConstant) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TooFewValuesInDeclaration struct {
  Specifier Snippet
  RequiredAmount int
}

func (tfvid *TooFewValuesInDeclaration) GetMessage() string {
  return fmt.Sprintf("%d more value(s) required in variable declaration", tfvid.RequiredAmount)
}

func (tfvid *TooFewValuesInDeclaration) GetSpecifier() Snippet {
  return tfvid.Specifier
}

func (tfvid *TooFewValuesInDeclaration) GetNote() string {
  return ""
}

func (tfvid *TooFewValuesInDeclaration) GetID() string {
  return "hyb021W"
}

func (tfvid *TooFewValuesInDeclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeRequiredInDeclaration struct {
  Specifier Snippet
}

func (etrid *ExplicitTypeRequiredInDeclaration) GetMessage() string {
  return "a variable declared without a value requires an explicit type"
}

func (etrid *ExplicitTypeRequiredInDeclaration) GetSpecifier() Snippet {
  return etrid.Specifier
}

func (etrid *ExplicitTypeRequiredInDeclaration) GetNote() string {
  return ""
}

func (etrid *ExplicitTypeRequiredInDeclaration) GetID() string {
  return "hyb022W"
}

func (etrid *ExplicitTypeRequiredInDeclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeMismatch struct {
  Specifier Snippet
  ExplicitType string
  ValueType string
}

func (etm *ExplicitTypeMismatch) GetMessage() string {
  return fmt.Sprintf("variable was given explicit type '%s', but its value is a '%s'", etm.ExplicitType, etm.ValueType)
}

func (etm *ExplicitTypeMismatch) GetSpecifier() Snippet {
  return etm.Specifier
}

func (etm *ExplicitTypeMismatch) GetNote() string {
  return ""
}

func (etm *ExplicitTypeMismatch) GetID() string {
  return "hyb023W"
}

func (etm *ExplicitTypeMismatch) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeNotAllowed struct {
  Specifier Snippet
  ExplicitType string
}

func (etna *ExplicitTypeNotAllowed) GetMessage() string {
  return fmt.Sprintf("cannot create a default value from the explicit type '%s'", etna.ExplicitType)
}

func (etna *ExplicitTypeNotAllowed) GetSpecifier() Snippet {
  return etna.Specifier
}

func (etna *ExplicitTypeNotAllowed) GetNote() string {
  return "some types don't have default values, like entities and classes"
}

func (etna *ExplicitTypeNotAllowed) GetID() string {
  return "hyb024W"
}

func (etna *ExplicitTypeNotAllowed) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TooManyValuesInDeclaration struct {
  Specifier Snippet
  ExtraAmount int
}

func (tmvid *TooManyValuesInDeclaration) GetMessage() string {
  return fmt.Sprintf("%d less value(s) required in variable declaration", tmvid.ExtraAmount)
}

func (tmvid *TooManyValuesInDeclaration) GetSpecifier() Snippet {
  return tmvid.Specifier
}

func (tmvid *TooManyValuesInDeclaration) GetNote() string {
  return ""
}

func (tmvid *TooManyValuesInDeclaration) GetID() string {
  return "hyb025W"
}

func (tmvid *TooManyValuesInDeclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ImportCycle struct {
  Specifier Snippet
  HybPath1 string
  HybPath2 string
}

func (ic *ImportCycle) GetMessage() string {
  return fmt.Sprintf("", ic.HybPath1, ic.HybPath1)
}

func (ic *ImportCycle) GetSpecifier() Snippet {
  return ic.Specifier
}

func (ic *ImportCycle) GetNote() string {
  return ""
}

func (ic *ImportCycle) GetID() string {
  return "hyb026W"
}

func (ic *ImportCycle) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UndeclaredVariableAccess struct {
  Specifier Snippet
  Var string
}

func (uva *UndeclaredVariableAccess) GetMessage() string {
  return fmt.Sprintf("'%s' is not a declared variable", uva.Var)
}

func (uva *UndeclaredVariableAccess) GetSpecifier() Snippet {
  return uva.Specifier
}

func (uva *UndeclaredVariableAccess) GetNote() string {
  return ""
}

func (uva *UndeclaredVariableAccess) GetID() string {
  return "hyb027W"
}

func (uva *UndeclaredVariableAccess) GetAlertType() Type {
  return Error
}

