package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
)

type Walker struct {
	Environment *Environment
	Walkers     *map[string]*Walker
	nodes       *[]ast.Node
	Errors      []ast.Error
	Warnings    []ast.Warning
	Context     ast.Node
}

func NewWalker(path string) *Walker {
	environment := NewEnvironment(path)
	walker := Walker{
		Environment: &environment,
		nodes:       &[]ast.Node{},
		Errors:      []ast.Error{},
		Warnings:    []ast.Warning{},
	}
	return &walker
}

func (w *Walker) error(token lexer.Token, msg string) {
	w.Errors = append(w.Errors, ast.Error{Token: token, Message: msg})
}

func (w *Walker) warn(token lexer.Token, msg string) {
	w.Warnings = append(w.Warnings, ast.Warning{Token: token, Message: msg})
}

func (w *Walker) addError(err ast.Error) {
	w.Errors = append(w.Errors, err)
}

func (w *Walker) GetValueFromType(typee Type) Value {
	switch typee.Type {
	case ast.Number:
		return &NumberVal{}
	case ast.FixedPoint, ast.Fixed, ast.Radian, ast.Degree:
		return &FixedVal{
			SpecificType: typee.Type,
		}
	case ast.Bool:
		return &BoolVal{}
	case ast.List:
		return &ListVal{
			ValueType: *typee.WrappedType,
		}
	case ast.Map:
		return &MapVal{
			MemberType: *typee.WrappedType,
		}
	case ast.Func:
		return &FunctionVal{
			params:    typee.Params,
			returnVal: typee.Returns,
		}
	case ast.String:
		return &StringVal{}
	case ast.Invalid:
		return &Invalid{}
	case 0:
		return &Unknown{}
	case ast.Struct:
		return w.Environment.Scope.GetStructType(&w.Environment.Scope, typee.Name)
	default:
		return &Invalid{}
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
	structType := scope.Environment.StructTypes[name]

	structType.IsUsed = true

	scope.Environment.StructTypes[name] = structType

	return scope.Environment.StructTypes[name]
}

func (s *Scope) AssignVariableByName(name string, value Value) (Value, *ast.Error) {
	scope := s.ResolveVariable(name)

	if scope == nil {
		return &Invalid{}, &ast.Error{Message: "cannot assign to an undeclared variable"}
	}

	variable := scope.Variables[name]
	if variable.IsConst {
		return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	scope.Variables[name] = variable

	temp := scope.Variables[name]

	return &temp, nil
}

func (s *Scope) AssignVariable(variable VariableVal, value Value) (Value, *ast.Error) {
	if variable.IsConst {
		return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	return &variable, nil
}

func (s *Scope) DeclareVariable(value VariableVal) (VariableVal, bool) {
	if _, found := s.Variables[value.Name]; found {
		return VariableVal{}, false
	}

	s.Variables[value.Name] = value
	s.VariableIndexes[value.Name] = len(s.VariableIndexes) + 1

	return value, true
}

func (s *Scope) DeclareStructType(structType *StructTypeVal) bool {
	if _, found := s.Environment.StructTypes[structType.Name.Lexeme]; found {
		return false
	}

	s.Environment.StructTypes[structType.Name.Lexeme] = structType
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

func (s *Scope) ResolveStructType(name string) *Scope {
	if _, found := s.Environment.StructTypes[name]; found {
		return s
	}

	if s.Parent == nil {
		return nil
	}

	return s.Parent.ResolveStructType(name)
}

func ResolveTagScope[T ScopeTag](sc *Scope) (*Scope, *ScopeTag, *T) {
	if tag, ok := sc.Tag.(T); ok {
		return sc, &sc.Tag, &tag
	}

	if sc.Parent == nil {
		return nil, nil, nil
	}

	return ResolveTagScope[T](sc.Parent)
}

func (sc *Scope) ResolveReturnable() *ExitableTag {
	if sc.Parent == nil {
		return nil
	}

	if returnable := helpers.GetValOfInterface[ExitableTag](sc.Tag); returnable != nil {
		return returnable
	}

	if helpers.IsZero(sc.Tag) {
		return nil
	}

	return sc.Parent.ResolveReturnable()
}

func (g *Environment) GetForeignType(str string) Value {
	return g.foreignTypes[str]
}

func (w *Walker) validateArithmeticOperands(left Type, right Type, expr ast.BinaryExpr) bool {
	//fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	if left.Type == ast.Invalid {
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	if right.Type == ast.Invalid {
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

func returnsAreValid(list1 []Type, list2 []Type) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i, v := range list1 {
		fmt.Printf("%s compared to %s\n", list1[i].ToString(), list2[i].ToString())
		if !((list2[i].WrappedType != nil && list2[i].WrappedType.Type == 0) ||
			(v.WrappedType != nil && v.WrappedType.Type == 0)) &&
			!TypeEquals(&list2[i], &v) {
			return false
		}
	}
	return true
}

func (w *Walker) validateReturnValues(_return Types, expectReturn Types) string {
	returnValues, expectedReturnValues := _return, expectReturn
	if len(returnValues) < len(expectedReturnValues) {
		return "not enough return values given"
	} else if len(returnValues) > len(expectedReturnValues) {
		return "too many return values given"
	}
	if !returnsAreValid(returnValues, expectedReturnValues) {
		return "invalid return type(s)"
	}
	return ""
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

func (w *Walker) Pass1(nodes *[]ast.Node, wlkrs *map[string]*Walker) []ast.Node {
	w.Walkers = wlkrs
	w.nodes = nodes

	newNodes := make([]ast.Node, 0)

	for _, node := range *nodes {
		w.WalkNode(&node, &w.Environment.Scope)
		newNodes = append(newNodes, node)
	}

	return newNodes
}

// func (w *Walker) Pass1(nodes *[]ast.Node) []ast.Node {
// 	w.nodes = nodes

// 	newNodes := make([]ast.Node, 0)

// 	scope := &w.Environment.Scope
// 	for _, node := range *nodes {
// 		switch newNode := node.(type) {
// 		case *ast.EnvironmentStmt:
// 			w.envStmt(newNode, scope)
// 		case *ast.VariableDeclarationStmt:
// 			w.variableDeclarationStmt(newNode, scope)
// 		case *ast.FunctionDeclarationStmt:
// 			w.functionDeclarationStmt(newNode, scope, Function)
// 		case *ast.StructDeclarationStmt:
// 			w.structDeclarationStmt(newNode, scope)
// 		case *ast.Improper:
// 			w.error(newNode.GetToken(), "Improper statement: parser fault")
// 		default:
// 			w.error(newNode.GetToken(), "Expected statement")
// 		}
// 		newNodes = append(newNodes, node)
// 	}

// 	return newNodes
// }

func (w *Walker) WalkNode(node *ast.Node, scope *Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		w.envStmt(newNode, scope)
	case *ast.VariableDeclarationStmt:
		w.variableDeclarationStmt(newNode, scope)
	case *ast.IfStmt:
		w.ifStmt(newNode, scope)
	case *ast.AssignmentStmt:
		w.assignmentStmt(newNode, scope)
	case *ast.FunctionDeclarationStmt:
		w.functionDeclarationStmt(newNode, scope, Function)
	case *ast.ReturnStmt:
		w.returnStmt(newNode, scope)
	case *ast.YieldStmt:
		w.yieldStmt(newNode, scope)
	case *ast.BreakStmt:
		w.breakStmt(newNode, scope)
	case *ast.ContinueStmt:
		w.continueStmt(newNode, scope)
	case *ast.RepeatStmt:
		w.repeatStmt(newNode, scope)
	case *ast.ForStmt:
		w.forStmt(newNode, scope)
	case *ast.TickStmt:
		w.tickStmt(newNode, scope)
	case *ast.CallExpr:
		w.callExpr(newNode, scope, Function)
	case *ast.MethodCallExpr:
		w.methodCallExpr(node, scope)
	case *ast.DirectiveExpr:
		w.directiveExpr(newNode, scope)
	case *ast.UseStmt:
		w.useStmt(newNode, scope)
	case *ast.StructDeclarationStmt:
		w.structDeclarationStmt(newNode, scope)
	case *ast.MatchStmt:
		w.matchStmt(newNode, false, scope)
	case *ast.Improper:
		w.error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.error(newNode.GetToken(), "Expected statement")
	}
}

func (w *Walker) GetNodeValue(node *ast.Node, scope *Scope) Value {
	var val Value

	switch newNode := (*node).(type) {
	case *ast.LiteralExpr:
		val = w.literalExpr(newNode)
	case *ast.BinaryExpr:
		val = w.binaryExpr(newNode, scope)
	case *ast.IdentifierExpr:
		val = w.identifierExpr(node, scope)
	case *ast.GroupExpr:
		val = w.groupingExpr(newNode, scope)
	case *ast.ListExpr:
		val = w.listExpr(newNode, scope)
	case *ast.UnaryExpr:
		val = w.unaryExpr(newNode, scope)
	case *ast.CallExpr:
		val = w.callExpr(newNode, scope, Function)
	case *ast.MapExpr:
		val = w.mapExpr(newNode, scope)
	case *ast.DirectiveExpr:
		val = w.directiveExpr(newNode, scope)
	case *ast.AnonFnExpr:
		val = w.anonFnExpr(newNode, scope)
	case *ast.AnonStructExpr:
		val = w.anonStructExpr(newNode, scope)
	case *ast.MethodCallExpr:
		val = w.methodCallExpr(node, scope)
	case *ast.MemberExpr:
		val = w.memberExpr(newNode, scope)
	case *ast.FieldExpr:
		val = w.fieldExpr(newNode, scope)
	case *ast.NewExpr:
		val = w.newExpr(newNode, scope)
	case *ast.SelfExpr:
		val = w.selfExpr(newNode, scope)
	case *ast.MatchExpr:
		val = w.matchExpr(newNode, scope)
	default:
		w.error(newNode.GetToken(), "Expected expression")
		return &Invalid{}
	}
	return val
}
