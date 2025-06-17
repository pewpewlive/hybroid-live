package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/generator/mapping"
	"hybroid/tokens"
)

func (gen *Generator) entityExpr(node ast.EntityEvaluationExpr) string {
	src := core.StringBuilder{}
	var op string
	switch node.Operator.Type {
	case tokens.Is:
		op = "~="
	case tokens.Isnt:
		op = "=="
	}
	if node.OfficialEntityType {
		switch node.Operator.Type {
		case tokens.Is:
			op = "=="
		case tokens.Isnt:
			op = "~="
		}
		src.Write("pewpew.get_entity_type(", gen.GenerateExpr(node.Expr), ") ", op, " ", "pewpew.EntityType.", mapping.PewpewEnums["EntityType"][node.Type.GetToken().Lexeme])
		return src.String()
	}
	expr := gen.GenerateExpr(node.Expr)

	src.Write(hyEntity, envMap[node.EnvName], node.EntityName, "[", expr, "] ", op, " nil")

	if node.ConvertedVarName != nil {
		gen.Twrite(gen.LatestSrc, "local ", gen.WriteVar(node.ConvertedVarName.Lexeme), " = ", expr, "\n")
	}
	return src.String()
}

func (gen *Generator) binaryExpr(node ast.BinaryExpr) string {
	left, right := gen.GenerateExpr(node.Left), gen.GenerateExpr(node.Right)
	var op string
	switch node.Operator.Type {
	case tokens.BangEqual:
		op = "~="
	case tokens.BackSlash:
		op = "//"
	default:
		op = node.Operator.Lexeme
	}
	return fmt.Sprintf("%s %s %s", left, op, right)
}

func (gen *Generator) literalExpr(node ast.LiteralExpr) string {
	switch node.GetToken().Type {
	case tokens.String:
		return fmt.Sprintf("\"%v\"", node.Value)
	case tokens.Fixed, tokens.Radian:
		return fmt.Sprintf("%vfx", fixedToFx(node.Value))
	case tokens.FixedPoint:
		return fmt.Sprintf("%vfx", node.Value)
	case tokens.Degree:
		return fmt.Sprintf("%vfx", degToRad(node.Value))
	default:
		return fmt.Sprintf("%v", node.Value)
	}
}

func (gen *Generator) identifierExpr(node ast.IdentifierExpr) string {
	if gen.env == ast.MeshEnv && node.Name.Lexeme == "meshes" {
		return "meshes"
	}
	if gen.env == ast.SoundEnv && node.Name.Lexeme == "sounds" {
		return "sounds"
	}
	if node.Type == ast.Raw || node.Name.Lexeme == "_" {
		return node.Name.Lexeme
	}
	return gen.WriteVar(node.Name.Lexeme)
}

func (gen *Generator) groupingExpr(node ast.GroupExpr) string {
	return fmt.Sprintf("(%s)", gen.GenerateExpr(node.Expr))
}

func (gen *Generator) listExpr(node ast.ListExpr) string {
	src := core.StringBuilder{}

	src.Write("{")
	for i, expr := range node.List {
		src.Write(gen.GenerateExpr(expr))

		if i != len(node.List)-1 {
			src.Write(", ")
		}
	}
	src.Write("}")

	return src.String()
}

func (gen *Generator) callExpr(node ast.CallExpr, tabbed bool) string {
	src := core.StringBuilder{}
	fn := gen.GenerateExpr(node.Caller)

	if tabbed {
		src.Write(gen.tabString(), fn, "(")
	} else {
		src.Write(fn, "(")
	}
	src.Write(gen.GenerateArgs(node.Args))

	return src.String()
}

func (gen *Generator) mapExpr(node ast.MapExpr) string {
	src := core.StringBuilder{}

	src.Write("{\n")
	gen.tabCount++
	index := 0
	for _, v := range node.KeyValueList {
		val := gen.GenerateExpr(v.Expr)

		token := v.Key.GetToken()
		ident := token.Lexeme

		if index != len(node.KeyValueList)-1 {
			src.Write(gen.tabString(), fmt.Sprintf("[%s] = %v,\n", ident, val))
		} else {
			src.Write(gen.tabString(), fmt.Sprintf("[%s] = %v\n", ident, val))
		}
		index++
	}
	gen.tabCount--
	src.Write(gen.tabString(), "}")

	return src.String()
}

func (gen *Generator) unaryExpr(node ast.UnaryExpr) string {
	var op string
	switch node.Operator.Lexeme {
	case "!":
		op = "not "
	default:
		op = node.Operator.Lexeme
	}
	return fmt.Sprintf("%s%s", op, gen.GenerateExpr(node.Value))
}

func (gen *Generator) accessExpr(node ast.AccessExpr) string {
	str := ""
	if node.Start.GetType() == ast.SelfExpression && node.Start.(*ast.SelfExpr).Type == ast.EntityMethod {
		str = "Self"
	} else {
		str = gen.GenerateExpr(node.Start)
	}

	for i := range node.Accessed {
		accessed := node.Accessed[i]
		switch expr := accessed.(type) {
		case *ast.FieldExpr:
			if expr.Index == 0 {
				str = fmt.Sprintf("%s[\"%s\"]", str, expr.Field.GetToken().Lexeme)
				break
			}
			str = fmt.Sprintf("%s[%v]", str, expr.Index)
		case *ast.MemberExpr:
			str = fmt.Sprintf("%s[%s]", str, gen.GenerateExpr(expr.Member))
		case *ast.EntityAccessExpr:
			tableAccess := hyEntity + envMap[expr.EnvName] + expr.EntityName
			str = fmt.Sprintf("%s[%s%s]", tableAccess, str, gen.GenerateExpr(expr.Expr))
		}
	}

	return str
}

func (gen *Generator) entityAccessExpr(node ast.EntityAccessExpr) string {
	src := core.StringBuilder{}
	src.Write(hyEntity, envMap[node.EnvName], node.EntityName, "[", gen.GenerateExpr(node.Expr), "]")
	return src.String()
}

func (gen *Generator) memberExpr(node ast.MemberExpr) string {
	src := core.StringBuilder{}
	src.Write("[", gen.GenerateExpr(node.Member), "]")
	return src.String()
}

func (gen *Generator) fieldExpr(node ast.FieldExpr) string {
	if node.Index == 0 {
		return fmt.Sprintf("[\"%s\"]", node.Field.GetToken().Lexeme)
	}
	return fmt.Sprintf("[%v]", node.Index)
}

func (gen *Generator) functionExpr(fn ast.FunctionExpr) string {
	src := core.StringBuilder{}
	src.Write("function(")
	gen.GenerateParams(&src, fn.Params)
	gen.GenerateBody(&src, fn.Body)
	src.Write(gen.tabString(), "end")

	return src.String()
}

func (gen *Generator) structExpr(node ast.StructExpr) string {
	src := core.StringBuilder{}

	src.Write("{\n")
	gen.tabCount++
	for i, v := range node.Fields {
		src.Write(gen.tabString(), v.Name.Lexeme, " = ", gen.GenerateExpr(node.Expressions[i]))
		if i != len(node.Fields)-1 {
			src.Write(", ")
		}
		src.Write("\n")
	}
	gen.tabCount--
	src.Write(gen.tabString(), "}")

	return src.String()
}

func (gen *Generator) selfExpr(self ast.SelfExpr) string {
	if self.Type == ast.ClassMethod {
		return "Self"
	}

	return "id"
}

func (gen *Generator) newExpr(new ast.NewExpr, stmt bool) string {
	src := core.StringBuilder{}

	if stmt {
		src.Write(gen.tabString())
	}

	gen.envPrefixName = envMap[new.EnvName]
	name := gen.GenerateExpr(new.Type.Name)
	src.Write(hyClass, name, "_New(")
	for i, arg := range new.Args {
		src.Write(gen.GenerateExpr(arg))
		if i != len(new.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")

	return src.String()
}

func (gen *Generator) spawnExpr(spawn ast.SpawnExpr, stmt bool) string {
	src := core.StringBuilder{}

	if stmt {
		src.Write(gen.tabString())
	}
	gen.envPrefixName = envMap[spawn.EnvName]
	name := gen.GenerateExpr(spawn.Type.Name)
	src.Write(hyEntity, name, "_Spawn(")
	for i, arg := range spawn.Args {
		src.Write(gen.GenerateExpr(arg))
		if i != len(spawn.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")
	return src.String()
}

func (gen *Generator) matchExpr(match ast.MatchExpr) string {
	src := core.StringBuilder{}
	varsSrc := core.StringBuilder{}
	vars := []string{}
	gotoLabel := GenerateVar(hyGotoLabel)

	for i := 0; i < match.ReturnAmount; i++ {
		helperVarName := GenerateVar(hyVar)
		if i == 0 {
			gen.Twrite(&src, "local ", helperVarName)
			varsSrc.Write(helperVarName)
		} else {
			src.Write(", ", helperVarName)
			varsSrc.Write(", ", helperVarName)
		}
		vars = append(vars, helperVarName)
	}
	ctx := NewYieldContext(vars, gotoLabel)

	src.Write("\n")

	node := match.MatchStmt
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
		} else if i == len(node.Cases)-1 {
			gen.Twrite(&src, "else\n")
		} else {
			gen.Twrite(&src, "elseif ", conditionsSrc.String(), " then\n")
		}
		gen.YieldContexts.Push("MatchExpr", ctx)

		gen.GenerateBody(&src, matchCase.Body)

		gen.YieldContexts.Pop("MatchExpr")
	}

	gen.Twrite(&src, "end\n")
	gen.Twrite(&src, "::", gotoLabel, "::\n")

	gen.Twrite(gen.LatestSrc, src.String())
	return varsSrc.String()
}

func (gen *Generator) envAccessExpr(node ast.EnvAccessExpr) string {
	envName := node.PathExpr.Path.Lexeme
	gen.envPrefixName = envMap[envName]
	accessed := gen.GenerateExpr(node.Accessed)

	var prefix string
	switch envName {
	case "Pewpew":
		prefix = "pewpew."
		accessed = mapping.PewpewVariables[accessed]
	case "Fmath":
		prefix = "fmath."
		accessed = mapping.FmathVariables[accessed]
	case "Math":
		prefix = "math."
		accessed = mapping.MathVariables[accessed]
	case "String":
		prefix = "string."
		accessed = mapping.StringVariables[accessed]
	case "Table":
		prefix = "table."
		accessed = mapping.TableVariables[accessed]
	default:
		prefix = ""
	}

	return prefix + accessed
}

func (gen *Generator) methodCallExpr(methodCall ast.MethodCallExpr, stmt bool) string {
	src := core.StringBuilder{}

	if stmt {
		src.Write(gen.tabString())
	}
	var extra string
	if methodCall.MethodType == ast.ClassMethod {
		extra = hyClass
	} else {
		extra = hyEntity
	}
	src.Write(extra, envMap[methodCall.EnvName], methodCall.TypeName, "_", methodCall.MethodName, "(", gen.GenerateExpr(methodCall.Caller))
	for i := range methodCall.Args {
		src.Write(", ", gen.GenerateExpr(methodCall.Args[i]))
	}
	if stmt {
		src.Write(")\n")
	} else {
		src.Write(")")
	}

	return src.String()
}

func (gen *Generator) methodExpr(method ast.MethodExpr) string {
	src := core.StringBuilder{}

	var extra string
	if method.MethodType == ast.ClassMethod {
		extra = hyClass
	} else {
		extra = hyEntity
	}
	src.Write(extra, envMap[method.EnvName], method.TypeName, "_", method.MethodName)

	return src.String()
}
