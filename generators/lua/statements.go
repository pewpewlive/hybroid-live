package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (gen *Generator) ifStmt(node ast.IfStmt) string {
	src := StringBuilder{}
	ifTabs := getTabs()

	TabsCount += 1

	src.Append(ifTabs, "if ", gen.GenerateNode(node.BoolExpr), " then\n")

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

	TabsCount -= 1

	return src.String()
}

func (gen *Generator) assignmentStmt(assginStmt ast.AssignmentStmt) string {
	src := StringBuilder{}
	src.AppendTabbed()

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
	fnTabs := getTabs()

	TabsCount += 1

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

	TabsCount -= 1

	return src.String()
}

func (gen *Generator) returnStmt(node ast.ReturnStmt) string {
	src := StringBuilder{}

	src.AppendTabbed("return ")
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
	repeatTabs := getTabs()

	TabsCount += 1
	tabs := getTabs()

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

	TabsCount -= 1

	return src.String()
}

func (gen *Generator) tickStmt(node ast.TickStmt) string {
	src := StringBuilder{}
	tickTabs := getTabs()

	TabsCount += 1

	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateNode(node.Variable)
		src.Append(tickTabs, "local ", variable, " = 0\n")
		src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
		src.AppendTabbed(variable, " = ", variable, " + 1\n")
	} else {
		src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
	}

	for _, stmt := range node.Body {
		value := gen.GenerateNode(stmt)
		src.AppendTabbed(value, "\n")
	}

	src.Append(tickTabs, "end)")

	TabsCount -= 1

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
		src.AppendTabbed("local ")
	} else {
		src.AppendTabbed("")
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

	for _, nodebody := range *node.Methods {
		src.WriteString(gen.methodDeclarationStmt(nodebody, node))
	}

	src.WriteString(gen.constructorDeclarationStmt(*node.Constructor, node))

	return src.String()
}

func (gen *Generator) constructorDeclarationStmt(node ast.ConstructorStmt, Struct ast.StructDeclarationStmt) string {
	src := StringBuilder{}

	if Struct.IsLocal {
		src.WriteString("local ")
	}

	src.Append("function Hybroid_", Struct.Name.Lexeme, "_New(")

	TabsCount += 1

	for i, param := range node.Params {
		src.Append(param.Name.Lexeme)
		if i != len(node.Params)-1 {
			src.Append(", ")
		}
	}
	src.Append(")\n")

	src.AppendTabbed("local Self = {")
	for i, field := range Struct.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(Struct.Fields)-1 {
				src.WriteString(gen.GenerateNode(value))
			} else {
				src.Append(gen.GenerateNode(value), ",")
			}
		}
	}
	src.WriteString("}\n")

	src.WriteString(gen.GenerateString(*node.Body))
	src.AppendTabbed("return Self\n")
	TabsCount -= 1
	src.AppendTabbed("end\n")
	return src.String()
}

func (gen *Generator) methodDeclarationStmt(node ast.MethodDeclarationStmt, Struct ast.StructDeclarationStmt) string {
	src := StringBuilder{}

	if Struct.IsLocal {
		src.WriteString("local ")
	}

	src.Append("function Hybroid_", Struct.Name.Lexeme, "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		src.WriteString(", ")
		src.Append(param.Name.Lexeme)
	}
	src.Append(")\n")

	TabsCount += 1

	src.Append(gen.GenerateString(node.Body))

	TabsCount -= 1

	src.AppendTabbed("end\n")

	return src.String()
}

func (gen *Generator) useStmt(node ast.UseStmt) string {
	src := StringBuilder{}

	fileName := strings.Replace(node.File.Literal, ".hyb", ".lua", 1)
	src.Append("local ", node.Variable.Name.Lexeme, " = require(\"/dynamic/", fileName, "\")")

	return src.String()
}

func (gen *Generator) matchStmt(node ast.MatchStmt) string {
	src := StringBuilder{}
	ifStmt := ast.IfStmt{
		BoolExpr: ast.BinaryExpr{Left: node.ExprToMatch, Operator: lexer.Token{Type:lexer.EqualEqual, Lexeme: "=="}, Right:node.Cases[0].Expression},
		Body: node.Cases[0].Body,
	}
	for i := range node.Cases {
		if i == 0 || i == len(node.Cases)-1 {
			continue
		}
		elseIfStmt := ast.IfStmt{
			BoolExpr: ast.BinaryExpr{Left: node.ExprToMatch, Operator: lexer.Token{Type:lexer.EqualEqual, Lexeme: "=="}, Right:node.Cases[i].Expression},
			Body: node.Cases[i].Body,
		}
		ifStmt.Elseifs = append(ifStmt.Elseifs, &elseIfStmt)
	}

	ifStmt.Else = &ast.IfStmt{
		Body: node.Cases[len(node.Cases)-1].Body,
	}

	src.WriteString(gen.ifStmt(ifStmt))

	return src.String()
}