package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
	"strings"
)

func (w *Walker) ifStatement(node *ast.IfStmt, scope *Scope) {
	w.context.EntityCasts.Clear()

	length := len(node.Elseifs) + 2
	mpt := NewMultiPathTag(length, scope.Attributes...)

	w.ifCondition(&node.BoolExpr, scope)
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
		w.ifCondition(&node.Elseifs[i].BoolExpr, scope)
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

	assignOp := assignStmt.AssignOp
	for i := range assignStmt.Identifiers {
		if i >= valuesLen {
			requiredAmount := identsLen - valuesLen
			w.AlertSingle(&alerts.TooFewValuesGiven{}, exprs[len(exprs)-1].GetToken(), requiredAmount, "in assignment")
			return
		}
		value := w.GetNodeValue(&idents[i], scope)
		variable, ok := value.(*VariableVal)
		if !ok {
			if value.GetType() != InvalidType {
				w.AlertSingle(&alerts.InvalidAssignment{}, idents[i].GetToken())
			}
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
		if valType == InvalidType || variableType == InvalidType {
			continue
		}

		if assignOp.Type == tokens.PipeEqual || assignOp.Type == tokens.AmpersandEqual || assignOp.Type == tokens.LeftShiftEqual || assignOp.Type == tokens.RightShiftEqual || assignOp.Type == tokens.TildeEqual {
			if valType.PVT() != ast.Number {
				w.AlertSingle(&alerts.InvalidType{}, exprs[values[i].Index].GetToken(), valType, "in bitwise compound assignment")
				continue
			}
			if variableType.PVT() != ast.Number {
				w.AlertSingle(&alerts.InvalidType{}, idents[i].GetToken(), variableType, "in bitwise compound assignment")
				continue
			}
		} else if assignOp.Type != tokens.Equal {
			if !isNumerical(valType.PVT()) {
				w.AlertSingle(&alerts.InvalidTypeInCompoundAssignment{}, exprs[values[i].Index].GetToken(), valType)
				continue
			}
			if !isNumerical(variableType.PVT()) {
				w.AlertSingle(&alerts.InvalidTypeInCompoundAssignment{}, idents[i].GetToken(), variableType)
				continue
			}
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
}

func (w *Walker) repeatStatement(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, repeatScope.Attributes...)
	repeatScope.Tag = lt

	end := w.GetActualNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	if endType != InvalidType && !isNumerical(endType.PVT()) {
		w.AlertSingle(&alerts.InvalidRepeatIterator{}, node.Iterator.GetToken(), endType)
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

	valType := w.GetActualNodeValue(&node.Iterator, scope).GetType()
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

func (w *Walker) matchStatement(node *ast.MatchStmt, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	valType := val.GetType()
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
	multiPathScope := NewScope(scope, mpt, BreakAllowing)

	for i := range node.Cases {
		caseScope := NewScope(multiPathScope, &UntaggedTag{})

		w.walkBody(&node.Cases[i].Body, mpt, caseScope)

		if node.Cases[i].Expressions[0].GetToken().Lexeme == "else" {
			if i != len(node.Cases)-1 {
				w.AlertSingle(&alerts.InvalidDefaultCasePlacement{}, node.Cases[i].Expressions[0].GetToken(), "in match statement")
			}
			continue
		}

		for j := range node.Cases[i].Expressions {
			caseValType := w.GetNodeValue(&node.Cases[i].Expressions[j], scope).GetType()
			if valType == InvalidType || caseValType == InvalidType {
				continue
			}
			if !TypeEquals(valType, caseValType) {
				w.AlertSingle(&alerts.InvalidCaseType{}, node.Cases[i].Expressions[j].GetToken(), valType, caseValType)
			}
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

	if walker.environment.Name == w.environment.Name {
		w.AlertSingle(&alerts.EnvironmentAccessToItself{}, node.PathExpr.GetToken())
		return
	}

	if walker.environment.Type != ast.SharedEnv && (w.environment.Type == ast.MeshEnv || w.environment.Type == ast.SoundEnv) {
		w.AlertSingle(&alerts.UnallowedEnvironmentAccess{}, node.PathExpr.GetToken(), "non Shared", "Mesh or Sound")
	} else if w.environment.Type == ast.LevelEnv && (walker.environment.Type == ast.MeshEnv || walker.environment.Type == ast.SoundEnv) {
		w.AlertSingle(&alerts.UnallowedEnvironmentAccess{}, node.PathExpr.GetToken(), "Mesh or Sound", "Level")
	}

	if paths, isCycle := w.ResolveImportCycle(walker); isCycle {
		paths = append([]string{w.environment.hybroidPath}, paths...)
		w.AlertSingle(&alerts.ImportCycle{}, node.PathExpr.Path, paths)
		return
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

	suppliedGenerics := w.getGenerics(node.GenericsArgs, entityVal.Destroy.Generics, scope)
	w.validateArguments(suppliedGenerics, args, entityVal.Destroy, node)
}
