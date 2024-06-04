package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (gen *Generator) ifStmt(node ast.IfStmt, scope *GenScope) {
	ifScope := NewGenScope(scope)
	ifTabs := getTabs()

	TabsCount += 1

	ifScope.Append(ifTabs, "if ", gen.GenerateExpr(node.BoolExpr, scope), " then\n")

	gen.GenerateString(node.Body, &ifScope)
	for _, elseif := range node.Elseifs {
		ifScope.Append(ifTabs, "elseif ", gen.GenerateExpr(elseif.BoolExpr, scope), " then\n")
		gen.GenerateString(elseif.Body, &ifScope)
	}
	if node.Else != nil {
		ifScope.Append(ifTabs, "else \n")
		gen.GenerateString(node.Else.Body, &ifScope)
	}

	ifScope.Append(ifTabs, "end\n")

	ifScope.ReplaceAll(ifScope.ReplaceSettings)

	TabsCount -= 1
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
	fnTabs := getTabs()

	TabsCount += 1

	if node.IsLocal {
		fnScope.Append(fnTabs, "local ")
	} else {
		fnScope.Append(fnTabs)
	}

	fnScope.Append("function ", "V", node.Name.Lexeme, "(")
	for i, param := range node.Params {
		fnScope.Append("V", param.Name.Lexeme)
		if i != len(node.Params)-1 {
			fnScope.Append(", ")
		}
	}
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

	gen.GenerateString(node.Body, &fnScope)

	fnScope.Append(fnTabs + "end")

	TabsCount -= 1

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
	repeatTabs := getTabs()

	TabsCount += 1
	//tabs := getTabs()

	end := gen.GenerateExpr(node.Iterator, scope)
	start := gen.GenerateExpr(node.Start, scope)
	skip := gen.GenerateExpr(node.Skip, scope)
	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateExpr(node.Variable, &repeatScope)
		repeatScope.Append(repeatTabs, "for ", variable, " = ", start, ", ", end, ", ", skip, " do\n")
	} else {
		repeatScope.Append(repeatTabs, "for _ = ", start, ", ", end, ", ", skip, " do\n")
	}

	gotoLabel := "hgtl" + GenerateVar()
	repeatScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateString(node.Body, &repeatScope)

	repeatScope.ReplaceAll(repeatScope.ReplaceSettings)

	scope.Write(repeatScope.Src)

	scope.AppendTabbed("::" + gotoLabel + "::\n")

	TabsCount -= 1

	scope.Append(repeatTabs, "end")
}

func (gen *Generator) forStmt(node ast.ForStmt, scope *GenScope) {
	forScope := NewGenScope(scope)
	forTabs := getTabs()

	TabsCount += 1

	iterator := gen.GenerateExpr(node.Iterator, scope)
	if node.Key.GetValueType() != 0 && node.Value.GetValueType() != 0 {
		key := gen.GenerateExpr(node.Key, &forScope)
		value := gen.GenerateExpr(node.Value, &forScope)
		forScope.Append(forTabs, "for ", key, ", ", value, " in ipairs(", iterator, ") do\n")
	} else {
		value := gen.GenerateExpr(node.Value, scope)
		forScope.Append(forTabs, "for _, ", value, " in ipairs(", iterator, ") do\n")
	}

	gotoLabel := "hgtl" + GenerateVar()
	forScope.ReplaceSettings = map[ReplaceType]string{
		ContinueReplacement: "goto " + gotoLabel,
	}

	gen.GenerateString(node.Body, &forScope)

	forScope.ReplaceAll(forScope.ReplaceSettings)

	scope.Write(forScope.Src)

	scope.AppendTabbed("::" + gotoLabel + "::\n")

	TabsCount -= 1

	scope.Append(forTabs, "end")
}

func (gen *Generator) tickStmt(node ast.TickStmt, scope *GenScope) {
	tickTabs := getTabs()

	tickScope := GenScope{Src: StringBuilder{}, Parent: scope}

	TabsCount += 1

	if node.Variable.GetValueType() != 0 {
		variable := gen.GenerateExpr(node.Variable, scope)
		tickScope.Src.Append(tickTabs, "local ", variable, " = 0\n")
		tickScope.Src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
		tickScope.Src.AppendTabbed(variable, " = ", variable, " + 1\n")
	} else {
		tickScope.Src.Append(tickTabs, "pewpew.add_update_callback(function()\n")
	}

	gen.GenerateString(node.Body, &tickScope)

	tickScope.Src.Append(tickTabs, "end)")

	TabsCount -= 1

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
			src.Append("V", fmt.Sprintf("%s = ", ident.Lexeme))
		} else if i == len(declaration.Identifiers)-1 {
			src.Append("V", ident.Lexeme)
		} else {
			src.Append("V", fmt.Sprintf("%s, ", ident.Lexeme))
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

func (gen *Generator) structDeclarationStmt(node ast.StructDeclarationStmt, scope *GenScope) {
	structScope := NewGenScope(scope)

	for _, nodebody := range *node.Methods {
		gen.methodDeclarationStmt(nodebody, node, &structScope)
	}

	gen.constructorDeclarationStmt(*node.Constructor, node, &structScope)

	scope.Write(structScope.Src)
}

func (gen *Generator) constructorDeclarationStmt(node ast.ConstructorStmt, Struct ast.StructDeclarationStmt, scope *GenScope) {
	src := StringBuilder{}

	constructorScope := NewGenScope(scope)

	if Struct.IsLocal {
		constructorScope.WriteString("local ")
	}

	constructorScope.Append("function Hybroid_", Struct.Name.Lexeme, "_New(")

	TabsCount += 1

	for i, param := range node.Params {
		constructorScope.Append("V", param.Name.Lexeme)
		if i != len(node.Params)-1 {
			constructorScope.Append(", ")
		}
	}
	constructorScope.Append(")\n")

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

	gen.GenerateString(*node.Body, &constructorScope)
	constructorScope.AppendTabbed("return Self\n")
	TabsCount -= 1
	constructorScope.AppendTabbed("end\n")

	scope.Write(constructorScope.Src)
}

func (gen *Generator) methodDeclarationStmt(node ast.MethodDeclarationStmt, Struct ast.StructDeclarationStmt, scope *GenScope) {
	methodScope := NewGenScope(scope)

	if Struct.IsLocal {
		methodScope.WriteString("local ")
	}

	methodScope.Append("function Hybroid_", Struct.Name.Lexeme, "_", node.Name.Lexeme, "(Self")
	for _, param := range node.Params {
		methodScope.WriteString(", ")
		methodScope.Append("V", param.Name.Lexeme)
	}
	methodScope.Append(")\n")

	TabsCount += 1

	gen.GenerateString(node.Body, &methodScope) // its constructor

	TabsCount -= 1

	methodScope.AppendTabbed("end\n")

	scope.Write(methodScope.Src)
}

func (gen *Generator) useStmt(node ast.UseStmt, scope *GenScope) {
	fileName := strings.Replace(node.File.Literal, ".hyb", ".lua", 1)
	scope.Append("local ", node.Variable.Name.Lexeme, " = require(\"/dynamic/", fileName, "\")")
}

func (gen *Generator) matchStmt(node ast.MatchStmt, scope *GenScope) {
	ifStmt := ast.IfStmt{
		BoolExpr: ast.BinaryExpr{Left: node.ExprToMatch, Operator: lexer.Token{Type: lexer.EqualEqual, Lexeme: "=="}, Right: node.Cases[0].Expression},
		Body:     node.Cases[0].Body,
	}
	has_default := false
	for i := range node.Cases {
		if node.Cases[i].Expression.GetToken().Lexeme == "_" {
			has_default = true
		}
		if i == 0 || i == len(node.Cases)-1 {
			continue
		}
		elseIfStmt := ast.IfStmt{
			BoolExpr: ast.BinaryExpr{Left: node.ExprToMatch, Operator: lexer.Token{Type: lexer.EqualEqual, Lexeme: "=="}, Right: node.Cases[i].Expression},
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
