package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
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
	Expr *ast.EnvAccessExpr
}

func NewUnresolvedVal(expr *ast.EnvAccessExpr) *UnresolvedVal {
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

type AnonStructVal struct {
	Fields map[string]*VariableVal
}

func (self *AnonStructVal) GetType() Type {
	return NewAnonStructType(self.Fields)
}

func (self *AnonStructVal) GetDefault() *ast.LiteralExpr {
	src := lua.StringBuilder{}

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

type EnumVal struct {
	Type   *EnumType
	Fields []*VariableVal
}

func (self *EnumVal) GetType() Type {
	return self.Type
}

func (self *EnumVal) GetDefault() *ast.LiteralExpr {
	src := lua.StringBuilder{}

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

type StructVal struct {
	Type    NamedType
	IsLocal bool
	Fields  []*VariableVal
	Methods map[string]*VariableVal
	Params  Types
}

func (self *StructVal) GetType() Type {
	return &self.Type
}

func (self *StructVal) GetDefault() *ast.LiteralExpr {
	src := lua.StringBuilder{}

	src.WriteString("{")
	length := len(self.Fields) - 1
	for i, v := range self.Fields {
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
func (self *StructVal) AddField(variable *VariableVal) {
	self.Fields = append(self.Fields, variable)
}

func (self *StructVal) AddMethod(variable *VariableVal) {
	self.Methods[variable.Name] = variable
}

func (self *StructVal) ContainsField(name string) (*VariableVal, int, bool) {
	if variable, ind, found := FindFromList(self.Fields, name); found {
		return variable, ind, true
	}

	return nil, -1, false
}

func (self *StructVal) ContainsMethod(name string) (*VariableVal, bool) {
	if variable, found := self.Methods[name]; found {
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
		if !(parser.IsFx(valTypes[prev].PVT()) && parser.IsFx(valTypes[curr].PVT())) && valTypes[prev].PVT() != valTypes[curr].PVT() {
			return NAType
		}
		index++
	}
	if parser.IsFx(valTypes[0].PVT()) {
		return NewBasicType(ast.FixedPoint)
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
	return NewBasicType(ast.FixedPoint)
}

func (f *FixedVal) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "0fx"}
}

func (f *FixedVal) GetSpecificType() ast.PrimitiveValueType {
	return f.SpecificType
}

var EmptyReturn = Types{}

type Types []Type

func (ts *Types) GetType() Type {
	if len(*ts) == 0 {
		return (&Invalid{}).GetType()
	}
	return (*ts)[0]
}

func (ts *Types) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "TYPES"}
}

func (rt *Types) Eq(otherRT *Types) bool {
	typesSame := true
	if len(*rt) == len(*otherRT) {
		for i, v := range *rt {
			if !TypeEquals(v, (*otherRT)[i]) {
				typesSame = false
				break
			}
		}
	} else {
		typesSame = false
	}
	return typesSame
}

type FunctionVal struct {
	Params  Types
	Returns Types
}

func (f *FunctionVal) GetType() Type {
	return NewFunctionType(f.Params, f.Returns)
}

func (f *FunctionVal) GetReturns() Types {
	return f.Returns
}

func (f *FunctionVal) GetDefault() *ast.LiteralExpr {
	src := lua.StringBuilder{}
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
	return NAType
}

func (u *Unknown) GetDefault() *ast.LiteralExpr {
	return &ast.LiteralExpr{Value: "unknown"}
}
