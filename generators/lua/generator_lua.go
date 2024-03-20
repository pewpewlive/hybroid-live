package lua

import (
	"fmt"
	"hybroid/lexer"
	"hybroid/parser"
	"math"
)

type GenError struct {
	Token   lexer.Token
	Message string
}

func (ge *GenError) generatorError() string {
	return fmt.Sprintf("Error: %v, at line: %v (%v)", ge.Message, ge.Token.Location.LineStart, ge.Token.ToString())
}

func (gen *Generator) error(token lexer.Token, message string) {
	gen.Errors = append(gen.Errors, GenError{token, message})
}

type Generator struct {
	Errors []GenError
	Src    string
	ctx    Context
}

type Value struct {
	Type parser.PrimitiveValueType
	Val  string
}

type Global struct {
	Scope        Scope
	foreignTypes map[string]Value
}

type Scope struct {
	Global    *Global
	Parent    *Scope
	Variables map[string]Value
}

// func (gen *Generator) validateOperands(left *Value, right *Value) bool {
// 	if (left.Type == 0 || left.Type == Nil) || (right.Type == 0 || right.Type == Nil) {
// 		gen.error(left.Token, "cannot perform arithmetic on nil value")
// 		return false
// 	} else if left.Type == Undefined || right.Type == Undefined {
// 		gen.error(left.Token, "cannot perform arithmetic on undefined value")
// 		return false
// 	} else {
// 		if (left.Type == List || left.Type == Map || left.Type == String || left.Type == Bool || left.ValTypeueType == Entity || left.Type == Struct) ||
// 			(right.Type == List || right.Type == Map || right.Type == String || right.Type == Bool || right.Type == Entity || right.Type == Struct) {
// 				gen.error(left.Token, "cannot perform arithmetic on a non-number value")
// 			return false
// 		}
// 	}
// 	return true
// }

type Context int

const (
	None Context = iota
	Expression
)

func (gen *Generator) GetErrors() []GenError {
	return gen.Errors
}

func (s *Scope) GetVariable(name string) Value {

	scope := s.Resolve(name)

	return scope.Variables[name]
}

func (s *Scope) AssignVariable(name string, value Value) (Value, bool) {
	scope := s.Resolve(name)

	// TODO: check if the value is a constant
	if scope == nil {
		return Value{}, false
	}

	scope.Variables[name] = value

	return value, true
}

func (s *Scope) DeclareVariable(name string, value Value) (Value, bool) {
	scope := s.Resolve(name)

	if scope == nil {
		s.Variables[name] = value
		return value, true
	} else {
		return Value{}, false
	}
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

func (gen *Generator) Generate(program parser.Program, environment *Scope) Value {
	var lastEvaluated Value

	for _, node := range program.Body {
		lastEvaluated = gen.GenerateNode(node, environment)
		gen.Src += lastEvaluated.Val + "\n"
	}

	return lastEvaluated
}

func (gen *Generator) variableDeclaration(declaration parser.Node, scope *Scope) Value {
	var value Value

	if declaration.Expression == nil {
		gen.error(declaration.Token, "expected expression after declaration")
	} else {
		value = gen.GenerateNode(*declaration.Expression, scope)
	}

	isLocal := declaration.Token.Type == lexer.Let
	src := ""
	if isLocal {
		src += fmt.Sprintf("local %s = %s", declaration.Identifier, value.Val)
	} else {
		src += fmt.Sprintf("%s = %s", declaration.Identifier, value.Val)
	}

	if _, success := scope.DeclareVariable(declaration.Identifier, value); !success {
		gen.error(lexer.Token{Lexeme: declaration.Identifier, Location: declaration.Token.Location},
			"cannot declare a value in the same scope twice")
	}

	return Value{Type: parser.Nil, Val: src}
}

func (gen *Generator) binaryExpr(node parser.Node, scope *Scope) Value {
	src := gen.GenerateNode(*node.Left, scope).Val
	src += fmt.Sprintf(" %s ", node.Token.Lexeme)
	src += gen.GenerateNode(*node.Right, scope).Val

	return Value{parser.Nil, src}
}

func (gen *Generator) literalExpr(node parser.Node) Value {
	var src string

	switch node.ValueType {
	case parser.String:
		src = "\"" + fmt.Sprintf("%v", node.Value) + "\""
	case parser.Fixed:
		src = fixedToFx(node.Value.(float64)) + "fx"
	case parser.FixedPoint:
		src = fmt.Sprintf("%vfx", node.Value)
	default:
		src = fmt.Sprintf("%v", node.Value)
	}

	return Value{node.ValueType, src}
}

func fixedToFx(float float64) string {
	abs_float := math.Abs(float)
	integer := math.Floor(abs_float)
	if integer > (2 << 51) {
		integer = (2 << 51)
	}
	var sign string
	if float < 0 {
		sign = "-"
	} else {
		sign = ""
	}

	frac := math.Floor((abs_float - integer) * 4096)
	frac_str := ""
	if frac != 0 {
		frac_str = "." + fmt.Sprintf("%d", frac)
	}

	// sign + int + frac_str + "fx"
	return fmt.Sprintf("%s%d%s", sign, integer, frac_str)
}

func (gen *Generator) identifierExpr(node parser.Node, scope *Scope) Value {
	scope.Resolve(node.Identifier)
	return Value{Type: node.ValueType, Val: node.Identifier}
}

func (gen *Generator) groupingExpr(node parser.Node, scope *Scope) Value {
	src := "("
	value := gen.GenerateNode(*node.Expression, scope)
	src += value.Val
	src += ")"

	return Value{value.Type, src}
}

func (gen *Generator) listExpr(node parser.Node, scope *Scope) Value {
	nodes, _ := node.Value.([]parser.Node)

	src := "{"
	for i, expr := range nodes {
		src += gen.GenerateNode(expr, scope).Val

		if i != len(nodes)-1 {
			src += ", "
		}
	}
	src += "}"

	return Value{parser.List, src}
}

func (gen *Generator) assignmentExpr(node parser.Node, scope *Scope) Value {
	if node.Expression.NodeType != parser.Identifier {
		gen.error(node.Expression.Token, "expected an identifier to assign to")
	}

	src := node.Expression.Identifier
	value := gen.GenerateNode(*node.Right, scope)
	if _, success := scope.AssignVariable(node.Expression.Identifier, value); !success { // for checking variable's existence and const checking
		gen.error(node.Expression.Token, "cannot assign a value to an undeclared variable")
	}
	src += fmt.Sprintf(" = %v", value.Val)

	return Value{value.Type, src}
}

func (gen *Generator) unaryExpr(node parser.Node, scope *Scope) Value {
	value := gen.GenerateNode(*node.Right, scope)
	src := fmt.Sprintf("%s%s", node.Token.Lexeme, value.Val)

	return Value{Type: value.Type, Val: src}
}

func (gen *Generator) functionDeclarationStmt(node parser.Node, scope *Scope) Value {
	return Value{}
}

func (gen *Generator) GenerateNode(node parser.Node, environment *Scope) Value {
	scope := environment
	switch node.NodeType {
	case parser.LiteralExpr:
		return gen.literalExpr(node)
	case parser.Prog:
		return gen.Generate(*node.Program, scope)
	case parser.VariableDeclarationStmt:
		return gen.variableDeclaration(node, scope)
	case parser.BinaryExpr:
		return gen.binaryExpr(node, scope)
	case parser.Identifier:
		return gen.identifierExpr(node, scope)
	case parser.GroupingExpr:
		return gen.groupingExpr(node, scope)
	case parser.ListExpr:
		return gen.listExpr(node, scope)
	case parser.AssignmentExpr:
		return gen.assignmentExpr(node, scope)
	case parser.UnaryExpr:
		return gen.unaryExpr(node, scope)
	case parser.FunctionDeclarationStmt:
		return gen.functionDeclarationStmt(node, scope)
	}

	return Value{}
}
