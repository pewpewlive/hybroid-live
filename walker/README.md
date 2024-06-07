# Walker

The walker walks through all the nodes in the AST (Abstract Syntax Tree), verifies their legitimacy and/or changes them.

## `values.go`

This section covers all the structs and interfaces used for abstracting values.

### Interfaces

#### **_Value:_**

```go
type Value interface {
  GetType() TypeVal
  GetDefault() ast.LiteralExpr
}
```

It's used to abstract any kind of value, including numbers, booleans, nil, strings, structs, maps, lists, etc.

##### **Methods:**

1. `TypeVal GetType()` - returns the type of the value in the form of a TypeVal value
2. `ast.LiteralExpr GetDefault()` - returns the default value in the form of a literal expression node

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
###### StructTypeVal
```go
type StructTypeVal struct {
  Name         lexer.Token
  Params       []TypeVal
  Fields       []VariableVal
  FieldIndexes map[string]int
  Methods      map[string]VariableVal
  IsUsed       bool
}
```

`StructTypeVal` is a TypeVal specifically for structs. It's used in the [Namespace](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#namespace) struct as `StructTypes`.

Some nodes such as field expressions are associated with a struct type. This struct type is pretty much the `StructTypeVal`. When declaring a struct, you also declare its body, with all its methods and fields. You pretty much declare a `StructTypeVal`, which gets added into the `StructTypes` of the [Namespace](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#namespace).

###### StructVal
```go
type StructVal struct {
  Type *StructTypeVal
}
```

`StructVal` is a struct value and it contains the type of the struct (i.e `StructTypeVal`).
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
###### MapVal
```go
type MapVal struct {
  MemberType TypeVal
  Members    map[string]VariableVal
}
```

**Extra methods:**

1. `GetContentsValueType() -> TypeVal` - checks the contents of the `MapVal` and, if all the values have the same type, returns the `TypeVal` that they all share. If they don't have the same value type it returns `Invalid`.

###### ListVal
```go
type ListVal struct {
  ValueType TypeVal
  Values    []Value
}
```

**Extra methods:**

1. a `GetContentsValueType() -> TypeVal` - same with `MapVal`'s method

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

#### **_Container:_**

```go
type Container interface {
  GetFields() map[string]VariableVal
  GetMethods() map[string]VariableVal
  Contains(name string) (Value, int, bool)
}
```

Used to abstract any kind of value that contains fields and methods (struct, entity, namespace).

##### **Methods:**

1. `GetFields() -> map[string]VariableVal` - returns a map of its fields.
2. `GetMethods() -> map[string]VariableVal` - returns a map of its methods.
3. `Contains(name string) -> (Value, int, bool)` - checks if any of its fields or methods contain _name_ and gives the value of it along with its index in the list. The boolean determines the success.

##### **Implementations:**

Only `Value`s implement `Container`, specifically:

###### [StructTypeVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#structtypeval)
###### [StructVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#structval)
###### [EnvironmentVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#environmentval)
###### [EntityVal](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#entityval)

### Types

#### Types
```go
type Types []Type
```

## `scope.go`

This section covers all the interfaces and structs used to make walking the nodes more organized and easier.

### Interfaces

#### **_ScopeTag:_**

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

### Structures

Here are the fundamental structs that are extremely important for the walking process:

#### **_Environment:_**

```go
type Environment struct {
	Name         string
	Ctx          Context
	Scope        Scope
	foreignTypes map[string]Value
	StructTypes  map[string]*StructTypeVal
}
```

##### **Constructor:**
`NewEnvironment() Environment`

#### **_Scope:_**

```go
type Scope struct {
  Environment *Environment
  Parent *Scope

  Tag        ScopeTag
  Attributes ScopeAttributes

  Variables       map[string]VariableVal
  VariableIndexes map[string]int
}
```

`Scope` is essentially a body that contains variables and has a tag and attributes. A scope has a parent which it stems from.

**Constructor:**
`NewScope(parent *Scope, tag ScopeTag) -> Scope` - returns a new scope with its parent being the _parent_ parameter and its tag the _tag_ param.

##### **Methods:**

1. `Is(types ...ScopeAttribute) bool` - checks whether the scope contains the given scope attributes.

### Types

#### ScopeAttributes
```go
type ScopeAttributes []ScopeAttribute
```

##### **Constructor:**
`NewScopeAttributes(types ...ScopeAttribute) -> ScopeAttributes`

### Enums

#### **_ScopeAttribute:_**

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

#### **_ExitType:_**

```go
type ExitType int

const (
  Yield ExitType = iota
  Return
  Continue
  Break
)
```

Expresses how the body is exiting (yielding, returning, continuing or breaking?).

#### **_ScopeTagType:_**

```go
type ScopeTagType int

const (
  Untagged ScopeTagType = iota
  Struct
  Entity
  Func
  MultiPath
  MatchExpr
)
```

`ScopeTagType` is the identity of the [ScopeTag](https://github.com/pewpewlive/hybroid/blob/master/walker/README.md#scopetag).

### `statements.go`
