package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/tokens"
)

func (gen *Generator) envStmt(node ast.EnvironmentDecl) {
	for i := range node.Requirements {
		gen.Write("require(\"", node.Requirements[i], "\")\n")
	}
}

func (gen *Generator) ifStmt(node ast.IfStmt) {
	expr := gen.GenerateExpr(node.BoolExpr) // very important that this is called before gen.WriteTabbed (entityExpr might write on gen)
	gen.WriteTabbed("if ", expr, " then\n")

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
	values := []string{}

	src.WriteTabbed()
	for i, v := range assignStmt.Identifiers {
		src.Write(gen.GenerateExpr(v))
		if i != len(assignStmt.Identifiers)-1 {
			src.Write(", ")
		}
	}
	src.Write(" = ")
	if call, ok := assignStmt.Values[0].(ast.CallNode); ok && call.GetReturnAmount() > 1 && assignStmt.AssignOp.Type != tokens.Equal {
		preSrc.WriteTabbed("local ")
		for i := range call.GetReturnAmount() {
			values = append(values, GenerateVar(hyVar))
			preSrc.Write(values[i])
			if i != call.GetReturnAmount()-1 {
				preSrc.Write(", ")
			}
		}
		preSrc.Write(" = ", gen.GenerateExpr(assignStmt.Values[0]), "\n")
	} else {
		for _, v := range assignStmt.Values {
			values = append(values, gen.GenerateExpr(v))
		}
	}
	for i := range values {
		if assignStmt.AssignOp.Type != tokens.Equal {
			op := tokens.TokenType(int(assignStmt.AssignOp.Type) - 1)
			src.Write(fmt.Sprintf("%s %s (%s)", gen.GenerateExpr(assignStmt.Identifiers[i]), op, values[i]))
		} else {
			src.Write(values[i])
		}
		if i != len(values)-1 {
			src.Write(", ")
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
	src.WriteTabbed("goto ", ctx.label)

	gen.Write(src.String())
}

func (gen *Generator) breakStmt(_ ast.BreakStmt) {
	if gen.BreakLabels.Top().Item != "" {
		gen.WriteTabbed("goto ", gen.BreakLabels.Top().Item)
		return
	}
	gen.WriteTabbed("break")
}

func (gen *Generator) continueStmt(_ ast.ContinueStmt) {
	label := gen.ContinueLabels.Top().Item

	gen.WriteTabbed("goto ", label)
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
	gen.BreakLabels.Push("RepeatStmt", "")

	gen.GenerateBody(node.Body)

	gen.BreakLabels.Pop("RepeatStmt")
	gen.ContinueLabels.Pop("RepeatStmt")

	gen.WriteTabbed("::" + gotoLabel + "::\n")
	gen.Write("end")
}

func (gen *Generator) whileStmt(node ast.WhileStmt) {
	gen.WriteTabbed("while ", gen.GenerateExpr(node.Condition), " do\n")

	gotoLabel := GenerateVar(hyGotoLabel)
	gen.ContinueLabels.Push("WhileStmt", gotoLabel)
	gen.BreakLabels.Push("WhileStmt", "")

	gen.GenerateBody(node.Body)

	gen.BreakLabels.Pop("WhileStmt")
	gen.ContinueLabels.Pop("WhileStmt")

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
	gen.BreakLabels.Push("ForStmt", "")

	gen.GenerateBody(node.Body)

	gen.BreakLabels.Pop("ForStmt")
	gen.ContinueLabels.Pop("ForStmt")

	gen.WriteTabbed("::" + gotoLabel + "::\n")
	gen.WriteTabbed("end")
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
	gotoLabel := GenerateVar(hyGotoLabel)

	toMatch := gen.GenerateExpr(node.ExprToMatch)
	for i, matchCase := range node.Cases {
		conditionsSrc := StringBuilder{}
		for j, expr := range matchCase.Expressions {
			conditionsSrc.Write(toMatch, " == ", gen.GenerateExpr(expr))
			if j != len(matchCase.Expressions)-1 {
				conditionsSrc.Write(" or ")
			}
		}
		if i == 0 {
			gen.WriteTabbed("if ", conditionsSrc.String(), " then\n")
		} else if i == len(node.Cases)-1 {
			gen.WriteTabbed("else\n")
		} else {
			gen.WriteTabbed("elseif ", conditionsSrc.String(), " then\n")
		}
		gen.BreakLabels.Push("MatchExpr", gotoLabel)
		gen.GenerateBody(matchCase.Body)
		gen.BreakLabels.Pop("MatchExpr")
	}

	gen.WriteTabbed("end\n")
	gen.WriteTabbed(fmt.Sprintf("::%s::\n", gotoLabel))
}

func (gen *Generator) destroyStmt(node ast.DestroyStmt) {
	src := StringBuilder{}

	src.WriteTabbed()
	src.Write(hyEntity, envMap[node.EnvName], node.EntityName, "_Destroy(", gen.GenerateExpr(node.Identifier))
	for _, arg := range node.Args {
		src.Write(", ")
		src.Write(gen.GenerateExpr(arg))
	}
	src.Write(")")
	gen.Write(src.String())
}
