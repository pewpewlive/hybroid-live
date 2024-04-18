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
		w.error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)",left.GetType().ToString(), right.GetType().ToString()))
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
	case ast.Fixed:
		return FixedVal{
			ast.Fixed,
			node.Value,
		}
	case ast.Radian:
		return FixedVal{
			ast.Radian,
			node.Value,
		}
	case ast.FixedPoint: 
		return FixedVal{
			ast.FixedPoint,
			node.Value,
		}
	case ast.Degree: 
		return FixedVal{
			ast.Degree,
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
		}
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
		value.Values = append(value.Values, w.GetNodeValue(&expr, scope))
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

func (w *Walker) memberExpr(array Value, node *ast.MemberExpr, scope *Scope) Value {
	if node.Owner == nil {
		sc := scope.Resolve(node.Identifier.GetToken().Lexeme)

		var array Value
		if sc == nil {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("undeclared variable \"%s\"", node.Identifier.GetToken().Lexeme))
			return Unknown{}
		} else {
			array = sc.GetVariable(node.Identifier.GetToken().Lexeme).Value
		}

		next, ok := node.Property.(ast.MemberExpr)

		if ok {
			if array.GetType() == ast.Namespace {
				return Undefined{}
			}
			if array.GetType() != ast.List && array.GetType() != ast.Map {
				w.error(node.Identifier.GetToken(), "variable is not a list, map or a namespace")
				return Unknown{}
			}
			return w.memberExpr(array, &next, scope)
		} else {
			w.error(node.GetToken(), "expected member expression")
			return Unknown{}
		}
	}

	val := w.GetNodeValue(&node.Identifier, scope)

	if node.Bracketed {
		if array.GetType() == ast.Map {
			if val.GetType() != ast.String && val.GetType() != 0 {
				w.error(node.Identifier.GetToken(), "variable is not a string")
				return Unknown{}
			}
		} else if array.GetType() == ast.List {
			if val.GetType() != ast.Number && val.GetType() != 0 {
				w.error(node.Identifier.GetToken(), "variable is not a number")
				return Unknown{}
			}
		}
	} 

	wrappedValType := ast.Undefined
	if list, ok := array.(ListVal); ok {
		wrappedValType = list.ValueType
	}else if mapp, ok := array.(ListVal); ok {
		wrappedValType = mapp.ValueType
	}
	
	return w.GetValue(wrappedValType)
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
