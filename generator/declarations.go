package generator

import (
	"fmt"
	"hybroid/ast"
)

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
		if i != 0 {
			gen.Write("\n")
		}
		gen.WriteTabbed(fmt.Sprintf("local function %sHCb%d", entityName, i), "(id")
		if len(v.Params) != 0 {
			gen.Write(", ")
		}
		gen.GenerateParams(v.Params)
		TabsCount++
		gen.WriteTabbed("local Self = ", entityName, "[id]\n")
		TabsCount--
		gen.GenerateBody(v.Body)
		gen.WriteTabbed("end\n")
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
}

func (gen *Generator) constructorDeclaration(node ast.ConstructorDecl, class ast.ClassDecl) {
	gen.Write("function ", gen.WriteVarExtra(class.Name.Lexeme, hyClass), "_New(")

	gen.GenerateParams(node.Params)

	TabsCount++
	gen.WriteTabbed("local Self = {}\n")
	counter := 1
	for _, fieldDecl := range class.Fields {
		gen.fieldDeclaration(fieldDecl, counter)
		gen.Write("\n")
		counter += len(fieldDecl.Identifiers)
	}
	TabsCount--
	gen.GenerateBody(node.Body)
	TabsCount++
	gen.WriteTabbed("return Self\n")
	TabsCount--
	gen.WriteTabbed("end")
}

func (gen *Generator) fieldDeclaration(node ast.VariableDecl, index int) {
	src := StringBuilder{}

	src.WriteTabbed()
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
		gen.WriteString(", ")
		gen.Write(gen.WriteVar(param.Name.Lexeme))
	}
	gen.Write(")\n")

	gen.GenerateBody(node.Body)

	gen.WriteTabbed("end")
}

func (gen *Generator) entityFunctionDeclaration(node ast.MethodDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	gen.Write("function ", entityName, "_", node.Name.Lexeme, "(id")
	for _, param := range node.Params {
		gen.WriteString(", ")
		gen.Write(gen.WriteVar(param.Name.Lexeme))
	}
	gen.Write(")\n")
	TabsCount++
	gen.WriteTabbed("local Self = ", entityName, "[id]\n")
	TabsCount--

	gen.GenerateBody(node.Body)

	gen.WriteTabbed("end")
}

func (gen *Generator) spawnDeclaration(node ast.EntityFunctionDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)

	gen.Write("function ", entityName, "_Spawn(")

	gen.GenerateParams(node.Params)

	TabsCount++

	gen.WriteTabbed("local id = pewpew.new_customizable_entity(", gen.WriteVar(node.Params[0].Name.Lexeme), ", ", gen.WriteVar(node.Params[1].Name.Lexeme), ")\n")
	tableAccess := entityName + "[id]"
	gen.WriteTabbed(tableAccess, " = {}\n")
	gen.WriteTabbed("local Self = ", tableAccess)
	counter := 1
	for _, field := range entity.Fields {
		gen.Write("\n")
		gen.fieldDeclaration(field, counter)
		counter += len(field.Identifiers)
	}

	TabsCount--
	gen.GenerateBody(node.Body)
	TabsCount++

	for i, v := range entity.Callbacks {
		switch v.Type {
		case ast.WallCollision:
			gen.WriteTabbed(fmt.Sprintf("\npewpew.customizable_entity_configure_wall_collision(id, true, %sHCb%d)\n", entityName, i))
		case ast.WeaponCollision:
			gen.WriteTabbed(fmt.Sprintf("\npewpew.customizable_entity_set_weapon_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.PlayerCollision:
			gen.WriteTabbed(fmt.Sprintf("\npewpew.customizable_entity_set_player_collision_callback(id, %sHCb%d)\n", entityName, i))
		case ast.Update:
			gen.WriteTabbed(fmt.Sprintf("\npewpew.entity_set_update_callback(id, %sHCb%d)\n", entityName, i))
		}
	}
	gen.WriteTabbed("return id\n")
	TabsCount--

	gen.WriteTabbed("end")
}

func (gen *Generator) destroyDeclaration(node ast.EntityFunctionDecl, entity ast.EntityDecl) {
	entityName := gen.WriteVarExtra(entity.Name.Lexeme, hyEntity)
	gen.Write("function ", entityName, "_Destroy(id")
	if len(node.Params) != 0 {
		gen.Write(", ")
	}
	gen.GenerateParams(node.Params)
	TabsCount++
	gen.WriteTabbed("local Self = ", entityName, "[id]\n")
	TabsCount--

	gen.GenerateBody(node.Body)

	gen.WriteString("end")
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
