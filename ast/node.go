package ast

import "hybroid/lexer"

type NodeType string

const (
	VariableDeclarationStatement  NodeType = "variableDeclarationStatement"
	FunctionDeclarationStatement  NodeType = "functionDeclarationStatement"
	StructureDeclarationStatement NodeType = "structureDeclarationStatement"
	FieldDeclarationStatement     NodeType = "fieldDeclarationStatement"
	MethodDeclarationStatement    NodeType = "methodDeclarationStatement"
	EnumDeclarationStatement      NodeType = "enumDeclarationStatement"
	MacroDeclarationStatement     NodeType = "macroDeclarationStatement"
	

	ConstructorStatement NodeType = "constructorStatement"
	EnvironmentStatement NodeType = "environmentStatement"
	AssignmentStatement  NodeType = "assignmentStatement"
	RepeatStatement      NodeType = "repeatStatement"
	WhileStatement       NodeType = "whileStatement"
	ForStatement         NodeType = "forStatement"
	TickStatement        NodeType = "tickStatement"
	IfStatement          NodeType = "ifStatement"
	UseStatement         NodeType = "useStatement"
	AddStatement         NodeType = "addStatement"
	RemoveStatement      NodeType = "removeStatement"
	BreakStatement       NodeType = "breakStatement"
	ContinueStatement    NodeType = "continueStatement"
	ReturnStatement      NodeType = "returnStatement"
	YieldStatement       NodeType = "yieldStatement"
	MatchStatement       NodeType = "matchStatement"

	EnvironmentPathExpression   NodeType = "environmentPathExpression"
	EnvironmentAccessExpression NodeType = "environmentAccessExpression"
	EnvironmentTypeExpression   NodeType = "environmentTypeExpression"
	AnonymousFunctionExpression NodeType = "anonymousFunctionExpression"
	AnonymousStructExpression   NodeType = "anonymousStructExpression"
	LiteralExpression           NodeType = "literalExpression"
	UnaryExpression             NodeType = "unaryExpression"
	BinaryExpression            NodeType = "binaryExpression"
	GroupingExpression          NodeType = "groupingExpression"
	ListExpression              NodeType = "listExpression"
	MapExpression               NodeType = "mapExpression"
	CallExpression              NodeType = "callExpression"
	MethodCallExpression        NodeType = "methodCallExpression"
	MacroCallExpression         NodeType = "macroCallExpression"
	MatchExpression             NodeType = "matchExpression"
	FieldExpression             NodeType = "fieldExpression"
	MemberExpression            NodeType = "memberExpression"
	ParentExpression            NodeType = "parentExpression"
	TypeExpression              NodeType = "typeExpression"
	SelfExpression              NodeType = "selfExpression"
	NewExpession                NodeType = "newExpession"

	Identifier NodeType = "identifier"

	NA NodeType = "notAssessed"
)

type PrimitiveValueType string

const (
	Unknown     PrimitiveValueType = "unknown"
	Number      PrimitiveValueType = "number"
	String      PrimitiveValueType = "string"
	Bool        PrimitiveValueType = "bool"
	FixedPoint  PrimitiveValueType = "fixedPoint"
	Fixed       PrimitiveValueType = "fixed"
	Radian      PrimitiveValueType = "radian"
	Degree      PrimitiveValueType = "degree"
	List        PrimitiveValueType = "list"
	Map         PrimitiveValueType = "map"
	Func        PrimitiveValueType = "func"
	Entity      PrimitiveValueType = "entity"
	Struct      PrimitiveValueType = "struct"
	AnonStruct  PrimitiveValueType = "anonStruct"
	Ident       PrimitiveValueType = "ident"
	Environment PrimitiveValueType = "environment"
	Enum        PrimitiveValueType = "enum"
	Unresolved  PrimitiveValueType = "unresolved"
	Invalid     PrimitiveValueType = "invalid"
)

type Node interface {
	GetType() NodeType
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}
