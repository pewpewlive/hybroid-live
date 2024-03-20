package generators

import (
	"hybroid/parser"
	"hybroid/generators/lua"
)

type Generator interface {
	Generate(program parser.Program, environment *lua.Scope) lua.Value
	GetErrors() []lua.GenError
}