// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "strings"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedStatement struct {
  Specifier Snippet
}

func (es *ExpectedStatement) GetMessage() string {
  return "expected statement"
}

func (es *ExpectedStatement) GetSpecifier() Snippet {
  return es.Specifier
}

func (es *ExpectedStatement) GetNote() string {
  return ""
}

func (es *ExpectedStatement) GetID() string {
  return "hyb001P"
}

func (es *ExpectedStatement) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
  Specifier Snippet
  Context string `default:""`
}

func (ee *ExpectedExpression) GetMessage() string {
  return fmt.Sprintf("expected expression %s", ee.Context)
}

func (ee *ExpectedExpression) GetSpecifier() Snippet {
  return ee.Specifier
}

func (ee *ExpectedExpression) GetNote() string {
  return ""
}

func (ee *ExpectedExpression) GetID() string {
  return "hyb002P"
}

func (ee *ExpectedExpression) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownStatement struct {
  Specifier Snippet
  Context string `default:""`
}

func (us *UnknownStatement) GetMessage() string {
  return fmt.Sprintf("unknown statement %s", us.Context)
}

func (us *UnknownStatement) GetSpecifier() Snippet {
  return us.Specifier
}

func (us *UnknownStatement) GetNote() string {
  return ""
}

func (us *UnknownStatement) GetID() string {
  return "hyb003P"
}

func (us *UnknownStatement) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedKeyword struct {
  Specifier Snippet
  Keyword string
  Context string `default:""`
}

func (ek *ExpectedKeyword) GetMessage() string {
  return fmt.Sprintf("expected keyword '%s' %s", ek.Keyword, ek.Context)
}

func (ek *ExpectedKeyword) GetSpecifier() Snippet {
  return ek.Specifier
}

func (ek *ExpectedKeyword) GetNote() string {
  return ""
}

func (ek *ExpectedKeyword) GetID() string {
  return "hyb004P"
}

func (ek *ExpectedKeyword) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
  Specifier Snippet
  Context string `default:""`
}

func (ei *ExpectedIdentifier) GetMessage() string {
  return fmt.Sprintf("expected identifier %s", ei.Context)
}

func (ei *ExpectedIdentifier) GetSpecifier() Snippet {
  return ei.Specifier
}

func (ei *ExpectedIdentifier) GetNote() string {
  return ""
}

func (ei *ExpectedIdentifier) GetID() string {
  return "hyb005P"
}

func (ei *ExpectedIdentifier) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedSymbol struct {
  Specifier Snippet
  Symbol string
  Context string `default:""`
}

func (es *ExpectedSymbol) GetMessage() string {
  return fmt.Sprintf("expected '%s' %s", es.Symbol, es.Context)
}

func (es *ExpectedSymbol) GetSpecifier() Snippet {
  return es.Specifier
}

func (es *ExpectedSymbol) GetNote() string {
  return ""
}

func (es *ExpectedSymbol) GetID() string {
  return "hyb006P"
}

func (es *ExpectedSymbol) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneElseBlock struct {
  Specifier Snippet
}

func (mtoeb *MoreThanOneElseBlock) GetMessage() string {
  return "cannot have more than one else block in an if statement"
}

func (mtoeb *MoreThanOneElseBlock) GetSpecifier() Snippet {
  return mtoeb.Specifier
}

func (mtoeb *MoreThanOneElseBlock) GetNote() string {
  return ""
}

func (mtoeb *MoreThanOneElseBlock) GetID() string {
  return "hyb007P"
}

func (mtoeb *MoreThanOneElseBlock) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneConstructor struct {
  Specifier Snippet
}

func (mtoc *MoreThanOneConstructor) GetMessage() string {
  return "cannot have more than one constructor in class declaration"
}

func (mtoc *MoreThanOneConstructor) GetSpecifier() Snippet {
  return mtoc.Specifier
}

func (mtoc *MoreThanOneConstructor) GetNote() string {
  return ""
}

func (mtoc *MoreThanOneConstructor) GetID() string {
  return "hyb008P"
}

func (mtoc *MoreThanOneConstructor) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneEntityFunction struct {
  Specifier Snippet
  FunctionType string
}

func (mtoef *MoreThanOneEntityFunction) GetMessage() string {
  return fmt.Sprintf("cannot have more than one %s in entity declaration", mtoef.FunctionType)
}

func (mtoef *MoreThanOneEntityFunction) GetSpecifier() Snippet {
  return mtoef.Specifier
}

func (mtoef *MoreThanOneEntityFunction) GetNote() string {
  return ""
}

func (mtoef *MoreThanOneEntityFunction) GetID() string {
  return "hyb009P"
}

func (mtoef *MoreThanOneEntityFunction) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultipleIdentifiersInCompoundAssignment struct {
  Specifier Snippet
}

func (miica *MultipleIdentifiersInCompoundAssignment) GetMessage() string {
  return "cannot have more than one left-hand identifier in a compound assignment"
}

func (miica *MultipleIdentifiersInCompoundAssignment) GetSpecifier() Snippet {
  return miica.Specifier
}

func (miica *MultipleIdentifiersInCompoundAssignment) GetNote() string {
  return "compound assignments include +=, -=, *=, /=, etc."
}

func (miica *MultipleIdentifiersInCompoundAssignment) GetID() string {
  return "hyb010P"
}

func (miica *MultipleIdentifiersInCompoundAssignment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ReturnsInConstructor struct {
  Specifier Snippet
}

func (ric *ReturnsInConstructor) GetMessage() string {
  return "cannot have return types in constructor"
}

func (ric *ReturnsInConstructor) GetSpecifier() Snippet {
  return ric.Specifier
}

func (ric *ReturnsInConstructor) GetNote() string {
  return ""
}

func (ric *ReturnsInConstructor) GetID() string {
  return "hyb011P"
}

func (ric *ReturnsInConstructor) GetAlertType() Type {
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
  return "hyb012P"
}

func (er *EnvironmentRedaclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironmentPathExpression struct {
  Specifier Snippet
}

func (eepe *ExpectedEnvironmentPathExpression) GetMessage() string {
  return "expected environment path expression"
}

func (eepe *ExpectedEnvironmentPathExpression) GetSpecifier() Snippet {
  return eepe.Specifier
}

func (eepe *ExpectedEnvironmentPathExpression) GetNote() string {
  return ""
}

func (eepe *ExpectedEnvironmentPathExpression) GetID() string {
  return "hyb013P"
}

func (eepe *ExpectedEnvironmentPathExpression) GetAlertType() Type {
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
  return "hyb014P"
}

func (ee *ExpectedEnvironment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
  Specifier Snippet
  Context string `default:""`
}

func (et *ExpectedType) GetMessage() string {
  return fmt.Sprintf("expected type %s", et.Context)
}

func (et *ExpectedType) GetSpecifier() Snippet {
  return et.Specifier
}

func (et *ExpectedType) GetNote() string {
  return "access expressions are: identifier, environment access, self, member and field expressions"
}

func (et *ExpectedType) GetID() string {
  return "hyb015P"
}

func (et *ExpectedType) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAssignmentSymbol struct {
  Specifier Snippet
}

func (eas *ExpectedAssignmentSymbol) GetMessage() string {
  return "expected assignment symbol"
}

func (eas *ExpectedAssignmentSymbol) GetSpecifier() Snippet {
  return eas.Specifier
}

func (eas *ExpectedAssignmentSymbol) GetNote() string {
  return "assignment symbols are: '=', '+=', '-=', '*=', '%%=', '/=', '\\='"
}

func (eas *ExpectedAssignmentSymbol) GetID() string {
  return "hyb016P"
}

func (eas *ExpectedAssignmentSymbol) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
  Specifier Snippet
}

func (eeob *ExpectedExpressionOrBody) GetMessage() string {
  return "expected expression or body"
}

func (eeob *ExpectedExpressionOrBody) GetSpecifier() Snippet {
  return eeob.Specifier
}

func (eeob *ExpectedExpressionOrBody) GetNote() string {
  return ""
}

func (eeob *ExpectedExpressionOrBody) GetID() string {
  return "hyb017P"
}

func (eeob *ExpectedExpressionOrBody) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallArgs struct {
  Specifier Snippet
}

func (eca *ExpectedCallArgs) GetMessage() string {
  return "expected call arguments"
}

func (eca *ExpectedCallArgs) GetSpecifier() Snippet {
  return eca.Specifier
}

func (eca *ExpectedCallArgs) GetNote() string {
  return ""
}

func (eca *ExpectedCallArgs) GetID() string {
  return "hyb018P"
}

func (eca *ExpectedCallArgs) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCall struct {
  Specifier Snippet
}

func (ic *InvalidCall) GetMessage() string {
  return "invalid expression to call"
}

func (ic *InvalidCall) GetSpecifier() Snippet {
  return ic.Specifier
}

func (ic *InvalidCall) GetNote() string {
  return ""
}

func (ic *InvalidCall) GetID() string {
  return "hyb019P"
}

func (ic *InvalidCall) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
  Specifier Snippet
}

func (iet *InvalidEnvironmentType) GetMessage() string {
  return "expected 'Level', 'Mesh' or 'Sound' as environment type"
}

func (iet *InvalidEnvironmentType) GetSpecifier() Snippet {
  return iet.Specifier
}

func (iet *InvalidEnvironmentType) GetNote() string {
  return ""
}

func (iet *InvalidEnvironmentType) GetID() string {
  return "hyb020P"
}

func (iet *InvalidEnvironmentType) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallAfterMacroSymbol struct {
  Specifier Snippet
}

func (ecams *ExpectedCallAfterMacroSymbol) GetMessage() string {
  return "expected a macro call after '@'"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetSpecifier() Snippet {
  return ecams.Specifier
}

func (ecams *ExpectedCallAfterMacroSymbol) GetNote() string {
  return ""
}

func (ecams *ExpectedCallAfterMacroSymbol) GetID() string {
  return "hyb021P"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetAlertType() Type {
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
  return "hyb022P"
}

func (ftie *ForbiddenTypeInEnvironment) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFieldDeclaration struct {
  Specifier Snippet
}

func (efd *ExpectedFieldDeclaration) GetMessage() string {
  return "expected field declaration inside struct"
}

func (efd *ExpectedFieldDeclaration) GetSpecifier() Snippet {
  return efd.Specifier
}

func (efd *ExpectedFieldDeclaration) GetNote() string {
  return ""
}

func (efd *ExpectedFieldDeclaration) GetID() string {
  return "hyb023P"
}

func (efd *ExpectedFieldDeclaration) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyWrappedType struct {
  Specifier Snippet
}

func (ewt *EmptyWrappedType) GetMessage() string {
  return "wrapped types must not be empty"
}

func (ewt *EmptyWrappedType) GetSpecifier() Snippet {
  return ewt.Specifier
}

func (ewt *EmptyWrappedType) GetNote() string {
  return ""
}

func (ewt *EmptyWrappedType) GetID() string {
  return "hyb024P"
}

func (ewt *EmptyWrappedType) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArgs struct {
  Specifier Snippet
}

func (era *ExpectedReturnArgs) GetMessage() string {
  return "expected return arguments after fat arrow (=>)"
}

func (era *ExpectedReturnArgs) GetSpecifier() Snippet {
  return era.Specifier
}

func (era *ExpectedReturnArgs) GetNote() string {
  return ""
}

func (era *ExpectedReturnArgs) GetID() string {
  return "hyb025P"
}

func (era *ExpectedReturnArgs) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAccessExpression struct {
  Specifier Snippet
}

func (eae *ExpectedAccessExpression) GetMessage() string {
  return "expected an access expression"
}

func (eae *ExpectedAccessExpression) GetSpecifier() Snippet {
  return eae.Specifier
}

func (eae *ExpectedAccessExpression) GetNote() string {
  return "access expression are: identifier, environment access, self, member and field expressions"
}

func (eae *ExpectedAccessExpression) GetID() string {
  return "hyb026P"
}

func (eae *ExpectedAccessExpression) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingIterator struct {
  Specifier Snippet
  Context string `default:""`
}

func (mi *MissingIterator) GetMessage() string {
  return fmt.Sprintf("missing iterator %s", mi.Context)
}

func (mi *MissingIterator) GetSpecifier() Snippet {
  return mi.Specifier
}

func (mi *MissingIterator) GetNote() string {
  return ""
}

func (mi *MissingIterator) GetID() string {
  return "hyb027P"
}

func (mi *MissingIterator) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateKeyword struct {
  Specifier Snippet
  Keyword string
}

func (dk *DuplicateKeyword) GetMessage() string {
  return fmt.Sprintf("cannot have multiple '%s' keywords", dk.Keyword)
}

func (dk *DuplicateKeyword) GetSpecifier() Snippet {
  return dk.Specifier
}

func (dk *DuplicateKeyword) GetNote() string {
  return ""
}

func (dk *DuplicateKeyword) GetID() string {
  return "hyb028P"
}

func (dk *DuplicateKeyword) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnexpectedKeyword struct {
  Specifier Snippet
  Keyword string
  Context string `default:""`
}

func (uk *UnexpectedKeyword) GetMessage() string {
  return fmt.Sprintf("unexpected keyword '%s' %s", uk.Keyword, uk.Context)
}

func (uk *UnexpectedKeyword) GetSpecifier() Snippet {
  return uk.Specifier
}

func (uk *UnexpectedKeyword) GetNote() string {
  return ""
}

func (uk *UnexpectedKeyword) GetID() string {
  return "hyb029P"
}

func (uk *UnexpectedKeyword) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type IteratorRedefinition struct {
  Specifier Snippet
  Context string `default:""`
}

func (ir *IteratorRedefinition) GetMessage() string {
  return fmt.Sprintf("redefinition of iterator %s", ir.Context)
}

func (ir *IteratorRedefinition) GetSpecifier() Snippet {
  return ir.Specifier
}

func (ir *IteratorRedefinition) GetNote() string {
  return ""
}

func (ir *IteratorRedefinition) GetID() string {
  return "hyb030P"
}

func (ir *IteratorRedefinition) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ElseIfBlockAfterElseBlock struct {
  Specifier Snippet
}

func (eibaeb *ElseIfBlockAfterElseBlock) GetMessage() string {
  return "cannot have an else if block after an else block"
}

func (eibaeb *ElseIfBlockAfterElseBlock) GetSpecifier() Snippet {
  return eibaeb.Specifier
}

func (eibaeb *ElseIfBlockAfterElseBlock) GetNote() string {
  return ""
}

func (eibaeb *ElseIfBlockAfterElseBlock) GetID() string {
  return "hyb031P"
}

func (eibaeb *ElseIfBlockAfterElseBlock) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneDefaultCase struct {
  Specifier Snippet
}

func (mtodc *MoreThanOneDefaultCase) GetMessage() string {
  return "cannot have more than one default case in match statement"
}

func (mtodc *MoreThanOneDefaultCase) GetSpecifier() Snippet {
  return mtodc.Specifier
}

func (mtodc *MoreThanOneDefaultCase) GetNote() string {
  return ""
}

func (mtodc *MoreThanOneDefaultCase) GetID() string {
  return "hyb032P"
}

func (mtodc *MoreThanOneDefaultCase) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InsufficientCases struct {
  Specifier Snippet
}

func (ic *InsufficientCases) GetMessage() string {
  return "match statement must have at least 1 case"
}

func (ic *InsufficientCases) GetSpecifier() Snippet {
  return ic.Specifier
}

func (ic *InsufficientCases) GetNote() string {
  return ""
}

func (ic *InsufficientCases) GetID() string {
  return "hyb033P"
}

func (ic *InsufficientCases) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DefaultCaseMissing struct {
  Specifier Snippet
}

func (dcm *DefaultCaseMissing) GetMessage() string {
  return "match statement must have 1 default case"
}

func (dcm *DefaultCaseMissing) GetSpecifier() Snippet {
  return dcm.Specifier
}

func (dcm *DefaultCaseMissing) GetNote() string {
  return "default cases start with 'else'"
}

func (dcm *DefaultCaseMissing) GetID() string {
  return "hyb034P"
}

func (dcm *DefaultCaseMissing) GetAlertType() Type {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnumVariantName struct {
  Specifier Snippet
}

func (ievn *InvalidEnumVariantName) GetMessage() string {
  return "enum variant name must be an identifier"
}

func (ievn *InvalidEnumVariantName) GetSpecifier() Snippet {
  return ievn.Specifier
}

func (ievn *InvalidEnumVariantName) GetNote() string {
  return ""
}

func (ievn *InvalidEnumVariantName) GetID() string {
  return "hyb035P"
}

func (ievn *InvalidEnumVariantName) GetAlertType() Type {
  return Error
}

