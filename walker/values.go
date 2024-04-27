package walker

import (
	"hybroid/ast"
	"hybroid/parser"
)

type Value interface {
	GetType() TypeVal
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

type NamespaceVal struct {
	Name string
}

func (n NamespaceVal) GetType() TypeVal {
	return TypeVal{Type: ast.Namespace}
}

/*
	type ListMemberVal struct {
		Val Value
		Owner ListVal
	}

	func (lm ListMemberVal) GetType() ast.PrimitiveValueType {
		return lm.Val.GetType()
	}
*/
type MapMemberVal struct {
	Var   VariableVal
	Owner MapVal
}

func (mm MapMemberVal) GetType() TypeVal {
	return mm.Var.GetType()
}

type MapVal struct {
	MemberType TypeVal
	Members    map[string]MapMemberVal
}

func (m MapVal) GetType() TypeVal {
	return TypeVal{Type: ast.Map, WrappedType: &m.MemberType}
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
	return TypeVal{Type: ast.List, WrappedType: &l.ValueType}
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
	return TypeVal{Type: ast.Number}
}

type DirectiveVal struct{}

func (d DirectiveVal) GetType() TypeVal {
	return TypeVal{Type: 0}
}

type FixedVal struct {
	SpecificType ast.PrimitiveValueType
}

func (f FixedVal) GetType() TypeVal {
	return TypeVal{Type: ast.FixedPoint}
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

type TypeVal struct { // fn(text, text) text
	WrappedType *TypeVal
	Type        ast.PrimitiveValueType
	Params      []TypeVal
	Returns     ReturnType
}

func (t TypeVal) Eq(otherT TypeVal) bool {
	paramsAreSame := true
	if len(t.Params) == len(otherT.Params) {
		for i, v := range t.Params {
			if !v.Eq(otherT.Params[i]) {
				paramsAreSame = false
				break
			}
		}
	} else {
		paramsAreSame = false
	}

	if (otherT.WrappedType == nil || t.WrappedType == nil) && !(otherT.WrappedType == nil && t.WrappedType == nil) {
		return false
	} else if otherT.WrappedType == nil && t.WrappedType == nil {
		return (t.Type == otherT.Type) && paramsAreSame && (t.Returns.Eq(&otherT.Returns))
	}

	return (t.Type == 0 || otherT.Type == 0 || t.Type == otherT.Type) && (t.WrappedType.Eq(*otherT.WrappedType)) && paramsAreSame && (t.Returns.Eq(&otherT.Returns))
}

func (t TypeVal) GetType() TypeVal {
	return t
}

type FunctionVal struct { // fn test(param map<fixed>)
	params    []TypeVal
	returnVal ReturnType
}

func (f FunctionVal) GetType() TypeVal {
	return TypeVal{Type: ast.Func, Params: f.params, Returns: f.returnVal}
}

func (f FunctionVal) GetReturnType() ReturnType {
	return f.returnVal
}

type CallVal struct {
	types ReturnType
}

func (f CallVal) GetType() TypeVal {
	if len(f.types.values) == 1 {
		return f.types.values[0]
	}
	return TypeVal{Type: ast.Invalid, Returns: f.types}
}

type BoolVal struct{}

func (b BoolVal) GetType() TypeVal {
	return TypeVal{Type: ast.Bool}
}

type StringVal struct{}

func (b StringVal) GetType() TypeVal {
	return TypeVal{Type: ast.String}
}

type NilVal struct{}

func (n NilVal) GetType() TypeVal {
	return TypeVal{Type: ast.Nil}
}

type Invalid struct{}

func (u Invalid) GetType() TypeVal {
	return TypeVal{Type: ast.Invalid}
}

type Unknown struct{}

func (u Unknown) GetType() TypeVal {
	return TypeVal{Type: 0}
}
