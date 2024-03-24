package ast

import "hybroid/lexer"

type Error struct {
	Token   lexer.Token
	Message string
}
