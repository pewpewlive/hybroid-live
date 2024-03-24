package err

import "hybroid/lexer"

type Error struct {
	Token   lexer.Token
	Message string
}

func (e *Error) New(token lexer.Token, msg string) Error {
	return Error{token, msg}
}