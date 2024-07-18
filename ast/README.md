# AST

Here lies all of the datatypes used by the parser to create the Abstract Syntax Tree. These datatypes are also used by the walker, notably [PrimitiveValueType](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#primitivevaluetype).

## `error.go`

This file contains the datatypes for errors and such.

### Interfaces

#### **_Alert_**

Alert is used to abstract errors and warnings.

```go
type Alert interface {
	GetToken() lexer.Token
	GetHeader() string
	GetMessage() string
}
```

##### **Methods:**

1. `GetToken() -> lexer.Token` - Returns the [lexer.Token]() (CITATION NEEDED) that the alert holds.
2. `GetHeader() -> string` - Returns "[red]Error" if it's an error and "[yellow]Warning" if it's a warning. This is used to prefix the message of the alert and it also colors the string appropriately.
3. `GetMessage() -> string` - Returns the message of the alert.

##### **Implementations:**

###### Error

This is used in the parser and walker.

```go
type Error struct {
	Token   lexer.Token
	Message string
}
```

###### Warning

This is used in the walker only.

```go
type Warning struct {
	Token   lexer.Token
	Message string
}
```

## `node.go`

### Enums

#### **_NodeType_**

```go
type NodeType string

const (
	VariableDeclarationStatement  NodeType = "variableDeclarationStatement"
	FunctionDeclarationStatement  NodeType = "functionDeclarationStatement"
	StructureDeclarationStatement NodeType = "structureDeclarationStatement"
	FieldDeclarationStatement     NodeType = "fieldDeclarationStatement"
	MethodDeclarationStatement    NodeType = "methodDeclarationStatement"
	EnumDeclarationStatement      NodeType = "enumDeclarationStatement"
	MacroDeclarationStatement     NodeType = "macroDeclarationStatement"

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
	
	PewpewExpression            NodeType = "pewpewExpression"
	PewpewCallExpression        NodeType = "pewpewCallExpression"
	FmathExpression         NodeType = "fmathExpression"
	BuiltinCallExpression       NodeType = "builtinCallExpession"
	StandardExpression      NodeType = "standardExpression"

	Identifier NodeType = "identifier"

	NA NodeType = "notAssessed"
)
```

#### **_PrimitiveValueType_**

```go
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
	Enum        PrimitiveValueType = "enum"
	Unresolved  PrimitiveValueType = "unresolved"
	Invalid     PrimitiveValueType = "invalid"
)
```

Some notable enum values here are `Unknown` and `Unresolved`. `Unknown` is used when we can't determine the type of a value, but don't want there to be errors thrown when being evaluated with other values. `Unresolved` is used whenever there is a value that is accessed from an [`Environment`](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#environment).

#### **_StandardLibrary_**

This enum refers to the lua standard libraries that [ppl-utils](https://github.com/pewpewlive/ppl-utils) allows to be used.

```go
type StandardLibrary int

const (
	MathLib StandardLibrary = iota
	StringLib
	TableLib 
)
```

#### **_EnvType_**

The type of the environment.

```go
type EnvType int

const (
	Mesh EnvType = iota
	Level
	Sound
	InvalidEnv
)
```

#### **_SelfExprType_**

Whenever using `self`, does it refer to an entity or a struct?

```go
type SelfExprType int

const (
	SelfStruct SelfExprType = iota
	SelfEntity
)
```

#### **_MacroType_**

Does the macro expand to an expression, or to a list of statements?

```go
type MacroType int

const (
	ExpressionExpansion MacroType = iota
	ProgramExpansion
)
```

#### **_EntityFunctionType_**

Used in [EntityFunctionDeclarationStmt](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#entityfunctiondeclarationstmt).

```go
type EntityFunctionType string

const (
	WeaponCollision EntityFunctionType = "weaponCollision"
	WallCollision   EntityFunctionType = "wallCollision"
	PlayerCollision EntityFunctionType = "playerCollision"
	Update          EntityFunctionType = "update"
	Destroy         EntityFunctionType = "destroy"
	Spawn           EntityFunctionType = "spawn"
)
```

### Types

#### **_Paths_**

Used in [EnvironmentStmt](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#environmentstmt).

```go
type Paths []string
```

### Interfaces

#### **_Node_**

Node is the abstracted essence of a collection of tokens.

```go
type Node interface {
	GetType() NodeType
	GetToken() lexer.Token
	GetValueType() PrimitiveValueType
}
```

##### **Methods:**

1. `GetType() -> NodeType` - Returns the type of the node.
2. `GetToken() -> lexer.Token` - Returns the main token of the node.
3. `GetValueType() -> PrimitiveValueType` - Returns the PVT of the node. When it's a statement it returns `Invalid`.

##### **Implementations:**

For the node statement implementations, please refer to [statements.go](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#expressionsgo).

For the node expression implementations, please refer to [expressions.go](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#statementsgo).

#### **_Accessor_**

`Accessor` is an extension of node. It essentially describes any node which has property which you can access.

```go
type Accessor interface {
	Node
	GetOwner() *Accessor
	GetProperty() *Node
	SetOwner(owner Accessor)
	SetProperty(prop Node)
	SetIdentifier(ident Node)
	Copy() Accessor
}
```

##### **Methods:**
1. `GetOwner() -> *Accessor` - Returns the owner of the node. Can be `nil`.
2. `GetProperty() -> *Node` - Returns the property of the node. Can be `nil`.
3. `SetOwner(owner Accessor)` - Set's the node's owner to the given `owner` parameter.
4. `SetProperty(prop Node)` - Set's the node's property to the given `prop` parameter.
5. `SetIdentifier(ident Node)` - Set's the node's identifier to the given `ident` parameter.
6. `Copy() -> Accessor` - Performs a shallow copy on the `*Accessor`.

##### **Implementations:**

###### [MemberExpr](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#memberexpr)
###### [FieldExpr](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#fieldExpr)

### Variables

#### **_Libraries_**

```go
var Libraries = map[string]StandardLibrary{
	"Math": MathLib,
	"String": StringLib,
	"Table": TableLib,
}
```

## `expressions.go`

This file contains the implementations of [Node](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#node), specifically the expressions.

### `Node` Implementations

#### **_EnvTypeExpr_**


```go
type EnvTypeExpr struct {
	Type  EnvType
	Token lexer.Token
}
```

In code:
```
env Path as Level
env Path as Mesh
env Path as Sound
			^^^^^
```

#### **_EnvPathExpr_**

```go
type EnvPathExpr struct {
	SubPaths []string
}
```

In code:
```
env Path as Level
	^^^^
env Path::Subpath as Level
	^^^^^^^^^^^^^
```

#### **_EnvAccessExpr_**

```go
type EnvAccessExpr struct {
	PathExpr *EnvPathExpr
	Accessed Node
}
```

In code:
```
let a = Path::expression
```

#### **_MacroCallExpr_**

```go
type MacroCallExpr struct {
	Caller *CallExpr
}

```

In code:
```
@MacroCall()
```

#### **_LiteralExpr_**

```go
type LiteralExpr struct {
	Value     string
	ValueType PrimitiveValueType
	Token     lexer.Token
}
```

In code:
```
let a = "text"
		^^^^^^
let b = 1
		^
let c, d = true, false
		   ^^^^  ^^^^^
```

#### **_UnaryExpr_**

```go
type UnaryExpr struct {
	Value     Node
	Operator  lexer.Token
	ValueType PrimitiveValueType
}
```

In code:
```
let a = -1
		^^
let b = !true
```

#### **_TypeExpr_**

```go
type TypeExpr struct {
	WrappedType *TypeExpr
	Name        Node
	Params      []*TypeExpr
	Returns     []*TypeExpr
	Fields      []Param
}
```

In code:
```
fn function(Type param) {}
			^^^^
let Type a
	^^^^

let list<Type> b
    ^^^^^^^^^^
let map<fixed> c
    ^^^^^^^^^^
```

#### **_GroupExpr_**

```go
type GroupExpr struct {
	Expr      Node
	ValueType PrimitiveValueType
	Token     lexer.Token
}
```

In code:
```
let a = (1f-5f)*9f
		^^^^^^^
```

#### **_BinaryExpr_**

```go
type BinaryExpr struct {
	Left, Right Node
	Operator    lexer.Token
	ValueType   PrimitiveValueType
}
```

In code:
```
let a = 1-5
		^^^
```

#### **_CallExpr_**

```go
type CallExpr struct {
	Name lexer.Token
	Caller     Node
	Args       []Node
}
```

In code:
```
Call()
```

#### **_PewpewExpr_**

```go
type PewpewExpr struct {
	Node Node
}
```

In code:
```
pewpew::expression
```

#### **_FmathExpr_**

```go
type FmathExpr struct {
	Node Node
}
```

In code:
```
fmath::expression
```

#### **_BuiltinCallExpr_**

```go
type BuiltinCallExpr struct {
	Name lexer.Token
	Args []Node
}
```

#### **_StandardExpr_**

```go
type StandardExpr struct {
	Library StandardLibrary
	Node Node
}
```

In code:
```
table::expression
math::expression
string::expression
```

#### **_AnonFnExpr_**

```go
type AnonFnExpr struct {
	Token  lexer.Token
	Return []*TypeExpr
	Params []Param
	Body   []Node
}
```

In code:
```
let a = fn(fixed param) {}
		^^^^^^^^^^^^^^^^^^
```

#### **_AnonStructExpr_**

```go
type AnonStructExpr struct {
	Token  lexer.Token
	Fields []*FieldDeclarationStmt
}
```

In code:
```
let a = struct { x = 0fx }
		^^^^^^^^^^^^^^^^^^
```

#### **_MatchExpr_**

```go
type MatchExpr struct {
	MatchStmt    MatchStmt
	ReturnAmount int
}
```

In code:
```
let a = match expression {
	expression => {
		yield 1
	}
	expression => 2
	_ => 3
}
```

#### **_SelfExpr_**

```go
type SelfExpr struct {
	Token lexer.Token
	Type  SelfExprType
}
```

In code:
```
self.expression

let a = self
        ^^^^
```

#### **_MethodCallExpr_**

```go
type MethodCallExpr struct {
	TypeName   string
	OwnerType  SelfExprType
	Owner      Node
	Call       Node
	MethodName string
	Args       []Node
	Token      lexer.Token
}
```

In code:
```
self.call()
```

#### **_NewExpr_**

```go
type NewExpr struct {
	Type  *TypeExpr
	Args  []Node
	Token lexer.Token
}
```

In code:
```
let instance = new Struct()
			   ^^^^^^^^^^^^
```

#### **_SpawnExpr_**

```go
type SpawnExpr struct {
	Type  *TypeExpr
	Args  []Node
	Token lexer.Token
}
```

In code:
```
let instance = spawn Entity()
			   ^^^^^^^^^^^^^^
```

#### **_FieldExpr_**

```go
type FieldExpr struct {
	Owner      Accessor
	Property   Node
	Identifier Node
	Index      int
}
```

In code:
```
owner.property
```

#### **_MemberExpr_**

```go
type MemberExpr struct {
	Owner      Accessor
	Property   Node
	Identifier Node
	IsList     bool
}
```

In code:
```
map["property"]
list[index]
```

#### **_MapExpr_**

Refer to [here](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#property) to see the definition of `Property`.

```go
type MapExpr struct {
	Token lexer.Token
	Map   map[lexer.Token]Property
}
```

In code:
```
let map<text> a = {
	property1 = "text",
	property2 = "text2"
}
```

#### **_ListExpr_**

```go
type ListExpr struct {
	List      []Node
	ValueType PrimitiveValueType
	Token     lexer.Token
}
```

In code:
```
let a = [1,2,3]
```

#### **_IdentifierExpr_**

```go
type IdentifierExpr struct {
	Name      lexer.Token
	ValueType PrimitiveValueType
}
```

In code:
```
let a = identifier
		^^^^^^^^^^
```

#### **_Improper_**

```go
type Improper struct {
	Token lexer.Token
}
```

##### **Constructor:**

`NewImproper(token lexer.Token) -> *Improper` - Creates an Improper node with the given `token` parameter.

### Structs

#### **_Property_**

Used in [MapExpr](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#mapexpr).

```go
type Property struct {
	Expr Node
	Type PrimitiveValueType
}
```

## `statements.go`

This file contains the implementations of [Node](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#node), specifically the statements.

### `Node` Implementations

#### **_EnvironmentStmt_**

```go
type EnvironmentStmt struct {
	EnvType      *EnvTypeExpr
	Env          *EnvPathExpr
	Requirements Paths
}
```

##### **Methods:**

1. `AddRequirement(path string) -> bool` - Takes a `path` parameter and adds that to the `Requirements`. Returns `true` if the `Requirements` already contained the path, otherwise `false`.

#### **_AssignmentStmt_**

```go
type AssignmentStmt struct {
	Values      []Node
	Identifiers []Node
	Token       lexer.Token
}
```

In code:
```
variable = expression
```

#### **_MacroDeclarationStmt_**

```go
type MacroDeclarationStmt struct {
	Name      lexer.Token
	Params    []lexer.Token
	MacroType MacroType
	Tokens    []lexer.Token
}
```

In code:
```
macro name() => {
	...
}

macro name() => ...
```

The first one describes a macro that expands to a list of statements, while the second one expands to an expression.

#### **_VariableDeclarationStmt_**

```go
type VariableDeclarationStmt struct {
	Identifiers []lexer.Token
	Types       []*TypeExpr
	Values      []Node
	IsLocal     bool
	Token       lexer.Token
}
```

In code:
```
let a = 1
let fixed b = 2
```

#### **_StructDeclarationStmt_**

```go
type StructDeclarationStmt struct {
	Token       lexer.Token
	Name        lexer.Token
	Fields      []FieldDeclarationStmt
	Constructor *ConstructorStmt
	Methods     []MethodDeclarationStmt
	IsLocal     bool
}
```

In code:
```
struct Name {
	new() {

	}

	fn method() {

	}
}
```


#### **_EntityDeclarationStmt_**

```go
type EntityDeclarationStmt struct {
	Token     lexer.Token
	Name      lexer.Token
	Fields    []FieldDeclarationStmt
	Spawner   *EntityFunctionDeclarationStmt
	Destroyer *EntityFunctionDeclarationStmt
	Callbacks []*EntityFunctionDeclarationStmt
	Methods   []MethodDeclarationStmt
	IsLocal   bool
}
```

In code:
```
entity Name {
	spawn() {

	}

	destroy() {

	}

	fn method() {
		
	}
}
```

#### **_EntityFunctionDeclarationStmt_**

```go
type EntityFunctionDeclarationStmt struct {
	Type   EntityFunctionType
	Params []Param
	Body   []Node
	Token  lexer.Token
}
```

In code:
```
spawn() {

}

destroy() {
	
}

WeaponCollision() {

}

PlayerCollision() {

}

WallCollision() {

}
```

#### **_EnumDeclarationStmt_**

```go
type EnumDeclarationStmt struct {
	Name    lexer.Token
	Fields  []lexer.Token
	IsLocal bool
}
```

In code:
```
enum Name {
	Field1,
	Field2,
	Field3
}
```

#### **_ConstructorStmt_**

```go
type ConstructorStmt struct {
	Token  lexer.Token
	Body   []Node
	Return []*TypeExpr
	Params []Param
}
```

In code:
```
new() {

}
```

#### **_FieldDeclarationStmt_**

```go
type FieldDeclarationStmt struct {
	Identifiers []lexer.Token
	Types       []*TypeExpr
	Values      []Node
	Token       lexer.Token
}
```

In code:
```
a = 1
fixed b = 2f
```

#### **_FunctionDeclarationStmt_**

Please refer to [here](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#param) for the definition of `Param`.

```go
type FunctionDeclarationStmt struct {
	Name    lexer.Token
	Return  []*TypeExpr
	Params  []Param
	IsLocal bool
	Body    []Node
}
```

In code:
```
fn name(text param1) bool {
}
```

#### **_MethodDeclarationStmt_**

```go
type MethodDeclarationStmt struct {
	Owner   lexer.Token
	Name    lexer.Token
	Return  []*TypeExpr
	Params  []Param
	IsLocal bool
	Body    []Node
}
```

In code:
```
fn name(text param1) bool {
}
```

#### **_MatchStmt_**

Please refer to [here](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#casestmt) for the definition of `CaseStmt`.

```go
type MatchStmt struct {
	ExprToMatch Node
	Cases       []CaseStmt
	HasDefault  bool
}
```

In code:
```
match expression {
	expression => {
		return 1
	}
	expression => 2
	_ => 3
}
```

#### **_IfStmt_**

```go
type IfStmt struct {
	BoolExpr Node
	Body     []Node
	Elseifs  []*IfStmt
	Else     *IfStmt
	Token    lexer.Token
}
```

In code:
```
if condtion {

}else if condtion {

}else {

}
```


#### **_RepeatStmt_**

```go
type RepeatStmt struct {
	Iterator Node
	Skip     Node
	Start    Node
	Variable *IdentifierExpr
	Body     []Node
	Token    lexer.Token
}
```

In code:
```
repeat 10 {

}

repeat to 10 from 5 {

}

repeat to 10 from 6 by 2 {

}
```

#### **_WhileStmt_**

```go
type WhileStmt struct {
	Condtion Node
	Body     []Node
	Token    lexer.Token
}
```

In code:
```
while condtion {

}
```

#### **_ForStmt_**

```go
type ForStmt struct {
	Iterator         Node
	KeyValuePair     [2]*IdentifierExpr
	OrderedIteration bool
	Body             []Node
	Token            lexer.Token
}
```

In code:
```
for index, value in list {

}

for key, value in map {

}
```

#### **_TickStmt_**

```go
type TickStmt struct {
	Variable IdentifierExpr
	Body     []Node
	Token    lexer.Token
}
```

In code:
```
tick {

}

tick with time {

}
```

#### **_ReturnStmt_**

```go
type ReturnStmt struct {
	Args  []Node
	Token lexer.Token
}
```

In code:
```
return expression
```

#### **_YieldStmt_**

```go
type YieldStmt struct {
	Args  []Node
	Token lexer.Token
}
```

In code:
```
yield expression
```

#### **_BreakStmt_**

```go
type BreakStmt struct {
	Token lexer.Token
}
```

In code:
```
break
```

#### **_ContinueStmt_**

```go
type ContinueStmt struct {
	Token lexer.Token
}
```

In code:
```
continue
```

#### **_AddStmt_**

```go
type AddStmt struct {
	Value      Node
	Identifier string
	Token      lexer.Token
}
```

In code:
```
add expression to list
```

#### **_RemoveStmt_**

```go
type RemoveStmt struct {
	Value      Node
	Identifier string
	Token      lexer.Token
}
```

In code:
```
remove expression from list
```

#### **_ContinueStmt_**

```go
type UseStmt struct {
	Path *EnvPathExpr
}
```

In code:
```
use Path
use Path::Subpath
```

### Structs

#### **_Param_**

Used in [FunctionDeclarationStmt](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#functiondeclarationstmt) and [MethodDeclarationStmt](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#methoddeclarationstmt).

```go
type Param struct {
	Type *TypeExpr
	Name lexer.Token
}
```

#### **_CaseStmt_**

Used in [MatchStmt](https://github.com/pewpewlive/hybroid/blob/master/ast/README.md#matchstmt).

```go
type CaseStmt struct {
	Expression Node
	Body       []Node
}
```