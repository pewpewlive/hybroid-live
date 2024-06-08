package ast

import (
	"hybroid/lexer"
)

type Accessor interface {
	Node
	GetOwner() *Accessor
	GetProperty() *Node
	SetOwner(owner Accessor) Accessor
	SetProperty(prop Node) Accessor
	SetIdentifier(ident Node) Accessor
}

type EnvType int

const (
	Mesh EnvType = iota
	Level
	Shared
	Sound
	InvalidEnv
)

type EnvTypeExpr struct {
	Type EnvType
	Token lexer.Token
}

func (le *EnvTypeExpr) GetType() NodeType {
	return EnvironmentExpression
}

func (le *EnvTypeExpr) GetToken() lexer.Token {
	return le.Token
}

func (le *EnvTypeExpr) GetValueType() PrimitiveValueType {
	return 0
}

type EnvExpr struct {
	Envs []lexer.Token
}

func (le *EnvExpr) GetType() NodeType {
	return EnvironmentExpression
}

func (le *EnvExpr) GetToken() lexer.Token {
	return le.Envs[0]
}

func (le *EnvExpr) GetValueType() PrimitiveValueType {
	return 0
}

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
	return ue.Operator
}

func (ue *UnaryExpr) GetValueType() PrimitiveValueType {
	return ue.ValueType
}

type TypeExpr struct {
	WrappedType *TypeExpr
	Name        lexer.Token
	Params      *[]*TypeExpr
	Returns     []*TypeExpr
}

func (te *TypeExpr) GetType() NodeType {
	return TypeExpression
}

func (te *TypeExpr) GetToken() lexer.Token {
	return te.Name
}

func (te *TypeExpr) GetValueType() PrimitiveValueType {
	return 0
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
	Identifier string
	Caller     Node
	Args       []Node
	Token      lexer.Token
}

func (ce *CallExpr) GetType() NodeType {
	return CallExpression
}

func (ce *CallExpr) GetToken() lexer.Token {
	return ce.Token
}

func (ce *CallExpr) GetValueType() PrimitiveValueType {
	return 0
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
	return 0
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
	return 0
}

type SelfExprType int

const (
	SelfStruct SelfExprType = iota
	SelfEntity
)

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
	return 0
}

type MethodCallExpr struct {
	TypeName   string
	Owner      Node
	Call       Node
	MethodName string
	Args       []Node
	Token      lexer.Token
}

func (mce *MethodCallExpr) GetType() NodeType {
	return MethodCallExpression
}

func (mce *MethodCallExpr) GetToken() lexer.Token {
	return mce.Token
}

func (mce *MethodCallExpr) GetValueType() PrimitiveValueType {
	return 0
}

type NewExpr struct {
	Type  lexer.Token
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
	return 0
}

type FieldExpr struct {
	Owner      Accessor
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
	return 0
}

func (fe *FieldExpr) SetProperty(prop Node) Accessor {
	fe.Property = prop
	return fe
}

func (fe *FieldExpr) GetProperty() *Node {
	return &fe.Property
}

func (fe *FieldExpr) GetOwner() *Accessor {
	return &fe.Owner
}

func (fe *FieldExpr) SetIdentifier(ident Node) Accessor {
	fe.Identifier = ident
	return fe
}

func (fe *FieldExpr) SetOwner(owner Accessor) Accessor {
	fe.Owner = owner
	return fe
}

type MemberExpr struct {
	Owner      Accessor
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
	return 0
}

func (me *MemberExpr) SetProperty(prop Node) Accessor {
	me.Property = prop
	return me
}

func (me *MemberExpr) GetProperty() *Node {
	return &me.Property
}

func (me *MemberExpr) GetOwner() *Accessor {
	return &me.Owner
}

func (me *MemberExpr) SetIdentifier(ident Node) Accessor {
	me.Identifier = ident
	return me
}

func (me *MemberExpr) SetOwner(owner Accessor) Accessor {
	me.Owner = owner
	return me
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
	return 0
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

type DirectiveExpr struct {
	Identifier lexer.Token
	Expr       Node
	Token      lexer.Token
	ValueType  PrimitiveValueType
}

func (de *DirectiveExpr) GetType() NodeType {
	return DirectiveExpression
}

func (de *DirectiveExpr) GetToken() lexer.Token {
	return de.Token
}

func (de *DirectiveExpr) GetValueType() PrimitiveValueType {
	return 0
}

type Improper struct {
	Token lexer.Token
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
