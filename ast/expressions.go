package ast

import (
	"hybroid/lexer"
)

type Accessor interface {
	Node 
	SetProperty(prop Node) Accessor 
	GetOwner() *Node
} 

type LiteralExpr struct {
	Value     string
	ValueType PrimitiveValueType
	Token     lexer.Token
}

func (le LiteralExpr) GetType() NodeType {
	return LiteralExpression
}

func (n LiteralExpr) GetToken() lexer.Token {
	return n.Token
}

func (n LiteralExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type UnaryExpr struct {
	Value     Node
	Operator  lexer.Token
	ValueType PrimitiveValueType
}

func (ue UnaryExpr) GetType() NodeType {
	return UnaryExpression
}

func (n UnaryExpr) GetToken() lexer.Token {
	return n.Operator
}

func (n UnaryExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type TypeExpr struct { //syntax: Type<WrappedType>
	WrappedType *TypeExpr
	Name        lexer.Token
	Params      []TypeExpr
	Returns     []TypeExpr
}

func (t TypeExpr) GetType() NodeType {
	return TypeExpression
}

func (t TypeExpr) GetToken() lexer.Token {
	return t.Name
}

func (t TypeExpr) GetValueType() PrimitiveValueType {
	return 0
}

type GroupExpr struct {
	Expr      Node
	ValueType PrimitiveValueType
	Token     lexer.Token
}

func (ge GroupExpr) GetType() NodeType {
	return GroupingExpression
}

func (n GroupExpr) GetToken() lexer.Token {
	return n.Token
}

func (n GroupExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type BinaryExpr struct {
	Left, Right Node
	Operator    lexer.Token
	ValueType   PrimitiveValueType
}

func (be BinaryExpr) GetType() NodeType {
	return BinaryExpression
}

func (n BinaryExpr) GetToken() lexer.Token {
	return n.Operator
}

func (n BinaryExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type CallExpr struct {
	Identifier string
	Caller     Node
	Args       []Node
	Token      lexer.Token
}

func (ce CallExpr) GetType() NodeType {
	return CallExpression
}

func (n CallExpr) GetToken() lexer.Token {
	return n.Token
}

func (n CallExpr) GetValueType() PrimitiveValueType {
	return 0
}

type AnonFnExpr struct {
	Token  lexer.Token
	Return []TypeExpr
	Params []Param
	Body   []Node
}

func (af AnonFnExpr) GetType() NodeType {
	return AnonymousFunctionExpression
}

func (af AnonFnExpr) GetToken() lexer.Token {
	return af.Token
}

func (af AnonFnExpr) GetValueType() PrimitiveValueType {
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

func (se SelfExpr) GetType() NodeType {
	return SelfExpression
}

func (se SelfExpr) GetToken() lexer.Token {
	return se.Token
}

func (se SelfExpr) GetValueType() PrimitiveValueType {
	return 0
}

type MethodCallExpr struct {
	TypeName string 
	Owner  	 Node 
	Call 	   Node
	MethodName string
	Args   	 []Node
	Token  	 lexer.Token
}

func (new MethodCallExpr) GetType() NodeType {
	return MethodCallExpression
}

func (new MethodCallExpr) GetToken() lexer.Token {
	return new.Token
}

func (new MethodCallExpr) GetValueType() PrimitiveValueType {
	return 0
}

type NewExpr struct {
	Type   lexer.Token
	Args []Node
	Token  lexer.Token
}

func (new NewExpr) GetType() NodeType {
	return NewExpession
}

func (new NewExpr) GetToken() lexer.Token {
	return new.Token
}

func (new NewExpr) GetValueType() PrimitiveValueType {
	return 0
}

type FieldExpr struct {
	Owner       Node
	Property    Node
	Identifier  Node
	Index       int
}

func (fi FieldExpr) GetType() NodeType {
	return FieldExpression
}

func (fi FieldExpr) GetToken() lexer.Token {
	return fi.Identifier.GetToken()
}

func (fi FieldExpr) GetValueType() PrimitiveValueType {
	return 0
}

func (fi FieldExpr) SetProperty(prop Node) Accessor {
	fi.Property = prop
	return fi
}

func (fi FieldExpr) GetOwner() *Node {
	return &fi.Owner 
}

type MemberExpr struct {
	Owner      Node
	Property   Node
	Identifier Node
	IsList     bool
}

func (me MemberExpr) GetType() NodeType {
	return MemberExpression
}

func (me MemberExpr) GetToken() lexer.Token {
	return me.Identifier.GetToken()
}

func (me MemberExpr) GetValueType() PrimitiveValueType {
	return 0
}

func (me MemberExpr) SetProperty(prop Node) Accessor {
	me.Property = prop
	return me
}

func (me MemberExpr) GetOwner() *Node {
	return &me.Owner 
}

type Property struct {
	Expr Node
	Type PrimitiveValueType
}

type MapExpr struct {
	Token lexer.Token
	Map   map[lexer.Token]Property
}

func (me MapExpr) GetType() NodeType {
	return MapExpression
}

func (n MapExpr) GetToken() lexer.Token {
	return n.Token
}

func (n MapExpr) GetValueType() PrimitiveValueType {
	return 0
}

type ListExpr struct {
	List      []Node
	ValueType PrimitiveValueType
	Token     lexer.Token
}

func (le ListExpr) GetType() NodeType {
	return ListExpression
}

func (n ListExpr) GetToken() lexer.Token {
	return n.Token
}

func (n ListExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type IdentifierExpr struct {
	Name      lexer.Token
	ValueType PrimitiveValueType
}

func (ie IdentifierExpr) GetType() NodeType {
	return Identifier
}

func (n IdentifierExpr) GetToken() lexer.Token {
	return n.Name
}

func (n IdentifierExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type DirectiveExpr struct {
	Identifier lexer.Token
	Expr       Node
	Token      lexer.Token
	ValueType  PrimitiveValueType
}

func (de DirectiveExpr) GetType() NodeType {
	return DirectiveExpression
}

func (de DirectiveExpr) GetToken() lexer.Token {
	return de.Token
}

func (de DirectiveExpr) GetValueType() PrimitiveValueType {
	return 0
}

type Improper struct {
	Token lexer.Token
}

func (un Improper) GetType() NodeType {
	return NA
}

func (n Improper) GetToken() lexer.Token {
	return n.Token
}

func (n Improper) GetValueType() PrimitiveValueType {
	return Invalid
}
