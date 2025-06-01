// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedStatement struct {
	Specifier Snippet
}

func (es *ExpectedStatement) Message() string {
	return "expected statement"
}

func (es *ExpectedStatement) SnippetSpecifier() Snippet {
	return es.Specifier
}

func (es *ExpectedStatement) Note() string {
	return ""
}

func (es *ExpectedStatement) ID() string {
	return "hyb001P"
}

func (es *ExpectedStatement) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
	Specifier Snippet
	Context   string `default:""`
}

func (ee *ExpectedExpression) Message() string {
	return fmt.Sprintf("expected expression %s", ee.Context)
}

func (ee *ExpectedExpression) SnippetSpecifier() Snippet {
	return ee.Specifier
}

func (ee *ExpectedExpression) Note() string {
	return ""
}

func (ee *ExpectedExpression) ID() string {
	return "hyb002P"
}

func (ee *ExpectedExpression) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownStatement struct {
	Specifier Snippet
	Context   string `default:""`
}

func (us *UnknownStatement) Message() string {
	return fmt.Sprintf("unknown statement %s", us.Context)
}

func (us *UnknownStatement) SnippetSpecifier() Snippet {
	return us.Specifier
}

func (us *UnknownStatement) Note() string {
	return ""
}

func (us *UnknownStatement) ID() string {
	return "hyb003P"
}

func (us *UnknownStatement) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedKeyword struct {
	Specifier Snippet
	Keyword   string
	Context   string `default:""`
}

func (ek *ExpectedKeyword) Message() string {
	return fmt.Sprintf("expected keyword '%s' %s", ek.Keyword, ek.Context)
}

func (ek *ExpectedKeyword) SnippetSpecifier() Snippet {
	return ek.Specifier
}

func (ek *ExpectedKeyword) Note() string {
	return ""
}

func (ek *ExpectedKeyword) ID() string {
	return "hyb004P"
}

func (ek *ExpectedKeyword) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
	Specifier Snippet
	Context   string `default:""`
}

func (ei *ExpectedIdentifier) Message() string {
	return fmt.Sprintf("expected identifier %s", ei.Context)
}

func (ei *ExpectedIdentifier) SnippetSpecifier() Snippet {
	return ei.Specifier
}

func (ei *ExpectedIdentifier) Note() string {
	return ""
}

func (ei *ExpectedIdentifier) ID() string {
	return "hyb005P"
}

func (ei *ExpectedIdentifier) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedSymbol struct {
	Specifier Snippet
	Symbol    string
	Context   string `default:""`
}

func (es *ExpectedSymbol) Message() string {
	return fmt.Sprintf("expected '%s' %s", es.Symbol, es.Context)
}

func (es *ExpectedSymbol) SnippetSpecifier() Snippet {
	return es.Specifier
}

func (es *ExpectedSymbol) Note() string {
	return ""
}

func (es *ExpectedSymbol) ID() string {
	return "hyb006P"
}

func (es *ExpectedSymbol) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneElseBlock struct {
	Specifier Snippet
}

func (mtoeb *MoreThanOneElseBlock) Message() string {
	return "cannot have more than one else block in an if statement"
}

func (mtoeb *MoreThanOneElseBlock) SnippetSpecifier() Snippet {
	return mtoeb.Specifier
}

func (mtoeb *MoreThanOneElseBlock) Note() string {
	return ""
}

func (mtoeb *MoreThanOneElseBlock) ID() string {
	return "hyb007P"
}

func (mtoeb *MoreThanOneElseBlock) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneConstructor struct {
	Specifier Snippet
}

func (mtoc *MoreThanOneConstructor) Message() string {
	return "cannot have more than one constructor in class declaration"
}

func (mtoc *MoreThanOneConstructor) SnippetSpecifier() Snippet {
	return mtoc.Specifier
}

func (mtoc *MoreThanOneConstructor) Note() string {
	return ""
}

func (mtoc *MoreThanOneConstructor) ID() string {
	return "hyb008P"
}

func (mtoc *MoreThanOneConstructor) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneEntityFunction struct {
	Specifier    Snippet
	FunctionType string
}

func (mtoef *MoreThanOneEntityFunction) Message() string {
	return fmt.Sprintf("cannot have more than one %s in entity declaration", mtoef.FunctionType)
}

func (mtoef *MoreThanOneEntityFunction) SnippetSpecifier() Snippet {
	return mtoef.Specifier
}

func (mtoef *MoreThanOneEntityFunction) Note() string {
	return ""
}

func (mtoef *MoreThanOneEntityFunction) ID() string {
	return "hyb009P"
}

func (mtoef *MoreThanOneEntityFunction) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultipleIdentifiersInCompoundAssignment struct {
	Specifier Snippet
}

func (miica *MultipleIdentifiersInCompoundAssignment) Message() string {
	return "cannot have more than one left-hand identifier in a compound assignment"
}

func (miica *MultipleIdentifiersInCompoundAssignment) SnippetSpecifier() Snippet {
	return miica.Specifier
}

func (miica *MultipleIdentifiersInCompoundAssignment) Note() string {
	return "compound assignments include +=, -=, *=, /=, etc."
}

func (miica *MultipleIdentifiersInCompoundAssignment) ID() string {
	return "hyb010P"
}

func (miica *MultipleIdentifiersInCompoundAssignment) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ReturnsInConstructor struct {
	Specifier Snippet
}

func (ric *ReturnsInConstructor) Message() string {
	return "cannot have return types in constructor"
}

func (ric *ReturnsInConstructor) SnippetSpecifier() Snippet {
	return ric.Specifier
}

func (ric *ReturnsInConstructor) Note() string {
	return ""
}

func (ric *ReturnsInConstructor) ID() string {
	return "hyb011P"
}

func (ric *ReturnsInConstructor) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironmentPathExpression struct {
	Specifier Snippet
}

func (eepe *ExpectedEnvironmentPathExpression) Message() string {
	return "expected environment path expression"
}

func (eepe *ExpectedEnvironmentPathExpression) SnippetSpecifier() Snippet {
	return eepe.Specifier
}

func (eepe *ExpectedEnvironmentPathExpression) Note() string {
	return ""
}

func (eepe *ExpectedEnvironmentPathExpression) ID() string {
	return "hyb012P"
}

func (eepe *ExpectedEnvironmentPathExpression) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
	Specifier Snippet
	Context   string `default:""`
}

func (et *ExpectedType) Message() string {
	return fmt.Sprintf("expected type %s", et.Context)
}

func (et *ExpectedType) SnippetSpecifier() Snippet {
	return et.Specifier
}

func (et *ExpectedType) Note() string {
	return ""
}

func (et *ExpectedType) ID() string {
	return "hyb013P"
}

func (et *ExpectedType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAssignmentSymbol struct {
	Specifier Snippet
}

func (eas *ExpectedAssignmentSymbol) Message() string {
	return "expected assignment symbol"
}

func (eas *ExpectedAssignmentSymbol) SnippetSpecifier() Snippet {
	return eas.Specifier
}

func (eas *ExpectedAssignmentSymbol) Note() string {
	return "assignment symbols are: '=', '+=', '-=', '*=', '%%=', '/=', '\\='"
}

func (eas *ExpectedAssignmentSymbol) ID() string {
	return "hyb014P"
}

func (eas *ExpectedAssignmentSymbol) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
	Specifier Snippet
}

func (eeob *ExpectedExpressionOrBody) Message() string {
	return "expected expression or body"
}

func (eeob *ExpectedExpressionOrBody) SnippetSpecifier() Snippet {
	return eeob.Specifier
}

func (eeob *ExpectedExpressionOrBody) Note() string {
	return ""
}

func (eeob *ExpectedExpressionOrBody) ID() string {
	return "hyb015P"
}

func (eeob *ExpectedExpressionOrBody) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallArgs struct {
	Specifier Snippet
}

func (eca *ExpectedCallArgs) Message() string {
	return "expected call arguments"
}

func (eca *ExpectedCallArgs) SnippetSpecifier() Snippet {
	return eca.Specifier
}

func (eca *ExpectedCallArgs) Note() string {
	return ""
}

func (eca *ExpectedCallArgs) ID() string {
	return "hyb016P"
}

func (eca *ExpectedCallArgs) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCall struct {
	Specifier Snippet
}

func (ic *InvalidCall) Message() string {
	return "invalid expression to call"
}

func (ic *InvalidCall) SnippetSpecifier() Snippet {
	return ic.Specifier
}

func (ic *InvalidCall) Note() string {
	return ""
}

func (ic *InvalidCall) ID() string {
	return "hyb017P"
}

func (ic *InvalidCall) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallAfterMacroSymbol struct {
	Specifier Snippet
}

func (ecams *ExpectedCallAfterMacroSymbol) Message() string {
	return "expected a macro call after '@'"
}

func (ecams *ExpectedCallAfterMacroSymbol) SnippetSpecifier() Snippet {
	return ecams.Specifier
}

func (ecams *ExpectedCallAfterMacroSymbol) Note() string {
	return ""
}

func (ecams *ExpectedCallAfterMacroSymbol) ID() string {
	return "hyb018P"
}

func (ecams *ExpectedCallAfterMacroSymbol) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFieldDeclaration struct {
	Specifier Snippet
}

func (efd *ExpectedFieldDeclaration) Message() string {
	return "expected field declaration inside struct"
}

func (efd *ExpectedFieldDeclaration) SnippetSpecifier() Snippet {
	return efd.Specifier
}

func (efd *ExpectedFieldDeclaration) Note() string {
	return ""
}

func (efd *ExpectedFieldDeclaration) ID() string {
	return "hyb019P"
}

func (efd *ExpectedFieldDeclaration) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyWrappedType struct {
	Specifier Snippet
}

func (ewt *EmptyWrappedType) Message() string {
	return "wrapped types must not be empty"
}

func (ewt *EmptyWrappedType) SnippetSpecifier() Snippet {
	return ewt.Specifier
}

func (ewt *EmptyWrappedType) Note() string {
	return ""
}

func (ewt *EmptyWrappedType) ID() string {
	return "hyb020P"
}

func (ewt *EmptyWrappedType) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArgs struct {
	Specifier Snippet
}

func (era *ExpectedReturnArgs) Message() string {
	return "expected return arguments after fat arrow (=>)"
}

func (era *ExpectedReturnArgs) SnippetSpecifier() Snippet {
	return era.Specifier
}

func (era *ExpectedReturnArgs) Note() string {
	return ""
}

func (era *ExpectedReturnArgs) ID() string {
	return "hyb021P"
}

func (era *ExpectedReturnArgs) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAccessExpression struct {
	Specifier Snippet
}

func (eae *ExpectedAccessExpression) Message() string {
	return "expected an access expression"
}

func (eae *ExpectedAccessExpression) SnippetSpecifier() Snippet {
	return eae.Specifier
}

func (eae *ExpectedAccessExpression) Note() string {
	return "access expression are: identifier, environment access, self, member and field expressions"
}

func (eae *ExpectedAccessExpression) ID() string {
	return "hyb022P"
}

func (eae *ExpectedAccessExpression) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingIterator struct {
	Specifier Snippet
	Context   string `default:""`
}

func (mi *MissingIterator) Message() string {
	return fmt.Sprintf("missing iterator %s", mi.Context)
}

func (mi *MissingIterator) SnippetSpecifier() Snippet {
	return mi.Specifier
}

func (mi *MissingIterator) Note() string {
	return ""
}

func (mi *MissingIterator) ID() string {
	return "hyb023P"
}

func (mi *MissingIterator) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateKeyword struct {
	Specifier Snippet
	Keyword   string
}

func (dk *DuplicateKeyword) Message() string {
	return fmt.Sprintf("cannot have multiple '%s' keywords", dk.Keyword)
}

func (dk *DuplicateKeyword) SnippetSpecifier() Snippet {
	return dk.Specifier
}

func (dk *DuplicateKeyword) Note() string {
	return ""
}

func (dk *DuplicateKeyword) ID() string {
	return "hyb024P"
}

func (dk *DuplicateKeyword) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnexpectedKeyword struct {
	Specifier Snippet
	Keyword   string
	Context   string `default:""`
}

func (uk *UnexpectedKeyword) Message() string {
	return fmt.Sprintf("unexpected keyword '%s' %s", uk.Keyword, uk.Context)
}

func (uk *UnexpectedKeyword) SnippetSpecifier() Snippet {
	return uk.Specifier
}

func (uk *UnexpectedKeyword) Note() string {
	return ""
}

func (uk *UnexpectedKeyword) ID() string {
	return "hyb025P"
}

func (uk *UnexpectedKeyword) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type IteratorRedefinition struct {
	Specifier Snippet
	Context   string `default:""`
}

func (ir *IteratorRedefinition) Message() string {
	return fmt.Sprintf("redefinition of iterator %s", ir.Context)
}

func (ir *IteratorRedefinition) SnippetSpecifier() Snippet {
	return ir.Specifier
}

func (ir *IteratorRedefinition) Note() string {
	return ""
}

func (ir *IteratorRedefinition) ID() string {
	return "hyb026P"
}

func (ir *IteratorRedefinition) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ElseIfBlockAfterElseBlock struct {
	Specifier Snippet
}

func (eibaeb *ElseIfBlockAfterElseBlock) Message() string {
	return "cannot have an else if block after an else block"
}

func (eibaeb *ElseIfBlockAfterElseBlock) SnippetSpecifier() Snippet {
	return eibaeb.Specifier
}

func (eibaeb *ElseIfBlockAfterElseBlock) Note() string {
	return ""
}

func (eibaeb *ElseIfBlockAfterElseBlock) ID() string {
	return "hyb027P"
}

func (eibaeb *ElseIfBlockAfterElseBlock) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneDefaultCase struct {
	Specifier Snippet
}

func (mtodc *MoreThanOneDefaultCase) Message() string {
	return "cannot have more than one default case in match statement"
}

func (mtodc *MoreThanOneDefaultCase) SnippetSpecifier() Snippet {
	return mtodc.Specifier
}

func (mtodc *MoreThanOneDefaultCase) Note() string {
	return ""
}

func (mtodc *MoreThanOneDefaultCase) ID() string {
	return "hyb028P"
}

func (mtodc *MoreThanOneDefaultCase) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnumVariantName struct {
	Specifier Snippet
}

func (ievn *InvalidEnumVariantName) Message() string {
	return "enum variant name must be an identifier"
}

func (ievn *InvalidEnumVariantName) SnippetSpecifier() Snippet {
	return ievn.Specifier
}

func (ievn *InvalidEnumVariantName) Note() string {
	return ""
}

func (ievn *InvalidEnumVariantName) ID() string {
	return "hyb029P"
}

func (ievn *InvalidEnumVariantName) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidExpression struct {
	Specifier Snippet
	Type      string
	Context   string `default:""`
}

func (ie *InvalidExpression) Message() string {
	return fmt.Sprintf("'%s' not allowed %s", ie.Type, ie.Context)
}

func (ie *InvalidExpression) SnippetSpecifier() Snippet {
	return ie.Specifier
}

func (ie *InvalidExpression) Note() string {
	return ""
}

func (ie *InvalidExpression) ID() string {
	return "hyb030P"
}

func (ie *InvalidExpression) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type SyntaxIncoherency struct {
	Specifier       Snippet
	ParsedSection   string
	PreviousSection string
	AllowsNextLine  bool
}

func (si *SyntaxIncoherency) Message() string {
	return fmt.Sprintf("'%s' needs to start in the same%s line as '%s'", si.ParsedSection, func(cond bool, str string) string {
		if !cond {
			return ""
		}
		return str
	}(si.AllowsNextLine, " or next"), si.PreviousSection)
}

func (si *SyntaxIncoherency) SnippetSpecifier() Snippet {
	return si.Specifier
}

func (si *SyntaxIncoherency) Note() string {
	return ""
}

func (si *SyntaxIncoherency) ID() string {
	return "hyb031P"
}

func (si *SyntaxIncoherency) AlertType() Type {
	return Error
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidMapKey struct {
	Specifier Snippet
}

func (imk *InvalidMapKey) Message() string {
	return "expected a string as a map key"
}

func (imk *InvalidMapKey) SnippetSpecifier() Snippet {
	return imk.Specifier
}

func (imk *InvalidMapKey) Note() string {
	return ""
}

func (imk *InvalidMapKey) ID() string {
	return "hyb032P"
}

func (imk *InvalidMapKey) AlertType() Type {
	return Error
}
