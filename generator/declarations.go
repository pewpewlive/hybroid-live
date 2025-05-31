package generator

import (
	"fmt"
	"hybroid/ast"
	"strconv"
)

func (gen *Generator) breakDownVariableDeclaration(declaration ast.VariableDecl) []ast.VariableDecl {
	emptyVarDecl := func() ast.VariableDecl {
		return ast.VariableDecl{
			Identifiers: []*ast.IdentifierExpr{},
			Expressions: []ast.Node{},
			IsPub:       declaration.IsPub,
			IsConst:     declaration.IsConst,
		}
	}
	decls := []ast.VariableDecl{}
	currentDeclIndex := -1
	for _, expr := range declaration.Expressions {
		if call, ok := expr.(*ast.CallExpr); ok && call.ReturnAmount > 1 {
			decls = append(decls, emptyVarDecl())
			currentDeclIndex = len(decls) - 1
			for range call.ReturnAmount {
				decls[currentDeclIndex].Identifiers = append(decls[currentDeclIndex].Identifiers, declaration.Identifiers[0])
				declaration.Identifiers = declaration.Identifiers[1:]
			}
			decls[currentDeclIndex].Expressions = append(decls[currentDeclIndex].Expressions, expr)
			currentDeclIndex = -1
			continue
		}
		if currentDeclIndex == -1 {
			decls = append(decls, emptyVarDecl())
			currentDeclIndex = len(decls) - 1
		}
		decls[currentDeclIndex].Expressions = append(decls[currentDeclIndex].Expressions, expr)
		decls[currentDeclIndex].Identifiers = append(decls[currentDeclIndex].Identifiers, declaration.Identifiers[0])
		declaration.Identifiers = declaration.Identifiers[1:]
	}

	return decls
}

func (gen *Generator) variableDeclaration(declaration ast.VariableDecl) {
	if declaration.IsConst {
		return
	}
	var values []string

	for _, expr := range declaration.Expressions {
		values = append(values, gen.GenerateExpr(expr))
	}

	src := StringBuilder{}
	src2 := StringBuilder{}

	src.WriteTabbed()
	if !declaration.IsPub {
		src.Write("local ")
	}
	for i, ident := range declaration.Identifiers {
		if i == len(declaration.Identifiers)-1 && len(values) != 0 {
			src.Write(fmt.Sprintf("%s = ", gen.GenerateExpr(ident)))
		} else if i == len(declaration.Identifiers)-1 {
			src.Write(gen.GenerateExpr(ident))
		} else {
			src.Write(fmt.Sprintf("%s, ", gen.GenerateExpr(ident)))
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

	gen.Write(src.String())
}

func (gen *Generator) enumDeclaration(node ast.EnumDecl) {
	if node.IsPub {
		gen.WriteTabbed("local ")
	} else {
		gen.WriteTabbed()
	}

	gen.Write(gen.WriteVar(node.Name.Lexeme), " = {\n")

	length := len(node.Fields)
	for i := range node.Fields {
		if i == length-1 {
			gen.WriteTabbed(strconv.Itoa(i), "\n")
		} else {
			gen.WriteTabbed(strconv.Itoa(i), ", \n")
		}
	}
	gen.WriteTabbed("}")
}

func (gen *Generator) classDeclaration(node ast.ClassDecl) {
	for _, nodebody := range node.Methods {
		gen.methodDeclaration(nodebody, node)
	}

	gen.constructorDeclaration(*node.Constructor, node)

	gen.Write(gen.String())
}

func (gen *Generator) entityDeclaration(node ast.EntityDecl) {
	entityName := gen.WriteVarExtra(node.Name.Lexeme, hyEntity)

	for i, v := range node.Callbacks {
		gen.WriteTabbed(fmt.Sprintf("local function %sHCb%d", entityName, i), "(id")
		if len(v.Params) != 0 {
			gen.Write(", ")
		}
		gen.GenerateParams(v.Params)
		gen.GenerateBody(v.Body)
		gen.WriteTabbed("end\n")
	}

	gen.spawnDeclaration(*node.Spawner, node)
	gen.destroyDeclaration(*node.Destroyer, node)

	for _, v := range node.Methods {
		gen.entityFunctionDeclaration(v, node)
	}

	gen.Write(gen.String())
}

func (gen *Generator) constructorDeclaration(node ast.ConstructorDecl, class ast.ClassDecl) {
	src := StringBuilder{}

	gen.Write("function ", gen.WriteVarExtra(class.Name.Lexeme, hyClass), "_New(")

	gen.GenerateParams(node.Params)

	TabsCount++
	src.WriteTabbed("local Self = {")
	for i, field := range class.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(class.Fields)-1 {
				src.Write(gen.GenerateExpr(value))
			} else {
				src.Write(gen.GenerateExpr(value), ",")
			}
		}
	}
	src.Write("}\n")
	gen.Write(src.String())
	TabsCount--
	gen.GenerateBody(node.Body)
	TabsCount++
	gen.WriteTabbed("return Self\n")
	TabsCount--
	gen.WriteTabbed("end\n")

	gen.Write(gen.String())
}

func (gen *Generator) fieldDeclaration(node ast.FieldDecl) string {
	src := StringBuilder{}

	// for i, v := range node.Identifiers {
	// 	src.Write(v.Name.Lexeme, " = ", gen.GenerateExpr(node.Values[i]))
	// 	if i != len(node.Identifiers)-1 {
	// 		src.Write(", ")
	// 	}
	// }

	return src.String()
}

func (gen *Generator) methodDeclaration(node ast.MethodDecl, Struct ast.ClassDecl) {
	gen.Write("function ", gen.WriteVarExtra(Struct.Name.Lexeme, hyClass), "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		gen.WriteString(", ")
		gen.Write(gen.WriteVar(param.Name.Lexeme))
	}
	gen.Write(")\n")

	gen.GenerateBody(node.Body)

	gen.WriteTabbed("end\n")

	gen.Write(gen.String())
}

func (gen *Generator) entityFunctionDeclaration(node ast.MethodDecl, entity ast.EntityDecl) {
	gen.Write("function ", gen.WriteVarExtra(entity.Name.Lexeme, hyEntity), "_", node.Name.Lexeme, "(id")
	for _, param := range node.Params {
		gen.WriteString(", ")
		gen.Write(gen.WriteVar(param.Name.Lexeme))
	}
	gen.Write(")\n")

	gen.GenerateBody(node.Body)

	gen.WriteTabbed("end\n")

	gen.Write(gen.String())
}

func (gen *Generator) spawnDeclaration(node ast.EntityFunctionDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	gen.Write(entityName, " = {}\n")
	gen.Write("function ", entityName, "_Spawn(")

	gen.GenerateParams(node.Params)

	TabsCount++

	gen.WriteTabbed("local id = pewpew.new_customizable_entity(", gen.WriteVar(node.Params[0].Name.Lexeme), ", ", gen.WriteVar(node.Params[1].Name.Lexeme), ")\n")
	gen.WriteTabbed(entityName, "[id] = {")
	for i, field := range entity.Fields {
		for j, value := range field.Values {
			if j == len(field.Values)-1 && i == len(entity.Fields)-1 {
				gen.WriteString(gen.GenerateExpr(value))
			} else {
				gen.Write(gen.GenerateExpr(value), ",")
			}
		}
	}
	gen.Write("}\n")

	TabsCount--
	gen.GenerateBody(node.Body)
	TabsCount++

	for i, v := range entity.Callbacks {
		switch v.Type {
		case ast.WallCollision:
			gen.WriteTabbed(fmt.Sprintf("pewpew.customizable_entity_configure_wall_collision(id, true, %sHCb%d)\n", entityName, i))
		case ast.WeaponCollision:
			gen.WriteTabbed(fmt.Sprintf("pewpew.customizable_entity_set_weapon_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.PlayerCollision:
			gen.WriteTabbed(fmt.Sprintf("pewpew.customizable_entity_set_player_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.Update:
			gen.WriteTabbed(fmt.Sprintf("pewpew.entity_set_update_callback(id, %sHCb%d)\n", entityName, i))
		}
	}
	gen.WriteTabbed("return id\n")
	TabsCount--

	gen.WriteTabbed("end\n")

	gen.Write(gen.String())
}

func (gen *Generator) destroyDeclaration(node ast.EntityFunctionDecl, entity ast.EntityDecl) {
	gen.Write("function ", gen.WriteVarExtra(entity.Name.Lexeme, hyEntity), "_Destroy(id")
	if len(node.Params) != 0 {
		gen.Write(", ")
	}

	gen.GenerateParams(node.Params)

	gen.GenerateBody(node.Body)

	gen.WriteString("end\n")
	gen.Write(gen.String())
}

func (gen *Generator) functionDeclaration(node ast.FunctionDecl) {
	if !node.IsPub {
		gen.WriteTabbed("local ")
	}

	gen.Write("function ", gen.WriteVar(node.Name.Lexeme), "(")
	gen.GenerateParams(node.Params)

	gen.GenerateBody(node.Body)

	gen.WriteTabbed("end")
}
