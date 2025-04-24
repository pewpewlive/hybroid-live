// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "strings"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedStatement struct {
  Specifier SnippetSpecifier
}

func (es *ExpectedStatement) GetMessage() string {
  return "expected statement"
}

func (es *ExpectedStatement) GetSpecifier() SnippetSpecifier {
  return es.Specifier
}

func (es *ExpectedStatement) GetNote() string {
  return ""
}

func (es *ExpectedStatement) GetID() string {
  return "hyb001"
}

func (es *ExpectedStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
  Specifier SnippetSpecifier
  Context string `default:""`
}

func (ee *ExpectedExpression) GetMessage() string {
  return fmt.Sprintf("expected expression %s", ee.Context)
}

func (ee *ExpectedExpression) GetSpecifier() SnippetSpecifier {
  return ee.Specifier
}

func (ee *ExpectedExpression) GetNote() string {
  return ""
}

func (ee *ExpectedExpression) GetID() string {
  return "hyb002"
}

func (ee *ExpectedExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownStatement struct {
  Specifier SnippetSpecifier
  Context string `default:""`
}

func (us *UnknownStatement) GetMessage() string {
  return fmt.Sprintf("unknown statement %s", us.Context)
}

func (us *UnknownStatement) GetSpecifier() SnippetSpecifier {
  return us.Specifier
}

func (us *UnknownStatement) GetNote() string {
  return ""
}

func (us *UnknownStatement) GetID() string {
  return "hyb003"
}

func (us *UnknownStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedKeyword struct {
  Specifier SnippetSpecifier
  Keyword string
  Context string `default:""`
}

func (ek *ExpectedKeyword) GetMessage() string {
  return fmt.Sprintf("expected keyword '%s' %s", ek.Keyword, ek.Context)
}

func (ek *ExpectedKeyword) GetSpecifier() SnippetSpecifier {
  return ek.Specifier
}

func (ek *ExpectedKeyword) GetNote() string {
  return ""
}

func (ek *ExpectedKeyword) GetID() string {
  return "hyb004"
}

func (ek *ExpectedKeyword) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
  Specifier SnippetSpecifier
  Context string `default:""`
}

func (ei *ExpectedIdentifier) GetMessage() string {
  return fmt.Sprintf("expected identifier %s", ei.Context)
}

func (ei *ExpectedIdentifier) GetSpecifier() SnippetSpecifier {
  return ei.Specifier
}

func (ei *ExpectedIdentifier) GetNote() string {
  return ""
}

func (ei *ExpectedIdentifier) GetID() string {
  return "hyb005"
}

func (ei *ExpectedIdentifier) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedSymbol struct {
  Specifier SnippetSpecifier
  Symbol string
  Context string `default:""`
}

func (es *ExpectedSymbol) GetMessage() string {
  return fmt.Sprintf("expected '%s' %s", es.Symbol, es.Context)
}

func (es *ExpectedSymbol) GetSpecifier() SnippetSpecifier {
  return es.Specifier
}

func (es *ExpectedSymbol) GetNote() string {
  return ""
}

func (es *ExpectedSymbol) GetID() string {
  return "hyb006"
}

func (es *ExpectedSymbol) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneElseStatement struct {
  Specifier SnippetSpecifier
}

func (mtoes *MoreThanOneElseStatement) GetMessage() string {
  return "cannot have more than one else statement in an if statement"
}

func (mtoes *MoreThanOneElseStatement) GetSpecifier() SnippetSpecifier {
  return mtoes.Specifier
}

func (mtoes *MoreThanOneElseStatement) GetNote() string {
  return ""
}

func (mtoes *MoreThanOneElseStatement) GetID() string {
  return "hyb007"
}

func (mtoes *MoreThanOneElseStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneConstructor struct {
  Specifier SnippetSpecifier
}

func (mtoc *MoreThanOneConstructor) GetMessage() string {
  return "cannot have more than one constructor in class declaration"
}

func (mtoc *MoreThanOneConstructor) GetSpecifier() SnippetSpecifier {
  return mtoc.Specifier
}

func (mtoc *MoreThanOneConstructor) GetNote() string {
  return ""
}

func (mtoc *MoreThanOneConstructor) GetID() string {
  return "hyb008"
}

func (mtoc *MoreThanOneConstructor) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneEntityFunction struct {
  Specifier SnippetSpecifier
  FunctionType string
}

func (mtoef *MoreThanOneEntityFunction) GetMessage() string {
  return fmt.Sprintf("cannot have more than one %s in entity declaration", mtoef.FunctionType)
}

func (mtoef *MoreThanOneEntityFunction) GetSpecifier() SnippetSpecifier {
  return mtoef.Specifier
}

func (mtoef *MoreThanOneEntityFunction) GetNote() string {
  return ""
}

func (mtoef *MoreThanOneEntityFunction) GetID() string {
  return "hyb009"
}

func (mtoef *MoreThanOneEntityFunction) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ReturnsInConstructor struct {
  Specifier SnippetSpecifier
}

func (ric *ReturnsInConstructor) GetMessage() string {
  return "cannot have return types in constructor"
}

func (ric *ReturnsInConstructor) GetSpecifier() SnippetSpecifier {
  return ric.Specifier
}

func (ric *ReturnsInConstructor) GetNote() string {
  return ""
}

func (ric *ReturnsInConstructor) GetID() string {
  return "hyb010"
}

func (ric *ReturnsInConstructor) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentRedaclaration struct {
  Specifier SnippetSpecifier
}

func (er *EnvironmentRedaclaration) GetMessage() string {
  return "cannot redeclare an environment"
}

func (er *EnvironmentRedaclaration) GetSpecifier() SnippetSpecifier {
  return er.Specifier
}

func (er *EnvironmentRedaclaration) GetNote() string {
  return ""
}

func (er *EnvironmentRedaclaration) GetID() string {
  return "hyb011"
}

func (er *EnvironmentRedaclaration) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironmentPathExpression struct {
  Specifier SnippetSpecifier
}

func (eepe *ExpectedEnvironmentPathExpression) GetMessage() string {
  return "expected environment path expression"
}

func (eepe *ExpectedEnvironmentPathExpression) GetSpecifier() SnippetSpecifier {
  return eepe.Specifier
}

func (eepe *ExpectedEnvironmentPathExpression) GetNote() string {
  return ""
}

func (eepe *ExpectedEnvironmentPathExpression) GetID() string {
  return "hyb012"
}

func (eepe *ExpectedEnvironmentPathExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
  Specifier SnippetSpecifier
}

func (ee *ExpectedEnvironment) GetMessage() string {
  return "expected environment statement"
}

func (ee *ExpectedEnvironment) GetSpecifier() SnippetSpecifier {
  return ee.Specifier
}

func (ee *ExpectedEnvironment) GetNote() string {
  return "environment statement has to be the first statement in any hybroid file. example: env HelloWorld as Level"
}

func (ee *ExpectedEnvironment) GetID() string {
  return "hyb013"
}

func (ee *ExpectedEnvironment) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
  Specifier SnippetSpecifier
}

func (et *ExpectedType) GetMessage() string {
  return "expected type"
}

func (et *ExpectedType) GetSpecifier() SnippetSpecifier {
  return et.Specifier
}

func (et *ExpectedType) GetNote() string {
  return fmt.Sprintf("this needs to be declared with a type. example: number %s", et.Specifier.GetTokens()[0].Lexeme)
}

func (et *ExpectedType) GetID() string {
  return "hyb014"
}

func (et *ExpectedType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAssignmentSymbol struct {
  Specifier SnippetSpecifier
}

func (eas *ExpectedAssignmentSymbol) GetMessage() string {
  return "expected assignment symbol"
}

func (eas *ExpectedAssignmentSymbol) GetSpecifier() SnippetSpecifier {
  return eas.Specifier
}

func (eas *ExpectedAssignmentSymbol) GetNote() string {
  return "assignment symbols are: '=', '+=', '-=', '*=', '%=', '/=', '\\='"
}

func (eas *ExpectedAssignmentSymbol) GetID() string {
  return "hyb015"
}

func (eas *ExpectedAssignmentSymbol) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
  Specifier SnippetSpecifier
}

func (eeob *ExpectedExpressionOrBody) GetMessage() string {
  return "expected expression or body"
}

func (eeob *ExpectedExpressionOrBody) GetSpecifier() SnippetSpecifier {
  return eeob.Specifier
}

func (eeob *ExpectedExpressionOrBody) GetNote() string {
  return ""
}

func (eeob *ExpectedExpressionOrBody) GetID() string {
  return "hyb016"
}

func (eeob *ExpectedExpressionOrBody) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallArgs struct {
  Specifier SnippetSpecifier
}

func (eca *ExpectedCallArgs) GetMessage() string {
  return "expected call arguments"
}

func (eca *ExpectedCallArgs) GetSpecifier() SnippetSpecifier {
  return eca.Specifier
}

func (eca *ExpectedCallArgs) GetNote() string {
  return ""
}

func (eca *ExpectedCallArgs) GetID() string {
  return "hyb017"
}

func (eca *ExpectedCallArgs) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCall struct {
  Specifier SnippetSpecifier
}

func (ic *InvalidCall) GetMessage() string {
  return "invalid expression to call"
}

func (ic *InvalidCall) GetSpecifier() SnippetSpecifier {
  return ic.Specifier
}

func (ic *InvalidCall) GetNote() string {
  return ""
}

func (ic *InvalidCall) GetID() string {
  return "hyb018"
}

func (ic *InvalidCall) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
  Specifier SnippetSpecifier
}

func (iet *InvalidEnvironmentType) GetMessage() string {
  return "expected 'Level', 'Mesh' or 'Sound' as environment type"
}

func (iet *InvalidEnvironmentType) GetSpecifier() SnippetSpecifier {
  return iet.Specifier
}

func (iet *InvalidEnvironmentType) GetNote() string {
  return ""
}

func (iet *InvalidEnvironmentType) GetID() string {
  return "hyb019"
}

func (iet *InvalidEnvironmentType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallAfterMacroSymbol struct {
  Specifier SnippetSpecifier
}

func (ecams *ExpectedCallAfterMacroSymbol) GetMessage() string {
  return "expected a macro call after '@'"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetSpecifier() SnippetSpecifier {
  return ecams.Specifier
}

func (ecams *ExpectedCallAfterMacroSymbol) GetNote() string {
  return ""
}

func (ecams *ExpectedCallAfterMacroSymbol) GetID() string {
  return "hyb020"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForbiddenTypeInEnvironment struct {
  Specifier SnippetSpecifier
  Type string
  Envs []string
}

func (ftie *ForbiddenTypeInEnvironment) GetMessage() string {
  return fmt.Sprintf("cannot have a %s in the following environments: %s", ftie.Type, strings.Join(ftie.Envs, ", "))
}

func (ftie *ForbiddenTypeInEnvironment) GetSpecifier() SnippetSpecifier {
  return ftie.Specifier
}

func (ftie *ForbiddenTypeInEnvironment) GetNote() string {
  return ""
}

func (ftie *ForbiddenTypeInEnvironment) GetID() string {
  return "hyb021"
}

func (ftie *ForbiddenTypeInEnvironment) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFieldDeclaration struct {
  Specifier SnippetSpecifier
}

func (efd *ExpectedFieldDeclaration) GetMessage() string {
  return "expected field declaration inside struct"
}

func (efd *ExpectedFieldDeclaration) GetSpecifier() SnippetSpecifier {
  return efd.Specifier
}

func (efd *ExpectedFieldDeclaration) GetNote() string {
  return ""
}

func (efd *ExpectedFieldDeclaration) GetID() string {
  return "hyb022"
}

func (efd *ExpectedFieldDeclaration) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyWrappedType struct {
  Specifier SnippetSpecifier
}

func (ewt *EmptyWrappedType) GetMessage() string {
  return "wrapped types must not be empty"
}

func (ewt *EmptyWrappedType) GetSpecifier() SnippetSpecifier {
  return ewt.Specifier
}

func (ewt *EmptyWrappedType) GetNote() string {
  return ""
}

func (ewt *EmptyWrappedType) GetID() string {
  return "hyb023"
}

func (ewt *EmptyWrappedType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArgs struct {
  Specifier SnippetSpecifier
}

func (era *ExpectedReturnArgs) GetMessage() string {
  return "expected return arguments after fat arrow (=>)"
}

func (era *ExpectedReturnArgs) GetSpecifier() SnippetSpecifier {
  return era.Specifier
}

func (era *ExpectedReturnArgs) GetNote() string {
  return ""
}

func (era *ExpectedReturnArgs) GetID() string {
  return "hyb024"
}

func (era *ExpectedReturnArgs) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAccessExpression struct {
  Specifier SnippetSpecifier
}

func (eae *ExpectedAccessExpression) GetMessage() string {
  return "expected an access expression"
}

func (eae *ExpectedAccessExpression) GetSpecifier() SnippetSpecifier {
  return eae.Specifier
}

func (eae *ExpectedAccessExpression) GetNote() string {
  return "access expression are: identifier, environment access, self, member and field expressions"
}

func (eae *ExpectedAccessExpression) GetID() string {
  return "hyb025"
}

func (eae *ExpectedAccessExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingIterator struct {
  Specifier SnippetSpecifier
  Context string `default:""`
}

func (mi *MissingIterator) GetMessage() string {
  return fmt.Sprintf("missing iterator %s", mi.Context)
}

func (mi *MissingIterator) GetSpecifier() SnippetSpecifier {
  return mi.Specifier
}

func (mi *MissingIterator) GetNote() string {
  return ""
}

func (mi *MissingIterator) GetID() string {
  return "hyb026"
}

func (mi *MissingIterator) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateKeyword struct {
  Specifier SnippetSpecifier
  Keyword string
}

func (dk *DuplicateKeyword) GetMessage() string {
  return fmt.Sprintf("cannot have multiple '%s' keywords", dk.Keyword)
}

func (dk *DuplicateKeyword) GetSpecifier() SnippetSpecifier {
  return dk.Specifier
}

func (dk *DuplicateKeyword) GetNote() string {
  return ""
}

func (dk *DuplicateKeyword) GetID() string {
  return "hyb027"
}

func (dk *DuplicateKeyword) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnexpectedKeyword struct {
  Specifier SnippetSpecifier
  Keyword string
  Context string `default:""`
}

func (uk *UnexpectedKeyword) GetMessage() string {
  return fmt.Sprintf("unexpected keyword '%s' %s", uk.Keyword, uk.Context)
}

func (uk *UnexpectedKeyword) GetSpecifier() SnippetSpecifier {
  return uk.Specifier
}

func (uk *UnexpectedKeyword) GetNote() string {
  return ""
}

func (uk *UnexpectedKeyword) GetID() string {
  return "hyb028"
}

func (uk *UnexpectedKeyword) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type IteratorRedefinition struct {
  Specifier SnippetSpecifier
  Context string `default:""`
}

func (ir *IteratorRedefinition) GetMessage() string {
  return fmt.Sprintf("redefinition of iterator %s", ir.Context)
}

func (ir *IteratorRedefinition) GetSpecifier() SnippetSpecifier {
  return ir.Specifier
}

func (ir *IteratorRedefinition) GetNote() string {
  return ""
}

func (ir *IteratorRedefinition) GetID() string {
  return "hyb029"
}

func (ir *IteratorRedefinition) GetAlertType() AlertType {
  return Error
}

