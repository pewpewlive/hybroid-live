package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/tokens"
	"strconv"
	"strings"
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

func (gen *Generator) assignmentStmt(assignStmt ast.AssignmentStmt, scope *GenScope) {
	src := StringBuilder{}

	preSrc := StringBuilder{}

	vars := []string{}

	index := 0
	for i := 0; i < len(assignStmt.Values); i++ {
		src.AppendTabbed()
		if assignStmt.Values[i].GetType() == ast.CallExpression {
			call := assignStmt.Values[i].(*ast.CallExpr)
			preSrc.AppendTabbed()
			for j := range call.ReturnAmount {
				src.Append(gen.GenerateExpr(assignStmt.Identifiers[index+j], scope))
				vars = append(vars, GenerateVar(hyVar))
				preSrc.WriteString(vars[j])
				if j != call.ReturnAmount-1 {
					preSrc.WriteString(", ")
					src.WriteString(", ")
				} else {
					preSrc.Append(" = ", gen.callExpr(*call, false, scope), "\n")
					src.Append(" = ", strings.Join(vars, ", "), "\n")
				}
			}
			vars = []string{}
			index += call.ReturnAmount
		}else {
			src.Append(gen.GenerateExpr(assignStmt.Identifiers[index], scope), " = ", gen.GenerateExpr(assignStmt.Values[i], scope), "\n")
			index++
		}
		if index >= len(assignStmt.Identifiers) {
			break
		}
	}

	scope.Write(preSrc)
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
	if node.Variable != nil {
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
	key := gen.GenerateExpr(node.First, &forScope)
	if node.Second == nil {
		forScope.Append(key, ", _ in  ", pairs, " (", iterator, ") do\n")
	}else {
		value := gen.GenerateExpr(node.Second, &forScope)
		forScope.Append(key, ", ", value, " in ", pairs, "(", iterator, ") do\n")
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

	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable, scope)
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

	src := StringBuilder{}
	src2 := StringBuilder{}
	if declaration.IsLocal {
		src.AppendTabbed("local ")
	} else {
		src.AppendTabbed("")
	}
	for i, ident := range declaration.Identifiers {
		if i == len(declaration.Identifiers)-1 && len(values) != 0 {
			src.Append(fmt.Sprintf("%s = ", gen.GenerateExpr(&ast.IdentifierExpr{Name: ident}, scope)))
		} else if i == len(declaration.Identifiers)-1 {
			src.Append(gen.GenerateExpr(&ast.IdentifierExpr{Name: ident}, scope))
		} else {
			src.Append(fmt.Sprintf("%s, ", gen.GenerateExpr(&ast.IdentifierExpr{Name: ident}, scope)))
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

func (gen *Generator) classDeclarationStmt(node ast.ClassDeclarationStmt, scope *GenScope) {
	classScope := NewGenScope(scope)

	for _, nodebody := range node.Methods {
		gen.methodDeclarationStmt(nodebody, node, &classScope)
	}

	gen.constructorDeclarationStmt(*node.Constructor, node, &classScope)

	scope.Write(classScope.Src)
}

func (gen *Generator) entityDeclarationStmt(node ast.EntityDeclarationStmt, scope *GenScope) {
	entityScope := NewGenScope(scope)

	entityName := gen.WriteVarExtra(node.Name.Lexeme, hyEntity)

	for i, v := range node.Callbacks {
		entityScope.AppendTabbed(fmt.Sprintf("local function %sHCb%d", entityName, i), "(id")
		if len(v.Params) != 0 {
			entityScope.Append(", ")
		}
		gen.GenerateParams(v.Params, &entityScope)
		gen.GenerateBody(v.Body, &entityScope)
		entityScope.AppendTabbed("end\n")
	}

	gen.spawnDeclarationStmt(*node.Spawner, node, &entityScope)
	gen.destroyDeclarationStmt(*node.Destroyer, node, &entityScope)

	for _, v := range node.Methods {
		gen.entityMethodDeclarationStmt(v, node, scope)
	}

	scope.Write(entityScope.Src)
}

func (gen *Generator) spawnDeclarationStmt(node ast.EntityFunctionDeclarationStmt, entity ast.EntityDeclarationStmt, scope *GenScope) {
	spawnScope := NewGenScope(scope)

	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	spawnScope.Append(entityName, " = {}\n")
	spawnScope.Append("function ", entityName, "_Spawn(")

	gen.GenerateParams(node.Params, &spawnScope)

	TabsCount++

	spawnScope.AppendTabbed("local id = pewpew.new_customizable_entity(", gen.WriteVar(node.Params[0].Name.Lexeme), ", ", gen.WriteVar(node.Params[1].Name.Lexeme), ")\n")
	spawnScope.AppendTabbed(entityName, "[id] = {")
	for i, field := range entity.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(entity.Fields)-1 {
				spawnScope.WriteString(gen.GenerateExpr(value, &spawnScope))
			} else {
				spawnScope.Append(gen.GenerateExpr(value, &spawnScope), ",")
			}
		}
	}
	spawnScope.Append("}\n")

	TabsCount--
	gen.GenerateBody(node.Body, &spawnScope)
	TabsCount++

	for i, v := range entity.Callbacks {
		switch v.Type {
		case ast.WallCollision:
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.customizable_entity_configure_wall_collision(id, true, %sHCb%d)\n", entityName, i))
		case ast.WeaponCollision:
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.customizable_entity_set_weapon_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.PlayerCollision:
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.customizable_entity_set_player_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.Update:
			spawnScope.AppendTabbed(fmt.Sprintf("pewpew.entity_set_update_callback(id, %sHCb%d)\n", entityName, i))
		}
	}
	spawnScope.AppendTabbed("return id\n")
	TabsCount--

	spawnScope.AppendTabbed("end\n")

	scope.Write(spawnScope.Src)
}

func (gen *Generator) destroyDeclarationStmt(node ast.EntityFunctionDeclarationStmt, entity ast.EntityDeclarationStmt, scope *GenScope) {
	spawnScope := NewGenScope(scope)

	spawnScope.Append("function ", gen.WriteVarExtra(entity.Name.Lexeme, hyEntity), "_Destroy(id")
	if len(node.Params) != 0 {
		spawnScope.Append(", ")
	}

	gen.GenerateParams(node.Params, &spawnScope)

	gen.GenerateBody(node.Body, &spawnScope)

	spawnScope.WriteString("end\n")

	scope.Write(spawnScope.Src)
}

func (gen *Generator) GenerateParams(params []ast.Param, scope *GenScope) {
	var variadicParam string
	for i, param := range params {
		if param.Type.IsVariadic {
			scope.WriteString("...")
			variadicParam = gen.WriteVar(param.Name.Lexeme)
		} else {
			scope.Append(gen.WriteVar(param.Name.Lexeme))
		}
		if i != len(params)-1 {
			scope.Append(", ")
		}
	}
	scope.Append(")\n")

	if variadicParam != "" {
		scope.AppendETabbed("local ", variadicParam, " = {...}\n")
	}
}

func (gen *Generator) constructorDeclarationStmt(node ast.ConstructorStmt, class ast.ClassDeclarationStmt, scope *GenScope) {
	src := StringBuilder{}

	constructorScope := NewGenScope(scope)

	constructorScope.Append("function ", gen.WriteVarExtra(class.Name.Lexeme, hyClass), "_New(")

	gen.GenerateParams(node.Params, &constructorScope)

	TabsCount++
	src.AppendTabbed("local Self = {")
	for i, field := range class.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(class.Fields)-1 {
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

func (gen *Generator) methodDeclarationStmt(node ast.MethodDeclarationStmt, Struct ast.ClassDeclarationStmt, scope *GenScope) {
	methodScope := NewGenScope(scope)

	methodScope.Append("function ", gen.WriteVarExtra(Struct.Name.Lexeme, hyClass), "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		methodScope.WriteString(", ")
		methodScope.Append(gen.WriteVar(param.Name.Lexeme))
	}
	methodScope.Append(")\n")

	gen.GenerateBody(node.Body, &methodScope)

	methodScope.AppendTabbed("end\n")

	scope.Write(methodScope.Src)
}

func (gen *Generator) entityMethodDeclarationStmt(node ast.MethodDeclarationStmt, entity ast.EntityDeclarationStmt, scope *GenScope) {
	methodScope := NewGenScope(scope)

	methodScope.Append("function ", gen.WriteVarExtra(entity.Name.Lexeme, hyEntity), "_", node.Name.Lexeme, "(id")
	for _, param := range node.Params {
		methodScope.WriteString(", ")
		methodScope.Append(gen.WriteVar(param.Name.Lexeme))
	}
	methodScope.Append(")\n")

	gen.GenerateBody(node.Body, &methodScope)

	methodScope.AppendTabbed("end\n")

	scope.Write(methodScope.Src)
}

func (gen *Generator) matchStmt(node ast.MatchStmt, scope *GenScope) {
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

	gen.ifStmt(ifStmt, scope)
}

func (gen *Generator) destroyStmt(node ast.DestroyStmt, scope *GenScope) {
	src := StringBuilder{}

	src.AppendTabbed()
	src.Append(envMap[node.EnvName], hyEntity, node.EntityName, "_Destroy(", gen.GenerateExpr(node.Identifier, scope))
	for _, arg := range node.Args {
		src.WriteString(", ")
		src.WriteString(gen.GenerateExpr(arg, scope))
	}
	src.WriteString(")")
	scope.Write(src)
}
