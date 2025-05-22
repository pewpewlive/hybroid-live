package ast

import "hybroid/tokens"

type AssignmentStmt struct {
	Values      []Node
	Identifiers []Node
	AssignOp    tokens.Token
	Token       tokens.Token
}

func (as *AssignmentStmt) GetType() NodeType {
	return AssignmentStatement
}

func (as *AssignmentStmt) GetToken() tokens.Token {
	return as.Token
}

func (as *AssignmentStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type DestroyStmt struct {
	Identifier  Node
	Args        []Node
	GenericArgs []*TypeExpr
	EntityName  string
	EnvName     string
	Token       tokens.Token
}

func (ne *DestroyStmt) GetType() NodeType {
	return DestroyStatement
}

func (ne *DestroyStmt) GetToken() tokens.Token {
	return ne.Token
}

func (ne *DestroyStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

func (ne *DestroyStmt) GetGenerics() []*TypeExpr {
	return ne.GenericArgs
}

func (ne *DestroyStmt) GetArgs() []Node {
	return ne.Args
}

type IfStmt struct {
	BoolExpr Node
	Body     []Node
	Elseifs  []*IfStmt
	Else     *IfStmt
	Token    tokens.Token
}

func (is *IfStmt) GetType() NodeType {
	return IfStatement
}

func (is *IfStmt) GetToken() tokens.Token {
	return is.Token
}

func (is *IfStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type CaseStmt struct {
	Expression Node
	Body       []Node
}

type MatchStmt struct {
	Token       tokens.Token
	ExprToMatch Node
	Cases       []CaseStmt
	HasDefault  bool
}

func (ms *MatchStmt) GetType() NodeType {
	return MatchStatement
}

func (ms *MatchStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

func (ms *MatchStmt) GetToken() tokens.Token {
	return ms.Token
}

type RepeatStmt struct {
	Iterator Node
	Skip     Node
	Start    Node
	Variable *IdentifierExpr
	Body     []Node
	Token    tokens.Token
}

func (rs *RepeatStmt) GetType() NodeType {
	return RepeatStatement
}

func (rs *RepeatStmt) GetToken() tokens.Token {
	return rs.Token
}

func (rs *RepeatStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type WhileStmt struct {
	Condition Node
	Body      []Node
	Token     tokens.Token
}

func (fs *WhileStmt) GetType() NodeType {
	return WhileStatement
}

func (fs *WhileStmt) GetToken() tokens.Token {
	return fs.Token
}

func (fs *WhileStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type ForStmt struct {
	Iterator         Node
	First            *IdentifierExpr
	Second           *IdentifierExpr
	OrderedIteration bool
	Body             []Node
	Token            tokens.Token
}

func (fs *ForStmt) GetType() NodeType {
	return ForStatement
}

func (fs *ForStmt) GetToken() tokens.Token {
	return fs.Token
}

func (fs *ForStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type TickStmt struct {
	Variable *IdentifierExpr
	Body     []Node
	Token    tokens.Token
}

func (ts *TickStmt) GetType() NodeType {
	return TickStatement
}

func (ts *TickStmt) GetToken() tokens.Token {
	return ts.Token
}

func (ts *TickStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type ReturnStmt struct {
	Args  []Node
	Token tokens.Token
}

func (rs *ReturnStmt) GetType() NodeType {
	return ReturnStatement
}

func (rs *ReturnStmt) GetToken() tokens.Token {
	return rs.Token
}

func (rs *ReturnStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type YieldStmt struct {
	Args  []Node
	Token tokens.Token
}

func (ys *YieldStmt) GetType() NodeType {
	return YieldStatement
}

func (ys *YieldStmt) GetToken() tokens.Token {
	return ys.Token
}

func (ys *YieldStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type BreakStmt struct {
	Token tokens.Token
}

func (bs *BreakStmt) GetType() NodeType {
	return BreakStatement
}

func (bs *BreakStmt) GetToken() tokens.Token {
	return bs.Token
}

func (bs *BreakStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type ContinueStmt struct {
	Token tokens.Token
}

func (cs *ContinueStmt) GetType() NodeType {
	return ContinueStatement
}

func (cs *ContinueStmt) GetToken() tokens.Token {
	return cs.Token
}

func (cs *ContinueStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type AddStmt struct {
	Value      Node
	Identifier string
	Token      tokens.Token
}

func (as *AddStmt) GetType() NodeType {
	return AddStatement
}

func (as *AddStmt) GetToken() tokens.Token {
	return as.Token
}

func (as *AddStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type RemoveStmt struct {
	Value      Node
	Identifier string
	Token      tokens.Token
}

func (rs *RemoveStmt) GetToken() tokens.Token {
	return rs.Token
}

func (rs *RemoveStmt) GetType() NodeType {
	return RemoveStatement
}

func (rs *RemoveStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type UseStmt struct {
	Path *EnvPathExpr
}

func (us *UseStmt) GetToken() tokens.Token {
	return us.Path.GetToken()
}

func (us *UseStmt) GetType() NodeType {
	return UseStatement
}

func (us *UseStmt) GetValueType() PrimitiveValueType {
	return Invalid
}
