// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnclosingMark struct {
  Specifier SnippetSpecifier
  Mark string
}

func (eem *ExpectedEnclosingMark) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", eem.Mark)
}

func (eem *ExpectedEnclosingMark) GetSpecifier() SnippetSpecifier {
  return eem.Specifier
}

func (eem *ExpectedEnclosingMark) GetNote() string {
  return ""
}

func (eem *ExpectedEnclosingMark) GetID() string {
  return "hyb001"
}

func (eem *ExpectedEnclosingMark) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedOpeningMark struct {
  Specifier SnippetSpecifier
  Mark string
}

func (eom *ExpectedOpeningMark) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", eom.Mark)
}

func (eom *ExpectedOpeningMark) GetSpecifier() SnippetSpecifier {
  return eom.Specifier
}

func (eom *ExpectedOpeningMark) GetNote() string {
  return ""
}

func (eom *ExpectedOpeningMark) GetID() string {
  return "hyb002"
}

func (eom *ExpectedOpeningMark) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedSymbol struct {
  Specifier SnippetSpecifier
  Symbol string
  Context string `default:""""`
}

func (es *ExpectedSymbol) GetMessage() string {
  return "Expected '%s' %s"
}

func (es *ExpectedSymbol) GetSpecifier() SnippetSpecifier {
  return es.Specifier
}

func (es *ExpectedSymbol) GetNote() string {
  return ""
}

func (es *ExpectedSymbol) GetID() string {
  return "hyb003"
}

func (es *ExpectedSymbol) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneElseStatement struct {
  Specifier SnippetSpecifier
}

func (mtoes *MoreThanOneElseStatement) GetMessage() string {
  return "Cannot have more than one else statement in an if statement"
}

func (mtoes *MoreThanOneElseStatement) GetSpecifier() SnippetSpecifier {
  return mtoes.Specifier
}

func (mtoes *MoreThanOneElseStatement) GetNote() string {
  return ""
}

func (mtoes *MoreThanOneElseStatement) GetID() string {
  return "hyb004"
}

func (mtoes *MoreThanOneElseStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedKeyword struct {
  Specifier SnippetSpecifier
  Keyword string
  Context string `default:""""`
}

func (ek *ExpectedKeyword) GetMessage() string {
  return "Expected keyword '%s' %s"
}

func (ek *ExpectedKeyword) GetSpecifier() SnippetSpecifier {
  return ek.Specifier
}

func (ek *ExpectedKeyword) GetNote() string {
  return ""
}

func (ek *ExpectedKeyword) GetID() string {
  return "hyb005"
}

func (ek *ExpectedKeyword) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EnvironmentRedaclaration struct {
  Specifier SnippetSpecifier
}

func (er *EnvironmentRedaclaration) GetMessage() string {
  return "Cannot redeclare an environment"
}

func (er *EnvironmentRedaclaration) GetSpecifier() SnippetSpecifier {
  return er.Specifier
}

func (er *EnvironmentRedaclaration) GetNote() string {
  return ""
}

func (er *EnvironmentRedaclaration) GetID() string {
  return "hyb006"
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
  return "hyb007"
}

func (eepe *ExpectedEnvironmentPathExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
  Specifier SnippetSpecifier
}

func (ee *ExpectedEnvironment) GetMessage() string {
  return "Expected environment statement"
}

func (ee *ExpectedEnvironment) GetSpecifier() SnippetSpecifier {
  return ee.Specifier
}

func (ee *ExpectedEnvironment) GetNote() string {
  return "environment statement has to be the first statement in any hybroid file. example: env HelloWorld as Level"
}

func (ee *ExpectedEnvironment) GetID() string {
  return "hyb008"
}

func (ee *ExpectedEnvironment) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
  Specifier SnippetSpecifier
  Context string `default:""""`
}

func (ei *ExpectedIdentifier) GetMessage() string {
  return fmt.Sprintf("Expected identifier %s", ei.Context)
}

func (ei *ExpectedIdentifier) GetSpecifier() SnippetSpecifier {
  return ei.Specifier
}

func (ei *ExpectedIdentifier) GetNote() string {
  return ""
}

func (ei *ExpectedIdentifier) GetID() string {
  return "hyb009"
}

func (ei *ExpectedIdentifier) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
  Specifier SnippetSpecifier
}

func (et *ExpectedType) GetMessage() string {
  return "Expected type"
}

func (et *ExpectedType) GetSpecifier() SnippetSpecifier {
  return et.Specifier
}

func (et *ExpectedType) GetNote() string {
  return fmt.Sprintf("this needs to be declared with a type. example: number %s", et.Specifier.GetTokens()[0].Lexeme)
}

func (et *ExpectedType) GetID() string {
  return "hyb010"
}

func (et *ExpectedType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
  Specifier SnippetSpecifier
}

func (ee *ExpectedExpression) GetMessage() string {
  return "Expected expression"
}

func (ee *ExpectedExpression) GetSpecifier() SnippetSpecifier {
  return ee.Specifier
}

func (ee *ExpectedExpression) GetNote() string {
  return ""
}

func (ee *ExpectedExpression) GetID() string {
  return "hyb011"
}

func (ee *ExpectedExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
  Specifier SnippetSpecifier
}

func (eeob *ExpectedExpressionOrBody) GetMessage() string {
  return "Expected expression or body"
}

func (eeob *ExpectedExpressionOrBody) GetSpecifier() SnippetSpecifier {
  return eeob.Specifier
}

func (eeob *ExpectedExpressionOrBody) GetNote() string {
  return ""
}

func (eeob *ExpectedExpressionOrBody) GetID() string {
  return "hyb012"
}

func (eeob *ExpectedExpressionOrBody) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArrow struct {
  Specifier SnippetSpecifier
}

func (era *ExpectedReturnArrow) GetMessage() string {
  return "Expected a return arrow (->)"
}

func (era *ExpectedReturnArrow) GetSpecifier() SnippetSpecifier {
  return era.Specifier
}

func (era *ExpectedReturnArrow) GetNote() string {
  return "return types on function declarations are written after a thin arrow (->) after the parameters"
}

func (era *ExpectedReturnArrow) GetID() string {
  return "hyb013"
}

func (era *ExpectedReturnArrow) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallArgs struct {
  Specifier SnippetSpecifier
}

func (eca *ExpectedCallArgs) GetMessage() string {
  return "Expected call arguments"
}

func (eca *ExpectedCallArgs) GetSpecifier() SnippetSpecifier {
  return eca.Specifier
}

func (eca *ExpectedCallArgs) GetNote() string {
  return ""
}

func (eca *ExpectedCallArgs) GetID() string {
  return "hyb014"
}

func (eca *ExpectedCallArgs) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCall struct {
  Specifier SnippetSpecifier
}

func (ic *InvalidCall) GetMessage() string {
  return "Invalid expression to call"
}

func (ic *InvalidCall) GetSpecifier() SnippetSpecifier {
  return ic.Specifier
}

func (ic *InvalidCall) GetNote() string {
  return ""
}

func (ic *InvalidCall) GetID() string {
  return "hyb015"
}

func (ic *InvalidCall) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
  Specifier SnippetSpecifier
}

func (iet *InvalidEnvironmentType) GetMessage() string {
  return "Expected 'Level', 'Mesh' or 'Sound' as environment type"
}

func (iet *InvalidEnvironmentType) GetSpecifier() SnippetSpecifier {
  return iet.Specifier
}

func (iet *InvalidEnvironmentType) GetNote() string {
  return ""
}

func (iet *InvalidEnvironmentType) GetID() string {
  return "hyb016"
}

func (iet *InvalidEnvironmentType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallAfterMacroSymbol struct {
  Specifier SnippetSpecifier
}

func (ecams *ExpectedCallAfterMacroSymbol) GetMessage() string {
  return "Expected an expression call after '@'"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetSpecifier() SnippetSpecifier {
  return ecams.Specifier
}

func (ecams *ExpectedCallAfterMacroSymbol) GetNote() string {
  return ""
}

func (ecams *ExpectedCallAfterMacroSymbol) GetID() string {
  return "hyb017"
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
  return fmt.Sprintf("Cannot have a %s in the following environments: %v", ftie.Type, ftie.Envs)
}

func (ftie *ForbiddenTypeInEnvironment) GetSpecifier() SnippetSpecifier {
  return ftie.Specifier
}

func (ftie *ForbiddenTypeInEnvironment) GetNote() string {
  return ""
}

func (ftie *ForbiddenTypeInEnvironment) GetID() string {
  return "hyb018"
}

func (ftie *ForbiddenTypeInEnvironment) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFieldDeclaration struct {
  Specifier SnippetSpecifier
}

func (efd *ExpectedFieldDeclaration) GetMessage() string {
  return "Expected field declaration inside struct"
}

func (efd *ExpectedFieldDeclaration) GetSpecifier() SnippetSpecifier {
  return efd.Specifier
}

func (efd *ExpectedFieldDeclaration) GetNote() string {
  return ""
}

func (efd *ExpectedFieldDeclaration) GetID() string {
  return "hyb019"
}

func (efd *ExpectedFieldDeclaration) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyWrappedType struct {
  Specifier SnippetSpecifier
}

func (ewt *EmptyWrappedType) GetMessage() string {
  return "Wrapped types must not be empty"
}

func (ewt *EmptyWrappedType) GetSpecifier() SnippetSpecifier {
  return ewt.Specifier
}

func (ewt *EmptyWrappedType) GetNote() string {
  return ""
}

func (ewt *EmptyWrappedType) GetID() string {
  return "hyb020"
}

func (ewt *EmptyWrappedType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArgs struct {
  Specifier SnippetSpecifier
}

func (era *ExpectedReturnArgs) GetMessage() string {
  return "Expected return arguments after fat arrow (=>)"
}

func (era *ExpectedReturnArgs) GetSpecifier() SnippetSpecifier {
  return era.Specifier
}

func (era *ExpectedReturnArgs) GetNote() string {
  return ""
}

func (era *ExpectedReturnArgs) GetID() string {
  return "hyb021"
}

func (era *ExpectedReturnArgs) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedStatement struct {
  Specifier SnippetSpecifier
}

func (es *ExpectedStatement) GetMessage() string {
  return "Expected statement"
}

func (es *ExpectedStatement) GetSpecifier() SnippetSpecifier {
  return es.Specifier
}

func (es *ExpectedStatement) GetNote() string {
  return ""
}

func (es *ExpectedStatement) GetID() string {
  return "hyb022"
}

func (es *ExpectedStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownStatementInsideClass struct {
  Specifier SnippetSpecifier
}

func (usic *UnknownStatementInsideClass) GetMessage() string {
  return "Unknown statement inside class"
}

func (usic *UnknownStatementInsideClass) GetSpecifier() SnippetSpecifier {
  return usic.Specifier
}

func (usic *UnknownStatementInsideClass) GetNote() string {
  return ""
}

func (usic *UnknownStatementInsideClass) GetID() string {
  return "hyb023"
}

func (usic *UnknownStatementInsideClass) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAccessExpression struct {
  Specifier SnippetSpecifier
}

func (eae *ExpectedAccessExpression) GetMessage() string {
  return "Expected an access expression"
}

func (eae *ExpectedAccessExpression) GetSpecifier() SnippetSpecifier {
  return eae.Specifier
}

func (eae *ExpectedAccessExpression) GetNote() string {
  return "Access expression are: identifier, environment access, self, member and field expressions"
}

func (eae *ExpectedAccessExpression) GetID() string {
  return "hyb024"
}

func (eae *ExpectedAccessExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type NoIteratorProvidedInForLoopStatement struct {
  Specifier SnippetSpecifier
}

func (nipifls *NoIteratorProvidedInForLoopStatement) GetMessage() string {
  return "No Iterator provided in for loop statement"
}

func (nipifls *NoIteratorProvidedInForLoopStatement) GetSpecifier() SnippetSpecifier {
  return nipifls.Specifier
}

func (nipifls *NoIteratorProvidedInForLoopStatement) GetNote() string {
  return ""
}

func (nipifls *NoIteratorProvidedInForLoopStatement) GetID() string {
  return "hyb025"
}

func (nipifls *NoIteratorProvidedInForLoopStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateKeywordInRepeatStatement struct {
  Specifier SnippetSpecifier
  Keyword string
}

func (dkirs *DuplicateKeywordInRepeatStatement) GetMessage() string {
  return fmt.Sprintf("Cannot have duplicate keyword (%s) in repeat statement", dkirs.Keyword)
}

func (dkirs *DuplicateKeywordInRepeatStatement) GetSpecifier() SnippetSpecifier {
  return dkirs.Specifier
}

func (dkirs *DuplicateKeywordInRepeatStatement) GetNote() string {
  return ""
}

func (dkirs *DuplicateKeywordInRepeatStatement) GetID() string {
  return "hyb026"
}

func (dkirs *DuplicateKeywordInRepeatStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type RedefinitionOfIteratorInRepeatStatement struct {
  Specifier SnippetSpecifier
}

func (roiirs *RedefinitionOfIteratorInRepeatStatement) GetMessage() string {
  return "Redefinition of iterator in repeat statement"
}

func (roiirs *RedefinitionOfIteratorInRepeatStatement) GetSpecifier() SnippetSpecifier {
  return roiirs.Specifier
}

func (roiirs *RedefinitionOfIteratorInRepeatStatement) GetNote() string {
  return ""
}

func (roiirs *RedefinitionOfIteratorInRepeatStatement) GetID() string {
  return "hyb027"
}

func (roiirs *RedefinitionOfIteratorInRepeatStatement) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFatArrowInMatchCase struct {
  Specifier SnippetSpecifier
}

func (efaimc *ExpectedFatArrowInMatchCase) GetMessage() string {
  return "Expected fat arrow (=>) in match case"
}

func (efaimc *ExpectedFatArrowInMatchCase) GetSpecifier() SnippetSpecifier {
  return efaimc.Specifier
}

func (efaimc *ExpectedFatArrowInMatchCase) GetNote() string {
  return ""
}

func (efaimc *ExpectedFatArrowInMatchCase) GetID() string {
  return "hyb028"
}

func (efaimc *ExpectedFatArrowInMatchCase) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownStatementInsideEntity struct {
  Specifier SnippetSpecifier
}

func (usie *UnknownStatementInsideEntity) GetMessage() string {
  return "Unknown statement inside class"
}

func (usie *UnknownStatementInsideEntity) GetSpecifier() SnippetSpecifier {
  return usie.Specifier
}

func (usie *UnknownStatementInsideEntity) GetNote() string {
  return ""
}

func (usie *UnknownStatementInsideEntity) GetID() string {
  return "hyb029"
}

func (usie *UnknownStatementInsideEntity) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingIteratorInRepeatStatement struct {
  Specifier SnippetSpecifier
}

func (miirs *MissingIteratorInRepeatStatement) GetMessage() string {
  return "Missing iterator in repeat statement"
}

func (miirs *MissingIteratorInRepeatStatement) GetSpecifier() SnippetSpecifier {
  return miirs.Specifier
}

func (miirs *MissingIteratorInRepeatStatement) GetNote() string {
  return ""
}

func (miirs *MissingIteratorInRepeatStatement) GetID() string {
  return "hyb030"
}

func (miirs *MissingIteratorInRepeatStatement) GetAlertType() AlertType {
  return Error
}

