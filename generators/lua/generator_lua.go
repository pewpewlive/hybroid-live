package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
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
	Errors    []GenError
	Src       strings.Builder
	TabsCount int
}

type Value struct {
	properties *[]Value
	Type       ast.PrimitiveValueType
	Token      lexer.Token
	Val        string
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

func (gen *Generator) validateOperands(left Value, right Value) bool {
	fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	if left.Type == ast.Nil {
		gen.error(left.Token, "cannot perform arithmetic on nil value")
		return false
	} else if right.Type == ast.Nil {
		gen.error(right.Token, "cannot perform arithmetic on nil value")
		return false
	} else if left.Type == ast.Undefined {
		gen.error(left.Token, "cannot perform arithmetic on undefined value")
		return false
	} else if right.Type == ast.Undefined {
		gen.error(right.Token, "cannot perform arithmetic on undefined value")
		return false
	} else {
		if left.Type == ast.List || left.Type == ast.Map || left.Type == ast.String || left.Type == ast.Bool || left.Type == ast.Entity || left.Type == ast.Struct {
			gen.error(left.Token, "cannot perform arithmetic on a non-number value")
			return false
		} else if right.Type == ast.List || right.Type == ast.Map || right.Type == ast.String || right.Type == ast.Bool || right.Type == ast.Entity || right.Type == ast.Struct {
			gen.error(right.Token, "cannot perform arithmetic on a non-number value")
			return false
		}
	}
	return true
}

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

func (gen *Generator) Generate(program []ast.Node, environment *Scope) Value {
	var lastEvaluated Value

	for _, node := range program {
		lastEvaluated = gen.GenerateNode(node, environment)
		gen.append(lastEvaluated.Val, "\n")
	}

	return lastEvaluated
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

func (gen *Generator) GenerateNode(node ast.Node, environment *Scope) Value {
	scope := environment

	switch newNode := node.(type) {
	case ast.VariableDeclarationStmt:
		return gen.variableDeclarationStmt(newNode, scope)
	case ast.IfStmt:
		return gen.ifStmt(newNode, scope)
	case ast.AssignmentStmt:
		return gen.assignmentStmt(newNode, scope)
	case ast.FunctionDeclarationStmt:
		return gen.functionDeclarationStmt(newNode, scope)
	case ast.ReturnStmt:
		return gen.returnStmt(newNode, scope)
	case ast.LiteralExpr:
		return gen.literalExpr(newNode)
	case ast.BinaryExpr:
		return gen.binaryExpr(newNode, scope)
	case ast.IdentifierExpr:
		return gen.identifierExpr(newNode, scope)
	case ast.GroupExpr:
		return gen.groupingExpr(newNode, scope)
	case ast.ListExpr:
		return gen.listExpr(newNode, scope)
	case ast.UnaryExpr:
		return gen.unaryExpr(newNode, scope)
	case ast.CallExpr:
		return gen.callExpr(newNode, scope)
	case ast.MapExpr:
		return gen.mapExpr(newNode, scope)
	case ast.MemberExpr:
		return gen.memberExpr(newNode, scope)
	}

	return Value{}
}
