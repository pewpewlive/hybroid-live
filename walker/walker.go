package walker

import (
	"hybroid/ast"
	"hybroid/lexer"
)

type Walker struct {
	nodes    *[]ast.Node
	Errors   []ast.Error
	Warnings []ast.Warning
}

type Global struct {
	Scope        Scope
	foreignTypes map[string]Value
}

func NewGlobal() Global {
	return Global{
		Scope: NewScope(nil, nil, ReturnProhibiting),
	}
}

type ScopeType int

const (
	ReturnAllowing ScopeType = iota
	ReturnProhibiting
)

type Scope struct {
	Global    *Global
	Parent    *Scope
	Type      ScopeType
	Variables map[string]VariableVal
}

func NewScope(global *Global, parent *Scope, typee ScopeType) Scope {
	return Scope{
		Global:    global,
		Parent:    parent,
		Type:      typee,
		Variables: map[string]VariableVal{},
	}
}

func (w *Walker) error(token lexer.Token, msg string) {
	w.Errors = append(w.Errors, ast.Error{Token: token, Message: msg})
}

func (w *Walker) addError(err ast.Error) {
	w.Errors = append(w.Errors, err)
}

func (w *Walker) GetValueFromType(typee TypeVal) Value {
	switch typee.Type {
	case ast.Number:
		return NumberVal{}
	case ast.Fixed:
		return FixedVal{
			SpecificType: ast.FixedPoint,
		}
	case ast.Bool:
		return BoolVal{}
	case ast.List:
		return ListVal{
			ValueType: *typee.WrappedType,
		}
	case ast.Map:
		return MapVal{
			MemberType: *typee.WrappedType,
		}
	case ast.Func:
		return FunctionVal{
			params:    typee.Params,
			returnVal: typee.Returns,
		}
	case ast.Nil:
		return NilVal{}
	case ast.String:
		return StringVal{}
	case ast.Invalid:
		return Invalid{}
	case 0:
		return Unknown{}
	// TODO: handle structs and entities in the future
	default:
		return Invalid{}
	}
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
		return Invalid{}, &ast.Error{Message: "cannot assign to an undeclared variable"}
	}

	variable := scope.Variables[name]
	if variable.IsConst {
		return Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	scope.Variables[name] = variable

	return scope.Variables[name], nil
}

func (s *Scope) AssignVariable(variable VariableVal, value Value) (Value, *ast.Error) {
	if variable.IsConst {
		return Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
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

func (w *Walker) validateArithmeticOperands(left TypeVal, right TypeVal, expr ast.BinaryExpr) bool {
	//fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	switch left.Type {
	case ast.Nil:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on nil value")
		return false
	case ast.Invalid:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	switch right.Type {
	case ast.Nil:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on nil value")
		return false
	case ast.Invalid:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	switch left.Type {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	switch right.Type {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	return true
}

func returnsAreValid(list1 []TypeVal, list2 []TypeVal) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i, v := range list1 {
		if !((list2[i].WrappedType != nil && list2[i].WrappedType.Type == 0) || 
			(v.WrappedType != nil && v.WrappedType.Type == 0)) && 
			!list2[i].Eq(v) {
			//fmt.Printf("%s : %s\n", v.WrappedType.Type.ToString(), list2[i].WrappedType.Type.ToString())

			return false
		}
	}
	return true
}

func (w *Walker) validateReturnValues(node ast.Node, returnValues []TypeVal, expectedReturnValues []TypeVal) {
	if !returnsAreValid(returnValues, expectedReturnValues) {
		w.error(node.GetToken(), "invalid return type(s)")
	}
	if len(returnValues) < len(expectedReturnValues) {
		w.error(node.GetToken(), "not enough return values given")
	} else if len(returnValues) > len(expectedReturnValues) {
		w.error(node.GetToken(), "too many return values given")
	}
}

func (w *Walker) getReturnFromNode(node *ast.Node, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	localScope := NewScope(scope.Global, scope, scope.Type)
	switch (*node).GetType() {
	case ast.IfStatement:
		converted := (*node).(ast.IfStmt)
		return w.ifReturns(&converted, expectedReturn, &localScope)
	case ast.RepeatStatement:
		converted := ((*node).(ast.RepeatStmt)).Body
		return w.bodyReturns(&converted, expectedReturn, &localScope)
	case ast.ReturnStatement:
		converted := (*node).(ast.ReturnStmt)
		return w.returnStmt(&converted, scope)
	default:
		return nil
	}
}

func (w *Walker) ifReturns(node *ast.IfStmt, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	var returns *ReturnType

	localScope := NewScope(scope.Global, scope, scope.Type)
	returns = w.bodyReturns(&node.Body, expectedReturn, &localScope)
	if returns == nil {
		return nil
	}

	for i := range node.Elseifs {
		localScope := NewScope(scope.Global, scope, scope.Type)
		returns = w.bodyReturns(&node.Elseifs[i].Body, expectedReturn, &localScope)
	}
	if returns == nil {
		return nil
	}

	if node.Else != nil {
		localScope := NewScope(scope.Global, scope, scope.Type)
		returns = w.bodyReturns(&node.Else.Body, expectedReturn, &localScope)
	} else {
		return nil
	}

	return returns
}

func (w *Walker) bodyReturns(body *[]ast.Node, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	var returns *ReturnType
	for _, node := range *body {
		returns = w.getReturnFromNode(&node, expectedReturn, scope)
		if returns == nil {
			continue
		}

		w.validateReturnValues(node, returns.values, expectedReturn.values)
	}
	if returns == nil {
		return returns
	}

	return returns
}

func (w *Walker) GetTypeFromString(str string) ast.PrimitiveValueType {
	switch str {
	case "number":
		return ast.Number
	case "fixed":
		return ast.FixedPoint
	case "text":
		return ast.String
	case "map":
		return ast.Map
	case "list":
		return ast.List
	case "fn":
		return ast.Func
	case "bool":
		return ast.Bool
	default:
		return ast.Invalid
	}
}

func (w *Walker) GetDefaultValue(typee TypeVal) string {
	switch typee.Type {
	case ast.Number:
		return "0"
	case ast.Fixed:
		return "0fx"
	case ast.Bool:
		return "false"
	case ast.String:
		return "\"\""
	case ast.List, ast.Map:
		return "{}"
	case ast.Func:
		return ""
	default:
		return "nil"
	}
}

func (w *Walker) Walk(nodes *[]ast.Node, global *Global) []ast.Node {
	w.nodes = nodes

	newNodes := make([]ast.Node, len(*nodes))

	for _, node := range *nodes {
		w.WalkNode(&node, &global.Scope)
		newNodes = append(newNodes, node)
	}

	return newNodes
}

func (w *Walker) WalkNode(node *ast.Node, scope *Scope) {
	switch newNode := (*node).(type) {
	case ast.VariableDeclarationStmt:
		w.variableDeclarationStmt(&newNode, scope)
		*node = newNode
	case ast.IfStmt:
		w.ifStmt(&newNode, scope)
		*node = newNode
	case ast.AssignmentStmt:
		w.assignmentStmt(&newNode, scope)
		*node = newNode
	case ast.FunctionDeclarationStmt:
		w.functionDeclarationStmt(&newNode, scope)
		*node = newNode
	case ast.ReturnStmt:
		w.returnStmt(&newNode, scope)
		*node = newNode
	case ast.RepeatStmt:
		w.repeatStmt(&newNode, scope)
		*node = newNode
	case ast.TickStmt:
		w.tickStmt(&newNode, scope)
		*node = newNode
	case ast.CallExpr:
		w.callExpr(&newNode, scope)
		*node = newNode
	case ast.DirectiveExpr:
		w.directiveExpr(&newNode, scope)
		*node = newNode
	case ast.UseStmt:
		w.useStmt(&newNode, scope)
		*node = newNode
	default:
		w.error(newNode.GetToken(), "Expected statement")
	}
}

func (w *Walker) GetNodeValue(node *ast.Node, scope *Scope) Value {
	switch newNode := (*node).(type) {
	case ast.LiteralExpr:
		return w.literalExpr(&newNode)
	case ast.BinaryExpr:
		return w.binaryExpr(&newNode, scope)
	case ast.IdentifierExpr:
		return w.identifierExpr(&newNode, scope)
	case ast.GroupExpr:
		return w.groupingExpr(&newNode, scope)
	case ast.ListExpr:
		return w.listExpr(&newNode, scope)
	case ast.UnaryExpr:
		return w.unaryExpr(&newNode, scope)
	case ast.CallExpr:
		return w.callExpr(&newNode, scope)
	case ast.MapExpr:
		return w.mapExpr(&newNode, scope)
	case ast.DirectiveExpr:
		return w.directiveExpr(&newNode, scope)
	case ast.AnonFnExpr:
		return w.anonFnExpr(&newNode, scope)
	case ast.MemberExpr:
		return w.memberExpr(nil, &newNode, scope)
	case ast.TypeExpr:
		return w.typeExpr(&newNode)
	default:
		w.error(newNode.GetToken(), "Expected expression")
		return NilVal{}
	}
}
