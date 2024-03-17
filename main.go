package main

import (
	"fmt"
	"hybroid/lexer"
	"os"
)

func main() {
	file, err := os.ReadFile("./example.hyb")
	if err != nil {
		fmt.Printf("[Error -> Main] reading file: %s\n", err.Error())
		return
	}

	l := lexer.New(file)
	l.Tokenize()
	// for _, token := range tokens {
	// 	fmt.Printf("Token { type: %v, lex: %v, lit: %v, line: %v }\n", token.Type.ToString(), token.Lexeme, token.Literal, token.Line)
	// }

	// actual compilation
	//var evaluators []evaluator.Evaluator

}
