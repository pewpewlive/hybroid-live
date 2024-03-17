package evaluator

import (
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

type Evaluator struct {
	file   FileCache
	lexer  lexer.Lexer
	parser parser.Parser
}

func New(fc FileCache) Evaluator {
	return Evaluator{
		fc,
		lexer.New([]byte(fc.DstPath)),
		parser.New(),
	}
}

func (e *Evaluator) HasValidSrc() bool {
	return e.file.SrcIsValid()
}

func (e *Evaluator) Action() {
	lcsrc, _ := os.ReadFile(e.file.SrcPath)
	e.lexer.ChangeSrc(lcsrc)
	e.lexer.Tokenize()
	e.parser.ParseTokens(e.lexer.Tokens)

}
