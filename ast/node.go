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

	EnvironmentStatement
	AssignmentStatement
	RepeatStatement
	ForStatement
	TickStatement
	IfStatement
	UseStatement
	AddStatement
	RemoveStatement
	BreakStatement
	ContinueStatement
	ReturnStatement
	YieldStatement
	MatchStatement

	EnvironmentExpression
	AnonymousFunctionExpression
	AnonymousStructExpression
	DirectiveExpression
	LiteralExpression
	UnaryExpression
	BinaryExpression
	GroupingExpression
	ListExpression
	MapExpression
	CallExpression
	MethodCallExpression
	MatchExpression
	FieldExpression
	MemberExpression
	ParentExpression
	TypeExpression
	SelfExpression
	NewExpession

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
	Func
	Entity
	Struct
	AnonStruct
	Ident
	Environment

	Invalid
)

var stringifiedPTV = [...]string{"unknown", "number", "text", "bool", "fixed", "fixed", "fixed", "fixed", "list", "map", "func", "entity", "struct", "anonymous struct", "identifier", "namespace", "invalid"}

func (pvt PrimitiveValueType) ToString() string {
	return stringifiedPTV[pvt]
}

type Node interface {
	GetType() NodeType
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}
