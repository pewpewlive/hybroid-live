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
	//TypeDeclarationStatement      NodeType = "typeDeclarationStatement"

	EntityDeclarationStatement        NodeType = "entityDeclarationStatement"
	EntityFunctionDeclarationStatemet NodeType = "entityFunctionDeclarationStatement"
	DestroyStmt                       NodeType = "spawnStatement"

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
	SpawnExpression             NodeType = "spawnExpression"

	PewpewExpression      NodeType = "pewpewExpression"
	PewpewCallExpression  NodeType = "pewpewCallExpression"
	FmathExpression       NodeType = "fmathExpression"
	BuiltinCallExpression NodeType = "builtinCallExpession"
	StandardExpression    NodeType = "standardExpression"

	Identifier NodeType = "identifier"

	NA NodeType = "notAssessed"
)

type PrimitiveValueType string

const (
	Unknown    PrimitiveValueType = "unknown"
	Number     PrimitiveValueType = "number"
	String     PrimitiveValueType = "string"
	Bool       PrimitiveValueType = "bool"
	FixedPoint PrimitiveValueType = "fixedPoint"
	Fixed      PrimitiveValueType = "fixed"
	Radian     PrimitiveValueType = "radian"
	Degree     PrimitiveValueType = "degree"
	List       PrimitiveValueType = "list"
	Map        PrimitiveValueType = "map"
	Func       PrimitiveValueType = "func"
	Entity     PrimitiveValueType = "entity"
	Struct     PrimitiveValueType = "struct"
	AnonStruct PrimitiveValueType = "anonStruct"
	Ident      PrimitiveValueType = "ident"
	Enum       PrimitiveValueType = "enum"
	Path       PrimitiveValueType = "path"
	Unresolved PrimitiveValueType = "unresolved"
	Invalid    PrimitiveValueType = "invalid"
)

type EnvType string

const (
	MeshEnv EnvType = "MeshEnv"
	Level EnvType = "LevelEnv"
	SoundEnv EnvType = "SoundEnv"
	InvalidEnv EnvType = "InvalidEnv"
)

type StandardLibrary int

const (
	MathLib StandardLibrary = iota
	StringLib
	TableLib
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
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}

type Accessor interface {
	Node
	GetOwner() *Accessor
	GetProperty() *Node
	SetOwner(owner Accessor)
	SetProperty(prop Node)
	SetIdentifier(ident Node)
	Copy() Accessor
}

var Libraries = map[string]StandardLibrary{
	"Math":   MathLib,
	"String": StringLib,
	"Table":  TableLib,
}
