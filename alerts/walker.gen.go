// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
	"hybroid/core"
	"strings"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForbiddenTypeInEnvironment struct {
	SnippetProvider
	envType string
	envs    []string
}

func NewForbiddenTypeInEnvironment(span core.Span, envType string, envs []string) *ForbiddenTypeInEnvironment {
	return &ForbiddenTypeInEnvironment{
		SnippetProvider: SnippetProvider{span: span},
		envType:         envType,
		envs:            envs,
	}
}
func (ftie *ForbiddenTypeInEnvironment) Message() string {
	return fmt.Sprintf("cannot have a %s in the following environments: %s", ftie.envType, strings.Join(ftie.envs, ", "))
}
func (ftie *ForbiddenTypeInEnvironment) Note() string { return "" }
func (ftie *ForbiddenTypeInEnvironment) Type() Type   { return Error }
func (ftie *ForbiddenTypeInEnvironment) ID() string   { return "hyb001W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
	SnippetProvider
	envType string
}

func NewInvalidEnvironmentType(span core.Span, envType string) *InvalidEnvironmentType {
	return &InvalidEnvironmentType{
		SnippetProvider: SnippetProvider{span: span},
		envType:         envType,
	}
}
func (iet *InvalidEnvironmentType) Message() string {
	return fmt.Sprintf("'%s' is not a valid environment type", iet.envType)
}
func (iet *InvalidEnvironmentType) Note() string {
	return "environment type can be 'Level', 'Mesh', 'Sound' or 'Shared'"
}
func (iet *InvalidEnvironmentType) Type() Type { return Error }
func (iet *InvalidEnvironmentType) ID() string { return "hyb002W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentRedaclaration struct {
	SnippetProvider
}

func NewEnvironmentRedaclaration(span core.Span) *EnvironmentRedaclaration {
	return &EnvironmentRedaclaration{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (er *EnvironmentRedaclaration) Message() string { return "cannot redeclare an environment" }
func (er *EnvironmentRedaclaration) Note() string    { return "" }
func (er *EnvironmentRedaclaration) Type() Type      { return Error }
func (er *EnvironmentRedaclaration) ID() string      { return "hyb003W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
	SnippetProvider
}

func NewExpectedEnvironment(span core.Span) *ExpectedEnvironment {
	return &ExpectedEnvironment{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ee *ExpectedEnvironment) Message() string { return "expected environment declaration" }
func (ee *ExpectedEnvironment) Note() string {
	return "the first declaration in any Hybroid file has to be an environment declaration"
}
func (ee *ExpectedEnvironment) Type() Type { return Error }
func (ee *ExpectedEnvironment) ID() string { return "hyb004W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateEnvironmentNames struct {
	SnippetProvider
	path1 string
	path2 string
}

func NewDuplicateEnvironmentNames(span core.Span, path1 string, path2 string) *DuplicateEnvironmentNames {
	return &DuplicateEnvironmentNames{
		SnippetProvider: SnippetProvider{span: span},
		path1:           path1,
		path2:           path2,
	}
}
func (den *DuplicateEnvironmentNames) Message() string {
	return fmt.Sprintf("duplicate environment names found between '%s' and '%s'", den.path1, den.path2)
}
func (den *DuplicateEnvironmentNames) Note() string { return "" }
func (den *DuplicateEnvironmentNames) Type() Type   { return Error }
func (den *DuplicateEnvironmentNames) ID() string   { return "hyb005W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidAccessValue struct {
	SnippetProvider
	valueType string
}

func NewInvalidAccessValue(span core.Span, valueType string) *InvalidAccessValue {
	return &InvalidAccessValue{
		SnippetProvider: SnippetProvider{span: span},
		valueType:       valueType,
	}
}
func (iav *InvalidAccessValue) Message() string {
	return fmt.Sprintf("value is of type '%s', so it cannot be accessed from", iav.valueType)
}
func (iav *InvalidAccessValue) Note() string {
	return "only lists, maps, classes, entities, structs and enums can be used to access values from"
}
func (iav *InvalidAccessValue) Type() Type { return Error }
func (iav *InvalidAccessValue) ID() string { return "hyb006W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type FieldAccessOnListOrMap struct {
	SnippetProvider
	field      string
	accessType string
}

func NewFieldAccessOnListOrMap(span core.Span, field string, accessType string) *FieldAccessOnListOrMap {
	return &FieldAccessOnListOrMap{
		SnippetProvider: SnippetProvider{span: span},
		field:           field,
		accessType:      accessType,
	}
}
func (faolom *FieldAccessOnListOrMap) Message() string {
	return fmt.Sprintf("cannot access field '%s' from the %s", faolom.field, faolom.accessType)
}
func (faolom *FieldAccessOnListOrMap) Note() string {
	return fmt.Sprintf("to access a value from a %s you use brackets, e.g. example[\"%s\"]", faolom.accessType, faolom.field)
}
func (faolom *FieldAccessOnListOrMap) Type() Type { return Error }
func (faolom *FieldAccessOnListOrMap) ID() string { return "hyb007W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MemberAccessOnNonListOrMap struct {
	SnippetProvider
	member     string
	accessType string
}

func NewMemberAccessOnNonListOrMap(span core.Span, member string, accessType string) *MemberAccessOnNonListOrMap {
	return &MemberAccessOnNonListOrMap{
		SnippetProvider: SnippetProvider{span: span},
		member:          member,
		accessType:      accessType,
	}
}
func (maonlom *MemberAccessOnNonListOrMap) Message() string {
	return fmt.Sprintf("cannot access member '[%s]' from the %s", maonlom.member, maonlom.accessType)
}
func (maonlom *MemberAccessOnNonListOrMap) Note() string {
	return "to access a value you use a dot and then an identifier, e.g. example.identifier"
}
func (maonlom *MemberAccessOnNonListOrMap) Type() Type { return Error }
func (maonlom *MemberAccessOnNonListOrMap) ID() string { return "hyb008W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidMemberIndex struct {
	SnippetProvider
	accessType string
	index      string
}

func NewInvalidMemberIndex(span core.Span, accessType string, index string) *InvalidMemberIndex {
	return &InvalidMemberIndex{
		SnippetProvider: SnippetProvider{span: span},
		accessType:      accessType,
		index:           index,
	}
}
func (imi *InvalidMemberIndex) Message() string {
	return fmt.Sprintf("'%s' is not of type number to be an index for the %s", imi.index, imi.accessType)
}
func (imi *InvalidMemberIndex) Note() string {
	return "for lists, an index (number) is used to access values, for maps, a key (text) is used"
}
func (imi *InvalidMemberIndex) Type() Type { return Error }
func (imi *InvalidMemberIndex) ID() string { return "hyb009W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidField struct {
	SnippetProvider
	accessType string
	fieldName  string
}

func NewInvalidField(span core.Span, accessType string, fieldName string) *InvalidField {
	return &InvalidField{
		SnippetProvider: SnippetProvider{span: span},
		accessType:      accessType,
		fieldName:       fieldName,
	}
}
func (_if *InvalidField) Message() string {
	return fmt.Sprintf("field '%s' does not belong to '%s'", _if.fieldName, _if.accessType)
}
func (_if *InvalidField) Note() string { return "" }
func (_if *InvalidField) Type() Type   { return Error }
func (_if *InvalidField) ID() string   { return "hyb010W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MixedMapOrListContents struct {
	SnippetProvider
	containerType string
	type1         string
	type2         string
}

func NewMixedMapOrListContents(span core.Span, containerType string, type1 string, type2 string) *MixedMapOrListContents {
	return &MixedMapOrListContents{
		SnippetProvider: SnippetProvider{span: span},
		containerType:   containerType,
		type1:           type1,
		type2:           type2,
	}
}
func (mmolc *MixedMapOrListContents) Message() string {
	return fmt.Sprintf("%s member is of type '%s', but the previous one was '%s'", mmolc.containerType, mmolc.type1, mmolc.type2)
}
func (mmolc *MixedMapOrListContents) Note() string { return "" }
func (mmolc *MixedMapOrListContents) Type() Type   { return Error }
func (mmolc *MixedMapOrListContents) ID() string   { return "hyb011W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCallerType struct {
	SnippetProvider
	callerType string
}

func NewInvalidCallerType(span core.Span, callerType string) *InvalidCallerType {
	return &InvalidCallerType{
		SnippetProvider: SnippetProvider{span: span},
		callerType:      callerType,
	}
}
func (ict *InvalidCallerType) Message() string {
	return fmt.Sprintf("cannot call value of type '%s' as a function", ict.callerType)
}
func (ict *InvalidCallerType) Note() string { return "" }
func (ict *InvalidCallerType) Type() Type   { return Error }
func (ict *InvalidCallerType) ID() string   { return "hyb012W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MethodOrFieldNotFound struct {
	SnippetProvider
	name string
}

func NewMethodOrFieldNotFound(span core.Span, name string) *MethodOrFieldNotFound {
	return &MethodOrFieldNotFound{
		SnippetProvider: SnippetProvider{span: span},
		name:            name,
	}
}
func (mofnf *MethodOrFieldNotFound) Message() string {
	return fmt.Sprintf("no method or field named '%s'", mofnf.name)
}
func (mofnf *MethodOrFieldNotFound) Note() string { return "" }
func (mofnf *MethodOrFieldNotFound) Type() Type   { return Error }
func (mofnf *MethodOrFieldNotFound) ID() string   { return "hyb013W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForeignLocalVariableAccess struct {
	SnippetProvider
	name string
}

func NewForeignLocalVariableAccess(span core.Span, name string) *ForeignLocalVariableAccess {
	return &ForeignLocalVariableAccess{
		SnippetProvider: SnippetProvider{span: span},
		name:            name,
	}
}
func (flva *ForeignLocalVariableAccess) Message() string {
	return fmt.Sprintf("cannot access local variable '%s' belonging to a different environment", flva.name)
}
func (flva *ForeignLocalVariableAccess) Note() string { return "" }
func (flva *ForeignLocalVariableAccess) Type() Type   { return Error }
func (flva *ForeignLocalVariableAccess) ID() string   { return "hyb014W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidArgumentType struct {
	SnippetProvider
	givenType    string
	expectedType string
}

func NewInvalidArgumentType(span core.Span, givenType string, expectedType string) *InvalidArgumentType {
	return &InvalidArgumentType{
		SnippetProvider: SnippetProvider{span: span},
		givenType:       givenType,
		expectedType:    expectedType,
	}
}
func (iat *InvalidArgumentType) Message() string {
	return fmt.Sprintf("argument was of type %s, but should be %s", iat.givenType, iat.expectedType)
}
func (iat *InvalidArgumentType) Note() string { return "" }
func (iat *InvalidArgumentType) Type() Type   { return Error }
func (iat *InvalidArgumentType) ID() string   { return "hyb015W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type PublicDeclarationInLocalScope struct {
	SnippetProvider
}

func NewPublicDeclarationInLocalScope(span core.Span) *PublicDeclarationInLocalScope {
	return &PublicDeclarationInLocalScope{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (pdils *PublicDeclarationInLocalScope) Message() string {
	return "cannot have a public declaration that is in a local scope"
}
func (pdils *PublicDeclarationInLocalScope) Note() string { return "" }
func (pdils *PublicDeclarationInLocalScope) Type() Type   { return Error }
func (pdils *PublicDeclarationInLocalScope) ID() string   { return "hyb016W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type Redeclaration struct {
	SnippetProvider
	varName  string
	declType string
}

func NewRedeclaration(span core.Span, varName string, declType string) *Redeclaration {
	return &Redeclaration{
		SnippetProvider: SnippetProvider{span: span},
		varName:         varName,
		declType:        declType,
	}
}
func (r *Redeclaration) Message() string {
	return fmt.Sprintf("a %s named '%s' already exists", r.declType, r.varName)
}
func (r *Redeclaration) Note() string { return "" }
func (r *Redeclaration) Type() Type   { return Error }
func (r *Redeclaration) ID() string   { return "hyb017W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type NoValueGivenForConstant struct {
	SnippetProvider
}

func NewNoValueGivenForConstant(span core.Span) *NoValueGivenForConstant {
	return &NoValueGivenForConstant{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (nvgfc *NoValueGivenForConstant) Message() string {
	return "constant must be declared with a value"
}
func (nvgfc *NoValueGivenForConstant) Note() string { return "" }
func (nvgfc *NoValueGivenForConstant) Type() Type   { return Error }
func (nvgfc *NoValueGivenForConstant) ID() string   { return "hyb018W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TooFewElementsGiven struct {
	SnippetProvider
	requiredAmount int
	elem           string
	context        string
}

func NewTooFewElementsGiven(span core.Span, requiredAmount int, elem string, context string) *TooFewElementsGiven {
	return &TooFewElementsGiven{
		SnippetProvider: SnippetProvider{span: span},
		requiredAmount:  requiredAmount,
		elem:            elem,
		context:         context,
	}
}
func (tfeg *TooFewElementsGiven) Message() string {
	return fmt.Sprintf("%d more %s(s) required %s", tfeg.requiredAmount, tfeg.elem, tfeg.context)
}
func (tfeg *TooFewElementsGiven) Note() string { return "" }
func (tfeg *TooFewElementsGiven) Type() Type   { return Error }
func (tfeg *TooFewElementsGiven) ID() string   { return "hyb019W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TooManyElementsGiven struct {
	SnippetProvider
	extraAmount int
	elem        string
	context     string
}

func NewTooManyElementsGiven(span core.Span, extraAmount int, elem string, context string) *TooManyElementsGiven {
	return &TooManyElementsGiven{
		SnippetProvider: SnippetProvider{span: span},
		extraAmount:     extraAmount,
		elem:            elem,
		context:         context,
	}
}
func (tmeg *TooManyElementsGiven) Message() string {
	return fmt.Sprintf("%d less %s(s) required %s", tmeg.extraAmount, tmeg.elem, tmeg.context)
}
func (tmeg *TooManyElementsGiven) Note() string { return "" }
func (tmeg *TooManyElementsGiven) Type() Type   { return Error }
func (tmeg *TooManyElementsGiven) ID() string   { return "hyb020W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeRequiredInDeclaration struct {
	SnippetProvider
	context string
}

func NewExplicitTypeRequiredInDeclaration(span core.Span, context string) *ExplicitTypeRequiredInDeclaration {
	return &ExplicitTypeRequiredInDeclaration{
		SnippetProvider: SnippetProvider{span: span},
		context:         context,
	}
}
func (etrid *ExplicitTypeRequiredInDeclaration) Message() string {
	return fmt.Sprintf("an explicit type is required %s", etrid.context)
}
func (etrid *ExplicitTypeRequiredInDeclaration) Note() string { return "" }
func (etrid *ExplicitTypeRequiredInDeclaration) Type() Type   { return Error }
func (etrid *ExplicitTypeRequiredInDeclaration) ID() string   { return "hyb021W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeMismatch struct {
	SnippetProvider
	explicitType string
	valueType    string
}

func NewExplicitTypeMismatch(span core.Span, explicitType string, valueType string) *ExplicitTypeMismatch {
	return &ExplicitTypeMismatch{
		SnippetProvider: SnippetProvider{span: span},
		explicitType:    explicitType,
		valueType:       valueType,
	}
}
func (etm *ExplicitTypeMismatch) Message() string {
	return fmt.Sprintf("variable was given explicit type '%s', but its value is a '%s'", etm.explicitType, etm.valueType)
}
func (etm *ExplicitTypeMismatch) Note() string { return "" }
func (etm *ExplicitTypeMismatch) Type() Type   { return Error }
func (etm *ExplicitTypeMismatch) ID() string   { return "hyb022W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExplicitTypeNotAllowed struct {
	SnippetProvider
	explicitType string
}

func NewExplicitTypeNotAllowed(span core.Span, explicitType string) *ExplicitTypeNotAllowed {
	return &ExplicitTypeNotAllowed{
		SnippetProvider: SnippetProvider{span: span},
		explicitType:    explicitType,
	}
}
func (etna *ExplicitTypeNotAllowed) Message() string {
	return fmt.Sprintf("cannot create a default value from the explicit type '%s'", etna.explicitType)
}
func (etna *ExplicitTypeNotAllowed) Note() string {
	return "some types don't have default values, like entities and classes"
}
func (etna *ExplicitTypeNotAllowed) Type() Type { return Error }
func (etna *ExplicitTypeNotAllowed) ID() string { return "hyb023W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ImportCycle struct {
	SnippetProvider
	paths []string
}

func NewImportCycle(span core.Span, paths []string) *ImportCycle {
	return &ImportCycle{
		SnippetProvider: SnippetProvider{span: span},
		paths:           paths,
	}
}
func (ic *ImportCycle) Message() string {
	return fmt.Sprintf("import cycle detected: %s", strings.Join(ic.paths, " -> "))
}
func (ic *ImportCycle) Note() string { return "" }
func (ic *ImportCycle) Type() Type   { return Error }
func (ic *ImportCycle) ID() string   { return "hyb024W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UndeclaredVariableAccess struct {
	SnippetProvider
	variable string
	context  string
}

func NewUndeclaredVariableAccess(span core.Span, variable string, context string) *UndeclaredVariableAccess {
	return &UndeclaredVariableAccess{
		SnippetProvider: SnippetProvider{span: span},
		variable:        variable,
		context:         context,
	}
}
func (uva *UndeclaredVariableAccess) Message() string {
	return fmt.Sprintf("'%s' is not a declared variable %s", uva.variable, uva.context)
}
func (uva *UndeclaredVariableAccess) Note() string { return "" }
func (uva *UndeclaredVariableAccess) Type() Type   { return Error }
func (uva *UndeclaredVariableAccess) ID() string   { return "hyb025W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ConstValueAssignment struct {
	SnippetProvider
}

func NewConstValueAssignment(span core.Span) *ConstValueAssignment {
	return &ConstValueAssignment{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (cva *ConstValueAssignment) Message() string { return "cannot modify a constant value" }
func (cva *ConstValueAssignment) Note() string    { return "" }
func (cva *ConstValueAssignment) Type() Type      { return Error }
func (cva *ConstValueAssignment) ID() string      { return "hyb026W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type AssignmentTypeMismatch struct {
	SnippetProvider
	varType string
	valType string
}

func NewAssignmentTypeMismatch(span core.Span, varType string, valType string) *AssignmentTypeMismatch {
	return &AssignmentTypeMismatch{
		SnippetProvider: SnippetProvider{span: span},
		varType:         varType,
		valType:         valType,
	}
}
func (atm *AssignmentTypeMismatch) Message() string {
	return fmt.Sprintf("variable is of type '%s', but a value of '%s' was assigned to it", atm.varType, atm.valType)
}
func (atm *AssignmentTypeMismatch) Note() string { return "" }
func (atm *AssignmentTypeMismatch) Type() Type   { return Error }
func (atm *AssignmentTypeMismatch) ID() string   { return "hyb027W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidTypeInCompoundAssignment struct {
	SnippetProvider
	valueType string
}

func NewInvalidTypeInCompoundAssignment(span core.Span, valueType string) *InvalidTypeInCompoundAssignment {
	return &InvalidTypeInCompoundAssignment{
		SnippetProvider: SnippetProvider{span: span},
		valueType:       valueType,
	}
}
func (itica *InvalidTypeInCompoundAssignment) Message() string {
	return fmt.Sprintf("the type '%s' is not allowed in compound assignment", itica.valueType)
}
func (itica *InvalidTypeInCompoundAssignment) Note() string {
	return "only numerical types are allowed, like number or fixed"
}
func (itica *InvalidTypeInCompoundAssignment) Type() Type { return Error }
func (itica *InvalidTypeInCompoundAssignment) ID() string { return "hyb028W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidUseOfSelf struct {
	SnippetProvider
}

func NewInvalidUseOfSelf(span core.Span) *InvalidUseOfSelf {
	return &InvalidUseOfSelf{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (iuos *InvalidUseOfSelf) Message() string { return "cannot use self outside of class or entity" }
func (iuos *InvalidUseOfSelf) Note() string {
	return "you're also not allowed to use self inside anonymous functions of class/entity fields"
}
func (iuos *InvalidUseOfSelf) Type() Type { return Error }
func (iuos *InvalidUseOfSelf) ID() string { return "hyb029W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnreachableCode struct {
	SnippetProvider
}

func NewUnreachableCode(span core.Span) *UnreachableCode {
	return &UnreachableCode{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (uc *UnreachableCode) Message() string { return "unreachable code detected" }
func (uc *UnreachableCode) Note() string    { return "" }
func (uc *UnreachableCode) Type() Type      { return Warning }
func (uc *UnreachableCode) ID() string      { return "hyb030W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidUseOfExitStmt struct {
	SnippetProvider
	exitNode string
	context  string
}

func NewInvalidUseOfExitStmt(span core.Span, exitNode string, context string) *InvalidUseOfExitStmt {
	return &InvalidUseOfExitStmt{
		SnippetProvider: SnippetProvider{span: span},
		exitNode:        exitNode,
		context:         context,
	}
}
func (iuoes *InvalidUseOfExitStmt) Message() string {
	return fmt.Sprintf("cannot use '%s' outside of %s", iuoes.exitNode, iuoes.context)
}
func (iuoes *InvalidUseOfExitStmt) Note() string { return "" }
func (iuoes *InvalidUseOfExitStmt) Type() Type   { return Error }
func (iuoes *InvalidUseOfExitStmt) ID() string   { return "hyb031W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypeMismatch struct {
	SnippetProvider
	type1   string
	type2   string
	context string
}

func NewTypeMismatch(span core.Span, type1 string, type2 string, context string) *TypeMismatch {
	return &TypeMismatch{
		SnippetProvider: SnippetProvider{span: span},
		type1:           type1,
		type2:           type2,
		context:         context,
	}
}
func (tm *TypeMismatch) Message() string {
	return fmt.Sprintf("expected %s, got '%s' %s", tm.type1, tm.type2, tm.context)
}
func (tm *TypeMismatch) Note() string { return "" }
func (tm *TypeMismatch) Type() Type   { return Error }
func (tm *TypeMismatch) ID() string   { return "hyb032W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidStmtInLocalBlock struct {
	SnippetProvider
	stmt string
}

func NewInvalidStmtInLocalBlock(span core.Span, stmt string) *InvalidStmtInLocalBlock {
	return &InvalidStmtInLocalBlock{
		SnippetProvider: SnippetProvider{span: span},
		stmt:            stmt,
	}
}
func (isilb *InvalidStmtInLocalBlock) Message() string {
	return fmt.Sprintf("%s must be in the global scope", isilb.stmt)
}
func (isilb *InvalidStmtInLocalBlock) Note() string { return "" }
func (isilb *InvalidStmtInLocalBlock) Type() Type   { return Error }
func (isilb *InvalidStmtInLocalBlock) ID() string   { return "hyb033W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnallowedLibraryUse struct {
	SnippetProvider
	library       string
	unallowedEnvs string
}

func NewUnallowedLibraryUse(span core.Span, library string, unallowedEnvs string) *UnallowedLibraryUse {
	return &UnallowedLibraryUse{
		SnippetProvider: SnippetProvider{span: span},
		library:         library,
		unallowedEnvs:   unallowedEnvs,
	}
}
func (ulu *UnallowedLibraryUse) Message() string {
	return fmt.Sprintf("cannot use the %s library in a %s environment", ulu.library, ulu.unallowedEnvs)
}
func (ulu *UnallowedLibraryUse) Note() string { return "" }
func (ulu *UnallowedLibraryUse) Type() Type   { return Error }
func (ulu *UnallowedLibraryUse) ID() string   { return "hyb034W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentAccess struct {
	SnippetProvider
	envName string
}

func NewInvalidEnvironmentAccess(span core.Span, envName string) *InvalidEnvironmentAccess {
	return &InvalidEnvironmentAccess{
		SnippetProvider: SnippetProvider{span: span},
		envName:         envName,
	}
}
func (iea *InvalidEnvironmentAccess) Message() string {
	return fmt.Sprintf("environment named '%s' does not exist", iea.envName)
}
func (iea *InvalidEnvironmentAccess) Note() string { return "" }
func (iea *InvalidEnvironmentAccess) Type() Type   { return Error }
func (iea *InvalidEnvironmentAccess) ID() string   { return "hyb035W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentReuse struct {
	SnippetProvider
	envName string
}

func NewEnvironmentReuse(span core.Span, envName string) *EnvironmentReuse {
	return &EnvironmentReuse{
		SnippetProvider: SnippetProvider{span: span},
		envName:         envName,
	}
}
func (er *EnvironmentReuse) Message() string {
	return fmt.Sprintf("environment named '%s' is already imported through use statement", er.envName)
}
func (er *EnvironmentReuse) Note() string { return "" }
func (er *EnvironmentReuse) Type() Type   { return Error }
func (er *EnvironmentReuse) ID() string   { return "hyb036W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidIteratorType struct {
	SnippetProvider
	value string
}

func NewInvalidIteratorType(span core.Span, value string) *InvalidIteratorType {
	return &InvalidIteratorType{
		SnippetProvider: SnippetProvider{span: span},
		value:           value,
	}
}
func (iit *InvalidIteratorType) Message() string {
	return fmt.Sprintf("a for loop iterator must be a map or a list (found: '%s')", iit.value)
}
func (iit *InvalidIteratorType) Note() string { return "" }
func (iit *InvalidIteratorType) Type() Type   { return Error }
func (iit *InvalidIteratorType) ID() string   { return "hyb037W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnnecessaryEmptyIdentifier struct {
	SnippetProvider
	context string
}

func NewUnnecessaryEmptyIdentifier(span core.Span, context string) *UnnecessaryEmptyIdentifier {
	return &UnnecessaryEmptyIdentifier{
		SnippetProvider: SnippetProvider{span: span},
		context:         context,
	}
}
func (uei *UnnecessaryEmptyIdentifier) Message() string {
	return fmt.Sprintf("unnecessary use of empty identifier ('_') %s", uei.context)
}
func (uei *UnnecessaryEmptyIdentifier) Note() string { return "" }
func (uei *UnnecessaryEmptyIdentifier) Type() Type   { return Warning }
func (uei *UnnecessaryEmptyIdentifier) ID() string   { return "hyb038W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentUsesItself struct {
	SnippetProvider
}

func NewEnvironmentUsesItself(span core.Span) *EnvironmentUsesItself {
	return &EnvironmentUsesItself{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eui *EnvironmentUsesItself) Message() string { return "an environment cannot 'use' itself" }
func (eui *EnvironmentUsesItself) Note() string    { return "" }
func (eui *EnvironmentUsesItself) Type() Type      { return Error }
func (eui *EnvironmentUsesItself) ID() string      { return "hyb039W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EntityConversionWithOrCondition struct {
	SnippetProvider
}

func NewEntityConversionWithOrCondition(span core.Span) *EntityConversionWithOrCondition {
	return &EntityConversionWithOrCondition{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ecwoc *EntityConversionWithOrCondition) Message() string {
	return "cannot convert an entity with an 'or' condition"
}
func (ecwoc *EntityConversionWithOrCondition) Note() string { return "" }
func (ecwoc *EntityConversionWithOrCondition) Type() Type   { return Error }
func (ecwoc *EntityConversionWithOrCondition) ID() string   { return "hyb040W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCondition struct {
	SnippetProvider
	context string
}

func NewInvalidCondition(span core.Span, context string) *InvalidCondition {
	return &InvalidCondition{
		SnippetProvider: SnippetProvider{span: span},
		context:         context,
	}
}
func (ic *InvalidCondition) Message() string { return fmt.Sprintf("invalid condition %s", ic.context) }
func (ic *InvalidCondition) Note() string {
	return "conditions always have to evaluate to either true or false"
}
func (ic *InvalidCondition) Type() Type { return Error }
func (ic *InvalidCondition) ID() string { return "hyb041W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidRepeatIterator struct {
	SnippetProvider
	valueType string
}

func NewInvalidRepeatIterator(span core.Span, valueType string) *InvalidRepeatIterator {
	return &InvalidRepeatIterator{
		SnippetProvider: SnippetProvider{span: span},
		valueType:       valueType,
	}
}
func (iri *InvalidRepeatIterator) Message() string {
	return fmt.Sprintf("invalid repeat iterator of type '%s'", iri.valueType)
}
func (iri *InvalidRepeatIterator) Note() string { return "repeat iterator must be a numerical type" }
func (iri *InvalidRepeatIterator) Type() Type   { return Error }
func (iri *InvalidRepeatIterator) ID() string   { return "hyb042W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InconsistentRepeatTypes struct {
	SnippetProvider
	from     string
	skip     string
	iterator string
}

func NewInconsistentRepeatTypes(span core.Span, from string, skip string, iterator string) *InconsistentRepeatTypes {
	return &InconsistentRepeatTypes{
		SnippetProvider: SnippetProvider{span: span},
		from:            from,
		skip:            skip,
		iterator:        iterator,
	}
}
func (irt *InconsistentRepeatTypes) Message() string {
	return fmt.Sprintf("repeat types are inconsistent (from:'%s', by:'%s', to:'%s')", irt.from, irt.skip, irt.iterator)
}
func (irt *InconsistentRepeatTypes) Note() string { return "" }
func (irt *InconsistentRepeatTypes) Type() Type   { return Error }
func (irt *InconsistentRepeatTypes) ID() string   { return "hyb043W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type OfficialEntityConversion struct {
	SnippetProvider
}

func NewOfficialEntityConversion(span core.Span) *OfficialEntityConversion {
	return &OfficialEntityConversion{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (oec *OfficialEntityConversion) Message() string {
	return "conversion of an official entity to a hybroid entity is not possible"
}
func (oec *OfficialEntityConversion) Note() string { return "" }
func (oec *OfficialEntityConversion) Type() Type   { return Error }
func (oec *OfficialEntityConversion) ID() string   { return "hyb044W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironment struct {
	SnippetProvider
}

func NewInvalidEnvironment(span core.Span) *InvalidEnvironment {
	return &InvalidEnvironment{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ie *InvalidEnvironment) Message() string { return "there is no environment with that path" }
func (ie *InvalidEnvironment) Note() string    { return "" }
func (ie *InvalidEnvironment) Type() Type      { return Error }
func (ie *InvalidEnvironment) ID() string      { return "hyb045W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentAccessAmbiguity struct {
	SnippetProvider
	envs    []string
	context string
}

func NewEnvironmentAccessAmbiguity(span core.Span, envs []string, context string) *EnvironmentAccessAmbiguity {
	return &EnvironmentAccessAmbiguity{
		SnippetProvider: SnippetProvider{span: span},
		envs:            envs,
		context:         context,
	}
}
func (eaa *EnvironmentAccessAmbiguity) Message() string {
	return fmt.Sprintf("the type '%s' can be found on multiple environments: %s", eaa.context, strings.Join(eaa.envs, ", "))
}
func (eaa *EnvironmentAccessAmbiguity) Note() string { return "" }
func (eaa *EnvironmentAccessAmbiguity) Type() Type   { return Error }
func (eaa *EnvironmentAccessAmbiguity) ID() string   { return "hyb046W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type NotAllCodePathsExit struct {
	SnippetProvider
	exitType string
}

func NewNotAllCodePathsExit(span core.Span, exitType string) *NotAllCodePathsExit {
	return &NotAllCodePathsExit{
		SnippetProvider: SnippetProvider{span: span},
		exitType:        exitType,
	}
}
func (nacpe *NotAllCodePathsExit) Message() string {
	return fmt.Sprintf("not all code paths %s", nacpe.exitType)
}
func (nacpe *NotAllCodePathsExit) Note() string { return "" }
func (nacpe *NotAllCodePathsExit) Type() Type   { return Error }
func (nacpe *NotAllCodePathsExit) ID() string   { return "hyb047W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InsufficientCases struct {
	SnippetProvider
}

func NewInsufficientCases(span core.Span) *InsufficientCases {
	return &InsufficientCases{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ic *InsufficientCases) Message() string {
	return "match statement must have at least 1 non-default case"
}
func (ic *InsufficientCases) Note() string { return "" }
func (ic *InsufficientCases) Type() Type   { return Error }
func (ic *InsufficientCases) ID() string   { return "hyb048W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DefaultCaseMissing struct {
	SnippetProvider
}

func NewDefaultCaseMissing(span core.Span) *DefaultCaseMissing {
	return &DefaultCaseMissing{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (dcm *DefaultCaseMissing) Message() string { return "match expression must have a default case" }
func (dcm *DefaultCaseMissing) Note() string    { return "default cases start with 'else'" }
func (dcm *DefaultCaseMissing) Type() Type      { return Error }
func (dcm *DefaultCaseMissing) ID() string      { return "hyb049W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCaseType struct {
	SnippetProvider
	matchValueType string
	caseValueType  string
}

func NewInvalidCaseType(span core.Span, matchValueType string, caseValueType string) *InvalidCaseType {
	return &InvalidCaseType{
		SnippetProvider: SnippetProvider{span: span},
		matchValueType:  matchValueType,
		caseValueType:   caseValueType,
	}
}
func (ict *InvalidCaseType) Message() string {
	return fmt.Sprintf("match value is of type '%s', but case value is of type '%s'", ict.matchValueType, ict.caseValueType)
}
func (ict *InvalidCaseType) Note() string { return "" }
func (ict *InvalidCaseType) Type() Type   { return Error }
func (ict *InvalidCaseType) ID() string   { return "hyb050W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type LiteralCondition struct {
	SnippetProvider
	conditionValue string
}

func NewLiteralCondition(span core.Span, conditionValue string) *LiteralCondition {
	return &LiteralCondition{
		SnippetProvider: SnippetProvider{span: span},
		conditionValue:  conditionValue,
	}
}
func (lc *LiteralCondition) Message() string {
	return fmt.Sprintf("condition is always %s", lc.conditionValue)
}
func (lc *LiteralCondition) Note() string { return "" }
func (lc *LiteralCondition) Type() Type   { return Warning }
func (lc *LiteralCondition) ID() string   { return "hyb051W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypesMismatch struct {
	SnippetProvider
	value1 string
	type1  string
	value2 string
	type2  string
}

func NewTypesMismatch(span core.Span, value1 string, type1 string, value2 string, type2 string) *TypesMismatch {
	return &TypesMismatch{
		SnippetProvider: SnippetProvider{span: span},
		value1:          value1,
		type1:           type1,
		value2:          value2,
		type2:           type2,
	}
}
func (tm *TypesMismatch) Message() string {
	return fmt.Sprintf("%s is of type '%s', but %s is of type '%s'", tm.value1, tm.type1, tm.value2, tm.type2)
}
func (tm *TypesMismatch) Note() string { return "" }
func (tm *TypesMismatch) Type() Type   { return Error }
func (tm *TypesMismatch) ID() string   { return "hyb052W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingConstructor struct {
	SnippetProvider
	constructorType string
	context         string
}

func NewMissingConstructor(span core.Span, constructorType string, context string) *MissingConstructor {
	return &MissingConstructor{
		SnippetProvider: SnippetProvider{span: span},
		constructorType: constructorType,
		context:         context,
	}
}
func (mc *MissingConstructor) Message() string {
	return fmt.Sprintf("missing '%s' constructor %s", mc.constructorType, mc.context)
}
func (mc *MissingConstructor) Note() string { return "" }
func (mc *MissingConstructor) Type() Type   { return Error }
func (mc *MissingConstructor) ID() string   { return "hyb053W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingDestroy struct {
	SnippetProvider
}

func NewMissingDestroy(span core.Span) *MissingDestroy {
	return &MissingDestroy{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (md *MissingDestroy) Message() string {
	return "missing 'destroy' destructor in entity declaration"
}
func (md *MissingDestroy) Note() string { return "" }
func (md *MissingDestroy) Type() Type   { return Error }
func (md *MissingDestroy) ID() string   { return "hyb054W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UninitializedFieldInConstructor struct {
	SnippetProvider
	varName string
	context string
}

func NewUninitializedFieldInConstructor(span core.Span, varName string, context string) *UninitializedFieldInConstructor {
	return &UninitializedFieldInConstructor{
		SnippetProvider: SnippetProvider{span: span},
		varName:         varName,
		context:         context,
	}
}
func (ufic *UninitializedFieldInConstructor) Message() string {
	return fmt.Sprintf("variable '%s' was not initialized in the constructor %s", ufic.varName, ufic.context)
}
func (ufic *UninitializedFieldInConstructor) Note() string { return "" }
func (ufic *UninitializedFieldInConstructor) Type() Type   { return Error }
func (ufic *UninitializedFieldInConstructor) ID() string   { return "hyb055W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type TypeRedeclaration struct {
	SnippetProvider
	typeName string
}

func NewTypeRedeclaration(span core.Span, typeName string) *TypeRedeclaration {
	return &TypeRedeclaration{
		SnippetProvider: SnippetProvider{span: span},
		typeName:        typeName,
	}
}
func (tr *TypeRedeclaration) Message() string {
	return fmt.Sprintf("type '%s' already exists", tr.typeName)
}
func (tr *TypeRedeclaration) Note() string { return "" }
func (tr *TypeRedeclaration) Type() Type   { return Error }
func (tr *TypeRedeclaration) ID() string   { return "hyb056W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCallAsArgument struct {
	SnippetProvider
}

func NewInvalidCallAsArgument(span core.Span) *InvalidCallAsArgument {
	return &InvalidCallAsArgument{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (icaa *InvalidCallAsArgument) Message() string {
	return "cannot have a call that returns more than 1 value as an argument"
}
func (icaa *InvalidCallAsArgument) Note() string { return "" }
func (icaa *InvalidCallAsArgument) Type() Type   { return Error }
func (icaa *InvalidCallAsArgument) ID() string   { return "hyb057W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneVariadicParameter struct {
	SnippetProvider
}

func NewMoreThanOneVariadicParameter(span core.Span) *MoreThanOneVariadicParameter {
	return &MoreThanOneVariadicParameter{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (mtovp *MoreThanOneVariadicParameter) Message() string {
	return "cannot have more than one variadic function parameter"
}
func (mtovp *MoreThanOneVariadicParameter) Note() string { return "" }
func (mtovp *MoreThanOneVariadicParameter) Type() Type   { return Error }
func (mtovp *MoreThanOneVariadicParameter) ID() string   { return "hyb058W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type VariadicParameterNotAtEnd struct {
	SnippetProvider
}

func NewVariadicParameterNotAtEnd(span core.Span) *VariadicParameterNotAtEnd {
	return &VariadicParameterNotAtEnd{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (vpnae *VariadicParameterNotAtEnd) Message() string {
	return "variadic parameters must be at the end of the function parameters"
}
func (vpnae *VariadicParameterNotAtEnd) Note() string { return "" }
func (vpnae *VariadicParameterNotAtEnd) Type() Type   { return Error }
func (vpnae *VariadicParameterNotAtEnd) ID() string   { return "hyb059W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateElement struct {
	SnippetProvider
	element  string
	elemName string
}

func NewDuplicateElement(span core.Span, element string, elemName string) *DuplicateElement {
	return &DuplicateElement{
		SnippetProvider: SnippetProvider{span: span},
		element:         element,
		elemName:        elemName,
	}
}
func (de *DuplicateElement) Message() string {
	return fmt.Sprintf("the %s '%s' already exists", de.element, de.elemName)
}
func (de *DuplicateElement) Note() string { return "" }
func (de *DuplicateElement) Type() Type   { return Error }
func (de *DuplicateElement) ID() string   { return "hyb060W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEntityFunctionSignature struct {
	SnippetProvider
	got            string
	expected       string
	entityFuncType string
}

func NewInvalidEntityFunctionSignature(span core.Span, got string, expected string, entityFuncType string) *InvalidEntityFunctionSignature {
	return &InvalidEntityFunctionSignature{
		SnippetProvider: SnippetProvider{span: span},
		got:             got,
		expected:        expected,
		entityFuncType:  entityFuncType,
	}
}
func (iefs *InvalidEntityFunctionSignature) Message() string {
	return fmt.Sprintf("expected '%s' for %s, got '%s'", iefs.expected, iefs.entityFuncType, iefs.got)
}
func (iefs *InvalidEntityFunctionSignature) Note() string { return "" }
func (iefs *InvalidEntityFunctionSignature) Type() Type   { return Error }
func (iefs *InvalidEntityFunctionSignature) ID() string   { return "hyb061W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidSpawnerParameters struct {
	SnippetProvider
}

func NewInvalidSpawnerParameters(span core.Span) *InvalidSpawnerParameters {
	return &InvalidSpawnerParameters{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (isp *InvalidSpawnerParameters) Message() string {
	return "the first two parameters of the spawner must be fixedpoints (x and y)"
}
func (isp *InvalidSpawnerParameters) Note() string { return "" }
func (isp *InvalidSpawnerParameters) Type() Type   { return Error }
func (isp *InvalidSpawnerParameters) ID() string   { return "hyb062W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidPewpewVariable struct {
	SnippetProvider
	pewpewVar string
	varType   string
}

func NewInvalidPewpewVariable(span core.Span, pewpewVar string, varType string) *InvalidPewpewVariable {
	return &InvalidPewpewVariable{
		SnippetProvider: SnippetProvider{span: span},
		pewpewVar:       pewpewVar,
		varType:         varType,
	}
}
func (ipv *InvalidPewpewVariable) Message() string {
	return fmt.Sprintf("'%s' variable should be global and of type 'list<%s>'", ipv.pewpewVar, ipv.varType)
}
func (ipv *InvalidPewpewVariable) Note() string { return "" }
func (ipv *InvalidPewpewVariable) Type() Type   { return Error }
func (ipv *InvalidPewpewVariable) ID() string   { return "hyb063W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingPewpewVariable struct {
	SnippetProvider
	pewpewVar string
	envType   string
}

func NewMissingPewpewVariable(span core.Span, pewpewVar string, envType string) *MissingPewpewVariable {
	return &MissingPewpewVariable{
		SnippetProvider: SnippetProvider{span: span},
		pewpewVar:       pewpewVar,
		envType:         envType,
	}
}
func (mpv *MissingPewpewVariable) Message() string {
	return fmt.Sprintf("A %s environment must have a '%s' variable", mpv.envType, mpv.pewpewVar)
}
func (mpv *MissingPewpewVariable) Note() string { return "" }
func (mpv *MissingPewpewVariable) Type() Type   { return Error }
func (mpv *MissingPewpewVariable) ID() string   { return "hyb064W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnallowedEnvironmentAccess struct {
	SnippetProvider
	unallowed string
	from      string
}

func NewUnallowedEnvironmentAccess(span core.Span, unallowed string, from string) *UnallowedEnvironmentAccess {
	return &UnallowedEnvironmentAccess{
		SnippetProvider: SnippetProvider{span: span},
		unallowed:       unallowed,
		from:            from,
	}
}
func (uea *UnallowedEnvironmentAccess) Message() string {
	return fmt.Sprintf("cannot access a %s environment from a %s environment", uea.unallowed, uea.from)
}
func (uea *UnallowedEnvironmentAccess) Note() string { return "" }
func (uea *UnallowedEnvironmentAccess) Type() Type   { return Error }
func (uea *UnallowedEnvironmentAccess) ID() string   { return "hyb065W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidDefaultCasePlacement struct {
	SnippetProvider
	context string
}

func NewInvalidDefaultCasePlacement(span core.Span, context string) *InvalidDefaultCasePlacement {
	return &InvalidDefaultCasePlacement{
		SnippetProvider: SnippetProvider{span: span},
		context:         context,
	}
}
func (idcp *InvalidDefaultCasePlacement) Message() string {
	return fmt.Sprintf("the default case must always be at the end %s", idcp.context)
}
func (idcp *InvalidDefaultCasePlacement) Note() string { return "" }
func (idcp *InvalidDefaultCasePlacement) Type() Type   { return Error }
func (idcp *InvalidDefaultCasePlacement) ID() string   { return "hyb066W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidType struct {
	SnippetProvider
	varType string
	context string
}

func NewInvalidType(span core.Span, varType string, context string) *InvalidType {
	return &InvalidType{
		SnippetProvider: SnippetProvider{span: span},
		varType:         varType,
		context:         context,
	}
}
func (it *InvalidType) Message() string {
	return fmt.Sprintf("cannot have a type '%s' %s", it.varType, it.context)
}
func (it *InvalidType) Note() string { return "" }
func (it *InvalidType) Type() Type   { return Error }
func (it *InvalidType) ID() string   { return "hyb067W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ListIndexOutOfBounds struct {
	SnippetProvider
}

func NewListIndexOutOfBounds(span core.Span) *ListIndexOutOfBounds {
	return &ListIndexOutOfBounds{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (lioob *ListIndexOutOfBounds) Message() string {
	return "list index is 0 or less, but it must be 1 or more"
}
func (lioob *ListIndexOutOfBounds) Note() string { return "" }
func (lioob *ListIndexOutOfBounds) Type() Type   { return Error }
func (lioob *ListIndexOutOfBounds) ID() string   { return "hyb068W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidListIndex struct {
	SnippetProvider
}

func NewInvalidListIndex(span core.Span) *InvalidListIndex {
	return &InvalidListIndex{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ili *InvalidListIndex) Message() string { return "a list index must be a whole number" }
func (ili *InvalidListIndex) Note() string    { return "" }
func (ili *InvalidListIndex) Type() Type      { return Error }
func (ili *InvalidListIndex) ID() string      { return "hyb069W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingGenericArgument struct {
	SnippetProvider
	genericType string
}

func NewMissingGenericArgument(span core.Span, genericType string) *MissingGenericArgument {
	return &MissingGenericArgument{
		SnippetProvider: SnippetProvider{span: span},
		genericType:     genericType,
	}
}
func (mga *MissingGenericArgument) Message() string {
	return fmt.Sprintf("generic type '%s' could not be inferred", mga.genericType)
}
func (mga *MissingGenericArgument) Note() string { return "" }
func (mga *MissingGenericArgument) Type() Type   { return Error }
func (mga *MissingGenericArgument) ID() string   { return "hyb070W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidAssignment struct {
	SnippetProvider
}

func NewInvalidAssignment(span core.Span) *InvalidAssignment {
	return &InvalidAssignment{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ia *InvalidAssignment) Message() string { return "left value was not a variable" }
func (ia *InvalidAssignment) Note() string    { return "" }
func (ia *InvalidAssignment) Type() Type      { return Error }
func (ia *InvalidAssignment) ID() string      { return "hyb071W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ConflictingVariableNameWithType struct {
	SnippetProvider
	name string
}

func NewConflictingVariableNameWithType(span core.Span, name string) *ConflictingVariableNameWithType {
	return &ConflictingVariableNameWithType{
		SnippetProvider: SnippetProvider{span: span},
		name:            name,
	}
}
func (cvnwt *ConflictingVariableNameWithType) Message() string {
	return fmt.Sprintf("variable name conflicts with type '%s'", cvnwt.name)
}
func (cvnwt *ConflictingVariableNameWithType) Note() string { return "" }
func (cvnwt *ConflictingVariableNameWithType) Type() Type   { return Error }
func (cvnwt *ConflictingVariableNameWithType) ID() string   { return "hyb072W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnusedElement struct {
	SnippetProvider
	elem string
}

func NewUnusedElement(span core.Span, elem string) *UnusedElement {
	return &UnusedElement{
		SnippetProvider: SnippetProvider{span: span},
		elem:            elem,
	}
}
func (ue *UnusedElement) Message() string { return fmt.Sprintf("%s is not used", ue.elem) }
func (ue *UnusedElement) Note() string    { return "" }
func (ue *UnusedElement) Type() Type      { return Warning }
func (ue *UnusedElement) ID() string      { return "hyb073W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyIdentifierOnSpawnParameters struct {
	SnippetProvider
}

func NewEmptyIdentifierOnSpawnParameters(span core.Span) *EmptyIdentifierOnSpawnParameters {
	return &EmptyIdentifierOnSpawnParameters{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eiosp *EmptyIdentifierOnSpawnParameters) Message() string {
	return "cannot use an empty identifier ('_') for the first two spawn parameters"
}
func (eiosp *EmptyIdentifierOnSpawnParameters) Note() string { return "" }
func (eiosp *EmptyIdentifierOnSpawnParameters) Type() Type   { return Error }
func (eiosp *EmptyIdentifierOnSpawnParameters) ID() string   { return "hyb074W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidListOrMapWrappedType struct {
	SnippetProvider
}

func NewInvalidListOrMapWrappedType(span core.Span) *InvalidListOrMapWrappedType {
	return &InvalidListOrMapWrappedType{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ilomwt *InvalidListOrMapWrappedType) Message() string {
	return "lists and maps have a singular wrapped type"
}
func (ilomwt *InvalidListOrMapWrappedType) Note() string { return "" }
func (ilomwt *InvalidListOrMapWrappedType) Type() Type   { return Error }
func (ilomwt *InvalidListOrMapWrappedType) ID() string   { return "hyb075W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type AssignmentToSelf struct {
	SnippetProvider
	varName string
}

func NewAssignmentToSelf(span core.Span, varName string) *AssignmentToSelf {
	return &AssignmentToSelf{
		SnippetProvider: SnippetProvider{span: span},
		varName:         varName,
	}
}
func (ats *AssignmentToSelf) Message() string {
	return fmt.Sprintf("the variable '%s' is assigned to itself", ats.varName)
}
func (ats *AssignmentToSelf) Note() string { return "" }
func (ats *AssignmentToSelf) Type() Type   { return Warning }
func (ats *AssignmentToSelf) ID() string   { return "hyb076W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownListOrMapContents struct {
	SnippetProvider
}

func NewUnknownListOrMapContents(span core.Span) *UnknownListOrMapContents {
	return &UnknownListOrMapContents{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ulomc *UnknownListOrMapContents) Message() string {
	return "lists or maps with no values need to have their wrapped type explicitly given"
}
func (ulomc *UnknownListOrMapContents) Note() string {
	return "this can be done like so: let exampleList = list<number>[] or let exampleMap = map<number>{}"
}
func (ulomc *UnknownListOrMapContents) Type() Type { return Error }
func (ulomc *UnknownListOrMapContents) ID() string { return "hyb077W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEntityForLoopType struct {
	SnippetProvider
}

func NewInvalidEntityForLoopType(span core.Span) *InvalidEntityForLoopType {
	return &InvalidEntityForLoopType{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ieflt *InvalidEntityForLoopType) Message() string {
	return "expected an entity type in the entity for loop"
}
func (ieflt *InvalidEntityForLoopType) Note() string { return "" }
func (ieflt *InvalidEntityForLoopType) Type() Type   { return Error }
func (ieflt *InvalidEntityForLoopType) ID() string   { return "hyb078W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidSpawnerParameter struct {
	SnippetProvider
	nth  string
	name string
}

func NewInvalidSpawnerParameter(span core.Span, nth string, name string) *InvalidSpawnerParameter {
	return &InvalidSpawnerParameter{
		SnippetProvider: SnippetProvider{span: span},
		nth:             nth,
		name:            name,
	}
}
func (isp *InvalidSpawnerParameter) Message() string {
	return fmt.Sprintf("the %s parameter has to be named '%s'", isp.nth, isp.name)
}
func (isp *InvalidSpawnerParameter) Note() string { return "" }
func (isp *InvalidSpawnerParameter) Type() Type   { return Error }
func (isp *InvalidSpawnerParameter) ID() string   { return "hyb079W" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnallowedNumberInEnvironment struct {
	SnippetProvider
	numberType string
	envType    string
}

func NewUnallowedNumberInEnvironment(span core.Span, numberType string, envType string) *UnallowedNumberInEnvironment {
	return &UnallowedNumberInEnvironment{
		SnippetProvider: SnippetProvider{span: span},
		numberType:      numberType,
		envType:         envType,
	}
}
func (unie *UnallowedNumberInEnvironment) Message() string {
	return fmt.Sprintf("%s numbers are not allowed in a %s environment", unie.numberType, unie.envType)
}
func (unie *UnallowedNumberInEnvironment) Note() string { return "" }
func (unie *UnallowedNumberInEnvironment) Type() Type   { return Error }
func (unie *UnallowedNumberInEnvironment) ID() string   { return "hyb080W" }
