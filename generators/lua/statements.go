package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (gen *Generator) ifStmt(node ast.IfStmt) string {
	src := StringBuilder{}
	ifTabs := gen.getTabs()

	gen.TabsCount += 1

	src.Append("if ", gen.GenerateNode(node.BoolExpr), " then\n")

	src.WriteString(gen.GenerateString(node.Body))
	for _, elseif := range node.Elseifs {
		src.Append(ifTabs, "elseif ", gen.GenerateNode(elseif.BoolExpr), " then\n")
		src.WriteString(gen.GenerateString(elseif.Body))
	}
	if node.Else != nil {
		src.Append(ifTabs, "else \n")
		src.WriteString(gen.GenerateString(node.Else.Body))
	}

	src.Append(ifTabs, "end")

	gen.TabsCount -= 1

	return src.String()
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
	src := StringBuilder{}
	fnTabs := gen.getTabs()

	gen.TabsCount += 1

	if node.IsLocal {
		src.Append(fnTabs, "local ")
	} else {
		src.Append(fnTabs)
	}

	src.Append("function ", node.Name.Lexeme, "(")
	for i, param := range node.Params {
		src.Append(param.Name.Lexeme)
		if i != len(node.Params)-1 {
			src.Append(", ")
		}
	}
	src.Append(")\n")

	src.WriteString(gen.GenerateString(node.Body))

	src.Append(fnTabs + "end")

	gen.TabsCount -= 1

	return src.String()
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
	src := StringBuilder{}
	repeatTabs := gen.getTabs()

	gen.TabsCount += 1
	tabs := gen.getTabs()

	end := gen.GenerateNode(node.Iterator)
	start := gen.GenerateNode(node.Start)
	skip := gen.GenerateNode(node.Skip)
	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateNode(node.Variable)
		src.Append(repeatTabs, "for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		src.Append(repeatTabs, "for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	for _, stmt := range node.Body {
		value := gen.GenerateNode(stmt)
		src.Append(tabs, value, "\n")
	}

	src.Append(repeatTabs, "end")

	gen.TabsCount -= 1

	return src.String()
}

func (gen *Generator) tickStmt(node ast.TickStmt) string {
	src := StringBuilder{}
	tickTabs := gen.getTabs()

	gen.TabsCount += 1
	tabs := gen.getTabs()

	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateNode(node.Variable)
		src.Append(tickTabs, "local ", variable, " = 0\n")
		src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
		src.Append(tabs, variable, " = ", variable, " + 1\n")
	} else {
		src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
	}

	for _, stmt := range node.Body {
		value := gen.GenerateNode(stmt)
		src.Append(tabs, value, "\n")
	}

	src.Append(tickTabs, "end)")

	gen.TabsCount -= 1

	return src.String()
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
			src.WriteString(fmt.Sprintf("%s = ", ident.Lexeme))
		} else if i == len(declaration.Identifiers)-1 {
			src.WriteString(ident.Lexeme)
		} else {
			src.WriteString(fmt.Sprintf("%s, ", ident.Lexeme))
		}
	}
	for i := range values {
		if i == len(values)-1 {
			src2.WriteString(values[i])
			break
		}
		src2.WriteString(fmt.Sprintf("%s, ", values[i]))

	}

	src.Append(src2.String(), "\n")

	return src.String()
}

func (gen *Generator) structDeclarationStmt(node ast.StructDeclarationStmt) string {
	src := StringBuilder{}

	structTabs := gen.getTabs()
	gen.TabsCount += 1;

	if node.IsLocal {
		src.Append(structTabs, "local ")
	}

	src.Append(structTabs, node.Name.Lexeme, " = {}\n")

	for _, nodebody := range *node.Body {
		switch newNode := nodebody.(type)  {
		case ast.FieldDeclarationStmt:
			src.WriteString(gen.fieldDeclarationStmt(newNode, node.Name.Lexeme))
		case ast.MethodDeclarationStmt:
			src.WriteString(gen.methodDeclarationStmt(newNode, node.Name.Lexeme))
		}
	}

	gen.TabsCount -= 1

	return src.String()
}

func (gen *Generator) fieldDeclarationStmt(node ast.FieldDeclarationStmt, structName string) string {
	src := StringBuilder{}

	var values []string
	for _, expr := range node.Values {
		values = append(values, gen.GenerateNode(expr))
	}

	for i, ident := range node.Identifiers {
		if i == len(node.Identifiers)-1 && len(values) != 0 {
			src.WriteString(fmt.Sprintf("%s.%s = ", structName, ident.Lexeme))
		} else if i == len(node.Identifiers)-1 {
			src.Append(structName, ".", ident.Lexeme)
		} else {
			src.WriteString(fmt.Sprintf("%s.%s, ", structName, ident.Lexeme))
		}
	}
	src2 := StringBuilder{}
	for i := range values {
		if i == len(values)-1 {
			src2.WriteString(values[i])
			break
		}
		src2.WriteString(fmt.Sprintf("%s, ", values[i]))

	}

	src.Append(src2.String(), "\n")

	return src.String()
}

func (gen *Generator) methodDeclarationStmt(node ast.MethodDeclarationStmt, structName string) string {
	src := StringBuilder{}

	gen.TabsCount += 1

	src.Append("function ", structName, ".", node.Name.Lexeme, "(")
	for i, param := range node.Params {
		src.Append(param.Name.Lexeme)
		if i != len(node.Params)-1 {
			src.Append(", ")
		}
	}
	src.Append(")\n")

	src.Append(gen.GenerateString(node.Body), "end")

	gen.TabsCount -= 1

	return src.String()
}

func (gen *Generator) useStmt(node ast.UseStmt) string {
	src := StringBuilder{}

	fileName := strings.Replace(node.File.Literal, ".hyb", ".lua", 1)
	src.Append("local ", node.Variable.Name.Lexeme, " = require(\"/dynamic/", fileName, "\")")

	return src.String()
}
