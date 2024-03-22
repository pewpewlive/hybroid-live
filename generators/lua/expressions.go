package lua

import (
	"fmt"
	"hybroid/ast"
	"strings"
)

func (gen *Generator) binaryExpr(node ast.BinaryExpr, scope *Scope) Value {
	src := strings.Builder{}
	left, right := gen.GenerateNode(node.Left, scope), gen.GenerateNode(node.Right, scope)
	gen.validateOperands(left, right)
	src.WriteString(left.Val)
	src.WriteString(fmt.Sprintf(" %s ", node.Operator.Lexeme))
	src.WriteString(right.Val)

	return Value{Type: node.ValueType, Token: node.Operator, Val: src.String()}
}

func (gen *Generator) literalExpr(node ast.LiteralExpr) Value {
	src := strings.Builder{}

	switch node.ValueType {
	case ast.String:
		src.WriteString("\"")
		src.WriteString(fmt.Sprintf("%v", node.Value))
		src.WriteString("\"")
	case ast.Fixed:
		src.WriteString(fixedToFx(node.Value))
		src.WriteString("fx")
	case ast.FixedPoint:
		src.WriteString(fmt.Sprintf("%vfx", node.Value))
	default:
		src.WriteString(fmt.Sprintf("%v", node.Value))
	}

	return Value{Type: node.ValueType, Token: node.Token, Val: src.String()}
}

func (gen *Generator) identifierExpr(node ast.IdentifierExpr, scope *Scope) Value {
	sc := scope.Resolve(node.Name)
	value := Value{Type: node.ValueType, Token: node.Token, Val: node.Name}
	if sc != nil {
		newValue := sc.GetVariable(node.Name)
		value.Type = newValue.Type
	} else {

		//error
	}
	return value
}

func (gen *Generator) groupingExpr(node ast.GroupExpr, scope *Scope) Value {
	src := strings.Builder{}
	src.WriteString("(")
	value := gen.GenerateNode(node.Expr, scope)
	src.WriteString(value.Val)
	src.WriteString(")")

	return Value{Type: value.Type, Token: node.Token, Val: src.String()}
}

func (gen *Generator) listExpr(node ast.ListExpr, scope *Scope) Value {

	src := strings.Builder{}
	src.WriteString("{")
	for i, expr := range node.List {
		src.WriteString(gen.GenerateNode(expr, scope).Val)

		if i != len(node.List)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString("}")

	return Value{Type: ast.List, Token: node.Token, Val: src.String()}
}

func (gen *Generator) callExpr(node ast.CallExpr, scope *Scope) Value {
	src := strings.Builder{}
	fn := gen.GenerateNode(node.Caller, scope)

	scope.Resolve(fn.Val)

	src.WriteString(fn.Val)
	src.WriteString("(")
	for i, arg := range node.Args {
		src.WriteString(gen.GenerateNode(arg, scope).Val)
		if i != len(node.Args)-1 {
			src.WriteString(", ")
		}
	}
	src.WriteString(")")

	//fnReturn := scope.GetVariable(node.Identifier)

	return Value{Type: ast.Bool, Token: node.Token, Val: src.String()}
}

func (gen *Generator) mapExpr(node ast.MapExpr, scope *Scope) Value {
	src := strings.Builder{}

	var tabs string
	for i := 0; i < scope.Count+gen.TabsCount; i++ {
		tabs += "\t"
	}

	var mapTabs string
	for i := 0; i < scope.Count+1+gen.TabsCount; i++ {
		mapTabs += "\t"
	}

	gen.TabsCount += 1

	src.WriteString("{\n")
	index := 0
	for k, v := range node.Map {
		val := gen.GenerateNode(v.Expr, scope)

		if index != len(node.Map)-1 {
			src.WriteString(fmt.Sprintf("%s%s = %v,\n", mapTabs, k, val.Val))
		} else {
			src.WriteString(fmt.Sprintf("%s%s = %v\n", mapTabs, k, val.Val))
		}
		index++
	}

	src.WriteString(tabs)
	src.WriteString("}")

	gen.TabsCount -= 1

	return Value{Type: ast.Map, Token: node.Token, Val: src.String()}
}

func (gen *Generator) unaryExpr(node ast.UnaryExpr, scope *Scope) Value {
	value := gen.GenerateNode(node.Value, scope)
	src := fmt.Sprintf("%s%s", node.Operator.Lexeme, value.Val)

	return Value{Type: value.Type, Token: node.Operator, Val: src}
}

func (gen *Generator) memberExpr(node ast.MemberExpr, scope *Scope) Value {
	src := strings.Builder{}

	expr := gen.GenerateNode(node.Identifier, scope)
	prop := gen.GenerateNode(node.Property, scope)

	src.WriteString(expr.Val)

	if prop.Type == ast.String {
		src.WriteString("[")
		src.WriteString(prop.Val)
		src.WriteString("]")
	} else {
		src.WriteString(".")
		src.WriteString(prop.Val)
	}

	return Value{Type: prop.Type, Token: node.Token, Val: src.String()}
}
