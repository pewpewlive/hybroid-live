package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
)

func (gen *Generator) ifStmt(node ast.IfStmt) string {
	ifTabs := gen.getTabs()

	gen.TabsCount += 1
	tabs := gen.getTabs()

	gen.Src.Append(ifTabs, "if ", gen.GenerateNode(node.BoolExpr), " then\n")

	body := node.Body
	for _, stmt := range body {
		gen.Src.Append(tabs, gen.GenerateNode(stmt), "\n")
	}

	gen.Src.Append(ifTabs, "end\n")

	gen.TabsCount -= 1

	return ""
}

func (gen *Generator) assignmentStmt(assginStmt ast.AssignmentStmt) string {
	src := StringBuilder{}

	for i, ident := range assginStmt.Identifiers {
		ident := gen.GenerateNode(ident)
		if i == len(assginStmt.Identifiers)-1 {
			src.Append(ident)
		} else {
			src.Append(ident, ", ")
		}
	}
	src.Append(" = ")

	for i, rightValue := range assginStmt.Values {
		value := gen.GenerateNode(rightValue)
		if i > len(assginStmt.Identifiers)-1 {
			src.Append(value)
			break
		}
		if i == len(assginStmt.Values)-1 {
			src.Append(value)
		} else {
			src.Append(value, ", ")
		}
	}

	return src.String()
}

func (gen *Generator) functionDeclarationStmt(node ast.FunctionDeclarationStmt) string {
	fnTabs := gen.getTabs()

	gen.TabsCount += 1
	tabs := gen.getTabs()

	if node.IsLocal {
		gen.Src.Append(fnTabs, "local ")
	} else {
		gen.Src.Append(fnTabs)
	}

	gen.Src.Append("function ", node.Name.Lexeme, "(")
	for i, param := range node.Params {
		gen.Src.Append(param.Lexeme)
		if i != len(node.Params)-1 {
			gen.Src.Append(", ")
		}
	}
	gen.Src.Append(")\n")

	for _, stmt := range node.Body {
		gen.Src.Append(tabs, gen.GenerateNode(stmt), "\n")
	}

	gen.Src.Append(fnTabs + "end\n")

	gen.TabsCount -= 1

	return ""
}

func (gen *Generator) returnStmt(node ast.ReturnStmt) string {
	src := StringBuilder{}

	src.Append("return ")
	for i, expr := range node.Args {
		val := gen.GenerateNode(expr)
		src.Append(val)
		if i != len(node.Args)-1 {
			src.Append(", ")
		}
	}

	return src.String()
}

func (gen *Generator) repeatStmt(node ast.RepeatStmt) string {

	repeatTabs := gen.getTabs()

	gen.TabsCount += 1
	tabs := gen.getTabs()

	end := gen.GenerateNode(node.Iterator)
	start := gen.GenerateNode(node.Start)
	skip := gen.GenerateNode(node.Skip)
	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateNode(node.Variable)
		gen.Src.Append(repeatTabs, "for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		gen.Src.Append(repeatTabs, "for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	for _, stmt := range node.Body {
		value := gen.GenerateNode(stmt)
		gen.Src.Append(tabs, value, "\n")
	}

	gen.Src.Append(repeatTabs, "end\n")

	gen.TabsCount -= 1

	return ""
}

func (gen *Generator) tickStmt(node ast.TickStmt) string {

	repeatTabs := gen.getTabs()

	gen.TabsCount += 1
	tabs := gen.getTabs()

	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateNode(node.Variable)
		gen.Src.Append(repeatTabs, "local ", variable, " = 0\n")
		gen.Src.Append(repeatTabs, "pewpew.add_update_callback(function()\n")
		gen.Src.Append(tabs, variable, " = ", variable, " + 1\n")
	} else {
		gen.Src.Append(repeatTabs, "pewpew.add_update_callback(function()\n")
	}

	for _, stmt := range node.Body {
		value := gen.GenerateNode(stmt)
		gen.Src.Append(tabs, value, "\n")
	}

	gen.Src.Append(repeatTabs, "end)\n")

	gen.TabsCount -= 1

	return ""
}

func (gen *Generator) variableDeclarationStmt(declaration ast.VariableDeclarationStmt) string {
	var values []string

	for _, expr := range declaration.Values {
		values = append(values, gen.GenerateNode(expr))
	}

	isLocal := declaration.Token.Type == lexer.Let
	src := StringBuilder{}
	src2 := StringBuilder{}
	if isLocal {
		src.WriteString("local ")
	}
	for i, ident := range declaration.Identifiers {
		if i == len(declaration.Identifiers)-1 && len(values) != 0 {
			src.WriteString(fmt.Sprintf("%s = ", ident))
		} else if i == len(declaration.Identifiers)-1 {
			src.WriteString(ident)
		} else {
			src.WriteString(fmt.Sprintf("%s, ", ident))
		}
	}
	for i := range declaration.Identifiers {
		if i == len(declaration.Identifiers)-1 {
			src2.WriteString(values[i])
			break
		}
		if i == len(declaration.Identifiers)-1 {
			src2.WriteString(values[i])
		} else {
			src2.WriteString(fmt.Sprintf("%s, ", values[i]))
		}
	}

	src.WriteString(src2.String())

	return src.String()
}
