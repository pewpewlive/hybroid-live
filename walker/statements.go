package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
	"strings"
)

func (w *Walker) ifStmt(node *ast.IfStmt, scope *Scope) {
	multiPathScope := NewScope(scope, &MultiPathTag{})
	ifScope := NewScope(&multiPathScope, &UntaggedTag{})
	boolExpr := w.GetNodeValue(&node.BoolExpr, scope)
	if boolExpr.GetType().Type != ast.Bool {
		w.error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	for i := range node.Body {
		w.Context = node
		w.WalkNode(&node.Body[i], &ifScope)
	}

	for i := range node.Elseifs {
		boolExpr := w.GetNodeValue(&node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().Type != ast.Bool {
			w.error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := NewScope(&multiPathScope, &UntaggedTag{})
		for j := range node.Elseifs[i].Body {
			w.WalkNode(&node.Elseifs[i].Body[j], &ifScope)
		}
	}

	has_else := false
	if node.Else != nil {
		elseScope := NewScope(&multiPathScope, &UntaggedTag{})
		for i := range node.Else.Body {
			w.WalkNode(&node.Else.Body[i], &elseScope) //lol
		}
		has_else = true
	}

	mpTag := multiPathScope.Tag.(*MultiPathTag)

	returnabl := scope.ResolveReturnable()

	if returnabl == nil {
		return
	}
	returnable := *returnabl
	if has_else {
		returnable.SetExit(mpTag.GetIfExits(Return), Return)
		returnable.SetExit(mpTag.GetIfExits(Yield), Yield)
		returnable.SetExit(mpTag.GetIfExits(Break), Break)
		returnable.SetExit(mpTag.GetIfExits(Continue), Continue)

		// if len(mpTag.Returns)+len(mpTag.Breaks)+len(mpTag.Continues)+len(mpTag.Yields) > 1 {
		// 	returnable.SetExit(false, Return)
		// }
	}else {
		returnable.SetExit(false, Return)
		returnable.SetExit(false, Yield)
		returnable.SetExit(false, Break)
		returnable.SetExit(false, Return)
	}
}

func SetReturnIfTrue(exits []bool, length int, returnable ExitableTag, gt ExitType) {
	if len(exits) == length {
		returnable.SetExit(true, gt)
	}
}

func (w *Walker) assignmentStmt(assignStmt *ast.AssignmentStmt, scope *Scope) {
	hasFuncs := false

	wIdents := []Value{}
	for i := range assignStmt.Identifiers {
		wIdents = append(wIdents, w.GetNodeValue(&assignStmt.Identifiers[i], scope))
	}

	for i := range assignStmt.Values {
		if assignStmt.Values[i].GetType() == ast.CallExpression {
			hasFuncs = true
		}
		value := w.GetNodeValue(&assignStmt.Values[i], scope)
		if i > len(wIdents)-1 {
			break
		}
		variableType := wIdents[i].GetType()
		valueType := value.GetType()
		if variableType.Type == ast.Invalid {
			w.error(assignStmt.Identifiers[i].GetToken(), "cannot assign a value to an undeclared variable")
			continue
		}

		if !TypeEquals(&variableType, &valueType) {
			w.error(assignStmt.Values[i].GetToken(), fmt.Sprintf("mismatched types: variable has a type of %s, but a value of %s was given to it.", variableType.ToString(), valueType.ToString()))
		}

		variable, ok := wIdents[i].(*VariableVal)

		if ok {
			if _, err := scope.AssignVariable(*variable, value); err != nil {
				err.Token = variable.Node.GetToken()
				w.addError(*err)
			}
		}
	}

	if hasFuncs {
		w.error(assignStmt.GetToken(), "cannot have a function call in assignment")
	} else if len(assignStmt.Values) < len(assignStmt.Identifiers) {
		w.error(assignStmt.Values[len(assignStmt.Values)-1].GetToken(), "not enough values provided in assignment")
	} else if len(assignStmt.Values) > len(assignStmt.Identifiers) {
		w.error(assignStmt.Values[len(assignStmt.Values)-1].GetToken(), "too many values provided in assignment")
	}
}

func (w *Walker) functionDeclarationStmt(node *ast.FunctionDeclarationStmt, scope *Scope, procType ProcedureType) VariableVal {
	ret := EmptyReturn
	for _, typee := range node.Return {
		ret = append(ret, w.typeExpr(typee))
		//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
	}

	fnScope := NewScope(scope, &FuncTag{ReturnType: ret})
	fnScope.Attributes.Add(ReturnAllowing)

	params := make([]Type, 0)
	for i, param := range node.Params {
		params = append(params, w.typeExpr(param.Type))
		value := w.GetValueFromType(params[i])
		fnScope.DeclareVariable(VariableVal{Name: param.Name.Lexeme, Value: value, Node: node})
	}

	variable := VariableVal{
		Name:  node.Name.Lexeme,
		Value: &FunctionVal{params: params, returnVal: ret},
		Node:  node,
	}
	if procType == Function {
		if _, success := scope.DeclareVariable(variable); !success {
			w.error(node.Name, fmt.Sprintf("variable with name '%s' already exists", variable.Name))
		}
	}

	if scope.Parent != nil && !node.IsLocal {
		w.error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	endIndex := -1
	for i := range node.Body {
		if funcTag, ok := fnScope.Tag.(*FuncTag); ok {
			if funcTag.GetIfExits(Return) {
				w.warn(node.Body[i].GetToken(), "unreachable code detected")
				endIndex = i
				break
			}
		}
		w.WalkNode(&node.Body[i], &fnScope)
	}
	if endIndex != -1 {
		node.Body = node.Body[:endIndex]
	}

	if funcTag, ok := fnScope.Tag.(*FuncTag); ok {
		if !funcTag.GetIfExits(Return) && !ret.Eq(&EmptyReturn) {
			w.error(node.GetToken(), "not all code paths return a value")
		}
	}

	return variable
}

func HasContents[T any](contents ...[]T) bool {
	sumContents := make([]T, 0)
	for _, v := range contents {
		sumContents = append(sumContents, v...)
	}
	return len(sumContents) != 0
}

func (w *Walker) returnStmt(node *ast.ReturnStmt, scope *Scope) *Types {
	if !scope.Is(ReturnAllowing) {
		w.error(node.GetToken(), "can't have a return statement outside of a function or method")
	}

	ret := EmptyReturn
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope)
		valType := val.GetType()
		if types, ok := val.(*Types); ok {
			ret = append(ret, *types...)
		} else {
			ret = append(ret, valType)
		}
	}
	sc, _, funcTag := ResolveTagScope[*FuncTag](scope)
	if sc == nil {
		return &ret
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Return)
	}

	errorMsg := w.validateReturnValues(ret, (*funcTag).ReturnType)
	if errorMsg != "" {
		w.error(node.GetToken(), errorMsg)
	}

	return &ret
}

func (w *Walker) yieldStmt(node *ast.YieldStmt, scope *Scope) *Types {
	if !scope.Is(YieldAllowing) {
		w.error(node.GetToken(), "cannot use yield outside of statement expressions") // wut
	}

	ret := EmptyReturn
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope)
		valType := val.GetType()
		if types, ok := val.(*Types); ok {
			ret = append(ret, *types...)
		} else {
			ret = append(ret, valType)
		}
	}

	sc, _, matchExprT := ResolveTagScope[*MatchExprTag](scope)

	if sc == nil {
		return &ret
	}

	matchExprTag := *matchExprT

	if matchExprTag.YieldValues == nil {
		matchExprTag.YieldValues = &ret
	} else {
		errorMsg := w.validateReturnValues(ret, *matchExprTag.YieldValues)
		if errorMsg != "" {
			errorMsg = strings.Replace(errorMsg, "return", "yield", -1)
			w.error(node.GetToken(), errorMsg)
		}
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Yield)
	}

	return &ret
}

func (w *Walker) breakStmt(node *ast.BreakStmt, scope *Scope) {
	if !scope.Is(BreakAllowing) {
		w.error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Break)
	}
}

func (w *Walker) continueStmt(node *ast.ContinueStmt, scope *Scope) {
	if !scope.Is(ContinueAllowing) {
		w.error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Continue)
	}
}

func (w *Walker) repeatStmt(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope, &MultiPathTag{})
	repeatScope.Attributes.Add(BreakAllowing)
	repeatScope.Attributes.Add(ContinueAllowing)

	end := w.GetNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	if !parser.IsFx(endType.Type) && endType.Type != ast.Number {
		w.error(node.Iterator.GetToken(), "invalid value type of iterator")
	} else if variable, ok := end.(*VariableVal); ok {
		if fixedpoint, ok := variable.Value.(*FixedVal); ok {
			endType = Type{Type: fixedpoint.SpecificType}
		}
	} else {
		if fixedpoint, ok := end.(*FixedVal); ok {
			endType = Type{Type: fixedpoint.SpecificType}
		}
	}
	if node.Start.GetType() == ast.NA {
		node.Start = &ast.LiteralExpr{Token: node.Start.GetToken(), ValueType: endType.Type, Value: "1"}
	}
	start := w.GetNodeValue(&node.Start, scope)
	if node.Skip.GetType() == ast.NA {
		node.Skip = &ast.LiteralExpr{Token: node.Skip.GetToken(), ValueType: endType.Type, Value: "1"}
	}
	skip := w.GetNodeValue(&node.Skip, scope)

	repeatType := end.GetType().Type
	startType := start.GetType().Type
	skipType := skip.GetType().Type

	if (repeatType != startType || startType == 0) &&
		(repeatType != skipType || skipType == 0) {
		w.error(node.Start.GetToken(), fmt.Sprintf("all value types must be the same (iter:%s, start:%s, by:%s)", repeatType.ToString(), startType.ToString(), skipType.ToString()))
	}

	if node.Variable.GetValueType() != 0 {
		repeatScope.
			DeclareVariable(VariableVal{Name: node.Variable.Name.Lexeme, Value: w.GetNodeValue(&node.Start, scope), Node: node})
	}

	endIndex := -1
	for i := range node.Body {
		loopTag := *helpers.GetValOfInterface[*MultiPathTag](repeatScope.Tag)
		if HasContents(loopTag.Breaks, loopTag.Returns, loopTag.Continues, loopTag.Yields) {
			w.warn(node.Body[i].GetToken(), "unreachable code detected")
			endIndex = i
		}
		w.WalkNode(&node.Body[i], &repeatScope)
	}
	if endIndex != -1 {
		node.Body = node.Body[:endIndex]
	}

	returnabl := repeatScope.Parent.ResolveReturnable()
	loopTag := *helpers.GetValOfInterface[*MultiPathTag](repeatScope.Tag)

	if returnabl == nil {
		return
	}
	returnable := *returnabl

	returnable.SetExit(loopTag.GetIfExits(Return), Return)
	returnable.SetExit(loopTag.GetIfExits(Yield), Yield)
	returnable.SetExit(loopTag.GetIfExits(Break), Break)
	returnable.SetExit(loopTag.GetIfExits(Continue), Continue)
}

func (w *Walker) forStmt(node *ast.ForStmt, scope *Scope) {
	forScope := NewScope(scope, &UntaggedTag{})
	forScope.Attributes.Add(BreakAllowing)
	forScope.Attributes.Add(ContinueAllowing)

	if node.Key.GetValueType() != 0 {
		forScope.DeclareVariable(VariableVal{Name: node.Key.Name.Lexeme, Value: &NumberVal{}})
	}

	if node.Value.GetValueType() != 0 {
		forScope.DeclareVariable(VariableVal{Name: node.Value.Name.Lexeme, Value: w.GetValueFromType(*w.GetNodeValue(&node.Iterator, scope).GetType().WrappedType)})
	}

	iteratorType := w.GetNodeValue(&node.Iterator, scope).GetType().Type
	if iteratorType != ast.List && iteratorType != ast.Map {
		w.error(node.Iterator.GetToken(), "iterator must be of type map or list")
	}

	for i := range node.Body {
		w.WalkNode(&node.Body[i], &forScope)
	}
}

func (w *Walker) tickStmt(node *ast.TickStmt, scope *Scope) {
	tickScope := NewScope(scope, &UntaggedTag{})

	if node.Variable.GetValueType() != 0 {
		tickScope.DeclareVariable(VariableVal{Name: node.Variable.Name.Lexeme, Value: &NumberVal{}})
	}

	for i := range node.Body {
		w.WalkNode(&node.Body[i], &tickScope)
	}
}

func (w *Walker) AddTypesToValues(list *[]Value, tys *Types) {
	for _, typ := range *tys {
		val := w.GetValueFromType(typ)
		*list = append(*list, val)
	}
}

func (w *Walker) variableDeclarationStmt(declaration *ast.VariableDeclarationStmt, scope *Scope) []VariableVal {
	declaredVariables := []VariableVal{}

	var values []Value

	for i := range declaration.Values {

		exprValue := w.GetNodeValue(&declaration.Values[i], scope)
		if declaration.Values[i].GetType() == ast.SelfExpression {
			w.error(declaration.Values[i].GetToken(), "cannot assign self to a variable")
		}
		if types, ok := exprValue.(*Types); ok {
			w.AddTypesToValues(&values, types)
		} else {
			values = append(values, exprValue)
		}
	}
	valuesLength := len(values)
	if valuesLength > len(declaration.Identifiers) {
		w.error(declaration.Token, "too many values provided in declaration")
		return declaredVariables
	} else if valuesLength < len(declaration.Identifiers) {
		w.error(declaration.Token, "too few values provided in declaration")
		return declaredVariables
	}

	if !declaration.IsLocal {
		w.error(declaration.Token, "cannot declare a global variable inside a local block")
	}
	if declaration.Token.Type == lexer.Const && scope.Parent != nil {
		w.error(declaration.Token, "cannot declare a global constant inside a local block")
	}

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}

		valType := values[i].GetType()

		if declaration.Types[i] != nil {
			if explicitType := w.typeExpr(declaration.Types[i]); !TypeEquals(&valType, &explicitType) {
				w.error(declaration.Token, fmt.Sprintf("mismatched types: explict type (%s) not the same with value type (%s)",
					valType.ToString(),
					explicitType.ToString()))
			}
		}

		variable := VariableVal{
			Value: values[i],
			Name:  ident.Lexeme,
			Node:  declaration,
		}
		declaredVariables = append(declaredVariables, variable)
		if _, success := scope.DeclareVariable(variable); !success {
			w.error(lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	return declaredVariables
}

func (w *Walker) structDeclarationStmt(node *ast.StructDeclarationStmt, scope *Scope) {
	structTypeVal := StructTypeVal{
		Name:         node.Name,
		Methods:      map[string]VariableVal{},
		Fields:       []VariableVal{},
		FieldIndexes: map[string]int{},
	}

	structScope := NewScope(scope, &StructTag{StructType: &structTypeVal})
	structScope.Attributes.Add(SelfAllowing)

	params := make([]Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, w.typeExpr(param.Type))
	}
	structTypeVal.Params = params

	scope.DeclareStructType(&structTypeVal)
	w.Environment.foreignTypes[structTypeVal.Name.Lexeme] = &structTypeVal

	funcDeclaration := ast.MethodDeclarationStmt{
		Name:    node.Constructor.Token,
		Params:  node.Constructor.Params,
		Return:  node.Constructor.Return,
		IsLocal: true,
		Body:    *node.Constructor.Body,
	}

	for i := range node.Fields {
		w.fieldDeclarationStmt(&node.Fields[i], &structTypeVal, &structScope)
	}

	structTypeVal.FieldIndexes = structScope.VariableIndexes

	for i := range *node.Methods {
		params := make([]Type, 0)
		for _, param := range (*node.Methods)[i].Params {
			params = append(params, w.typeExpr(param.Type))
		}

		ret := EmptyReturn
		for _, typee := range (*node.Methods)[i].Return {
			ret = append(ret, w.typeExpr(typee))
			//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
		}
		variable := VariableVal{
			Name:  (*node.Methods)[i].Name.Lexeme,
			Value: &FunctionVal{params: params, returnVal: ret},
			Node:  &(*node.Methods)[i],
		}
		if _, success := structScope.DeclareVariable(variable); !success {
			w.error((*node.Methods)[i].Name, fmt.Sprintf("variable with name '%s' already exists", variable.Name))
		}
		structTypeVal.Methods[variable.Name] = variable
	}

	for i := range *node.Methods {
		w.methodDeclarationStmt(&(*node.Methods)[i], &structTypeVal, &structScope)
	}

	w.methodDeclarationStmt(&funcDeclaration, &structTypeVal, &structScope)
}

func (w *Walker) fieldDeclarationStmt(node *ast.FieldDeclarationStmt, container Container, scope *Scope) {
	varDecl := ast.VariableDeclarationStmt{
		Identifiers: node.Identifiers,
		Types:       node.Types,
		Values:      node.Values,
		IsLocal:     true,
		Token:       node.Token,
	}
	structType := container.GetType()
	if len(node.Types) != 0 {
		for i := range node.Types {
			explicitType := w.typeExpr(node.Types[i])
			if TypeEquals(&explicitType, &structType) {
				w.error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	} else if len(node.Types) != 0 {
		for i := range node.Values {
			valType := w.GetNodeValue(&node.Values[i], scope).GetType()
			if TypeEquals(&valType, &structType) {
				w.error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	}

	variables := w.variableDeclarationStmt(&varDecl, scope)
	node.Values = varDecl.Values
	for i := range variables {
		container.AddField(variables[i])
	}
}

func (w *Walker) methodDeclarationStmt(node *ast.MethodDeclarationStmt, container Container, scope *Scope) {
	funcExpr := ast.FunctionDeclarationStmt{
		Name:    node.Name,
		Return:  node.Return,
		Params:  node.Params,
		Body:    node.Body,
		IsLocal: true,
	}

	variable := w.functionDeclarationStmt(&funcExpr, scope, Method)
	node.Body = funcExpr.Body
	container.AddMethod(variable)
}

func (w *Walker) useStmt(node *ast.UseStmt, scope *Scope) {
	variable := VariableVal{Name: node.Variable.Name.Lexeme, Value: &EnvironmentVal{Name: node.Variable.Name.Lexeme}, Node: node}

	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Variable.Name, "cannot declare a value in the same scope twice")
	}
}

func (w *Walker) matchStmt(node *ast.MatchStmt, isExpr bool, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	valType := val.GetType()

	multiPathScope := NewScope(scope, &MultiPathTag{})

	var has_default bool
	for _, matchCase := range node.Cases {
		caseScope := NewScope(&multiPathScope, &UntaggedTag{})

		if matchCase.Expression.GetToken().Lexeme == "_" {
			has_default = true
		}

		if !isExpr {
			for i := range matchCase.Body {
				w.WalkNode(&matchCase.Body[i], &caseScope)
			}
		}
		if caseScope.Tag.GetType() == Untagged {
			continue
		}
		caseValType := w.GetNodeValue(&matchCase.Expression, scope).GetType()
		if !TypeEquals(&valType, &caseValType) {
			w.error(
				matchCase.Expression.GetToken(),
				fmt.Sprintf("mismatched types: arm expression (%s) and match expression (%s)",
					caseValType.ToString(),
					valType.ToString()))
		}
	}

	if has_default && len(node.Cases) == 1 {
		w.error(node.Cases[0].Expression.GetToken(), "cannot have a match statement/expression with one arm that is default")
	}

	if !has_default && isExpr {
		w.error(node.GetToken(), "match expression has no default arm")
	}

	if isExpr {
		return
	}

	returnabl := scope.ResolveReturnable()

	if returnabl == nil {
		return
	}

	returnable := *returnabl

	mpTag := multiPathScope.Tag.(*MultiPathTag)

	if has_default {
		returnable.SetExit(mpTag.GetIfExits(Yield), Yield)
		returnable.SetExit(mpTag.GetIfExits(Return), Yield)
		returnable.SetExit(mpTag.GetIfExits(Break), Yield)
		returnable.SetExit(mpTag.GetIfExits(Continue), Yield)

		// if (len(mpTag.Returns) + len(mpTag.Breaks) +
		// 	len(mpTag.Continues) + len(mpTag.Yields)) > 1 {
		// 	returnable.SetExit(false, Return)
		// }
	}else {
		returnable.SetExit(false, Return)
		returnable.SetExit(false, Yield)
		returnable.SetExit(false, Break)
		returnable.SetExit(false, Continue)
	}
}

func (w *Walker) envStmt(node *ast.EnvironmentStmt, scope *Scope) {
	if scope.Environment.Name != "" {
		w.error(node.GetToken(), "can't have more than one environment statement in a file")
		return
	}

	for i, v := range node.Env.Envs {
		if i < len(node.Env.Envs)-1 {
			scope.Environment.Name += v.Lexeme + "::"
		}else {
			scope.Environment.Name += v.Lexeme
		}
	}
	 
	if wlkr, found := (*w.Walkers)[w.Environment.Name]; found {
		w.error(node.GetToken(), fmt.Sprintf("cannot have two environments with the same name, path: %s",wlkr.Environment.Path))
		return
	}

	(*w.Walkers)[w.Environment.Name] = w
}