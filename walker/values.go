package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/generators/lua"
	"hybroid/lexer"
	"hybroid/parser"
)

type Value interface {
	GetType() TypeVal
	GetDefault() ast.LiteralExpr
}

type Container interface {
	GetFields() map[string]VariableVal
	GetMethods() map[string]VariableVal
	Contains(name string) (Value, int, bool)
}

type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsConst bool
	Node    ast.Node
}

func (v VariableVal) GetType() TypeVal {
	return v.Value.GetType()
}

func (v VariableVal) GetDefault() ast.LiteralExpr {
	return v.Value.GetDefault()
}

type StructTypeVal struct {
	Name         lexer.Token
	Params       []TypeVal
	Fields       []VariableVal
	FieldIndexes map[string]int
	Methods      map[string]VariableVal
	IsUsed       bool
}

func (st StructTypeVal) GetField(name string) (*VariableVal, bool) {
	if found, ok := FindFromList(st.Fields, name); ok {
		return found, true
	}

	return nil, false
}

func (st StructTypeVal) GetFields() map[string]VariableVal {
	fields := map[string]VariableVal{}

	for _, v := range st.Fields {
		fields[v.Name] = v
	}
	
	return fields
}

func (st StructTypeVal) GetMethods() map[string]VariableVal {
	return st.Methods
}

func FindFromList(list []VariableVal, name string) (*VariableVal, bool) {
	for _, v := range list {
		if v.Name == name {
			return &v, true
		}
	}
	return nil, false
}

func (st StructTypeVal) Contains(name string) (Value, int, bool) {
	if variable, found := FindFromList(st.Fields, name); found {
		return *variable, st.FieldIndexes[name], true
	}

	if variable, found := st.Methods[name]; found {
		return variable, st.FieldIndexes[name], true
	}

	return nil, -1, false
}

func (st StructTypeVal) GetType() TypeVal {
	return TypeVal{Type: ast.Struct, Name: st.Name.Lexeme}
}

func (st StructTypeVal) GetDefault() ast.LiteralExpr {
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

func (s StructVal) GetType() TypeVal {
	return s.Type.GetType()
}

func (s StructVal) GetFields() map[string]VariableVal {
	return s.Type.GetFields()
}

func (s StructVal) GetMethods() map[string]VariableVal {
	return s.Type.GetMethods()
}

func (s StructVal) Contains(name string) (Value, int, bool) {
	return s.Type.Contains(name)
}

func (s StructVal) GetDefault() ast.LiteralExpr {
	return s.Type.GetDefault()
}

type NamespaceVal struct {
	Name         string
	Fields       map[string]VariableVal
	Methods      map[string]VariableVal
	FieldIndexes map[string]int
}

func (n NamespaceVal) GetType() TypeVal {
	return TypeVal{Name: "namespace", Type: ast.Namespace}
}

func (n NamespaceVal) GetFields() map[string]VariableVal {
	return n.Fields
}

func (n NamespaceVal) GetMethods() map[string]VariableVal {
	return n.Methods
}

func (n NamespaceVal) Contains(name string) (Value, int, bool) {
	if variable, found := n.Fields[name]; found {
		return variable, n.FieldIndexes[name], true
	}

	if variable, found := n.Methods[name]; found {
		return variable, n.FieldIndexes[name], true
	}

	return nil, -1, false
}

func (n NamespaceVal) GetDefault() ast.LiteralExpr {
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

type MapMemberVal struct {
	Var   VariableVal
	Owner MapVal
}

func (mm MapMemberVal) GetType() TypeVal {
	return mm.Var.GetType()
}

func (mm MapMemberVal) GetDefault() ast.Node {
	return mm.Var.GetDefault()
}

type MapVal struct {
	MemberType TypeVal
	Members    map[string]MapMemberVal
}

func (m MapVal) GetType() TypeVal {
	return TypeVal{Name: "map", Type: ast.Map, WrappedType: &m.MemberType}
}

func (m MapVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "{}"}
}

func (l MapVal) GetContentsValueType() TypeVal {
	valTypes := []TypeVal{}
	index := 0
	if len(l.Members) == 0 {
		return TypeVal{Type: 0}
	}
	for _, v := range l.Members {
		if index == 0 {
			valTypes = append(valTypes, v.GetType())
			index++
			continue
		}
		valTypes = append(valTypes, v.GetType())
		prev, curr := index-1, len(valTypes)-1
		if !(parser.IsFx(valTypes[prev].Type) && parser.IsFx(valTypes[curr].Type)) && valTypes[prev].Type != valTypes[curr].Type {
			return TypeVal{Type: ast.Invalid}
		}
		index++
	}
	if parser.IsFx(valTypes[0].Type) {
		return TypeVal{Type: ast.FixedPoint}
	}
	return valTypes[0]
}

type ListVal struct {
	ValueType TypeVal
	Values    []Value
}

func (l ListVal) GetType() TypeVal {
	return TypeVal{Name: "list", Type: ast.List, WrappedType: &l.ValueType}
}

func (l ListVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "{}"}
}

func (l ListVal) GetContentsValueType() TypeVal {
	valTypes := []TypeVal{}
	index := 0
	if len(l.Values) == 0 {
		return TypeVal{Type: ast.Invalid}
	}
	for _, v := range l.Values {
		if index == 0 {
			valTypes = append(valTypes, v.GetType())
			index++
			continue
		}
		valTypes = append(valTypes, v.GetType())
		prev, curr := index-1, len(valTypes)-1
		if !(parser.IsFx(valTypes[prev].Type) && parser.IsFx(valTypes[curr].Type)) && valTypes[prev].Type != valTypes[curr].Type {
			return TypeVal{Type: ast.Invalid}
		}
		index++
	}
	if parser.IsFx(valTypes[0].Type) {
		return TypeVal{Type: ast.FixedPoint}
	}
	return valTypes[0]
}

type NumberVal struct{}

func (n NumberVal) GetType() TypeVal {
	return TypeVal{Type: ast.Number, Name: "number"}
}

func (n NumberVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "0"}
}

type DirectiveVal struct{}

func (d DirectiveVal) GetType() TypeVal {
	return TypeVal{Type: 0, Name: "directive"}
}

func (d DirectiveVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "DEFAULT_DIRECTIVE_CALL"}
}

type FixedVal struct {
	SpecificType ast.PrimitiveValueType
}

func (f FixedVal) GetType() TypeVal {
	return TypeVal{Type: ast.FixedPoint, Name: "fixed"}
}

func (f FixedVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "0fx"}
}

func (f FixedVal) GetSpecificType() ast.PrimitiveValueType {
	return f.SpecificType
}

type ReturnType struct {
	values []TypeVal
}

func (rt *ReturnType) Eq(otherRT *ReturnType) bool {
	typesSame := true
	if len(rt.values) == len(otherRT.values) {
		for i, v := range rt.values {
			if !v.Eq(otherRT.values[i]) {
				typesSame = false
				break
			}
		}
	} else {
		typesSame = false
	}
	return typesSame
}

func (n ReturnType) GetType() TypeVal {
	return TypeVal{Type: 0}
}

type TypeVal struct {
	WrappedType *TypeVal
	Name        string
	Type        ast.PrimitiveValueType
	Params      *[]TypeVal
	Returns     ReturnType
}

func (t TypeVal) Eq(otherT TypeVal) bool {
	paramsAreSame := true

	if (otherT.Params == nil || t.Params == nil) && !(otherT.Params == nil && t.Params == nil) {
		return false
	} else {
		if otherT.Params == nil && t.Params == nil {
			paramsAreSame = true
		} else if len(*t.Params) == len(*otherT.Params) {
			for i, v := range *t.Params {
				if !v.Eq((*otherT.Params)[i]) {
					paramsAreSame = false
					break
				}
			}
		} else {
			paramsAreSame = false
		}
	}

	if (otherT.WrappedType == nil || t.WrappedType == nil) && !(otherT.WrappedType == nil && t.WrappedType == nil) {
		return false
	} else if otherT.WrappedType == nil && t.WrappedType == nil {
		return (t.Type == otherT.Type) && paramsAreSame && (t.Returns.Eq(&otherT.Returns))
	}

	return (t.Type == 0 || otherT.Type == 0 || t.Type == otherT.Type) && (t.WrappedType.Eq(*otherT.WrappedType)) && paramsAreSame && (t.Returns.Eq(&otherT.Returns))
}

func (t TypeVal) ToString() string {
	src := lua.StringBuilder{}

	src.Append(t.Name)

	if t.Params != nil {
		src.Append("(")
		for i := range *t.Params {
			if i == len(*t.Params)-1 {
				src.Append((*t.Params)[i].ToString())
			} else {
				src.Append((*t.Params)[i].ToString(), ", ")
			}
		}
		src.Append(")")
		if len(t.Returns.values) != 0 {
			src.Append(" ")
			for i := range t.Returns.values {
				if i == len(t.Returns.values)-1 {
					src.Append(t.Returns.values[i].ToString())
				} else {
					src.Append(t.Returns.values[i].ToString(), ", ")
				}
			}
		}
	}

	if t.WrappedType != nil {
		src.Append("<", t.WrappedType.ToString(), ">")
	}

	return src.String()
}

func (t TypeVal) GetType() TypeVal {
	return t
}

func (t TypeVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "DEFAULT_TYPE_VALUE"}
}

type FunctionVal struct { // fn test(param map<fixed>)
	params    []TypeVal
	returnVal ReturnType
}

func (f FunctionVal) GetType() TypeVal {
	return TypeVal{Name: "function", Type: ast.Func, Params: &f.params, Returns: f.returnVal}
}

func (f FunctionVal) GetReturnType() ReturnType {
	return f.returnVal
}

func (f FunctionVal) GetDefault() ast.LiteralExpr {
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

type CallVal struct {
	types ReturnType
}

func (f CallVal) GetType() TypeVal {
	if len(f.types.values) == 1 {
		return f.types.values[0]
	}
	return TypeVal{Name: "invalid", Type: ast.Invalid, Returns: f.types}
}

func (f CallVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "DEFAULT_CALL_VALUE"}
}

type BoolVal struct{}

func (b BoolVal) GetType() TypeVal {
	return TypeVal{Name: "boolean", Type: ast.Bool}
}

func (b BoolVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "false"}
}

type StringVal struct{}

func (s StringVal) GetType() TypeVal {
	return TypeVal{Name: "string", Type: ast.String}
}

func (s StringVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "\"\""}
}

type NilVal struct{}

func (n NilVal) GetType() TypeVal {
	return TypeVal{Name: "nil", Type: ast.Nil}
}

func (n NilVal) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "nil"}
}

type Invalid struct{}

func (u Invalid) GetType() TypeVal {
	return TypeVal{Name: "invalid", Type: ast.Invalid}
}

func (n Invalid) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "invalid"}
}

type Unknown struct{}

func (u Unknown) GetType() TypeVal {
	return TypeVal{Name: "unknown", Type: 0}
}

func (u Unknown) GetDefault() ast.LiteralExpr {
	return ast.LiteralExpr{Value: "unknown"}
}
