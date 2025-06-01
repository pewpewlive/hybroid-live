package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/tokens"
)

func (gen *Generator) entityExpr(node ast.EntityEvaluationExpr) string {
	src := StringBuilder{}
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
		src.Write("pewpew.get_entity_type(", gen.GenerateExpr(node.Expr), ") ", op, " ", "pewpew.EntityType.", PewpewEnums["EntityType"][node.Type.GetToken().Lexeme])
		return src.String()
	}
	expr := gen.GenerateExpr(node.Expr)

	src.Write(envMap[node.EnvName], hyEntity, node.EntityName, "[", expr, "] ", op, " nil")

	if node.ConvertedVarName != nil {
		preSrc := StringBuilder{}

		preSrc.Write("local ", gen.WriteVar(node.ConvertedVarName.Lexeme), " = ", expr, "\n")
		gen.WriteTabbed(preSrc.String())
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
	return gen.WriteVar(node.Name.Lexeme)
}

func (gen *Generator) groupingExpr(node ast.GroupExpr) string {
	return fmt.Sprintf("(%s)", gen.GenerateExpr(node.Expr))
}

func (gen *Generator) listExpr(node ast.ListExpr) string {
	src := StringBuilder{}

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
	src := StringBuilder{}
	fn := gen.GenerateExpr(node.Caller)

	if tabbed {
		src.WriteTabbed(fn, "(")
	} else {
		src.Write(fn, "(")
	}
	src.Write(gen.GenerateArgs(node.Args))

	return src.String()
}

func (gen *Generator) GenerateArgs(args []ast.Node) string {
	src := StringBuilder{}

	for i, arg := range args {
		src.Write(gen.GenerateExpr(arg))
		if i != len(args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")

	return src.String()
}

func (gen *Generator) mapExpr(node ast.MapExpr) string {
	src := StringBuilder{}

	src.Write("{\n")
	TabsCount += 1
	index := 0
	for _, v := range node.KeyValueList {
		val := gen.GenerateExpr(v.Expr)

		token := v.Key.GetToken()
		ident := token.Lexeme

		if index != len(node.KeyValueList)-1 {
			src.WriteTabbed(fmt.Sprintf("[%s] = %v,\n", ident, val))
		} else {
			src.WriteTabbed(fmt.Sprintf("[%s] = %v\n", ident, val))
		}
		index++
	}
	TabsCount -= 1
	src.WriteTabbed("}")

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
	src := StringBuilder{}
	src.Write(gen.GenerateExpr(node.Start))
	for _, accessed := range node.Accessed {
		switch expr := accessed.(type) {
		case *ast.FieldExpr:
			if gen.pewpewEnum != nil {
				src.Write(".", (*gen.pewpewEnum)[expr.Field.GetToken().Lexeme])
				gen.pewpewEnum = nil
				continue
			}
			src.Write(fmt.Sprintf("[%v]", expr.Index))
		case *ast.MemberExpr:
			src.Write(fmt.Sprintf("[%s]", expr.GetToken().Lexeme))
		}
	}

	return src.String()
}

func (gen *Generator) functionExpr(fn ast.FunctionExpr) string {
	src := StringBuilder{}
	src.WriteString("function (")
	for i, param := range fn.Params {
		gen.Write(param.Name.Lexeme)
		if i != len(fn.Params)-1 {
			gen.Write(", ")
		}
	}
	src.Write(")")
	if len(fn.Body) == 0 {
		src.Write(" end")
		return src.String()
	} else {
		src.Write("\n")
	}
	src.Write(gen.GenerateBodyValue(fn.Body))
	src.WriteTabbed("end")

	return src.String()
}

func (gen *Generator) structExpr(node ast.StructExpr) string {
	src := StringBuilder{}

	src.Write("{\n")
	TabsCount += 1
	for i, v := range node.Fields {
		src.WriteTabbed(gen.GenerateExpr(v), " = ", gen.GenerateExpr(node.Expressions[i]))
		if i != len(node.Fields)-1 {
			src.Write(", ")
		}
		src.Write("\n")
	}
	TabsCount -= 1
	src.WriteTabbed("}")

	return src.String()
}

func (gen *Generator) selfExpr(self ast.SelfExpr) string {
	if self.Type == ast.ClassMethod {
		return "Self"
	} else if self.Type == ast.EntityMethod {
		return "id"
	}
	return ""
}

func (gen *Generator) newExpr(new ast.NewExpr, stmt bool) string {
	src := StringBuilder{}

	if stmt {
		src.WriteTabbed()
	}

	length := len(new.Type.Name.GetToken().Lexeme)
	name := gen.GenerateExpr(new.Type.Name)
	fullLength := len(name)
	cut := name[fullLength-length:]
	src.Write(name[:fullLength-length], hyClass, cut, "_New(")
	for i, arg := range new.Args {
		src.Write(gen.GenerateExpr(arg))
		if i != len(new.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")

	return src.String()
}

func (gen *Generator) matchExpr(match ast.MatchExpr) string {
	varsSrc := StringBuilder{}
	vars := []string{}
	gotoLabel := GenerateVar(hyGotoLabel)

	for i := 0; i < match.ReturnAmount; i++ {
		helperVarName := GenerateVar(hyVar)
		if i == 0 {
			gen.WriteTabbed("local ", helperVarName)
			varsSrc.Write(helperVarName)
		} else {
			gen.Write(", ", helperVarName)
			varsSrc.Write(", ", helperVarName)
		}
		vars = append(vars, helperVarName)
	}
	ctx := NewYieldContext(vars, gotoLabel)

	gen.Write("\n")

	node := match.MatchStmt

	toMatch := gen.GenerateExpr(node.ExprToMatch)

	for i, matchCase := range node.Cases {
		if i == 0 {
			gen.WriteTabbed("if ", toMatch, " == ", gen.GenerateExpr(matchCase.Expression), " then\n")
		} else if i == len(node.Cases)-1 {
			gen.WriteTabbed("else\n")
		} else {
			gen.WriteTabbed("elseif ", toMatch, " == ", gen.GenerateExpr(matchCase.Expression), " then\n")
		}
		gen.YieldContexts.Push("MatchExpr", ctx)

		gen.GenerateBody(matchCase.Body)

		gen.YieldContexts.Pop("MatchExpr")
	}

	gen.WriteTabbed("end\n")

	gen.WriteTabbed(fmt.Sprintf("::%s::\n", gotoLabel))

	return varsSrc.String()
}

func (gen *Generator) envAccessExpr(node ast.EnvAccessExpr) string {
	accessed := gen.GenerateExpr(node.Accessed)
	accessed = accessed[len(gen.envName):]

	envName := node.PathExpr.Path.Lexeme
	var prefix string
	switch envName {
	case "Pewpew":
		prefix = "pewpew."
		accessed = PewpewVariables[accessed]
		temp := PewpewEnums[accessed]
		gen.pewpewEnum = &temp
	case "Fmath":
		prefix = "fmath."
	case "Math":
		prefix = "math."
	case "String":
		prefix = "string."
	case "Table":
		prefix = "table."
	default:
		prefix = envMap[envName]
	}

	return prefix + accessed
}

func (gen *Generator) spawnExpr(spawn ast.SpawnExpr, stmt bool) string {
	src := StringBuilder{}

	if stmt {
		src.WriteTabbed()
	}
	length := len(spawn.Type.Name.GetToken().Lexeme)
	name := gen.GenerateExpr(spawn.Type.Name)
	fullLength := len(name)
	cut := name[fullLength-length:]
	src.Write(name[:fullLength-length], hyEntity, cut, "_Spawn(")
	for i, arg := range spawn.Args {
		src.Write(gen.GenerateExpr(arg))
		if i != len(spawn.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")
	return src.String()
}

func (gen *Generator) methodCallExpr(methodCall ast.MethodCallExpr, stmt bool) string {
	src := StringBuilder{}

	if stmt {
		src.WriteTabbed()
	}
	var extra string
	if methodCall.MethodType == ast.ClassMethod {
		extra = hyClass
	} else {
		extra = hyEntity
	}
	src.Write(envMap[methodCall.EnvName], extra, methodCall.TypeName, "_", methodCall.MethodName, "(", gen.GenerateExpr(methodCall.Caller))
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
