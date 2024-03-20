package generators

import (
	"hybroid/generators/lua"
	"hybroid/parser"
)

type Generator interface {
	Generate(program parser.Program, environment *lua.Scope) lua.Value
	GetErrors() []lua.GenError
	GetSrc() []byte
}
