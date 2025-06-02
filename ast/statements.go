package ast

import "hybroid/tokens"

type AssignmentStmt struct {
	Values      []Node
	Identifiers []Node
	AssignOp    tokens.Token
	Token       tokens.Token
}

func (as *AssignmentStmt) GetType() NodeType      { return AssignmentStatement }
func (as *AssignmentStmt) GetToken() tokens.Token { return as.Token }

type DestroyStmt struct {
	Identifier  Node
	Args        []Node
	GenericArgs []*TypeExpr
	EntityName  string
	EnvName     string
	Token       tokens.Token
}

func (ne *DestroyStmt) GetType() NodeType      { return DestroyStatement }
func (ne *DestroyStmt) GetToken() tokens.Token { return ne.Token }

func (ne *DestroyStmt) GetGenerics() []*TypeExpr {
	return ne.GenericArgs
}

func (ne *DestroyStmt) GetArgs() []Node {
	return ne.Args
}

type IfStmt struct {
	Body

	BoolExpr Node
	Elseifs  []*IfStmt
	Else     *IfStmt
	Token    tokens.Token
}

func (is *IfStmt) GetType() NodeType      { return IfStatement }
func (is *IfStmt) GetToken() tokens.Token { return is.Token }

type CaseStmt struct {
	Body
	Expressions []Node
}

func (ms *CaseStmt) GetType() NodeType      { return CaseStatement }
func (ms *CaseStmt) GetToken() tokens.Token { return ms.Expressions[0].GetToken() }

type MatchStmt struct {
	Token       tokens.Token
	ExprToMatch Node
	Cases       []*CaseStmt
	HasDefault  bool
}

func (ms *MatchStmt) GetType() NodeType      { return MatchStatement }
func (ms *MatchStmt) GetToken() tokens.Token { return ms.Token }

type RepeatStmt struct {
	Body

	Iterator Node
	Skip     Node
	Start    Node
	Variable *IdentifierExpr
	Token    tokens.Token
}

func (rs *RepeatStmt) GetType() NodeType      { return RepeatStatement }
func (rs *RepeatStmt) GetToken() tokens.Token { return rs.Token }

type WhileStmt struct {
	Body

	Condition Node
	Token     tokens.Token
}

func (fs *WhileStmt) GetType() NodeType      { return WhileStatement }
func (fs *WhileStmt) GetToken() tokens.Token { return fs.Token }

type ForStmt struct {
	Body

	Iterator         Node
	First            *IdentifierExpr
	Second           *IdentifierExpr
	OrderedIteration bool
	Token            tokens.Token
}

func (fs *ForStmt) GetType() NodeType      { return ForStatement }
func (fs *ForStmt) GetToken() tokens.Token { return fs.Token }

type TickStmt struct {
	Body

	Variable *IdentifierExpr
	Token    tokens.Token
}

func (ts *TickStmt) GetType() NodeType      { return TickStatement }
func (ts *TickStmt) GetToken() tokens.Token { return ts.Token }

type ReturnStmt struct {
	Args  []Node
	Token tokens.Token
}

func (rs *ReturnStmt) GetType() NodeType      { return ReturnStatement }
func (rs *ReturnStmt) GetToken() tokens.Token { return rs.Token }

type YieldStmt struct {
	Args  []Node
	Token tokens.Token
}

func (ys *YieldStmt) GetType() NodeType      { return YieldStatement }
func (ys *YieldStmt) GetToken() tokens.Token { return ys.Token }

type BreakStmt struct {
	Token tokens.Token
}

func (bs *BreakStmt) GetType() NodeType      { return BreakStatement }
func (bs *BreakStmt) GetToken() tokens.Token { return bs.Token }

type ContinueStmt struct {
	Token tokens.Token
}

func (cs *ContinueStmt) GetType() NodeType      { return ContinueStatement }
func (cs *ContinueStmt) GetToken() tokens.Token { return cs.Token }

type AddStmt struct {
	Value      Node
	Identifier string
	Token      tokens.Token
}

func (as *AddStmt) GetType() NodeType      { return AddStatement }
func (as *AddStmt) GetToken() tokens.Token { return as.Token }

type RemoveStmt struct {
	Value      Node
	Identifier string
	Token      tokens.Token
}

func (rs *RemoveStmt) GetType() NodeType      { return RemoveStatement }
func (rs *RemoveStmt) GetToken() tokens.Token { return rs.Token }

type UseStmt struct {
	Token    tokens.Token
	PathExpr *EnvPathExpr
}

func (us *UseStmt) GetType() NodeType      { return UseStatement }
func (us *UseStmt) GetToken() tokens.Token { return us.Token }
