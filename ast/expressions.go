package ast

import "hybroid/tokens"

type EnvTypeExpr struct {
	Type  EnvType
	Token tokens.Token
}

func (ete *EnvTypeExpr) GetType() NodeType {
	return EnvironmentTypeExpression
}

func (ete *EnvTypeExpr) GetToken() tokens.Token {
	return ete.Token
}

func (ete *EnvTypeExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type EnvPathExpr struct {
	Path tokens.Token
}

func (epe *EnvPathExpr) GetType() NodeType {
	return EnvironmentPathExpression
}

func (epe *EnvPathExpr) GetToken() tokens.Token {
	return tokens.Token{Lexeme: epe.Path.Lexeme}
}

func (epe *EnvPathExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

func (epe *EnvPathExpr) Combine(token tokens.Token) {
	epe.Path.Lexeme += ":" + token.Lexeme
	epe.Path.Column.End = token.Column.End
	epe.Path.Line.End = token.Line.End
}

type EnvAccessExpr struct {
	PathExpr *EnvPathExpr
	Accessed Node
}

func (eae *EnvAccessExpr) GetType() NodeType {
	return EnvironmentAccessExpression
}

func (eae *EnvAccessExpr) GetToken() tokens.Token {
	return eae.Accessed.GetToken()
}

func (eae *EnvAccessExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type MacroCallExpr struct {
	Caller *CallExpr
}

func (mce *MacroCallExpr) GetType() NodeType {
	return MacroCallExpression
}

func (mce *MacroCallExpr) GetToken() tokens.Token {
	return mce.Caller.GetToken()
}

func (mce *MacroCallExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type LiteralExpr struct {
	Value     string
	ValueType PrimitiveValueType
	Token     tokens.Token
}

func (le *LiteralExpr) GetType() NodeType {
	return LiteralExpression
}

func (le *LiteralExpr) GetToken() tokens.Token {
	return le.Token
}

func (le *LiteralExpr) GetValueType() PrimitiveValueType {
	return le.ValueType
}

type UnaryExpr struct {
	Value     Node
	Operator  tokens.Token
	ValueType PrimitiveValueType
}

func (ue *UnaryExpr) GetType() NodeType {
	return UnaryExpression
}

func (ue *UnaryExpr) GetToken() tokens.Token {
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

func (te *TypeExpr) GetToken() tokens.Token {
	return te.Name.GetToken()
}

func (te *TypeExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type EntityExpr struct {
	Expr               Node
	Type               *TypeExpr
	ConvertedVarName   *tokens.Token
	OfficialEntityType bool
	EntityName         string
	EnvName            string
	Operator           tokens.Token
	Token              tokens.Token
}

func (ge *EntityExpr) GetType() NodeType {
	return GroupingExpression
}

func (ge *EntityExpr) GetToken() tokens.Token {
	return ge.Token
}

func (ge *EntityExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type GroupExpr struct {
	Expr      Node
	ValueType PrimitiveValueType
	Token     tokens.Token
}

func (ge *GroupExpr) GetType() NodeType {
	return GroupingExpression
}

func (ge *GroupExpr) GetToken() tokens.Token {
	return ge.Token
}

func (ge *GroupExpr) GetValueType() PrimitiveValueType {
	return ge.ValueType
}

type BinaryExpr struct {
	Left, Right Node
	Operator    tokens.Token
	ValueType   PrimitiveValueType
}

func (be *BinaryExpr) GetType() NodeType {
	return BinaryExpression
}

func (be *BinaryExpr) GetToken() tokens.Token {
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

func (ce *CallExpr) GetToken() tokens.Token {
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

func (mce *MethodCallExpr) GetToken() tokens.Token {
	return mce.Call.GetToken()
}

func (mce *MethodCallExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type BuiltinExpr struct {
	Name tokens.Token
}

func (ce *BuiltinExpr) GetType() NodeType {
	return BuiltinExpression
}

func (ce *BuiltinExpr) GetToken() tokens.Token {
	return ce.Name
}

func (ce *BuiltinExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type FunctionExpr struct {
	Token  tokens.Token
	Return []*TypeExpr
	Params []Param
	Body   []Node
}

func (fe *FunctionExpr) GetType() NodeType {
	return FunctionExpression
}

func (fe *FunctionExpr) GetToken() tokens.Token {
	return fe.Token
}

func (fe *FunctionExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type StructExpr struct {
	Token  tokens.Token
	Fields []*FieldDecl
}

func (ase *StructExpr) GetType() NodeType {
	return StructExpression
}

func (ase *StructExpr) GetToken() tokens.Token {
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

func (me *MatchExpr) GetToken() tokens.Token {
	return me.MatchStmt.GetToken()
}

func (me *MatchExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type SelfExpr struct {
	Token      tokens.Token
	EntityName string
	Type       SelfExprType
}

func (se *SelfExpr) GetType() NodeType {
	return SelfExpression
}

func (se *SelfExpr) GetToken() tokens.Token {
	return se.Token
}

func (se *SelfExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type NewExpr struct {
	Type     *TypeExpr
	Generics []*TypeExpr
	Args     []Node
	Token    tokens.Token
}

func (ne *NewExpr) GetType() NodeType {
	return NewExpession
}

func (ne *NewExpr) GetToken() tokens.Token {
	return ne.Token
}

func (ne *NewExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type SpawnExpr struct {
	Type     *TypeExpr
	Args     []Node
	Generics []*TypeExpr
	Token    tokens.Token
}

func (ne *SpawnExpr) GetType() NodeType {
	return SpawnExpression
}

func (ne *SpawnExpr) GetToken() tokens.Token {
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

func (fe *FieldExpr) GetToken() tokens.Token {
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

func (me *MemberExpr) GetToken() tokens.Token {
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
	Token tokens.Token
	Map   map[tokens.Token]Property
}

func (me *MapExpr) GetType() NodeType {
	return MapExpression
}

func (me *MapExpr) GetToken() tokens.Token {
	return me.Token
}

func (me *MapExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type ListExpr struct {
	List      []Node
	ValueType PrimitiveValueType
	Token     tokens.Token
}

func (le *ListExpr) GetType() NodeType {
	return ListExpression
}

func (le *ListExpr) GetToken() tokens.Token {
	return le.Token
}

func (le *ListExpr) GetValueType() PrimitiveValueType {
	return le.ValueType
}

type IdentifierExpr struct {
	Name      tokens.Token
	ValueType PrimitiveValueType
}

func (ie *IdentifierExpr) GetType() NodeType {
	return Identifier
}

func (ie *IdentifierExpr) GetToken() tokens.Token {
	return ie.Name
}

func (ie *IdentifierExpr) GetValueType() PrimitiveValueType {
	return ie.ValueType
}

type Improper struct {
	Token tokens.Token
	Type  NodeType
}

func NewImproper(token tokens.Token, nodeType NodeType) *Improper {
	return &Improper{
		Token: token,
		Type:  nodeType,
	}
}

// Make sure to check if node.GetType() == ast.NA before trying to get the node type
func ImproperToNodeType(improper Node) NodeType {
	return improper.(*Improper).Type
}

func (i *Improper) GetType() NodeType {
	return NA
}

func (i *Improper) GetToken() tokens.Token {
	return i.Token
}

func (i *Improper) GetValueType() PrimitiveValueType {
	return Invalid
}
