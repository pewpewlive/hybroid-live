package lua

import (
	"fmt"
	"hybroid/lexer"
	"hybroid/parser"
	"strings"
)

func (gen *Generator) assignmentStmt(node parser.Node, scope *Scope) Value {
	//if node.Expression.NodeType != parser.Identifier {
	//	gen.error(node.Expression.Token, "expected an identifier to assign to")
	//}

	src := strings.Builder{}

	idents := []parser.Node{}

	hasFuncs := false

	switch node.Value.(type) {
	case []parser.Node:
		idents = node.Value.([]parser.Node)
	case parser.Node:
		idents = append(idents, node.Value.(parser.Node))
	}
	genIdents := []Value{}
	for i, ident := range idents {
		ident := gen.GenerateNode(ident, scope)
		genIdents = append(genIdents, ident)
		if i == len(idents)-1 {
			src.WriteString(ident.Val)
		} else {
			src.WriteString(fmt.Sprintf("%s, ", ident.Val))
		}
	}
	src.WriteString(" = ")

	rightValues := []parser.Node{}
	switch node.Value2.(type) {
	case []parser.Node:
		rightValues = node.Value2.([]parser.Node)
	case parser.Node:
		rightValues = append(rightValues, node.Value2.(parser.Node))
	}

	for i, rightValue := range rightValues {
		if rightValue.NodeType == parser.CallExpr {
			hasFuncs = true
		}
		value := gen.GenerateNode(rightValue, scope)// mpathingthing
		if i > len(genIdents)-1 {
			src.WriteString(value.Val)
			break
		}
		if i == len(rightValues)-1 {
			src.WriteString(value.Val)
		} else {
			src.WriteString(fmt.Sprintf("%s, ", value.Val))
		}
		if idents[i].NodeType != parser.MemberExpr {
			if _, success := scope.AssignVariable(genIdents[i].Val, value); !success {
				gen.error(node.Expression.Token, "cannot assign a value to an undeclared variable")
			}
		}
	}

	if len(rightValues) < len(idents) && !hasFuncs {
		gen.error(rightValues[len(rightValues)-1].Token, "not enough values provided in assignment")
	} else if len(rightValues) > len(idents) {
		gen.error(rightValues[len(rightValues)-1].Token, "too many values provided in assignment")
	}

	return Value{parser.Nil, src.String()}
}

func (gen *Generator) functionDeclarationStmt(node parser.Node, scope *Scope) Value {
	fnScope := Scope{Global: scope.Global, Parent: scope, Count: scope.Count + 1, Variables: map[string]Value{}}
	var returnValType parser.PrimitiveValueType
	scope.DeclareVariable(node.Identifier, Value{})

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
		gen.error(node.Token, "cannot declare a global function inside a local block")
	}

	gen.append("function ", node.Identifier, "(")
	params := node.Value.([]lexer.Token)
	for i, param := range params {
		gen.append(param.Lexeme)
		fnScope.DeclareVariable(param.Lexeme, Value{})
		if i != len(params)-1 {
			gen.append(", ")
		}
	}
	gen.append(")\n")

	body := node.Program.Body

	for _, stmt := range body {
		value := gen.GenerateNode(stmt, &fnScope)
		if stmt.NodeType == parser.ReturnStmt {
			returnValType = value.Type
		}
		gen.append(tabs, value.Val, "\n")
	}

	gen.append(fnTabs + "end\n")

	fnScope.AssignVariable(node.Identifier, Value{returnValType, ""})
	return Value{returnValType, ""}
}

func (gen *Generator) returnStmt(node parser.Node, scope *Scope) Value {
	src := strings.Builder{}

	src.WriteString("return ")

	args := node.Value.([]parser.Node)
	for i, expr := range args {
		val := gen.GenerateNode(expr, scope)
		src.WriteString(val.Val)
		if i != len(args)-1 {
			src.WriteString(", ")
		}
	}

	return Value{parser.Nil, src.String()}
}

func (gen *Generator) variableDeclarationStmt(declaration parser.Node, scope *Scope) Value {
	var values []Value

	hasFuncs := false
	if declaration.Value2 == nil {
		gen.error(declaration.Token, "expected expression(s) after declaration")
	} else {
		exprs := declaration.Value2.([]parser.Node)
		for _, expr := range exprs {
			if expr.NodeType == parser.CallExpr {
				hasFuncs = true
			}
			values = append(values,gen.GenerateNode(expr, scope))
		}
	}

	isLocal := declaration.Token.Type == lexer.Let
	src := strings.Builder{}
	src2 := strings.Builder{}
	idents := declaration.Value.([]string)
	if isLocal {
		src.WriteString("local ")
	}else {
		if scope.Parent != nil {
			gen.error(declaration.Token, "cannot declare a global variable inside a local block")
		}
	}
	for i, ident := range idents {
		if i == len(idents)-1 {
			src.WriteString(fmt.Sprintf("%s = ", ident))
		}else {
			src.WriteString(fmt.Sprintf("%s, ", ident))
		}
	}
	for i, value := range values {
		if i > len(idents)-1 {
			src2.WriteString(value.Val)
			break
		}
		if i == len(values)-1{
			src2.WriteString(value.Val)
		}else {
			src2.WriteString(fmt.Sprintf("%s, ", value.Val))
		}
		if _, success := scope.DeclareVariable(idents[i], value); !success {
			gen.error(lexer.Token{Lexeme: declaration.Identifier, Location: declaration.Token.Location},
				"cannot declare a value in the same scope twice")
		}
	}

	if len(values) < len(idents) && !hasFuncs {
		gen.error(declaration.Token, "not enough values provided in declaration")
	}else if len(values) > len(idents) {
		gen.error(declaration.Token, "too many values provided in declaration")
	}

	src.WriteString(src2.String())

	return Value{Type: parser.Nil, Val: src.String()}
}