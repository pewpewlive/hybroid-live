package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
	"strconv"
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

func (w *Walker) binaryExpr(node *ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(&node.Left, scope), w.GetNodeValue(&node.Right, scope)
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.validateArithmeticOperands(left, right, *node)
	default:
		if left.GetType() != right.GetType() {
			w.error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)",left.GetType().ToString(), right.GetType().ToString()))
		}else {
			return BoolVal{}
		}
	}
	val := w.GetValue(w.determineValueType(left, right))

	if val.GetType() == ast.Undefined {
		w.error(node.GetToken(), "invalid binary expression")
		return val
	} else {
		return val
	}
}

func (w *Walker) literalExpr(node *ast.LiteralExpr) Value {

	switch node.ValueType {
	case ast.String:
		return StringVal{
			node.Value,
		}
	case ast.FixedPoint, ast.Radian, ast.Fixed, ast.Degree:
		return FixedVal{
			node.Value,
		}
	case ast.Bool:
		return BoolVal{
			node.Value,
		}
	case ast.Nil:
		return NilVal{}
	case ast.Number:
		return NumberVal{
			node.Value,
		} // ok
	default:
		return Unknown{}
	}

}

func (w *Walker) identifierExpr(node *ast.IdentifierExpr, scope *Scope) Value {
	sc := scope.Resolve(node.Name.Lexeme)

	if sc != nil {
		newValue := sc.GetVariable(node.Name.Lexeme)
		return newValue
	} else {
		//w.error(node.Name, "unknown identifier")
		return Unknown{}
	}
}

func (w *Walker) groupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) listExpr(node *ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for _, expr := range node.List {
		value.values = append(value.values, w.GetNodeValue(&expr, scope))
	}
	return value
}

func (w *Walker) callExpr(node *ast.CallExpr, scope *Scope) Value {
	callerToken := node.Caller.GetToken()
	sc := scope.Resolve(callerToken.Lexeme)

	if sc == nil { //make sure in the future member calls are also taken into account
		w.error(node.Token, "undeclared function")
		return Unknown{}
	} else {
		fn := sc.GetVariable(callerToken.Lexeme)
		fun, ok := fn.Value.(FunctionVal)
		arguments := make([]ast.PrimitiveValueType, 0)
		for _, arg := range node.Args {
			val := w.GetNodeValue(&arg, scope)
			if function, ok := val.(FunctionVal); ok {
				arguments = append(arguments, function.returnVal.values...)
			}else {
				arguments = append(arguments, val.GetType())
			}
		}
		if !ok {
			w.error(callerToken, "variable used as if it's a function")
		} else if len(fun.params) < len(arguments) {
			w.error(callerToken, "too many arguments given in function call")
		} else if len(fun.params) > len(arguments) {
			w.error(callerToken, "too few arguments given in function call")
		}

		return fun
	}
}

func (w *Walker) mapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: map[string]MapMemberVal{}}
	for k, v := range node.Map {
		val := w.GetNodeValue(&v.Expr, scope)

		mapVal.Members[k.Lexeme] = MapMemberVal{
			Var: VariableVal{
				Name:  k.Lexeme,
				Value: val,
				Node:  v.Expr,
			},
			Owner: mapVal,
		}
	}
	return mapVal
}

func (w *Walker) unaryExpr(node *ast.UnaryExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Value, scope)
}

func (w *Walker) memberExpr(mapp Value, node *ast.MemberExpr, scope *Scope) Value {
	if node.Owner == nil {
		sc := scope.Resolve(node.Identifier.GetToken().Lexeme)

		var mapp Value
		if sc == nil {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("undeclared variable \"%s\"", node.Identifier.GetToken().Lexeme))
			return Unknown{}
		} else {
			mapp = sc.GetVariable(node.Identifier.GetToken().Lexeme).Value
		}

		next, ok := node.Property.(ast.MemberExpr)

		if ok {
			if mapp.GetType() == ast.Namespace {
				return Undefined{}
			}
			if mapp.GetType() != ast.List && mapp.GetType() != ast.Map {
				w.error(node.Identifier.GetToken(), "variable is not a list, map or a namespace")
				return Unknown{}
			}
			return w.memberExpr(mapp, &next, scope)
		} else {
			w.error(node.GetToken(), "expected member expression")
			return Unknown{}
		}
	}

	var val Value

	if node.Bracketed && node.GetToken().Type != lexer.Identifier {
		val = StringVal{node.GetToken().Literal}
	} else if node.Bracketed && node.GetToken().Type == lexer.Identifier {
		if mapp.GetType() == ast.Map {
			val = w.GetNodeValue(&node.Identifier, scope)
			if val.GetType() != ast.String && val.GetType() != 0 {
				w.error(node.Identifier.GetToken(), "variable is not a string")
				return Unknown{}
			}
		} else if mapp.GetType() == ast.List {
			val = w.GetNodeValue(&node.Identifier, scope)
			if val.GetType() != ast.Number && val.GetType() != 0 {
				w.error(node.Identifier.GetToken(), "variable is not a number")
				return Unknown{}
			}
			varia, isVar := val.(VariableVal)
			if isVar {
				val = varia.Value
			}
		}
	} else {
		val = StringVal{node.GetToken().Lexeme}
	}

	var mem Value
	mem, _ = findMember(mapp, val)

	mapMember, isVariable := mem.(MapMemberVal)

	var value Value
	if isVariable {
		value = mapMember.Var.Value
	} else {
		value = mem
	}

	next, isMember := node.Property.(ast.MemberExpr)

	if isMember {
		if value.GetType() == ast.Map {
			return w.memberExpr(value, &next, scope)
		} else if value.GetType() == ast.List {
			return w.memberExpr(value, &next, scope)
		} else {
			w.error(node.Identifier.GetToken(), "variable is being treated as a map or list, but isn't one")
			return Unknown{}
		}
	} else {
		return mem
	}
}

func findMember(val Value, detection Value) (Value, bool) {
	list, isList := val.(ListVal)
	mapp, isMap := val.(MapVal)

	if isList {
		if detection.GetType() == ast.Number {
			parsedNum, ok := strconv.Atoi(detection.(NumberVal).Val)
			if ok != nil {
				return Unknown{}, false
			}
			if parsedNum > len(list.values) {
				return Undefined{}, false
			}
			mem := list.values[parsedNum-1]
			return mem, false
		} else if detection.GetType() == 0 {
			return Undefined{}, false
		}
	} else if isMap {

		mem, found := mapp.Members[detection.(StringVal).Val]
		if found {
			return mem, false
		}
		return MapMemberVal{Var: VariableVal{Value: Undefined{}}, Owner: mapp}, false
	}
	return Unknown{}, false
}

func (w *Walker) directiveExpr(node *ast.DirectiveExpr, scope *Scope) Value {

	if node.Identifier.Lexeme != "Environment" {
		variable := w.GetNodeValue(&node.Expr, scope)
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
