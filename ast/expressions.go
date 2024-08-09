package ast

import (
	"hybroid/lexer"
	"strings"
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
	SubPaths []string
}

func (epe *EnvPathExpr) GetType() NodeType {
	return EnvironmentPathExpression
}

func (epe *EnvPathExpr) GetToken() lexer.Token {
	return lexer.Token{Lexeme: epe.SubPaths[0]}
}

func (epe *EnvPathExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

func (epe *EnvPathExpr) Nameify() string {
	return strings.Join(epe.SubPaths, "::")
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
	Caller Node
	GenericArgs []*TypeExpr
	Args   []Node
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

type PewpewExpr struct {
	Node Node
}

func (ce *PewpewExpr) GetType() NodeType {
	return PewpewExpression
}

func (ce *PewpewExpr) GetToken() lexer.Token {
	return ce.Node.GetToken()
}

func (ce *PewpewExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type FmathExpr struct {
	Node Node
}

func (ce *FmathExpr) GetType() NodeType {
	return FmathExpression
}

func (ce *FmathExpr) GetToken() lexer.Token {
	return ce.Node.GetToken()
}

func (ce *FmathExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type BuiltinCallExpr struct {
	Name lexer.Token
	Args []Node
}

func (ce *BuiltinCallExpr) GetType() NodeType {
	return BuiltinCallExpression
}

func (ce *BuiltinCallExpr) GetToken() lexer.Token {
	return ce.Name
}

func (ce *BuiltinCallExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type StandardExpr struct {
	Library StandardLibrary
	Node    Node
}

func (ce *StandardExpr) GetType() NodeType {
	return StandardExpression
}

func (ce *StandardExpr) GetToken() lexer.Token {
	return ce.Node.GetToken()
}

func (ce *StandardExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type AnonFnExpr struct {
	Token  lexer.Token
	Return []*TypeExpr
	Params []Param
	Body   []Node
}

func (afe *AnonFnExpr) GetType() NodeType {
	return AnonymousFunctionExpression
}

func (afe *AnonFnExpr) GetToken() lexer.Token {
	return afe.Token
}

func (afe *AnonFnExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type AnonStructExpr struct {
	Token  lexer.Token
	Fields []*FieldDeclarationStmt
}

func (ase *AnonStructExpr) GetType() NodeType {
	return AnonymousStructExpression
}

func (ase *AnonStructExpr) GetToken() lexer.Token {
	return ase.Token
}

func (ase *AnonStructExpr) GetValueType() PrimitiveValueType {
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
	Token lexer.Token
	Type  SelfExprType
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
	Type  *TypeExpr
	Args  []Node
	Token lexer.Token
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
	Type  *TypeExpr
	Args  []Node
	Token lexer.Token
}

func (ne *SpawnExpr) GetType() NodeType {
	return NewExpession
}

func (ne *SpawnExpr) GetToken() lexer.Token {
	return ne.Token
}

func (ne *SpawnExpr) GetValueType() PrimitiveValueType {
	return Invalid
}

type FieldExpr struct {
	Property   Node
	Identifier Node
	Index      int
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

func (fe *FieldExpr) SetProperty(prop Node) {
	fe.Property = prop
}

func (fe *FieldExpr) GetProperty() *Node {
	return &fe.Property
}

func (fe *FieldExpr) SetIdentifier(ident Node) {
	fe.Identifier = ident
}

type MemberExpr struct {
	Property   Node
	Identifier Node
	IsList     bool
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

func (me *MemberExpr) SetProperty(prop Node) {
	me.Property = prop
}

func (me *MemberExpr) GetProperty() *Node {
	return &me.Property
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
