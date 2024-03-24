package generators

import (
	"hybroid/ast"
	"hybroid/generators/lua"
)

type Generator interface {
	Generate(program []ast.Node, environment *lua.Scope) lua.Value
	GetErrors() []ast.Error
	GetSrc() []byte
}
