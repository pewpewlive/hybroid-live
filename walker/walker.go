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
		return FunctionVal{}
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

type ComparableType interface {
	ast.PrimitiveValueType
}

func listsAreValid[T ComparableType](list1 []T, list2 []T) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i, v := range list1 {
		if list2[i] != v {
			return false
		}
	}

	return true
}

func (s *Scope) GetVariable(name string) VariableVal {

	scope := s.Resolve(name)

	variable := scope.Variables[name]

	variable.IsUsed = true

	scope.Variables[name] = variable

	return scope.Variables[name]
}

func (s *Scope) AssignVariableByName(name string, value Value) (Value, *ast.Error) {
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

func (s *Scope) AssignVariable(variable VariableVal, value Value) (Value, *ast.Error) {
	if variable.IsConst {
		return Unknown{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	return variable, nil
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

func (w *Walker) GetTypeFromString(str string) ast.PrimitiveValueType {
	switch str {
	case "map":
		return ast.Map
	case "list":
		return ast.List
	case "number":
		return ast.Number
	case "Bool":
		return ast.Bool
	case "fixed":
		return ast.FixedPoint
	case "text":
		return ast.String
	case "fn":
		return ast.Func
	default:
		return ast.Undefined
	}
}

func (w *Walker) validateReturnValues(node ast.Node, returnValues []ast.PrimitiveValueType, expectedReturnValues []ast.PrimitiveValueType) {
	if !listsAreValid(returnValues, expectedReturnValues) {
		w.error(node.GetToken(), "invalid return type(s)")
	}
	if len(returnValues) < len(expectedReturnValues) {
		w.error(node.GetToken(), "not enough return values given")
	} else if len(returnValues) > len(expectedReturnValues) {
		w.error(node.GetToken(), "too many return values given")
	}
}

func (w *Walker) getReturnFromNode(node ast.Node, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	localScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
	switch node.GetType() {
	case ast.IfStatement:
		return w.ifReturns(node.(ast.IfStmt), expectedReturn, &localScope)
	case ast.RepeatStatement:
		return w.bodyReturns(node.(ast.RepeatStmt).Body, expectedReturn, &localScope)
	case ast.ReturnStatement:
		return w.returnStmt(node.(ast.ReturnStmt), scope)
	default:
		return nil
	}
}

func (w *Walker) ifReturns(node ast.IfStmt, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	var returns *ReturnType

	for _, bodynode := range node.Body {
		returns = w.getReturnFromNode(bodynode, expectedReturn, scope)
		if returns == nil {
			continue
		} else if bodynode.GetToken() != node.Body[len(node.Body)-1].GetToken() {
			w.error(bodynode.GetToken(), "unreachable code detected")
		}
		
		w.validateReturnValues(node, returns.values, expectedReturn.values)
	}
	if returns == nil {
		return returns
	}

	for _, elseif := range node.Elseifs {
		for _, node := range elseif.Body {
			returns = w.getReturnFromNode(node, expectedReturn, scope)
			if returns == nil {
				continue
			} else if node.GetToken() != elseif.Body[len(elseif.Body)-1].GetToken() {
				w.error(node.GetToken(), "unreachable code detected")
			}

			w.validateReturnValues(node, returns.values, expectedReturn.values)
		}
		if returns == nil {
			return returns
		}
	}

	if node.Else != nil {
		localScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		returns = w.bodyReturns(node.Else.Body, expectedReturn, &localScope)
	}else {
		return nil
	}

	return returns
}

func (w *Walker) bodyReturns(body []ast.Node, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	var returns *ReturnType
	for _, node := range body {
		returns = w.getReturnFromNode(node, expectedReturn, scope)
		if returns == nil {
			continue
		} else if node.GetToken() != body[len(body)-1].GetToken() {
			w.error(node.GetToken(), "unreachable code detected")
		}

		w.validateReturnValues(node, returns.values, expectedReturn.values)
	}
	if returns == nil {
		return returns
	}

	return returns
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
	case ast.MemberExpr:
		return w.memberExpr(nil, newNode, scope)
	default:
		w.error(newNode.GetToken(), "Expected expression")
		return NilVal{}
	}
}
