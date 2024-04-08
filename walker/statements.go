package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) ifStmt(node ast.IfStmt, scope *Scope) {
	
	for _, node := range node.Body {
		ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		w.WalkNode(node, &ifScope)
		// if stmt.GetType() == ast.ReturnStatement {
		// 	returnStmt := stmt.(ast.ReturnStmt)
		// 	for _, arg := range returnStmt.Args {
		// 		value := w.GetNodeValue(arg, scope)
		// 	}
		// }
	}

	for _, elseif := range node.Elseifs {
		ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		for _, stmt := range elseif.Body {
			w.WalkNode(stmt, &ifScope)
			// if stmt.GetType() == ast.ReturnStatement {
			// 	returnStmt := stmt.(ast.ReturnStmt)
			// 	for _, arg := range returnStmt.Args {
			// 		value := w.GetNodeValue(arg, scope)
			// 	}
			// }
		}
	}

	if node.Else != nil {
		ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		for _, stmt := range node.Else.Body {
			w.WalkNode(stmt, &ifScope)
		}
	}
}

func (w *Walker) assignmentStmt(assignStmt ast.AssignmentStmt, scope *Scope) {
	//if node.Expression.NodeType != parser.Identifier {
	//	w.error(node.Expression.Token, "expected an identifier to assign to")
	//}

	hasFuncs := false

	wIdents := []Value{}
	for _, ident := range assignStmt.Identifiers {
		ident := w.GetNodeValue(ident, scope)
		wIdents = append(wIdents, ident)
	}

	for i, rightValue := range assignStmt.Values {
		if rightValue.GetType() == ast.CallExpression {
			hasFuncs = true
		}
		value := w.GetNodeValue(rightValue, scope)
		if i > len(wIdents)-1 {
			break
		}
		if assignStmt.Identifiers[i].GetType() != ast.MemberExpression {
			variable, ok := wIdents[i].(VariableVal)
			if ok {
				if _, err := scope.AssignVariable(variable.Name, value); err != nil {
					err.Token = variable.Node.GetToken()
					w.addError(*err)
				}
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

func (w *Walker) functionDeclarationStmt(node ast.FunctionDeclarationStmt, scope *Scope) {
	fnScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

	for i, param := range node.Params {
		value := w.GetValue(w.GetTypeFromString(node.Params[i].Type.Lexeme))
		fnScope.DeclareVariable(VariableVal{Name: param.Name.Lexeme, Value: value, Node: node})
	}

	var ret ReturnType
	for _, token := range node.Return {
		ret.values = append(ret.values, w.GetTypeFromString(token.Lexeme))
	}
	if len(ret.values) == 0 {
		ret.values = append(ret.values, ast.Nil)
	}

	variable := VariableVal{ 
		Name:  node.Name.Lexeme, 
		Value: FunctionVal{params: node.Params, returnVal: ret},
		Node:  node,
	}
	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Name, "cannot redeclare a function")
	}

	if scope.Parent != nil && !node.IsLocal {
		w.error(node.GetToken(), "cannot declare a global function inside a local block")
	}
	

	for _, stmt := range node.Body {
		if stmt.GetType() == ast.ReturnStatement {
			returnStmt := stmt.(ast.ReturnStmt)
			returnValue := ReturnType{}
			for _, arg := range returnStmt.Args {
				value := w.GetNodeValue(arg, &fnScope)
				returnValue.values = append(returnValue.values, value.GetType())
			}
		}
	}

	if w.bodyReturns(node.Body, &ret, &fnScope) == nil && ret.values[0] != ast.Nil {
		w.error(node.GetToken(), "not all function paths return a value")
	}
}

func (w *Walker) validateReturnValues(node ast.Node, returnValues []ast.PrimitiveValueType, expectedReturnValues []ast.PrimitiveValueType) {
	if !listsAreValid(returnValues, expectedReturnValues) {
		w.error(node.GetToken(), "invalid return type(s)")
	}
	if len(returnValues) < len(expectedReturnValues) {
		w.error(node.GetToken(), "not enough return values given")
	} else if len(returnValues) > len(expectedReturnValues) {
		w.error(node.GetToken(), "too many return values given")
	}
}

func (w *Walker) ifReturns(node ast.IfStmt, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	var returns *ReturnType

	for _, bodynode := range node.Body {
		localScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		switch bodynode.GetType() {
		case ast.IfStatement:
			returns = w.ifReturns(bodynode.(ast.IfStmt), expectedReturn, &localScope)
		case ast.RepeatStatement:
			returns = w.bodyReturns(bodynode.(ast.RepeatStmt).Body, expectedReturn, &localScope)
		case ast.ReturnStatement:
			returns = w.returnStmt(bodynode.(ast.ReturnStmt), scope)
		}
		if returns == nil {
			continue
		} else if bodynode.GetToken() != node.Body[len(node.Body)-1].GetToken() {
			w.error(bodynode.GetToken(), "unreachable code detected")
		}
		
		w.validateReturnValues(node, returns.values, expectedReturn.values)
	}
	if returns == nil {
		return returns
	}

	for _, elseif := range node.Elseifs {
		localScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		for _, node := range elseif.Body {
			switch node.GetType() {
			case ast.IfStatement:
				returns = w.ifReturns(node.(ast.IfStmt), expectedReturn, &localScope)
			case ast.RepeatStatement:
				returns = w.bodyReturns(node.(ast.RepeatStmt).Body, expectedReturn, &localScope)
			case ast.ReturnStatement:
				returns = w.returnStmt(node.(ast.ReturnStmt), scope)
			}
			if returns == nil {
				continue
			} else if node.GetToken() != elseif.Body[len(elseif.Body)-1].GetToken() {
				w.error(node.GetToken(), "unreachable code detected")
			}

			w.validateReturnValues(node, returns.values, expectedReturn.values)
		}
		if returns == nil {
			return returns
		}
	}

	if node.Else != nil {
		localScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
		returns = w.bodyReturns(node.Else.Body, expectedReturn, &localScope)
	}else {
		return nil
	}

	return returns
}

func (w *Walker) bodyReturns(body []ast.Node, expectedReturn *ReturnType, scope *Scope) *ReturnType {
	var returns *ReturnType
	localScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
	for _, node := range body {
		switch node.GetType() {
		case ast.IfStatement:
			fmt.Printf("called\n")
			returns = w.ifReturns(node.(ast.IfStmt), expectedReturn, &localScope)
		case ast.RepeatStatement:
			returns = w.bodyReturns(node.(ast.RepeatStmt).Body, expectedReturn, &localScope)
		case ast.ReturnStatement:
			returns = w.returnStmt(node.(ast.ReturnStmt), scope)
		}	
		if returns == nil {
			continue
		} else if node.GetToken() != body[len(body)-1].GetToken() {
			w.error(node.GetToken(), "unreachable code detected")
		}

		w.validateReturnValues(node, returns.values, expectedReturn.values)
	}
	if returns == nil {
		return returns
	}

	return returns
}

func (w *Walker) returnStmt(node ast.ReturnStmt, scope *Scope) *ReturnType {
	var ret ReturnType
	for _, expr := range node.Args {
		val := w.GetNodeValue(expr, scope)
		ret.values = append(ret.values, val.GetType())
	}
	if len(ret.values) == 0 {
		ret.values = append(ret.values, ast.Nil)
	}
	return &ret
}

func (w *Walker) repeatStmt(node ast.RepeatStmt, scope *Scope) {
	repeatScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

	end := w.GetNodeValue(node.Iterator, scope)
	start := w.GetNodeValue(node.Start, scope)
	skip := w.GetNodeValue(node.Skip, scope)

	if !parser.IsFx(end.GetType()) && end.GetType() != ast.Number {
		w.error(node.Iterator.GetToken(), "invalid value type of iterator")
	}

	repeatType := end.GetType()

	if (repeatType != start.GetType() || start.GetType() == 0) &&
		(repeatType != skip.GetType() || skip.GetType() == 0) {
		w.error(node.Start.GetToken(), "all value types must be the same")
	}

	if node.Variable.GetValueType() != 0 {
		repeatScope.DeclareVariable(VariableVal{Name: node.Variable.Name.Lexeme, Value: w.GetNodeValue(node.Start, scope), Node: node})
	}

	body := node.Body
	for _, stmt := range body {
		w.WalkNode(stmt, &repeatScope)
	}
}

func (w *Walker) tickStmt(node ast.TickStmt, scope *Scope) {
	tickScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

	if node.Variable.GetValueType() != 0 {
		tickScope.DeclareVariable(VariableVal{Name: node.Variable.Name.Lexeme})
	}

	for _, nod := range node.Body {
		w.WalkNode(nod, &tickScope)
	}
}

func GetValue(values []Value, index int) Value {
	if index <= len(values)-1 {
		return values[index]
	} else {
		return NilVal{}
	}
}

func (w *Walker) variableDeclarationStmt(declaration ast.VariableDeclarationStmt, scope *Scope) {
	var values []Value

	hasFuncs := false
	for i, expr := range declaration.Values {
		if expr.GetType() == ast.CallExpression {
			hasFuncs = true
		}
		exprValue := w.GetNodeValue(expr, scope)
		if declaration.Types[i] != nil {
			typee := w.GetTypeFromString(declaration.Types[i].Name.Lexeme)
			exprValueType := exprValue.GetType()
			if typee != exprValueType {
				w.error(expr.GetToken(), fmt.Sprintf("mismatched types: a type '%s' is given, but a value of type '%s' is assigned", typee.ToString(), exprValueType.ToString()))
			}
			var valueType ast.PrimitiveValueType
			value := ""
			switch val := exprValue.(type) {
			case MapVal:
				valueType = val.GetMemberType()
				value = "map"
			case ListVal:
				valueType = val.GetValuesType()
				value = "list"
			}
			if value != "" {
				if declaration.Types[i].WrappedType == nil {
					w.error(declaration.Types[i].GetToken(), value+"s require a wrapped type to be given")
				} else {

					wrappedType := w.GetTypeFromString(declaration.Types[i].WrappedType.Name.Lexeme)
					if wrappedType == ast.Undefined {
						w.error(declaration.Types[i].WrappedType.Name, "wrapped type given is undefined")
					} else if valueType != wrappedType {
						w.error(expr.GetToken(), value+" contents must the same type as the wrapped type given")
					}
				}
			}

		} else {
			var valueType ast.PrimitiveValueType
			switch val := exprValue.(type) {
			case MapVal:
				valueType = val.GetMemberType()
			case ListVal:
				valueType = val.GetValuesType()
			}

			if valueType == ast.Undefined {

			}
		}
		values = append(values, exprValue)
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
		variable := VariableVal{
			Value: GetValue(values, i),
			Name:  ident,
			Node:  declaration,
		}
		if _, success := scope.DeclareVariable(variable); !success {
			w.error(lexer.Token{Lexeme: ident, Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	if len(values) > len(declaration.Identifiers) {
		w.error(declaration.Token, "too many values provided in declaration")
	} else if len(values) < len(declaration.Identifiers) && !hasFuncs && !isLocal {
		w.error(declaration.Token, "too few values provided in declaration")
	}
}

func (w *Walker) useStmt(node ast.UseStmt, scope *Scope) {
	variable := VariableVal{Name: node.Variable.Name.Lexeme, Value: NamespaceVal{Name: node.Variable.Name.Lexeme}, Node: node}

	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Variable.Name, "cannot declare a value in the same scope twice")
	}
}
