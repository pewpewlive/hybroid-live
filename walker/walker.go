package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
)

type Walker struct {
	Global   *Global
	nodes    *[]ast.Node
	Errors   []ast.Error
	Warnings []ast.Warning
	Context  ast.Node
}

// type Context struct {
// 	Node  *ast.Node
// 	Value *Value
// 	Ret   *ReturnType
// } 
// we may or may not need it
// we should use it if absolutely necessary
// or if there is no other way
type Global struct {
	//Ctx          Context
	Scope        Scope
	foreignTypes map[string]Value
	StructTypes  map[string]*StructTypeVal
	//EntityTypes map[string]EntityTypeVal
}

func NewGlobal() Global {
	scope := NewScope(nil, nil, ReturnProhibiting)
	global := Global{
		Scope:        scope,
		foreignTypes: map[string]Value{},
		StructTypes:  map[string]*StructTypeVal{},
	}

	global.Scope.Global = &global
	return global
}

type ScopeType int

const (
	ReturnAllowing ScopeType = iota
	ReturnProhibiting
	Structure
	Entity
)

type Scope struct {
	Global               *Global
	Parent               *Scope
	Type                 ScopeType
	Variables            map[string]VariableVal
	VariableIndexes      map[string]int
	currentVariableIndex int
	WrappedType          TypeVal
}

func NewScope(global *Global, parent *Scope, typee ScopeType) Scope {
	return Scope{
		Global:               global,
		Parent:               parent,
		Type:                 typee,
		Variables:            map[string]VariableVal{},
		VariableIndexes:      map[string]int{},
		currentVariableIndex: 1,
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
	case ast.FixedPoint, ast.Fixed, ast.Radian, ast.Degree:
		return FixedVal{
			SpecificType: typee.Type,
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
			params:    *typee.Params,
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
	case ast.Struct:
		return w.Global.Scope.GetStructType(&w.Global.Scope, typee.Name)
	default:
		return Invalid{}
	}
}

func (s *Scope) GetVariable(scope *Scope, name string) VariableVal {
	variable := scope.Variables[name]

	variable.IsUsed = true

	scope.Variables[name] = variable

	return scope.Variables[name]
}

func (s *Scope) GetVariableIndex(scope *Scope, name string) int {
	//variable := scope.Variables[name]

	//variable.IsUsed = true

	//scope.Variables[name] = variable

	return scope.VariableIndexes[name]
}

func (s *Scope) GetStructType(scope *Scope, name string) *StructTypeVal {
	structType := scope.Global.StructTypes[name]

	structType.IsUsed = true

	scope.Global.StructTypes[name] = structType

	return scope.Global.StructTypes[name]
}

func (s *Scope) AssignVariableByName(name string, value Value) (Value, *ast.Error) {
	scope := s.ResolveVariable(name)

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
	s.VariableIndexes[value.Name] = s.currentVariableIndex
	s.currentVariableIndex++

	return value, true
}

func (s *Scope) DeclareStructType(structType *StructTypeVal) bool {
	if _, found := s.Global.StructTypes[structType.Name.Lexeme]; found {
		return false
	}

	s.Global.StructTypes[structType.Name.Lexeme] = structType
	return true
}

func (s *Scope) ResolveVariable(name string) *Scope {
	if _, found := s.Variables[name]; found {
		return s
	}

	if s.Parent == nil {
		return nil
	}

	return s.Parent.ResolveVariable(name)
}

func (s *Scope) ResolveStructType(name string) *Scope { // for new expression, i.e new Rectangle
	if _, found := s.Global.StructTypes[name]; found { //yes
		return s
	}

	if s.Parent == nil {
		return nil
	}

	return s.Parent.ResolveStructType(name)
}

func (s *Scope) ResolveStructScope() *Scope { // for self.call()
	if s.Type == Structure {
		return s
	}

	if s.Parent == nil {
		return nil
	}

	return s.Parent.ResolveStructScope()
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
		fmt.Printf("%s compared to %s\n", list1[i].ToString(), list2[i].ToString())
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
		converted := (*node).(ast.RepeatStmt).Body
		return w.bodyReturns(&converted, expectedReturn, &localScope)
	case ast.MatchStatement:
		converted := (*node).(ast.MatchStmt)
		var returns *ReturnType
		for i := range converted.Cases {
			returns = w.bodyReturns(&converted.Cases[i].Body, expectedReturn, scope)
			if returns == nil {
				return nil
			}
		}
		return returns
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
		if expectedReturn == nil {
			*expectedReturn = *returns
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

func (w *Walker) Walk(nodes *[]ast.Node, global *Global) []ast.Node {
	w.nodes = nodes

	newNodes := make([]ast.Node, 0)

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
		w.functionDeclarationStmt(&newNode, scope, Function)
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
		w.callExpr(&newNode, scope, Function)
		*node = newNode
	case ast.MethodCallExpr:
		w.methodCallExpr(node, scope)
	case ast.DirectiveExpr:
		w.directiveExpr(&newNode, scope)
		*node = newNode
	case ast.UseStmt:
		w.useStmt(&newNode, scope)
		*node = newNode
	case ast.StructDeclarationStmt:
		w.structDeclarationStmt(&newNode, scope)
	case ast.MatchStmt:
		w.matchStmt(&newNode, false, scope)
		*node = newNode
	case ast.MatchExpr:
		w.matchExpr(&newNode, scope) // u can start making the generation of the match expr ima eat for a bit
	case ast.Improper:
		w.error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.error(newNode.GetToken(), "Expected statement")
	}
}

var owner Value

func (w *Walker) GetNodeValue(node *ast.Node, scope *Scope) Value {
	var val Value

	switch newNode := (*node).(type) {
	case ast.LiteralExpr:
		val = w.literalExpr(&newNode)
		*node = newNode
	case ast.BinaryExpr:
		val = w.binaryExpr(&newNode, scope)
		*node = newNode
	case ast.IdentifierExpr:
		val = w.identifierExpr(node, scope)
	case ast.GroupExpr:
		val = w.groupingExpr(&newNode, scope)
		*node = newNode
	case ast.ListExpr:
		val = w.listExpr(&newNode, scope)
		*node = newNode
	case ast.UnaryExpr:
		val = w.unaryExpr(&newNode, scope)
		*node = newNode
	case ast.CallExpr:
		val = w.callExpr(&newNode, scope, Function)
		*node = newNode
	case ast.MapExpr:
		val = w.mapExpr(&newNode, scope)
		*node = newNode
	case ast.DirectiveExpr:
		val = w.directiveExpr(&newNode, scope)
		*node = newNode
	case ast.AnonFnExpr:
		val = w.anonFnExpr(&newNode, scope)
		*node = newNode
	case ast.MethodCallExpr:
		val = w.methodCallExpr(node, scope)
	case ast.MemberExpr:
		val = w.memberExpr(owner, &newNode, scope)
		*node = newNode
	case ast.FieldExpr:
		val = w.fieldExpr(owner, &newNode, scope)
		*node = newNode
	case ast.TypeExpr:
		val = w.typeExpr(&newNode)
		*node = newNode
	case ast.NewExpr:
		val = w.newExpr(&newNode, scope)
		*node = newNode
	case ast.SelfExpr:
		val = w.selfExpr(&newNode, scope)
		*node = newNode
	default:
		w.error(newNode.GetToken(), "Expected expression")
		return NilVal{}
	}

	return val
}
