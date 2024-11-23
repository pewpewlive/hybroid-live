package ast

import "hybroid/tokens"

type Error struct {
	Token   tokens.Token
	Message string
}
