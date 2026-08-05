// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
	"fmt"
	"hybroid/core"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedStatement struct {
	SnippetProvider
}

func NewExpectedStatement(span core.Span) *ExpectedStatement {
	return &ExpectedStatement{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (es *ExpectedStatement) Message() string { return "expected statement" }
func (es *ExpectedStatement) Note() string    { return "" }
func (es *ExpectedStatement) Type() Type      { return Error }
func (es *ExpectedStatement) ID() string      { return "hyb001P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpression struct {
	SnippetProvider
	ContextProvider
}

func NewExpectedExpression(span core.Span) *ExpectedExpression {
	return &ExpectedExpression{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
	}
}
func (ee *ExpectedExpression) Message() string {
	return fmt.Sprintf("expected expression %s", ee.context)
}
func (ee *ExpectedExpression) Note() string { return "" }
func (ee *ExpectedExpression) Type() Type   { return Error }
func (ee *ExpectedExpression) ID() string   { return "hyb002P" }
func (ee *ExpectedExpression) WithContext(ctx string) *ExpectedExpression {
	ee.context = ctx
	return ee
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnknownStatement struct {
	SnippetProvider
	ContextProvider
}

func NewUnknownStatement(span core.Span) *UnknownStatement {
	return &UnknownStatement{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
	}
}
func (us *UnknownStatement) Message() string { return fmt.Sprintf("unknown statement %s", us.context) }
func (us *UnknownStatement) Note() string    { return "" }
func (us *UnknownStatement) Type() Type      { return Error }
func (us *UnknownStatement) ID() string      { return "hyb003P" }
func (us *UnknownStatement) WithContext(ctx string) *UnknownStatement {
	us.context = ctx
	return us
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedKeyword struct {
	SnippetProvider
	ContextProvider
	keyword string
}

func NewExpectedKeyword(span core.Span, keyword string) *ExpectedKeyword {
	return &ExpectedKeyword{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
		keyword:         keyword,
	}
}
func (ek *ExpectedKeyword) Message() string {
	return fmt.Sprintf("expected keyword '%s' %s", ek.keyword, ek.context)
}
func (ek *ExpectedKeyword) Note() string { return "" }
func (ek *ExpectedKeyword) Type() Type   { return Error }
func (ek *ExpectedKeyword) ID() string   { return "hyb004P" }
func (ek *ExpectedKeyword) WithContext(ctx string) *ExpectedKeyword {
	ek.context = ctx
	return ek
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedIdentifier struct {
	SnippetProvider
	ContextProvider
}

func NewExpectedIdentifier(span core.Span) *ExpectedIdentifier {
	return &ExpectedIdentifier{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
	}
}
func (ei *ExpectedIdentifier) Message() string {
	return fmt.Sprintf("expected identifier %s", ei.context)
}
func (ei *ExpectedIdentifier) Note() string { return "" }
func (ei *ExpectedIdentifier) Type() Type   { return Error }
func (ei *ExpectedIdentifier) ID() string   { return "hyb005P" }
func (ei *ExpectedIdentifier) WithContext(ctx string) *ExpectedIdentifier {
	ei.context = ctx
	return ei
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedSymbol struct {
	SnippetProvider
	ContextProvider
	symbol string
}

func NewExpectedSymbol(span core.Span, symbol string) *ExpectedSymbol {
	return &ExpectedSymbol{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
		symbol:          symbol,
	}
}
func (es *ExpectedSymbol) Message() string {
	return fmt.Sprintf("expected '%s' %s", es.symbol, es.context)
}
func (es *ExpectedSymbol) Note() string { return "" }
func (es *ExpectedSymbol) Type() Type   { return Error }
func (es *ExpectedSymbol) ID() string   { return "hyb006P" }
func (es *ExpectedSymbol) WithContext(ctx string) *ExpectedSymbol {
	es.context = ctx
	return es
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneElseBlock struct {
	SnippetProvider
}

func NewMoreThanOneElseBlock(span core.Span) *MoreThanOneElseBlock {
	return &MoreThanOneElseBlock{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (mtoeb *MoreThanOneElseBlock) Message() string {
	return "cannot have more than one else block in an if statement"
}
func (mtoeb *MoreThanOneElseBlock) Note() string { return "" }
func (mtoeb *MoreThanOneElseBlock) Type() Type   { return Error }
func (mtoeb *MoreThanOneElseBlock) ID() string   { return "hyb007P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneConstructor struct {
	SnippetProvider
}

func NewMoreThanOneConstructor(span core.Span) *MoreThanOneConstructor {
	return &MoreThanOneConstructor{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (mtoc *MoreThanOneConstructor) Message() string {
	return "cannot have more than one constructor in class declaration"
}
func (mtoc *MoreThanOneConstructor) Note() string { return "" }
func (mtoc *MoreThanOneConstructor) Type() Type   { return Error }
func (mtoc *MoreThanOneConstructor) ID() string   { return "hyb008P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneEntityFunction struct {
	SnippetProvider
	functionType string
}

func NewMoreThanOneEntityFunction(span core.Span, functionType string) *MoreThanOneEntityFunction {
	return &MoreThanOneEntityFunction{
		SnippetProvider: SnippetProvider{span: span},
		functionType:    functionType,
	}
}
func (mtoef *MoreThanOneEntityFunction) Message() string {
	return fmt.Sprintf("cannot have more than one %s in entity declaration", mtoef.functionType)
}
func (mtoef *MoreThanOneEntityFunction) Note() string { return "" }
func (mtoef *MoreThanOneEntityFunction) Type() Type   { return Error }
func (mtoef *MoreThanOneEntityFunction) ID() string   { return "hyb009P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MultipleIdentifiersInCompoundAssignment struct {
	SnippetProvider
}

func NewMultipleIdentifiersInCompoundAssignment(span core.Span) *MultipleIdentifiersInCompoundAssignment {
	return &MultipleIdentifiersInCompoundAssignment{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (miica *MultipleIdentifiersInCompoundAssignment) Message() string {
	return "cannot have more than one left-hand identifier in a compound assignment"
}
func (miica *MultipleIdentifiersInCompoundAssignment) Note() string {
	return "compound assignments include +=, -=, *=, /=, etc."
}
func (miica *MultipleIdentifiersInCompoundAssignment) Type() Type { return Error }
func (miica *MultipleIdentifiersInCompoundAssignment) ID() string { return "hyb010P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ReturnsInConstructor struct {
	SnippetProvider
}

func NewReturnsInConstructor(span core.Span) *ReturnsInConstructor {
	return &ReturnsInConstructor{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ric *ReturnsInConstructor) Message() string { return "cannot have return types in constructor" }
func (ric *ReturnsInConstructor) Note() string    { return "" }
func (ric *ReturnsInConstructor) Type() Type      { return Error }
func (ric *ReturnsInConstructor) ID() string      { return "hyb011P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedEnvironmentPathExpression struct {
	SnippetProvider
}

func NewExpectedEnvironmentPathExpression(span core.Span) *ExpectedEnvironmentPathExpression {
	return &ExpectedEnvironmentPathExpression{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eepe *ExpectedEnvironmentPathExpression) Message() string {
	return "expected environment path expression"
}
func (eepe *ExpectedEnvironmentPathExpression) Note() string { return "" }
func (eepe *ExpectedEnvironmentPathExpression) Type() Type   { return Error }
func (eepe *ExpectedEnvironmentPathExpression) ID() string   { return "hyb012P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedType struct {
	SnippetProvider
	ContextProvider
}

func NewExpectedType(span core.Span) *ExpectedType {
	return &ExpectedType{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
	}
}
func (et *ExpectedType) Message() string { return fmt.Sprintf("expected type %s", et.context) }
func (et *ExpectedType) Note() string    { return "" }
func (et *ExpectedType) Type() Type      { return Error }
func (et *ExpectedType) ID() string      { return "hyb013P" }
func (et *ExpectedType) WithContext(ctx string) *ExpectedType {
	et.context = ctx
	return et
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAssignmentSymbol struct {
	SnippetProvider
}

func NewExpectedAssignmentSymbol(span core.Span) *ExpectedAssignmentSymbol {
	return &ExpectedAssignmentSymbol{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eas *ExpectedAssignmentSymbol) Message() string { return "expected assignment symbol" }
func (eas *ExpectedAssignmentSymbol) Note() string {
	return "assignment symbols are: '=', '+=', '-=', '*=', '%%=', '/=', '\\='"
}
func (eas *ExpectedAssignmentSymbol) Type() Type { return Error }
func (eas *ExpectedAssignmentSymbol) ID() string { return "hyb014P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedExpressionOrBody struct {
	SnippetProvider
}

func NewExpectedExpressionOrBody(span core.Span) *ExpectedExpressionOrBody {
	return &ExpectedExpressionOrBody{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eeob *ExpectedExpressionOrBody) Message() string { return "expected expression or body" }
func (eeob *ExpectedExpressionOrBody) Note() string    { return "" }
func (eeob *ExpectedExpressionOrBody) Type() Type      { return Error }
func (eeob *ExpectedExpressionOrBody) ID() string      { return "hyb015P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallArgs struct {
	SnippetProvider
}

func NewExpectedCallArgs(span core.Span) *ExpectedCallArgs {
	return &ExpectedCallArgs{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eca *ExpectedCallArgs) Message() string { return "expected call arguments" }
func (eca *ExpectedCallArgs) Note() string    { return "" }
func (eca *ExpectedCallArgs) Type() Type      { return Error }
func (eca *ExpectedCallArgs) ID() string      { return "hyb016P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidCall struct {
	SnippetProvider
}

func NewInvalidCall(span core.Span) *InvalidCall {
	return &InvalidCall{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ic *InvalidCall) Message() string { return "invalid expression to call" }
func (ic *InvalidCall) Note() string    { return "" }
func (ic *InvalidCall) Type() Type      { return Error }
func (ic *InvalidCall) ID() string      { return "hyb017P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedCallAfterMacroSymbol struct {
	SnippetProvider
}

func NewExpectedCallAfterMacroSymbol(span core.Span) *ExpectedCallAfterMacroSymbol {
	return &ExpectedCallAfterMacroSymbol{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ecams *ExpectedCallAfterMacroSymbol) Message() string { return "expected a macro call after '@'" }
func (ecams *ExpectedCallAfterMacroSymbol) Note() string    { return "" }
func (ecams *ExpectedCallAfterMacroSymbol) Type() Type      { return Error }
func (ecams *ExpectedCallAfterMacroSymbol) ID() string      { return "hyb018P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedFieldDeclaration struct {
	SnippetProvider
}

func NewExpectedFieldDeclaration(span core.Span) *ExpectedFieldDeclaration {
	return &ExpectedFieldDeclaration{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (efd *ExpectedFieldDeclaration) Message() string {
	return "expected field declaration inside struct"
}
func (efd *ExpectedFieldDeclaration) Note() string { return "" }
func (efd *ExpectedFieldDeclaration) Type() Type   { return Error }
func (efd *ExpectedFieldDeclaration) ID() string   { return "hyb019P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type EmptyWrappedType struct {
	SnippetProvider
}

func NewEmptyWrappedType(span core.Span) *EmptyWrappedType {
	return &EmptyWrappedType{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ewt *EmptyWrappedType) Message() string { return "wrapped types must not be empty" }
func (ewt *EmptyWrappedType) Note() string    { return "" }
func (ewt *EmptyWrappedType) Type() Type      { return Error }
func (ewt *EmptyWrappedType) ID() string      { return "hyb020P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedReturnArgs struct {
	SnippetProvider
}

func NewExpectedReturnArgs(span core.Span) *ExpectedReturnArgs {
	return &ExpectedReturnArgs{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (era *ExpectedReturnArgs) Message() string {
	return "expected return arguments after fat arrow (=>)"
}
func (era *ExpectedReturnArgs) Note() string { return "" }
func (era *ExpectedReturnArgs) Type() Type   { return Error }
func (era *ExpectedReturnArgs) ID() string   { return "hyb021P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ExpectedAccessExpression struct {
	SnippetProvider
}

func NewExpectedAccessExpression(span core.Span) *ExpectedAccessExpression {
	return &ExpectedAccessExpression{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eae *ExpectedAccessExpression) Message() string { return "expected an access expression" }
func (eae *ExpectedAccessExpression) Note() string {
	return "access expression are: identifier, environment access, self, member and field expressions"
}
func (eae *ExpectedAccessExpression) Type() Type { return Error }
func (eae *ExpectedAccessExpression) ID() string { return "hyb022P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MissingIterator struct {
	SnippetProvider
	ContextProvider
}

func NewMissingIterator(span core.Span) *MissingIterator {
	return &MissingIterator{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
	}
}
func (mi *MissingIterator) Message() string { return fmt.Sprintf("missing iterator %s", mi.context) }
func (mi *MissingIterator) Note() string    { return "" }
func (mi *MissingIterator) Type() Type      { return Error }
func (mi *MissingIterator) ID() string      { return "hyb023P" }
func (mi *MissingIterator) WithContext(ctx string) *MissingIterator {
	mi.context = ctx
	return mi
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type DuplicateKeyword struct {
	SnippetProvider
	keyword string
}

func NewDuplicateKeyword(span core.Span, keyword string) *DuplicateKeyword {
	return &DuplicateKeyword{
		SnippetProvider: SnippetProvider{span: span},
		keyword:         keyword,
	}
}
func (dk *DuplicateKeyword) Message() string {
	return fmt.Sprintf("cannot have multiple '%s' keywords", dk.keyword)
}
func (dk *DuplicateKeyword) Note() string { return "" }
func (dk *DuplicateKeyword) Type() Type   { return Error }
func (dk *DuplicateKeyword) ID() string   { return "hyb024P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type UnexpectedKeyword struct {
	SnippetProvider
	ContextProvider
	keyword string
}

func NewUnexpectedKeyword(span core.Span, keyword string) *UnexpectedKeyword {
	return &UnexpectedKeyword{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
		keyword:         keyword,
	}
}
func (uk *UnexpectedKeyword) Message() string {
	return fmt.Sprintf("unexpected keyword '%s' %s", uk.keyword, uk.context)
}
func (uk *UnexpectedKeyword) Note() string { return "" }
func (uk *UnexpectedKeyword) Type() Type   { return Error }
func (uk *UnexpectedKeyword) ID() string   { return "hyb025P" }
func (uk *UnexpectedKeyword) WithContext(ctx string) *UnexpectedKeyword {
	uk.context = ctx
	return uk
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type IteratorRedefinition struct {
	SnippetProvider
	ContextProvider
}

func NewIteratorRedefinition(span core.Span) *IteratorRedefinition {
	return &IteratorRedefinition{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
	}
}
func (ir *IteratorRedefinition) Message() string {
	return fmt.Sprintf("redefinition of iterator %s", ir.context)
}
func (ir *IteratorRedefinition) Note() string { return "" }
func (ir *IteratorRedefinition) Type() Type   { return Error }
func (ir *IteratorRedefinition) ID() string   { return "hyb026P" }
func (ir *IteratorRedefinition) WithContext(ctx string) *IteratorRedefinition {
	ir.context = ctx
	return ir
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type ElseIfBlockAfterElseBlock struct {
	SnippetProvider
}

func NewElseIfBlockAfterElseBlock(span core.Span) *ElseIfBlockAfterElseBlock {
	return &ElseIfBlockAfterElseBlock{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (eibaeb *ElseIfBlockAfterElseBlock) Message() string {
	return "cannot have an else if block after an else block"
}
func (eibaeb *ElseIfBlockAfterElseBlock) Note() string { return "" }
func (eibaeb *ElseIfBlockAfterElseBlock) Type() Type   { return Error }
func (eibaeb *ElseIfBlockAfterElseBlock) ID() string   { return "hyb027P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type MoreThanOneDefaultCase struct {
	SnippetProvider
}

func NewMoreThanOneDefaultCase(span core.Span) *MoreThanOneDefaultCase {
	return &MoreThanOneDefaultCase{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (mtodc *MoreThanOneDefaultCase) Message() string {
	return "cannot have more than one default case in match statement"
}
func (mtodc *MoreThanOneDefaultCase) Note() string { return "" }
func (mtodc *MoreThanOneDefaultCase) Type() Type   { return Error }
func (mtodc *MoreThanOneDefaultCase) ID() string   { return "hyb028P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidEnumVariantName struct {
	SnippetProvider
}

func NewInvalidEnumVariantName(span core.Span) *InvalidEnumVariantName {
	return &InvalidEnumVariantName{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (ievn *InvalidEnumVariantName) Message() string {
	return "enum variant name must be an identifier"
}
func (ievn *InvalidEnumVariantName) Note() string { return "" }
func (ievn *InvalidEnumVariantName) Type() Type   { return Error }
func (ievn *InvalidEnumVariantName) ID() string   { return "hyb029P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidExpression struct {
	SnippetProvider
	ContextProvider
	exprType string
}

func NewInvalidExpression(span core.Span, exprType string) *InvalidExpression {
	return &InvalidExpression{
		SnippetProvider: SnippetProvider{span: span},
		ContextProvider: ContextProvider{context: ""}, // Default value
		exprType:        exprType,
	}
}
func (ie *InvalidExpression) Message() string {
	return fmt.Sprintf("'%s' not allowed %s", ie.exprType, ie.context)
}
func (ie *InvalidExpression) Note() string { return "" }
func (ie *InvalidExpression) Type() Type   { return Error }
func (ie *InvalidExpression) ID() string   { return "hyb030P" }
func (ie *InvalidExpression) WithContext(ctx string) *InvalidExpression {
	ie.context = ctx
	return ie
}

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type SyntaxIncoherency struct {
	SnippetProvider
	parsedSection   string
	previousSection string
	allowsNextLine  bool
}

func NewSyntaxIncoherency(span core.Span, parsedSection string, previousSection string, allowsNextLine bool) *SyntaxIncoherency {
	return &SyntaxIncoherency{
		SnippetProvider: SnippetProvider{span: span},
		parsedSection:   parsedSection,
		previousSection: previousSection,
		allowsNextLine:  allowsNextLine,
	}
}
func (si *SyntaxIncoherency) Message() string {
	return fmt.Sprintf("'%s' needs to start in the same%s line as '%s'", si.parsedSection, func(cond bool, str string) string {
		if !cond {
			return ""
		}
		return str
	}(si.allowsNextLine, " or next"), si.previousSection)
}
func (si *SyntaxIncoherency) Note() string { return "" }
func (si *SyntaxIncoherency) Type() Type   { return Error }
func (si *SyntaxIncoherency) ID() string   { return "hyb031P" }

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
type InvalidMapKey struct {
	SnippetProvider
}

func NewInvalidMapKey(span core.Span) *InvalidMapKey {
	return &InvalidMapKey{
		SnippetProvider: SnippetProvider{span: span},
	}
}
func (imk *InvalidMapKey) Message() string { return "expected a string as a map key" }
func (imk *InvalidMapKey) Note() string    { return "" }
func (imk *InvalidMapKey) Type() Type      { return Error }
func (imk *InvalidMapKey) ID() string      { return "hyb032P" }
