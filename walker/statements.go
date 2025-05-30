package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
)

func (w *Walker) ifStatement(node *ast.IfStmt, scope *Scope) {
	w.context.EntityCasts.Clear()

	length := len(node.Elseifs) + 2
	mpt := NewMultiPathTag(length, scope.Attributes...)

	condition := w.GetActualNodeValue(&node.BoolExpr, scope)
	if condition.GetType().PVT() != ast.Bool {
		w.AlertSingle(&alerts.InvalidCondition{}, node.BoolExpr.GetToken(), "in if statement")
	} else if condition.(*BoolVal).Value != "unknown" {
		w.AlertSingle(&alerts.LiteralCondition{}, node.BoolExpr.GetToken(), condition.(*BoolVal).Value)
	}

	multiPathScope := NewScope(scope, mpt)
	ifScope := NewScope(multiPathScope, &UntaggedTag{})

	for w.context.EntityCasts.Count() != 0 {
		cast := w.context.EntityCasts.Pop()
		w.declareVariable(scope, &VariableVal{
			Name:   cast.Name.Lexeme,
			Value:  cast.Entity,
			IsInit: true,
		})
	}
	w.walkBody(&node.Body, mpt, ifScope)

	for i := range node.Elseifs {
		boolExpr := node.Elseifs[i].BoolExpr
		condition := w.GetActualNodeValue(&node.Elseifs[i].BoolExpr, scope)
		if condition.GetType().PVT() != ast.Bool {
			w.AlertSingle(&alerts.InvalidCondition{}, boolExpr.GetToken(), "in if statement")
		} else if condition.(*BoolVal).Value != "unknown" {
			w.AlertSingle(&alerts.LiteralCondition{}, boolExpr.GetToken(), condition.(*BoolVal).Value)
		}
		ifScope := NewScope(multiPathScope, &UntaggedTag{})
		for w.context.EntityCasts.Count() != 0 {
			cast := w.context.EntityCasts.Pop()
			w.declareVariable(scope, &VariableVal{
				Name:   cast.Name.Lexeme,
				Value:  cast.Entity,
				IsInit: true,
			})
		}
		w.walkBody(&node.Elseifs[i].Body, mpt, ifScope)
	}

	if node.Else != nil {
		elseScope := NewScope(multiPathScope, &UntaggedTag{})
		w.walkBody(&node.Else.Body, mpt, elseScope)
	}

	w.reportExits(mpt, scope)
}

// Rewrote
func (w *Walker) assignmentStatement(assignStmt *ast.AssignmentStmt, scope *Scope) {
	values := []Value2{}
	for i := range assignStmt.Values {
		exprValue := w.GetActualNodeValue(&assignStmt.Values[i], scope)
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
					valType.String(),
				)
				continue
			}
			if !isNumerical(variableType.PVT()) {
				w.AlertSingle(&alerts.InvalidTypeInCompoundAssignment{}, idents[i].GetToken(),
					variableType.String(),
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
				variableType.String(),
				valType.String(),
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

	end := w.GetActualNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	if !isNumerical(endType.PVT()) {
		w.AlertSingle(&alerts.InvalidRepeatIterator{}, node.Iterator.GetToken())
	}
	if node.Start == nil {
		node.Start = &ast.LiteralExpr{ValueType: endType.PVT(), Value: "1"}
	}
	start := w.GetNodeValue(&node.Start, scope)
	if node.Skip == nil {
		node.Skip = &ast.LiteralExpr{ValueType: endType.PVT(), Value: "1"}
	}
	skip := w.GetNodeValue(&node.Skip, scope)

	startType := start.GetType()
	skipType := skip.GetType()

	if !(TypeEquals(endType, startType) && TypeEquals(startType, skipType)) {
		w.AlertSingle(&alerts.InconsistentRepeatTypes{}, node.Token,
			startType.String(),
			skipType.String(),
			endType.String(),
		)
	}

	if node.Variable != nil {
		w.declareVariable(repeatScope, NewVariable(node.Variable.Name, end))
	}

	w.walkBody(&node.Body, lt, repeatScope)
	w.reportExits(lt, scope)
}

func (w *Walker) whileStatement(node *ast.WhileStmt, scope *Scope) {
	whileScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, whileScope.Attributes...)
	whileScope.Tag = lt
	w.GetNodeValue(&node.Condition, scope)

	w.walkBody(&node.Body, lt, whileScope)
	w.reportExits(lt, scope)
}

func (w *Walker) forStatement(node *ast.ForStmt, scope *Scope) {
	forScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, forScope.Attributes...)
	forScope.Tag = lt

	valType := w.GetNodeValue(&node.Iterator, scope).GetType()
	wrapper, ok := valType.(*WrapperType)
	if !ok {
		w.AlertSingle(&alerts.InvalidIteratorType{}, node.Iterator.GetToken(), valType.String())
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
			w.AlertSingle(&alerts.UnnecessaryEmptyIdentifier{}, node.Second.Name, "in for loop statement")
		}
		if node.Second.Name.Lexeme != "_" {
			w.declareVariable(forScope, NewVariable(node.Second.Name, w.typeToValue(wrapper.WrappedType)))
		}
	}

	w.walkBody(&node.Body, lt, forScope)
	w.reportExits(lt, scope)
}

func (w *Walker) tickStatement(node *ast.TickStmt, scope *Scope) {
	funcTag := &FuncTag{ReturnTypes: EmptyReturn}
	tickScope := NewScope(scope, funcTag, ReturnAllowing)

	if node.Variable != nil {
		w.declareVariable(tickScope, NewVariable(node.Variable.Name, &NumberVal{}))
	}

	w.walkBody(&node.Body, funcTag, tickScope)
	w.reportExits(funcTag, scope)
}

func (w *Walker) matchStatement(node *ast.MatchStmt, isExpr bool, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	casesLength := len(node.Cases)
	if !node.HasDefault {
		casesLength++
		if casesLength < 1 {
			w.AlertSingle(&alerts.InsufficientCases{}, node.Token)
		}
	} else if casesLength < 2 {
		w.AlertSingle(&alerts.InsufficientCases{}, node.Token)
	}
	mpt := NewMultiPathTag(casesLength, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)

	for i := range node.Cases {
		caseScope := NewScope(multiPathScope, &UntaggedTag{})

		if !isExpr {
			w.walkBody(&node.Cases[i].Body, mpt, caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "else" {
			continue
		}

		caseValType := w.GetNodeValue(&node.Cases[i].Expression, scope).GetType()
		if !TypeEquals(val.GetType(), caseValType) {
			w.AlertSingle(&alerts.InvalidCaseType{}, node.Cases[i].Expression.GetToken(), val.GetType(), caseValType)
		}
	}
	w.reportExits(mpt, scope)
}

func (w *Walker) breakStatement(node *ast.BreakStmt, scope *Scope) {
	if !scope.Is(BreakAllowing) {
		w.AlertSingle(&alerts.InvalidUseOfExitStmt{}, node.Token,
			"break",
			"for loops",
		)
	}

	if returnable := scope.resolveReturnable(); returnable != nil {
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

	if returnable := scope.resolveReturnable(); returnable != nil {
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
	sc, _, funcTag := resolveTagScope[*FuncTag](scope)
	if sc == nil {
		return ret.Types()
	}

	w.validateReturnValues(node.Args, ret, (*funcTag).ReturnTypes, "in return arguments") // wait

	if returnable := scope.resolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Return)
		(*returnable).SetExit(true, All)
	}

	return ret.Types()
}

func (w *Walker) yieldStatement(node *ast.YieldStmt, scope *Scope) *[]Type {
	if !scope.Is(YieldAllowing) {
		w.AlertSingle(&alerts.InvalidUseOfExitStmt{}, node.Token,
			"yield",
			"a match expression",
		)
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
	sc, _, matchExprT := resolveTagScope[*MatchExprTag](scope)

	if sc == nil {
		return ret.Types()
	}

	matchExprTag := *matchExprT

	if core.ListsAreSame(matchExprTag.YieldTypes, EmptyReturn) {
		matchExprTag.YieldTypes = *ret.Types()
	} else {
		w.validateReturnValues(node.Args, ret, matchExprTag.YieldTypes, "in yield arguments")
	}

	if returnable := scope.resolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Yield)
		(*returnable).SetExit(true, All)
	}

	return ret.Types()
}

func (w *Walker) useStatement(node *ast.UseStmt, scope *Scope) {
	if scope.Parent != nil {
		w.AlertSingle(&alerts.InvalidStmtInLocalBlock{}, node.Token, "use statement")
		return
	}

	envName := node.PathExpr.Path.Lexeme

	switch envName {
	case "Pewpew":
		if w.environment.Type != ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Pewpew", "Mesh or Sound")
		}
		w.environment.UsedLibraries = append(w.environment.UsedLibraries, Pewpew)
		return
	case "Fmath":
		if w.environment.Type != ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Fmath", "Mesh or Sound")
		}
		w.environment.UsedLibraries = append(w.environment.UsedLibraries, Fmath)
		return
	case "Math":
		if w.environment.Type == ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Math", "Level")
		}
		w.environment.UsedLibraries = append(w.environment.UsedLibraries, Math)
		return
	case "String":
		w.environment.UsedLibraries = append(w.environment.UsedLibraries, String)
		return
	case "Table":
		w.environment.UsedLibraries = append(w.environment.UsedLibraries, Table)
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
		w.AlertSingle(&alerts.TypeMismatch{}, node.Identifier.GetToken(), "entity", valType.String(), "in destroy statement")
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

	suppliedGenerics := w.getGenerics(node.GenericArgs, entityVal.DestroyGenerics, scope)
	w.validateArguments(suppliedGenerics, args, entityVal.DestroyParams, node)
}
