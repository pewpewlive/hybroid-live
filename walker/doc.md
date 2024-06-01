# WALKER

The walker walks through all the nodes, verifies their legitimacy and/or changes them.

## Values.go

This section covers all the structs and interfaces used for abstracting values.

### Value

It's used to abstract any kind of value, including numbers, booleans, nil, strings, structs, maps, lists, etc.

```go
type Value interface {
    GetType() TypeVal
    GetDefault() ast.LiteralExpr
}
```

**Methods:**

1. `TypeVal GetType()` - returns the type of the value in the form of a TypeVal value
2. `ast.LiteralExpr GetType()` - returns the default value in the form of a literal expression node


**Implementations:**

```go
type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsConst bool
	Node    ast.Node
}
```

```go
type TypeVal struct {
	WrappedType *TypeVal
	Name        string
	Type        ast.PrimitiveValueType
	Params      *[]TypeVal
	Returns     ReturnType
}
```

Extra methods:
1. `TypeVal.Eq(otherT TypeVal) -> bool` - returns true if the given TypeVal is the same with self.
2. `TypeVal.ToString() -> string` - returns a stringified version of the TypeVal.

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

`StructTypeVal` is a TypeVal specifically for structs. It's used in the `Global` struct as `StructTypes`. 

Some nodes such as field expressions are associated with a struct type. This struct type is pretty much the `StructTypeVal`. When declaring a struct, you also declare its body, with all its methods and fields. You pretty much declare a `StructTypeVal`, which gets added into the `StructTypes` of the Global.

```go
type StructVal struct {
	Type *StructTypeVal
}
```

`StructVal` is a struct value and it contains the type of the struct (i.e `StructTypeVal`).

```go
type NamespaceVal struct {
	Name         string
	Fields       map[string]VariableVal
	Methods      map[string]VariableVal
	FieldIndexes map[string]int
}
```

```go
type MapMemberVal struct {
	Var   VariableVal
	Owner MapVal
}
```

```go
type MapVal struct {
	MemberType TypeVal
	Members    map[string]MapMemberVal
}
```

Extra methods:
1. `GetContentsValueType() -> TypeVal` - checks the contents of the `MapVal` and, if all the values have the same type, returns the `TypeVal` that they all share. If they don't have the same value type it returns `Invalid`.

```go
type ListVal struct {
	ValueType TypeVal
	Values    []Value
}
```

Extra methods:
1. a `GetContentsValueType() -> TypeVal` - checks the contents of the `MapVal` and, if all the values have the same type, returns the `TypeVal` that they all share. If they don't have the same value type it returns `Invalid`.

```go
type NumberVal struct{}
```

```go
type DirectiveVal struct{}
```

```go
type FixedVal struct {
	SpecificType ast.PrimitiveValueType
}
```

```go
type ReturnType []TypeVal
```

```go
type FunctionVal struct {
	params    []TypeVal
	returnVal ReturnType
}
```

```go
type CallVal struct {
	types ReturnType
}
```

```go
type BoolVal struct{}
```

```go
type StringVal struct{}
```

```go
type NilVal struct{}
```

```go
type Invalid struct{}
```

```go
type Unknown struct{}
```

##
### Container

Used to abstract any kind of value that contains fields and methods (struct, entity, namespace).

```go
type Container interface {
	GetFields() map[string]VariableVal
	GetMethods() map[string]VariableVal
	Contains(name string) (Value, int, bool)
}
```
##
**Methods:**

1. `GetFields() -> map[string]VariableVal` - returns a map of its fields.
2. `GetMethods() -> map[string]VariableVal` - returns a map of its methods.
3. `Contains(name string) -> (Value, int, bool)` - checks if any of its fields or methods contain *name* and gives the value of it along with its index in the list. The boolean determines the success. 

**Implementations:**

Only `Value`s implement `Container`, specifically:
1. `StructTypeVal`
2. `StructVal`
3. `NamespaceVal`
4. `EntityVal` (doesn't exist yet)

## Scope.go

This section covers all the interfaces and structs used to make walking the nodes more organized and easier.

Here are the fundamental structs that are extremely important for the walking process:

```go
type Global struct {
	Ctx          Context
	Scope        Scope
	foreignTypes map[string]Value
	StructTypes  map[string]*StructTypeVal
}
```

`Global` is es

```go
type Scope struct {
	Global *Global
	Parent *Scope

	Tag        ScopeTag
	Attributes ScopeAttributes

	Variables       map[string]VariableVal
	VariableIndexes map[string]int
}
```