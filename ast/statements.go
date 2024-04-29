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
	return 0
}

type VariableDeclarationStmt struct {
	Identifiers []lexer.Token
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
	return 0
}

type StructDeclarationStmt struct {
	Token lexer.Token
	Name  lexer.Token
	Body  *[]Node
	IsLocal bool
}

func (sds StructDeclarationStmt) GetType() NodeType {
	return StructureDeclarationStatement
}

func (sds StructDeclarationStmt) GetToken() lexer.Token {
	return sds.Token
}

func (sds StructDeclarationStmt) GetValueType() PrimitiveValueType {
	return 0
}

type FieldDeclarationStmt struct {
	Identifiers []lexer.Token
	Types       []*TypeExpr
	Values      []Node
	IsLocal     bool
	Token       lexer.Token
}

func (f FieldDeclarationStmt) GetType() NodeType {
	return FieldDeclarationStatement
}

func (f FieldDeclarationStmt) GetToken() lexer.Token {
	return f.Token
}

func (f FieldDeclarationStmt) GetValueType() PrimitiveValueType {
	return 0
}

type Param struct {
	Type TypeExpr
	Name lexer.Token
}

type FunctionDeclarationStmt struct {
	Name    lexer.Token
	Return  []TypeExpr
	Params  []Param
	IsLocal bool
	Body    []Node
}

func (fds FunctionDeclarationStmt) GetType() NodeType {
	return MethodDeclarationStatement
}

func (fds FunctionDeclarationStmt) GetToken() lexer.Token {
	return fds.Name
}

func (fds FunctionDeclarationStmt) GetValueType() PrimitiveValueType {
	return 0
}


type MethodDeclarationStmt struct {
	Name    lexer.Token
	Return  []TypeExpr
	Params  []Param
	IsLocal bool
	Body    []Node
}

func (fds MethodDeclarationStmt) GetType() NodeType {
	return FunctionDeclarationStatement
}

func (fds MethodDeclarationStmt) GetToken() lexer.Token {
	return fds.Name
}

func (fds MethodDeclarationStmt) GetValueType() PrimitiveValueType {
	return 0
}

type IfStmt struct {
	BoolExpr Node
	Body     []Node
	Elseifs  []*IfStmt
	Else     *IfStmt
	Token    lexer.Token
}

func (is IfStmt) GetType() NodeType {
	return IfStatement
}

func (is IfStmt) GetToken() lexer.Token {
	return is.Token
}

func (is IfStmt) GetValueType() PrimitiveValueType {
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
}
