package ast

import "hybroid/tokens"

type Alert interface {
	GetToken() tokens.Token
	GetHeader() string
	GetMessage() string
}

type Error struct {
	Token   tokens.Token
	Message string
}

func (e Error) GetToken() tokens.Token {
	return e.Token
}

func (e Error) GetMessage() string {
	return e.Message
}

func (e Error) GetHeader() string {
	return "[red]Error"
}

type Warning struct {
	Token   tokens.Token
	Message string
}

func (w Warning) GetToken() tokens.Token {
	return w.Token
}

func (w Warning) GetMessage() string {
	return w.Message
}

func (e Warning) GetHeader() string {
	return "[yellow]Warning"
}
