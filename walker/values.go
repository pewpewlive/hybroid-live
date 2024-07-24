package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/lexer"
)

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

type UnresolvedVal struct {
	Expr ast.Node
}

func NewUnresolvedVal(expr ast.Node) *UnresolvedVal {
	return &UnresolvedVal{
		Expr: expr,
	}
}

func (uv *UnresolvedVal) GetType() Type {
	return &UnresolvedType{
		Expr: uv.Expr,
	}
}

func (v *UnresolvedVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "UNRESOLVED"}
}

type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsConst bool
	IsLocal bool
	Token   lexer.Token
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
	Path string
	EnvType ast.EnvType
}

func NewPathVal(path string, envType ast.EnvType) *PathVal {
	return &PathVal{
		Path: path,
		EnvType: envType,
	}
} 

func (self *PathVal) GetType() Type {
	return NewPathType(self.EnvType)
}

func (self *PathVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "DEFAULT_PATH"}
}

type AnonStructVal struct {
	Fields map[string]*VariableVal
	Lenient bool
}

func NewAnonStructVal(fields map[string]*VariableVal, lenient bool)  *AnonStructVal {
	return &AnonStructVal{
		Fields: fields,
		Lenient: lenient,
	}
}

func (self *AnonStructVal) GetType() Type {
	return NewAnonStructType(self.Fields, self.Lenient)
}

func (self *AnonStructVal) GetDefault() *ast.LiteralExpr {
	src := generator.StringBuilder{}

	src.WriteString("{")
	length := len(self.Fields) - 1
	index := 0
	for k, v := range self.Fields {
		if index == length {
			src.Append(k, " = ", v.GetDefault().Value)
		} else {
			src.Append(k, " = ", v.GetDefault().Value, ", ")
		}
		index++
	}
	src.WriteString("}")

	return &ast.LiteralExpr{Value: src.String()}
}

func (self *AnonStructVal) AddField(variable *VariableVal) {
	self.Fields[variable.Name] = variable
}

func (self *AnonStructVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, found := self.Fields[name]; found {
		return variable, -1, true
	}

	return nil, -1, false
}

type CustomVal struct {
	Type *CustomType
}

func NewCustomVal(typ *CustomType) *CustomVal {
	return &CustomVal{
		Type: typ,
	}
}

func (self *CustomVal) GetType() Type {
	return self.Type
}

func (self *CustomVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "CUSTOM_VALUE"}
}

type EnumVal struct {
	Type   *EnumType
	Fields []*VariableVal
}

func NewEnumVal(name string, isLocal bool, fields ...string) *EnumVal {
	val := &EnumVal{
		Type: NewEnumType(name),
		Fields: []*VariableVal{},
	}
	if len(fields) == 0 {
		return val
	}
	for i := range fields {
		val.Fields = append(val.Fields, NewEnumFieldVar(fields[i], name, isLocal))
	}
	return val
}

func (self *EnumVal) GetType() Type {
	return self.Type
}

func (self *EnumVal) GetDefault() *ast.LiteralExpr {
	src := generator.StringBuilder{}

	src.Append(self.Type.Name, "[1]")

	return &ast.LiteralExpr{Value: src.String()}
}

func (self *EnumVal) AddField(variable *VariableVal) {
	self.Fields = append(self.Fields, variable)
}

func (self *EnumVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, ind, found := FindFromList(self.Fields, name); found {
		return variable, ind, true
	}

	return nil, -1, false
}

type EnumFieldVal struct {
	Type *EnumType
}

func (self *EnumFieldVal) GetType() Type {
	return self.Type
}

func (self *EnumFieldVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "ENUM_FIELD_VAL"}
}

func NewEnumFieldVar(name string, enumName string, isLocal bool) *VariableVal {
	return &VariableVal{
		Name: name,
		Value: &EnumFieldVal{
			Type: NewEnumType(enumName),
		},
		IsLocal: isLocal,
		IsConst: true,
	}
}


type RawEntityVal struct{}

func (self *RawEntityVal) GetType() Type {
	return &RawEntityType{}
}

func (ev *RawEntityVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "ID"}
}

type EntityVal struct {
	Type          NamedType
	IsLocal       bool
	Fields        []*VariableVal
	Methods       map[string]*VariableVal
	SpawnParams   Types
	DestroyParams Types
}

func NewEntityVal(name string, isLocal bool) *EntityVal {
	return &EntityVal{
		Type:    *NewNamedType(name, ast.Entity),
		IsLocal: isLocal,
		Methods: make(map[string]*VariableVal),
		Fields: make([]*VariableVal, 0),
	}
}

func (ev *EntityVal) GetType() Type {
	return &ev.Type
}

func (ev *EntityVal) GetDefault() *ast.LiteralExpr {
	src := generator.StringBuilder{}

	src.WriteString("{")
	length := len(ev.Fields) - 1
	for i, v := range ev.Fields {
		if i == length {
			src.Append(v.GetDefault().Value)
		} else {
			src.Append(v.GetDefault().Value, ", ")
		}
	}
	src.WriteString("}")

	return &ast.LiteralExpr{Value: src.String()}
}

// Container
func (ev *EntityVal) AddField(variable *VariableVal) {
	ev.Fields = append(ev.Fields, variable)
}

func (ev *EntityVal) AddMethod(variable *VariableVal) {
	ev.Methods[variable.Name] = variable
}

func (ev *EntityVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, ind, found := FindFromList(ev.Fields, name); found {
		return variable, ind, true
	}

	return nil, -1, false
}

func (ev *EntityVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := ev.Methods[name]; found {
		return variable, true
	}

	return nil, false
}

type StructVal struct {
	Type    NamedType
	IsLocal bool
	Fields  []*VariableVal
	Methods map[string]*VariableVal
	Params  Types
}

func (sv *StructVal) GetType() Type {
	return &sv.Type
}

func (sv *StructVal) GetDefault() *ast.LiteralExpr {
	src := generator.StringBuilder{}

	src.WriteString("{")
	length := len(sv.Fields) - 1
	for i, v := range sv.Fields {
		if i == length {
			src.Append(v.GetDefault().Value)
		} else {
			src.Append(v.GetDefault().Value, ", ")
		}
	}
	src.WriteString("}")

	return &ast.LiteralExpr{Value: src.String()}
}

// Container
func (sv *StructVal) AddField(variable *VariableVal) {
	sv.Fields = append(sv.Fields, variable)
}

func (sv *StructVal) AddMethod(variable *VariableVal) {
	sv.Methods[variable.Name] = variable
}

func (sv *StructVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, ind, found := FindFromList(sv.Fields, name); found {
		return variable, ind, true
	}

	return nil, -1, false
}

func (sv *StructVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := sv.Methods[name]; found {
		return variable, true
	}

	return nil, false
}

type MapVal struct {
	MemberType Type
	Members    []Value
}

func (m *MapVal) GetType() Type {
	return NewWrapperType(NewBasicType(ast.Map), m.MemberType)
}

func (m *MapVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "{}"}
}

func GetContentsValueType(values []Value) Type {
	valTypes := []Type{}
	index := 0
	if len(values) == 0 {
		return NAType
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
			return NAType
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
	return &ast.LiteralExpr{Value: "TYPES"}
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
		src.WriteString(types[i].ToString())
		if i != len(types)-1 {
			src.WriteRune('\n')
		}
	}

	return src.String()
}

type FunctionVal struct {
	Params  Types
	Returns Types
}

func NewFunction(params ...Type) *FunctionVal {
	return &FunctionVal{
		Params: params,
	}
}

func (fn FunctionVal) WithReturns(returns ...Type) *FunctionVal {
	fn.Returns = returns
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
	src.WriteString("function(")
	for i := range f.Params {
		src.WriteString(fmt.Sprintf("param%v", i))
		if i != len(f.Params)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(") end")
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

type Invalid struct{}

func (u *Invalid) GetType() Type {
	return InvalidType
}

func (n *Invalid) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "invalid"}
}

type Unknown struct{}

func (u *Unknown) GetType() Type {
	return NewBasicType(ast.Unknown)
}

func (u *Unknown) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "unknown"}
}
