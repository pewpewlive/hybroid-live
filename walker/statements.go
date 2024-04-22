package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) ifStmt(node *ast.IfStmt, scope *Scope) {
	ifScope := NewScope(scope.Global, scope, scope.Type)
	w.GetNodeValue(&node.BoolExpr, scope)
	for _, node := range node.Body {
		w.WalkNode(&node, &ifScope)
		// if stmt.GetType() == ast.ReturnStatement {
		// 	returnStmt := stmt.(ast.ReturnStmt)
		// 	for _, arg := range returnStmt.Args {
		// 		value := w.GetNodeValue(arg, scope)
		// 	}
		// }
	}

	for _, elseif := range node.Elseifs {
		w.GetNodeValue(&elseif.BoolExpr, scope)
		ifScope := NewScope(scope.Global, scope, scope.Type)
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
		ifScope := NewScope(scope.Global, scope, scope.Type)
		for _, stmt := range node.Else.Body {
			w.WalkNode(&stmt, &ifScope)
		}
	}
}

func (w *Walker) assignmentStmt(assignStmt *ast.AssignmentStmt, scope *Scope) {
	//if node.Expression.NodeType != parser.Identifier {
	//	w.error(node.Expression.Token, "expected an identifier to assign to")
	//}

	hasFuncs := false

	wIdents := []Value{}
	for _, ident := range assignStmt.Identifiers {
		wIdents = append(wIdents, w.GetNodeValue(&ident, scope))
	}

	for i, rightValue := range assignStmt.Values {
		if rightValue.GetType() == ast.CallExpression {
			hasFuncs = true
		}
		value := w.GetNodeValue(&rightValue, scope)
		if i > len(wIdents)-1 {
			break
		}
		if assignStmt.Identifiers[i].GetType() == ast.MemberExpression {
			/*memberType := wIdents[i].(MapMemberVal).Owner.MemberType
			valueType := value.GetType()
			if value.GetType() != 0 && memberType != valueType {
				w.error(rightValue.GetToken(), fmt.Sprintf("map accepts only type of %s but a value of type %s is assigned to its member", memberType.ToString(), valueType.ToString()))
			}*/
			continue
		}

		if wIdents[i].GetType().Type == ast.Undefined {
			w.error(assignStmt.Identifiers[i].GetToken(), "cannot assign a value to an undeclared variable")
			continue
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

func (w *Walker) functionDeclarationStmt(node *ast.FunctionDeclarationStmt, scope *Scope) {
	fnScope := NewScope(scope.Global, scope, ReturnAllowing)

	params := make([]TypeVal, 0)
	for i, param := range node.Params {
		params = append(params, w.typeExpr(&param.Type))
		value := w.GetValueFromType(params[i])
		fnScope.DeclareVariable(VariableVal{Name: param.Name.Lexeme, Value: value, Node: node})
	}

	var ret ReturnType
	for _, typee := range node.Return {
		ret.values = append(ret.values, w.typeExpr(&typee))
	}
	if len(ret.values) == 0 {
		ret.values = append(ret.values, TypeVal{Type: ast.Nil})
	}

	variable := VariableVal{
		Name:  node.Name.Lexeme,
		Value: FunctionVal{params: params, returnVal: ret},
		Node:  node,
	}
	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Name, "cannot redeclare a function")
	}

	if scope.Parent != nil && !node.IsLocal {
		w.error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	for _, node := range node.Body {
		w.WalkNode(&node, &fnScope)
	}

	if w.bodyReturns(&node.Body, &ret, &fnScope) == nil && ret.values[0].Type != ast.Nil {
		w.error(node.GetToken(), "not all function paths return a value")
	}
}

func (w *Walker) returnStmt(node *ast.ReturnStmt, scope *Scope) *ReturnType {
	var ret ReturnType
	if scope.Type == ReturnProhibiting {
		w.error(node.GetToken(), "can't have a return statement outside of a function")
	}
	for _, expr := range node.Args {
		val := w.GetNodeValue(&expr, scope)
		ret.values = append(ret.values, val.GetType())
	}
	if len(ret.values) == 0 {
		ret.values = append(ret.values, TypeVal{Type: ast.Nil})
	}
	return &ret
}

func (w *Walker) repeatStmt(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope.Global, scope, scope.Type)

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
	tickScope := NewScope(scope.Global, scope, scope.Type)

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
		return Undefined{}
	}
}

func (w *Walker) variableDeclarationStmt(declaration *ast.VariableDeclarationStmt, scope *Scope) {
	var values []Value

	for _, expr := range declaration.Values {

		exprValue := w.GetNodeValue(&expr, scope)
		if function, ok := exprValue.(FunctionVal); ok {
			for _, returnVal := range function.returnVal.values {
				values = append(values, w.GetValueFromType(returnVal))
			}
		} else {
			values = append(values, exprValue)
		}
	}

	isLocal := declaration.Token.Type == lexer.Let
	if !isLocal {
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
			if explicitType.Type == ast.Undefined {
				w.error(declaration.Identifiers[i], "uninitialized variable must have its type declared")
			}
			declaration.Values = append(declaration.Values, ast.LiteralExpr{Value: w.GetDefaultValue(explicitType), ValueType: explicitType.Type})
			val = w.GetValueFromType(explicitType)
			valType = val.GetType()
			values = append(values, val)
		}
		if wasMapOrList {
			if declaration.Types[i] == nil {
				if valType.WrappedType.Type == ast.Undefined || valType.WrappedType.Type == 0 {
					w.error(ident, "cannot infer the wrapped type of the map/list: empty or mixed value types")
				}
			} else if declaration.Types[i].WrappedType == nil {
				w.error(declaration.Types[i].GetToken(), "expected a wrapped type in map/list declaration")
			}else if valType.WrappedType.Type != 0 && !valType.Eq(explicitType) {
				w.error(ident, fmt.Sprintf("given value for '%s' does not match with the type given", ident.Lexeme))
			}

		}else if valType.Type != 0 && explicitType.Type != 0 && !valType.Eq(explicitType) {
			w.error(ident, fmt.Sprintf("given value for '%s' does not match with the type given", ident.Lexeme))
		}
		//fmt.Printf("%s\n", valType.Type.ToString())
		//fmt.Printf("%s\n", explicitType.Type.ToString())
		

		if _, success := scope.DeclareVariable(variable); !success {
			w.error(lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	if len(values) == 0 {
		return
	}
	if len(values) > len(declaration.Identifiers) {
		w.error(declaration.Token, "too many values provided in declaration")
	} else if len(values) < len(declaration.Identifiers) {
		w.error(declaration.Token, "too few values provided in declaration")
	}
}

func (w *Walker) useStmt(node *ast.UseStmt, scope *Scope) {
	variable := VariableVal{Name: node.Variable.Name.Lexeme, Value: NamespaceVal{Name: node.Variable.Name.Lexeme}, Node: node}

	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Variable.Name, "cannot declare a value in the same scope twice")
	}
}
