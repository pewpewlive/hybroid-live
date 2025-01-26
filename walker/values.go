package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/tokens"
)

type ScopeableValue interface {
	Value
	Scopify(parent *Scope, expr *ast.FieldExpr) *Scope
}

type Value interface {
	GetType() Type
	GetDefault() *ast.LiteralExpr
}

type FieldContainer interface {
	Value
	AddField(variable *VariableVal)
	ContainsField(name string) (*VariableVal, int, bool)
}

type MethodContainer interface {
	Value
	AddMethod(variable *VariableVal)
	ContainsMethod(name string) (*VariableVal, bool)
}

type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsConst bool
	IsLocal bool
	IsInit  bool
	Token   tokens.Token
}

func (v *VariableVal) GetType() Type {
	return v.Value.GetType()
}

func (v *VariableVal) GetDefault() *ast.LiteralExpr {
	return v.Value.GetDefault()
}

func FindFromList(list []*VariableVal, name string) (*VariableVal, int, bool) {
	for i, v := range list {
		if v.Name == name {
			return v, i + 1, true
		}
	}
	return nil, -1, false
}

// this will check if it actually has a scopable value inside of it if not then its invalid
//func (v *VariableVal) Scopify()

type PathVal struct {
	Path    string
	EnvType ast.EnvType
}

func NewPathVal(path string, envType ast.EnvType) *PathVal {
	return &PathVal{
		Path:    path,
		EnvType: envType,
	}
}

func (pv *PathVal) GetType() Type {
	return NewPathType(pv.EnvType)
}

func (pv *PathVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type AnonStructVal struct {
	Fields  map[string]Field
	Lenient bool
}

func NewAnonStructVal(fields map[string]Field, lenient bool) *AnonStructVal {
	return &AnonStructVal{
		Fields:  fields,
		Lenient: lenient,
	}
}

func (asv *AnonStructVal) GetType() Type {
	return NewStructTypeWithFields(asv.Fields, asv.Lenient)
}

func (asv *AnonStructVal) GetDefault() *ast.LiteralExpr {
	src := generator.StringBuilder{}

	src.Write("{")
	length := len(asv.Fields) - 1
	index := 0
	for k, v := range asv.Fields {
		if index == length {
			src.Write(k, " = ", v.Var.GetDefault().Value)
		} else {
			src.Write(k, " = ", v.Var.GetDefault().Value, ", ")
		}
		index++
	}
	src.Write("}")

	return &ast.LiteralExpr{Value: src.String()}
}

func (asv *AnonStructVal) AddField(variable *VariableVal) {
	asv.Fields[variable.Name] = NewField(len(asv.Fields), variable)
}

func (asv *AnonStructVal) ContainsField(name string) (*VariableVal, int, bool) {
	if v, found := asv.Fields[name]; found {
		return v.Var, v.Index + 1, true
	}

	return nil, -1, false
}

func (asv *AnonStructVal) Scopify(parent *Scope, expr *ast.FieldExpr) *Scope {
	scope := NewScope(parent, &UntaggedTag{})

	for k, v := range asv.Fields {
		scope.Variables[k] = v.Var
	}
	return scope
}

type CustomVal struct {
	Type *CustomType
}

func NewCustomVal(typ *CustomType) *CustomVal {
	return &CustomVal{
		Type: typ,
	}
}

func (cv *CustomVal) GetType() Type {
	return cv.Type
}

func (cv *CustomVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type EnumVal struct {
	Type   *EnumType
	Fields map[string]*VariableVal
}

func NewEnumVal(envName string, name string, isLocal bool, fields ...string) *EnumVal {
	val := &EnumVal{
		Type:   NewEnumType(envName, name),
		Fields: map[string]*VariableVal{},
	}
	if len(fields) == 0 {
		return val
	}
	for i := range fields {
		val.Fields[fields[i]] = NewEnumFieldVar(fields[i], *val.Type, i)
	}
	return val
}

func (ev *EnumVal) GetType() Type {
	return ev.Type
}

func (ev *EnumVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

func (ev *EnumVal) AddField(variable *VariableVal) {
	if ev == nil {
		return
	}
	enumFieldVal := variable.Value.(*EnumFieldVal)
	enumFieldVal.Index = len(ev.Fields)
	ev.Fields[variable.Name] = variable
}

func (ev *EnumVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := ev.Fields[name]; found {
		return variable, variable.Value.(*EnumFieldVal).Index + 1, true
	}

	return nil, -1, false
}

func (ev *EnumVal) Scopify(parent *Scope, expr *ast.FieldExpr) *Scope {
	scope := NewScope(parent, &UntaggedTag{})

	scope.Variables = ev.Fields

	return scope
}

type EnumFieldVal struct {
	Index int
	Type  *EnumType
}

func (efv *EnumFieldVal) GetType() Type {
	return efv.Type
}

func (efv *EnumFieldVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

func NewEnumFieldVar(name string, enumType EnumType, index int) *VariableVal {
	return &VariableVal{
		Name: name,
		Value: &EnumFieldVal{
			Type: &enumType,
		},
		IsLocal: true,
		IsConst: true,
	}
}

type RawEntityVal struct{}

func (rev *RawEntityVal) GetType() Type {
	return &RawEntityType{}
}

func (rev *RawEntityVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type Field struct {
	Index int
	Var   *VariableVal
}

func NewField(index int, val *VariableVal) Field {
	return Field{
		Index: index,
		Var:   val,
	}
}

type EntityVal struct {
	Type            NamedType
	IsLocal         bool
	Fields          map[string]Field
	Methods         map[string]*VariableVal
	SpawnParams     Types
	DestroyParams   Types
	SpawnGenerics   []*GenericType
	DestroyGenerics []*GenericType
}

func NewEntityVal(envName string, name string, isLocal bool) *EntityVal {
	return &EntityVal{
		Type:    *NewNamedType(envName, name, ast.Entity),
		IsLocal: isLocal,
		Methods: make(map[string]*VariableVal),
		Fields:  make(map[string]Field, 0),
	}
}

func (ev *EntityVal) GetType() Type {
	return &ev.Type
}

func (ev *EntityVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

// Container
func (ev *EntityVal) AddField(variable *VariableVal) {
	ev.Fields[variable.Name] = NewField(len(ev.Fields), variable)
}

func (ev *EntityVal) AddMethod(variable *VariableVal) {
	ev.Methods[variable.Name] = variable
}

func (ev *EntityVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := ev.Fields[name]; found {
		return variable.Var, variable.Index + 1, true
	}

	return nil, -1, false
}

func (ev *EntityVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := ev.Methods[name]; found {
		return variable, true
	}

	return nil, false
}

func (ev *EntityVal) Scopify(parent *Scope, expr *ast.FieldExpr) *Scope {
	scope := NewScope(parent, &UntaggedTag{})

	for _, v := range ev.Fields {
		scope.Variables[v.Var.Name] = v.Var
	}

	for _, v := range ev.Methods {
		scope.Variables[v.Name] = v
	}

	return scope
}

type ClassVal struct {
	Type     NamedType
	IsLocal  bool
	Fields   map[string]Field
	Methods  map[string]*VariableVal
	Params   Types
	Generics []*GenericType
}

func (cv *ClassVal) GetType() Type {
	return &cv.Type
}

func (cv *ClassVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

// Container
func (cv *ClassVal) AddField(variable *VariableVal) {
	cv.Fields[variable.Name] = NewField(len(cv.Fields), variable)
}

func (cv *ClassVal) AddMethod(variable *VariableVal) {
	cv.Methods[variable.Name] = variable
}

func (cv *ClassVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := cv.Fields[name]; found {
		return variable.Var, variable.Index + 1, true
	}

	return nil, -1, false
}

func (cv *ClassVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := cv.Methods[name]; found {
		return variable, true
	}

	return nil, false
}

func (cv *ClassVal) Scopify(parent *Scope, expr *ast.FieldExpr) *Scope {
	scope := NewScope(parent, &UntaggedTag{})

	for _, v := range cv.Fields {
		scope.Variables[v.Var.Name] = v.Var
	}

	for _, v := range cv.Methods {
		scope.Variables[v.Name] = v
	}

	return scope
}

type MapVal struct {
	MemberType Type
	Members    []Value
}

func (mv *MapVal) GetType() Type {
	return NewWrapperType(NewBasicType(ast.Map), mv.MemberType)
}

func (mv *MapVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "{}"}
}

func GetContentsValueType(values []Value) Type {
	valTypes := []Type{}
	index := 0
	if len(values) == 0 {
		return ObjectTyp
	}
	for _, v := range values {
		if index == 0 {
			valTypes = append(valTypes, v.GetType())
			index++
			continue
		}
		valTypes = append(valTypes, v.GetType())
		prev, curr := index-1, len(valTypes)-1
		if !TypeEquals(values[prev].GetType(), values[curr].GetType()) {
			return InvalidType
		}
		index++
	}

	return valTypes[0]
}

type ListVal struct {
	ValueType Type
	Values    []Value
}

func (l *ListVal) GetType() Type {
	return NewWrapperType(NewBasicType(ast.List), l.ValueType)
}

func (l *ListVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "{}"}
}

type NumberVal struct{}

func (n *NumberVal) GetType() Type {
	return NewBasicType(ast.Number)
}

func (n *NumberVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "0"}
}

type FixedVal struct {
	SpecificType ast.PrimitiveValueType
}

func (f *FixedVal) GetType() Type {
	return NewFixedPointType(f.SpecificType)
}

func (f *FixedVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "0fx"}
}

func (f *FixedVal) GetSpecificType() ast.PrimitiveValueType {
	return f.SpecificType
}

var EmptyReturn = Types{}

type Types []Type

func (ts Types) GetType() Type {
	if len(ts) == 0 {
		return (&Invalid{}).GetType()
	}
	return (ts)[0]
}

func (ts Types) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

func (rt Types) Eq(otherRT Types) bool {
	typesSame := true
	if len(rt) == len(otherRT) {
		for i, v := range rt {
			if !TypeEquals(v, otherRT[i]) {
				typesSame = false
				break
			}
		}
	} else {
		typesSame = false
	}
	return typesSame
}

func TypesToString(types []Type) string {
	src := generator.StringBuilder{}

	for i := range types {
		src.Write(types[i].ToString())
		if i != len(types)-1 {
			src.WriteRune('\n')
		}
	}

	return src.String()
}

type FunctionVal struct {
	Generics []*GenericType
	Params   Types
	Returns  Types
}

func NewFunction(params ...Type) *FunctionVal {
	return &FunctionVal{
		Params:  params,
		Returns: EmptyReturn,
	}
}

func (fn FunctionVal) WithReturns(returns ...Type) *FunctionVal {
	fn.Returns = returns
	return &fn
}

func (fn FunctionVal) WithGenerics(generics ...*GenericType) *FunctionVal {
	fn.Generics = generics
	return &fn
}

func (f *FunctionVal) GetType() Type {
	return NewFunctionType(f.Params, f.Returns)
}

func (f *FunctionVal) GetReturns() Types {
	return f.Returns
}

func (f *FunctionVal) GetDefault() *ast.LiteralExpr {
	src := generator.StringBuilder{}
	src.Write("function(")
	for i := range f.Params {
		src.Write(fmt.Sprintf("param%v", i))
		if i != len(f.Params)-1 {
			src.Write(", ")
		}
	}
	src.Write(") end")
	return &ast.LiteralExpr{Value: src.String()}
}

type BoolVal struct{}

func (b *BoolVal) GetType() Type {
	return NewBasicType(ast.Bool)
}

func (b *BoolVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "false"}
}

type StringVal struct{}

func (s *StringVal) GetType() Type {
	return NewBasicType(ast.String)
}

func (s *StringVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "\"\""}
}

type GenericVal struct {
	Type *GenericType
}

func (gv *GenericVal) GetType() Type {
	return gv.Type
}

func (gv *GenericVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type Invalid struct{}

func (u *Invalid) GetType() Type {
	return InvalidType
}

func (n *Invalid) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type Unknown struct{}

func (u *Unknown) GetType() Type {
	return ObjectTyp
}

func (u *Unknown) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}
