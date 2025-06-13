package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/core"
)

func (gen *Generator) variableDeclaration(declaration ast.VariableDecl) {
	if declaration.IsConst {
		return
	}
	var values []string

	for _, expr := range declaration.Expressions {
		values = append(values, gen.GenerateExpr(expr))
	}

	src := core.StringBuilder{}
	src2 := core.StringBuilder{}

	src.Write(gen.tabString())
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

	src.Write(src2.String())
	gen.Write(src.String())
}

func (gen *Generator) classDeclaration(node ast.ClassDecl) {
	for _, nodebody := range node.Methods {
		gen.methodDeclaration(nodebody, node)
		gen.Write("\n")
	}

	totalFieldDecls := make([]ast.VariableDecl, 0)
	for i := range node.Fields {
		fieldDecls := gen.breakDownVariableDeclaration(node.Fields[i])
		totalFieldDecls = append(totalFieldDecls, fieldDecls...)
	}
	node.Fields = totalFieldDecls

	gen.constructorDeclaration(*node.Constructor, node)
}

func (gen *Generator) entityDeclaration(node ast.EntityDecl) {
	entityName := gen.WriteVarExtra(node.Name.Lexeme, hyEntity)

	gen.Write(entityName, " = {}\n")
	for i, v := range node.Callbacks {
		gen.Twrite(fmt.Sprintf("local function %sHCb%d", entityName, i), "(id")
		if len(v.Params) != 0 {
			gen.Write(", ")
		}
		gen.GenerateParams(v.Params)
		gen.tabCount++
		gen.Twrite("local Self = ", entityName, "[id]\n")
		gen.tabCount--
		gen.GenerateBody(v.Body)
		gen.Twrite("end\n")
	}

	totalFieldDecls := make([]ast.VariableDecl, 0)
	for i := range node.Fields {
		fieldDecls := gen.breakDownVariableDeclaration(node.Fields[i])
		totalFieldDecls = append(totalFieldDecls, fieldDecls...)
	}
	node.Fields = totalFieldDecls

	gen.spawnDeclaration(*node.Spawner, node)
	gen.Write("\n")
	gen.destroyDeclaration(*node.Destroyer, node)

	for _, v := range node.Methods {
		gen.Write("\n")
		gen.entityFunctionDeclaration(v, node)
	}
	gen.Write("\n")
	gen.Twrite("local function check() for k in pairs(", entityName, ") do if not pewpew.entity_get_is_alive(k) then ", entityName, "[k] = nil end end end\n")
	gen.Twrite("pewpew.add_update_callback(check)")
}

func (gen *Generator) constructorDeclaration(node ast.ConstructorDecl, class ast.ClassDecl) {
	gen.Write("function ", gen.WriteVarExtra(class.Name.Lexeme, hyClass), "_New(")

	gen.GenerateParams(node.Params)

	gen.tabCount++
	gen.Twrite("local Self = {}\n")
	counter := 1
	for _, fieldDecl := range class.Fields {
		gen.fieldDeclaration(fieldDecl, counter)
		gen.Write("\n")
		counter += len(fieldDecl.Identifiers)
	}
	gen.tabCount--
	gen.GenerateBody(node.Body)
	gen.tabCount++
	gen.Twrite("return Self\n")
	gen.tabCount--
	gen.Twrite("end")
}

func (gen *Generator) fieldDeclaration(node ast.VariableDecl, index int) {
	src := core.StringBuilder{}

	src.Write(gen.tabString())
	for i := range node.Identifiers {
		src.Write(fmt.Sprintf("Self[%v]", index+i))
		if i != len(node.Identifiers)-1 {
			src.Write(", ")
		}
	}
	src.Write(" = ")
	for i, v := range node.Expressions {
		src.Write(gen.GenerateExpr(v))
		if i != len(node.Expressions)-1 {
			src.Write(", ")
		}
	}
	gen.Write(src.String())
}

func (gen *Generator) methodDeclaration(node ast.MethodDecl, Struct ast.ClassDecl) {
	gen.Write("function ", gen.WriteVarExtra(Struct.Name.Lexeme, hyClass), "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		gen.Write(", ")
		gen.Write(gen.WriteVar(param.Name.Lexeme))
	}
	gen.Write(")\n")

	gen.GenerateBody(node.Body)

	gen.Twrite("end")
}

func (gen *Generator) entityFunctionDeclaration(node ast.MethodDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	gen.Write("function ", entityName, "_", node.Name.Lexeme, "(id")
	for _, param := range node.Params {
		gen.Write(", ")
		gen.Write(gen.WriteVar(param.Name.Lexeme))
	}
	gen.Write(")\n")
	gen.tabCount++
	gen.Twrite("local Self = ", entityName, "[id]\n")
	gen.tabCount--

	gen.GenerateBody(node.Body)

	gen.Twrite("end")
}

func (gen *Generator) spawnDeclaration(node ast.EntityFunctionDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	gen.Write("function ", entityName, "_Spawn(")

	gen.GenerateParams(node.Params)

	gen.tabCount++

	gen.Twrite("local id = pewpew.new_customizable_entity(", gen.WriteVar(node.Params[0].Name.Lexeme), ", ", gen.WriteVar(node.Params[1].Name.Lexeme), ")\n")
	tableAccess := entityName + "[id]"
	gen.Twrite(tableAccess, " = {}\n")
	gen.Twrite("local Self = ", tableAccess, "\n")
	counter := 1
	for _, field := range entity.Fields {
		gen.fieldDeclaration(field, counter)
		gen.Write("\n")
		counter += len(field.Identifiers)
	}

	gen.tabCount--
	gen.GenerateBody(node.Body)
	gen.tabCount++

	for i, v := range entity.Callbacks {
		switch v.Type {
		case ast.WallCollision:
			gen.Twrite(fmt.Sprintf("pewpew.customizable_entity_configure_wall_collision(id, true, %sHCb%d)\n", entityName, i))
		case ast.WeaponCollision:
			gen.Twrite(fmt.Sprintf("pewpew.customizable_entity_set_weapon_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.PlayerCollision:
			gen.Twrite(fmt.Sprintf("pewpew.customizable_entity_set_player_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.Update:
			gen.Twrite(fmt.Sprintf("pewpew.entity_set_update_callback(id, %sHCb%d)\n", entityName, i))
		}
	}
	gen.Twrite("return id\n")
	gen.tabCount--

	gen.Twrite("end")
}

func (gen *Generator) destroyDeclaration(node ast.EntityFunctionDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)
	gen.Write("function ", entityName, "_Destroy(id")
	if len(node.Params) != 0 {
		gen.Write(", ")
	}
	gen.GenerateParams(node.Params)
	gen.tabCount++
	gen.Twrite("local Self = ", entityName, "[id]\n")
	gen.tabCount--

	gen.GenerateBody(node.Body)

	gen.Write("end")
}

func (gen *Generator) functionDeclaration(node ast.FunctionDecl) {
	if !node.IsPub {
		gen.Twrite("local ")
	}

	gen.Write("function ", gen.WriteVar(node.Name.Lexeme), "(")
	gen.GenerateParams(node.Params)

	gen.GenerateBody(node.Body)

	gen.Twrite("end")
}
