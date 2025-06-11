package generator

import (
	"fmt"
	"hybroid/ast"
	"hybroid/core"
	"math"
	"strconv"
)

func (gen *Generator) GenerateParams(params []ast.FunctionParam) {
	var variadicParam string
	for i, param := range params {
		if param.Type.IsVariadic {
			gen.Write("...")
			variadicParam = gen.GenerateExpr(&ast.IdentifierExpr{Name: param.Name})
		} else {
			gen.Write(gen.GenerateExpr(&ast.IdentifierExpr{Name: param.Name}))
		}
		if i != len(params)-1 {
			gen.Write(", ")
		}
	}
	gen.Write(")\n")

	if variadicParam != "" {
		gen.Twrite("local ", variadicParam, " = {...}\n")
	}
}

func (gen *Generator) GenerateArgs(args []ast.Node) string {
	src := core.StringBuilder{}

	for i, arg := range args {
		src.Write(gen.GenerateExpr(arg))
		if i != len(args)-1 {
			src.Write(", ")
		}
	}
	src.Write(")")

	return src.String()
}

func (gen *Generator) breakDownAssignStmt(stmt ast.AssignmentStmt) []ast.AssignmentStmt {
	emptyVarDecl := func() ast.AssignmentStmt {
		return ast.AssignmentStmt{
			Identifiers: []ast.Node{},
			Values:      []ast.Node{},
			AssignOp:    stmt.AssignOp,
		}
	}
	stmts := []ast.AssignmentStmt{}
	currentDeclIndex := -1
	for _, expr := range stmt.Values {
		if call, ok := expr.(ast.CallNode); ok && call.GetReturnAmount() > 1 {
			stmts = append(stmts, emptyVarDecl())
			currentDeclIndex = len(stmts) - 1
			for range call.GetReturnAmount() {
				stmts[currentDeclIndex].Identifiers = append(stmts[currentDeclIndex].Identifiers, stmt.Identifiers[0])
				stmt.Identifiers = stmt.Identifiers[1:]
			}
			stmts[currentDeclIndex].Values = append(stmts[currentDeclIndex].Values, expr)
			currentDeclIndex = -1
			continue
		}
		if currentDeclIndex == -1 {
			stmts = append(stmts, emptyVarDecl())
			currentDeclIndex = len(stmts) - 1
		}
		stmts[currentDeclIndex].Values = append(stmts[currentDeclIndex].Values, expr)
		stmts[currentDeclIndex].Identifiers = append(stmts[currentDeclIndex].Identifiers, stmt.Identifiers[0])
		stmt.Identifiers = stmt.Identifiers[1:]
	}

	return stmts
}

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
		if call, ok := expr.(ast.CallNode); ok && call.GetReturnAmount() > 1 {
			decls = append(decls, emptyVarDecl())
			currentDeclIndex = len(decls) - 1
			for range call.GetReturnAmount() {
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

func (gen *Generator) GenerateBodyValue(body ast.Body) string {
	gen.writeToBuffer = true
	gen.tabCount++
	for _, node := range body {
		gen.GenerateStmt(node)
	}
	gen.tabCount--
	gen.writeToBuffer = false
	defer gen.buffer.Reset()
	return gen.buffer.String()
}

func fixedToFx(floatstr string) string {
	float, _ := strconv.ParseFloat(floatstr, 64)
	abs_float := math.Abs(float)
	integer := min(math.Floor(abs_float), (2 << 51))
	var sign string
	if float < 0 {
		sign = "-"
	} else {
		sign = ""
	}

	frac := math.Floor((abs_float - integer) * 4096)
	frac_str := ""
	if frac != 0 {
		frac_str = "." + fmt.Sprintf("%v", frac)
	}

	// sign + int + frac_str + "fx"
	return fmt.Sprintf("%s%v%s", sign, integer, frac_str)
}

func degToRad(floatstr string) string {
	float, _ := strconv.ParseFloat(floatstr, 64)
	radians := float * math.Pi / 180
	return fixedToFx(fmt.Sprintf("%v", radians))
}
