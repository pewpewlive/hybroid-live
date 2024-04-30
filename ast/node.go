package ast

import "hybroid/lexer"

type NodeType int

const (
	VariableDeclarationStatement NodeType = iota + 1
	FunctionDeclarationStatement
	StructureDeclarationStatement
	FieldDeclarationStatement
	MethodDeclarationStatement
	ConstructorStatement

	AssignmentStatement
	RepeatStatement
	TickStatement
	IfStatement
	UseStatement
	AddStatement
	RemoveStatement
	ReturnStatement
	MatchStatement

	Progr

	AnonymousFunctionExpression
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
	SelfExpression

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

	Invalid
)

func (pvt PrimitiveValueType) ToString() string {
	return [...]string{"unknown", "number", "string", "bool", "fixedpoint", "fixed", "radian", "degree", "list", "map", "nil", "func", "entity", "struct", "identifier", "namespace", "undefined"}[pvt]
}

type Node interface {
	GetType() NodeType
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}
