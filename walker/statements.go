package walker

import (
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) ifStmt(node *ast.IfStmt, scope *Scope) {
	ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
	w.GetNodeValue(&node.BoolExpr,scope)
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
		w.GetNodeValue(&elseif.BoolExpr,scope)
		ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
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
		ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}
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

		if wIdents[i].GetType() == ast.Undefined {
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

	for _, node := range node.Body {
		w.WalkNode(&node, &fnScope)
	}

	if w.bodyReturns(&node.Body, &ret, &fnScope) == nil && ret.values[0] != ast.Nil {
		w.error(node.GetToken(), "not all function paths return a value")
	}
}

func (w *Walker) returnStmt(node *ast.ReturnStmt, scope *Scope) *ReturnType {
	var ret ReturnType
	for _, expr := range node.Args {
		val := w.GetNodeValue(&expr, scope)
		ret.values = append(ret.values, val.GetType())
	}
	if len(ret.values) == 0 {
		ret.values = append(ret.values, ast.Nil)
	}
	return &ret
}

func (w *Walker) repeatStmt(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

	end := w.GetNodeValue(&node.Iterator, scope)
	start := w.GetNodeValue(&node.Start, scope)
	skip := w.GetNodeValue(&node.Skip, scope)

	if !parser.IsFx(end.GetType()) && end.GetType() != ast.Number {
		w.error(node.Iterator.GetToken(), "invalid value type of iterator")
	}

	repeatType := end.GetType()

	if (repeatType != start.GetType() || start.GetType() == 0) &&
		(repeatType != skip.GetType() || skip.GetType() == 0) {
		w.error(node.Start.GetToken(), "all value types must be the same")
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
	tickScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

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
		return NilVal{}
	}
}

func (w *Walker) variableDeclarationStmt(declaration *ast.VariableDeclarationStmt, scope *Scope) {
	var values []Value

	hasFuncs := false
	for _, expr := range declaration.Values {
		if expr.GetType() == ast.CallExpression {
			hasFuncs = true
		}
		exprValue := w.GetNodeValue(&expr, scope)

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

func (w *Walker) useStmt(node *ast.UseStmt, scope *Scope) {
	variable := VariableVal{Name: node.Variable.Name.Lexeme, Value: NamespaceVal{Name: node.Variable.Name.Lexeme}, Node: node}

	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Variable.Name, "cannot declare a value in the same scope twice")
	}
}
