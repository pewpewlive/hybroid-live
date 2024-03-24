package generators

import (
	"hybroid/ast"
)

type Generator interface {
	Generate(program []ast.Node) string
	GetErrors() []ast.Error
	GetSrc() []byte
}
