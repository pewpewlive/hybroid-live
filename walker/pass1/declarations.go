package pass1

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	. "hybroid/walker"
)

func AnonFnExprPass1(w *Walker, fn *ast.AnonFnExpr, scope *Scope) *FunctionVal {
	ret := EmptyReturn
	for _, typee := range fn.Return {
		ret = append(ret, w.TypeExpr(typee))
	}

	funcTag := &FuncTag{ReturnType: ret}
	fnScope := NewScope(scope, funcTag)
	fnScope.Attributes.Add(ReturnAllowing)

	params := make([]Type, 0)
	for i, param := range fn.Params {
		params = append(params, w.TypeExpr(param.Type))
		value := w.TypeToValue(params[i])
		w.DeclareVariable(&fnScope, &VariableVal{Name: param.Name.Lexeme, Value: value, Token: param.Name}, param.Name)
	}

	WalkBodyPass1(w, &fn.Body, funcTag, &fnScope)

	if !funcTag.GetIfExits(Return) && !ret.Eq(&EmptyReturn) {
		w.Error(fn.GetToken(), "not all code paths return a value")
	}

	return &FunctionVal{
		Params:  params,
		Returns: ret,
	}
}

func MatchExprPass1(w *Walker, node *ast.MatchExpr, scope *Scope) Value {
	casesLength := len(node.MatchStmt.Cases)+1
	if node.MatchStmt.HasDefault {
		casesLength--
	}
	matchScope := NewScope(scope, &MatchExprTag{})
	matchScope.Attributes.Add(YieldAllowing)
	mtt := &MatchExprTag{Mpt:NewMultiPathTag(casesLength, matchScope.Attributes...)}
	matchScope.Tag = mtt

	w.Match(&node.MatchStmt, true, scope)
	for i := range node.MatchStmt.Cases {
		caseScope := NewScope(&matchScope, &UntaggedTag{})
		WalkBodyPass1(w, &node.MatchStmt.Cases[i].Body, mtt, &caseScope)
	}

	return mtt.YieldValues
}

func GetNodeValuePass1(w *Walker, node *ast.Node, scope *Scope) Value {
	var val Value

	switch newNode := (*node).(type) {
	case *ast.LiteralExpr:
		val = w.LiteralExpr(newNode)
	case *ast.BinaryExpr:
		val = w.BinaryExpr(newNode, scope)
	case *ast.IdentifierExpr:
		val = w.IdentifierExpr(node, scope)
	case *ast.GroupExpr:
		val = w.GroupingExpr(newNode, scope)
	case *ast.ListExpr:
		val = w.ListExpr(newNode, scope)
	case *ast.UnaryExpr:
		val = w.UnaryExpr(newNode, scope)
	case *ast.CallExpr:
		val = w.CallExpr(newNode, scope, Function)
	case *ast.MapExpr:
		val = w.MapExpr(newNode, scope)
	case *ast.DirectiveExpr:
		val = w.DirectiveExpr(newNode, scope)
	case *ast.AnonFnExpr:
		val = AnonFnExprPass1(w, newNode, scope)
	case *ast.AnonStructExpr:
		val = w.AnonStructExpr(newNode, scope)
	case *ast.MethodCallExpr:
		val = w.MethodCallExpr(node, scope)
	case *ast.MemberExpr:
		val = w.MemberExpr(newNode, scope)
	case *ast.FieldExpr:
		val = w.FieldExpr(newNode, scope)
	case *ast.NewExpr:
		val = w.NewExpr(newNode, scope)
	case *ast.SelfExpr:
		val = w.SelfExpr(newNode, scope)
	case *ast.MatchExpr:
		val = MatchExprPass1(w, newNode, scope)
	default:
		w.Error(newNode.GetToken(), "Expected expression")
		return &Invalid{}
	}
	return val
}

func WalkBodyPass1(w *Walker, body *[]ast.Node, tag ExitableTag, scope *Scope) {
	endIndex := -1
	for i := range *body {
		if tag.GetIfExits(All) {
			w.Warn((*body)[i].GetToken(), "unreachable code detected")
			endIndex = i
			break
		}
		WalkNodePass1(w, &(*body)[i], scope)
	}
	if endIndex != -1 {
		*body = (*body)[:endIndex]
	}
}

func WalkNodePass1(w *Walker, node *ast.Node, scope *Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		w.EnvStmt(newNode, scope)
	case *ast.VariableDeclarationStmt:
		variableDeclarationPass1(w, newNode, scope)
	case *ast.FunctionDeclarationStmt:
		w.FunctionDeclaration(newNode, scope, Function)
	case *ast.StructDeclarationStmt:
		w.StructDeclaration(newNode, scope)
	case *ast.EnumDeclarationStmt:
		w.EnumDeclarationStmt(newNode, scope)
	case *ast.Improper:
		w.Error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.Error(newNode.GetToken(), "Expected statement")
	}
}

func variableDeclarationPass1(w *Walker, declaration *ast.VariableDeclarationStmt, scope *Scope) []*VariableVal {
	declaredVariables := []*VariableVal{}

	idents := len(declaration.Identifiers)
	values := make([]Value, idents)

	for i := range values {
		values[i] = &Invalid{}
	}

	valuesLength := len(declaration.Values)
	if valuesLength > idents {
		w.Error(declaration.Token, "too many values provided in declaration")
		return declaredVariables
	}

	for i := range declaration.Values {

		exprValue := GetNodeValuePass1(w, &declaration.Values[i], scope)
		if declaration.Values[i].GetType() == ast.SelfExpression {
			w.Error(declaration.Values[i].GetToken(), "cannot assign self to a variable")
		}
		if types, ok := exprValue.(*Types); ok { 
			temp := values[i:]
			values = values[:i]
			w.AddTypesToValues(&values, types)
			values = append(values, temp...)
		} else {
			values[i] = exprValue
		}
	}

	if !declaration.IsLocal {
		w.Error(declaration.Token, "cannot declare a global variable inside a local block")
	}
	if declaration.Token.Type == lexer.Const && scope.Parent != nil {
		w.Error(declaration.Token, "cannot declare a global constant inside a local block")
	}

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}

		valType := values[i].GetType()

		if declaration.Types[i] != nil {
			explicitType := w.TypeExpr(declaration.Types[i])
			if valType == InvalidType && explicitType != InvalidType {
				values[i] = w.TypeToValue(explicitType)
				declaration.Values = append(declaration.Values, values[i].GetDefault())
			} else if !TypeEquals(valType, explicitType) {
				w.Error(declaration.Token, fmt.Sprintf("mismatched types: value type (%s) not the same with explict type (%s)",
					valType.ToString(),
					explicitType.ToString()))
			}
		}

		variable := &VariableVal{
			Value: values[i],
			Name:  ident.Lexeme,
			Token:  ident,
		}
		declaredVariables = append(declaredVariables, variable)
		w.DeclareVariable(scope, variable, lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location})
	}

	return declaredVariables
}