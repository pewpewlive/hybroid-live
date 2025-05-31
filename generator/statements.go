package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/tokens"
	"strings"
)

func (gen *Generator) envStmt(node ast.EnvironmentDecl) {
	for i := range node.Requirements {
		gen.Write("require(\"", node.Requirements[i], "\")\n")
	}
}

func (gen *Generator) ifStmt(node ast.IfStmt) {
	gen.WriteTabbed("if ", gen.GenerateExpr(node.BoolExpr), " then\n")

	gen.GenerateBody(node.Body)
	for _, elseif := range node.Elseifs {
		gen.WriteTabbed("elseif ", gen.GenerateExpr(elseif.BoolExpr), " then\n")
		gen.GenerateBody(elseif.Body)
	}
	if node.Else != nil {
		gen.WriteTabbed("else \n")
		gen.GenerateBody(node.Else.Body)
	}

	gen.WriteTabbed("end\n")
}

func (gen *Generator) assignmentStmt(assignStmt ast.AssignmentStmt) {
	src := StringBuilder{}
	preSrc := StringBuilder{}
	vars := []string{}

	index := 0
	for i := range assignStmt.Values {
		src.WriteTabbed()
		if assignStmt.Values[i].GetType() == ast.CallExpression {
			call := assignStmt.Values[i].(*ast.CallExpr)
			preSrc.WriteTabbed()
			for j := range call.ReturnAmount {
				src.Write(gen.GenerateExpr(assignStmt.Identifiers[index+j]))
				vars = append(vars, GenerateVar(hyVar))
				preSrc.Write(vars[j])
				if j != call.ReturnAmount-1 {
					preSrc.Write(", ")
					src.Write(", ")
				} else {
					preSrc.Write(" = ", gen.callExpr(*call, false), "\n")
					src.Write(" = ", strings.Join(vars, ", "), "\n")
				}
			}
			vars = []string{}
			index += call.ReturnAmount
		} else {
			src.Write(gen.GenerateExpr(assignStmt.Identifiers[index]), " = ", gen.GenerateExpr(assignStmt.Values[i]), "\n")
			index++
		}
		if index >= len(assignStmt.Identifiers) {
			break
		}
	}

	gen.Write(preSrc.String())
	gen.Write(src.String())
}

func (gen *Generator) returnStmt(node ast.ReturnStmt) {
	src := StringBuilder{}

	src.WriteTabbed("return ")
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr)
		src.Write(val)
		if i != len(node.Args)-1 {
			src.Write(", ")
		}
	}

	gen.Write(src.String())
}

func (gen *Generator) yieldStmt(node ast.YieldStmt) {
	src := StringBuilder{}

	ctx := gen.YieldContexts.Top().Item
	lenVars := len(ctx.vars)

	src.WriteTabbed()
	for i, v := range ctx.vars {
		if i == lenVars-1 {
			src.Write(fmt.Sprintf("%s = ", v))
		} else {
			src.Write(fmt.Sprintf("%s, ", v))
		}
	}
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr)
		src.Write(val)
		if i != len(node.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write("\n")
	src.WriteTabbed("goto", ctx.label)

	gen.Write(src.String())
}

func (gen *Generator) breakStmt(_ ast.BreakStmt) {
	gen.WriteTabbed("break")
}

func (gen *Generator) continueStmt(_ ast.ContinueStmt) {
	label := gen.ContinueLabels.Top().Item

	gen.WriteTabbed(fmt.Sprintf("goto %s", label))
}

func (gen *Generator) repeatStmt(node ast.RepeatStmt) {

	end := gen.GenerateExpr(node.Iterator)
	start := gen.GenerateExpr(node.Start)
	skip := gen.GenerateExpr(node.Skip)
	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable)
		gen.WriteTabbed("for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		gen.WriteTabbed("for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("RepeatStmt", gotoLabel)

	gen.GenerateBody(node.Body)

	gen.ContinueLabels.Pop("RepeatStmt")

	gen.WriteTabbed("::" + gotoLabel + "::\n")
	gen.Write("end")
}

func (gen *Generator) whileStmt(node ast.WhileStmt) {
	gen.WriteTabbed("while ", gen.GenerateExpr(node.Condition), " do\n")

	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("WhileStmt", gotoLabel)

	gen.GenerateBody(node.Body)

	gen.ContinueLabels.Pop("WhileStmt")

	gen.Write(gen.String())
	gen.WriteTabbed("::" + gotoLabel + "::\n")
	gen.WriteTabbed("end")
}

func (gen *Generator) forStmt(node ast.ForStmt) {
	gen.WriteTabbed("for ")

	pairs := "pairs"
	if node.OrderedIteration {
		pairs = "ipairs"
	}
	iterator := gen.GenerateExpr(node.Iterator)
	key := gen.GenerateExpr(node.First)
	if node.Second == nil {
		gen.Write(key, ", _ in  ", pairs, " (", iterator, ") do\n")
	} else {
		value := gen.GenerateExpr(node.Second)
		gen.Write(key, ", ", value, " in ", pairs, "(", iterator, ") do\n")
	}
	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("ForStmt", gotoLabel)

	gen.GenerateBody(node.Body)

	gen.ContinueLabels.Pop("ForStmt")

	gen.WriteTabbed("::" + gotoLabel + "::\n")
	gen.WriteTabbed("end")
	gen.Write(gen.String())
}

func (gen *Generator) tickStmt(node ast.TickStmt) {
	tickTabs := getTabs()

	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable)
		gen.Write(tickTabs, "local ", variable, " = 0\n")
		gen.Write(tickTabs, "pewpew.add_update_callback(function()\n")
		gen.WriteTabbed(variable, " = ", variable, " + 1\n")
	} else {
		gen.Write(tickTabs, "pewpew.add_update_callback(function()\n")
	}

	gen.GenerateBody(node.Body)

	gen.Write(tickTabs, "end)")
	gen.Write(gen.String())
}

func (gen *Generator) GenerateParams(params []ast.FunctionParam) {
	var variadicParam string
	for i, param := range params {
		if param.Type.IsVariadic {
			gen.WriteString("...")
			variadicParam = gen.WriteVar(param.Name.Lexeme)
		} else {
			gen.Write(gen.WriteVar(param.Name.Lexeme))
		}
		if i != len(params)-1 {
			gen.Write(", ")
		}
	}
	gen.Write(")\n")

	if variadicParam != "" {
		gen.WriteTabbed("local ", variadicParam, " = {...}\n")
	}
}

func (gen *Generator) matchStmt(node ast.MatchStmt) {
	ifStmt := ast.IfStmt{
		BoolExpr: &ast.BinaryExpr{Left: node.ExprToMatch, Operator: tokens.Token{Type: tokens.EqualEqual, Lexeme: "=="}, Right: node.Cases[0].Expression},
		Body:     node.Cases[0].Body,
	}
	has_default := false
	for i := range node.Cases {
		if node.Cases[i].Expression.GetToken().Lexeme == "_" {
			has_default = true
		}
		if i == 0 || (i == len(node.Cases)-1 && has_default) {
			continue
		}
		elseIfStmt := ast.IfStmt{
			BoolExpr: &ast.BinaryExpr{Left: node.ExprToMatch, Operator: tokens.Token{Type: tokens.EqualEqual, Lexeme: "=="}, Right: node.Cases[i].Expression},
			Body:     node.Cases[i].Body,
		}
		ifStmt.Elseifs = append(ifStmt.Elseifs, &elseIfStmt)
	}

	if has_default {
		ifStmt.Else = &ast.IfStmt{
			Body: node.Cases[len(node.Cases)-1].Body,
		}
	}

	gen.ifStmt(ifStmt)
}

func (gen *Generator) destroyStmt(node ast.DestroyStmt) {
	src := StringBuilder{}

	src.WriteTabbed()
	src.Write(envMap[node.EnvName], hyEntity, node.EntityName, "_Destroy(", gen.GenerateExpr(node.Identifier))
	for _, arg := range node.Args {
		src.Write(", ")
		src.Write(gen.GenerateExpr(arg))
	}
	src.Write(")")
	gen.Write(src.String())
}
