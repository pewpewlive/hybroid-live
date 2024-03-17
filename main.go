package main

import (
	"fmt"
	"livecode/lexer"
	"os"
)

func main() {
	file, err := os.ReadFile("./example.lc")
	if err != nil {
		fmt.Printf("[Error -> Main] reading file: %s\n", err.Error())
		return
	}

	tokens := lexer.Tokenize(file)
	for _, token := range tokens {
		fmt.Printf("Token { type: %v, lex: %v, lit: %v, line: %v }\n", token.Type.ToString(), token.Lexeme, token.Literal, token.Line)
	}
}
