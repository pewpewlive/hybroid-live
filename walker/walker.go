package walker

import (
	"hybroid/ast"
	"hybroid/lexer"
)

type Walker struct {
	nodes  []ast.Node
	Errors []ast.Error
}

type Global struct {
	Scope        Scope
	foreignTypes map[string]Value
}

func NewGlobal() Global {
	return Global{
		Scope: Scope{
			Global:    nil,
			Parent:    nil,
			Variables: make(map[string]VariableVal),
		},
	}
}

type Scope struct {
	Global    *Global
	Parent    *Scope
	Variables map[string]VariableVal
}

func (w *Walker) error(token lexer.Token, msg string) {
	w.Errors = append(w.Errors, ast.Error{Token: token, Message: msg})
}

func (w *Walker) addError(err ast.Error) {
	w.Errors = append(w.Errors, err)
}

func (w *Walker) GetValue(pvt ast.PrimitiveValueType) Value {
	switch pvt {
	case ast.Number:
		return NumberVal{}
	case ast.FixedPoint:
		return FixedVal{}
	case ast.Bool:
		return BoolVal{}
	case ast.List:
		return ListVal{}
	case ast.Map:
		return MapVal{}
	case ast.Func:
		return CallVal{}
	case ast.Nil:
		return NilVal{}
	case ast.String:
		return StringVal{}
	case ast.Ident:
		return VariableVal{}
	case ast.Undefined:
		return Unknown{}
	case 0:
		return Undefined{}
	// TODO: handle structs and entities in the future
	default:
		return Unknown{}
	}
}

func (s *Scope) GetVariable(name string) VariableVal {

	scope := s.Resolve(name)

	variable := scope.Variables[name]

	variable.IsUsed = true

	scope.Variables[name] = variable

	return scope.Variables[name]
}

func (s *Scope) AssignVariable(name string, value Value) (Value, *ast.Error) {
	scope := s.Resolve(name)

	if scope == nil {
		return Unknown{}, &ast.Error{Message: "cannot assign to an undeclared variable"}
	}

	variable := scope.Variables[name]
	if variable.IsConst {
		return Unknown{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	scope.Variables[name] = variable

	return scope.Variables[name], nil
}

func (s *Scope) AssignVariableType(name string, pvt ast.PrimitiveValueType) {

}

func (s *Scope) DeclareVariable(value VariableVal) (VariableVal, bool) {
	if _, found := s.Variables[value.Name]; found {
		return VariableVal{}, false
	}

	s.Variables[value.Name] = value
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

func (w *Walker) validateArithmeticOperands(left Value, right Value, expr ast.BinaryExpr) bool {
	//fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	switch left.GetType() {
	case ast.Nil:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on nil value")
		return false
	case ast.Undefined:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on undefined value")
		return false
	}

	switch right.GetType() {
	case ast.Nil:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on nil value")
		return false
	case ast.Undefined:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on undefined value")
		return false
	}

	switch left.GetType() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	switch right.GetType() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	return true
}

func (w *Walker) Walk(nodes []ast.Node, global *Global) []ast.Node {
	w.nodes = nodes

	newNodes := make([]ast.Node, len(nodes))

	for _, node := range nodes {
		w.WalkNode(node, &global.Scope)
	}

	return newNodes
}

func (w *Walker) WalkNode(node ast.Node, scope *Scope) {
	switch newNode := node.(type) {
	case ast.VariableDeclarationStmt:
		w.variableDeclarationStmt(newNode, scope)
	case ast.IfStmt:
		w.ifStmt(newNode, scope)
	case ast.AssignmentStmt:
		w.assignmentStmt(newNode, scope)
	case ast.FunctionDeclarationStmt:
		w.functionDeclarationStmt(newNode, scope)
	case ast.ReturnStmt:
		w.returnStmt(newNode, scope)
	case ast.RepeatStmt:
		w.repeatStmt(newNode, scope)
	case ast.TickStmt:
		w.tickStmt(newNode, scope)
	case ast.CallExpr:
		w.callExpr(newNode, scope)
	case ast.DirectiveExpr:
		w.directiveExpr(newNode, scope)
	case ast.UseStmt:
		w.useStmt(newNode, scope)
	default:
		w.error(newNode.GetToken(), "Expected statement")
	}
}

func (w *Walker) GetNodeValue(node ast.Node, scope *Scope) Value {
	switch newNode := node.(type) {
	case ast.LiteralExpr:
		return w.literalExpr(newNode)
	case ast.BinaryExpr:
		return w.binaryExpr(newNode, scope)
	case ast.IdentifierExpr:
		return w.identifierExpr(newNode, scope)
	case ast.GroupExpr:
		return w.groupingExpr(newNode, scope)
	case ast.ListExpr:
		return w.listExpr(newNode, scope)
	case ast.UnaryExpr:
		return w.unaryExpr(newNode, scope)
	case ast.CallExpr:
		return w.callExpr(newNode, scope)
	case ast.MapExpr:
		return w.mapExpr(newNode, scope)
	case ast.DirectiveExpr:
		return w.directiveExpr(newNode, scope)
	case ast.ParentExpr:
		return w.parentExpr(newNode, scope)
	default:
		w.error(newNode.GetToken(), "Expected expression")
		return NilVal{}
	}
}
