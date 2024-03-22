package lua

import (
	"fmt"
	"hybroid/lexer"
	"hybroid/parser"
	"math"
	"strconv"
	"strings"
)

type GenError struct {
	Token   lexer.Token
	Message string
}

// func (ge *GenError) generatorError() string {
// 	return fmt.Sprintf("Error: %v, at line: %v (%v)", ge.Message, ge.Token.Location.LineStart, ge.Token.ToString())
// }

func (gen *Generator) error(token lexer.Token, message string) {
	gen.Errors = append(gen.Errors, GenError{token, message})
}

type Generator struct {
	Errors []GenError
	Src    strings.Builder
	TabsCount int
}

type Value struct {//properties for maps, structs, entities
	//properties *[]Value
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
	Count     int
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

func (gen Generator) GetErrors() []GenError {
	return gen.Errors
}

func (gen *Generator) GetSrc() string {
	return gen.Src.String()
}

func (gen *Generator) append(strings ...string) {
	for _, str := range strings {
		gen.Src.WriteString(str)
	}
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
	if _, found := s.Variables[name]; found {
		return Value{}, false
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

func (gen *Generator) Generate(program parser.Program, environment *Scope) Value {
	var lastEvaluated Value

	for _, node := range program.Body {
		lastEvaluated = gen.GenerateNode(node, environment)
		gen.append(lastEvaluated.Val, "\n")
	}

	return lastEvaluated
}

func (gen *Generator) binaryExpr(node parser.Node, scope *Scope) Value {
	src := strings.Builder{}
	src.WriteString(gen.GenerateNode(*node.Left, scope).Val)
	src.WriteString(fmt.Sprintf(" %s ", node.Token.Lexeme))
	src.WriteString(gen.GenerateNode(*node.Right, scope).Val)

	return Value{parser.Nil, src.String()}
}

func (gen *Generator) literalExpr(node parser.Node) Value {
	src := strings.Builder{}

	switch node.ValueType {
	case parser.String:
		src.WriteString("\"")
		src.WriteString(fmt.Sprintf("%v", node.Value))
		src.WriteString("\"")
	case parser.Fixed:
		src.WriteString(fixedToFx(node.Value.(string)))
		src.WriteString("fx")
	case parser.FixedPoint:
		src.WriteString(fmt.Sprintf("%vfx", node.Value))
	default:
		src.WriteString(fmt.Sprintf("%v", node.Value))
	}

	return Value{node.ValueType, src.String()}
}

func fixedToFx(floatstr string) string {
	float, _ := strconv.ParseFloat(floatstr, 64)
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
		frac_str = "." + fmt.Sprintf("%v", frac)
	}

	// sign + int + frac_str + "fx"
	return fmt.Sprintf("%s%v%s", sign, integer, frac_str)
}

func (gen *Generator) identifierExpr(node parser.Node, scope *Scope) Value {
	scope.Resolve(node.Identifier)
	return Value{Type: node.ValueType, Val: node.Identifier}
}

func (gen *Generator) groupingExpr(node parser.Node, scope *Scope) Value {
	src := strings.Builder{}
	src.WriteString("(")
	value := gen.GenerateNode(*node.Expression, scope)
	src.WriteString(value.Val)
	src.WriteString(")")

	return Value{value.Type, src.String()}
}

func (gen *Generator) listExpr(node parser.Node, scope *Scope) Value {
	nodes, _ := node.Value.([]parser.Node)

	src := strings.Builder{}
	src.WriteString("{")
	for i, expr := range nodes {
		src.WriteString(gen.GenerateNode(expr, scope).Val)

		if i != len(nodes)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString("}")

	return Value{parser.List, src.String()}
}

func (gen *Generator) callExpr(node parser.Node, scope *Scope) Value {
	src := strings.Builder{}
	fn := gen.GenerateNode(*node.Expression, scope)
	args := node.Value.([]parser.Node)

	src.WriteString(fn.Val)
	src.WriteString("(")
	for i, arg := range args {
		src.WriteString(gen.GenerateNode(arg, scope).Val)
		if i != len(args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")

	//fnReturn := scope.GetVariable(node.Identifier)

	return Value{parser.Bool, src.String()}
}

func (gen *Generator) mapExpr(node parser.Node, scope *Scope) Value {
	src := strings.Builder{}

	var tabs string
	for i := 0; i < scope.Count + gen.TabsCount; i++ {
		tabs += "\t"
	}

	var mapTabs string
	for i := 0; i < scope.Count + 1 + gen.TabsCount; i++ {
		mapTabs += "\t"
	}

	gen.TabsCount += 1

	src.WriteString("{\n")
	kv := node.Value.(map[string]parser.Node)
	index := 0
	for k, v := range kv {
		val := gen.GenerateNode(v, scope)

		if index != len(kv)-1 {
			src.WriteString(fmt.Sprintf("%s%s = %v,\n", mapTabs, k, val.Val))		
		}else {
			src.WriteString(fmt.Sprintf("%s%s = %v\n", mapTabs, k, val.Val))		
		}
		index++
	}

	src.WriteString(tabs)
	src.WriteString("}")

	gen.TabsCount -= 1

	return Value{parser.Map, src.String()}
}

func (gen *Generator) unaryExpr(node parser.Node, scope *Scope) Value {
	value := gen.GenerateNode(*node.Right, scope)
	src := fmt.Sprintf("%s%s", node.Token.Lexeme, value.Val)

	return Value{Type: value.Type, Val: src}
}

func (gen *Generator) memberExpr(node parser.Node, scope *Scope) Value {
	src := strings.Builder{}

	expr := gen.GenerateNode(*node.Expression, scope)
	prop := gen.GenerateNode(node.Value.(parser.Node), scope)

	src.WriteString(expr.Val)
	
	if prop.Type == parser.String {
		src.WriteString("[")
		src.WriteString(prop.Val)
		src.WriteString("]")
	} else {
		src.WriteString(".")
		src.WriteString(prop.Val)
	}

	return Value{prop.Type, src.String()}
}

func (gen *Generator) GenerateNode(node parser.Node, environment *Scope) Value {
	scope := environment

	switch node.NodeType {
	case parser.LiteralExpr:
		return gen.literalExpr(node)
	case parser.Prog:
		return gen.Generate(*node.Program, scope)
	case parser.VariableDeclarationStmt:
		return gen.variableDeclarationStmt(node, scope)
	case parser.BinaryExpr:
		return gen.binaryExpr(node, scope)
	case parser.Identifier:
		return gen.identifierExpr(node, scope)
	case parser.GroupingExpr:
		return gen.groupingExpr(node, scope)
	case parser.ListExpr:
		return gen.listExpr(node, scope)
	case parser.AssignmentStmt:
		return gen.assignmentStmt(node, scope)
	case parser.UnaryExpr:
		return gen.unaryExpr(node, scope)
	case parser.FunctionDeclarationStmt:
		return gen.functionDeclarationStmt(node, scope)
	case parser.CallExpr:
		return gen.callExpr(node, scope)
	case parser.ReturnStmt:
		return gen.returnStmt(node, scope)
	case parser.MapExpr:
		return gen.mapExpr(node, scope)
	case parser.MemberExpr:
		return gen.memberExpr(node, scope)
	}

	return Value{}
}
