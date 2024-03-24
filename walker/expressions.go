package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) determineValueType(left Value, right Value) ast.PrimitiveValueType {
	if left.GetType() == 0 || right.GetType() == 0 {
		return 0
	}
	if left.GetType() == right.GetType() {
		return left.GetType()
	}
	if parser.IsFx(left.GetType()) && parser.IsFx(right.GetType()) {
		return ast.FixedPoint
	}

	return ast.Undefined
}

func (w *Walker) binaryExpr(node ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(node.Left, scope), w.GetNodeValue(node.Right, scope)
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.validateArithmeticOperands(left, right, node)
	}
	val := w.GetValue(w.determineValueType(left, right))

	if val.GetType() == ast.Undefined {
		w.error(node.GetToken(), "invalid binary expression")
		return val
	} else {
		return val
	}
}

func (w *Walker) literalExpr(node ast.LiteralExpr) Value {

	switch node.ValueType {
	case ast.String:
		return StringVal{}
	case ast.FixedPoint, ast.Radian, ast.Fixed, ast.Degree:
		return FixedVal{} //but that means that we have to modify the nodes before adding them to the new nodes list
	case ast.Bool:
		return BoolVal{}
	case ast.Nil: // map expr was messing up the if stmt
		return NilVal{} // list and map also
	case ast.Number:
		return NumberVal{} // ok
	default:
		return Unknown{}
	}

}

func (w *Walker) identifierExpr(node ast.IdentifierExpr, scope *Scope) Value {
	sc := scope.Resolve(node.Name.Lexeme)

	if sc != nil {
		newValue := sc.GetVariable(node.Name.Lexeme)
		return newValue
	} else {
		w.error(node.Name, "undeclared identifier in the current scope")
		return Unknown{}
	}
}

func (w *Walker) groupingExpr(node ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(node.Expr, scope)
}

func (w *Walker) listExpr(node ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for _, expr := range node.List {
		value.values = append(value.values, w.GetNodeValue(expr, scope))
	}
	return value
}

func (w *Walker) callExpr(node ast.CallExpr, scope *Scope) Value {
	callerToken := node.Caller.GetToken()
	sc := scope.Resolve(callerToken.Lexeme)

	if sc == nil { //make sure in the future member calls are also taken into account
		w.error(node.Token, "undeclared function")
		return Unknown{}
	} else {
		fn := sc.GetVariable(callerToken.Lexeme)
		call, ok := fn.Value.(CallVal)
		if !ok {
			w.error(callerToken, "variable used as if it's a function")
		} else if len(call.params) < len(node.Args) {
			w.error(callerToken, "too many arguments given in function call")
		} else if len(call.params) > len(node.Args) {
			w.error(callerToken, "too few arguments given in function call")
		}

		return call
	}
}

func (w *Walker) mapExpr(node ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: map[string]VariableVal{}}
	for k, v := range node.Map {
		val := w.GetNodeValue(v.Expr, scope)

		mapVal.Members[k.Lexeme] = VariableVal{
			Name: k.Lexeme,
			Value: val,
			Node: v.Expr,
		}
	}
	return mapVal
}

func (w *Walker) unaryExpr(node ast.UnaryExpr, scope *Scope) Value {
	return w.GetNodeValue(node.Value, scope)
}

func (w *Walker) memberExpr(node ast.MemberExpr, scope *Scope) Value {
	identToken := node.Identifier.GetToken()
	sc := scope.Resolve(identToken.Lexeme)

	if sc == nil {
		w.error(identToken, "undeclared map")
		return Unknown{}
	}

	variable := sc.GetVariable(identToken.Lexeme)
	mapp, ok := variable.Value.(MapVal)
	_, ok2 := variable.Value.(ListVal)

	if !ok && !ok2 {
		w.error(identToken, fmt.Sprintf("variable %s is not a map nor a list", identToken.Lexeme))
		return Unknown{}
	}
	propToken := node.Property.GetToken()

	val := w.GetNodeValue(node.Property, scope)
	if ok {
		if node.Bracketed {
			switch val.GetType() {
			case ast.Bool, ast.Entity, ast.Nil, ast.Struct, ast.List, ast.Map, ast.FixedPoint, ast.Number:
				w.error(propToken, "invalid expression inside brackets")
				return Unknown{}
			}
		}
		member, _ := findMember(mapp, propToken.Lexeme)
	
		return member
	}else if ok2 {
		if !node.Bracketed {
			w.error(propToken, "invalid member expression of a list")
			return Unknown{}
		}
		switch val.GetType() {
		case ast.Bool, ast.Entity, ast.Nil, ast.Struct, ast.List, ast.Map, ast.FixedPoint:
			w.error(propToken, "invalid expression inside brackets")
			return Unknown{}
		}
		return Undefined{}
	}

	return Unknown{}
}

func findMember(mapp MapVal, name string) (VariableVal, bool) {
	mem, found := mapp.Members[name]
	if found {
		return mem, true
	}
	return VariableVal{}, false
}

func (w *Walker) directiveExpr(node ast.DirectiveExpr, scope *Scope) Value {

	if node.Identifier.Lexeme != "Environment" {
		variable := w.GetNodeValue(node.Expr, scope)
		variableToken := node.Expr.GetToken()
		switch node.Identifier.Lexeme {
		case "Len":
			node.ValueType = ast.Number
			if variable.GetType() != ast.Map && variable.GetType() != ast.List && variable.GetType() != ast.String {
				w.error(variableToken, "invalid expression in '@Len' directive")
			}
		case "MapToStr":
			node.ValueType = ast.String
			if variable.GetType() != ast.Map {
				w.error(variableToken, "expected a map in '@MapToStr' directive")
			}
		case "ListToStr":
			node.ValueType = ast.List
			if variable.GetType() != ast.List {
				w.error(variableToken, "expected a list in '@ListToStr' directive")
			}
		default:
			// TODO: Implement custom directives

			w.error(node.Token, "unknown directive")
		}

	} else {

		ident, ok := node.Expr.(ast.IdentifierExpr)
		if !ok {
			w.error(node.Expr.GetToken(), "expected an identifier in '@Environment' directive")
		} else {
			name := ident.Name.Lexeme
			if name != "Level" && name != "Mesh" && name != "Sound" && name != "Shared" && name != "LuaGeneric" {
				w.error(node.Expr.GetToken(), "invalid identifier in '@Environment' directive")
			}
		}
	}
	return DirectiveVal{}
}
