package ast

import "hybroid/lexer"

type Alert interface {
	GetToken() lexer.Token
	GetMessage() string
}

type Error struct {
	Token   lexer.Token
	Message string
}

func (e Error) GetToken() lexer.Token {
	return e.Token
}

func (e Error) GetMessage() string {
	return e.Message
}

type Warning struct {
	Token   lexer.Token
	Message string
}

func (w Warning) GetToken() lexer.Token {
	return w.Token
}

func (w Warning) GetMessage() string {
	return w.Message
}
