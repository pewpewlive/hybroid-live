package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
)

func (gen *Generator) envStmt(node ast.EnvironmentDecl) string {
	src := core.StringBuilder{}
	for i := range node.Requirements {
		src.Write("require(\"", node.Requirements[i], "\")")
		if i != len(node.Requirements)-1 {
			src.Write("\n")
		}
	}
	return src.String()
}

func (gen *Generator) ifStmt(node ast.IfStmt) string {
	src := core.StringBuilder{}
	expr := gen.GenerateExpr(node.BoolExpr) // very important that this is called before gen.Twrite (entityExpr might write on gen)
	gen.Twrite(&src, "if ", expr, " then\n")

	gen.GenerateBody(&src, node.Body)
	for _, elseif := range node.Elseifs {
		gen.Twrite(&src, "elseif ", gen.GenerateExpr(elseif.BoolExpr), " then\n")
		gen.GenerateBody(&src, elseif.Body)
	}
	if node.Else != nil {
		gen.Twrite(&src, "else \n")
		gen.GenerateBody(&src, node.Else.Body)
	}

	gen.Twrite(&src, "end")

	return src.String()
}

func (gen *Generator) assignmentStmt(assignStmt ast.AssignmentStmt) string {
	src := core.StringBuilder{}
	preSrc := core.StringBuilder{}
	values := []string{}

	src.Write(gen.tabString())
	for i, v := range assignStmt.Identifiers {
		src.Write(gen.GenerateExpr(v))
		if i != len(assignStmt.Identifiers)-1 {
			src.Write(", ")
		}
	}
	src.Write(" = ")
	if call, ok := assignStmt.Values[0].(ast.CallNode); ok && call.GetReturnAmount() > 1 && assignStmt.AssignOp.Type != tokens.Equal {
		preSrc.Write(gen.tabString(), "local ")
		for i := range call.GetReturnAmount() {
			values = append(values, GenerateVar(hyVar))
			preSrc.Write(values[i])
			if i != call.GetReturnAmount()-1 {
				preSrc.Write(", ")
			}
		}
		preSrc.Write(" = ", gen.GenerateExpr(assignStmt.Values[0]))
	} else {
		for _, v := range assignStmt.Values {
			values = append(values, gen.GenerateExpr(v))
		}
	}
	for i := range values {
		if assignStmt.AssignOp.Type != tokens.Equal {
			op := tokens.TokenType(int(assignStmt.AssignOp.Type) - 1)
			src.Writef("%s %s (%s)", gen.GenerateExpr(assignStmt.Identifiers[i]), op, values[i])
		} else {
			src.Write(values[i])
		}
		if i != len(values)-1 {
			src.Write(", ")
		}
	}

	src.Write(preSrc.String())
	return src.String()
}

func (gen *Generator) returnStmt(node ast.ReturnStmt) string {
	src := core.StringBuilder{}

	src.Write(gen.tabString(), "return ")
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr)
		src.Write(val)
		if i != len(node.Args)-1 {
			src.Write(", ")
		}
	}

	return src.String()
}

func (gen *Generator) yieldStmt(node ast.YieldStmt) string {
	src := core.StringBuilder{}

	ctx := gen.YieldContexts.Top().Item
	lenVars := len(ctx.vars)

	src.Write(gen.tabString())
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
	src.Write(gen.tabString(), "goto ", ctx.label)

	return src.String()
}

func (gen *Generator) breakStmt(_ ast.BreakStmt) string {
	src := core.StringBuilder{}
	if gen.BreakLabels.Top().Item == "" {
		gen.Twrite(&src, "break")
		return src.String()
	}

	gen.Twrite(&src, "goto ", gen.BreakLabels.Top().Item)
	return src.String()
}

func (gen *Generator) continueStmt(_ ast.ContinueStmt) string {
	label := gen.ContinueLabels.Top().Item
	return gen.tabString() + "goto " + label
}

func (gen *Generator) repeatStmt(node ast.RepeatStmt) string {
	src := core.StringBuilder{}
	end := gen.GenerateExpr(node.Iterator)
	start := gen.GenerateExpr(node.Start)
	skip := gen.GenerateExpr(node.Skip)
	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable)
		gen.Twrite(&src, "for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		gen.Twrite(&src, "for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("RepeatStmt", gotoLabel)
	gen.BreakLabels.Push("RepeatStmt", "")

	gen.GenerateBody(&src, node.Body)

	gen.BreakLabels.Pop("RepeatStmt")
	gen.ContinueLabels.Pop("RepeatStmt")

	gen.tabCount++
	gen.Twrite(&src, "::"+gotoLabel+"::\n")
	gen.tabCount--
	gen.Twrite(&src, "end")
	return src.String()
}

func (gen *Generator) whileStmt(node ast.WhileStmt) string {
	src := core.StringBuilder{}
	gen.Twrite(&src, "while ", gen.GenerateExpr(node.Condition), " do\n")

	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("WhileStmt", gotoLabel)
	gen.BreakLabels.Push("WhileStmt", "")

	gen.GenerateBody(&src, node.Body)

	gen.BreakLabels.Pop("WhileStmt")
	gen.ContinueLabels.Pop("WhileStmt")

	gen.tabCount++
	gen.Twrite(&src, "::"+gotoLabel+"::\n")
	gen.tabCount--
	gen.Twrite(&src, "end")
	return src.String()
}

func (gen *Generator) forStmt(node ast.ForStmt) string {
	src := core.StringBuilder{}
	gen.Twrite(&src, "for ")

	pairs := "pairs"
	if node.OrderedIteration {
		pairs = "ipairs"
	}
	iterator := ""
	if node.IsEntity {
		iterator = hyEntity + envMap[node.EnvName] + node.EntityName
	} else {
		iterator = gen.GenerateExpr(node.Iterator)
	}
	key := gen.GenerateExpr(node.First)
	if node.Second == nil {
		src.Write(key, ", _ in  ", pairs, " (", iterator, ") do\n")
	} else {
		value := gen.GenerateExpr(node.Second)
		src.Write(key, ", ", value, " in ", pairs, "(", iterator, ") do\n")
	}
	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("ForStmt", gotoLabel)
	gen.BreakLabels.Push("ForStmt", "")

	gen.GenerateBody(&src, node.Body)

	gen.BreakLabels.Pop("ForStmt")
	gen.ContinueLabels.Pop("ForStmt")

	gen.tabCount++
	gen.Twrite(&src, "::"+gotoLabel+"::\n")
	gen.tabCount--
	gen.Twrite(&src, "end")
	return src.String()
}

func (gen *Generator) tickStmt(node ast.TickStmt) string {
	src := core.StringBuilder{}
	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable)
		gen.Twrite(&src, "local ", variable, " = 0\n")
		gen.Twrite(&src, "pewpew.add_update_callback(function()\n")
		gen.Twrite(&src, variable, " = ", variable, " + 1\n")
	} else {
		gen.Twrite(&src, "pewpew.add_update_callback(function()\n")
	}

	gen.GenerateBody(&src, node.Body)

	gen.Twrite(&src, "end)")
	return src.String()
}

func (gen *Generator) matchStmt(node ast.MatchStmt) string {
	src := core.StringBuilder{}
	label := GenerateVar(hyGotoLabel)
	toMatch := gen.GenerateExpr(node.ExprToMatch)
	hyVar := GenerateVar(hyVar)
	gen.Twrite(&src, "local ", hyVar, " = ", toMatch, "\n")
	for i, matchCase := range node.Cases {
		conditionsSrc := core.StringBuilder{}
		for j, expr := range matchCase.Expressions {
			conditionsSrc.Write(hyVar, " == ", gen.GenerateExpr(expr))
			if j != len(matchCase.Expressions)-1 {
				conditionsSrc.Write(" or ")
			}
		}
		if i == 0 {
			gen.Twrite(&src, "if ", conditionsSrc.String(), " then\n")
		} else if i == len(node.Cases)-1 && node.HasDefault {
			gen.Twrite(&src, "else\n")
		} else {
			gen.Twrite(&src, "elseif ", conditionsSrc.String(), " then\n")
		}
		gen.BreakLabels.Push("MatchStmt", label)
		gen.GenerateBody(&src, matchCase.Body)
		gen.BreakLabels.Pop("MatchStmt")
	}

	gen.Twrite(&src, "end\n")
	gen.Twrite(&src, "::", label, "::")
	return src.String()
}

func (gen *Generator) destroyStmt(node ast.DestroyStmt) string {
	src := core.StringBuilder{}

	src.Write(gen.tabString())
	src.Write(hyEntity, envMap[node.EnvName], node.EntityName, "_Destroy(", gen.GenerateExpr(node.Identifier))
	for _, arg := range node.Args {
		src.Write(", ")
		src.Write(gen.GenerateExpr(arg))
	}
	src.Write(")")
	return src.String()
}
