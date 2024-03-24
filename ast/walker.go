package ast

import "hybroid/lexer"

type Walker struct {
	nodes   []Node
	current int
	Errors  []Error
}

type Value interface {
	GetType() PrimitiveValueType
	GetToken() lexer.Token
}

type MapVal struct {
	Properties *[]Value
}

func (mv *MapVal) GetType() PrimitiveValueType {
	return Map
}

func (mv *MapVal) GetToken() lexer.Token {
	return
}

type ListVal struct {
	values []Value
}

type NumberVal struct {
	Val string
}

// r, d, f, fx
type FixedVal struct {
	Val string
}

type ReturnValue struct {
	values []Value
}

type CallVal struct {
	params     []string
	returnVals []ReturnValue
} //remove this file, we have the walker.go already and the types why is this here

type MemberVal struct {
}

type Global struct {
	Scope        Scope // look at the values.go
	foreignTypes map[string]Value
}

type Variable struct {
	Name       string
	Val        Value
	IsUsed     bool
	IsConstant bool
	Nodes      []*Node
}

type Scope struct {
	Global    *Global
	Parent    *Scope
	Variables map[string]Variable
}

func (s *Scope) GetVariable(name string) Variable {

	scope := s.Resolve(name)

	return scope.Variables[name]
}

func (s *Scope) AssignVariable(name string, value Variable) (Variable, bool) {
	scope := s.Resolve(name)

	// TODO: check if the value is a constant
	if scope == nil {
		return Variable{}, false
	}

	scope.Variables[name] = value

	return value, true
}

func (s *Scope) DeclareVariable(name string, value Variable) (Variable, bool) {
	if _, found := s.Variables[name]; found {
		return Variable{}, false
	}

	s.Variables[name] = value
	return value, true
}

func (s *Scope) Resolve(name string) *Scope {
	if _, found := s.Variables[name]; found {
		return s
	}

	if s.Parent == nil {
		return nil
	}

	return s.Parent.Resolve(name)
}

func (g *Global) GetForeignType(str string) Value {
	return g.foreignTypes[str]
}

func (w *Walker) error(token lexer.Token, msg string) {

}

func (w *Walker) validateArithmeticOperands(left Variable, right Variable, expr BinaryExpr) bool {
	//fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	switch left.Val.GetType() {
	case Nil:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on nil value")
		return false
	case Undefined:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on undefined value")
		return false
	}

	switch right.Val.GetType() {
	case Nil:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on nil value")
		return false
	case Undefined:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on undefined value")
		return false
	}

	switch left.Val.GetType() {
	case List, Map, String, Bool, Entity, Struct:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	switch right.Val.GetType() {
	case List, Map, String, Bool, Entity, Struct:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	return true
}

func (w *Walker) Walk(nodes []Node) []Node {
	w.nodes = nodes

	newNodes := make([]Node, len(nodes))

	return newNodes
}
