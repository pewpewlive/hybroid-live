package walker

import (
	"hybroid/ast"
	"hybroid/lexer"
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

type MapVal struct {
	MemberType ast.PrimitiveValueType
	Members map[string]VariableVal
}

func (m MapVal) GetType() ast.PrimitiveValueType {
	return ast.Map
}

func (m MapVal) GetMemberType() ast.PrimitiveValueType {
	var prev *VariableVal
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
	ValuesType ast.PrimitiveValueType
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

type ReturnValue struct {
	values []ast.PrimitiveValueType
}

func (n ReturnValue) GetType() ast.PrimitiveValueType {
	return 0
}

type CallVal struct {
	params     []lexer.Token
	returnVals []ReturnValue
}

func (f CallVal) GetType() ast.PrimitiveValueType {
	return 0
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
