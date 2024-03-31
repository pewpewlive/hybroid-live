package ast

import "hybroid/lexer"

type NodeType int

const (
	VariableDeclarationStatement NodeType = iota + 1
	FunctionDeclarationStatement

	AssignmentStatement
	RepeatStatement
	TickStatement
	IfStatement
	UseStatement
	AddStatement
	RemoveStatement
	ReturnStatement

	Progr

	DirectiveExpression
	LiteralExpression
	UnaryExpression
	BinaryExpression
	GroupingExpression
	ListExpression
	MapExpression
	CallExpression
	MemberExpression
	ParentExpression
	TypeExpression

	Identifier

	NA
)

type PrimitiveValueType int

const (
	Number PrimitiveValueType = iota + 1
	String
	Bool
	FixedPoint
	Fixed
	Radian
	Degree
	List
	Map
	Nil
	Func
	Entity
	Struct
	Ident
	Namespace

	Undefined
)

func (pvt PrimitiveValueType) ToString() string {
	return [...]string{"unknown", "number", "string", "bool", "fixedpoint", "fixed", "radian", "degree", "list", "map", "nil", "func", "entity", "struct", "identifier", "undefined"}[pvt]
}

type Node interface {
	GetType() NodeType
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}
