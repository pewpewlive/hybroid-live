// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnclosingMark struct {
  Specifier Multiline
  Mark string
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
  Mark string
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
}

func (ei *ExpectedIdentifier) GetMessage() string {
  return "Expected identifier"
}

func (ei *ExpectedIdentifier) GetSpecifier() SnippetSpecifier {
  return &ei.Specifier
}

func (ei *ExpectedIdentifier) GetNote() string {
  return ""
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
  return "hyb008"
}

func (es *ExpectedStatement) GetAlertType() AlertType {
  return Error
}

