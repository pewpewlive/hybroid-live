package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
)

// let a = {a:{b:0},b:[2,3,4]}

// Member Node
// a.a.b

func (gen *Generator) binaryExpr(node ast.BinaryExpr, scope *GenScope) string {
	left, right := gen.GenerateExpr(node.Left, scope), gen.GenerateExpr(node.Right, scope)
	return fmt.Sprintf("%s %s %s", left, node.Operator.Lexeme, right)
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
	return node.Name.Lexeme
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

func (gen *Generator) callExpr(node ast.CallExpr, scope *GenScope) string {
	src := StringBuilder{}
	fn := gen.GenerateExpr(node.Caller, scope)

	src.AppendTabbed(fn, "(")
	for i, arg := range node.Args {
		src.WriteString(gen.GenerateExpr(arg, scope))
		if i != len(node.Args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")

	return src.String()
}

func (gen *Generator) methodCallExpr(node ast.MethodCallExpr, scope *GenScope) string {
	src := StringBuilder{}
	src.AppendTabbed("Hybroid_", node.TypeName, "_", node.MethodName)

	src.Append("(", gen.GenerateExpr(node.Owner, scope))
	if len(node.Args) != 0 {
		src.WriteString(", ")
	}
	for i, arg := range node.Args {
		src.WriteString(gen.GenerateExpr(arg, scope))
		if i != len(node.Args)-1 {
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

		if k.Type == lexer.String {
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
	return fmt.Sprintf("%s%s", node.Operator.Lexeme, gen.GenerateExpr(node.Value, scope))
}

func (gen *Generator) fieldExpr(node ast.FieldExpr, scope *GenScope) string {
	src := StringBuilder{}

	var prop string
	if node.Property != nil {
		prop = gen.GenerateExpr(node.Property, scope)
	}

	var expr string
	if node.Owner == nil {
		expr = gen.GenerateExpr(node.Identifier, scope)
	} else {
		src.Append("[", fmt.Sprintf("%v", node.Index), "]", prop)
		return src.String()
	} // Self.rect
	src.Append(expr, prop)

	return src.String()
}

func (gen *Generator) memberExpr(node ast.MemberExpr, scope *GenScope) string {
	src := StringBuilder{}

	var prop string
	if node.Property != nil {
		prop = gen.GenerateExpr(node.Property, scope)
	}

	if node.Owner == nil {
		src.Append(gen.GenerateExpr(node.Identifier, scope), prop)
		return src.String()
	}

	expr := gen.GenerateExpr(node.Identifier, scope)

	src.Append("[", expr, "]", prop)

	return src.String()
}

func (gen *Generator) directiveExpr(node ast.DirectiveExpr, scope *GenScope) string {
	src := StringBuilder{}

	if node.Identifier.Lexeme != "Environment" {
		src.Append(node.Identifier.Lexeme, "(", gen.GenerateExpr(node.Expr, scope), ")")
	}

	return src.String()
}

func (gen *Generator) anonFnExpr(fn ast.AnonFnExpr, scope *GenScope) string {
	fnScope := NewGenScope(scope)

	TabsCount += 1

	fnScope.WriteString("function (")
	for i, param := range fn.Params {
		fnScope.Append(param.Name.Lexeme)
		if i != len(fn.Params)-1 {
			fnScope.Append(", ")
		}
	}
	fnScope.Append(")\n")

	gen.GenerateString(fn.Body, &fnScope)

	TabsCount -= 1

	fnScope.AppendTabbed("end")

	return fnScope.Src.String()
}

func (gen *Generator) selfExpr(self ast.SelfExpr, _ *GenScope) string {
	if self.Type == ast.SelfStruct {
		return "Self"
	}
	return ""
}

func (gen *Generator) newExpr(new ast.NewExpr, scope *GenScope) string {
	src := StringBuilder{}

	src.Append("Hybroid_", new.Type.Lexeme, "_New(")
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

	gotoLabel := "glab" + RandStr(5)

	for i := 0; i < match.ReturnAmount; i++ {
		helperVarName := "hv" + RandStr(5) // "hv" stands for hybroid variable
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

		TabsCount += 1

		caseScope := NewGenScope(scope)

		gen.GenerateString(matchCase.Body, &caseScope)

		caseScope.DoTheDos(map[DoType]string{
			YieldReplacement: vars.String() + " = ",
			GotoReplacement:  "goto " + gotoLabel,
		})

		scope.Write(caseScope.Src)
		//scope.TransferDos(&caseScope)

		TabsCount -= 1
	}

	scope.AppendTabbed("end\n")

	scope.AppendTabbed(fmt.Sprintf("::%s::\n", gotoLabel))

	return vars.String()
}
