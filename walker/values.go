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
	MemberType  ast.PrimitiveValueType
	Members     map[string]MapMemberVal
}

func (m MapVal) GetType() ast.PrimitiveValueType {
	return ast.Map
}

func (l MapVal) GetContentsValueType() ast.PrimitiveValueType {
	valTypes := []ast.PrimitiveValueType{}
	index := 0
	for _, v := range l.Members {
		if index == 0 {
			valTypes = append(valTypes, v.GetType())
			index++
			continue
		}
		valTypes = append(valTypes, v.GetType())
		if valTypes[index-1] != valTypes[len(valTypes)-1] {
			return ast.Undefined
		}
		index++
	}
	return valTypes[1]
}

type ListVal struct {
	ValueType   ast.PrimitiveValueType
	Values      []Value
}

func (l ListVal) GetType() ast.PrimitiveValueType {
	return ast.List
}

func (l ListVal) GetContentsValueType() ast.PrimitiveValueType {
	valTypes := []ast.PrimitiveValueType{}
	index := 0
	for _, v := range l.Values {
		if index == 0 {
			valTypes = append(valTypes, v.GetType())
			index++
			continue
		}
		valTypes = append(valTypes, v.GetType())
		if valTypes[index-1] != valTypes[len(valTypes)-1] {
			return ast.Undefined
		}
		index++
	}
	return valTypes[1]
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
	SpecificType ast.PrimitiveValueType
	Val string
}

func (f FixedVal) GetType() ast.PrimitiveValueType {
	return ast.FixedPoint
}

func (f FixedVal) GetSpecificType() ast.PrimitiveValueType {
	return f.SpecificType
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
