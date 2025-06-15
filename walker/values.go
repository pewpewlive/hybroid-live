package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"

	"github.com/mitchellh/copystructure"
)

type Value interface {
	GetType() Type
	GetDefault() *ast.LiteralExpr
}

type ScopeableValue interface {
	Value
	Scopify(parent *Scope) *Scope
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

type FullContainer interface {
	FieldContainer
	MethodContainer
}

type ConstVal struct {
	Node ast.Node
	Val  Value
}

func (c *ConstVal) GetType() Type {
	return c.Val.GetType()
}

func (c *ConstVal) GetDefault() *ast.LiteralExpr {
	return c.Val.GetDefault()
}

type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsPub   bool
	IsInit  bool
	IsConst bool
	Token   tokens.Token
}

func (w *Walker) SetVarToUsed(v *VariableVal) {
	v.IsUsed = true
	if fn, ok := v.Value.(*FunctionVal); ok && fn.ProcType == Method {
		if fn.MethodName == "spawn" && fn.MethodType == ast.EntityMethod {
			val := w.walkers[fn.EnvName].environment.Entities[fn.TypeName]
			val.Type.IsUsed = true
			w.walkers[fn.EnvName].environment.Entities[fn.TypeName] = val
		} else if fn.MethodName == "new" && fn.MethodType == ast.ClassMethod {
			val := w.walkers[fn.EnvName].environment.Classes[fn.TypeName]
			val.Type.IsUsed = true
			w.walkers[fn.EnvName].environment.Classes[fn.TypeName] = val
		}
	} else if num, ok := v.Value.(*NumberVal); ok {
		num.Value = ""
	} else if boolean, ok := v.Value.(*BoolVal); ok {
		boolean.Value = ""
	}
}

// args go as follows: value Value, isPub bool
func NewVariable(token tokens.Token, args ...any) *VariableVal {
	var value Value
	isPub := false
	if args != nil {
		if args[0] != nil {
			value = args[0].(Value)
		}

		if len(args) == 2 {
			isPub = args[1].(bool)
		}
	}

	return &VariableVal{
		Name:   token.Lexeme,
		Token:  token,
		Value:  value,
		IsPub:  isPub,
		IsInit: value != nil,
	}
}

// Makes the variable const
func (v *VariableVal) Const() *VariableVal {
	v.IsConst = true
	return v
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

type PathVal struct {
	Path    string
	Env     ast.Env
	EnvName string
}

func NewPathVal(path string, envType ast.Env, envName string) *PathVal {
	return &PathVal{
		Path:    path,
		Env:     envType,
		EnvName: envName,
	}
}

func (pv *PathVal) GetType() Type {
	return NewPathType(pv.Env)
}

func (pv *PathVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type StructVal struct {
	Fields  map[string]StructField
	Lenient bool
}

func NewStructVal(fields map[string]StructField) *StructVal {
	return &StructVal{
		Fields: fields,
	}
}

func (sv *StructVal) GetType() Type {
	return &StructType{
		Fields: sv.Fields,
	}
}

func (sv *StructVal) GetDefault() *ast.LiteralExpr {
	src := core.StringBuilder{}

	src.Write("{")
	length := len(sv.Fields) - 1
	index := 0
	for k, v := range sv.Fields {
		val := v.Var.GetDefault().Value
		if val == "nil" {
			return &ast.LiteralExpr{Value: "nil"}
		}
		if index == length {
			src.Write(k, " = ", val)
		} else {
			src.Write(k, " = ", val, ", ")
		}
		index++
	}
	src.Write("}")

	return &ast.LiteralExpr{Value: src.String()}
}

func (sv *StructVal) AddField(variable *VariableVal) {
	sv.Fields[variable.Name] = StructField{Var: variable, Lenient: false}
}

func (sv *StructVal) ContainsField(name string) (*VariableVal, int, bool) {
	if v, found := sv.Fields[name]; found {
		return v.Var, -1, true
	}

	return nil, -1, false
}

func (sv *StructVal) Scopify(parent *Scope) *Scope {
	scope := NewScope(parent, &UntaggedTag{})

	for k, v := range sv.Fields {
		scope.Variables[k] = v.Var
	}
	return scope
}

type EnumVal struct {
	Token  tokens.Token
	Type   *EnumType
	Fields map[string]*VariableVal
	IsPub  bool
}

func NewEnumVal(envName string, name string, isPub bool, fields ...string) *EnumVal {
	val := &EnumVal{
		Type:   NewEnumType(envName, name),
		Fields: map[string]*VariableVal{},
		IsPub:  isPub,
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
	enumFieldVal := variable.Value.(*EnumFieldVal)
	enumFieldVal.Index = len(ev.Fields) + 1
	ev.Fields[variable.Name] = variable
}

func (ev *EnumVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := ev.Fields[name]; found {
		return variable, variable.Value.(*EnumFieldVal).Index, true
	}

	return nil, -1, false
}

func (ev *EnumVal) Scopify(parent *Scope) *Scope {
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
	}
}

type RawEntityVal struct{}

func (rev *RawEntityVal) GetType() Type {
	return &RawEntityType{}
}

func (rev *RawEntityVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

type StructField struct {
	Lenient bool
	Var     *VariableVal
}

func NewStructField(name string, val Value, lenient ...bool) StructField {
	isLenient := false
	if lenient != nil && len(lenient) == 1 {
		isLenient = lenient[0]
	}

	return StructField{
		Var:     &VariableVal{Name: name, Value: val},
		Lenient: isLenient,
	}
}

type Field struct {
	Index int
	Var   *VariableVal
}

type EntityVal struct {
	Token   tokens.Token
	Type    NamedType
	IsLocal bool
	Fields  map[string]Field
	Methods map[string]*VariableVal

	Spawn   *FunctionVal
	Destroy *FunctionVal
}

func CopyEntityVal(ref *EntityVal) EntityVal {
	val, err := copystructure.Copy(*ref)
	if err != nil {
		panic(err)
	}
	switch newVal := val.(type) {
	case EntityVal:
		return newVal
	default:
		panic(fmt.Sprintf("Attempt to copy entityVal, got: %T", val))
	}
}

func NewEntityVal(envName string, node *ast.EntityDecl) *EntityVal {
	name := node.Name.Lexeme
	return &EntityVal{
		Token:   node.Name,
		Type:    *NewNamedType(envName, name, ast.Entity),
		IsLocal: !node.IsPub,
		Methods: make(map[string]*VariableVal),
		Fields:  make(map[string]Field, 0),
		Destroy: NewMethod(ast.NewMethodInfo(ast.EntityMethod, "destroy", name, envName)),
		Spawn:   NewMethod(ast.NewMethodInfo(ast.EntityMethod, "spawn", name, envName)),
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
	ev.Fields[variable.Name] = Field{
		Var:   variable,
		Index: len(ev.Fields) + 1,
	}
}

func (ev *EntityVal) AddMethod(variable *VariableVal) {
	ev.Methods[variable.Name] = variable
}

func (ev *EntityVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := ev.Fields[name]; found {
		return variable.Var, variable.Index, true
	}

	return nil, -1, false
}

func (ev *EntityVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := ev.Methods[name]; found {
		return variable, true
	}

	return nil, false
}

func (ev *EntityVal) Scopify(parent *Scope) *Scope {
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
	Token       tokens.Token
	Type        NamedType
	IsLocal     bool
	Fields      map[string]Field
	Methods     map[string]*VariableVal
	GenericArgs []Type

	New *FunctionVal
}

func CopyClassVal(ref *ClassVal) ClassVal {
	val, err := copystructure.Copy(*ref)
	if err != nil {
		panic(err)
	}
	switch newVal := val.(type) {
	case ClassVal:
		return newVal
	default:
		panic(fmt.Sprintf("Attempt to copy classVal, got: %T", val))
	}
}

func (cv *ClassVal) GetType() Type {
	return &cv.Type
}

func (cv *ClassVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

// Container
func (cv *ClassVal) AddField(variable *VariableVal) {
	cv.Fields[variable.Name] = Field{
		Var:   variable,
		Index: len(cv.Fields) + 1,
	}
}

func (cv *ClassVal) AddMethod(variable *VariableVal) {
	cv.Methods[variable.Name] = variable
}

func (cv *ClassVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := cv.Fields[name]; found {
		return variable.Var, variable.Index, true
	}

	return nil, -1, false
}

func (cv *ClassVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := cv.Methods[name]; found {
		return variable, true
	}

	return nil, false
}

func (cv *ClassVal) Scopify(parent *Scope) *Scope {
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
}

func (mv *MapVal) GetType() Type {
	return NewWrapperType(NewBasicType(ast.Map), mv.MemberType)
}

func (mv *MapVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "{}"}
}

type ListVal struct {
	ValueType Type
}

func (l *ListVal) GetType() Type {
	return NewWrapperType(NewBasicType(ast.List), l.ValueType)
}

func (l *ListVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "{}"}
}

type NumberVal struct {
	Value string
}

func (n *NumberVal) GetType() Type {
	return NewBasicType(ast.Number)
}

func (n *NumberVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "0"}
}

func NewNumberVal(value string) *NumberVal {
	return &NumberVal{
		Value: value,
	}
}

type FixedVal struct {
}

func (f *FixedVal) GetType() Type {
	return NewFixedPointType()
}

func (f *FixedVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "0fx"}
}

type Values []Value

func (ts Values) GetType() Type {
	if len(ts) == 0 || len(ts) > 1 {
		return (&Unknown{}).GetType()
	}
	return (ts)[0].GetType()
}

func (ts Values) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}

func (rt Values) Eq(otherRT Values) bool {
	typesSame := true
	if len(rt) == len(otherRT) {
		for i, v := range rt {
			if !TypeEquals(v.GetType(), otherRT[i].GetType()) {
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
	src := core.StringBuilder{}

	for i := range types {
		src.Write(types[i].String())
		if i != len(types)-1 {
			src.WriteRune('\n')
		}
	}

	return src.String()
}

var EmptyReturn = []Type{}

type FunctionVal struct {
	Generics       []*GenericType
	Params         []Type
	Returns        []Type
	ProcType       ProcedureType
	ast.MethodInfo // check if ProcType == Method before accessing this
}

func NewFunction(params ...Type) *FunctionVal {
	return &FunctionVal{
		ProcType: Function,
		Params:   params,
		Returns:  EmptyReturn,
	}
}

func NewMethod(mi ast.MethodInfo, params ...Type) *FunctionVal {
	return &FunctionVal{
		ProcType:   Method,
		Params:     params,
		Returns:    EmptyReturn,
		MethodInfo: mi,
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
	return NewFunctionType(f.Params, f.Returns, f.ProcType)
}

func (f *FunctionVal) GetReturns() []Type {
	return f.Returns
}

func (f *FunctionVal) GetDefault() *ast.LiteralExpr {
	src := core.StringBuilder{}
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

type BoolVal struct {
	Value string
}

func (b *BoolVal) GetType() Type {
	return NewBasicType(ast.Bool)
}

func (b *BoolVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "false"}
}

func NewBoolVal(value string) *BoolVal {
	return &BoolVal{
		Value: value,
	}
}

type StringVal struct{}

func (s *StringVal) GetType() Type {
	return NewBasicType(ast.Text)
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
	return UnknownTyp
}

func (u *Unknown) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "nil"}
}
