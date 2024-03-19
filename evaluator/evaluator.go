package evaluator

import (
	"fmt"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
)

func New(fc FileCache) LuaGenerator {
	return LuaGenerator{
		fc,
		lexer.New([]byte(fc.DstPath)),
		parser.Parser{},
		nil,
	}
}

func (e *LuaGenerator) HasValidSrc() bool {
	return e.file.SrcIsValid()
}

func (e *LuaGenerator) Action() {
	lcsrc, _ := os.ReadFile(e.file.SrcPath)
	e.lexer.ChangeSrc(lcsrc)
	e.lexer.Tokenize()
	if len(e.lexer.Errors) != 0 {
		for _, err := range e.lexer.Errors {
			fmt.Errorf("Test failed tokenizing: %v\n", err)
		}
		return
	}
	e.parser.ParseTokens(e.lexer.Tokens)
}
