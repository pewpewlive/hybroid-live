package ast

import (
	"hybroid/lexer"
)

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
	Types       []*TypeExpr
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

type Param struct {
	Type lexer.Token
	Name lexer.Token
}

type FunctionDeclarationStmt struct {
	Name    lexer.Token
	Return  []lexer.Token
	Params  []Param
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

/*

let a = match ExprToMatch {
"a" => {

}
"b" => {

}

_ => {

}

}

let a = if condition {
	a =  1
}

local a
if condtion then
	a = 1
end

let o = 0
repeat to 10 with i {

	if i == 10 {
		o = i+20
	}else {
		o = 0
	}
}

local o

for i = 1, 10 do
	if i == 10 then
		o = i+20
		break
	else
	 o = 0
	 break
	end
end

*/

type MatchStmt struct {
	ExprToMatch Node
	Cases       []CaseStmt
}

type CaseStmt struct {
	Body  []Node
	Cases []Node
}

func (ms MatchStmt) GetType() NodeType {
	return MatchStatement
}

func (ms MatchStmt) GetValueType() PrimitiveValueType {
	return Undefined
}

func (ms MatchStmt) GetToken() lexer.Token {
	return ms.ExprToMatch.GetToken()
}

type RepeatStmt struct {
	Iterator Node
	Skip     Node
	Start    Node
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

type TickStmt struct {
	Variable IdentifierExpr
	Body     []Node
	Token    lexer.Token
}

func (ts TickStmt) GetType() NodeType {
	return TickStatement
}

func (ts TickStmt) GetToken() lexer.Token {
	return ts.Token
}

func (ts TickStmt) GetValueType() PrimitiveValueType {
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

type UseStmt struct {
	File     lexer.Token
	Variable IdentifierExpr
}

func (us UseStmt) GetToken() lexer.Token {
	return us.File
}

func (us UseStmt) GetType() NodeType {
	return UseStatement
}

func (us UseStmt) GetValueType() PrimitiveValueType {
	return Undefined
}
