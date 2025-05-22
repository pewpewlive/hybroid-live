package walker

import (
	"hybroid/ast"
	"hybroid/helpers"
	"strings"
)

func (w *Walker) DeclareConversion(scope *Scope) {
	if len(w.context.Conversions) == 1 {
		conv := w.context.Conversions[0]
		w.DeclareVariable(scope, &VariableVal{
			Name:   conv.Name.Lexeme,
			Value:  conv.Entity,
			IsInit: true,
		})
	}
	w.context.Conversions = make([]EntityConversion, 0)
}

func (w *Walker) ifStatement(node *ast.IfStmt, scope *Scope) {
	length := len(node.Elseifs) + 2
	mpt := NewMultiPathTag(length, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)
	ifScope := NewScope(multiPathScope, &UntaggedTag{})

	w.context.Conversions = make([]EntityConversion, 0)

	boolExpr := w.GetNodeValue(&node.BoolExpr, scope)
	if boolExpr.GetType().PVT() != ast.Bool {
		// w.Error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	w.DeclareConversion(ifScope)
	w.WalkBody(&node.Body, mpt, ifScope)

	for i := range node.Elseifs {
		boolExpr := w.GetNodeValue(&node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().PVT() != ast.Bool {
			// w.Error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := NewScope(multiPathScope, &UntaggedTag{})
		w.DeclareConversion(ifScope)
		w.WalkBody(&node.Elseifs[i].Body, mpt, ifScope)
	}

	if node.Else != nil {
		elseScope := NewScope(multiPathScope, &UntaggedTag{})
		w.WalkBody(&node.Else.Body, mpt, elseScope)
	}
}

func (w *Walker) assignmentStatement(assignStmt *ast.AssignmentStmt, scope *Scope) {
	type Value2 struct {
		Value
		Index int
	}

	values := []Value2{}
	index := 1
	for i := range assignStmt.Values { // function()
		index++
		exprValue := w.GetNodeValue(&assignStmt.Values[i], scope)
		// if assignStmt.Values[i].GetType() == ast.SelfExpression {
		// 	// w.Error(assignStmt.Values[i].GetToken(), "cannot assign self to a variable")
		// }
		if vals, ok := exprValue.(Values); ok {
			for j := range vals {
				values = append(values, Value2{vals[j], i})
			}
		} else {
			values = append(values, Value2{exprValue, i})
		}
	}

	variablesLength := len(assignStmt.Identifiers)
	valuesLength := len(values)
	if variablesLength < valuesLength {
		// w.Error(assignStmt.Token, "too many values given in variable declaration")
	} else if variablesLength > valuesLength {
		// w.Error(assignStmt.Token, "too few values given in variable declaration")
	}

	for i := index; i < variablesLength; i++ {
		values = append(values, Value2{&Invalid{}, i})
	}

	for i := range assignStmt.Identifiers {
		value := w.GetNodeValue(&assignStmt.Identifiers[i], scope)
		variable, ok := value.(*VariableVal)
		if !ok {
			variable = &VariableVal{
				Name:   "",
				Value:  value,
				IsInit: true,
			}
		}
		if variable.IsConst {
			// variableToken := assignStmt.Identifiers[i].GetToken()
			// w.Error(variableToken, "cannot modify '%s' because it is const", variableToken.Lexeme)
			continue
		}

		variableType := variable.GetType()
		if !variable.IsInit {
			variable.IsInit = true
		}

		valType := values[i].GetType()

		if !TypeEquals(variableType, valType) {
			// variableName := assignStmt.Identifiers[i].GetToken().Lexeme
			// w.Error(assignStmt.Values[values[i].Index].GetToken(), "mismatched types: '%s' is of type %s but a value of %s was given to it", variableName, variableType.ToString(), valType.ToString())
		}

		if vr, ok := values[i].Value.(*VariableVal); ok {
			values[i] = Value2{vr.Value, values[i].Index}
		}

		//variable.Value = values[i]
	}
}

func (w *Walker) repeatStatement(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, repeatScope.Attributes...)
	repeatScope.Tag = lt

	end := w.GetNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	//if !parser.IsFx(endType.PVT()) && endType.PVT() != ast.Number {
	// w.Error(node.Iterator.GetToken(), "invalid value type of iterator")
	//} else if variable, ok := end.(*VariableVal); ok {
	//if fixedpoint, ok := variable.Value.(*FixedVal); ok {
	//	endType = NewBasicType(fixedpoint.SpecificType)
	//}
	//} else {
	if fixedpoint, ok := end.(*FixedVal); ok {
		endType = NewBasicType(fixedpoint.SpecificType)
	}
	//}
	if node.Start.GetType() == ast.NA {
		node.Start = &ast.LiteralExpr{Token: node.Start.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	start := w.GetNodeValue(&node.Start, scope)
	if node.Skip.GetType() == ast.NA {
		node.Skip = &ast.LiteralExpr{Token: node.Skip.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	skip := w.GetNodeValue(&node.Skip, scope)

	repeatType := end.GetType()
	startType := start.GetType()
	skipType := skip.GetType()

	if !(TypeEquals(repeatType, startType) && TypeEquals(startType, skipType)) {
		// w.Error(node.Start.GetToken(), fmt.Sprintf("all value types must be the same (iter:%s, start:%s, by:%s)", repeatType.ToString(), startType.ToString(), skipType.ToString()))
	}

	if node.Variable != nil {
		w.DeclareVariable(repeatScope, NewVariable(node.Variable.Name, w.TypeToValue(repeatType)))
	}

	w.WalkBody(&node.Body, lt, repeatScope)
}

func (w *Walker) whileStatement(node *ast.WhileStmt, scope *Scope) {
	whileScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, whileScope.Attributes...)
	whileScope.Tag = lt

	_ = w.GetNodeValue(&node.Condition, scope)

	w.WalkBody(&node.Body, lt, whileScope)
}

func (w *Walker) forStatement(node *ast.ForStmt, scope *Scope) {
	forScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, forScope.Attributes...)
	forScope.Tag = lt

	w.DeclareVariable(forScope, NewVariable(node.First.Name, &NumberVal{}))

	valType := w.GetNodeValue(&node.Iterator, scope).GetType()
	wrapper, ok := valType.(*WrapperType)
	if !ok {
		// w.Error(node.Iterator.GetToken(), "iterator must be of type map or list")
	} else if node.Second != nil {
		node.OrderedIteration = wrapper.PVT() == ast.List
		w.DeclareVariable(forScope, NewVariable(node.Second.Name, w.TypeToValue(wrapper)))
	}

	w.WalkBody(&node.Body, lt, forScope)
}

func (w *Walker) tickStatement(node *ast.TickStmt, scope *Scope) {
	funcTag := &FuncTag{ReturnTypes: EmptyReturn}
	tickScope := NewScope(scope, funcTag, ReturnAllowing)

	if node.Variable != nil {
		w.DeclareVariable(tickScope, NewVariable(node.Variable.Name, &NumberVal{}))
	}

	w.WalkBody(&node.Body, funcTag, tickScope)
}

func (w *Walker) matchStatement(node *ast.MatchStmt, isExpr bool, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	if val.GetType().PVT() == ast.Invalid {
		// w.Error(node.ExprToMatch.GetToken(), "variable is of type invalid")
	}
	valType := val.GetType()
	casesLength := len(node.Cases) + 1
	if node.HasDefault {
		casesLength--
	}
	mpt := NewMultiPathTag(casesLength, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)

	for i := range node.Cases {
		caseScope := NewScope(multiPathScope, &UntaggedTag{})

		if !isExpr {
			w.WalkBody(&node.Cases[i].Body, mpt, caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "else" {
			continue
		}

		caseValType := w.GetNodeValue(&node.Cases[i].Expression, scope).GetType()
		if !TypeEquals(valType, caseValType) {
			// w.Error(
			// node.Cases[i].Expression.GetToken(),
			// fmt.Sprintf("mismatched types: arm expression (%s) and match expression (%s)",
			// 	caseValType.ToString(),
			// 	valType.ToString()))
		}
	}
}

func (w *Walker) breakStatement(node *ast.BreakStmt, scope *Scope) {
	if !scope.Is(BreakAllowing) {
		// w.Error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Break)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) continueStatement(node *ast.ContinueStmt, scope *Scope) {
	if !scope.Is(ContinueAllowing) {
		// w.Error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Continue)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) returnStatement(node *ast.ReturnStmt, scope *Scope) *[]Type {
	if !scope.Is(ReturnAllowing) {
		// w.Error(node.GetToken(), "can't have a return statement outside of a function or method")
	}

	ret := EmptyReturn
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope) // we need to check waht happens here
		valType := val.GetType()
		if vls, ok := val.(Values); ok {
			for i := range vls {
				ret = append(ret, vls[i].GetType())
			}
		} else {
			ret = append(ret, valType)
		}
	}
	sc, _, funcTag := ResolveTagScope[*FuncTag](scope)
	if sc == nil {
		return &ret
	}

	errorMsg := w.ValidateReturnValues(ret, (*funcTag).ReturnTypes) // wait
	if errorMsg != "" {
		// w.Error(node.GetToken(), errorMsg)
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Return)
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) yieldStatement(node *ast.YieldStmt, scope *Scope) *[]Type {
	if !scope.Is(YieldAllowing) {
		// w.Error(node.GetToken(), "cannot use yield outside of statement expressions") // wut
	}

	ret := EmptyReturn
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope)
		valType := val.GetType()
		if vls, ok := val.(Values); ok {
			for i := range vls {
				ret = append(ret, vls[i].GetType())
			}
		} else {
			ret = append(ret, valType)
		}
	}

	sc, _, matchExprT := ResolveTagScope[*MatchExprTag](scope)

	if sc == nil {
		return &ret
	}

	matchExprTag := *matchExprT

	if helpers.ListsAreSame(matchExprTag.YieldTypes, EmptyReturn) {
		matchExprTag.YieldTypes = ret
	} else {
		errorMsg := w.ValidateReturnValues(ret, matchExprTag.YieldTypes)
		if errorMsg != "" {
			errorMsg = strings.Replace(errorMsg, "return", "yield", -1)
			// w.Error(node.GetToken(), errorMsg)
		}
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Yield)
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) useStatement(node *ast.UseStmt, scope *Scope) {
	if scope.Parent != nil {
		// w.Error(node.GetToken(), "cannot have a use statement inside a local block")
		return
	}

	path := node.Path.Path.Lexeme

	switch path {
	case "Pewpew":
		w.environment.UsedLibraries[Pewpew] = true
		return
	}

	switch path {
	case "Pewpew":
		if w.environment.Type != ast.LevelEnv {
			// w.Error(node.GetToken(), "cannot use the pewpew library in a non-level environment")
		}
		w.environment.UsedLibraries[Pewpew] = true
		return
	case "Fmath":
		if w.environment.Type != ast.LevelEnv {
			// w.Error(node.GetToken(), "cannot use the fmath library in a non-level environment")
		}
		w.environment.UsedLibraries[Fmath] = true
		return
	case "Math":
		w.environment.UsedLibraries[Math] = true
		return
	case "String":
		w.environment.UsedLibraries[String] = true
		return
	case "Table":
		w.environment.UsedLibraries[Table] = true
		return
	}

	envName := node.Path.Path.Lexeme
	walker, found := w.walkers[envName]

	if !found {
		// w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return
	}

	if walker.environment.luaPath == "/dynamic/level.lua" {
		w.environment.importedWalkers = append(w.environment.importedWalkers, walker)
		return
	}

	for _, v := range walker.environment.Requirements() {
		if v == w.environment.luaPath {
			// w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
			return
		}
	}

	success := w.environment.AddRequirement(walker.environment.luaPath)

	if !success {
		// w.Error(node.GetToken(), fmt.Sprintf("Environment '%s' is already used", envName))
		return
	}

	w.environment.importedWalkers = append(w.environment.importedWalkers, walker)
}

func (w *Walker) destroyStatement(node *ast.DestroyStmt, scope *Scope) {
	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		// w.Error(node.Identifier.GetToken(), "invalid variable given in destroy expression")
		return
	} else if valType.PVT() != ast.Entity {
		// w.Error(node.Identifier.GetToken(), "variable given in destroy statement is not an entity")
		return
	}

	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}

	entityVal := val.(*EntityVal)

	node.EnvName = entityVal.Type.EnvName
	node.EntityName = entityVal.Type.Name

	args := make([]Type, 0)
	for i := range node.Args {
		args = append(args, w.GetNodeValue(&node.Args[i], scope).GetType())
	}

	suppliedGenerics := w.GetGenerics(node, node.GenericArgs, entityVal.DestroyGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, entityVal.DestroyParams, node)
}
