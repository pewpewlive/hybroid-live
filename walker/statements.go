package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) ifStmt(node *ast.IfStmt, scope *Scope) {
	ifScope := NewScope(scope.Global, scope)
	boolExpr := w.GetNodeValue(&node.BoolExpr, scope)
	if boolExpr.GetType().Type != ast.Bool {
		w.error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	for _, node := range node.Body {
		w.Context = node
		w.WalkNode(&node, &ifScope)
		// if stmt.GetType() == ast.ReturnStatement {
		// 	returnStmt := stmt.(ast.ReturnStmt)
		// 	for _, arg := range returnStmt.Args {
		// 		value := w.GetNodeValue(arg, scope)
		// 	}
		// }
	}

	for _, elseif := range node.Elseifs {
		boolExpr := w.GetNodeValue(&elseif.BoolExpr, scope)
		if boolExpr.GetType().Type != ast.Bool {
			w.error(elseif.BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := NewScope(scope.Global, scope)
		for _, stmt := range elseif.Body {
			w.WalkNode(&stmt, &ifScope)
			// if stmt.GetType() == ast.ReturnStatement {
			// 	returnStmt := stmt.(ast.ReturnStmt)
			// 	for _, arg := range returnStmt.Args {
			// 		value := w.GetNodeValue(arg, scope)
			// 	}
			// }
		}
	}

	if node.Else != nil {
		ifScope := NewScope(scope.Global, scope)
		for _, stmt := range node.Else.Body {
			w.WalkNode(&stmt, &ifScope)
		}
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

		if !variableType.Eq(valueType) {
			w.error(assignStmt.Values[i].GetToken(), fmt.Sprintf("mismatched types: variable has a type of %s, but a value of %s was given to it.", variableType.ToString(), valueType.ToString()))
		}

		variable, ok := wIdents[i].(VariableVal)

		if ok {
			if _, err := scope.AssignVariable(variable, value); err != nil {
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
	fnScope := NewScopeWithAttrs(scope.Global, scope, ReturnAllowing)

	params := make([]TypeVal, 0)
	for i, param := range node.Params {
		params = append(params, w.typeExpr(&param.Type))
		value := w.GetValueFromType(params[i])
		fnScope.DeclareVariable(VariableVal{Name: param.Name.Lexeme, Value: value, Node: node})
	}

	ret := ReturnType{
		values: []TypeVal{},
	}
	for _, typee := range node.Return {
		ret.values = append(ret.values, w.typeExpr(&typee))
		//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
	}

	variable := VariableVal{
		Name:  node.Name.Lexeme,
		Value: FunctionVal{params: params, returnVal: ret},
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

	for i := range node.Body {
		w.WalkNode(&node.Body[i], &fnScope)
	}

	if w.bodyReturns(&node.Body, &ret, &fnScope) == nil && len(ret.values) != 0 {
		w.error(node.GetToken(), "not all function paths return a value")
	}

	return variable
}

func (w *Walker) returnStmt(node *ast.ReturnStmt, scope *Scope) *ReturnType {
	var ret ReturnType
	if !scope.Is(ReturnAllowing) { // Structure ReturnAllowing false
		w.error(node.GetToken(), "can't have a return statement outside of a function or method")
	}
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope)
		valType := val.GetType()
		if valType.Type == ast.Func {
			ret.values = append(ret.values, valType.Returns.values...)
		} else {
			ret.values = append(ret.values, valType)
		}
	}
	if len(ret.values) == 0 {
		ret.values = append(ret.values, TypeVal{Type: ast.Nil})
	}
	return &ret
}

func (w *Walker) yieldStmt(node *ast.YieldStmt, scope *Scope) *ReturnType {
	if !scope.Is(YieldAllowing) {
		w.error(node.GetToken(), "cannot use yield outside of ternary operators") // wut
	}

	var ret ReturnType

	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope)
		valType := val.GetType()
		if valType.Type == ast.Func {
			ret.values = append(ret.values, valType.Returns.values...)
		} else {
			ret.values = append(ret.values, valType)
		}
	}

	if len(ret.values) == 0 {
		ret.values = append(ret.values, TypeVal{Type: ast.Nil})
	}

	return &ret
}

func (w *Walker) repeatStmt(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope.Global, scope)

	end := w.GetNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	if !parser.IsFx(endType.Type) && endType.Type != ast.Number {
		w.error(node.Iterator.GetToken(), "invalid value type of iterator")
	} else if variable, ok := end.(VariableVal); ok {
		if fixedpoint, ok := variable.Value.(FixedVal); ok {
			endType = TypeVal{Type: fixedpoint.SpecificType}
		}
	} else {
		if fixedpoint, ok := end.(FixedVal); ok {
			endType = TypeVal{Type: fixedpoint.SpecificType}
		}
	}
	if node.Start.GetType() == ast.NA {
		node.Start = ast.LiteralExpr{Token: node.Start.GetToken(), ValueType: endType.Type, Value: "1"}
	}
	start := w.GetNodeValue(&node.Start, scope)
	if node.Skip.GetType() == ast.NA {
		node.Skip = ast.LiteralExpr{Token: node.Skip.GetToken(), ValueType: endType.Type, Value: "1"}
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

	body := node.Body
	for _, stmt := range body {
		w.WalkNode(&stmt, &repeatScope)
	}
}

func (w *Walker) tickStmt(node *ast.TickStmt, scope *Scope) {
	tickScope := NewScope(scope.Global, scope)

	if node.Variable.GetValueType() != 0 {
		tickScope.DeclareVariable(VariableVal{Name: node.Variable.Name.Lexeme})
	}

	for _, nod := range node.Body {
		w.WalkNode(&nod, &tickScope)
	}
}

func GetValue(values []Value, index int) Value {
	if index <= len(values)-1 {
		return values[index]
	} else {
		return Unknown{}
	}
}

func (w *Walker) GetReturnVals(list *[]Value, ret ReturnType) {
	for _, returnVal := range ret.values {
		val := w.GetValueFromType(returnVal)
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
		if call, ok := exprValue.(CallVal); ok {
			w.GetReturnVals(&values, call.types)
		} else if ret, ok := exprValue.(ReturnType); ok {
			w.GetReturnVals(&values, ret)
		} else {
			values = append(values, exprValue)
		}
	}

	if !declaration.IsLocal {
		if scope.Parent != nil {
			w.error(declaration.Token, "cannot declare a global variable inside a local block")
		}
		if len(values) == 0 {
			w.error(declaration.Token, "cannot declare a global without a value")
		}
	}

	isConstant := declaration.Token.Type == lexer.Const
	if isConstant {
		if scope.Parent != nil {
			w.error(declaration.Token, "cannot declare a global constant inside a local block")
		}
		if len(values) == 0 {
			w.error(declaration.Token, "cannot declare a global constant without a value")
		}
	}
	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}
		val := GetValue(values, i)
		variable := VariableVal{
			Value: val,
			Name:  ident.Lexeme,
			Node:  declaration,
		}

		wasMapOrList := false
		valType := val.GetType()
		explicitType := w.typeExpr(declaration.Types[i])
		if valType.Type == ast.Map || valType.Type == ast.List {
			wasMapOrList = true
		}

		if valType.Type == 0 {
			if explicitType.Type == ast.Invalid {
				w.error(declaration.Identifiers[i], "uninitialized variable must have its type declared")
			} else if explicitType.Type == ast.Func {
				w.error(declaration.Identifiers[i], "cannot declare an uninitialized function")
			}
			declaration.Values = append(declaration.Values, w.GetValueFromType(explicitType).GetDefault())
			val = w.GetValueFromType(explicitType)
			variable.Value = val
			valType = val.GetType()
			values = append(values, val)
		}

		if wasMapOrList {
			if declaration.Types[i] == nil {
				if valType.WrappedType.Type == ast.Invalid || valType.WrappedType.Type == 0 {
					w.error(ident, "cannot infer the wrapped type of the map/list: empty or mixed value types")
				}
			} else if declaration.Types[i].WrappedType == nil {
				w.error(declaration.Types[i].GetToken(), "expected a wrapped type in map/list declaration")
			} else if valType.WrappedType.Type != 0 && !valType.Eq(explicitType) {
				w.error(ident, fmt.Sprintf("given value for '%s' does not match with the type given, (explicit:%s, inferred:%s)", ident.Lexeme, explicitType.ToString(), valType.ToString()))
			}
		} else if valType.Type != 0 && explicitType.Type != ast.Invalid && !valType.Eq(explicitType) {
			w.error(ident, fmt.Sprintf("given value for '%s' does not match with the type given", ident.Lexeme))
		}

		declaredVariables = append(declaredVariables, variable)
		if _, success := scope.DeclareVariable(variable); !success {
			w.error(lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	if len(values) != 0 {
		return declaredVariables
	}
	if len(values) > len(declaration.Identifiers) {
		w.error(declaration.Token, "too many values provided in declaration")
	} else if len(values) < len(declaration.Identifiers) {
		w.error(declaration.Token, "too few values provided in declaration")
	}

	return declaredVariables
}

func (w *Walker) structDeclarationStmt(node *ast.StructDeclarationStmt, scope *Scope) {
	structScope := NewScopeWithAttrs(scope.Global, scope, Structure)

	structTypeVal := StructTypeVal{
		Name:         node.Name,
		Methods:      map[string]VariableVal{},
		Fields:       []VariableVal{},
		FieldIndexes: map[string]int{},
	}
	structScope.WrappedType = structTypeVal.GetType()

	params := make([]TypeVal, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, w.typeExpr(&param.Type))
	}
	structTypeVal.Params = params

	scope.DeclareStructType(&structTypeVal)
	w.Global.foreignTypes[structTypeVal.Name.Lexeme] = &structTypeVal

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
		params := make([]TypeVal, 0)
		for _, param := range (*node.Methods)[i].Params {
			params = append(params, w.typeExpr(&param.Type))
		}

		ret := ReturnType{
			values: []TypeVal{},
		}
		for _, typee := range (*node.Methods)[i].Return {
			ret.values = append(ret.values, w.typeExpr(&typee))
			//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
		}
		variable := VariableVal{
			Name:  (*node.Methods)[i].Name.Lexeme,
			Value: FunctionVal{params: params, returnVal: ret},
			Node:  (*node.Methods)[i],
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

func (w *Walker) fieldDeclarationStmt(node *ast.FieldDeclarationStmt, structTypeVal *StructTypeVal, scope *Scope) {
	varDecl := ast.VariableDeclarationStmt{
		Identifiers: node.Identifiers,
		Types:       node.Types,
		Values:      node.Values,
		IsLocal:     true,
		Token:       node.Token,
	}
	structType := structTypeVal.GetType()
	if len(node.Types) != 0 {
		for i := range node.Types {
			explicitType := w.typeExpr(node.Types[i])
			if explicitType.Eq(structType) {
				w.error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	} else if len(node.Types) != 0 {
		for i := range node.Values {
			valType := w.GetNodeValue(&node.Values[i], scope).GetType()
			if valType.Eq(structType) {
				w.error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	}

	variables := w.variableDeclarationStmt(&varDecl, scope)
	node.Values = varDecl.Values
	structTypeVal.Fields = append(structTypeVal.Fields, variables...)
}

func (w *Walker) methodDeclarationStmt(node *ast.MethodDeclarationStmt, structType *StructTypeVal, scope *Scope) {
	funcExpr := ast.FunctionDeclarationStmt{
		Name:    node.Name,
		Return:  node.Return,
		Params:  node.Params,
		Body:    node.Body,
		IsLocal: true,
	}

	variable := w.functionDeclarationStmt(&funcExpr, scope, Method)
	node.Body = funcExpr.Body
	structType.Methods[variable.Name] = variable
}

func (w *Walker) useStmt(node *ast.UseStmt, scope *Scope) {
	variable := VariableVal{Name: node.Variable.Name.Lexeme, Value: NamespaceVal{Name: node.Variable.Name.Lexeme}, Node: node}

	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Variable.Name, "cannot declare a value in the same scope twice")
	}
}

func (w *Walker) matchStmt(node *ast.MatchStmt, needsDefault bool, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	valType := val.GetType()

	matchScope := NewScope(scope.Global, scope)

	var has_default bool
	for i := range node.Cases {
		for j := range node.Cases[i].Body {
			w.WalkNode(&node.Cases[i].Body[j], &matchScope)
		}
		if node.Cases[i].Expression.GetToken().Lexeme == "_" {
			has_default = true
			continue
		}
		caseValType := w.GetNodeValue(&node.Cases[i].Expression, scope).GetType()
		if !valType.Eq(caseValType) {
			w.error(
				node.Cases[i].Expression.GetToken(),
				fmt.Sprintf("mismatched types: arm expression (%s) and match expression (%s)",
					caseValType.ToString(),
					valType.ToString()))
		}
	}

	if !has_default && needsDefault {
		w.error(node.GetToken(), "match statement has no default arm")
	}
}
