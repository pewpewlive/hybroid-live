package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
	"strconv"
	"strings"
)

func (gen *Generator) envStmt(node ast.EnvironmentDecl, scope *GenScope) {
	for i := range node.Requirements {
		scope.Write("require(\"", node.Requirements[i], "\")\n")
	}
}

func (gen *Generator) ifStmt(node ast.IfStmt, scope *GenScope) {
	ifScope := NewGenScope(scope)

	ifScope.WriteTabbed("if ", gen.GenerateExpr(node.BoolExpr, scope), " then\n")

	gen.GenerateBody(node.Body, &ifScope)
	for _, elseif := range node.Elseifs {
		ifScope.WriteTabbed("elseif ", gen.GenerateExpr(elseif.BoolExpr, scope), " then\n")
		gen.GenerateBody(elseif.Body, &ifScope)
	}
	if node.Else != nil {
		ifScope.WriteTabbed("else \n")
		gen.GenerateBody(node.Else.Body, &ifScope)
	}

	ifScope.WriteTabbed("end\n")

	ifScope.ReplaceAll()

	scope.Write(ifScope.String())
}

func (gen *Generator) assignmentStmt(assignStmt ast.AssignmentStmt, scope *GenScope) {
	src := StringBuilder{}

	preSrc := StringBuilder{}

	vars := []string{}

	index := 0
	for i := 0; i < len(assignStmt.Values); i++ {
		src.WriteTabbed()
		if assignStmt.Values[i].GetType() == ast.CallExpression {
			call := assignStmt.Values[i].(*ast.CallExpr)
			preSrc.WriteTabbed()
			for j := range call.ReturnAmount {
				src.Write(gen.GenerateExpr(assignStmt.Identifiers[index+j], scope))
				vars = append(vars, GenerateVar(hyVar))
				preSrc.Write(vars[j])
				if j != call.ReturnAmount-1 {
					preSrc.Write(", ")
					src.Write(", ")
				} else {
					preSrc.Write(" = ", gen.callExpr(*call, false, scope), "\n")
					src.Write(" = ", strings.Join(vars, ", "), "\n")
				}
			}
			vars = []string{}
			index += call.ReturnAmount
		} else {
			src.Write(gen.GenerateExpr(assignStmt.Identifiers[index], scope), " = ", gen.GenerateExpr(assignStmt.Values[i], scope), "\n")
			index++
		}
		if index >= len(assignStmt.Identifiers) {
			break
		}
	}

	scope.Write(preSrc.String())
	scope.Write(src.String())
}

func (gen *Generator) functionDeclarationStmt(node ast.FunctionDecl, scope *GenScope) {
	fnScope := NewGenScope(scope)

	if !node.IsPub {
		fnScope.WriteTabbed("local ")
	}

	fnScope.Write("function ", gen.WriteVar(node.Name.Lexeme), "(")
	gen.GenerateParams(node.Params, &fnScope)

	gen.GenerateBody(node.Body, &fnScope)

	fnScope.WriteTabbed("end")

	scope.Write(fnScope.String())
}

func (gen *Generator) returnStmt(node ast.ReturnStmt, scope *GenScope) {
	src := StringBuilder{}

	src.WriteTabbed("return ")
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr, scope)
		src.Write(val)
		if i != len(node.Args)-1 {
			src.Write(", ")
		}
	}

	scope.Write(src.String())
}

func (gen *Generator) yieldStmt(node ast.YieldStmt, scope *GenScope) {
	src := StringBuilder{}

	src.WriteTabbed()
	startIndex := src.Len()
	src.Write("yield ")
	endIndex := src.Len()
	for i, expr := range node.Args {
		val := gen.GenerateExpr(expr, scope)
		src.Write(val)
		if i != len(node.Args)-1 {
			src.Write(", ")
		}
	}

	src.Write("\n")
	src.WriteTabbed()
	startIndex2 := src.Len()
	src.Write("goto hyl")
	endIndex2 := src.Len()
	src.Write("\n")

	scopeLength := scope.Len()

	scope.Write(src.String())

	scope.AddReplacement(YieldReplacement, helpers.NewSpan(startIndex+scopeLength, endIndex+scopeLength))
	scope.AddReplacement(GotoReplacement, helpers.NewSpan(startIndex2+scopeLength, endIndex2+scopeLength))
}

func (gen *Generator) breakStmt(_ ast.BreakStmt, scope *GenScope) {
	scope.WriteTabbed("break")
}

func (gen *Generator) continueStmt(_ ast.ContinueStmt, scope *GenScope) {
	src := StringBuilder{}

	src.WriteTabbed()
	startIndex := src.Len()
	src.Write("continue")
	endIndex := src.Len()

	scopeLength := scope.Len()

	scope.Write(src.String())

	scope.AddReplacement(ContinueReplacement, helpers.NewSpan(startIndex+scopeLength, endIndex+scopeLength))
}

func (gen *Generator) repeatStmt(node ast.RepeatStmt, scope *GenScope) {
	repeatScope := NewGenScope(scope)

	end := gen.GenerateExpr(node.Iterator, scope)
	start := gen.GenerateExpr(node.Start, scope)
	skip := gen.GenerateExpr(node.Skip, scope)
	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable, &repeatScope)
		repeatScope.WriteTabbed("for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		repeatScope.WriteTabbed("for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	gotoLabel := GenerateVar(hyGTL)
	repeatScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateBody(node.Body, &repeatScope)

	repeatScope.ReplaceAll()

	repeatScope.WriteTabbed("::" + gotoLabel + "::\n")

	repeatScope.Write("end")

	scope.Write(repeatScope.String())
}

func (gen *Generator) whileStmt(node ast.WhileStmt, scope *GenScope) {
	whileScope := NewGenScope(scope)
	whileScope.WriteTabbed("while ", gen.GenerateExpr(node.Condition, &whileScope), " do\n")

	gotoLabel := GenerateVar(hyGTL)
	whileScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateBody(node.Body, &whileScope)

	whileScope.ReplaceAll()

	scope.Write(whileScope.String())
	scope.WriteTabbed("::" + gotoLabel + "::\n")

	scope.WriteTabbed("end")
}

func (gen *Generator) forStmt(node ast.ForStmt, scope *GenScope) {
	forScope := NewGenScope(scope)

	forScope.WriteTabbed("for ")

	pairs := "pairs"
	if node.OrderedIteration {
		pairs = "ipairs"
	}
	iterator := gen.GenerateExpr(node.Iterator, scope)
	key := gen.GenerateExpr(node.First, &forScope)
	if node.Second == nil {
		forScope.Write(key, ", _ in  ", pairs, " (", iterator, ") do\n")
	} else {
		value := gen.GenerateExpr(node.Second, &forScope)
		forScope.Write(key, ", ", value, " in ", pairs, "(", iterator, ") do\n")
	}
	gotoLabel := GenerateVar(hyGTL)
	forScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateBody(node.Body, &forScope)

	forScope.ReplaceAll()

	forScope.WriteTabbed("::" + gotoLabel + "::\n")

	forScope.WriteTabbed("end")

	scope.Write(forScope.String())
}

func (gen *Generator) tickStmt(node ast.TickStmt, scope *GenScope) {
	tickTabs := getTabs()

	tickScope := NewGenScope(scope)

	if node.Variable != nil {
		variable := gen.GenerateExpr(node.Variable, scope)
		tickScope.Write(tickTabs, "local ", variable, " = 0\n")
		tickScope.Write(tickTabs, "pewpew.add_update_callback(function()\n")
		tickScope.WriteTabbed(variable, " = ", variable, " + 1\n")
	} else {
		tickScope.Write(tickTabs, "pewpew.add_update_callback(function()\n")
	}

	gen.GenerateBody(node.Body, &tickScope)

	tickScope.Write(tickTabs, "end)")

	scope.Write(tickScope.String())
}

func (gen *Generator) variableDeclarationStmt(declaration ast.VariableDecl, scope *GenScope) {
	var values []string

	for _, expr := range declaration.Expressions {
		values = append(values, gen.GenerateExpr(expr, scope))
	}

	src := StringBuilder{}
	src2 := StringBuilder{}
	if !declaration.IsPub {
		src.WriteTabbed("local ")
	}
	for i, ident := range declaration.Identifiers {
		if i == len(declaration.Identifiers)-1 && len(values) != 0 {
			src.Write(fmt.Sprintf("%s = ", gen.GenerateExpr(ident, scope)))
		} else if i == len(declaration.Identifiers)-1 {
			src.Write(gen.GenerateExpr(ident, scope))
		} else {
			src.Write(fmt.Sprintf("%s, ", gen.GenerateExpr(ident, scope)))
		}
	}
	for i := range values {
		if i == len(values)-1 {
			src2.Write(values[i])
			break
		}
		src2.Write(fmt.Sprintf("%s, ", values[i]))

	}

	src.Write(src2.String(), "\n")

	scope.Write(src.String())
}

func (gen *Generator) enumDeclarationStmt(node ast.EnumDecl, scope *GenScope) {
	if node.IsPub {
		scope.WriteTabbed("local ")
	} else {
		scope.WriteTabbed()
	}

	scope.Write(gen.WriteVar(node.Name.Lexeme), " = {\n")

	length := len(node.Fields)
	for i := range node.Fields {
		if i == length-1 {
			scope.WriteTabbed(strconv.Itoa(i), "\n")
		} else {
			scope.WriteTabbed(strconv.Itoa(i), ", \n")
		}
	}
	scope.WriteTabbed("}")
}

func (gen *Generator) classDeclarationStmt(node ast.ClassDecl, scope *GenScope) {
	classScope := NewGenScope(scope)

	for _, nodebody := range node.Methods {
		gen.methodDeclarationStmt(nodebody, node, &classScope)
	}

	gen.constructorDeclarationStmt(*node.Constructor, node, &classScope)

	scope.Write(classScope.String())
}

func (gen *Generator) entityDeclarationStmt(node ast.EntityDecl, scope *GenScope) {
	entityScope := NewGenScope(scope)

	entityName := gen.WriteVarExtra(node.Name.Lexeme, hyEntity)

	for i, v := range node.Callbacks {
		entityScope.WriteTabbed(fmt.Sprintf("local function %sHCb%d", entityName, i), "(id")
		if len(v.Params) != 0 {
			entityScope.Write(", ")
		}
		gen.GenerateParams(v.Params, &entityScope)
		gen.GenerateBody(v.Body, &entityScope)
		entityScope.WriteTabbed("end\n")
	}

	gen.spawnDeclarationStmt(*node.Spawner, node, &entityScope)
	gen.destroyDeclarationStmt(*node.Destroyer, node, &entityScope)

	for _, v := range node.Methods {
		gen.entityMethodDeclarationStmt(v, node, scope)
	}

	scope.Write(entityScope.String())
}

func (gen *Generator) spawnDeclarationStmt(node ast.EntityFunctionDecl, entity ast.EntityDecl, scope *GenScope) {
	spawnScope := NewGenScope(scope)

	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	spawnScope.Write(entityName, " = {}\n")
	spawnScope.Write("function ", entityName, "_Spawn(")

	gen.GenerateParams(node.Params, &spawnScope)

	TabsCount++

	spawnScope.WriteTabbed("local id = pewpew.new_customizable_entity(", gen.WriteVar(node.Params[0].Name.Lexeme), ", ", gen.WriteVar(node.Params[1].Name.Lexeme), ")\n")
	spawnScope.WriteTabbed(entityName, "[id] = {")
	for i, field := range entity.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(entity.Fields)-1 {
				spawnScope.WriteString(gen.GenerateExpr(value, &spawnScope))
			} else {
				spawnScope.Write(gen.GenerateExpr(value, &spawnScope), ",")
			}
		}
	}
	spawnScope.Write("}\n")

	TabsCount--
	gen.GenerateBody(node.Body, &spawnScope)
	TabsCount++

	for i, v := range entity.Callbacks {
		switch v.Type {
		case ast.WallCollision:
			spawnScope.WriteTabbed(fmt.Sprintf("pewpew.customizable_entity_configure_wall_collision(id, true, %sHCb%d)\n", entityName, i))
		case ast.WeaponCollision:
			spawnScope.WriteTabbed(fmt.Sprintf("pewpew.customizable_entity_set_weapon_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.PlayerCollision:
			spawnScope.WriteTabbed(fmt.Sprintf("pewpew.customizable_entity_set_player_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.Update:
			spawnScope.WriteTabbed(fmt.Sprintf("pewpew.entity_set_update_callback(id, %sHCb%d)\n", entityName, i))
		}
	}
	spawnScope.WriteTabbed("return id\n")
	TabsCount--

	spawnScope.WriteTabbed("end\n")

	scope.Write(spawnScope.String())
}

func (gen *Generator) destroyDeclarationStmt(node ast.EntityFunctionDecl, entity ast.EntityDecl, scope *GenScope) {
	spawnScope := NewGenScope(scope)

	spawnScope.Write("function ", gen.WriteVarExtra(entity.Name.Lexeme, hyEntity), "_Destroy(id")
	if len(node.Params) != 0 {
		spawnScope.Write(", ")
	}

	gen.GenerateParams(node.Params, &spawnScope)

	gen.GenerateBody(node.Body, &spawnScope)

	spawnScope.WriteString("end\n")

	scope.Write(spawnScope.String())
}

func (gen *Generator) GenerateParams(params []ast.FunctionParam, scope *GenScope) {
	var variadicParam string
	for i, param := range params {
		if param.Type.IsVariadic {
			scope.WriteString("...")
			variadicParam = gen.WriteVar(param.Name.Lexeme)
		} else {
			scope.Write(gen.WriteVar(param.Name.Lexeme))
		}
		if i != len(params)-1 {
			scope.Write(", ")
		}
	}
	scope.Write(")\n")

	if variadicParam != "" {
		scope.WriteTabbed("local ", variadicParam, " = {...}\n")
	}
}

func (gen *Generator) constructorDeclarationStmt(node ast.ConstructorDecl, class ast.ClassDecl, scope *GenScope) {
	src := StringBuilder{}

	constructorScope := NewGenScope(scope)

	constructorScope.Write("function ", gen.WriteVarExtra(class.Name.Lexeme, hyClass), "_New(")

	gen.GenerateParams(node.Params, &constructorScope)

	TabsCount++
	src.WriteTabbed("local Self = {")
	for i, field := range class.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(class.Fields)-1 {
				src.Write(gen.GenerateExpr(value, &constructorScope))
			} else {
				src.Write(gen.GenerateExpr(value, &constructorScope), ",")
			}
		}
	}
	src.Write("}\n")
	constructorScope.Write(src.String())
	TabsCount--
	gen.GenerateBody(node.Body, &constructorScope)
	TabsCount++
	constructorScope.WriteTabbed("return Self\n")
	TabsCount--
	constructorScope.WriteTabbed("end\n")

	scope.Write(constructorScope.String())
}

func (gen *Generator) fieldDeclarationStmt(node ast.FieldDecl, scope *GenScope) string {
	src := StringBuilder{}

	// for i, v := range node.Identifiers {
	// 	src.Write(v.Name.Lexeme, " = ", gen.GenerateExpr(node.Values[i], scope))
	// 	if i != len(node.Identifiers)-1 {
	// 		src.Write(", ")
	// 	}
	// }

	return src.String()
}

func (gen *Generator) methodDeclarationStmt(node ast.MethodDecl, Struct ast.ClassDecl, scope *GenScope) {
	methodScope := NewGenScope(scope)

	methodScope.Write("function ", gen.WriteVarExtra(Struct.Name.Lexeme, hyClass), "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		methodScope.WriteString(", ")
		methodScope.Write(gen.WriteVar(param.Name.Lexeme))
	}
	methodScope.Write(")\n")

	gen.GenerateBody(node.Body, &methodScope)

	methodScope.WriteTabbed("end\n")

	scope.Write(methodScope.String())
}

func (gen *Generator) entityMethodDeclarationStmt(node ast.MethodDecl, entity ast.EntityDecl, scope *GenScope) {
	methodScope := NewGenScope(scope)

	methodScope.Write("function ", gen.WriteVarExtra(entity.Name.Lexeme, hyEntity), "_", node.Name.Lexeme, "(id")
	for _, param := range node.Params {
		methodScope.WriteString(", ")
		methodScope.Write(gen.WriteVar(param.Name.Lexeme))
	}
	methodScope.Write(")\n")

	gen.GenerateBody(node.Body, &methodScope)

	methodScope.WriteTabbed("end\n")

	scope.Write(methodScope.String())
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

	src.WriteTabbed()
	src.Write(envMap[node.EnvName], hyEntity, node.EntityName, "_Destroy(", gen.GenerateExpr(node.Identifier, scope))
	for _, arg := range node.Args {
		src.Write(", ")
		src.Write(gen.GenerateExpr(arg, scope))
	}
	src.Write(")")
	scope.Write(src.String())
}
