// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
	"strings"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForbiddenTypeInEnvironment struct {
	Specifier Snippet
	Type      string
	Envs      []string
}

func (ftie *ForbiddenTypeInEnvironment) Message() string {
	return fmt.Sprintf("cannot have a %s in the following environments: %s", ftie.Type, strings.Join(ftie.Envs, ", "))
}

func (ftie *ForbiddenTypeInEnvironment) SnippetSpecifier() Snippet {
	return ftie.Specifier
}

func (ftie *ForbiddenTypeInEnvironment) Note() string {
	return ""
}

func (ftie *ForbiddenTypeInEnvironment) ID() string {
	return "hyb001W"
}

func (ftie *ForbiddenTypeInEnvironment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
	Specifier Snippet
	Type      string
}

func (iet *InvalidEnvironmentType) Message() string {
	return fmt.Sprintf("'%s' is not a valid environment type", iet.Type)
}

func (iet *InvalidEnvironmentType) SnippetSpecifier() Snippet {
	return iet.Specifier
}

func (iet *InvalidEnvironmentType) Note() string {
	return "environment type can be 'Level', 'Mesh', 'Sound' or 'Shared'"
}

func (iet *InvalidEnvironmentType) ID() string {
	return "hyb002W"
}

func (iet *InvalidEnvironmentType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentRedaclaration struct {
	Specifier Snippet
}

func (er *EnvironmentRedaclaration) Message() string {
	return "cannot redeclare an environment"
}

func (er *EnvironmentRedaclaration) SnippetSpecifier() Snippet {
	return er.Specifier
}

func (er *EnvironmentRedaclaration) Note() string {
	return ""
}

func (er *EnvironmentRedaclaration) ID() string {
	return "hyb003W"
}

func (er *EnvironmentRedaclaration) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
	Specifier Snippet
}

func (ee *ExpectedEnvironment) Message() string {
	return "expected environment declaration"
}

func (ee *ExpectedEnvironment) SnippetSpecifier() Snippet {
	return ee.Specifier
}

func (ee *ExpectedEnvironment) Note() string {
	return "the first declaration in any Hybroid file has to be an environment declaration"
}

func (ee *ExpectedEnvironment) ID() string {
	return "hyb004W"
}

func (ee *ExpectedEnvironment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateEnvironmentNames struct {
	Specifier Snippet
	Path1     string
	Path2     string
}

func (den *DuplicateEnvironmentNames) Message() string {
	return fmt.Sprintf("duplicate environment names found between '%s' and '%s'", den.Path1, den.Path2)
}

func (den *DuplicateEnvironmentNames) SnippetSpecifier() Snippet {
	return den.Specifier
}

func (den *DuplicateEnvironmentNames) Note() string {
	return ""
}

func (den *DuplicateEnvironmentNames) ID() string {
	return "hyb005W"
}

func (den *DuplicateEnvironmentNames) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidAccessValue struct {
	Specifier Snippet
	Type      string
}

func (iav *InvalidAccessValue) Message() string {
	return fmt.Sprintf("value is of type '%s', so it cannot be accessed from", iav.Type)
}

func (iav *InvalidAccessValue) SnippetSpecifier() Snippet {
	return iav.Specifier
}

func (iav *InvalidAccessValue) Note() string {
	return "only lists, maps, classes, entities, structs and enums can be used to access values from"
}

func (iav *InvalidAccessValue) ID() string {
	return "hyb006W"
}

func (iav *InvalidAccessValue) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type FieldAccessOnListOrMap struct {
	Specifier  Snippet
	Field      string
	AccessType string
}

func (faolom *FieldAccessOnListOrMap) Message() string {
	return fmt.Sprintf("cannot access field '%s' from the %s", faolom.Field, faolom.AccessType)
}

func (faolom *FieldAccessOnListOrMap) SnippetSpecifier() Snippet {
	return faolom.Specifier
}

func (faolom *FieldAccessOnListOrMap) Note() string {
	return fmt.Sprintf("to access a value from a %s you use brackets, e.g. example[\"%s\"]", faolom.AccessType, faolom.Field)
}

func (faolom *FieldAccessOnListOrMap) ID() string {
	return "hyb007W"
}

func (faolom *FieldAccessOnListOrMap) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MemberAccessOnNonListOrMap struct {
	Specifier  Snippet
	Member     string
	AccessType string
}

func (maonlom *MemberAccessOnNonListOrMap) Message() string {
	return fmt.Sprintf("cannot access member '[%s]' from the %s", maonlom.Member, maonlom.AccessType)
}

func (maonlom *MemberAccessOnNonListOrMap) SnippetSpecifier() Snippet {
	return maonlom.Specifier
}

func (maonlom *MemberAccessOnNonListOrMap) Note() string {
	return "to access a value you use a dot and then an identifier, e.g. example.identifier"
}

func (maonlom *MemberAccessOnNonListOrMap) ID() string {
	return "hyb008W"
}

func (maonlom *MemberAccessOnNonListOrMap) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidMemberIndex struct {
	Specifier  Snippet
	AccessType string
	Index      string
}

func (imi *InvalidMemberIndex) Message() string {
	return fmt.Sprintf("'%s' is not of type number to be an index for the %s", imi.Index, imi.AccessType)
}

func (imi *InvalidMemberIndex) SnippetSpecifier() Snippet {
	return imi.Specifier
}

func (imi *InvalidMemberIndex) Note() string {
	return "for lists, an index (number) is used to access values, for maps, a key (text) is used"
}

func (imi *InvalidMemberIndex) ID() string {
	return "hyb009W"
}

func (imi *InvalidMemberIndex) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidField struct {
	Specifier  Snippet
	AccessType string
	FieldName  string
}

func (_if *InvalidField) Message() string {
	return fmt.Sprintf("field '%s' does not belong to '%s'", _if.FieldName, _if.AccessType)
}

func (_if *InvalidField) SnippetSpecifier() Snippet {
	return _if.Specifier
}

func (_if *InvalidField) Note() string {
	return ""
}

func (_if *InvalidField) ID() string {
	return "hyb010W"
}

func (_if *InvalidField) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MixedMapOrListContents struct {
	Specifier     Snippet
	ContainerType string
	Type1         string
	Type2         string
}

func (mmolc *MixedMapOrListContents) Message() string {
	return fmt.Sprintf("%s member is of type '%s', but the previous one was '%s'", mmolc.ContainerType, mmolc.Type1, mmolc.Type2)
}

func (mmolc *MixedMapOrListContents) SnippetSpecifier() Snippet {
	return mmolc.Specifier
}

func (mmolc *MixedMapOrListContents) Note() string {
	return ""
}

func (mmolc *MixedMapOrListContents) ID() string {
	return "hyb011W"
}

func (mmolc *MixedMapOrListContents) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCallerType struct {
	Specifier Snippet
	Type      string
}

func (ict *InvalidCallerType) Message() string {
	return fmt.Sprintf("cannot call value of type '%s' as a function", ict.Type)
}

func (ict *InvalidCallerType) SnippetSpecifier() Snippet {
	return ict.Specifier
}

func (ict *InvalidCallerType) Note() string {
	return ""
}

func (ict *InvalidCallerType) ID() string {
	return "hyb012W"
}

func (ict *InvalidCallerType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MethodOrFieldNotFound struct {
	Specifier Snippet
	Name      string
}

func (mofnf *MethodOrFieldNotFound) Message() string {
	return fmt.Sprintf("no method or field named '%s'", mofnf.Name)
}

func (mofnf *MethodOrFieldNotFound) SnippetSpecifier() Snippet {
	return mofnf.Specifier
}

func (mofnf *MethodOrFieldNotFound) Note() string {
	return ""
}

func (mofnf *MethodOrFieldNotFound) ID() string {
	return "hyb013W"
}

func (mofnf *MethodOrFieldNotFound) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForeignLocalVariableAccess struct {
	Specifier Snippet
	Name      string
}

func (flva *ForeignLocalVariableAccess) Message() string {
	return fmt.Sprintf("cannot access local variable '%s' belonging to a different environment", flva.Name)
}

func (flva *ForeignLocalVariableAccess) SnippetSpecifier() Snippet {
	return flva.Specifier
}

func (flva *ForeignLocalVariableAccess) Note() string {
	return ""
}

func (flva *ForeignLocalVariableAccess) ID() string {
	return "hyb014W"
}

func (flva *ForeignLocalVariableAccess) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidArgumentType struct {
	Specifier    Snippet
	GivenType    string
	ExpectedType string
}

func (iat *InvalidArgumentType) Message() string {
	return fmt.Sprintf("argument was of type %s, but should be %s", iat.GivenType, iat.ExpectedType)
}

func (iat *InvalidArgumentType) SnippetSpecifier() Snippet {
	return iat.Specifier
}

func (iat *InvalidArgumentType) Note() string {
	return ""
}

func (iat *InvalidArgumentType) ID() string {
	return "hyb015W"
}

func (iat *InvalidArgumentType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type PublicDeclarationInLocalScope struct {
	Specifier Snippet
}

func (pdils *PublicDeclarationInLocalScope) Message() string {
	return "cannot have a public declaration that is in a local scope"
}

func (pdils *PublicDeclarationInLocalScope) SnippetSpecifier() Snippet {
	return pdils.Specifier
}

func (pdils *PublicDeclarationInLocalScope) Note() string {
	return ""
}

func (pdils *PublicDeclarationInLocalScope) ID() string {
	return "hyb016W"
}

func (pdils *PublicDeclarationInLocalScope) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type Redeclaration struct {
	Specifier Snippet
	VarName   string
	DeclType  string
}

func (r *Redeclaration) Message() string {
	return fmt.Sprintf("a %s named '%s' already exists", r.DeclType, r.VarName)
}

func (r *Redeclaration) SnippetSpecifier() Snippet {
	return r.Specifier
}

func (r *Redeclaration) Note() string {
	return ""
}

func (r *Redeclaration) ID() string {
	return "hyb017W"
}

func (r *Redeclaration) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnnecessaryTypeInConstDeclaration struct {
	Specifier Snippet
}

func (uticd *UnnecessaryTypeInConstDeclaration) Message() string {
	return "an explicit type is not necessary for a const declaration"
}

func (uticd *UnnecessaryTypeInConstDeclaration) SnippetSpecifier() Snippet {
	return uticd.Specifier
}

func (uticd *UnnecessaryTypeInConstDeclaration) Note() string {
	return ""
}

func (uticd *UnnecessaryTypeInConstDeclaration) ID() string {
	return "hyb018W"
}

func (uticd *UnnecessaryTypeInConstDeclaration) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type NoValueGivenForConstant struct {
	Specifier Snippet
}

func (nvgfc *NoValueGivenForConstant) Message() string {
	return "constant must be declared with a value"
}

func (nvgfc *NoValueGivenForConstant) SnippetSpecifier() Snippet {
	return nvgfc.Specifier
}

func (nvgfc *NoValueGivenForConstant) Note() string {
	return ""
}

func (nvgfc *NoValueGivenForConstant) ID() string {
	return "hyb019W"
}

func (nvgfc *NoValueGivenForConstant) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TooFewValuesGiven struct {
	Specifier      Snippet
	RequiredAmount int
	Context        string
}

func (tfvg *TooFewValuesGiven) Message() string {
	return fmt.Sprintf("%d more value(s) required %s", tfvg.RequiredAmount, tfvg.Context)
}

func (tfvg *TooFewValuesGiven) SnippetSpecifier() Snippet {
	return tfvg.Specifier
}

func (tfvg *TooFewValuesGiven) Note() string {
	return ""
}

func (tfvg *TooFewValuesGiven) ID() string {
	return "hyb020W"
}

func (tfvg *TooFewValuesGiven) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeRequiredInDeclaration struct {
	Specifier Snippet
}

func (etrid *ExplicitTypeRequiredInDeclaration) Message() string {
	return "a variable declared without a value requires an explicit type"
}

func (etrid *ExplicitTypeRequiredInDeclaration) SnippetSpecifier() Snippet {
	return etrid.Specifier
}

func (etrid *ExplicitTypeRequiredInDeclaration) Note() string {
	return ""
}

func (etrid *ExplicitTypeRequiredInDeclaration) ID() string {
	return "hyb021W"
}

func (etrid *ExplicitTypeRequiredInDeclaration) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeMismatch struct {
	Specifier    Snippet
	ExplicitType string
	ValueType    string
}

func (etm *ExplicitTypeMismatch) Message() string {
	return fmt.Sprintf("variable was given explicit type '%s', but its value is a '%s'", etm.ExplicitType, etm.ValueType)
}

func (etm *ExplicitTypeMismatch) SnippetSpecifier() Snippet {
	return etm.Specifier
}

func (etm *ExplicitTypeMismatch) Note() string {
	return ""
}

func (etm *ExplicitTypeMismatch) ID() string {
	return "hyb022W"
}

func (etm *ExplicitTypeMismatch) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeNotAllowed struct {
	Specifier    Snippet
	ExplicitType string
}

func (etna *ExplicitTypeNotAllowed) Message() string {
	return fmt.Sprintf("cannot create a default value from the explicit type '%s'", etna.ExplicitType)
}

func (etna *ExplicitTypeNotAllowed) SnippetSpecifier() Snippet {
	return etna.Specifier
}

func (etna *ExplicitTypeNotAllowed) Note() string {
	return "some types don't have default values, like entities and classes"
}

func (etna *ExplicitTypeNotAllowed) ID() string {
	return "hyb023W"
}

func (etna *ExplicitTypeNotAllowed) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TooManyValuesGiven struct {
	Specifier   Snippet
	ExtraAmount int
	Context     string
}

func (tmvg *TooManyValuesGiven) Message() string {
	return fmt.Sprintf("%d less value(s) required %s", tmvg.ExtraAmount, tmvg.Context)
}

func (tmvg *TooManyValuesGiven) SnippetSpecifier() Snippet {
	return tmvg.Specifier
}

func (tmvg *TooManyValuesGiven) Note() string {
	return ""
}

func (tmvg *TooManyValuesGiven) ID() string {
	return "hyb024W"
}

func (tmvg *TooManyValuesGiven) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ImportCycle struct {
	Specifier Snippet
	HybPaths  []string
}

func (ic *ImportCycle) Message() string {
	return fmt.Sprintf("import cycle detected: %s", strings.Join(ic.HybPaths, " -> "))
}

func (ic *ImportCycle) SnippetSpecifier() Snippet {
	return ic.Specifier
}

func (ic *ImportCycle) Note() string {
	return ""
}

func (ic *ImportCycle) ID() string {
	return "hyb025W"
}

func (ic *ImportCycle) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UndeclaredVariableAccess struct {
	Specifier Snippet
	Var       string
	Context   string
}

func (uva *UndeclaredVariableAccess) Message() string {
	return fmt.Sprintf("'%s' is not a declared variable %s", uva.Var, uva.Context)
}

func (uva *UndeclaredVariableAccess) SnippetSpecifier() Snippet {
	return uva.Specifier
}

func (uva *UndeclaredVariableAccess) Note() string {
	return ""
}

func (uva *UndeclaredVariableAccess) ID() string {
	return "hyb026W"
}

func (uva *UndeclaredVariableAccess) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ConstValueAssignment struct {
	Specifier Snippet
}

func (cva *ConstValueAssignment) Message() string {
	return "cannot modify a constant value"
}

func (cva *ConstValueAssignment) SnippetSpecifier() Snippet {
	return cva.Specifier
}

func (cva *ConstValueAssignment) Note() string {
	return ""
}

func (cva *ConstValueAssignment) ID() string {
	return "hyb027W"
}

func (cva *ConstValueAssignment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type AssignmentTypeMismatch struct {
	Specifier Snippet
	VarType   string
	ValType   string
}

func (atm *AssignmentTypeMismatch) Message() string {
	return fmt.Sprintf("variable is of type '%s', but a value of '%s' was assigned to it", atm.VarType, atm.ValType)
}

func (atm *AssignmentTypeMismatch) SnippetSpecifier() Snippet {
	return atm.Specifier
}

func (atm *AssignmentTypeMismatch) Note() string {
	return ""
}

func (atm *AssignmentTypeMismatch) ID() string {
	return "hyb028W"
}

func (atm *AssignmentTypeMismatch) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidTypeInCompoundAssignment struct {
	Specifier Snippet
	Type      string
}

func (itica *InvalidTypeInCompoundAssignment) Message() string {
	return fmt.Sprintf("the type '%s' is not allowed in compound assignment", itica.Type)
}

func (itica *InvalidTypeInCompoundAssignment) SnippetSpecifier() Snippet {
	return itica.Specifier
}

func (itica *InvalidTypeInCompoundAssignment) Note() string {
	return "only numerical types are allowed, like number or fixed"
}

func (itica *InvalidTypeInCompoundAssignment) ID() string {
	return "hyb029W"
}

func (itica *InvalidTypeInCompoundAssignment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidUseOfSelf struct {
	Specifier Snippet
}

func (iuos *InvalidUseOfSelf) Message() string {
	return "cannot use self outside of class or entity"
}

func (iuos *InvalidUseOfSelf) SnippetSpecifier() Snippet {
	return iuos.Specifier
}

func (iuos *InvalidUseOfSelf) Note() string {
	return "you're also not allowed to use self inside anonymous functions of class/entity fields"
}

func (iuos *InvalidUseOfSelf) ID() string {
	return "hyb030W"
}

func (iuos *InvalidUseOfSelf) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnreachableCode struct {
	Specifier Snippet
}

func (uc *UnreachableCode) Message() string {
	return "unreachable code detected"
}

func (uc *UnreachableCode) SnippetSpecifier() Snippet {
	return uc.Specifier
}

func (uc *UnreachableCode) Note() string {
	return ""
}

func (uc *UnreachableCode) ID() string {
	return "hyb031W"
}

func (uc *UnreachableCode) AlertType() Type {
	return Warning
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidUseOfExitStmt struct {
	Specifier Snippet
	ExitNode  string
	Context   string
}

func (iuoes *InvalidUseOfExitStmt) Message() string {
	return fmt.Sprintf("cannot use '%s' outside of %s", iuoes.ExitNode, iuoes.Context)
}

func (iuoes *InvalidUseOfExitStmt) SnippetSpecifier() Snippet {
	return iuoes.Specifier
}

func (iuoes *InvalidUseOfExitStmt) Note() string {
	return ""
}

func (iuoes *InvalidUseOfExitStmt) ID() string {
	return "hyb032W"
}

func (iuoes *InvalidUseOfExitStmt) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypeMismatch struct {
	Specifier Snippet
	Type1     string
	Type2     string
	Context   string
}

func (tm *TypeMismatch) Message() string {
	return fmt.Sprintf("expected %s, got '%s' %s", tm.Type1, tm.Type2, tm.Context)
}

func (tm *TypeMismatch) SnippetSpecifier() Snippet {
	return tm.Specifier
}

func (tm *TypeMismatch) Note() string {
	return ""
}

func (tm *TypeMismatch) ID() string {
	return "hyb033W"
}

func (tm *TypeMismatch) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidStmtInLocalBlock struct {
	Specifier Snippet
	StmtType  string
}

func (isilb *InvalidStmtInLocalBlock) Message() string {
	return fmt.Sprintf("%s must be in the global scope", isilb.StmtType)
}

func (isilb *InvalidStmtInLocalBlock) SnippetSpecifier() Snippet {
	return isilb.Specifier
}

func (isilb *InvalidStmtInLocalBlock) Note() string {
	return ""
}

func (isilb *InvalidStmtInLocalBlock) ID() string {
	return "hyb034W"
}

func (isilb *InvalidStmtInLocalBlock) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnallowedLibraryUse struct {
	Specifier     Snippet
	Library       string
	UnallowedEnvs string
}

func (ulu *UnallowedLibraryUse) Message() string {
	return fmt.Sprintf("cannot use the %s library in a %s environment", ulu.Library, ulu.UnallowedEnvs)
}

func (ulu *UnallowedLibraryUse) SnippetSpecifier() Snippet {
	return ulu.Specifier
}

func (ulu *UnallowedLibraryUse) Note() string {
	return ""
}

func (ulu *UnallowedLibraryUse) ID() string {
	return "hyb035W"
}

func (ulu *UnallowedLibraryUse) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentAccess struct {
	Specifier Snippet
	EnvName   string
}

func (iea *InvalidEnvironmentAccess) Message() string {
	return fmt.Sprintf("environment named '%s' does not exist", iea.EnvName)
}

func (iea *InvalidEnvironmentAccess) SnippetSpecifier() Snippet {
	return iea.Specifier
}

func (iea *InvalidEnvironmentAccess) Note() string {
	return ""
}

func (iea *InvalidEnvironmentAccess) ID() string {
	return "hyb036W"
}

func (iea *InvalidEnvironmentAccess) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentReuse struct {
	Specifier Snippet
	EnvName   string
}

func (er *EnvironmentReuse) Message() string {
	return fmt.Sprintf("environment named '%s' is already imported through use statement", er.EnvName)
}

func (er *EnvironmentReuse) SnippetSpecifier() Snippet {
	return er.Specifier
}

func (er *EnvironmentReuse) Note() string {
	return ""
}

func (er *EnvironmentReuse) ID() string {
	return "hyb037W"
}

func (er *EnvironmentReuse) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidIteratorType struct {
	Specifier Snippet
	Type      string
}

func (iit *InvalidIteratorType) Message() string {
	return fmt.Sprintf("a for loop iterator must be a map or a list (found: '%s')", iit.Type)
}

func (iit *InvalidIteratorType) SnippetSpecifier() Snippet {
	return iit.Specifier
}

func (iit *InvalidIteratorType) Note() string {
	return ""
}

func (iit *InvalidIteratorType) ID() string {
	return "hyb038W"
}

func (iit *InvalidIteratorType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnnecessaryEmptyIdentifier struct {
	Specifier Snippet
	Context   string
}

func (uei *UnnecessaryEmptyIdentifier) Message() string {
	return fmt.Sprintf("unnecessary use of empty identifier ('_') %s", uei.Context)
}

func (uei *UnnecessaryEmptyIdentifier) SnippetSpecifier() Snippet {
	return uei.Specifier
}

func (uei *UnnecessaryEmptyIdentifier) Note() string {
	return ""
}

func (uei *UnnecessaryEmptyIdentifier) ID() string {
	return "hyb039W"
}

func (uei *UnnecessaryEmptyIdentifier) AlertType() Type {
	return Warning
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentAccessToItself struct {
	Specifier Snippet
}

func (eati *EnvironmentAccessToItself) Message() string {
	return "an environment cannot access itself"
}

func (eati *EnvironmentAccessToItself) SnippetSpecifier() Snippet {
	return eati.Specifier
}

func (eati *EnvironmentAccessToItself) Note() string {
	return ""
}

func (eati *EnvironmentAccessToItself) ID() string {
	return "hyb040W"
}

func (eati *EnvironmentAccessToItself) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EntityConversionWithOrCondition struct {
	Specifier Snippet
}

func (ecwoc *EntityConversionWithOrCondition) Message() string {
	return "cannot convert an entity with an 'or' condition"
}

func (ecwoc *EntityConversionWithOrCondition) SnippetSpecifier() Snippet {
	return ecwoc.Specifier
}

func (ecwoc *EntityConversionWithOrCondition) Note() string {
	return ""
}

func (ecwoc *EntityConversionWithOrCondition) ID() string {
	return "hyb041W"
}

func (ecwoc *EntityConversionWithOrCondition) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCondition struct {
	Specifier Snippet
	Context   string
}

func (ic *InvalidCondition) Message() string {
	return fmt.Sprintf("invalid condition %s", ic.Context)
}

func (ic *InvalidCondition) SnippetSpecifier() Snippet {
	return ic.Specifier
}

func (ic *InvalidCondition) Note() string {
	return "conditions always have to evaluate to either true or false"
}

func (ic *InvalidCondition) ID() string {
	return "hyb042W"
}

func (ic *InvalidCondition) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidRepeatIterator struct {
	Specifier Snippet
	Type      string
}

func (iri *InvalidRepeatIterator) Message() string {
	return fmt.Sprintf("invalid repeat iterator of type '%s'", iri.Type)
}

func (iri *InvalidRepeatIterator) SnippetSpecifier() Snippet {
	return iri.Specifier
}

func (iri *InvalidRepeatIterator) Note() string {
	return "repeat iterator must be a numerical type"
}

func (iri *InvalidRepeatIterator) ID() string {
	return "hyb043W"
}

func (iri *InvalidRepeatIterator) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InconsistentRepeatTypes struct {
	Specifier Snippet
	From      string
	Skip      string
	Iterator  string
}

func (irt *InconsistentRepeatTypes) Message() string {
	return fmt.Sprintf("repeat types are inconsistent (from:'%s', by:'%s', to:'%s')", irt.From, irt.Skip, irt.Iterator)
}

func (irt *InconsistentRepeatTypes) SnippetSpecifier() Snippet {
	return irt.Specifier
}

func (irt *InconsistentRepeatTypes) Note() string {
	return ""
}

func (irt *InconsistentRepeatTypes) ID() string {
	return "hyb044W"
}

func (irt *InconsistentRepeatTypes) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type OfficialEntityConversion struct {
	Specifier Snippet
}

func (oec *OfficialEntityConversion) Message() string {
	return "conversion of an official entity to a hybroid entity is not possible"
}

func (oec *OfficialEntityConversion) SnippetSpecifier() Snippet {
	return oec.Specifier
}

func (oec *OfficialEntityConversion) Note() string {
	return ""
}

func (oec *OfficialEntityConversion) ID() string {
	return "hyb045W"
}

func (oec *OfficialEntityConversion) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironment struct {
	Specifier Snippet
}

func (ie *InvalidEnvironment) Message() string {
	return "there is no environment with that path"
}

func (ie *InvalidEnvironment) SnippetSpecifier() Snippet {
	return ie.Specifier
}

func (ie *InvalidEnvironment) Note() string {
	return ""
}

func (ie *InvalidEnvironment) ID() string {
	return "hyb046W"
}

func (ie *InvalidEnvironment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentAccessAmbiguity struct {
	Specifier Snippet
	Envs      []string
	Context   string
}

func (eaa *EnvironmentAccessAmbiguity) Message() string {
	return fmt.Sprintf("the type '%s' can be found on multiple environments: %s", eaa.Context, strings.Join(eaa.Envs, ", "))
}

func (eaa *EnvironmentAccessAmbiguity) SnippetSpecifier() Snippet {
	return eaa.Specifier
}

func (eaa *EnvironmentAccessAmbiguity) Note() string {
	return ""
}

func (eaa *EnvironmentAccessAmbiguity) ID() string {
	return "hyb047W"
}

func (eaa *EnvironmentAccessAmbiguity) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type NotAllCodePathsExit struct {
	Specifier Snippet
	ExitType  string
}

func (nacpe *NotAllCodePathsExit) Message() string {
	return fmt.Sprintf("not all code paths %s", nacpe.ExitType)
}

func (nacpe *NotAllCodePathsExit) SnippetSpecifier() Snippet {
	return nacpe.Specifier
}

func (nacpe *NotAllCodePathsExit) Note() string {
	return ""
}

func (nacpe *NotAllCodePathsExit) ID() string {
	return "hyb048W"
}

func (nacpe *NotAllCodePathsExit) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InsufficientCases struct {
	Specifier Snippet
}

func (ic *InsufficientCases) Message() string {
	return "match statement must have at least 1 non-default case"
}

func (ic *InsufficientCases) SnippetSpecifier() Snippet {
	return ic.Specifier
}

func (ic *InsufficientCases) Note() string {
	return ""
}

func (ic *InsufficientCases) ID() string {
	return "hyb049W"
}

func (ic *InsufficientCases) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DefaultCaseMissing struct {
	Specifier Snippet
}

func (dcm *DefaultCaseMissing) Message() string {
	return "match expression must have a default case"
}

func (dcm *DefaultCaseMissing) SnippetSpecifier() Snippet {
	return dcm.Specifier
}

func (dcm *DefaultCaseMissing) Note() string {
	return "default cases start with 'else'"
}

func (dcm *DefaultCaseMissing) ID() string {
	return "hyb050W"
}

func (dcm *DefaultCaseMissing) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCaseType struct {
	Specifier      Snippet
	MatchValueType string
	CaseValueType  string
}

func (ict *InvalidCaseType) Message() string {
	return fmt.Sprintf("match value is of type '%s', but case value is of type '%s'", ict.MatchValueType, ict.CaseValueType)
}

func (ict *InvalidCaseType) SnippetSpecifier() Snippet {
	return ict.Specifier
}

func (ict *InvalidCaseType) Note() string {
	return ""
}

func (ict *InvalidCaseType) ID() string {
	return "hyb051W"
}

func (ict *InvalidCaseType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type LiteralCondition struct {
	Specifier      Snippet
	ConditionValue string
}

func (lc *LiteralCondition) Message() string {
	return fmt.Sprintf("condition is always %s", lc.ConditionValue)
}

func (lc *LiteralCondition) SnippetSpecifier() Snippet {
	return lc.Specifier
}

func (lc *LiteralCondition) Note() string {
	return ""
}

func (lc *LiteralCondition) ID() string {
	return "hyb052W"
}

func (lc *LiteralCondition) AlertType() Type {
	return Warning
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypesMismatch struct {
	Specifier Snippet
	Value1    string
	Type1     string
	Value2    string
	Type2     string
}

func (tm *TypesMismatch) Message() string {
	return fmt.Sprintf("%s is of type '%s', but %s is of type '%s'", tm.Value1, tm.Type1, tm.Value2, tm.Type2)
}

func (tm *TypesMismatch) SnippetSpecifier() Snippet {
	return tm.Specifier
}

func (tm *TypesMismatch) Note() string {
	return ""
}

func (tm *TypesMismatch) ID() string {
	return "hyb053W"
}

func (tm *TypesMismatch) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingConstructor struct {
	Specifier       Snippet
	ConstructorType string
	Context         string
}

func (mc *MissingConstructor) Message() string {
	return fmt.Sprintf("missing '%s' constructor %s", mc.ConstructorType, mc.Context)
}

func (mc *MissingConstructor) SnippetSpecifier() Snippet {
	return mc.Specifier
}

func (mc *MissingConstructor) Note() string {
	return ""
}

func (mc *MissingConstructor) ID() string {
	return "hyb054W"
}

func (mc *MissingConstructor) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingDestroy struct {
	Specifier Snippet
}

func (md *MissingDestroy) Message() string {
	return "missing 'destroy' destructor in entity declaration"
}

func (md *MissingDestroy) SnippetSpecifier() Snippet {
	return md.Specifier
}

func (md *MissingDestroy) Note() string {
	return ""
}

func (md *MissingDestroy) ID() string {
	return "hyb055W"
}

func (md *MissingDestroy) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UninitializedFieldInConstructor struct {
	Specifier Snippet
	VarName   string
	Context   string
}

func (ufic *UninitializedFieldInConstructor) Message() string {
	return fmt.Sprintf("variable '%s' was not initialized in the constructor %s", ufic.VarName, ufic.Context)
}

func (ufic *UninitializedFieldInConstructor) SnippetSpecifier() Snippet {
	return ufic.Specifier
}

func (ufic *UninitializedFieldInConstructor) Note() string {
	return ""
}

func (ufic *UninitializedFieldInConstructor) ID() string {
	return "hyb056W"
}

func (ufic *UninitializedFieldInConstructor) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypeRedeclaration struct {
	Specifier Snippet
	TypeName  string
}

func (tr *TypeRedeclaration) Message() string {
	return fmt.Sprintf("type '%s' already exists", tr.TypeName)
}

func (tr *TypeRedeclaration) SnippetSpecifier() Snippet {
	return tr.Specifier
}

func (tr *TypeRedeclaration) Note() string {
	return ""
}

func (tr *TypeRedeclaration) ID() string {
	return "hyb057W"
}

func (tr *TypeRedeclaration) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCallAsArgument struct {
	Specifier Snippet
}

func (icaa *InvalidCallAsArgument) Message() string {
	return "cannot have a call that returns more than 1 value as an argument"
}

func (icaa *InvalidCallAsArgument) SnippetSpecifier() Snippet {
	return icaa.Specifier
}

func (icaa *InvalidCallAsArgument) Note() string {
	return ""
}

func (icaa *InvalidCallAsArgument) ID() string {
	return "hyb058W"
}

func (icaa *InvalidCallAsArgument) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneVariadicParameter struct {
	Specifier Snippet
}

func (mtovp *MoreThanOneVariadicParameter) Message() string {
	return "cannot have more than one variadic function parameter"
}

func (mtovp *MoreThanOneVariadicParameter) SnippetSpecifier() Snippet {
	return mtovp.Specifier
}

func (mtovp *MoreThanOneVariadicParameter) Note() string {
	return ""
}

func (mtovp *MoreThanOneVariadicParameter) ID() string {
	return "hyb059W"
}

func (mtovp *MoreThanOneVariadicParameter) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type VariadicParameterNotAtEnd struct {
	Specifier Snippet
}

func (vpnae *VariadicParameterNotAtEnd) Message() string {
	return "variadic parameters must be at the end of the function parameters"
}

func (vpnae *VariadicParameterNotAtEnd) SnippetSpecifier() Snippet {
	return vpnae.Specifier
}

func (vpnae *VariadicParameterNotAtEnd) Note() string {
	return ""
}

func (vpnae *VariadicParameterNotAtEnd) ID() string {
	return "hyb060W"
}

func (vpnae *VariadicParameterNotAtEnd) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateElement struct {
	Specifier Snippet
	Element   string
	ElemName  string
}

func (de *DuplicateElement) Message() string {
	return fmt.Sprintf("the %s '%s' already exists", de.Element, de.ElemName)
}

func (de *DuplicateElement) SnippetSpecifier() Snippet {
	return de.Specifier
}

func (de *DuplicateElement) Note() string {
	return ""
}

func (de *DuplicateElement) ID() string {
	return "hyb061W"
}

func (de *DuplicateElement) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEntityFunctionSignature struct {
	Specifier      Snippet
	Got            string
	Expected       string
	EntityFuncType string
}

func (iefs *InvalidEntityFunctionSignature) Message() string {
	return fmt.Sprintf("expected '%s' for %s, got '%s'", iefs.Expected, iefs.EntityFuncType, iefs.Got)
}

func (iefs *InvalidEntityFunctionSignature) SnippetSpecifier() Snippet {
	return iefs.Specifier
}

func (iefs *InvalidEntityFunctionSignature) Note() string {
	return ""
}

func (iefs *InvalidEntityFunctionSignature) ID() string {
	return "hyb062W"
}

func (iefs *InvalidEntityFunctionSignature) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidSpawnerParameters struct {
	Specifier Snippet
}

func (isp *InvalidSpawnerParameters) Message() string {
	return "the first two parameters of the spawner must be fixedpoints"
}

func (isp *InvalidSpawnerParameters) SnippetSpecifier() Snippet {
	return isp.Specifier
}

func (isp *InvalidSpawnerParameters) Note() string {
	return ""
}

func (isp *InvalidSpawnerParameters) ID() string {
	return "hyb063W"
}

func (isp *InvalidSpawnerParameters) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidPewpewVariable struct {
	Specifier Snippet
	PewpewVar string
	Type      string
}

func (ipv *InvalidPewpewVariable) Message() string {
	return fmt.Sprintf("'%s' variable should be global and of type 'list<%s>'", ipv.PewpewVar, ipv.Type)
}

func (ipv *InvalidPewpewVariable) SnippetSpecifier() Snippet {
	return ipv.Specifier
}

func (ipv *InvalidPewpewVariable) Note() string {
	return ""
}

func (ipv *InvalidPewpewVariable) ID() string {
	return "hyb064W"
}

func (ipv *InvalidPewpewVariable) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingPewpewVariable struct {
	Specifier Snippet
	PewpewVar string
	EnvType   string
}

func (mpv *MissingPewpewVariable) Message() string {
	return fmt.Sprintf("A %s environment must have a '%s' variable", mpv.EnvType, mpv.PewpewVar)
}

func (mpv *MissingPewpewVariable) SnippetSpecifier() Snippet {
	return mpv.Specifier
}

func (mpv *MissingPewpewVariable) Note() string {
	return ""
}

func (mpv *MissingPewpewVariable) ID() string {
	return "hyb065W"
}

func (mpv *MissingPewpewVariable) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnallowedEnvironmentAccess struct {
	Specifier Snippet
	Unallowed string
	From      string
}

func (uea *UnallowedEnvironmentAccess) Message() string {
	return fmt.Sprintf("cannot access a %s environment from a %s environment", uea.Unallowed, uea.From)
}

func (uea *UnallowedEnvironmentAccess) SnippetSpecifier() Snippet {
	return uea.Specifier
}

func (uea *UnallowedEnvironmentAccess) Note() string {
	return ""
}

func (uea *UnallowedEnvironmentAccess) ID() string {
	return "hyb066W"
}

func (uea *UnallowedEnvironmentAccess) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidDefaultCasePlacement struct {
	Specifier Snippet
	Context   string
}

func (idcp *InvalidDefaultCasePlacement) Message() string {
	return fmt.Sprintf("the default case must always be at the end %s", idcp.Context)
}

func (idcp *InvalidDefaultCasePlacement) SnippetSpecifier() Snippet {
	return idcp.Specifier
}

func (idcp *InvalidDefaultCasePlacement) Note() string {
	return ""
}

func (idcp *InvalidDefaultCasePlacement) ID() string {
	return "hyb067W"
}

func (idcp *InvalidDefaultCasePlacement) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidType struct {
	Specifier Snippet
	Type      string
	Context   string
}

func (it *InvalidType) Message() string {
	return fmt.Sprintf("cannot have a type '%s' %s", it.Type, it.Context)
}

func (it *InvalidType) SnippetSpecifier() Snippet {
	return it.Specifier
}

func (it *InvalidType) Note() string {
	return ""
}

func (it *InvalidType) ID() string {
	return "hyb068W"
}

func (it *InvalidType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ListIndexOutOfBounds struct {
	Specifier Snippet
	Value     string
}

func (lioob *ListIndexOutOfBounds) Message() string {
	return fmt.Sprintf("list index is '%s', but it must be 1 or more", lioob.Value)
}

func (lioob *ListIndexOutOfBounds) SnippetSpecifier() Snippet {
	return lioob.Specifier
}

func (lioob *ListIndexOutOfBounds) Note() string {
	return ""
}

func (lioob *ListIndexOutOfBounds) ID() string {
	return "hyb069W"
}

func (lioob *ListIndexOutOfBounds) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidListIndex struct {
	Specifier Snippet
}

func (ili *InvalidListIndex) Message() string {
	return "a list index must be a whole number"
}

func (ili *InvalidListIndex) SnippetSpecifier() Snippet {
	return ili.Specifier
}

func (ili *InvalidListIndex) Note() string {
	return ""
}

func (ili *InvalidListIndex) ID() string {
	return "hyb070W"
}

func (ili *InvalidListIndex) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingGenericArgument struct {
	Specifier Snippet
	Type      string
}

func (mga *MissingGenericArgument) Message() string {
	return fmt.Sprintf("generic type '%s' could not be inferred", mga.Type)
}

func (mga *MissingGenericArgument) SnippetSpecifier() Snippet {
	return mga.Specifier
}

func (mga *MissingGenericArgument) Note() string {
	return ""
}

func (mga *MissingGenericArgument) ID() string {
	return "hyb071W"
}

func (mga *MissingGenericArgument) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidAssignment struct {
	Specifier Snippet
}

func (ia *InvalidAssignment) Message() string {
	return "left value was not a variable"
}

func (ia *InvalidAssignment) SnippetSpecifier() Snippet {
	return ia.Specifier
}

func (ia *InvalidAssignment) Note() string {
	return ""
}

func (ia *InvalidAssignment) ID() string {
	return "hyb072W"
}

func (ia *InvalidAssignment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ConflictingVariableNameWithType struct {
	Specifier Snippet
	Type      string
}

func (cvnwt *ConflictingVariableNameWithType) Message() string {
	return fmt.Sprintf("variable name conflicts with type '%s'", cvnwt.Type)
}

func (cvnwt *ConflictingVariableNameWithType) SnippetSpecifier() Snippet {
	return cvnwt.Specifier
}

func (cvnwt *ConflictingVariableNameWithType) Note() string {
	return ""
}

func (cvnwt *ConflictingVariableNameWithType) ID() string {
	return "hyb073W"
}

func (cvnwt *ConflictingVariableNameWithType) AlertType() Type {
	return Error
}
