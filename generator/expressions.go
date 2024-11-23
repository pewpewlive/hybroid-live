package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/tokens"
)

func (gen *Generator) entityExpr(node ast.EntityExpr, scope *GenScope) string {
	src := StringBuilder{}
	var op string
	switch node.Operator.Type {
	case tokens.Is:
		op = "~="
	case tokens.Isnt:
		op = "=="
	default:
		op = node.Operator.Lexeme
	}
	if node.OfficialEntityType {
		src.Append("pewpew.get_entity_type(", gen.GenerateExpr(node.Expr, scope), ") ", op, " ", "pewpew.EntityType.", PewpewEnums["EntityType"][node.Type.GetToken().Lexeme])
		return src.String()
	}
	expr := gen.GenerateExpr(node.Expr, scope)

	src.Append(envMap[node.EnvName], hyEntity, node.EntityName, "[", expr, "] ", op, " nil")

	if node.ConvertedVarName != nil {
		preSrc := StringBuilder{}

		preSrc.Append("local ", gen.WriteVar(node.ConvertedVarName.Lexeme), " = ", expr, "\n")
		gen.Future = preSrc.String()
	}
	return src.String()
}

func (gen *Generator) binaryExpr(node ast.BinaryExpr, scope *GenScope) string {
	left, right := gen.GenerateExpr(node.Left, scope), gen.GenerateExpr(node.Right, scope)
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
	switch node.ValueType {
	case ast.String:
		return fmt.Sprintf("\"%v\"", node.Value)
	case ast.Fixed:
		return fmt.Sprintf("%vfx", fixedToFx(node.Value))
	case ast.FixedPoint, ast.Radian:
		return fmt.Sprintf("%vfx", node.Value)
	case ast.Degree:
		return fmt.Sprintf("%vfx", degToRad(node.Value))
	default:
		return fmt.Sprintf("%v", node.Value)
	}
}

func (gen *Generator) identifierExpr(node ast.IdentifierExpr, _ *GenScope) string {
	if gen.envType == ast.MeshEnv && node.Name.Lexeme == "meshes" {
		return "meshes"
	}
	if gen.envType == ast.SoundEnv && node.Name.Lexeme == "sounds" {
		return "sounds"
	}
	return gen.WriteVar(node.Name.Lexeme)
}

func (gen *Generator) groupingExpr(node ast.GroupExpr, scope *GenScope) string {
	return fmt.Sprintf("(%s)", gen.GenerateExpr(node.Expr, scope))
}

func (gen *Generator) listExpr(node ast.ListExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.WriteString("{")
	for i, expr := range node.List {
		src.WriteString(gen.GenerateExpr(expr, scope))

		if i != len(node.List)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString("}")

	return src.String()
}

func (gen *Generator) callExpr(node ast.CallExpr, tabbed bool, scope *GenScope) string {
	src := StringBuilder{}
	fn := gen.GenerateExpr(node.Caller, scope)

	if tabbed {
		src.AppendTabbed(fn, "(")
	} else {
		src.Append(fn, "(")
	}
	src.WriteString(gen.GenerateArgs(node.Args, scope))

	return src.String()
}

func (gen *Generator) GenerateArgs(args []ast.Node, scope *GenScope) string {
	src := StringBuilder{}

	for i, arg := range args {
		src.WriteString(gen.GenerateExpr(arg, scope))
		if i != len(args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")

	return src.String()
}

func (gen *Generator) mapExpr(node ast.MapExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.WriteString("{\n")
	TabsCount += 1
	index := 0
	for k, v := range node.Map {
		val := gen.GenerateExpr(v.Expr, scope)

		ident := k.Lexeme

		if k.Type == tokens.String {
			ident = k.Literal
		}

		if index != len(node.Map)-1 {
			src.AppendTabbed(fmt.Sprintf("%s = %v,\n", ident, val))
		} else {
			src.AppendTabbed(fmt.Sprintf("%s = %v\n", ident, val))
		}
		index++
	}
	TabsCount -= 1
	src.AppendTabbed("}")

	return src.String()
}

func (gen *Generator) unaryExpr(node ast.UnaryExpr, scope *GenScope) string {
	var op string
	switch node.Operator.Lexeme {
	case "!":
		op = "not "
	default:
		op = node.Operator.Lexeme
	}
	return fmt.Sprintf("%s%s", op, gen.GenerateExpr(node.Value, scope))
}

func (gen *Generator) fieldExpr(node ast.FieldExpr, scope *GenScope) string {
	src := StringBuilder{}

	if node.ExprType == ast.SelfEntity {
		src.Append(envMap[node.EnvName], hyEntity, node.EntityName, "[", gen.GenerateExpr(node.Identifier, scope), "]")
	} else {
		src.WriteString(gen.GenerateExpr(node.Identifier, scope))
	}

	val := gen.GenerateExpr(node.Property, scope)
	cut := ""
	for i := range val {
		if val[i] == '[' {
			cut = val[i:]
			break
		}
	}
	if node.Index >= 0 {
		val = fmt.Sprintf("[%v]%s", node.Index, cut)
	} else {
		val = "." + val[len(gen.envName):]
	}
	src.WriteString(val)
	return src.String()
}

func (gen *Generator) memberExpr(node ast.MemberExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.WriteString(gen.GenerateExpr(node.Identifier, scope))
	val := gen.GenerateExpr(node.Property, scope)
	name := gen.GenerateExpr(node.PropertyIdentifier, scope)
	val = fmt.Sprintf("[%s]%s", name, val[len(name):])
	src.WriteString(val)

	return src.String()
}

func (gen *Generator) functionExpr(fn ast.FunctionExpr, scope *GenScope) string {
	fnScope := NewGenScope(scope)

	fnScope.WriteString("function (")
	for i, param := range fn.Params {
		fnScope.Append(param.Name.Lexeme)
		if i != len(fn.Params)-1 {
			fnScope.Append(", ")
		}
	}
	fnScope.Append(")\n")

	gen.GenerateBody(fn.Body, &fnScope)

	fnScope.AppendTabbed("end")

	return fnScope.Src.String()
}

func (gen *Generator) structExpr(node ast.StructExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.WriteString("{\n")
	TabsCount += 1
	for i, v := range node.Fields {
		src.AppendTabbed(gen.fieldDeclarationStmt(*v, scope))
		if i != len(node.Fields)-1 {
			src.WriteString(", ")
		}
		src.WriteString("\n")
	}
	TabsCount -= 1
	src.AppendTabbed("}")

	return src.String()
}

func (gen *Generator) selfExpr(self ast.SelfExpr, _ *GenScope) string {
	if self.Type == ast.SelfStruct {
		return "Self"
	} else if self.Type == ast.SelfEntity {
		return "id"
	}
	return ""
}

func (gen *Generator) newExpr(new ast.NewExpr, stmt bool, scope *GenScope) string {
	src := StringBuilder{}

	if stmt {
		src.AppendTabbed()
	}

	length := len(new.Type.Name.GetToken().Lexeme)
	name := gen.GenerateExpr(new.Type.Name, scope)
	fullLength := len(name)
	cut := name[fullLength-length:]
	src.Append(name[:fullLength-length], hyClass, cut, "_New(")
	for i, arg := range new.Args {
		src.WriteString(gen.GenerateExpr(arg, scope))
		if i != len(new.Args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")

	return src.String()
}

func (gen *Generator) matchExpr(match ast.MatchExpr, scope *GenScope) string {
	vars := StringBuilder{}

	gotoLabel := GenerateVar(hyGTL)

	for i := 0; i < match.ReturnAmount; i++ {
		helperVarName := GenerateVar(hyVar)
		if i == 0 {
			scope.Src.AppendTabbed("local ", helperVarName)
			vars.WriteString(helperVarName)
		} else {
			scope.Src.Append(", ", helperVarName)
			vars.Append(", ", helperVarName)
		}
	}

	scope.Src.WriteString("\n")

	node := match.MatchStmt

	toMatch := gen.GenerateExpr(node.ExprToMatch, scope)

	for i, matchCase := range node.Cases {
		if i == 0 {
			scope.AppendTabbed("if ", toMatch, " == ", gen.GenerateExpr(matchCase.Expression, scope), " then\n")
		} else if i == len(node.Cases)-1 {
			scope.AppendTabbed("else\n")
		} else {
			scope.AppendTabbed("elseif ", toMatch, " == ", gen.GenerateExpr(matchCase.Expression, scope), " then\n")
		}

		caseScope := NewGenScope(scope)

		caseScope.ReplaceSettings = map[ReplaceType]string{
			YieldReplacement: vars.String() + " = ",
			GotoReplacement:  "goto " + gotoLabel,
		}

		gen.GenerateBody(matchCase.Body, &caseScope)

		caseScope.ReplaceAll()

		scope.Write(caseScope.Src)
	}

	scope.AppendTabbed("end\n")

	scope.AppendTabbed(fmt.Sprintf("::%s::\n", gotoLabel))

	return vars.String()
}

func (gen *Generator) envAccessExpr(node ast.EnvAccessExpr, scope *GenScope) string {
	accessed := gen.GenerateExpr(node.Accessed, scope)
	accessed = accessed[len(gen.envName):]

	envName := node.PathExpr.Path.Lexeme
	var prefix string
	switch envName {
	case "Pewpew":
		prefix = "pewpew."
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

	gen.libraryVars = nil
	return prefix + accessed
}

func (gen *Generator) spawnExpr(spawn ast.SpawnExpr, stmt bool, scope *GenScope) string {
	src := StringBuilder{}

	if stmt {
		src.AppendTabbed()
	}
	length := len(spawn.Type.Name.GetToken().Lexeme)
	name := gen.GenerateExpr(spawn.Type.Name, scope)
	fullLength := len(name)
	cut := name[fullLength-length:]
	src.Append(name[:fullLength-length], hyEntity, cut, "_Spawn(")
	for i, arg := range spawn.Args {
		src.WriteString(gen.GenerateExpr(arg, scope))
		if i != len(spawn.Args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")
	return src.String()
}

func (gen *Generator) methodCallExpr(methodCall ast.MethodCallExpr, stmt bool, scope *GenScope) string {
	src := StringBuilder{}

	if stmt {
		src.AppendTabbed()
	}
	var extra string
	if methodCall.ExprType == ast.SelfStruct {
		extra = hyClass
	} else {
		extra = hyEntity
	}
	src.Append(envMap[methodCall.EnvName], extra, methodCall.TypeName, "_", methodCall.MethodName, "(", gen.GenerateExpr(methodCall.Identifier, scope))
	for i := range methodCall.Call.Args {
		src.Append(", ", gen.GenerateExpr(methodCall.Call.Args[i], scope))
	}
	if stmt {
		src.WriteString(")\n")
	} else {
		src.WriteString(")")
	}

	return src.String()
}

// func (gen *Generator) builtinExpr(builtin ast.BuiltinExpr) string {
// 	return builtin.Name.Lexeme
// }

// func (gen *Generator) castExpr(cast ast.CastExpr, scope *GenScope) string {
// 	return gen.GenerateExpr(cast.Value, scope)
// }
