package ast

import (
	"hybroid/lexer"
)

type EnvTypeExpr struct {
	Type  EnvType
	Token lexer.Token
}

func (ete *EnvTypeExpr) GetType() NodeType {
	return EnvironmentTypeExpression
}

func (ete *EnvTypeExpr) GetToken() lexer.Token {
	return ete.Token
}

func (ete *EnvTypeExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

// func (ete *EnvTypeExpr) DrawNode(str *strings.Builder, depth int) *strings.Builder {
// 	str.WriteString(ete.Type)
// }

type EnvPathExpr struct {
	Path lexer.Token
}

func (epe *EnvPathExpr) GetType() NodeType {
	return EnvironmentPathExpression
}

func (epe *EnvPathExpr) GetToken() lexer.Token {
	return lexer.Token{Lexeme: epe.Path.Lexeme}
}

func (epe *EnvPathExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

func (epe *EnvPathExpr) Combine(token lexer.Token) {
	epe.Path.Lexeme += ":" + token.Lexeme
	epe.Path.Location.ColEnd = token.Location.ColEnd
	epe.Path.Location.LineEnd = token.Location.LineEnd
}

type EnvAccessExpr struct {
	PathExpr *EnvPathExpr
	Accessed Node
}

func (eae *EnvAccessExpr) GetType() NodeType {
	return EnvironmentAccessExpression
}

func (eae *EnvAccessExpr) GetToken() lexer.Token {
	return eae.Accessed.GetToken()
}

func (eae *EnvAccessExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type MacroCallExpr struct {
	Caller *CallExpr
}

func (self *MacroCallExpr) GetType() NodeType {
	return MacroCallExpression
}

func (self *MacroCallExpr) GetToken() lexer.Token {
	return self.Caller.GetToken()
}

func (self *MacroCallExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

// type CastExpr struct{
// 	Value Node
// 	Type  *TypeExpr
// }

// func (le *CastExpr) GetType() NodeType {
// 	return LiteralExpression
// }

// func (le *CastExpr) GetToken() lexer.Token {
// 	return le.Value.GetToken()
// }

// func (le *CastExpr) GetValueType() PrimitiveValueType {
// 	return le.Value.GetValueType()
// }

type LiteralExpr struct {
	Value     string
	ValueType PrimitiveValueType
	Token     lexer.Token
}

func (le *LiteralExpr) GetType() NodeType {
	return LiteralExpression
}

func (le *LiteralExpr) GetToken() lexer.Token {
	return le.Token
}

func (le *LiteralExpr) GetValueType() PrimitiveValueType {
	return le.ValueType
}

type UnaryExpr struct {
	Value     Node
	Operator  lexer.Token
	ValueType PrimitiveValueType
}

func (ue *UnaryExpr) GetType() NodeType {
	return UnaryExpression
}

func (ue *UnaryExpr) GetToken() lexer.Token {
	return ue.Value.GetToken()
}

func (ue *UnaryExpr) GetValueType() PrimitiveValueType {
	return ue.ValueType
}

type TypeExpr struct {
	WrappedType *TypeExpr
	Name        Node
	Params      []*TypeExpr
	Returns     []*TypeExpr
	Fields      []Param
	IsVariadic  bool
}

func (te *TypeExpr) GetType() NodeType {
	return TypeExpression
}

func (te *TypeExpr) GetToken() lexer.Token {
	return te.Name.GetToken()
}

func (te *TypeExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type EntityExpr struct {
	Expr Node
	Type *TypeExpr
	ConvertedVarName *lexer.Token
	OfficialEntityType bool
	EntityName string
	EnvName string
	Operator lexer.Token
	Token lexer.Token
}

func (ge *EntityExpr) GetType() NodeType {
	return GroupingExpression
}

func (ge *EntityExpr) GetToken() lexer.Token {
	return ge.Token
}

func (ge *EntityExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type GroupExpr struct {
	Expr      Node
	ValueType PrimitiveValueType
	Token     lexer.Token
}

func (ge *GroupExpr) GetType() NodeType {
	return GroupingExpression
}

func (ge *GroupExpr) GetToken() lexer.Token {
	return ge.Token
}

func (ge *GroupExpr) GetValueType() PrimitiveValueType {
	return ge.ValueType
}

type BinaryExpr struct {
	Left, Right Node
	Operator    lexer.Token
	ValueType   PrimitiveValueType
}

func (be *BinaryExpr) GetType() NodeType {
	return BinaryExpression
}

func (be *BinaryExpr) GetToken() lexer.Token {
	return be.Operator
}

func (be *BinaryExpr) GetValueType() PrimitiveValueType {
	return be.ValueType
}

type CallExpr struct {
	Caller       Node
	GenericArgs  []*TypeExpr
	Args         []Node
	ReturnAmount int
}

func (ce *CallExpr) GetType() NodeType {
	return CallExpression
}

func (ce *CallExpr) GetToken() lexer.Token {
	return ce.Caller.GetToken()
}

func (ce *CallExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type MethodCallExpr struct {
	EnvName    string
	TypeName   string
	ExprType   SelfExprType
	Identifier Node
	Call       *CallExpr
	MethodName string
}

func (mce *MethodCallExpr) GetType() NodeType {
	return MethodCallExpression
}

func (mce *MethodCallExpr) GetToken() lexer.Token {
	return mce.Call.GetToken()
}

func (mce *MethodCallExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type BuiltinExpr struct {
	Name lexer.Token
}

func (ce *BuiltinExpr) GetType() NodeType {
	return BuiltinExpression
}

func (ce *BuiltinExpr) GetToken() lexer.Token {
	return ce.Name
}

func (ce *BuiltinExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type FunctionExpr struct {
	Token  lexer.Token
	Return []*TypeExpr
	Params []Param
	Body   []Node
}

func (fe *FunctionExpr) GetType() NodeType {
	return FunctionExpression
}

func (fe *FunctionExpr) GetToken() lexer.Token {
	return fe.Token
}

func (fe *FunctionExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type StructExpr struct {
	Token  lexer.Token
	Fields []*FieldDeclarationStmt
}

func (ase *StructExpr) GetType() NodeType {
	return StructExpression
}

func (ase *StructExpr) GetToken() lexer.Token {
	return ase.Token
}

func (ase *StructExpr) GetValueType() PrimitiveValueType {
	return Struct
}

type MatchExpr struct {
	MatchStmt    MatchStmt
	ReturnAmount int
}

func (me *MatchExpr) GetType() NodeType {
	return MatchExpression
}

func (me *MatchExpr) GetToken() lexer.Token {
	return me.MatchStmt.GetToken()
}

func (me *MatchExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type SelfExpr struct {
	Token      lexer.Token
	EntityName string
	Type       SelfExprType
}

func (se *SelfExpr) GetType() NodeType {
	return SelfExpression
}

func (se *SelfExpr) GetToken() lexer.Token {
	return se.Token
}

func (se *SelfExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type NewExpr struct {
	Type     *TypeExpr
	Generics []*TypeExpr
	Args     []Node
	Token    lexer.Token
}

func (ne *NewExpr) GetType() NodeType {
	return NewExpession
}

func (ne *NewExpr) GetToken() lexer.Token {
	return ne.Token
}

func (ne *NewExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type SpawnExpr struct {
	Type     *TypeExpr
	Args     []Node
	Generics []*TypeExpr
	Token    lexer.Token
}

func (ne *SpawnExpr) GetType() NodeType {
	return SpawnExpression
}

func (ne *SpawnExpr) GetToken() lexer.Token {
	return ne.Token
}

func (ne *SpawnExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type FieldExpr struct {
	Property           Node
	PropertyIdentifier Node
	Identifier         Node
	ExprType           SelfExprType
	EnvName            string
	EntityName         string
	Index              int
}

func (fe *FieldExpr) GetType() NodeType {
	return FieldExpression
}

func (fe *FieldExpr) GetToken() lexer.Token {
	return fe.Identifier.GetToken()
}

func (fe *FieldExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

func (fe *FieldExpr) GetIdentifier() Node {
	return fe.Identifier
}
func (fe *FieldExpr) SetProperty(prop Node) {
	fe.Property = prop
}

func (fe *FieldExpr) SetPropertyIdentifier(ident Node) {
	fe.PropertyIdentifier = ident
}

func (fe *FieldExpr) GetPropertyIdentifier() Node {
	return fe.PropertyIdentifier
}

func (fe *FieldExpr) GetProperty() Node {
	return fe.Property
}

func (fe *FieldExpr) SetIdentifier(ident Node) {
	fe.Identifier = ident
}

type MemberExpr struct {
	Property           Node
	PropertyIdentifier Node
	Identifier         Node
	IsList             bool
}

func (me *MemberExpr) GetType() NodeType {
	return MemberExpression
}

func (me *MemberExpr) GetToken() lexer.Token {
	return me.Identifier.GetToken()
}

func (me *MemberExpr) GetValueType() PrimitiveValueType {
	return me.Identifier.GetValueType()
}

func (me *MemberExpr) GetIdentifier() Node {
	return me.Identifier
}

func (me *MemberExpr) SetProperty(prop Node) {
	me.Property = prop
}

func (me *MemberExpr) SetPropertyIdentifier(ident Node) {
	me.PropertyIdentifier = ident
}

func (me *MemberExpr) GetPropertyIdentifier() Node {
	return me.PropertyIdentifier
}

func (me *MemberExpr) GetProperty() Node {
	return me.Property
}

func (me *MemberExpr) SetIdentifier(ident Node) {
	me.Identifier = ident
}

type Property struct {
	Expr Node
	Type PrimitiveValueType
}

type MapExpr struct {
	Token lexer.Token
	Map   map[lexer.Token]Property
}

func (me *MapExpr) GetType() NodeType {
	return MapExpression
}

func (me *MapExpr) GetToken() lexer.Token {
	return me.Token
}

func (me *MapExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type ListExpr struct {
	List      []Node
	ValueType PrimitiveValueType
	Token     lexer.Token
}

func (le *ListExpr) GetType() NodeType {
	return ListExpression
}

func (le *ListExpr) GetToken() lexer.Token {
	return le.Token
}

func (le *ListExpr) GetValueType() PrimitiveValueType {
	return le.ValueType
}

type IdentifierExpr struct {
	Name      lexer.Token
	ValueType PrimitiveValueType
}

func (ie *IdentifierExpr) GetType() NodeType {
	return Identifier
}

func (ie *IdentifierExpr) GetToken() lexer.Token {
	return ie.Name
}

func (ie *IdentifierExpr) GetValueType() PrimitiveValueType {
	return ie.ValueType
}

type Improper struct {
	Token lexer.Token
}

func NewImproper(token lexer.Token) *Improper {
	return &Improper{
		Token: token,
	}
}

func (i *Improper) GetType() NodeType {
	return NA
}

func (i *Improper) GetToken() lexer.Token {
	return i.Token
}

func (i *Improper) GetValueType() PrimitiveValueType {
	return Invalid
}
