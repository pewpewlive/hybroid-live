package ast

import (
	"hybroid/lexer"
)

type Paths []string

type EnvironmentStmt struct {
	EnvType      *EnvTypeExpr
	Env          *EnvPathExpr
	Requirements Paths
}

func (as *EnvironmentStmt) GetType() NodeType {
	return EnvironmentStatement
}

func (as *EnvironmentStmt) GetToken() lexer.Token {
	return as.Env.GetToken()
}

func (as *EnvironmentStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type AssignmentStmt struct {
	Values      []Node
	Identifiers []Node
	Token       lexer.Token
}

func (as *AssignmentStmt) GetType() NodeType {
	return AssignmentStatement
}

func (as *AssignmentStmt) GetToken() lexer.Token {
	return as.Token
}

func (as *AssignmentStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type VariableDeclarationStmt struct {
	Identifiers []lexer.Token
	Types       []*TypeExpr
	Values      []Node
	IsLocal     bool
	Token       lexer.Token
}

func (vds *VariableDeclarationStmt) GetType() NodeType {
	return VariableDeclarationStatement
}

func (vds *VariableDeclarationStmt) GetToken() lexer.Token {
	return vds.Token
}

func (vds *VariableDeclarationStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type StructDeclarationStmt struct {
	Token       lexer.Token
	Name        lexer.Token
	Fields      []FieldDeclarationStmt
	Constructor *ConstructorStmt
	Methods     *[]MethodDeclarationStmt
	IsLocal     bool
}

func (sds *StructDeclarationStmt) GetType() NodeType {
	return StructureDeclarationStatement
}

func (sds *StructDeclarationStmt) GetToken() lexer.Token {
	return sds.Token
}

func (sds *StructDeclarationStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type EnumDeclarationStmt struct {
	Name    lexer.Token
	Fields  []lexer.Token
	IsLocal bool
}

func (sds *EnumDeclarationStmt) GetType() NodeType {
	return StructureDeclarationStatement
}

func (sds *EnumDeclarationStmt) GetToken() lexer.Token {
	return sds.Name
}

func (sds *EnumDeclarationStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type ConstructorStmt struct {
	Token  lexer.Token
	Body   []Node
	Return []*TypeExpr
	Params []Param
}

func (cs *ConstructorStmt) GetType() NodeType {
	return ConstructorStatement
}

func (cs *ConstructorStmt) GetToken() lexer.Token {
	return cs.Token
}

func (cs *ConstructorStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type FieldDeclarationStmt struct {
	Identifiers []lexer.Token
	Types       []*TypeExpr
	Values      []Node
	Token       lexer.Token
}

func (fds *FieldDeclarationStmt) GetType() NodeType {
	return FieldDeclarationStatement
}

func (fds *FieldDeclarationStmt) GetToken() lexer.Token {
	return fds.Token
}

func (fds *FieldDeclarationStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type Param struct {
	Type *TypeExpr
	Name lexer.Token
}

type FunctionDeclarationStmt struct {
	Name    lexer.Token
	Return  []*TypeExpr
	Params  []Param
	IsLocal bool
	Body    []Node
}

func (fds *FunctionDeclarationStmt) GetType() NodeType {
	return MethodDeclarationStatement
}

func (fds *FunctionDeclarationStmt) GetToken() lexer.Token {
	return fds.Name
}

func (fds *FunctionDeclarationStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type MethodDeclarationStmt struct {
	Owner   lexer.Token
	Name    lexer.Token
	Return  []*TypeExpr
	Params  []Param
	IsLocal bool
	Body    []Node
}

func (mds *MethodDeclarationStmt) GetType() NodeType {
	return FunctionDeclarationStatement
}

func (mds *MethodDeclarationStmt) GetToken() lexer.Token {
	return mds.Name
}

func (mds *MethodDeclarationStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type IfStmt struct {
	BoolExpr Node
	Body     []Node
	Elseifs  []*IfStmt
	Else     *IfStmt
	Token    lexer.Token
}

func (is *IfStmt) GetType() NodeType {
	return IfStatement
}

func (is *IfStmt) GetToken() lexer.Token {
	return is.Token
}

func (is *IfStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type MatchStmt struct {
	ExprToMatch Node
	Cases       []CaseStmt
	HasDefault  bool
}

type CaseStmt struct {
	Expression Node
	Body       []Node
}

func (ms *MatchStmt) GetType() NodeType {
	return MatchStatement
}

func (ms *MatchStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

func (ms *MatchStmt) GetToken() lexer.Token {
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

func (rs *RepeatStmt) GetType() NodeType {
	return RepeatStatement
}

func (rs *RepeatStmt) GetToken() lexer.Token {
	return rs.Token
}

func (rs *RepeatStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type WhileStmt struct {
	Condtion Node
	Body     []Node
	Token    lexer.Token
}

func (fs *WhileStmt) GetType() NodeType {
	return WhileStatement
}

func (fs *WhileStmt) GetToken() lexer.Token {
	return fs.Token
}

func (fs *WhileStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type ForStmt struct {
	Iterator         Node
	KeyValuePair     [2]*IdentifierExpr
	OrderedIteration bool
	Body             []Node
	Token            lexer.Token
}

func (fs *ForStmt) GetType() NodeType {
	return ForStatement
}

func (fs *ForStmt) GetToken() lexer.Token {
	return fs.Token
}

func (fs *ForStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type TickStmt struct {
	Variable IdentifierExpr
	Body     []Node
	Token    lexer.Token
}

func (ts *TickStmt) GetType() NodeType {
	return TickStatement
}

func (ts *TickStmt) GetToken() lexer.Token {
	return ts.Token
}

func (ts *TickStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type ReturnStmt struct {
	Args  []Node
	Token lexer.Token
}

func (rs *ReturnStmt) GetType() NodeType {
	return ReturnStatement
}

func (rs *ReturnStmt) GetToken() lexer.Token {
	return rs.Token
}

func (rs *ReturnStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type YieldStmt struct {
	Args  []Node
	Token lexer.Token
}

func (ys *YieldStmt) GetType() NodeType {
	return YieldStatement
}

func (ys *YieldStmt) GetToken() lexer.Token {
	return ys.Token
}

func (ys *YieldStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type BreakStmt struct {
	Token lexer.Token
}

func (bs *BreakStmt) GetType() NodeType {
	return BreakStatement
}

func (bs *BreakStmt) GetToken() lexer.Token {
	return bs.Token
}

func (bs *BreakStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type ContinueStmt struct {
	Token lexer.Token
}

func (cs *ContinueStmt) GetType() NodeType {
	return ContinueStatement
}

func (cs *ContinueStmt) GetToken() lexer.Token {
	return cs.Token
}

func (cs *ContinueStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type AddStmt struct {
	Value      Node
	Identifier string
	Token      lexer.Token
}

func (as *AddStmt) GetType() NodeType {
	return AddStatement
}

func (as *AddStmt) GetToken() lexer.Token {
	return as.Token
}

func (as *AddStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type RemoveStmt struct {
	Value      Node
	Identifier string
	Token      lexer.Token
}

func (rs *RemoveStmt) GetToken() lexer.Token {
	return rs.Token
}

func (rs *RemoveStmt) GetType() NodeType {
	return RemoveStatement
}

func (rs *RemoveStmt) GetValueType() PrimitiveValueType {
	return Unknown
}

type UseStmt struct {
	Path     *EnvPathExpr
}

func (us *UseStmt) GetToken() lexer.Token {
	return us.Path.GetToken()
}

func (us *UseStmt) GetType() NodeType {
	return UseStatement
}

func (us *UseStmt) GetValueType() PrimitiveValueType {
	return Unknown
}
