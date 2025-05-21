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
	}
	if node.OfficialEntityType {
		switch node.Operator.Type {
		case tokens.Is:
			op = "=="
		case tokens.Isnt:
			op = "~="
		}
		src.Write("pewpew.get_entity_type(", gen.GenerateExpr(node.Expr, scope), ") ", op, " ", "pewpew.EntityType.", PewpewEnums["EntityType"][node.Type.GetToken().Lexeme])
		return src.String()
	}
	expr := gen.GenerateExpr(node.Expr, scope)

	src.Write(envMap[node.EnvName], hyEntity, node.EntityName, "[", expr, "] ", op, " nil")

	if node.ConvertedVarName != nil {
		preSrc := StringBuilder{}

		preSrc.Write("local ", gen.WriteVar(node.ConvertedVarName.Lexeme), " = ", expr, "\n")
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
	case ast.Fixed, ast.Radian:
		return fmt.Sprintf("%vfx", fixedToFx(node.Value))
	case ast.FixedPoint:
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

	src.Write("{")
	for i, expr := range node.List {
		src.Write(gen.GenerateExpr(expr, scope))

		if i != len(node.List)-1 {
			src.Write(", ")
		}
	}
	src.Write("}")

	return src.String()
}

func (gen *Generator) callExpr(node ast.CallExpr, tabbed bool, scope *GenScope) string {
	src := StringBuilder{}
	fn := gen.GenerateExpr(node.Caller, scope)

	if tabbed {
		src.WriteTabbed(fn, "(")
	} else {
		src.Write(fn, "(")
	}
	src.Write(gen.GenerateArgs(node.Args, scope))

	return src.String()
}

func (gen *Generator) GenerateArgs(args []ast.Node, scope *GenScope) string {
	src := StringBuilder{}

	for i, arg := range args {
		src.Write(gen.GenerateExpr(arg, scope))
		if i != len(args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")

	return src.String()
}

func (gen *Generator) mapExpr(node ast.MapExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.Write("{\n")
	TabsCount += 1
	//index := 0
	// for k, v := range node.Map {
	// 	val := gen.GenerateExpr(v.Expr, scope)

	// 	ident := k.Lexeme

	// 	if k.Type == tokens.String {
	// 		ident = k.Literal
	// 	}

	// 	if index != len(node.Map)-1 {
	// 		src.WriteTabbed(fmt.Sprintf("%s = %v,\n", ident, val))
	// 	} else {
	// 		src.WriteTabbed(fmt.Sprintf("%s = %v\n", ident, val))
	// 	}
	// 	index++
	// }
	TabsCount -= 1
	src.WriteTabbed("}")

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

	// if node.ExprType == ast.SelfEntity {
	// 	src.Write(envMap[node.EnvName], hyEntity, node.EntityName, "[", gen.GenerateExpr(node.Identifier, scope), "]")
	// } else {
	// 	src.Write(gen.GenerateExpr(node.Identifier, scope))
	// }

	// val := gen.GenerateExpr(node.Property, scope)
	// cut := ""
	// for i := range val {
	// 	if val[i] == '[' {
	// 		cut = val[i:]
	// 		break
	// 	}
	// }
	// if node.Index >= 0 {
	// 	val = fmt.Sprintf("[%v]%s", node.Index, cut)
	// } else {
	// 	val = "." + val[len(gen.envName):]
	// }
	// src.Write(val)
	return src.String()
}

func (gen *Generator) memberExpr(node ast.MemberExpr, scope *GenScope) string {
	src := StringBuilder{}

	// src.Write(gen.GenerateExpr(node.Identifier, scope))
	// val := gen.GenerateExpr(node.Property, scope)
	// name := gen.GenerateExpr(node.PropertyIdentifier, scope)
	// val = fmt.Sprintf("[%s]%s", name, val[len(name):])
	// src.Write(val)

	return src.String()
}

func (gen *Generator) functionExpr(fn ast.FunctionExpr, scope *GenScope) string {
	fnScope := NewGenScope(scope)

	fnScope.WriteString("function (")
	for i, param := range fn.Params {
		fnScope.Write(param.Name.Lexeme)
		if i != len(fn.Params)-1 {
			fnScope.Write(", ")
		}
	}
	fnScope.Write(")\n")

	gen.GenerateBody(fn.Body, &fnScope)

	fnScope.WriteTabbed("end")

	return fnScope.String()
}

func (gen *Generator) structExpr(node ast.StructExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.Write("{\n")
	TabsCount += 1
	for i, v := range node.Fields {
		src.WriteTabbed(gen.fieldDeclarationStmt(*v, scope))
		if i != len(node.Fields)-1 {
			src.Write(", ")
		}
		src.Write("\n")
	}
	TabsCount -= 1
	src.WriteTabbed("}")

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
		src.WriteTabbed()
	}

	length := len(new.Type.Name.GetToken().Lexeme)
	name := gen.GenerateExpr(new.Type.Name, scope)
	fullLength := len(name)
	cut := name[fullLength-length:]
	src.Write(name[:fullLength-length], hyClass, cut, "_New(")
	for i, arg := range new.Args {
		src.Write(gen.GenerateExpr(arg, scope))
		if i != len(new.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")

	return src.String()
}

func (gen *Generator) matchExpr(match ast.MatchExpr, scope *GenScope) string {
	vars := StringBuilder{}

	gotoLabel := GenerateVar(hyGTL)

	for i := 0; i < match.ReturnAmount; i++ {
		helperVarName := GenerateVar(hyVar)
		if i == 0 {
			scope.WriteTabbed("local ", helperVarName)
			vars.Write(helperVarName)
		} else {
			scope.Write(", ", helperVarName)
			vars.Write(", ", helperVarName)
		}
	}

	scope.Write("\n")

	node := match.MatchStmt

	toMatch := gen.GenerateExpr(node.ExprToMatch, scope)

	for i, matchCase := range node.Cases {
		if i == 0 {
			scope.WriteTabbed("if ", toMatch, " == ", gen.GenerateExpr(matchCase.Expression, scope), " then\n")
		} else if i == len(node.Cases)-1 {
			scope.WriteTabbed("else\n")
		} else {
			scope.WriteTabbed("elseif ", toMatch, " == ", gen.GenerateExpr(matchCase.Expression, scope), " then\n")
		}

		caseScope := NewGenScope(scope)

		caseScope.ReplaceSettings = map[ReplaceType]string{
			YieldReplacement: vars.String() + " = ",
			GotoReplacement:  "goto " + gotoLabel,
		}

		gen.GenerateBody(matchCase.Body, &caseScope)

		caseScope.ReplaceAll()

		scope.Write(caseScope.String())
	}

	scope.WriteTabbed("end\n")

	scope.WriteTabbed(fmt.Sprintf("::%s::\n", gotoLabel))

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

	return prefix + accessed
}

func (gen *Generator) spawnExpr(spawn ast.SpawnExpr, stmt bool, scope *GenScope) string {
	src := StringBuilder{}

	if stmt {
		src.WriteTabbed()
	}
	length := len(spawn.Type.Name.GetToken().Lexeme)
	name := gen.GenerateExpr(spawn.Type.Name, scope)
	fullLength := len(name)
	cut := name[fullLength-length:]
	src.Write(name[:fullLength-length], hyEntity, cut, "_Spawn(")
	for i, arg := range spawn.Args {
		src.Write(gen.GenerateExpr(arg, scope))
		if i != len(spawn.Args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")
	return src.String()
}

func (gen *Generator) methodCallExpr(methodCall ast.MethodCallExpr, stmt bool, scope *GenScope) string {
	src := StringBuilder{}

	// if stmt {
	// 	src.WriteTabbed()
	// }
	// var extra string
	// if methodCall.ExprType == ast.SelfStruct {
	// 	extra = hyClass
	// } else {
	// 	extra = hyEntity
	// }
	// src.Write(envMap[methodCall.EnvName], extra, methodCall.TypeName, "_", methodCall.MethodName, "(", gen.GenerateExpr(methodCall.Identifier, scope))
	// for i := range methodCall.Call.Args {
	// 	src.Write(", ", gen.GenerateExpr(methodCall.Call.Args[i], scope))
	// }
	// if stmt {
	// 	src.Write(")\n")
	// } else {
	// 	src.Write(")")
	// }

	return src.String()
}

// func (gen *Generator) builtinExpr(builtin ast.BuiltinExpr) string {
// 	return builtin.Name.Lexeme
// }

// func (gen *Generator) castExpr(cast ast.CastExpr, scope *GenScope) string {
// 	return gen.GenerateExpr(cast.Value, scope)
// }
