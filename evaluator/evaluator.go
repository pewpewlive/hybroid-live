package evaluator

import (
	"fmt"
	"hybroid/lexer"
	"hybroid/parser"
	"os"
	"strings"
)

type FileCache struct {
	SrcPath string
	DstPath string
}

func NewFileCache(lcSrcPath string) FileCache {
	luaSrcPath := lcSrcPath
	luaSrcPath = strings.Replace(luaSrcPath, "hybsrc", "luasrc", 1)
	return FileCache{
		lcSrcPath,
		luaSrcPath,
	}
}

func (fc *FileCache) SrcIsValid() bool {
	_, err := os.ReadFile(fc.SrcPath)

	return err == nil
}

type LuaGenerator struct {
	file   FileCache
	lexer  lexer.Lexer
	parser parser.Parser
}

func New(fc FileCache) LuaGenerator {
	return LuaGenerator{
		fc,
		lexer.New([]byte(fc.DstPath)),
		parser.New(),
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
