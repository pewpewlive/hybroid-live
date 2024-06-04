package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
)

type Walker struct {
	Namespace *Namespace
	nodes     *[]ast.Node
	Errors    []ast.Error
	Warnings  []ast.Warning
	Context   ast.Node
}

func NewWalker() *Walker {
	namespace := NewNamespace()
	walker := Walker{
		Namespace: &namespace,
		nodes:     &[]ast.Node{},
		Errors:    []ast.Error{},
		Warnings:  []ast.Warning{},
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
			params:    typee.Params,
			returnVal: typee.Returns,
		}
	case ast.String:
		return StringVal{}
	case ast.Invalid:
		return Invalid{}
	case 0:
		return Unknown{}
	case ast.Struct:
		return w.Namespace.Scope.GetStructType(&w.Namespace.Scope, typee.Name)
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
	structType := scope.Namespace.StructTypes[name]

	structType.IsUsed = true

	scope.Namespace.StructTypes[name] = structType

	return scope.Namespace.StructTypes[name]
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
	s.VariableIndexes[value.Name] = len(s.VariableIndexes) + 1

	return value, true
}

func (s *Scope) DeclareStructType(structType *StructTypeVal) bool {
	if _, found := s.Namespace.StructTypes[structType.Name.Lexeme]; found {
		return false
	}

	s.Namespace.StructTypes[structType.Name.Lexeme] = structType
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
	if _, found := s.Namespace.StructTypes[name]; found {
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

func (sc *Scope) ResolveReturnable() (*Scope, *ReturnableTag) {
	if sc.Parent == nil {
		return nil, nil
	}

	if returnable := helpers.GetValOfInterface[ReturnableTag](sc.Tag); returnable != nil {
		return sc, returnable
	}

	if helpers.IsZero(sc.Tag) {
		return nil, nil
	}

	return sc.Parent.ResolveReturnable()
}

func (g *Namespace) GetForeignType(str string) Value {
	return g.foreignTypes[str]
}

func (w *Walker) validateArithmeticOperands(left TypeVal, right TypeVal, expr ast.BinaryExpr) bool {
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

func returnsAreValid(list1 []TypeVal, list2 []TypeVal) bool {
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

func (w *Walker) validateReturnValues(_return Returns, expectReturn Returns) string {
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

func (w *Walker) Stage1(nodes *[]ast.Node) []ast.Node {
	w.nodes = nodes

	newNodes := make([]ast.Node, 0)

	for _, node := range *nodes {
		w.WalkNode(&node, &w.Namespace.Scope)
		newNodes = append(newNodes, node)
	}

	return newNodes
}

func (w *Walker) Stage2(nodes *[]ast.Node, namespaces *map[string]*Namespace) []ast.Node {
	return *nodes
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
	case ast.YieldStmt:
		w.yieldStmt(&newNode, scope)
		*node = newNode
	case ast.BreakStmt:
		w.breakStmt(&newNode, scope)
		*node = newNode
	case ast.ContinueStmt:
		w.continueStmt(&newNode, scope)
		*node = newNode
	case ast.RepeatStmt:
		w.repeatStmt(&newNode, scope)
		*node = newNode
	case ast.ForStmt:
		w.forStmt(&newNode, scope)
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
	case ast.Improper:
		w.error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.error(newNode.GetToken(), "Expected statement")
	}
}

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
		val = w.memberExpr(&newNode, scope)
		*node = newNode
	case ast.FieldExpr:
		val = w.fieldExpr(&newNode, scope)
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
	case ast.MatchExpr:
		val = w.matchExpr(&newNode, scope)
		*node = newNode
	default:
		w.error(newNode.GetToken(), "Expected expression")
		return Invalid{}
	}
	return val
}
