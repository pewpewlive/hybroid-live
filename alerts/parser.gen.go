// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnclosingMark struct {
	Specifier Multiline
	Mark      string
}

func (eem *ExpectedEnclosingMark) GetMessage() string {
	return fmt.Sprintf("Expected '%s'", eem.Mark)
}

func (eem *ExpectedEnclosingMark) GetSpecifier() SnippetSpecifier {
	return &eem.Specifier
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
	Specifier Singleline
	Mark      string
}

func (eom *ExpectedOpeningMark) GetMessage() string {
	return fmt.Sprintf("Expected '%s'", eom.Mark)
}

func (eom *ExpectedOpeningMark) GetSpecifier() SnippetSpecifier {
	return &eom.Specifier
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
type ExpectedEnvironment struct {
	Specifier Singleline
}

func (ee *ExpectedEnvironment) GetMessage() string {
	return "Expected environment statement"
}

func (ee *ExpectedEnvironment) GetSpecifier() SnippetSpecifier {
	return &ee.Specifier
}

func (ee *ExpectedEnvironment) GetNote() string {
	return "environment statement has to be the first statement in any hybroid file. example: env HelloWorld as Level"
}

func (ee *ExpectedEnvironment) GetID() string {
	return "hyb003"
}

func (ee *ExpectedEnvironment) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
	Specifier Singleline
	Context   string `default:""""`
}

func (ei *ExpectedIdentifier) GetMessage() string {
	return "Expected identifier"
}

func (ei *ExpectedIdentifier) GetSpecifier() SnippetSpecifier {
	return &ei.Specifier
}

func (ei *ExpectedIdentifier) GetNote() string {
	return fmt.Sprintf("%s", ei.Context) // = ""
}

func (ei *ExpectedIdentifier) GetID() string {
	return "hyb004"
}

func (ei *ExpectedIdentifier) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
	Specifier Singleline
}

func (et *ExpectedType) GetMessage() string {
	return "Expected type"
}

func (et *ExpectedType) GetSpecifier() SnippetSpecifier {
	return &et.Specifier
}

func (et *ExpectedType) GetNote() string {
	return fmt.Sprintf("this needs to be declared with a type. example: number %s", et.Specifier.GetTokens()[0].Lexeme)
}

func (et *ExpectedType) GetID() string {
	return "hyb005"
}

func (et *ExpectedType) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
	Specifier Singleline
}

func (ee *ExpectedExpression) GetMessage() string {
	return "Expected expression"
}

func (ee *ExpectedExpression) GetSpecifier() SnippetSpecifier {
	return &ee.Specifier
}

func (ee *ExpectedExpression) GetNote() string {
	return ""
}

func (ee *ExpectedExpression) GetID() string {
	return "hyb006"
}

func (ee *ExpectedExpression) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
	Specifier Singleline
}

func (eeob *ExpectedExpressionOrBody) GetMessage() string {
	return "Expected expression or body"
}

func (eeob *ExpectedExpressionOrBody) GetSpecifier() SnippetSpecifier {
	return &eeob.Specifier
}

func (eeob *ExpectedExpressionOrBody) GetNote() string {
	return ""
}

func (eeob *ExpectedExpressionOrBody) GetID() string {
	return "hyb007"
}

func (eeob *ExpectedExpressionOrBody) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArrow struct {
	Specifier Singleline
}

func (era *ExpectedReturnArrow) GetMessage() string {
	return "Expected a return arrow (->)"
}

func (era *ExpectedReturnArrow) GetSpecifier() SnippetSpecifier {
	return &era.Specifier
}

func (era *ExpectedReturnArrow) GetNote() string {
	return "return types on function declarations are written after a thin arrow (->) after the parameters"
}

func (era *ExpectedReturnArrow) GetID() string {
	return "hyb008"
}

func (era *ExpectedReturnArrow) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallArgs struct {
	Specifier Singleline
}

func (eca *ExpectedCallArgs) GetMessage() string {
	return "Expected call arguments"
}

func (eca *ExpectedCallArgs) GetSpecifier() SnippetSpecifier {
	return &eca.Specifier
}

func (eca *ExpectedCallArgs) GetNote() string {
	return ""
}

func (eca *ExpectedCallArgs) GetID() string {
	return "hyb009"
}

func (eca *ExpectedCallArgs) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCall struct {
	Specifier Singleline
}

func (ic *InvalidCall) GetMessage() string {
	return "Invalid expression to call"
}

func (ic *InvalidCall) GetSpecifier() SnippetSpecifier {
	return &ic.Specifier
}

func (ic *InvalidCall) GetNote() string {
	return ""
}

func (ic *InvalidCall) GetID() string {
	return "hyb010"
}

func (ic *InvalidCall) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnvironmentType struct {
	Specifier Singleline
}

func (iet *InvalidEnvironmentType) GetMessage() string {
	return "Expected 'Level', 'Mesh' or 'Sound' as environment type"
}

func (iet *InvalidEnvironmentType) GetSpecifier() SnippetSpecifier {
	return &iet.Specifier
}

func (iet *InvalidEnvironmentType) GetNote() string {
	return ""
}

func (iet *InvalidEnvironmentType) GetID() string {
	return "hyb011"
}

func (iet *InvalidEnvironmentType) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallAfterMacroSymbol struct {
	Specifier Singleline
}

func (ecams *ExpectedCallAfterMacroSymbol) GetMessage() string {
	return "Expected an expression call after '@'"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetSpecifier() SnippetSpecifier {
	return &ecams.Specifier
}

func (ecams *ExpectedCallAfterMacroSymbol) GetNote() string {
	return ""
}

func (ecams *ExpectedCallAfterMacroSymbol) GetID() string {
	return "hyb012"
}

func (ecams *ExpectedCallAfterMacroSymbol) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ForbiddenTypeInEnvironment struct {
	Specifier Singleline
	Type      string
	Envs      []string
}

func (ftie *ForbiddenTypeInEnvironment) GetMessage() string {
	return fmt.Sprintf("Cannot have a %s in the following environments: %v", ftie.Type, ftie.Envs)
}

func (ftie *ForbiddenTypeInEnvironment) GetSpecifier() SnippetSpecifier {
	return &ftie.Specifier
}

func (ftie *ForbiddenTypeInEnvironment) GetNote() string {
	return ""
}

func (ftie *ForbiddenTypeInEnvironment) GetID() string {
	return "hyb013"
}

func (ftie *ForbiddenTypeInEnvironment) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFieldDeclaration struct {
	Specifier Singleline
}

func (efd *ExpectedFieldDeclaration) GetMessage() string {
	return "Expected field declaration inside struct"
}

func (efd *ExpectedFieldDeclaration) GetSpecifier() SnippetSpecifier {
	return &efd.Specifier
}

func (efd *ExpectedFieldDeclaration) GetNote() string {
	return ""
}

func (efd *ExpectedFieldDeclaration) GetID() string {
	return "hyb014"
}

func (efd *ExpectedFieldDeclaration) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyWrappedType struct {
	Specifier Singleline
}

func (ewt *EmptyWrappedType) GetMessage() string {
	return "Wrapped types must not be empty"
}

func (ewt *EmptyWrappedType) GetSpecifier() SnippetSpecifier {
	return &ewt.Specifier
}

func (ewt *EmptyWrappedType) GetNote() string {
	return ""
}

func (ewt *EmptyWrappedType) GetID() string {
	return "hyb015"
}

func (ewt *EmptyWrappedType) GetAlertType() AlertType {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedStatement struct {
	Specifier Singleline
}

func (es *ExpectedStatement) GetMessage() string {
	return "Expected statement"
}

func (es *ExpectedStatement) GetSpecifier() SnippetSpecifier {
	return &es.Specifier
}

func (es *ExpectedStatement) GetNote() string {
	return ""
}

func (es *ExpectedStatement) GetID() string {
	return "hyb016"
}

func (es *ExpectedStatement) GetAlertType() AlertType {
	return Error
}
