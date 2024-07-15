package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strconv"
)

func (gen *Generator) envStmt(node ast.EnvironmentStmt, scope *GenScope) {
	for i := range node.Requirements {
		scope.Append("require(\"", node.Requirements[i], "\")\n")
	}
}

func (gen *Generator) ifStmt(node ast.IfStmt, scope *GenScope) {
	ifScope := NewGenScope(scope)

	ifScope.AppendTabbed("if ", gen.GenerateExpr(node.BoolExpr, scope), " then\n")

	gen.GenerateBody(node.Body, &ifScope)
	for _, elseif := range node.Elseifs {
		ifScope.AppendTabbed("elseif ", gen.GenerateExpr(elseif.BoolExpr, scope), " then\n")
		gen.GenerateBody(elseif.Body, &ifScope)
	}
	if node.Else != nil {
		ifScope.AppendTabbed("else \n")
		gen.GenerateBody(node.Else.Body, &ifScope)
	}

	ifScope.AppendTabbed("end\n")

	ifScope.ReplaceAll()

	scope.Write(ifScope.Src)
}

func (gen *Generator) assignmentStmt(assginStmt ast.AssignmentStmt, scope *GenScope) {
	src := StringBuilder{}

	src.AppendTabbed()

	for i, ident := range assginStmt.Identifiers {
		ident := gen.GenerateExpr(ident, scope)
		if i == len(assginStmt.Identifiers)-1 {
			src.Append(ident)
		} else {
			src.Append(ident, ", ")
		}
	}
	src.Append(" = ")

	for i, rightValue := range assginStmt.Values {
		value := gen.GenerateExpr(rightValue, scope)
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

	scope.Write(src)
}

func (gen *Generator) functionDeclarationStmt(node ast.FunctionDeclarationStmt, scope *GenScope) {
	fnScope := NewGenScope(scope)

	if node.IsLocal {
		fnScope.AppendTabbed("local ")
	} else {
		fnScope.AppendTabbed()
	}

	fnScope.Append("function ", gen.WriteVar(node.Name.Lexeme), "(")
	gen.GenerateParams(node.Params, &fnScope)
	fnScope.Append(")\n")

	/* scope src, node src

	function name(a)

		call(1)

		local hy12455
		if a == "a" {
			hy12455 = 1
			goto Leave
		}else {
			hy12455 = 2
			goto Leave
		}
		::Leave::
		call(hy12455)


	*/

	gen.GenerateBody(node.Body, &fnScope)

	fnScope.AppendTabbed("end")

	scope.Write(fnScope.Src)
}

func (gen *Generator) returnStmt(node ast.ReturnStmt, scope *GenScope) {
	src := StringBuilder{}

	src.AppendTabbed("return ")
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr, scope)
		src.Append(val)
		if i != len(node.Args)-1 {
			src.Append(", ")
		}
	}

	scope.Write(src)
}

func (gen *Generator) yieldStmt(node ast.YieldStmt, scope *GenScope) {
	src := StringBuilder{}

	src.AppendTabbed()
	startIndex := src.Len()
	src.Append("yield ")
	endIndex := src.Len()
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr, scope)
		src.Append(val)
		if i != len(node.Args)-1 {
			src.Append(", ")
		}
	}

	src.WriteString("\n")
	src.AppendTabbed()
	startIndex2 := src.Len()
	src.Append("goto hyl")
	endIndex2 := src.Len()
	src.WriteString("\n")

	scopeLength := scope.Src.Len()

	scope.Write(src)

	_range := NewRange(startIndex+scopeLength, endIndex+scopeLength)
	scope.AddReplacement(YieldReplacement, _range)
	_range2 := NewRange(startIndex2+scopeLength, endIndex2+scopeLength)
	scope.AddReplacement(GotoReplacement, _range2)
}

func (gen *Generator) breakStmt(_ ast.BreakStmt, scope *GenScope) {
	scope.AppendTabbed("break")
}

func (gen *Generator) continueStmt(_ ast.ContinueStmt, scope *GenScope) {
	src := StringBuilder{}

	src.AppendTabbed()
	startIndex := src.Len()
	src.Append("continue")
	endIndex := src.Len()

	scopeLength := scope.Src.Len()

	scope.Write(src)

	_range := NewRange(startIndex+scopeLength, endIndex+scopeLength)
	scope.AddReplacement(ContinueReplacement, _range)
}

func (gen *Generator) repeatStmt(node ast.RepeatStmt, scope *GenScope) {
	repeatScope := NewGenScope(scope)

	end := gen.GenerateExpr(node.Iterator, scope)
	start := gen.GenerateExpr(node.Start, scope)
	skip := gen.GenerateExpr(node.Skip, scope)
	if node.Variable.GetValueType() != ast.Unknown {
		variable := gen.GenerateExpr(node.Variable, &repeatScope)
		repeatScope.AppendTabbed("for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		repeatScope.AppendTabbed("for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	gotoLabel := GenerateVar(hyGTL)
	repeatScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateBody(node.Body, &repeatScope)

	repeatScope.ReplaceAll()

	repeatScope.AppendETabbed("::" + gotoLabel + "::\n")

	repeatScope.Append("end")

	scope.Write(repeatScope.Src)
}

func (gen *Generator) whileStmt(node ast.WhileStmt, scope *GenScope) {
	whileScope := NewGenScope(scope)
	whileScope.AppendTabbed("while ", gen.GenerateExpr(node.Condtion, &whileScope), " do\n")

	gotoLabel := GenerateVar(hyGTL)
	whileScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateBody(node.Body, &whileScope)

	whileScope.ReplaceAll()

	scope.Write(whileScope.Src)
	scope.AppendETabbed("::" + gotoLabel + "::\n")

	scope.AppendTabbed("end")
}

func (gen *Generator) forStmt(node ast.ForStmt, scope *GenScope) {
	forScope := NewGenScope(scope)

	forScope.AppendTabbed("for ")

	pairs := "pairs"
	if node.OrderedIteration {
		pairs = "ipairs"
	}
	iterator := gen.GenerateExpr(node.Iterator, scope)
	if len(node.KeyValuePair) == 1 {
		key := gen.GenerateExpr(node.KeyValuePair[0], &forScope)
		forScope.Append(key, ", _ in  ", pairs, " (", iterator, ") do\n")
	} else if len(node.KeyValuePair) == 2 {
		key := gen.GenerateExpr(node.KeyValuePair[0], &forScope)
		value := gen.GenerateExpr(node.KeyValuePair[1], &forScope)
		forScope.Append(key, ", ", value, " in ", pairs, "(", iterator, ") do\n")
	} else {
		forScope.Append("_, _ in ", pairs, "(", iterator, ") do\n")
	}
	gotoLabel := GenerateVar(hyGTL)
	forScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateBody(node.Body, &forScope)

	forScope.ReplaceAll()

	forScope.AppendETabbed("::" + gotoLabel + "::\n")

	forScope.AppendTabbed("end")

	scope.Write(forScope.Src)
}

func (gen *Generator) tickStmt(node ast.TickStmt, scope *GenScope) {
	tickTabs := getTabs()

	tickScope := GenScope{Src: StringBuilder{}, Parent: scope}

	if node.Variable.GetValueType() != ast.Unknown {
		variable := gen.GenerateExpr(&node.Variable, scope)
		tickScope.Src.Append(tickTabs, "local ", variable, " = 0\n")
		tickScope.Src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
		tickScope.Src.AppendTabbed(variable, " = ", variable, " + 1\n")
	} else {
		tickScope.Src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
	}

	gen.GenerateBody(node.Body, &tickScope)

	tickScope.Src.Append(tickTabs, "end)")

	scope.Write(tickScope.Src)
}

func (gen *Generator) variableDeclarationStmt(declaration ast.VariableDeclarationStmt, scope *GenScope) {
	var values []string

	for _, expr := range declaration.Values {
		values = append(values, gen.GenerateExpr(expr, scope))
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
			src.Append(fmt.Sprintf("%s = ", gen.WriteVar(ident.Lexeme)))
		} else if i == len(declaration.Identifiers)-1 {
			src.Append(gen.WriteVar(ident.Lexeme))
		} else {
			src.Append(fmt.Sprintf("%s, ", gen.WriteVar(ident.Lexeme)))
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

	scope.Write(src)
}

func (gen *Generator) enumDeclarationStmt(node ast.EnumDeclarationStmt, scope *GenScope) {
	if node.IsLocal {
		scope.AppendTabbed("local ")
	} else {
		scope.AppendTabbed()
	}

	scope.Append(gen.WriteVar(node.Name.Lexeme), " = {\n")

	length := len(node.Fields)
	for i := range node.Fields {
		if i == length-1 {
			scope.AppendETabbed(strconv.Itoa(i), "\n")
		} else {
			scope.AppendETabbed(strconv.Itoa(i), ", \n")
		}
	}
	scope.AppendTabbed("}")
}

func (gen *Generator) structDeclarationStmt(node ast.StructDeclarationStmt, scope *GenScope) {
	structScope := NewGenScope(scope)

	for _, nodebody := range node.Methods {
		gen.methodDeclarationStmt(nodebody, node, &structScope)
	}

	gen.constructorDeclarationStmt(*node.Constructor, node, &structScope)

	scope.Write(structScope.Src)
}

func (gen *Generator) entityDeclarationStmt(node ast.EntityDeclarationStmt, scope *GenScope) {
	entityScope := NewGenScope(scope)

	gen.spawnDeclarationStmt(*node.Spawner, node, &entityScope)
	gen.destroyDeclarationStmt(*node.Destroyer, node, &entityScope)

	scope.Write(entityScope.Src)
}

func (gen *Generator) spawnDeclarationStmt(node ast.SpawnDeclarationStmt, entity ast.EntityDeclarationStmt, scope *GenScope) {
	spawnScope :=  NewGenScope(scope)

	// if entity.IsLocal {
	// 	spawnScope.WriteString("local ")
	// }

	spawnScope.Append("HS_", gen.WriteVar(entity.Name.Lexeme), " = {}\n")

	// if entity.IsLocal {
	// 	spawnScope.WriteString("local ")
	// }

	spawnScope.Append("function Hy_", gen.WriteVar(entity.Name.Lexeme), "_Spawn(")
	gen.GenerateParams(node.Params, &spawnScope)

	spawnScope.Append(")\n")

	TabsCount++
	spawnScope.AppendTabbed("local instance = pewpew.new_customizable_entity(", node.Params[0].Name.Lexeme,", ", node.Params[1].Name.Lexeme,")\n")

	spawnScope.AppendTabbed("local Self = {")
	for i, v := range entity.Fields {
		gen.fieldDeclarationStmt(v, &spawnScope)
		if i != len(entity.Fields)-1 {
			spawnScope.WriteString(", ")
		}
		spawnScope.WriteString("\n")
	}
	spawnScope.WriteString("}\n")
	spawnScope.AppendTabbed("HS_", gen.WriteVar(entity.Name.Lexeme), "[instance] = Self\n\n")

	for i, v := range entity.Callbacks {
		spawnScope.AppendTabbed(fmt.Sprintf("local function HCb%v", i),"(")
		gen.GenerateParams(v.Params, &spawnScope)
		spawnScope.WriteString(")\n")
		gen.GenerateBody(v.Body, &spawnScope)
		spawnScope.AppendTabbed("end\n")
		if v.Callback == ast.WallCollision {
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.customizable_entity_configure_wall_collision(instance, true, HCb%v)\n", i))
		}else if v.Callback == ast.WeaponCollision {
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.customizable_entity_set_weapon_collision(instance, HCb%v)\n", i))
		}else if v.Callback == ast.PlayerCollision {
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.customizable_entity_set_player_collision(instance, HCb%v)\n", i))
		}
	}
	TabsCount--
	gen.GenerateBody(node.Body, &spawnScope)
	spawnScope.WriteString("end\n")

	scope.Write(spawnScope.Src)
} 

func (gen *Generator) destroyDeclarationStmt(node ast.DestroyDeclarationStmt, entity ast.EntityDeclarationStmt, scope *GenScope) {
	spawnScope :=  NewGenScope(scope)

	// if entity.IsLocal {
	// 	spawnScope.WriteString("local ")
	// }

	spawnScope.Append("function Hy_", gen.WriteVar(entity.Name.Lexeme), "_Destroy(")

	gen.GenerateParams(node.Params, &spawnScope)

	spawnScope.Append(")\n")

	gen.GenerateBody(node.Body, &spawnScope)

	spawnScope.WriteString("end\n")

	scope.Write(spawnScope.Src)
}

func (gen *Generator) GenerateParams(params []ast.Param, scope *GenScope) {
	for i, param := range params {
		scope.Append(gen.WriteVar(param.Name.Lexeme))
		if i != len(params)-1 {
			scope.Append(", ")
		}
	}
} 

func (gen *Generator) constructorDeclarationStmt(node ast.ConstructorStmt, Struct ast.StructDeclarationStmt, scope *GenScope) {
	src := StringBuilder{}

	constructorScope := NewGenScope(scope)

	// if Struct.IsLocal {
	// 	constructorScope.WriteString("local ")
	// }

	constructorScope.Append("function Hy_", gen.WriteVar(Struct.Name.Lexeme), "_New(")

	gen.GenerateParams(node.Params, &constructorScope)

	constructorScope.Append(")\n")

	TabsCount++
	src.AppendTabbed("local Self = {")
	for i, field := range Struct.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(Struct.Fields)-1 {
				src.WriteString(gen.GenerateExpr(value, &constructorScope))
			} else {
				src.Append(gen.GenerateExpr(value, &constructorScope), ",")
			}
		}
	}
	src.WriteString("}\n")
	constructorScope.Write(src)
	TabsCount--
	gen.GenerateBody(node.Body, &constructorScope)
	TabsCount++
	constructorScope.AppendTabbed("return Self\n")
	TabsCount--
	constructorScope.AppendTabbed("end\n")

	scope.Write(constructorScope.Src)
}

func (gen *Generator) fieldDeclarationStmt(node ast.FieldDeclarationStmt, scope *GenScope) string {
	src := StringBuilder{}

	for i, v := range node.Identifiers {
		src.Append(v.Lexeme, " = ", gen.GenerateExpr(node.Values[i], scope))
		if i != len(node.Identifiers)-1 {
			src.WriteString(", ")
		}
	}

	return src.String()
}

func (gen *Generator) methodDeclarationStmt(node ast.MethodDeclarationStmt, Struct ast.StructDeclarationStmt, scope *GenScope) {
	methodScope := NewGenScope(scope)

	// if Struct.IsLocal {
	// 	methodScope.WriteString("local ")
	// }

	methodScope.Append("function Hybroid_", Struct.Name.Lexeme, "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		methodScope.WriteString(", ")
		methodScope.Append(gen.WriteVar(param.Name.Lexeme))
	}
	methodScope.Append(")\n")

	gen.GenerateBody(node.Body, &methodScope) // its constructor

	methodScope.AppendTabbed("end\n")

	scope.Write(methodScope.Src)
}

func (gen *Generator) matchStmt(node ast.MatchStmt, scope *GenScope) {
	ifStmt := ast.IfStmt{
		BoolExpr: &ast.BinaryExpr{Left: node.ExprToMatch, Operator: lexer.Token{Type: lexer.EqualEqual, Lexeme: "=="}, Right: node.Cases[0].Expression},
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
			BoolExpr: &ast.BinaryExpr{Left: node.ExprToMatch, Operator: lexer.Token{Type: lexer.EqualEqual, Lexeme: "=="}, Right: node.Cases[i].Expression},
			Body:     node.Cases[i].Body,
		}
		ifStmt.Elseifs = append(ifStmt.Elseifs, &elseIfStmt)
	}

	if has_default {
		ifStmt.Else = &ast.IfStmt{
			Body: node.Cases[len(node.Cases)-1].Body,
		}
	}

	gen.ifStmt(ifStmt, scope)
}
