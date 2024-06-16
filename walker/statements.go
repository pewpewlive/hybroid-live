package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
	"strings"
)

func (w *Walker) ifStmt(node *ast.IfStmt, scope *Scope) {
	length := len(node.Elseifs)+2
	mpt := NewMultiPathTag(length, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)
	ifScope := NewScope(&multiPathScope, &UntaggedTag{})
	boolExpr := w.GetNodeValue(&node.BoolExpr, scope)
	if boolExpr.GetType().PVT() != ast.Bool {
		w.error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	w.WalkBody(&node.Body, mpt, &ifScope)

	for i := range node.Elseifs {
		boolExpr := w.GetNodeValue(&node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().PVT() != ast.Bool {
			w.error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := NewScope(&multiPathScope, &UntaggedTag{})
		w.WalkBody(&node.Elseifs[i].Body, mpt, &ifScope)
	}

	
	if node.Else != nil {
		elseScope := NewScope(&multiPathScope, &UntaggedTag{})
		w.WalkBody(&node.Else.Body, mpt, &elseScope)
	}

	returnabl := scope.ResolveReturnable()

	if returnabl == nil {
		return
	}
	returnable := *returnabl

	returnable.SetExit(mpt.GetIfExits(Return), Return)
	returnable.SetExit(mpt.GetIfExits(Yield), Yield)
	returnable.SetExit(mpt.GetIfExits(Break), Break)
	returnable.SetExit(mpt.GetIfExits(Continue), Continue)
	returnable.SetExit(mpt.GetIfExits(All), All)
}

func (w *Walker) assignment(assignStmt *ast.AssignmentStmt, scope *Scope) {
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
		if variableType.PVT() == ast.Invalid {
			w.error(assignStmt.Identifiers[i].GetToken(), "cannot assign a value to an undeclared variable")
			continue
		}

		if !TypeEquals(variableType, valueType) {
			w.error(assignStmt.Values[i].GetToken(), fmt.Sprintf("mismatched types: variable has a type of %s, but a value of %s was given to it.", variableType.ToString(), valueType.ToString()))
		}

		variable, ok := wIdents[i].(*VariableVal)

		if ok {
			if _, err := scope.AssignVariable(variable, value); err != nil {
				err.Token = variable.Token
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

func (w *Walker) functionDeclaration(node *ast.FunctionDeclarationStmt, scope *Scope, procType ProcedureType) *VariableVal {
	ret := EmptyReturn
	for _, typee := range node.Return {
		ret = append(ret, w.typeExpr(typee))
		//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
	}
	funcTag := &FuncTag{ReturnType: ret}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)

	params := make([]Type, 0)
	for i, param := range node.Params {
		params = append(params, w.typeExpr(param.Type))
		value := w.TypeToValue(params[i])
		fnScope.DeclareVariable(&VariableVal{Name: param.Name.Lexeme, Value: value, Token: node.GetToken()})
	}

	variable := &VariableVal{
		Name:  node.Name.Lexeme,
		Value: &FunctionVal{params: params, returns: ret},
		Token:  node.GetToken(),
	}
	if procType == Function {
		if _, success := scope.DeclareVariable(variable); !success {
			w.error(node.Name, fmt.Sprintf("variable with name '%s' already exists", variable.Name))
		}
	}

	if scope.Parent != nil && !node.IsLocal {
		w.error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	w.WalkBody(&node.Body, funcTag, &fnScope)

	if funcTag, ok := fnScope.Tag.(*FuncTag); ok {
		if !funcTag.GetIfExits(Return) && !ret.Eq(&EmptyReturn) {
			w.error(node.GetToken(), "not all code paths return a value")
		}
	}

	return variable
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
		(*returnable).SetExit(true, All)
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
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) breakStmt(node *ast.BreakStmt, scope *Scope) {
	if !scope.Is(BreakAllowing) {
		w.error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Break)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) continueStmt(node *ast.ContinueStmt, scope *Scope) {
	if !scope.Is(ContinueAllowing) {
		w.error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Continue)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) repeat(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope, &LoopTag{}, BreakAllowing, ContinueAllowing)
	lt := NewLoopTag(repeatScope.Attributes...)
	repeatScope.Tag = lt

	end := w.GetNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	if !parser.IsFx(endType.PVT()) && endType.PVT() != ast.Number {
		w.error(node.Iterator.GetToken(), "invalid value type of iterator")
	} else if variable, ok := end.(*VariableVal); ok {
		if fixedpoint, ok := variable.Value.(*FixedVal); ok {
			endType = NewBasicType(fixedpoint.SpecificType)
		}
	} else {
		if fixedpoint, ok := end.(*FixedVal); ok {
			endType = NewBasicType(fixedpoint.SpecificType)
		}
	}
	if node.Start.GetType() == ast.NA {
		node.Start = &ast.LiteralExpr{Token: node.Start.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	start := w.GetNodeValue(&node.Start, scope)
	if node.Skip.GetType() == ast.NA {
		node.Skip = &ast.LiteralExpr{Token: node.Skip.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	skip := w.GetNodeValue(&node.Skip, scope)

	repeatType := end.GetType().PVT()
	startType := start.GetType().PVT()
	skipType := skip.GetType().PVT()

	if (repeatType != startType || startType == 0) &&
		(repeatType != skipType || skipType == 0) {
		w.error(node.Start.GetToken(), fmt.Sprintf("all value types must be the same (iter:%s, start:%s, by:%s)", repeatType.ToString(), startType.ToString(), skipType.ToString()))
	}

	if node.Variable.GetValueType() != 0 {
		repeatScope.
			DeclareVariable(&VariableVal{Name: node.Variable.Name.Lexeme, Value: w.GetNodeValue(&node.Start, scope), Token: node.GetToken()})
	}

	w.WalkBody(&node.Body, lt, &repeatScope)

	w.ReportExits(lt, scope)
}

func (w *Walker) while(node *ast.WhileStmt, scope *Scope) {
	whileScope := NewScope(scope, &LoopTag{}, BreakAllowing, ContinueAllowing)
	lt := NewLoopTag(whileScope.Attributes...)
	whileScope.Tag = lt

	_ = w.GetNodeValue(&node.Condtion, scope)

	w.WalkBody(&node.Body, lt, &whileScope)
}

func (w *Walker) forloop(node *ast.ForStmt, scope *Scope) {
	forScope := NewScope(scope, &LoopTag{}, BreakAllowing, ContinueAllowing)
	lt := NewLoopTag(forScope.Attributes...)
	forScope.Tag = lt

	if len(node.KeyValuePair) != 0 {
		forScope.DeclareVariable(&VariableVal{Name: node.KeyValuePair[0].Name.Lexeme, Value: &NumberVal{}})
	}
	valType := w.GetNodeValue(&node.Iterator, scope).GetType()
	wrapper, ok := valType.(*WrapperType)
	if !ok {
		w.error(node.Iterator.GetToken(), "iterator must be of type map or list")
	}else if len(node.KeyValuePair) == 2 {
		node.OrderedIteration = wrapper.PVT() == ast.List
		forScope.DeclareVariable(&VariableVal{Name: node.KeyValuePair[1].Name.Lexeme, Value: w.TypeToValue(wrapper.WrappedType)})
	}

	w.WalkBody(&node.Body, lt, &forScope)

	w.ReportExits(lt, scope)
}

func (w *Walker) tick(node *ast.TickStmt, scope *Scope) {
	tickScope := NewScope(scope, &UntaggedTag{})

	if node.Variable.GetValueType() != 0 {
		tickScope.DeclareVariable(&VariableVal{Name: node.Variable.Name.Lexeme, Value: &NumberVal{}})
	}

	for i := range node.Body {
		w.WalkNode(&node.Body[i], &tickScope)
	}
}

func (w *Walker) AddTypesToValues(list *[]Value, tys *Types) {
	for _, typ := range *tys {
		val := w.TypeToValue(typ)
		*list = append(*list, val)
	}
}

func (w *Walker) variableDeclaration(declaration *ast.VariableDeclarationStmt, scope *Scope) []*VariableVal {
	declaredVariables := []*VariableVal{}

	idents := len(declaration.Identifiers)
	values := make([]Value, idents)

	for i := range values {
		values[i] = &Invalid{}
	}

	valuesLength := len(declaration.Values)
	if valuesLength > idents {
		w.error(declaration.Token, "too many values provided in declaration")
		return declaredVariables
	}

	for i := range declaration.Values {

		exprValue := w.GetNodeValue(&declaration.Values[i], scope)
		if declaration.Values[i].GetType() == ast.SelfExpression {
			w.error(declaration.Values[i].GetToken(), "cannot assign self to a variable")
		}
		if types, ok := exprValue.(*Types); ok { // i want to make it so that when you dont supply a variable with a value
			temp := values[i:] // it will get the default value if it was given an explicit type
			values = values[:i] // but shit doesnt want to shit
			w.AddTypesToValues(&values, types)
			values = append(values, temp...) // what are they?
		} else {
			values[i] = exprValue
		}
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
			explicitType := w.typeExpr(declaration.Types[i])
			if valType == InvalidType && explicitType != InvalidType {
				values[i] = w.TypeToValue(explicitType)
				declaration.Values = append(declaration.Values, values[i].GetDefault())
			} else if !TypeEquals(valType, explicitType) {
				w.error(declaration.Token, fmt.Sprintf("mismatched types: value type (%s) not the same with explict type (%s)",
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
		if _, success := scope.DeclareVariable(variable); !success {
			w.error(lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	return declaredVariables
}

func (w *Walker) enumDeclarationStmt(node *ast.EnumDeclarationStmt, scope *Scope) {
	enumVal := &EnumVal{
		Type:NewEnumType(node.Name.Lexeme),
	}

	if len(node.Fields) == 0 {
		w.error(node.GetToken(), "can't declare an enum with no fields")
	}
	for _, v := range node.Fields {
		variable := &VariableVal{
			Name: v.Lexeme,
			Value: &EnumFieldVal{Type:enumVal.Type},
			IsConst: true,
		}
		enumVal.AddField(variable)
	}

	enumVar := &VariableVal{
		Name: enumVal.Type.Name,
		Value: enumVal,
		IsConst: true,
	}

	if _, ok := scope.DeclareVariable(enumVar); !ok {
		w.error(node.GetToken(), "cannot declare an enum with the same name as another variable")
	}
}

func (w *Walker) structDeclaration(node *ast.StructDeclarationStmt, scope *Scope) {
	structVal := &StructVal{
		Type: *NewNamedType(node.Name.Lexeme),
		Fields: make([]*VariableVal, 0),
		Methods: map[string]*VariableVal{},
		Params: Types{},
	}

	structScope := NewScope(scope, &StructTag{StructVal: structVal}, SelfAllowing)

	params := make([]Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, w.typeExpr(param.Type))
	}
	structVal.Params = params

	w.DeclareStruct(structVal)

	funcDeclaration := ast.MethodDeclarationStmt{
		Name:    node.Constructor.Token,
		Params:  node.Constructor.Params,
		Return:  node.Constructor.Return,
		IsLocal: true,
		Body:    node.Constructor.Body,
	}

	for i := range node.Fields {
		w.fieldDeclaration(&node.Fields[i], structVal, &structScope)
	}

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
		variable := &VariableVal{
			Name:  (*node.Methods)[i].Name.Lexeme,
			Value: &FunctionVal{params: params, returns: ret},
			Token:  (*node.Methods)[i].GetToken(),
		}
		if _, success := structScope.DeclareVariable(variable); !success {
			w.error((*node.Methods)[i].Name, fmt.Sprintf("variable with name '%s' already exists", variable.Name))
		}
		structVal.Methods[variable.Name] = variable
	}

	for i := range *node.Methods {
		w.methodDeclaration(&(*node.Methods)[i], structVal, &structScope)
	}

	w.methodDeclaration(&funcDeclaration, structVal, &structScope)
}

func (w *Walker) fieldDeclaration(node *ast.FieldDeclarationStmt, container FieldContainer, scope *Scope) {
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
			if TypeEquals(explicitType, structType) {
				w.error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	} else if len(node.Types) != 0 {
		for i := range node.Values {
			valType := w.GetNodeValue(&node.Values[i], scope).GetType()
			if TypeEquals(valType, structType) {
				w.error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	}

	variables := w.variableDeclaration(&varDecl, scope)
	node.Values = varDecl.Values
	for i := range variables {
		container.AddField(variables[i])
	}
}

func (w *Walker) methodDeclaration(node *ast.MethodDeclarationStmt, container MethodContainer, scope *Scope) {
	funcExpr := ast.FunctionDeclarationStmt{
		Name:    node.Name,
		Return:  node.Return,
		Params:  node.Params,
		Body:    node.Body,
		IsLocal: true,
	}

	variable := w.functionDeclaration(&funcExpr, scope, Method)
	node.Body = funcExpr.Body
	container.AddMethod(variable)
}

func (w *Walker) use(node *ast.UseStmt, scope *Scope) {
	variable := &VariableVal{
		Name: node.Variable.Name.Lexeme, 
		Value: &EnvironmentVal{
			Type:&EnvironmentType{
				Name: node.Variable.Name.Lexeme,
			},
		},
		Token: node.GetToken(),
	}

	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Variable.Name, "cannot declare a value in the same scope twice")
	}
}

func (w *Walker) match(node *ast.MatchStmt, isExpr bool, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	if val.GetType().PVT() == ast.Invalid {
		w.error(node.ExprToMatch.GetToken(), "variable is of type invalid")
	}
	valType := val.GetType()
	casesLength := len(node.Cases)+1
	if node.HasDefault {
		casesLength--
	}
	mpt := NewMultiPathTag(casesLength, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)

	var has_default bool
	for i := range node.Cases {
		caseScope := NewScope(&multiPathScope, &UntaggedTag{})

		if !isExpr {
			w.WalkBody(&node.Cases[i].Body, mpt, &caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "_" {
			has_default = true
			continue
		}

		caseValType := w.GetNodeValue(&node.Cases[i].Expression, scope).GetType()
		if !TypeEquals(valType, caseValType) {
			w.error(
				node.Cases[i].Expression.GetToken(),
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

	w.ReportExits(mpt, scope)
}

func (w *Walker) env(node *ast.EnvironmentStmt, scope *Scope) {
	if scope.Environment.Type.Name != "UNKNOWN" {
		w.error(node.GetToken(), "can't have more than one environment statement in a file")
		return
	}

	for i, v := range node.Env.Envs {
		if i < len(node.Env.Envs)-1 {
			scope.Environment.Type.Name += v.Lexeme + "::"
		}else {
			scope.Environment.Type.Name += v.Lexeme
		}
	}
	 
	if wlkr, found := (*w.Walkers)[w.Environment.Type.Name]; found {
		w.error(node.GetToken(), fmt.Sprintf("cannot have two environments with the same name, path: %s",wlkr.Environment.Type.Path))
		return
	}

	(*w.Walkers)[w.Environment.Type.Name] = w
}