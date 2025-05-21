package ast

import "hybroid/tokens"

type NodeType string

const (
	EnvironmentDeclaration    NodeType = "environmentDeclaration"
	VariableDeclaration       NodeType = "variableDeclaration"
	FunctionDeclaration       NodeType = "functionDeclaration"
	ClassDeclaration          NodeType = "classDeclaration"
	ConstructorDeclaration    NodeType = "constructorDeclaration"
	FieldDeclaration          NodeType = "fieldDeclaration"
	MethodDeclaration         NodeType = "methodDeclaration"
	EnumDeclaration           NodeType = "enumDeclaration"
	MacroDeclaration          NodeType = "macroDeclaration"
	AliasDeclaration          NodeType = "aliasDeclaration"
	EntityDeclaration         NodeType = "entityDeclaration"
	EntityFunctionDeclaration NodeType = "entityFunctionDeclaration"

	DestroyStatement    NodeType = "destroyStatement"
	AssignmentStatement NodeType = "assignmentStatement"
	RepeatStatement     NodeType = "repeatStatement"
	WhileStatement      NodeType = "whileStatement"
	ForStatement        NodeType = "forStatement"
	TickStatement       NodeType = "tickStatement"
	IfStatement         NodeType = "ifStatement"
	UseStatement        NodeType = "useStatement"
	AddStatement        NodeType = "addStatement"
	RemoveStatement     NodeType = "removeStatement"
	BreakStatement      NodeType = "breakStatement"
	ContinueStatement   NodeType = "continueStatement"
	ReturnStatement     NodeType = "returnStatement"
	YieldStatement      NodeType = "yieldStatement"
	MatchStatement      NodeType = "matchStatement"

	EnvironmentPathExpression   NodeType = "environmentPathExpression"
	EnvironmentAccessExpression NodeType = "environmentAccessExpression"
	EnvironmentTypeExpression   NodeType = "environmentTypeExpression"
	FunctionExpression          NodeType = "functionExpression"
	StructExpression            NodeType = "structExpression"
	LiteralExpression           NodeType = "literalExpression"
	UnaryExpression             NodeType = "unaryExpression"
	BinaryExpression            NodeType = "binaryExpression"
	GroupExpression             NodeType = "groupExpression"
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
	TupleExpression             NodeType = "tupleExpression"
	SelfExpression              NodeType = "selfExpression"
	NewExpession                NodeType = "newExpession"
	SpawnExpression             NodeType = "spawnExpression"
	EntityExpression            NodeType = "entityExpression"

	PewpewExpression   NodeType = "pewpewExpression"
	FmathExpression    NodeType = "fmathExpression"
	BuiltinExpression  NodeType = "builtinExpression"
	StandardExpression NodeType = "standardExpression"

	Identifier NodeType = "identifier"

	NA NodeType = "notAssessed"
)

type PrimitiveValueType string

const (
	Object        PrimitiveValueType = "object"
	Number        PrimitiveValueType = "number"
	String        PrimitiveValueType = "string"
	Bool          PrimitiveValueType = "bool"
	FixedPoint    PrimitiveValueType = "fixedPoint"
	Fixed         PrimitiveValueType = "fixed"
	Radian        PrimitiveValueType = "radian"
	Degree        PrimitiveValueType = "degree"
	List          PrimitiveValueType = "list"
	Map           PrimitiveValueType = "map"
	Func          PrimitiveValueType = "func"
	Entity        PrimitiveValueType = "entity"
	Struct        PrimitiveValueType = "struct"
	AnonStruct    PrimitiveValueType = "anonStruct"
	Ident         PrimitiveValueType = "ident"
	Enum          PrimitiveValueType = "enum"
	Path          PrimitiveValueType = "path"
	Generic       PrimitiveValueType = "generic"
	Tuple         PrimitiveValueType = "tuple"
	Invalid       PrimitiveValueType = "invalid"
	Uninitialized PrimitiveValueType = "uninitialized"
)

type EnvType string

const (
	MeshEnv    EnvType = "MeshEnv"
	LevelEnv   EnvType = "LevelEnv"
	SoundEnv   EnvType = "SoundEnv"
	InvalidEnv EnvType = "InvalidEnv"
)

type SelfExprType int

const (
	SelfStruct SelfExprType = iota
	SelfEntity
)

type MacroType int

const (
	ExpressionExpansion MacroType = iota
	ProgramExpansion
)

type EntityFunctionType string

const (
	WeaponCollision EntityFunctionType = "weaponCollision"
	WallCollision   EntityFunctionType = "wallCollision"
	PlayerCollision EntityFunctionType = "playerCollision"
	Update          EntityFunctionType = "update"
	Destroy         EntityFunctionType = "destroy"
	Spawn           EntityFunctionType = "spawn"
)

type Paths []string

type Node interface {
	GetType() NodeType
	GetToken() tokens.Token
	GetValueType() PrimitiveValueType
}
