package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
)

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

	gen.TabsCount += 1
	tabs := gen.getTabs()

	src.WriteString("{\n")
	index := 0
	for k, v := range node.Map {
		val := gen.GenerateNode(v.Expr)

		ident := k.Lexeme

		if k.Type == lexer.String {
			ident = k.Literal
		}

		if index != len(node.Map)-1 {
			src.WriteString(fmt.Sprintf("%s%s = %v,\n", tabs, ident, val))
		} else {
			src.WriteString(fmt.Sprintf("%s%s = %v\n", tabs, ident, val))
		}
		index++
	}
	gen.TabsCount -= 1
	tabs = gen.getTabs()
	src.Append(tabs, "}")

	return src.String()
}

func (gen *Generator) unaryExpr(node ast.UnaryExpr) string {
	return fmt.Sprintf("%s%s", node.Operator.Lexeme, gen.GenerateNode(node.Value))
}

func (gen *Generator) parentExpr(node ast.ParentExpr) string {
	return gen.GenerateNode(node.Member)
}

func (gen *Generator) memberExpr(node ast.MemberExpr) string {
	src := StringBuilder{}

	if node.Property.GetType() == ast.MemberExpression {
		return gen.memberExpr(node.Property.(ast.MemberExpr))
	}
	expr := gen.GenerateNode(node.Owner)
	prop := gen.GenerateNode(node.Property)

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
