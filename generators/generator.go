package generators

import (
	"hybroid/ast"
	"hybroid/generators/lua"
	"hybroid/err"
)

type Generator interface {
	Generate(program []ast.Node, environment *lua.Scope) lua.Value
	GetErrors() []err.Error
	GetSrc() []byte
}
