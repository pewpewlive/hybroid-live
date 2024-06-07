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
	GetDefault() ast.LiteralExpr
}

type Container interface {
	Value
	GetFields() map[string]VariableVal
	GetMethods() map[string]VariableVal
	AddField(variable VariableVal)
	AddMethod(variable VariableVal)
	Contains(name string) (Value, int, bool)
}

type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsConst bool
	Node    ast.Node
}

func (v *VariableVal) GetType() Type {
	return v.Value.GetType()
}

func (v *VariableVal) GetDefault() ast.LiteralExpr {
	return v.Value.GetDefault()
}

type StructTypeVal struct {
	Name         lexer.Token
	Params       []Type
	Fields       []VariableVal
	FieldIndexes map[string]int
	Methods      map[string]VariableVal
	IsUsed       bool
}

func (st *StructTypeVal) GetField(name string) (*VariableVal, bool) {
	if found, ok := FindFromList(st.Fields, name); ok {
		return found, true
	}

	return nil, false
}

func (st *StructTypeVal) GetFields() map[string]VariableVal {
	fields := map[string]VariableVal{}

	for _, v := range st.Fields {
		fields[v.Name] = v
	}

	return fields
}

func (st *StructTypeVal) GetMethods() map[string]VariableVal {
	return st.Methods
}

func (st *StructTypeVal) AddField(variable VariableVal) {
	st.Fields = append(st.Fields, variable)

	index := 0
	for range st.FieldIndexes {
		index++
	}
	st.FieldIndexes[variable.Name] = index + 1
}

func (st *StructTypeVal) AddMethod(variable VariableVal) {
	st.Methods[variable.Name] = variable
}

func FindFromList(list []VariableVal, name string) (*VariableVal, bool) {
	for _, v := range list {
		if v.Name == name {
			return &v, true
		}
	}
	return nil, false
}

func (st *StructTypeVal) Contains(name string) (Value, int, bool) {
	if variable, found := FindFromList(st.Fields, name); found {
		return variable, st.FieldIndexes[name], true
	}

	if variable, found := st.Methods[name]; found {
		return &variable, st.FieldIndexes[name], true
	}

	return nil, -1, false
}

func (st *StructTypeVal) GetType() Type {
	return Type{Type: ast.Struct, Name: st.Name.Lexeme}
}

func (st *StructTypeVal) GetDefault() ast.LiteralExpr {
	src := lua.StringBuilder{}

	src.WriteString("{")

	index := 0
	for _, v := range st.Fields {
		src.WriteString(v.GetDefault().Value)
		if index != len(st.Fields)-1 {
			src.WriteString(", ")
		}
		index += 1
	}
	src.WriteString("}")

	return ast.LiteralExpr{Value: src.String()}
}

type StructVal struct {
	Type *StructTypeVal
}

func (s *StructVal) GetType() Type {
	return s.Type.GetType()
}

func (s *StructVal) GetFields() map[string]VariableVal {
	return s.Type.GetFields()
}

func (s *StructVal) GetMethods() map[string]VariableVal {
	return s.Type.GetMethods()
}

func (s *StructVal) Contains(name string) (Value, int, bool) {
	return s.Type.Contains(name)
}

func (s *StructVal) GetDefault() ast.LiteralExpr {
	return s.Type.GetDefault()
}

func (s *StructVal) AddField(variable VariableVal) {
	s.Type.AddField(variable)
}

func (s *StructVal) AddMethod(variable VariableVal) {
	s.Type.AddMethod(variable)
}

type AnonStructTypeVal struct {
	Fields []VariableVal
}

func (astv *AnonStructTypeVal) GetField(name string) (*VariableVal, bool) {
	if found, ok := FindFromList(astv.Fields, name); ok {
		return found, true
	}

	return nil, false
}

func (astv *AnonStructTypeVal) GetFields() map[string]VariableVal {
	fields := map[string]VariableVal{}

	for _, v := range astv.Fields {
		fields[v.Name] = v
	}

	return fields
}

func (astv *AnonStructTypeVal) GetMethods() map[string]VariableVal {
	return astv.GetFields()
}

func (astv *AnonStructTypeVal) Contains(name string) (Value, int, bool) {
	if variable, found := FindFromList(astv.Fields, name); found {
		return variable, -1, true
	}

	return nil, -1, false
}

func (astv *AnonStructTypeVal) AddField(variable VariableVal) {
	astv.Fields = append(astv.Fields, variable)
}

func (astv *AnonStructTypeVal) AddMethod(variable VariableVal) {
}

func (astv *AnonStructTypeVal) GetType() Type {
	types := make(map[string]Type, len(astv.Fields))

	for _, v := range astv.Fields {
		types[v.Name] = v.GetType()
	}

	return Type{Type: ast.AnonStruct, Fields: types}
}

func (astv *AnonStructTypeVal) GetDefault() ast.LiteralExpr {
	src := lua.StringBuilder{}

	src.WriteString("{")

	index := 0
	for _, v := range astv.Fields {
		src.WriteString(v.GetDefault().Value)
		if index != len(astv.Fields)-1 {
			src.WriteString(", ")
		}
		index += 1
	}
	src.WriteString("}")

	return ast.LiteralExpr{Value: src.String()}
}

type EnvironmentVal struct {
	Path         string
	Name         string
	Fields       map[string]VariableVal
	Methods      map[string]VariableVal
	FieldIndexes map[string]int
}

func (n *EnvironmentVal) GetType() Type {
	return Type{Name: "namespace", Type: ast.Namespace}
}

func (n EnvironmentVal) GetDefault() ast.LiteralExpr {
	src := lua.StringBuilder{}

	src.WriteString("{")

	index := 0
	for _, v := range n.Fields {
		src.WriteString(v.GetDefault().Value)
		if index != len(n.Fields)-1 {
			src.WriteString(", ")
		}
		index += 1
	}
	src.WriteString("}")

	return ast.LiteralExpr{Value: src.String()}
}

type PathVal struct {
	Path string
	Env  ast.EnvExpr
}

type MapVal struct {
	MemberType Type
	Members    []Value
}

func (m *MapVal) GetType() Type {
	return Type{Name: "map", Type: ast.Map, WrappedType: &m.MemberType}
}

func (m *MapVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "{}"}
}

func GetContentsValueType(values []Value) Type {
	valTypes := []Type{}
	index := 0
	if len(values) == 0 {
		return Type{Type: 0}
	}
	for _, v := range values {
		if index == 0 {
			valTypes = append(valTypes, v.GetType())
			index++
			continue
		}
		valTypes = append(valTypes, v.GetType())
		prev, curr := index-1, len(valTypes)-1
		if !(parser.IsFx(valTypes[prev].Type) && parser.IsFx(valTypes[curr].Type)) && valTypes[prev].Type != valTypes[curr].Type {
			return Type{Type: ast.Invalid}
		}
		index++
	}
	if parser.IsFx(valTypes[0].Type) {
		return Type{Type: ast.FixedPoint}
	}
	return valTypes[0]
}

type ListVal struct {
	ValueType Type
	Values    []Value
}

func (l *ListVal) GetType() Type {
	return Type{Name: "list", Type: ast.List, WrappedType: &l.ValueType}
}

func (l *ListVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "{}"}
}

type NumberVal struct{}

func (n *NumberVal) GetType() Type {
	return Type{Type: ast.Number, Name: "number"}
}

func (n *NumberVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "0"}
}

type DirectiveVal struct{}

func (d *DirectiveVal) GetType() Type {
	return Type{Type: 0, Name: "directive"}
}

func (d *DirectiveVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "DEFAULT_DIRECTIVE_CALL"}
}

type FixedVal struct {
	SpecificType ast.PrimitiveValueType
}

func (f *FixedVal) GetType() Type {
	return Type{Type: ast.FixedPoint, Name: "fixed"}
}

func (f *FixedVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "0fx"}
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

func (ts *Types) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "TYPES"}
}

func (rt *Types) Eq(otherRT *Types) bool {
	typesSame := true
	if len(*rt) == len(*otherRT) {
		for i, v := range *rt {
			if !TypeEquals(&v, &(*otherRT)[i]) {
				typesSame = false
				break
			}
		}
	} else {
		typesSame = false
	}
	return typesSame
}

type Type struct {
	WrappedType *Type
	Name        string
	Type        ast.PrimitiveValueType
	Params      []Type
	Returns     Types
	Fields      map[string]Type
}

func TypeEquals(t *Type, otherT *Type) bool {
	if t == nil && otherT == nil {
		return true
	} else if t != nil && otherT != nil {
	} else {
		return false
	}

	// unknown checking
	if t.Type == 0 || otherT.Type == 0 {
		return true
	}

	// name checking
	if t.Name != otherT.Name {
		return false
	}

	// param checking
	if len(t.Params) != len(otherT.Params) {
		return false
	}
	for i, v := range t.Params {
		if !TypeEquals(&v, &otherT.Params[i]) {
			return false
		}
	}

	// return checking
	if len(t.Returns) != len(otherT.Returns) {
		return false
	}
	for i, v := range t.Returns {
		if !TypeEquals(&v, &otherT.Returns[i]) {
			return false
		}
	}

	// field checking
	if len(t.Fields) != len(otherT.Fields) {
		return false
	}
	for i, v := range t.Fields {
		temp := otherT.Fields[i]
		if !TypeEquals(&v, &temp) {
			return false
		}
	}

	// pvt checking
	if t.Type != otherT.Type {
		return false
	}

	// wrapped type checking
	if !TypeEquals(t.WrappedType, otherT.WrappedType) {
		return false
	}

	return true
}

func (t *Type) ToString() string {
	src := lua.StringBuilder{}

	src.Append(t.Name)

	if t.Params != nil {
		src.Append("(")
		for i := range t.Params {
			if i == len(t.Params)-1 {
				src.Append(t.Params[i].ToString())
			} else {
				src.Append(t.Params[i].ToString(), ", ")
			}
		}
		src.Append(")")
		if len(t.Returns) != 0 {
			src.Append(" ")
			for i := range t.Returns {
				if i == len(t.Returns)-1 {
					src.Append(t.Returns[i].ToString())
				} else {
					src.Append(t.Returns[i].ToString(), ", ")
				}
			}
		}
	}

	if t.WrappedType != nil {
		src.Append("<", t.WrappedType.ToString(), ">")
	}

	return src.String()
}

type FunctionVal struct {
	params    []Type
	returnVal Types
}

func (f *FunctionVal) GetType() Type {
	return Type{Name: "function", Type: ast.Func, Params: f.params, Returns: f.returnVal}
}

func (f *FunctionVal) GetReturns() Types {
	return f.returnVal
}

func (f *FunctionVal) GetDefault() ast.LiteralExpr {
	src := lua.StringBuilder{}
	src.WriteString("function(")
	for i := range f.params {
		src.WriteString(fmt.Sprintf("param%v", i))
		if i != len(f.params)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(") end")
	return ast.LiteralExpr{Value: src.String()}
}

type BoolVal struct{}

func (b *BoolVal) GetType() Type {
	return Type{Name: "bool", Type: ast.Bool}
}

func (b *BoolVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "false"}
}

type StringVal struct{}

func (s *StringVal) GetType() Type {
	return Type{Name: "string", Type: ast.String}
}

func (s *StringVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "\"\""}
}

type Invalid struct{}

func (u *Invalid) GetType() Type {
	return Type{Name: "invalid", Type: ast.Invalid}
}

func (n *Invalid) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "invalid"}
}

type Unknown struct{}

func (u *Unknown) GetType() Type {
	return Type{Name: "unknown", Type: 0}
}

func (u *Unknown) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "unknown"}
}
