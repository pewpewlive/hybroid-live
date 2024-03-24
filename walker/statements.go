package walker

import (
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) ifStmt(node ast.IfStmt, scope *Scope) {
	ifScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

	for _, stmt := range node.Body {
		w.WalkNode(stmt, &ifScope)
		// if stmt.GetType() == ast.ReturnStatement {
		// 	returnStmt := stmt.(ast.ReturnStmt)
		// 	for _, arg := range returnStmt.Args {
		// 		value := w.GetNodeValue(arg, scope)
		// 	}
		// }
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
			ident := w.GetNodeValue(assignStmt.Identifiers[i], scope)
			if ident.GetType() != ast.Ident {
				w.error(assignStmt.Identifiers[i].GetToken(), "expected an identifier to assign to")
			} else {
				if _, err := scope.AssignVariable(ident.(VariableVal).Name, value); err != nil {
					err.Token = ident.(VariableVal).Node.GetToken()
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
	var returnStmts []ReturnValue
	variable := VariableVal{
		Name:  node.Name.Lexeme,
		Value: CallVal{params: node.Params},
		Node:  node,
	}
	if _, success := scope.DeclareVariable(variable); !success {
		w.error(node.Name, "cannot redeclare a function")
	}

	if scope.Parent != nil && !node.IsLocal {
		w.error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	for _, param := range node.Params {
		fnScope.DeclareVariable(VariableVal{Name: param.Lexeme, Node: node})
	}

	for _, stmt := range node.Body {
		if stmt.GetType() == ast.ReturnStatement {
			returnStmt := stmt.(ast.ReturnStmt)
			returnValue := ReturnValue{}
			for _, arg := range returnStmt.Args {
				value := w.GetNodeValue(arg, &fnScope)
				returnValue.values = append(returnValue.values, value.GetType())
			}
			returnStmts = append(returnStmts, returnValue)
		}
	}

	varr := fnScope.GetVariable(node.Name.Lexeme).Value.(CallVal)
	varr.returnVals = returnStmts
}

func (w *Walker) returnStmt(node ast.ReturnStmt, scope *Scope) {
	// for i, expr := range node.Args {
	// 	val := w.GetNodeValue(expr, scope)
	// }
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
		repeatScope.DeclareVariable(VariableVal{Name: node.Variable.Name, Value: w.GetNodeValue(node.Start, scope), Node: node})
	}

	body := node.Body
	for _, stmt := range body {
		w.WalkNode(stmt, &repeatScope)
	}
}

func (w *Walker) tickStmt(node ast.TickStmt, scope *Scope) {
	tickScope := Scope{Global: scope.Global, Parent: scope, Variables: map[string]VariableVal{}}

	if node.Variable.GetValueType() != 0 {
		tickScope.DeclareVariable(VariableVal{Name: node.Variable.Name})
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
	for _, expr := range declaration.Values {
		if expr.GetType() == ast.CallExpression {
			hasFuncs = true
		}
		values = append(values, w.GetNodeValue(expr, scope))
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
