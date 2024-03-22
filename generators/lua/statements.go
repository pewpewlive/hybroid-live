package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (gen *Generator) ifStmt(node ast.IfStmt, scope *Scope) Value {
	ifScope := Scope{Global: scope.Global, Parent: scope, Count: scope.Count + 1, Variables: map[string]Value{}}
	var returnValType ast.PrimitiveValueType

	var tabs string
	for i := 0; i < ifScope.Count; i++ {
		tabs += "\t"
	}

	var ifTabs string
	for i := 0; i < scope.Count; i++ {
		ifTabs += "\t"
	}

	gen.append(ifTabs, "if ")

	expr := gen.GenerateNode(node.BoolExpr, scope)
	gen.append(expr.Val, " then\n")

	body := node.Body

	for _, stmt := range body {
		value := gen.GenerateNode(stmt, &ifScope)
		if stmt.GetType() == ast.ReturnStatement {
			returnValType = value.Type
		}
		gen.append(tabs, value.Val, "\n")
	}

	gen.append(ifTabs, "end\n")

	return Value{Type: returnValType, Token: node.Token, Val: ""}
}

func (gen *Generator) assignmentStmt(assginStmt ast.AssignmentStmt, scope *Scope) Value {
	//if node.Expression.NodeType != parser.Identifier {
	//	gen.error(node.Expression.Token, "expected an identifier to assign to")
	//}

	src := strings.Builder{}

	hasFuncs := false

	genIdents := []Value{}
	for i, ident := range assginStmt.Identifiers {
		ident := gen.GenerateNode(ident, scope)
		genIdents = append(genIdents, ident)
		if i == len(assginStmt.Identifiers)-1 {
			src.WriteString(ident.Val)
		} else {
			src.WriteString(fmt.Sprintf("%s, ", ident.Val))
		}
	}
	src.WriteString(" = ")

	for i, rightValue := range assginStmt.Values {
		if rightValue.GetType() == ast.CallExpression {
			hasFuncs = true
		}
		value := gen.GenerateNode(rightValue, scope) // mpathingthing
		if i > len(genIdents)-1 {
			src.WriteString(value.Val)
			break
		}
		if i == len(assginStmt.Values)-1 {
			src.WriteString(value.Val)
		} else {
			src.WriteString(fmt.Sprintf("%s, ", value.Val))
		}
		if assginStmt.Identifiers[i].GetType() != ast.MemberExpression {
			if _, success := scope.AssignVariable(genIdents[i].Val, value); !success {
				gen.error(assginStmt.Token, "cannot assign a value to an undeclared variable")
			}
		}
	}

	if len(assginStmt.Values) < len(assginStmt.Identifiers) && !hasFuncs {
		gen.error(assginStmt.Values[len(assginStmt.Values)-1].GetToken(), "not enough values provided in assignment")
	} else if len(assginStmt.Values) > len(assginStmt.Identifiers) {
		gen.error(assginStmt.Values[len(assginStmt.Values)-1].GetToken(), "too many values provided in assignment")
	}

	return Value{Type: ast.Undefined, Token: assginStmt.Token, Val: src.String()}
}

func (gen *Generator) functionDeclarationStmt(node ast.FunctionDeclarationStmt, scope *Scope) Value {
	fnScope := Scope{Global: scope.Global, Parent: scope, Count: scope.Count + 1, Variables: map[string]Value{}}
	var returnValType ast.PrimitiveValueType
	scope.DeclareVariable(node.Name.Lexeme, Value{})

	var tabs string
	for i := 0; i < fnScope.Count; i++ {
		tabs += "\t"
	}

	var fnTabs string
	for i := 0; i < scope.Count; i++ {
		fnTabs += "\t"
	}

	if node.IsLocal {
		gen.append(fnTabs, "local ")
	} else {
		gen.append(fnTabs)
	}

	if scope.Parent != nil && !node.IsLocal {
		gen.error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	gen.append("function ", node.Name.Lexeme, "(")
	for i, param := range node.Params {
		gen.append(param.Lexeme)
		fnScope.DeclareVariable(param.Lexeme, Value{})
		if i != len(node.Params)-1 {
			gen.append(", ")
		}
	}
	gen.append(")\n")

	for _, stmt := range node.Body {
		value := gen.GenerateNode(stmt, &fnScope)
		if stmt.GetType() == ast.ReturnStatement {
			returnValType = value.Type
		}
		gen.append(tabs, value.Val, "\n")
	}

	gen.append(fnTabs + "end\n")

	fnScope.AssignVariable(node.Name.Lexeme, Value{Type: returnValType, Val: ""})
	return Value{Type: returnValType, Token: node.GetToken(), Val: ""}
}

func (gen *Generator) returnStmt(node ast.ReturnStmt, scope *Scope) Value {
	src := strings.Builder{}

	src.WriteString("return ")
	for i, expr := range node.Args {
		val := gen.GenerateNode(expr, scope)
		src.WriteString(val.Val)
		if i != len(node.Args)-1 {
			src.WriteString(", ")
		}
	}

	// TODO: Make it not undefined
	return Value{Type: ast.Undefined, Token: node.Token, Val: src.String()}
}

func (gen *Generator) variableDeclarationStmt(declaration ast.VariableDeclarationStmt, scope *Scope) Value {
	var values []Value

	hasFuncs := false
	if len(declaration.Values) == 0 {
		gen.error(declaration.Token, "expected expression(s) after declaration")
	} else {
		for _, expr := range declaration.Values {
			if expr.GetType() == ast.CallExpression {
				hasFuncs = true
			}
			values = append(values, gen.GenerateNode(expr, scope))
		}
	}

	isLocal := declaration.Token.Type == lexer.Let
	src := strings.Builder{}
	src2 := strings.Builder{}
	if isLocal {
		src.WriteString("local ")
	} else {
		if scope.Parent != nil {
			gen.error(declaration.Token, "cannot declare a global variable inside a local block")
		}
	}
	for i, ident := range declaration.Identifiers {
		if i == len(declaration.Identifiers)-1 {
			src.WriteString(fmt.Sprintf("%s = ", ident))
		} else {
			src.WriteString(fmt.Sprintf("%s, ", ident))
		}
	}
	for i, value := range values {
		if i > len(declaration.Identifiers)-1 {
			src2.WriteString(value.Val)
			break
		}
		if i == len(values)-1 {
			src2.WriteString(value.Val)
		} else {
			src2.WriteString(fmt.Sprintf("%s, ", value.Val))
		}
		if _, success := scope.DeclareVariable(declaration.Identifiers[i], value); !success {
			gen.error(lexer.Token{Lexeme: declaration.Identifiers[i], Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	if len(values) < len(declaration.Identifiers) && !hasFuncs {
		gen.error(declaration.Token, "not enough values provided in declaration")
	} else if len(values) > len(declaration.Identifiers) {
		gen.error(declaration.Token, "too many values provided in declaration")
	}

	src.WriteString(src2.String())

	return Value{Type: ast.Nil, Token: declaration.Token, Val: src.String()}
}
