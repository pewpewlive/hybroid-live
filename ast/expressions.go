package ast

import "hybroid/lexer"

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
	Caller     Node //identifier
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

type MemberExpr struct {
	Identifier Node
	Property   Node
	Bracketed  bool
	Token      lexer.Token
}

func (me MemberExpr) GetType() NodeType {
	return MemberExpression
}

func (n MemberExpr) GetToken() lexer.Token {
	return n.Token
}

func (n MemberExpr) GetValueType() PrimitiveValueType {
	return n.Identifier.GetValueType()
}

type Property struct {
	Expr Node
	Type PrimitiveValueType
}

type MapExpr struct {
	ValueType PrimitiveValueType
	Token     lexer.Token
	Map       map[string]Property
}

func (me MapExpr) GetType() NodeType {
	return MapExpression
}

func (n MapExpr) GetToken() lexer.Token {
	return n.Token
}

func (n MapExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
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
	Name      string
	Token     lexer.Token
	ValueType PrimitiveValueType
}

func (ie IdentifierExpr) GetType() NodeType {
	return Identifier
}

func (n IdentifierExpr) GetToken() lexer.Token {
	return n.Token
}

func (n IdentifierExpr) GetValueType() PrimitiveValueType {
	return n.ValueType
}

type DirectiveExpr struct {
	Identifier string
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
	return Undefined
}

type Unknown struct {
	Token lexer.Token
}

func (un Unknown) GetType() NodeType {
	return NA
}

func (n Unknown) GetToken() lexer.Token {
	return n.Token
}

func (n Unknown) GetValueType() PrimitiveValueType {
	return Undefined
}
