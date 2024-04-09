package walker

import (
	"hybroid/ast"
)

type Value interface {
	GetType() ast.PrimitiveValueType
}

type VariableVal struct {
	Name    string
	Value   Value
	IsUsed  bool
	IsConst bool
	Node    ast.Node
}

func (v VariableVal) GetType() ast.PrimitiveValueType {
	return v.Value.GetType()
}

type NamespaceVal struct {
	Name string
}

func (n NamespaceVal) GetType() ast.PrimitiveValueType {
	return ast.Namespace;
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
	Var VariableVal
	Owner MapVal
}

func (mm MapMemberVal) GetType() ast.PrimitiveValueType {
	return mm.Var.GetType()
}

type MapVal struct {
	MemberType ast.PrimitiveValueType
	Members map[string]MapMemberVal
}

func (m MapVal) GetType() ast.PrimitiveValueType {
	return ast.Map
}

func (m MapVal) GetMembersType() ast.PrimitiveValueType {
	var prev *MapMemberVal
	for _, val := range m.Members {
		if prev == nil {
			prev = &val
			continue
		}
		if prev.GetType() != val.GetType() {
			return ast.Undefined
		}
	}
	return prev.GetType()
}

type ListVal struct {
	ValueType ast.PrimitiveValueType
	values []Value
}

func (l ListVal) GetType() ast.PrimitiveValueType {
	return ast.List
}

func (m ListVal) GetValuesType() ast.PrimitiveValueType {
	var prev *Value
	for _, val := range m.values {
		if prev == nil {
			prev = &val
			continue
		}
		if (*prev).GetType() != val.GetType() {
			return ast.Undefined
		}
	}
	return (*prev).GetType()
}

type NumberVal struct {
	Val string
}

func (n NumberVal) GetType() ast.PrimitiveValueType {
	return ast.Number
}

type DirectiveVal struct{}

func (d DirectiveVal) GetType() ast.PrimitiveValueType {
	return 0
}

type FixedVal struct {
	Val string
}

func (f FixedVal) GetType() ast.PrimitiveValueType {
	return ast.FixedPoint
}

type ReturnType struct {
	values []ast.PrimitiveValueType
}

func (n ReturnType) GetType() ast.PrimitiveValueType {
	return 0
}

type FunctionVal struct { 
	params     []ast.Param
	returnVal ReturnType
}

func (f FunctionVal) GetType() ast.PrimitiveValueType {
	return 0
}

func (f FunctionVal) GetReturnType() ReturnType {
	return f.returnVal
}

type BoolVal struct {
	Val string
}

func (b BoolVal) GetType() ast.PrimitiveValueType {
	return ast.Bool
}

type StringVal struct {
	Val string
}

func (b StringVal) GetType() ast.PrimitiveValueType {
	return ast.String
}

type NilVal struct{}

func (n NilVal) GetType() ast.PrimitiveValueType {
	return ast.Nil
}

type Unknown struct {
}

func (u Unknown) GetType() ast.PrimitiveValueType {
	return ast.Undefined
}

type Undefined struct {
}

func (u Undefined) GetType() ast.PrimitiveValueType {
	return 0
}
