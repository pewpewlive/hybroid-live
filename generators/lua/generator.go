package lua

import (
	"fmt"
	"hybroid/lexer"
	"hybroid/parser"
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

func (s *Scope) GetVariable(name string) Value {
	scope := s.Resolve(name)

	return scope.Variables[name]
}

func (s *Scope) AssignVariable(name string, value Value) Value {
	scope := s.Resolve(name)

	// TODO: check if the value is a constant

	scope.Variables[name] = value

	return value
}

func (s *Scope) DeclareVariable(name string, value Value) Value {
	//if _, found := s.variables[name]; found {
	// error: variable with this name already exists
	//}

	s.Variables[name] = value

	return value
}

func (s *Scope) Resolve(name string) *Scope {
	if _, found := s.Variables[name]; found {
		return s
	}

	//if s.parent == nil {
	// error: variable does not exist
	//}

	return s.Parent.Resolve(name)
}

func (g *Global) GetForeignType(str string) Value {
	return g.foreignTypes[str]
}

func (gen *Generator) Program(program parser.Program, environment *Scope) Value {
	var lastEvaluated Value

	for _, node := range program.Body {
		lastEvaluated = gen.Generate(node, environment)
	}

	return lastEvaluated
}

func (gen *Generator) variableDeclaration(declaration parser.Node, scope *Scope) Value {
	var value Value

	if declaration.Expression == nil {
		gen.error(declaration.Token, "expected expression after declaration")
	} else {
		value = gen.Generate(*declaration.Expression, scope)
	}

	isLocal := declaration.Token.Type == lexer.Let

	if isLocal {
		gen.Src += fmt.Sprintf("local %s = %s\n", declaration.Identifier, value.Val)
	} else {
		gen.Src += fmt.Sprintf("%s = %s\n", declaration.Identifier, value.Val)
	}

	return scope.DeclareVariable(declaration.Identifier, value)
}

func (gen *Generator) binaryExpr(node parser.Node, scope *Scope) Value {
	src := gen.Generate(*node.Left, scope).Val
	src += fmt.Sprintf(" %s ", node.Token.Lexeme)
	src += gen.Generate(*node.Right, scope).Val

	return Value{parser.Nil, src}
}

func (gen *Generator) literalExpr(node parser.Node) Value {
	src := fmt.Sprintf("%v", node.Value)

	return Value{node.ValueType, src}
}

func (gen *Generator) identifierExpr(node parser.Node, scope *Scope) Value {
	scope.Resolve(node.Identifier)
	return Value{Type: node.ValueType, Val: node.Identifier}
}

func (gen *Generator) groupingExpr(node parser.Node, scope *Scope) Value {
	src := "("
	value := gen.Generate(*node.Expression, scope)
	src += value.Val
	src += ")"

	return Value{value.Type, src}
}

func (gen *Generator) listExpr(node parser.Node, scope *Scope) Value {
	nodes, _ := node.Value.([]parser.Node)

	src := "{"
	for i, expr := range nodes {
		src += gen.Generate(expr, scope).Val

		if i != len(nodes)-1 {
			src += ", "
		}
	}
	src += "}"

	return Value{parser.List, src}
}

func (gen *Generator) Generate(node parser.Node, environment *Scope) Value {
	scope := environment
	switch node.NodeType {
	case parser.LiteralExpr:
		return gen.literalExpr(node)
	case parser.Prog:
		return gen.Program(*node.Program, scope)
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
	}

	return Value{}
}
