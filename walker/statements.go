package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
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

// Rewrote
func (w *Walker) assignmentStatement(assignStmt *ast.AssignmentStmt, scope *Scope) {
	type Value2 struct {
		Value
		Index int
	}

	values := []Value2{}
	for i := range assignStmt.Values {
		exprValue := w.GetNodeActualValue(&assignStmt.Values[i], scope)
		if vals, ok := exprValue.(Values); ok {
			for j := range vals {
				values = append(values, Value2{vals[j], i})
			}
		} else {
			values = append(values, Value2{exprValue, i})
		}
	}

	idents := assignStmt.Identifiers

	identsLen := len(idents)
	valuesLen := len(values)

	exprs := assignStmt.Values
	binExprs := make([]ast.Node, 0) // for compound assignment

	assignOp := assignStmt.AssignOp
	for i := range assignStmt.Identifiers {
		if i >= valuesLen {
			requiredAmount := identsLen - valuesLen
			w.AlertSingle(&alerts.TooFewValuesGiven{}, exprs[len(exprs)-1].GetToken(), requiredAmount, "assignment")
			return
		}
		value := w.GetNodeValue(&idents[i], scope)
		variable, ok := value.(*VariableVal)
		if !ok {
			continue
		}
		if variable.IsConst {
			w.AlertSingle(&alerts.ConstValueAssignment{}, idents[i].GetToken())
			continue
		}
		if !variable.IsInit {
			variable.IsInit = true
		}

		variableType := variable.GetType()
		valType := values[i].GetType()
		if valType == InvalidType {
			continue
		}

		if assignOp.Type != tokens.Equal {
			if !isNumerical(valType.PVT()) {
				w.AlertSingle(&alerts.InvalidTypeInCompoundAssignment{}, exprs[values[i].Index].GetToken(),
					valType.ToString(),
				)
				continue
			}
			if !isNumerical(variableType.PVT()) {
				w.AlertSingle(&alerts.InvalidTypeInCompoundAssignment{}, idents[i].GetToken(),
					variableType.ToString(),
				)
				continue
			}

			binExpr := &ast.BinaryExpr{
				Left:     idents[i],
				Right:    exprs[values[i].Index],
				Operator: assignOp,
			}
			binExprs = append(binExprs, binExpr)
		}

		if !TypeEquals(variableType, valType) {
			w.AlertSingle(&alerts.AssignmentTypeMismatch{}, exprs[values[i].Index].GetToken(),
				variableType.ToString(),
				valType.ToString(),
			)
			continue
		}
	}
	if valuesLen > identsLen {
		extraAmount := valuesLen - identsLen
		w.AlertMulti(&alerts.TooManyValuesGiven{},
			exprs[values[valuesLen-1].Index].GetToken(),
			exprs[values[valuesLen-extraAmount].Index].GetToken(),
			extraAmount,
			"in assignment",
		)
		return
	}
	if len(binExprs) != 0 {
		assignStmt.Values = binExprs
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
	w.GetNodeValue(&node.Condition, scope)

	w.WalkBody(&node.Body, lt, whileScope)
}

func (w *Walker) forStatement(node *ast.ForStmt, scope *Scope) {
	forScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, forScope.Attributes...)
	forScope.Tag = lt

	valType := w.GetNodeValue(&node.Iterator, scope).GetType()
	wrapper, ok := valType.(*WrapperType)
	if !ok {
		w.AlertSingle(&alerts.InvalidIteratorType{}, node.Iterator.GetToken(), valType.ToString())
		return
	}
	node.OrderedIteration = wrapper.PVT() == ast.List

	if node.First.Name.Lexeme != "_" {
		var firstValue Value
		if node.OrderedIteration {
			firstValue = &NumberVal{}
		} else {
			firstValue = &StringVal{}
		}
		w.DeclareVariable(forScope, NewVariable(node.First.Name, firstValue))
	}

	if node.Second != nil && ok {
		if node.Second.Name.Lexeme == "_" && node.First.Name.Lexeme != "_" {
			w.AlertSingle(&alerts.UnnecessaryEmptyIdentifier{}, node.Second.Name, "in for loop statement")
		}
		if node.Second.Name.Lexeme != "_" {
			w.DeclareVariable(forScope, NewVariable(node.Second.Name, w.TypeToValue(wrapper.WrappedType)))
		}
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
		w.AlertSingle(&alerts.InvalidUseOfExitStmt{}, node.Token,
			"break",
			"for loops",
		)
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Break)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) continueStatement(node *ast.ContinueStmt, scope *Scope) {
	if !scope.Is(ContinueAllowing) {
		w.AlertSingle(&alerts.InvalidUseOfExitStmt{}, node.Token,
			"continue",
			"for loops",
		)
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Continue)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) returnStatement(node *ast.ReturnStmt, scope *Scope) *[]Type {
	if !scope.Is(ReturnAllowing) {
		w.AlertSingle(&alerts.InvalidUseOfExitStmt{}, node.Token,
			"return",
			"a function or method",
		)
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

	w.ValidateReturnValues(node.Args, ret, (*funcTag).ReturnTypes, "in return arguments") // wait

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Return)
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) yieldStatement(node *ast.YieldStmt, scope *Scope) *[]Type {
	if !scope.Is(YieldAllowing) {
		w.AlertSingle(&alerts.InvalidUseOfExitStmt{}, node.Token,
			"yield",
			"a match expression",
		)
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
		w.ValidateReturnValues(node.Args, ret, matchExprTag.YieldTypes, "in yield arguments")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Yield)
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) useStatement(node *ast.UseStmt, scope *Scope) {
	if scope.Parent != nil {
		w.AlertSingle(&alerts.UseStmtInLocalBlock{}, node.Token)
		return
	}

	envName := node.PathExpr.Path.Lexeme

	switch envName {
	case "Pewpew":
		w.environment.UsedLibraries[Pewpew] = true
		return
	}

	switch envName {
	case "Pewpew":
		if w.environment.Type != ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Pewpew", "Mesh or Sound")
		}
		w.environment.UsedLibraries[Pewpew] = true
		return
	case "Fmath":
		if w.environment.Type != ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Fmath", "Mesh or Sound")
		}
		w.environment.UsedLibraries[Fmath] = true
		return
	case "Math":
		if w.environment.Type == ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Math", "Level")
		}
		w.environment.UsedLibraries[Math] = true
		return
	case "String":
		w.environment.UsedLibraries[String] = true
		return
	case "Table":
		w.environment.UsedLibraries[Table] = true
		return
	}

	walker, found := w.walkers[envName]

	if !found {
		w.AlertSingle(&alerts.InvalidEnvironmentAccess{}, node.PathExpr.Path, envName)
		return
	}

	for i := range walker.environment.importedWalkers {
		if walker.environment.importedWalkers[i].environment.Name == w.environment.Name {
			w.AlertSingle(&alerts.ImportCycle{}, node.PathExpr.Path, w.environment.hybroidPath, walker.environment.hybroidPath)
			return
		}
	}

	if walker.environment.luaPath == "/dynamic/level.lua" {
		w.environment.importedWalkers = append(w.environment.importedWalkers, walker)
		return // we don't put level.hyb in requirements as that would break things
	}

	success := w.environment.AddRequirement(walker.environment.luaPath)

	if !success {
		w.AlertSingle(&alerts.EnvironmentReuse{}, node.PathExpr.Path, envName)
		return
	}

	w.environment.importedWalkers = append(w.environment.importedWalkers, walker)

	if !walker.Walked {
		walker.Pass2()
	}
}

func (w *Walker) destroyStatement(node *ast.DestroyStmt, scope *Scope) {
	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		return
	}
	if valType.PVT() != ast.Entity {
		w.AlertSingle(&alerts.InvalidTypeInDestroyStatement{}, node.Identifier.GetToken(), valType.ToString())
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

	suppliedGenerics := w.GetGenerics(node.GenericArgs, entityVal.DestroyGenerics, scope)
	w.ValidateArguments(suppliedGenerics, args, entityVal.DestroyParams, node)
}
