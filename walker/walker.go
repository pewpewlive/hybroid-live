package walker

import (
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
)

type Walker struct {
	Environment *EnvironmentVal
	Walkers     *map[string]*Walker
	nodes       *[]ast.Node
	Errors      []ast.Error
	Warnings    []ast.Warning
	Context     Context
}

func NewWalker(path string) *Walker {
	environment := NewEnvironment(path)
	walker := Walker{
		Environment: &environment,
		nodes:       &[]ast.Node{},
		Errors:      []ast.Error{},
		Warnings:    []ast.Warning{},
		Context: 	 Context{
			Node:  &ast.Improper{},
			Value: &Unknown{},
			Ret:   Types{},
		},
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

func (s *Scope) GetVariable( name string) *VariableVal {
	variable := s.Variables[name]

	variable.IsUsed = true

	s.Variables[name] = variable

	return s.Variables[name]
}

// ONLY CALL WHEN 100% SURE YOU'RE GONNA GET A STRUCT BACK
func (w *Walker) GetStruct(name string) (Value, bool) {
	structType, found := w.Environment.Structs[name]
	if !found {
		//w.error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s", name, " exists"))
		return nil, false
	}

	structType.Type.IsUsed = true

	w.Environment.Structs[name] = structType

	return structType, true
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

	return temp, nil
}

func (s *Scope) AssignVariable(variable *VariableVal, value Value) (Value, *ast.Error) {
	if variable.IsConst {
		return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	return variable, nil
}

func (s *Scope) DeclareVariable(value *VariableVal) (*VariableVal, bool) {
	if varFound, found := s.Variables[value.Name]; found {
		return varFound, false
	}

	s.Variables[value.Name] = value
	return value, true
}

func (w *Walker) DeclareStruct(structVal *StructVal) bool {
	if _, found := w.Environment.Structs[structVal.Type.Name]; found {
		return false
	}

	w.Environment.Structs[structVal.Type.Name] = structVal
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

func (w *Walker) validateArithmeticOperands(left Type, right Type, expr ast.BinaryExpr) bool {
	//fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	if left.PVT() == ast.Invalid {
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	if right.PVT() == ast.Invalid {
		w.error(expr.Right.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	switch left.PVT() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.error(expr.Left.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	switch right.PVT() {
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
		//fmt.Printf("%s compared to %s\n", list1[i].ToString(), list2[i].ToString())
		if !TypeEquals(v, list2[i]) {
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

func (w *Walker) TypeToValue(_type Type) Value {
	switch _type.PVT() {
	case ast.Radian, ast.Fixed, ast.FixedPoint, ast.Degree:
		return &FixedVal{SpecificType: _type.PVT()}
	case ast.Bool:
		return &BoolVal{}
	case ast.String:
		return &StringVal{}
	case ast.Number:
		return &NumberVal{}
	case ast.List:
		return &ListVal{
			ValueType: _type.(*WrapperType).WrappedType,
		}
	case ast.Map:
		return &MapVal{
			MemberType: _type.(*WrapperType).WrappedType,
		}
	case ast.Struct:
		val, _ := w.GetStruct(_type.ToString())
		return val
	case ast.AnonStruct:
		return &AnonStructVal{
			Fields: _type.(*AnonStructType).Fields,
		}
	case ast.Environment:
		return (*w.Walkers)[_type.(*EnvironmentType).Name].Environment
	case ast.Enum:
		return w.Environment.Scope.GetVariable(_type.(*EnumType).Name)
	default:
		return &Invalid{}
	}
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
	case "struct":
		return ast.AnonStruct
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

func (w *Walker) ReportExits(sender ExitableTag, scope *Scope) {
	receiver_ := scope.ResolveReturnable()

	if receiver_ == nil {
		return
	}

	receiver := *receiver_

	receiver.SetExit(sender.GetIfExits(Yield), Yield)
	receiver.SetExit(sender.GetIfExits(Return), Return)
	receiver.SetExit(sender.GetIfExits(Break), Break)
	receiver.SetExit(sender.GetIfExits(Continue), Continue)
	receiver.SetExit(sender.GetIfExits(Yield), All)
} 

func (w *Walker) WalkBody(body *[]ast.Node, tag ExitableTag, scope *Scope) {
	endIndex := -1
	for i := range *body {
		if tag.GetIfExits(All) {
			w.warn((*body)[i].GetToken(), "unreachable code detected")
			endIndex = i
			break
		}
		w.WalkNode(&(*body)[i], scope)
	}
	if endIndex != -1 {
		*body = (*body)[:endIndex]
	}
}

func (w *Walker) WalkNode(node *ast.Node, scope *Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		w.env(newNode, scope)
	case *ast.VariableDeclarationStmt:
		w.variableDeclaration(newNode, scope)
	case *ast.IfStmt:
		w.ifStmt(newNode, scope)
	case *ast.AssignmentStmt:
		w.assignment(newNode, scope)
	case *ast.FunctionDeclarationStmt:
		w.functionDeclaration(newNode, scope, Function)
	case *ast.ReturnStmt:
		w.returnStmt(newNode, scope)
	case *ast.YieldStmt:
		w.yieldStmt(newNode, scope)
	case *ast.BreakStmt:
		w.breakStmt(newNode, scope)
	case *ast.ContinueStmt:
		w.continueStmt(newNode, scope)
	case *ast.RepeatStmt:
		w.repeat(newNode, scope)
	case *ast.WhileStmt:
		w.while(newNode, scope)
	case *ast.ForStmt:
		w.forloop(newNode, scope)
	case *ast.TickStmt:
		w.tick(newNode, scope)
	case *ast.CallExpr:
		w.callExpr(newNode, scope, Function)
	case *ast.MethodCallExpr:
		w.methodCallExpr(node, scope)
	case *ast.DirectiveExpr:
		w.directiveExpr(newNode, scope)
	case *ast.UseStmt:
		w.use(newNode, scope)
	case *ast.EnumDeclarationStmt:
		w.enumDeclarationStmt(newNode, scope)
	case *ast.StructDeclarationStmt:
		w.structDeclaration(newNode, scope)
	case *ast.MatchStmt:
		w.match(newNode, false, scope)
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
