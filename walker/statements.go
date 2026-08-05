package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
	"slices"
	"strings"
)

func (w *Walker) ifStatement(node *ast.IfStmt, scope *Scope) {
	w.context.EntityCasts.Clear()

	w.ifCondition(&node.BoolExpr, scope)

	for w.context.EntityCasts.Count() != 0 {
		cast := w.context.EntityCasts.Pop()
		w.declareVariable(scope, NewVariable(cast.Name, cast.Entity))
	}

	pt := NewPathTag()
	ifScope := w.NewScope(scope, pt)
	w.RegisterScope(ifScope, node.Token, w.GetBodyEndToken(&node.Body))

	w.walkBody(&node.Body, pt, ifScope)

	prevPathTag := *pt
	for i := range node.Elseifs {
		w.ifCondition(&node.Elseifs[i].BoolExpr, scope)
		pt := NewPathTag()
		ifScope := w.NewScope(scope, pt)
		w.RegisterScope(ifScope, node.Elseifs[i].Token, w.GetBodyEndToken(&node.Elseifs[i].Body))
		for w.context.EntityCasts.Count() != 0 {
			cast := w.context.EntityCasts.Pop()
			w.declareVariable(scope, NewVariable(cast.Name, cast.Entity))
		}
		w.walkBody(&node.Elseifs[i].Body, pt, ifScope)
		prevPathTag.SetAllExitAND(pt)
	}

	if node.Else != nil {
		pt := NewPathTag()
		elseScope := w.NewScope(scope, pt)
		w.RegisterScope(elseScope, node.Else.Token, w.GetBodyEndToken(&node.Else.Body))
		w.walkBody(&node.Else.Body, pt, elseScope)
		prevPathTag.SetAllExitAND(pt)
	} else {
		prevPathTag.SetAllFalse()
	}

	w.reportExits(&prevPathTag, scope)
}

// Rewrote
func (w *Walker) assignmentStatement(assignStmt *ast.AssignmentStmt, scope *Scope) {
	values := []Value2{}
	idents := assignStmt.Identifiers
	exprs := assignStmt.Values
	identsLen := len(idents)
	assignOp := assignStmt.AssignOp
	exprCounter := 0
	for i := range assignStmt.Identifiers {
		w.context.DontSetToUsed = true
		value := w.GetNodeValue(&idents[i], scope)
		w.context.DontSetToUsed = false
		variable, ok := value.(*VariableVal)
		if !ok {
			if value.GetType() != InvalidType {
				w.Report(alerts.NewInvalidAssignment(idents[i].GetToken().Span))
			}
			continue
		}
		if variable.IsConst {
			w.Report(alerts.NewConstValueAssignment(idents[i].GetToken().Span))
			continue
		}
		if !variable.IsInit {
			variable.IsInit = true
		}
		variableType := variable.GetType()

		var valType Type
		if i <= len(values)-1 {
			valType = values[i].GetType()
		} else if exprCounter <= len(assignStmt.Values)-1 {
			exprValue := w.GetNodeValue(&assignStmt.Values[exprCounter], scope)
			if variable2, ok := exprValue.(*VariableVal); ok {
				if variable2 == variable && assignStmt.Values[exprCounter].GetType() == ast.Identifier && assignOp.Type == tokens.Equal {
					w.Report(alerts.NewAssignmentToSelf(assignStmt.Values[exprCounter].GetToken().Span, variable.Name))
				}
				exprValue = variable2.Value
			}
			if constVal, ok := exprValue.(*ConstVal); ok {
				exprValue = constVal.Val
			}
			if vals, ok := exprValue.(Values); ok {
				for j := range vals {
					values = append(values, Value2{vals[j], i})
				}
			} else {
				values = append(values, Value2{exprValue, i})
			}
			valType = values[i].GetType()
			exprCounter++
		} else {
			requiredAmount := identsLen - len(values)
			w.Report(alerts.NewTooFewElementsGiven(exprs[len(exprs)-1].GetToken().Span, requiredAmount, "value", "in assignment"))
			return
		}

		if valType == InvalidType || variableType == InvalidType {
			continue
		}

		if assignOp.Type == tokens.PipeEqual || assignOp.Type == tokens.AmpersandEqual || assignOp.Type == tokens.LeftShiftEqual || assignOp.Type == tokens.RightShiftEqual || assignOp.Type == tokens.TildeEqual {
			if valType.PVT() != ast.Number {
				w.Report(alerts.NewInvalidType(exprs[values[i].Index].GetToken().Span, valType.String(), "in bitwise compound assignment"))
				continue
			}
			if variableType.PVT() != ast.Number {
				w.Report(alerts.NewInvalidType(idents[i].GetToken().Span, variableType.String(), "in bitwise compound assignment"))
				continue
			}
		} else if assignOp.Type != tokens.Equal {
			if !isNumerical(valType.PVT()) {
				w.Report(alerts.NewInvalidTypeInCompoundAssignment(exprs[values[i].Index].GetToken().Span, valType.String()))
				continue
			}
			if !isNumerical(variableType.PVT()) {
				w.Report(alerts.NewInvalidTypeInCompoundAssignment(idents[i].GetToken().Span, variableType.String()))
				continue
			}
		}

		if !TypeEquals(variableType, valType) {
			w.Report(alerts.NewAssignmentTypeMismatch(exprs[values[i].Index].GetToken().Span, variableType.String(), valType.String()))
			continue
		}
	}
	exprsLen := len(exprs)
	if exprCounter != exprsLen {
		for range exprsLen - exprCounter {
			val := w.GetActualNodeValue(&assignStmt.Values[exprCounter], scope)
			if vls, ok := val.(Values); ok {
				for _, v := range vls {
					values = append(values, Value2{v, exprCounter})
				}
			} else {
				values = append(values, Value2{val, exprCounter})
			}
			exprCounter++
		}
	}
	valuesLen := len(values)
	if valuesLen > identsLen {
		extraAmount := valuesLen - identsLen
		w.Report(alerts.NewTooManyElementsGiven(core.MergeSpans(exprs[values[valuesLen-extraAmount].Index].GetToken().Span, exprs[values[valuesLen-1].Index].GetToken().Span), extraAmount, "value", "in assignment"))
		return
	}
}

func (w *Walker) repeatStatement(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := w.NewScope(scope, &PathTag{}, BreakAllowing, ContinueAllowing)
	w.RegisterScope(repeatScope, node.Token, w.GetBodyEndToken(&node.Body))
	lt := NewPathTag()
	repeatScope.Tag = lt

	end := w.GetActualNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	if endType != InvalidType && !isNumerical(endType.PVT()) {
		w.Report(alerts.NewInvalidRepeatIterator(node.Iterator.GetToken().Span, endType.String()))
	}
	var start Value
	if node.Start == nil {
		start = w.typeToValue(endType)
		literal := start.GetDefault()
		literal.Value = strings.Replace(literal.Value, "0", "1", 1)
		node.Start = literal
	} else {
		start = w.GetNodeValue(&node.Start, scope)
	}
	var skip Value
	if node.Skip == nil {
		skip = w.typeToValue(endType)
		literal := skip.GetDefault()
		literal.Value = strings.Replace(literal.Value, "0", "1", 1)
		node.Skip = literal
	} else {
		skip = w.GetNodeValue(&node.Skip, scope)
	}

	startType := start.GetType()
	skipType := skip.GetType()
	if !(TypeEquals(endType, startType) && TypeEquals(startType, skipType)) {
		w.Report(alerts.NewInconsistentRepeatTypes(node.Token.Span, startType.String(), skipType.String(), endType.String()))
	}

	if node.Variable != nil {
		w.declareVariable(repeatScope, NewVariable(node.Variable.Name, end))
	}

	w.walkBody(&node.Body, lt, repeatScope)
	w.reportExits(lt, scope)
}

func (w *Walker) whileStatement(node *ast.WhileStmt, scope *Scope) {
	whileScope := w.NewScope(scope, &PathTag{}, BreakAllowing, ContinueAllowing)
	w.RegisterScope(whileScope, node.Token, w.GetBodyEndToken(&node.Body))
	lt := NewPathTag()
	whileScope.Tag = lt
	w.GetNodeValue(&node.Condition, scope)

	w.walkBody(&node.Body, lt, whileScope)
	w.reportExits(lt, scope)
}

func (w *Walker) forStatement(node *ast.ForStmt, scope *Scope) {
	forScope := w.NewScope(scope, &PathTag{}, BreakAllowing, ContinueAllowing)
	w.RegisterScope(forScope, node.Token, w.GetBodyEndToken(&node.Body))
	lt := NewPathTag()
	forScope.Tag = lt

	if node.IsEntity {
		valType := w.typeExpression(node.Iterator.(*ast.TypeExpr), scope)
		if nt, ok := valType.(*NamedType); ok && nt.Pvt == ast.Entity {
			node.EnvName = nt.EnvName
			node.EntityName = nt.Name
			if node.Second != nil {
				w.Report(alerts.NewTooManyElementsGiven(node.Second.Name.Span, 1, "for loop variable", ""))
			}
			if node.First.Name.Lexeme != "_" {
				w.declareVariable(forScope, NewVariable(node.First.Name, w.typeToValue(valType)))
			}
			w.walkBody(&node.Body, lt, forScope)
			w.reportExits(lt, scope)
			return
		} else {
			w.Report(alerts.NewInvalidEntityForLoopType(node.Iterator.GetToken().Span))
		}
	}

	valType := w.GetActualNodeValue(&node.Iterator, scope).GetType()
	wrapper, ok := valType.(*WrapperType)
	if !ok {
		w.Report(alerts.NewInvalidIteratorType(node.Iterator.GetToken().Span, valType.String()))
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
		w.declareVariable(forScope, NewVariable(node.First.Name, firstValue))
	}

	if node.Second != nil && ok {
		if node.Second.Name.Lexeme == "_" && node.First.Name.Lexeme != "_" {
			w.Report(alerts.NewUnnecessaryEmptyIdentifier(node.Second.Name.Span, "in for loop statement"))
		}
		if node.Second.Name.Lexeme != "_" {
			w.declareVariable(forScope, NewVariable(node.Second.Name, w.typeToValue(wrapper.WrappedType)))
		}
	}

	w.walkBody(&node.Body, lt, forScope)
	w.reportExits(lt, scope)
}

func (w *Walker) tickStatement(node *ast.TickStmt, scope *Scope) {
	tickScope := w.NewScope(scope, &PathTag{}, ReturnAllowing)
	w.RegisterScope(tickScope, node.Token, w.GetBodyEndToken(&node.Body))
	tt := NewPathTag()
	tickScope.Tag = tt

	if node.Variable != nil {
		w.declareVariable(tickScope, NewVariable(node.Variable.Name, &NumberVal{}))
	}

	w.walkBody(&node.Body, tt, tickScope)
	w.reportExits(tt, scope)
}

func (w *Walker) matchStatement(node *ast.MatchStmt, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	valType := val.GetType()
	casesLength := len(node.Cases)
	if !node.HasDefault {
		if casesLength < 1 {
			w.Report(alerts.NewInsufficientCases(node.Token.Span))
		}
	} else if casesLength < 2 {
		w.Report(alerts.NewInsufficientCases(node.Token.Span))
	}

	var prevPathTag PathTag
	for i := range node.Cases {
		pt := NewPathTag()
		caseScope := w.NewScope(scope, pt, BreakAllowing)
		w.RegisterScope(caseScope, node.Cases[i].GetToken(), w.GetBodyEndToken(&node.Cases[i].Body))
		w.walkBody(&node.Cases[i].Body, pt, caseScope)
		if i != 0 {
			prevPathTag.SetAllExitAND(pt)
		} else {
			prevPathTag = *pt
		}

		if node.Cases[i].Expressions[0].GetToken().Lexeme == "else" {
			if i != len(node.Cases)-1 {
				w.Report(alerts.NewInvalidDefaultCasePlacement(node.Cases[i].Expressions[0].GetToken().Span, "in match statement"))
			}
			continue
		}

		for j := range node.Cases[i].Expressions {
			caseValType := w.GetNodeValue(&node.Cases[i].Expressions[j], scope).GetType()
			if valType == InvalidType || caseValType == InvalidType {
				continue
			}
			if !TypeEquals(valType, caseValType) {
				w.Report(alerts.NewInvalidCaseType(node.Cases[i].Expressions[j].GetToken().Span, valType.String(), caseValType.String()))
			}
		}
	}
	if !node.HasDefault {
		prevPathTag.SetAllFalse()
	}
	w.reportExits(&prevPathTag, scope)
}

func (w *Walker) breakStatement(node *ast.BreakStmt, scope *Scope) {
	if !scope.Is(BreakAllowing) {
		w.Report(alerts.NewInvalidUseOfExitStmt(node.Token.Span, "break", "for loops or match statement"))
	}

	if returnable := scope.resolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, ControlFlow)
	}
}

func (w *Walker) continueStatement(node *ast.ContinueStmt, scope *Scope) {
	if !scope.Is(ContinueAllowing) {
		w.Report(alerts.NewInvalidUseOfExitStmt(node.Token.Span, "continue", "for loops"))
	}

	if returnable := scope.resolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, ControlFlow)
	}
}

func (w *Walker) returnStatement(node *ast.ReturnStmt, scope *Scope) *[]Type {
	if !scope.Is(ReturnAllowing) {
		w.Report(alerts.NewInvalidUseOfExitStmt(node.Token.Span, "return", "a function or method"))
	}

	ret := make(Values2, 0)
	for i := range node.Args {
		val := w.GetActualNodeValue(&node.Args[i], scope)
		if vls, ok := val.(Values); ok {
			for j := range vls {
				ret = append(ret, Value2{vls[j], i})
			}
		} else {
			ret = append(ret, Value2{val, i})
		}
	}
	sc, funcTag := resolveTagScope[*FuncTag](scope)
	if sc == nil {
		return ret.Types()
	}

	w.validateReturnValues(node.Args, ret, (*funcTag).ReturnTypes, node.Token, "in return arguments") // wait

	if returnable := scope.resolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Return)
		(*returnable).SetExit(true, ControlFlow)
	}

	return ret.Types()
}

func (w *Walker) yieldStatement(node *ast.YieldStmt, scope *Scope) *[]Type {
	if !scope.Is(YieldAllowing) {
		w.Report(alerts.NewInvalidUseOfExitStmt(node.Token.Span, "yield", "a match expression"))
	}

	ret := make(Values2, 0)
	for i := range node.Args {
		val := w.GetActualNodeValue(&node.Args[i], scope)
		if vls, ok := val.(Values); ok {
			for i := range vls {
				ret = append(ret, Value2{vls[i], i})
			}
		} else {
			ret = append(ret, Value2{val, i})
		}
	}
	sc, matchExprT := resolveTagScope[*MatchExprTag](scope)

	if sc == nil {
		return ret.Types()
	}

	matchExprTag := *matchExprT

	if slices.Equal(matchExprTag.YieldTypes, EmptyReturn) {
		matchExprTag.YieldTypes = *ret.Types()
	} else {
		w.validateReturnValues(node.Args, ret, matchExprTag.YieldTypes, node.Token, "in yield arguments")
	}

	if returnable := scope.resolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Yield)
	}

	return ret.Types()
}

func (w *Walker) useStatement(node *ast.UseStmt, scope *Scope) {
	if scope.Parent != nil {
		w.Report(alerts.NewInvalidStmtInLocalBlock(node.Token.Span, "use statement"))
		return
	}

	envName := node.PathExpr.Path.Lexeme

	switch envName {
	case "Pewpew":
		if w.environment.Type != ast.LevelEnv {
			w.Report(alerts.NewUnallowedLibraryUse(node.PathExpr.Path.Span, "Pewpew", "Mesh/Sound/Shared"))
		}
		if !w.AddLibrary(ast.Pewpew) {
			w.Report(alerts.NewEnvironmentReuse(node.PathExpr.Path.Span, envName))
		}
		w.ImportLibrary(ast.Pewpew)
		w.AddReference("env", envName, node.PathExpr.Path)
		return
	case "Fmath":
		if w.environment.Type != ast.LevelEnv {
			w.Report(alerts.NewUnallowedLibraryUse(node.PathExpr.Path.Span, "Fmath", "Mesh or Sound"))
		}
		if !w.AddLibrary(ast.Fmath) {
			w.Report(alerts.NewEnvironmentReuse(node.PathExpr.Path.Span, envName))
		}
		w.ImportLibrary(ast.Fmath)
		w.AddReference("env", envName, node.PathExpr.Path)
		return
	case "Math":
		if w.environment.Type == ast.LevelEnv {
			w.Report(alerts.NewUnallowedLibraryUse(node.PathExpr.Path.Span, "Math", "Level"))
		}
		if !w.AddLibrary(ast.Math) {
			w.Report(alerts.NewEnvironmentReuse(node.PathExpr.Path.Span, envName))
		}
		w.ImportLibrary(ast.Math)
		w.AddReference("env", envName, node.PathExpr.Path)
		return
	case "String":
		if !w.AddLibrary(ast.String) {
			w.Report(alerts.NewEnvironmentReuse(node.PathExpr.Path.Span, envName))
		}
		w.ImportLibrary(ast.String)
		w.AddReference("env", envName, node.PathExpr.Path)
		return
	case "Table":
		if !w.AddLibrary(ast.Table) {
			w.Report(alerts.NewEnvironmentReuse(node.PathExpr.Path.Span, envName))
		}
		w.ImportLibrary(ast.Table)
		w.AddReference("env", envName, node.PathExpr.Path)
		return
	}

	walker, found := w.walkers[envName]
	if !found {
		w.Report(alerts.NewInvalidEnvironmentAccess(node.PathExpr.Path.Span, envName))
		return
	}

	if walker.environment.Name == w.environment.Name {
		w.Report(alerts.NewEnvironmentUsesItself(node.PathExpr.GetToken().Span))
		return
	}

	if walker.environment.Type != ast.SharedEnv && (w.environment.Type == ast.MeshEnv || w.environment.Type == ast.SoundEnv) {
		w.Report(alerts.NewUnallowedEnvironmentAccess(node.PathExpr.GetToken().Span, "non Shared", "Mesh or Sound"))
	} else if w.environment.Type == ast.LevelEnv && (walker.environment.Type == ast.MeshEnv || walker.environment.Type == ast.SoundEnv) {
		w.Report(alerts.NewUnallowedEnvironmentAccess(node.PathExpr.GetToken().Span, "Mesh or Sound", "Level"))
	}

	if paths, isCycle := w.ResolveImportCycle(walker); isCycle {
		paths = append([]string{w.environment.hybroidPath}, paths...)
		w.Report(alerts.NewImportCycle(node.PathExpr.Path.Span, paths))
		return
	}

	success := w.environment.AddRequirement(walker.environment.luaPath)

	if !success {
		w.Report(alerts.NewEnvironmentReuse(node.PathExpr.Path.Span, envName))
		return
	}
	w.environment.imports = append(w.environment.imports, Import{
		Walker:     walker,
		ThroughUse: true,
	})
	w.AddReference("env", envName, node.PathExpr.Path)

	if walker.environment.luaPath == "/dynamic/level.lua" {
		return // we don't put level.hyb in requirements as that would break things
	}

	if !walker.Walked {
		walker.Walk()
	}
}

func (w *Walker) destroyStatement(node *ast.DestroyStmt, scope *Scope) {
	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		return
	}
	if valType.PVT() != ast.Entity {
		w.Report(alerts.NewTypeMismatch(node.Identifier.GetToken().Span, "entity", valType.String(), "in destroy statement"))
		return
	}

	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}

	entityVal := val.(*EntityVal)

	node.EnvName = entityVal.Type.EnvName
	node.EntityName = entityVal.Type.Name

	args := make([]Value, 0)
	for i := range node.Args {
		args = append(args, w.GetActualNodeValue(&node.Args[i], scope))
	}

	suppliedGenerics := w.getGenerics(node.GenericsArgs, entityVal.Destroy.Generics, scope)
	w.validateArguments(suppliedGenerics, args, entityVal.Destroy, node)
}
