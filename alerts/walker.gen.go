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
  Type string
}

func (iet *InvalidEnvironmentType) GetMessage() string {
  return fmt.Sprintf("'%s' is not a valid environment type", iet.Type)
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
  Type string
}

func (iav *InvalidAccessValue) GetMessage() string {
  return fmt.Sprintf("value is of type '%s', so it cannot be accessed from", iav.Type)
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
  return fmt.Sprintf("cannot call value of of type '%s' as a function", ict.Type)
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
type TooFewValuesGiven struct {
  Specifier Snippet
  RequiredAmount int
  Context string
}

func (tfvg *TooFewValuesGiven) GetMessage() string {
  return fmt.Sprintf("%d more value(s) required %s", tfvg.RequiredAmount, tfvg.Context)
}

func (tfvg *TooFewValuesGiven) GetSpecifier() Snippet {
  return tfvg.Specifier
}

func (tfvg *TooFewValuesGiven) GetNote() string {
  return ""
}

func (tfvg *TooFewValuesGiven) GetID() string {
  return "hyb021W"
}

func (tfvg *TooFewValuesGiven) GetAlertType() Type {
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
type TooManyValuesGiven struct {
  Specifier Snippet
  ExtraAmount int
  Context string
}

func (tmvg *TooManyValuesGiven) GetMessage() string {
  return fmt.Sprintf("%d less value(s) required %s", tmvg.ExtraAmount, tmvg.Context)
}

func (tmvg *TooManyValuesGiven) GetSpecifier() Snippet {
  return tmvg.Specifier
}

func (tmvg *TooManyValuesGiven) GetNote() string {
  return ""
}

func (tmvg *TooManyValuesGiven) GetID() string {
  return "hyb025W"
}

func (tmvg *TooManyValuesGiven) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ImportCycle struct {
  Specifier Snippet
  HybPath1 string
  HybPath2 string
}

func (ic *ImportCycle) GetMessage() string {
  return fmt.Sprintf("import cycle detected: cycling paths: '%s' and '%s'", ic.HybPath1, ic.HybPath2)
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

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ConstValueAssignment struct {
  Specifier Snippet
}

func (cva *ConstValueAssignment) GetMessage() string {
  return "cannot modify a constant value"
}

func (cva *ConstValueAssignment) GetSpecifier() Snippet {
  return cva.Specifier
}

func (cva *ConstValueAssignment) GetNote() string {
  return ""
}

func (cva *ConstValueAssignment) GetID() string {
  return "hyb028W"
}

func (cva *ConstValueAssignment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type AssignmentTypeMismatch struct {
  Specifier Snippet
  VarType string
  ValType string
}

func (atm *AssignmentTypeMismatch) GetMessage() string {
  return fmt.Sprintf("variable is of type '%s', but a value of '%s' was assigned to it", atm.VarType, atm.ValType)
}

func (atm *AssignmentTypeMismatch) GetSpecifier() Snippet {
  return atm.Specifier
}

func (atm *AssignmentTypeMismatch) GetNote() string {
  return ""
}

func (atm *AssignmentTypeMismatch) GetID() string {
  return "hyb029W"
}

func (atm *AssignmentTypeMismatch) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidTypeInCompoundAssignment struct {
  Specifier Snippet
  Type string
}

func (itica *InvalidTypeInCompoundAssignment) GetMessage() string {
  return fmt.Sprintf("the type '%s' is not allowed in compound assignment", itica.Type)
}

func (itica *InvalidTypeInCompoundAssignment) GetSpecifier() Snippet {
  return itica.Specifier
}

func (itica *InvalidTypeInCompoundAssignment) GetNote() string {
  return "only numerical types are allowed, like numbers, fixeds, fixedpoints, degrees and radians"
}

func (itica *InvalidTypeInCompoundAssignment) GetID() string {
  return "hyb030W"
}

func (itica *InvalidTypeInCompoundAssignment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidUseOfSelf struct {
  Specifier Snippet
}

func (iuos *InvalidUseOfSelf) GetMessage() string {
  return "cannot use self outside of class or entity"
}

func (iuos *InvalidUseOfSelf) GetSpecifier() Snippet {
  return iuos.Specifier
}

func (iuos *InvalidUseOfSelf) GetNote() string {
  return ""
}

func (iuos *InvalidUseOfSelf) GetID() string {
  return "hyb031W"
}

func (iuos *InvalidUseOfSelf) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateGenericParameter struct {
  Specifier Snippet
  Name string
}

func (dgp *DuplicateGenericParameter) GetMessage() string {
  return "the generic parameter '%s' is given more than once"
}

func (dgp *DuplicateGenericParameter) GetSpecifier() Snippet {
  return dgp.Specifier
}

func (dgp *DuplicateGenericParameter) GetNote() string {
  return ""
}

func (dgp *DuplicateGenericParameter) GetID() string {
  return "hyb032W"
}

func (dgp *DuplicateGenericParameter) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnreachableCode struct {
  Specifier Snippet
}

func (uc *UnreachableCode) GetMessage() string {
  return "unreachable code detected"
}

func (uc *UnreachableCode) GetSpecifier() Snippet {
  return uc.Specifier
}

func (uc *UnreachableCode) GetNote() string {
  return ""
}

func (uc *UnreachableCode) GetID() string {
  return "hyb033W"
}

func (uc *UnreachableCode) GetAlertType() Type {
  return Warning
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidUseOfExitStmt struct {
  Specifier Snippet
  ExitNode string
  Context string
}

func (iuoes *InvalidUseOfExitStmt) GetMessage() string {
  return fmt.Sprintf("cannot use '%s' outside of %s", iuoes.ExitNode, iuoes.Context)
}

func (iuoes *InvalidUseOfExitStmt) GetSpecifier() Snippet {
  return iuoes.Specifier
}

func (iuoes *InvalidUseOfExitStmt) GetNote() string {
  return ""
}

func (iuoes *InvalidUseOfExitStmt) GetID() string {
  return "hyb034W"
}

func (iuoes *InvalidUseOfExitStmt) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypeMismatch struct {
  Specifier Snippet
  Type1 string
  Type2 string
  Context string
}

func (tm *TypeMismatch) GetMessage() string {
  return fmt.Sprintf("expected '%s', got '%s' %s", tm.Type1, tm.Type2, tm.Context)
}

func (tm *TypeMismatch) GetSpecifier() Snippet {
  return tm.Specifier
}

func (tm *TypeMismatch) GetNote() string {
  return ""
}

func (tm *TypeMismatch) GetID() string {
  return "hyb035W"
}

func (tm *TypeMismatch) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UseStmtInLocalBlock struct {
  Specifier Snippet
}

func (usilb *UseStmtInLocalBlock) GetMessage() string {
  return "use statements must be in the global scope"
}

func (usilb *UseStmtInLocalBlock) GetSpecifier() Snippet {
  return usilb.Specifier
}

func (usilb *UseStmtInLocalBlock) GetNote() string {
  return ""
}

func (usilb *UseStmtInLocalBlock) GetID() string {
  return "hyb036W"
}

func (usilb *UseStmtInLocalBlock) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UsedPewpewInNonLevelEnvironment struct {
  Specifier Snippet
}

func (upinle *UsedPewpewInNonLevelEnvironment) GetMessage() string {
  return "cannot use the Pewpew environment in a Mesh or Sound type environment"
}

func (upinle *UsedPewpewInNonLevelEnvironment) GetSpecifier() Snippet {
  return upinle.Specifier
}

func (upinle *UsedPewpewInNonLevelEnvironment) GetNote() string {
  return ""
}

func (upinle *UsedPewpewInNonLevelEnvironment) GetID() string {
  return "hyb037W"
}

func (upinle *UsedPewpewInNonLevelEnvironment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnallowedLibraryUse struct {
  Specifier Snippet
  Library string
  UnallowedEnvs string
}

func (ulu *UnallowedLibraryUse) GetMessage() string {
  return fmt.Sprintf("cannot use the %s library in a %s environment", ulu.Library, ulu.UnallowedEnvs)
}

func (ulu *UnallowedLibraryUse) GetSpecifier() Snippet {
  return ulu.Specifier
}

func (ulu *UnallowedLibraryUse) GetNote() string {
  return ""
}

func (ulu *UnallowedLibraryUse) GetID() string {
  return "hyb038W"
}

func (ulu *UnallowedLibraryUse) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentAccess struct {
  Specifier Snippet
  EnvName string
}

func (iea *InvalidEnvironmentAccess) GetMessage() string {
  return fmt.Sprintf("environment named '%s' does not exist", iea.EnvName)
}

func (iea *InvalidEnvironmentAccess) GetSpecifier() Snippet {
  return iea.Specifier
}

func (iea *InvalidEnvironmentAccess) GetNote() string {
  return ""
}

func (iea *InvalidEnvironmentAccess) GetID() string {
  return "hyb039W"
}

func (iea *InvalidEnvironmentAccess) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentReuse struct {
  Specifier Snippet
  EnvName string
}

func (er *EnvironmentReuse) GetMessage() string {
  return fmt.Sprintf("environment named '%s' is already imported through use statement", er.EnvName)
}

func (er *EnvironmentReuse) GetSpecifier() Snippet {
  return er.Specifier
}

func (er *EnvironmentReuse) GetNote() string {
  return ""
}

func (er *EnvironmentReuse) GetID() string {
  return "hyb040W"
}

func (er *EnvironmentReuse) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidIteratorType struct {
  Specifier Snippet
  Type string
}

func (iit *InvalidIteratorType) GetMessage() string {
  return fmt.Sprintf("a for loop iterator must be a map or a list (found: '%s')", iit.Type)
}

func (iit *InvalidIteratorType) GetSpecifier() Snippet {
  return iit.Specifier
}

func (iit *InvalidIteratorType) GetNote() string {
  return ""
}

func (iit *InvalidIteratorType) GetID() string {
  return "hyb041W"
}

func (iit *InvalidIteratorType) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnnecessaryEmptyIdentifier struct {
  Specifier Snippet
  Context string
}

func (uei *UnnecessaryEmptyIdentifier) GetMessage() string {
  return fmt.Sprintf("unnecessary use of empty identifier ('_') %s", uei.Context)
}

func (uei *UnnecessaryEmptyIdentifier) GetSpecifier() Snippet {
  return uei.Specifier
}

func (uei *UnnecessaryEmptyIdentifier) GetNote() string {
  return ""
}

func (uei *UnnecessaryEmptyIdentifier) GetID() string {
  return "hyb042W"
}

func (uei *UnnecessaryEmptyIdentifier) GetAlertType() Type {
  return Warning
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentAccessToItself struct {
  Specifier Snippet
}

func (eati *EnvironmentAccessToItself) GetMessage() string {
  return "an environment cannot access itself"
}

func (eati *EnvironmentAccessToItself) GetSpecifier() Snippet {
  return eati.Specifier
}

func (eati *EnvironmentAccessToItself) GetNote() string {
  return ""
}

func (eati *EnvironmentAccessToItself) GetID() string {
  return "hyb043W"
}

func (eati *EnvironmentAccessToItself) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EntityConversionWithOrCondition struct {
  Specifier Snippet
}

func (ecwoc *EntityConversionWithOrCondition) GetMessage() string {
  return "cannot convert an entity with an 'or' condition"
}

func (ecwoc *EntityConversionWithOrCondition) GetSpecifier() Snippet {
  return ecwoc.Specifier
}

func (ecwoc *EntityConversionWithOrCondition) GetNote() string {
  return ""
}

func (ecwoc *EntityConversionWithOrCondition) GetID() string {
  return "hyb044W"
}

func (ecwoc *EntityConversionWithOrCondition) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCondition struct {
  Specifier Snippet
  Context string
}

func (ic *InvalidCondition) GetMessage() string {
  return fmt.Sprintf("invalid condition %s", ic.Context)
}

func (ic *InvalidCondition) GetSpecifier() Snippet {
  return ic.Specifier
}

func (ic *InvalidCondition) GetNote() string {
  return "conditions always have to evaluate to either true or false"
}

func (ic *InvalidCondition) GetID() string {
  return "hyb045W"
}

func (ic *InvalidCondition) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidRepeatIterator struct {
  Specifier Snippet
  Type string
}

func (iri *InvalidRepeatIterator) GetMessage() string {
  return fmt.Sprintf("invalid repeat iterator of type '%s'", iri.Type)
}

func (iri *InvalidRepeatIterator) GetSpecifier() Snippet {
  return iri.Specifier
}

func (iri *InvalidRepeatIterator) GetNote() string {
  return "repeat iterator must be a numerical type"
}

func (iri *InvalidRepeatIterator) GetID() string {
  return "hyb046W"
}

func (iri *InvalidRepeatIterator) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InconsistentRepeatTypes struct {
  Specifier Snippet
  From string
  Skip string
  Iterator string
}

func (irt *InconsistentRepeatTypes) GetMessage() string {
  return fmt.Sprintf("repeat types are inconsistent (from:'%s', by:'%s', to:'%s')", irt.From, irt.Skip, irt.Iterator)
}

func (irt *InconsistentRepeatTypes) GetSpecifier() Snippet {
  return irt.Specifier
}

func (irt *InconsistentRepeatTypes) GetNote() string {
  return ""
}

func (irt *InconsistentRepeatTypes) GetID() string {
  return "hyb047W"
}

func (irt *InconsistentRepeatTypes) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type OfficialEntityConversion struct {
  Specifier Snippet
}

func (oec *OfficialEntityConversion) GetMessage() string {
  return "conversion of an official entity to a hybroid entity is not possible"
}

func (oec *OfficialEntityConversion) GetSpecifier() Snippet {
  return oec.Specifier
}

func (oec *OfficialEntityConversion) GetNote() string {
  return ""
}

func (oec *OfficialEntityConversion) GetID() string {
  return "hyb048W"
}

func (oec *OfficialEntityConversion) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidType struct {
  Specifier Snippet
  Expected string
  Got string
  Context string
}

func (it *InvalidType) GetMessage() string {
  return fmt.Sprintf("expected %s, got '%s' %s", it.Expected, it.Got, it.Context)
}

func (it *InvalidType) GetSpecifier() Snippet {
  return it.Specifier
}

func (it *InvalidType) GetNote() string {
  return ""
}

func (it *InvalidType) GetID() string {
  return "hyb049W"
}

func (it *InvalidType) GetAlertType() Type {
  return Error
}

