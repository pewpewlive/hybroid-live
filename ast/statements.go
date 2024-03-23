package ast

import "hybroid/lexer"

type AssignmentStmt struct {
	Values      []Node
	Identifiers []Node
	Token       lexer.Token
}

func (as AssignmentStmt) GetType() NodeType {
	return AssignmentStatement
}

func (as AssignmentStmt) GetToken() lexer.Token {
	return as.Token
}

func (as AssignmentStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type VariableDeclarationStmt struct {
	Identifiers []string
	Values      []Node
	IsLocal     bool
	Token       lexer.Token
}

func (vds VariableDeclarationStmt) GetType() NodeType {
	return VariableDeclarationStatement
}

func (vds VariableDeclarationStmt) GetToken() lexer.Token {
	return vds.Token
}

func (vds VariableDeclarationStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type FunctionDeclarationStmt struct {
	Name    lexer.Token
	Params  []lexer.Token
	IsLocal bool
	Body    []Node
}

func (fds FunctionDeclarationStmt) GetType() NodeType {
	return FunctionDeclarationStatement
}

func (fds FunctionDeclarationStmt) GetToken() lexer.Token {
	return fds.Name
}

func (fds FunctionDeclarationStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type IfStmt struct {
	BoolExpr Node
	Body     []Node
	Token    lexer.Token
}

func (is IfStmt) GetType() NodeType {
	return IfStatement
}

func (is IfStmt) GetToken() lexer.Token {
	return is.Token
}

func (is IfStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type MatchStmt struct {
	ExprToMatch Node
	Match       PrimitiveValueType
	Cases       [][]Node // EXPRESSIONS
	Bodies      []Node   // STATEMENTS
}

type RepeatStmt struct {
	Iterator Node
	Variable IdentifierExpr
	Body     []Node
	Token    lexer.Token
}

func (rs RepeatStmt) GetType() NodeType {
	return RepeatStatement
}

func (rs RepeatStmt) GetToken() lexer.Token {
	return rs.Token
}

func (rs RepeatStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type ReturnStmt struct {
	Args  []Node
	Token lexer.Token
}

func (rs ReturnStmt) GetType() NodeType {
	return ReturnStatement
}

func (rs ReturnStmt) GetToken() lexer.Token {
	return rs.Token
}

func (rs ReturnStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type AddStmt struct {
	Value      Node
	Identifier string
	Token      lexer.Token
}

func (as AddStmt) GetType() NodeType {
	return AddStatement
}

func (as AddStmt) GetToken() lexer.Token {
	return as.Token
}

func (as AddStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

type RemoveStmt struct {
	Value      Node
	Identifier string
	Token      lexer.Token
}

func (rs RemoveStmt) GetToken() lexer.Token {
	return rs.Token
}

func (rs RemoveStmt) GetType() NodeType {
	return RemoveStatement
}

func (rs RemoveStmt) GetValueType() PrimitiveValueType {
	return Undefined
}
