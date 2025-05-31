package generator

import (
	"fmt"
	"hybroid/ast"
	"strconv"
)

func (gen *Generator) variableDeclaration(declaration ast.VariableDecl) {
	var values []string

	for _, expr := range declaration.Expressions {
		values = append(values, gen.GenerateExpr(expr))
	}

	src := StringBuilder{}
	src2 := StringBuilder{}
	if !declaration.IsPub {
		src.WriteTabbed("local ")
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

	//gen.spawnDeclarationStmt(*node.Spawner, node)
	//gen.destroyDeclarationStmt(*node.Destroyer, node)

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
