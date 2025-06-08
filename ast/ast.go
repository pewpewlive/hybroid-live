package ast

import "hybroid/tokens"

type NodeType string

const (
	EnvironmentDeclaration    NodeType = "environmentDeclaration"
	VariableDeclaration       NodeType = "variableDeclaration"
	FunctionDeclaration       NodeType = "functionDeclaration"
	ClassDeclaration          NodeType = "classDeclaration"
	ConstructorDeclaration    NodeType = "constructorDeclaration"
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
	CaseStatement       NodeType = "caseStatement"

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
	EntityEvaluationExpression  NodeType = "entityEvaluationExpression"

	EntityAccessExpression NodeType = "entityAccessExpression"

	Identifier NodeType = "identifier"

	NA NodeType = "notAssessed"
)

type PrimitiveValueType string

const (
	Object        PrimitiveValueType = "object"
	Number        PrimitiveValueType = "number"
	Text          PrimitiveValueType = "text"
	Bool          PrimitiveValueType = "bool"
	Fixed         PrimitiveValueType = "fixed"
	List          PrimitiveValueType = "list"
	Map           PrimitiveValueType = "map"
	Func          PrimitiveValueType = "func"
	Entity        PrimitiveValueType = "entity"
	Class         PrimitiveValueType = "class"
	Struct        PrimitiveValueType = "struct"
	Ident         PrimitiveValueType = "ident"
	Enum          PrimitiveValueType = "enum"
	Path          PrimitiveValueType = "path"
	Generic       PrimitiveValueType = "generic"
	Invalid       PrimitiveValueType = "invalid"
	Uninitialized PrimitiveValueType = "uninitialized"
)

type Env string

const (
	MeshEnv    Env = "MeshEnv"
	LevelEnv   Env = "LevelEnv"
	SoundEnv   Env = "SoundEnv"
	GenericEnv Env = "GenericEnv"
	InvalidEnv Env = "InvalidEnv"
)

type MethodCallType int

const (
	ClassMethod MethodCallType = iota
	EntityMethod
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

// used only for map expression
type Property struct {
	Key  Node
	Expr Node
}

type Node interface {
	GetType() NodeType
	GetToken() tokens.Token
}

type NodeCall interface {
	Node
	GetGenerics() []*TypeExpr
	GetCaller() Node
	GetArgs() []Node
}

type Body []Node

func NewBody() Body {
	return make(Body, 0)
}

func (b Body) Size() int {
	return len(b)
}

func (b Body) Node(i int) *Node {
	return &b[i]
}

func (b *Body) Append(node Node) {
	*b = append(*b, node)
}

type MethodInfo struct {
	MethodType MethodCallType
	MethodName string
	TypeName   string
	EnvName    string
}

func NewMethodInfo(methodType MethodCallType, methodName string, typeName string, envName string) MethodInfo {
	return MethodInfo{
		MethodType: methodType,
		MethodName: methodName,
		TypeName:   typeName,
		EnvName:    envName,
	}
}
