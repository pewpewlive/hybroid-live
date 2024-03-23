package ast

import "hybroid/lexer"

type NodeType int

const (
	VariableDeclarationStatement NodeType = iota + 1
	FunctionDeclarationStatement

	AssignmentStatement
	RepeatStatement
	IfStatement

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

	Undefined
)

type Node interface {
	GetType() NodeType
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}
