// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedParenthesis struct {
  Token tokens.Token
  Location tokens.TokenLocation
  Symbol string
}

func (ep *ExpectedParenthesis) GetMessage() string {
  return fmt.Sprintf("Expected '%s'", ep.Symbol)
}

func (ep *ExpectedParenthesis) GetTokens() []tokens.Token {
  return []tokens.Token{ep.Token}
}

func (ep *ExpectedParenthesis) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ep.Location}
}

func (ep *ExpectedParenthesis) GetNote() string {
  return ""
}

func (ep *ExpectedParenthesis) GetName() string {
  return "ExpectedParenthesis"
}

func (ep *ExpectedParenthesis) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironment struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (ee *ExpectedEnvironment) GetMessage() string {
  return "Expected environment statement"
}

func (ee *ExpectedEnvironment) GetTokens() []tokens.Token {
  return []tokens.Token{ee.Token}
}

func (ee *ExpectedEnvironment) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ee.Location}
}

func (ee *ExpectedEnvironment) GetNote() string {
  return "environment statement has to be the first statement in any hybroid file. example: env HelloWorld as Level"
}

func (ee *ExpectedEnvironment) GetName() string {
  return "ExpectedEnvironment"
}

func (ee *ExpectedEnvironment) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (ei *ExpectedIdentifier) GetMessage() string {
  return "Expected identifier"
}

func (ei *ExpectedIdentifier) GetTokens() []tokens.Token {
  return []tokens.Token{ei.Token}
}

func (ei *ExpectedIdentifier) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ei.Location}
}

func (ei *ExpectedIdentifier) GetNote() string {
  return ""
}

func (ei *ExpectedIdentifier) GetName() string {
  return "ExpectedIdentifier"
}

func (ei *ExpectedIdentifier) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (et *ExpectedType) GetMessage() string {
  return "Expected type"
}

func (et *ExpectedType) GetTokens() []tokens.Token {
  return []tokens.Token{et.Token}
}

func (et *ExpectedType) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{et.Location}
}

func (et *ExpectedType) GetNote() string {
  return fmt.Sprintf("this needs to be declared with a type. example: number %s", et.Token.Lexeme)
}

func (et *ExpectedType) GetName() string {
  return "ExpectedType"
}

func (et *ExpectedType) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (ee *ExpectedExpression) GetMessage() string {
  return "Expected expression"
}

func (ee *ExpectedExpression) GetTokens() []tokens.Token {
  return []tokens.Token{ee.Token}
}

func (ee *ExpectedExpression) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{ee.Location}
}

func (ee *ExpectedExpression) GetNote() string {
  return ""
}

func (ee *ExpectedExpression) GetName() string {
  return "ExpectedExpression"
}

func (ee *ExpectedExpression) GetAlertType() AlertType {
  return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
  Token tokens.Token
  Location tokens.TokenLocation
}

func (eeob *ExpectedExpressionOrBody) GetMessage() string {
  return "Expected expression"
}

func (eeob *ExpectedExpressionOrBody) GetTokens() []tokens.Token {
  return []tokens.Token{eeob.Token}
}

func (eeob *ExpectedExpressionOrBody) GetLocations() []tokens.TokenLocation {
  return []tokens.TokenLocation{eeob.Location}
}

func (eeob *ExpectedExpressionOrBody) GetNote() string {
  return ""
}

func (eeob *ExpectedExpressionOrBody) GetName() string {
  return "ExpectedExpressionOrBody"
}

func (eeob *ExpectedExpressionOrBody) GetAlertType() AlertType {
  return Error
}

