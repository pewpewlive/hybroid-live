package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
)

// let a = {a:{b:0},b:[2,3,4]}

// Member Node
// a.a.b

func (gen *Generator) binaryExpr(node ast.BinaryExpr) string {
	left, right := gen.GenerateNode(node.Left), gen.GenerateNode(node.Right)
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

func (gen *Generator) identifierExpr(node ast.IdentifierExpr) string {
	return node.Name.Lexeme
}

func (gen *Generator) groupingExpr(node ast.GroupExpr) string {
	return fmt.Sprintf("(%s)", gen.GenerateNode(node.Expr))
}

func (gen *Generator) listExpr(node ast.ListExpr) string {
	src := StringBuilder{}

	src.WriteString("{")
	for i, expr := range node.List {
		src.WriteString(gen.GenerateNode(expr))

		if i != len(node.List)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString("}")

	return src.String()
}

func (gen *Generator) callExpr(node ast.CallExpr) string {
	src := StringBuilder{}
	fn := gen.GenerateNode(node.Caller)

	src.Append(fn, "(")
	for i, arg := range node.Args {
		src.WriteString(gen.GenerateNode(arg))
		if i != len(node.Args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")

	return src.String()
}

func (gen *Generator) mapExpr(node ast.MapExpr) string {
	src := StringBuilder{}

	src.WriteString("{\n")
	TabsCount += 1
	index := 0
	for k, v := range node.Map {
		val := gen.GenerateNode(v.Expr)

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

func (gen *Generator) unaryExpr(node ast.UnaryExpr) string {
	return fmt.Sprintf("%s%s", node.Operator.Lexeme, gen.GenerateNode(node.Value))
}

func (gen *Generator) memberExpr(node ast.MemberExpr) string {
	src := StringBuilder{}

	if node.Property.GetType() == ast.MemberExpression {
		return gen.memberExpr(node.Property.(ast.MemberExpr))
	}

	expr := gen.GenerateNode(node.Owner)
	prop := gen.GenerateNode(node.Property)

	if expr == "" {
		return prop
	}
	src.WriteString(expr)

	if node.Bracketed {
		src.Append("[", prop, "]")
	} else {
		src.Append(".", prop)
	}

	return src.String()
}

func (gen *Generator) directiveExpr(node ast.DirectiveExpr) string {
	src := StringBuilder{}

	if node.Identifier.Lexeme != "Environment" {
		src.Append(node.Identifier.Lexeme, "(", gen.GenerateNode(node.Expr), ")")
	}

	return src.String()
}

func (gen *Generator) anonFnExpr(fn ast.AnonFnExpr) string {
	src := StringBuilder{}

	TabsCount += 1

	src.WriteString("function (")
	for i, param := range fn.Params {
		src.Append(param.Name.Lexeme)
		if i != len(fn.Params)-1 {
			src.Append(", ")
		}
	}
	src.Append(")\n")

	src.WriteString(gen.GenerateString(fn.Body))

	TabsCount -= 1

	src.AppendTabbed("end")

	return src.String()
}

func (gen *Generator) selfExpr(self ast.SelfExpr) string {
	if self.Type == ast.SelfStruct {
		if self.Value.GetType() == ast.SelfExpression {
			return "Self"
		}else {
			return "Self."+gen.GenerateNode(self.Value)
		}
	}
	return ""
}