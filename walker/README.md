# Walker

The walker walks through all the nodes in the AST (Abstract Syntax Tree), verifying their legitimacy and/or changing them.

## `values.go`

This section covers all the structs and interfaces used for abstracting values.

### Interfaces

#### **_Value_**

```go
type Value interface {
  GetType() Type
  GetDefault() ast.LiteralExpr
}
```

It's used to abstract any kind of value, including numbers, booleans, nil, strings, structs, maps, lists, etc.

##### **Methods:**

1. `Type GetType()` - returns the type of the value in the form of a Type value
2. `GetDefault() ast.LiteralExpr` - returns the default value in the form of a literal expression node

##### **Implementations:**

###### VariableVal
```go
type VariableVal struct {
  Name    string
  Value   Value
  IsUsed  bool
  IsConst bool
  Node    ast.Node
}
```
###### StructVal
```go
type StructVal struct {
  Type         *NamedType
  Params       []Type // of the constructor
  Fields       []*VariableVal
  FieldIndexes map[string]int
  Methods      map[string]*VariableVal
  IsUsed       bool
}
```

`StructVal` is a value that contains all the data of a struct type. It's used in the [Environment](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#environment) struct as `Structs`.

Some nodes such as field expressions are associated with a struct type, which is associated with a struct. This struct is pretty much the `StructVal`. When declaring a struct type, you also declare its body, with all its methods and fields. You pretty much declare a `StructVal`, which gets added into the `Structs` of the [Environment](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#environment).

###### EnvironmentVal
```go
type EnvironmentVal struct {
	Path         string
	Name         string
	Fields       map[string]VariableVal
	Methods      map[string]VariableVal
	FieldIndexes map[string]int
}
```

###### **_Environment_**

```go
type EnvironmentVal struct {
	Name         string
  Path         string
	Ctx          Context
	Scope        Scope
	StructTypes  map[string]*StructVal
}
```

The Environment is like the Global scope of the file. The `Name` of it is provided through the EnvStmt (CITATION NEEDED).

`Path` is the hybroid path.

**Constructor:**
`NewEnvironment() EnvironmentVal`

###### MapVal
```go
type MapVal struct {
  MemberType Type
  Members    map[string]Value
}
```

**Extra methods:**

1. `GetContentsValueType() -> Type` - Checks the contents of the `MapVal` and, if all the values have the same type, returns the `Type` that they all share. If they don't have the same value type it returns `Invalid`.

###### ListVal
```go
type ListVal struct {
  ValueType Type
  Values    []Value
}
```

**Extra methods:**

1. a `GetContentsValueType() -> Type` - Same with `MapVal`'s method

###### Types

```go
type Types []Type
```

It is important to note here that if it contains more than 1 type, `GetType` returns invalid, otherwise it returns the only type it has.

`GetDefault` just returns a literal expression with a value of "TYPES".

###### NumberVal
```go
type NumberVal struct{}
```
###### DirectiveVal
```go
type DirectiveVal struct{}
```
###### FixedVal
```go
type FixedVal struct {
  SpecificType ast.PrimitiveValueType
}
```
###### FunctionVal
```go
type FunctionVal struct {
  params    Returns
  returnVal Returns
}
```
###### BoolVal
```go
type BoolVal struct{}
```
###### StringVal
```go
type StringVal struct{}
```
###### NilVal
```go
type NilVal struct{}
```
###### Invalid
```go
type Invalid struct{}
```
###### Unknown
```go
type Unknown struct{}
```

#### **_Container_**

```go
type Container interface {
	Value
	GetFields() map[string]VariableVal
	GetMethods() map[string]VariableVal
	AddField(variable VariableVal)
	AddMethod(variable VariableVal)
	Contains(name string) (Value, int, bool)
}
```

Used to abstract any kind of value that contains fields and methods (struct, entity). Extends [Value](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#value).

##### **Methods:**

1. `GetFields() -> map[string]VariableVal` - returns a map of its fields.
2. `GetMethods() -> map[string]VariableVal` - returns a map of its methods.
3. `AddField(variable VariableVal)` - adds a fields to its fields
4. `AddMethod(variable VariableVal)` - adds a method to its methods 
5. `Contains(name string) -> (Value, int, bool)` - checks if any of its fields or methods contain _name_ and gives the value of it along with its index in the list. The boolean determines the success.

##### **Implementations:**

###### [StructTypeVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#structtypeval)
###### [StructVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#structval)
###### [EnvironmentVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#environmentval)
###### [EntityVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#entityval)

### Global variables

#### **_EmptyReturn_**
```go
var EmptyReturn = Types{}
```

## `scope.go`

This section covers all the interfaces and structs used to make walking the nodes more organized and easier.

### Interfaces

#### **_ScopeTag_**

```go
type ScopeTag interface {
  GetType() ScopeTagType
}
```

`ScopeTag` is like [ScopeAttribute](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#scopeattribute), but it holds some special information with it (depending on the interface implementation).

_When creating a new scope, the tag of the parent scope does not get carried onto the new one._

##### **Implementations:**

###### UntaggedTag
```go
type UntaggedTag struct{}
```
###### StructTag
```go
type StructTag struct {
  StructType *StructTypeVal
}
```
###### EntityTag
```go
//to be used
type EntityTag struct {
  //EntityType *StructTypeVal
}
```
###### FuncTag
```go
type FuncTag struct {
  Returns    []bool
  ReturnType Returns
}
```
###### MatchExprTag
```go
type MatchExprTag struct {
  mpt         MultiPathTag
  ArmsYielded int
  YieldValues *Returns
}
```
###### MultiPathTag
```go
type MultiPathTag struct {
  ReturnAmount   []bool
  YieldAmount    []bool
  ContinueAmount []bool
  BreakAmount    []bool
}
```

The values here express how many times the `Scope` (i.e. the body) has returned, yielded, continued and breaked. These values are used by many nodes (usually statements like `IfStmt`) and then evaluated.

#### **_ExitableTag_**

```go
type ExitableTag interface {
	ScopeTag
	SetExit(state bool, _type ExitType)
	GetIfExits(_type ExitType) bool
}
```

ExitableTag extends [ScopeTag](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#scopetag).

##### **Methods:**

1. `SetExit(state bool, _type ExitType)` - appends to `state` to a list of booleans. The list to be used for the append depends on the ExitType. For example the `[]bool` named `Returns` would be used in this case if `_type` is `Return`.
2. `GetIfExits(_type ExitType) -> bool` - evaluates the list of booleans corresponding to the `_type` and returns true if it does exit.

##### **Implementations:**

###### [FuncTag](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#functag)
###### [MatchTag](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#matchtag)
###### [MultiPathTag](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#multipathtag)

### Structures

Here are the fundamental structs that are extremely important for the walking process:

#### **_Context_**

```go
type Context struct {
	Node  ast.Node
	Value Value
	Ret   Types
}
```

`Context` provides info about the previous node and only the previous node. Thus, it is short-sighted but quite useful for some cases (see fieldExpr (CITATION NEED) and memberExpr (CITATION NEEDED) walker methods).

It is possible that this struct might get removed later, as it might not be needed, but that's up to debate.

#### **_Scope_**

```go
type Scope struct {
  Environment *Environment
  Parent *Scope

  Tag        ScopeTag
  Attributes ScopeAttributes

  Variables       map[string]*VariableVal
  VariableIndexes map[string]int
}
```

`Scope` is essentially a body that contains variables and has a tag and attributes. A scope has a parent which it stems from.

**Constructor:**
`NewScope(parent *Scope, tag ScopeTag, attrs ...ScopeAttribute) -> Scope` - Returns a new scope with its parent being the _parent_ parameter, its tag the _tag_ param. The attributes of the _parent_ get passed onto the new scope plus any extra attributes defined (_attrs_).

##### **Methods:**

1. `Is(types ...ScopeAttribute) bool` - Checks whether the scope contains the given scope attributes.

### Types

#### ScopeAttributes
```go
type ScopeAttributes []ScopeAttribute
```

##### **Constructor:**
`NewScopeAttributes(types ...ScopeAttribute) -> ScopeAttributes`

##### **Methods:**

1. `Add(attr ScopeAttribute)` - Adds the _attr_ into the _ScopeAttributes_ if it doesn't exist.

### Enums

#### **_ScopeAttribute_**

```go
type ScopeAttribute int

const (
  ReturnAllowing ScopeAttribute = iota + 1
  YieldAllowing
  SelfAllowing
  BreakAllowing
  ContinueAllowing
)
```

`ScopeAttribute` allows scopes to be ascribed with a specific property. Sometimes, when the nodes are being walked, the walker needs to know about the scope in more detail, especially when we want to prohibit the coder from writing some specific nodes (e.g. Only in Struct and Entity scopes do we allow Self to be used, so naturally you want that scope to carry that priviledge, hence SelfAllowing).

_It is important to note that, when creating a new scope, the attributes of the parent scope are carried onto the new scope._

#### **_ExitType_**

```go
type ExitType int

const (
  Yield ExitType = iota
  Return
  Continue
  Break
  All
)
```

Expresses how the body is exiting (yielding, returning, continuing or breaking?). 

`All` is used for unreachable code detection. When we check if a scope exits, _All_ allows us to check if it exits through _All_ its paths that it could have, regardless of *how* it exits. Whether it does exit on all paths or not it still gets passed onto the parent scope tag.

#### **_ScopeTagType_**

```go
type ScopeTagType int

const (
  Untagged ScopeTagType = iota
  Struct
  Entity
  Func
  MultiPath
  MatchExpr
  Loop
)
```

`ScopeTagType` is the identity of the [ScopeTag](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#scopetag).

### `statements.go`

