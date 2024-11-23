package ast

import "hybroid/tokens"

type EnvironmentStmt struct {
	EnvType      *EnvTypeExpr
	Env          *EnvPathExpr
	Requirements Paths
}

func (as *EnvironmentStmt) AddRequirement(path string) bool {
	for i := range as.Requirements {
		if as.Requirements[i] == path {
			return false
		}
	}
	as.Requirements = append(as.Requirements, path)
	return true
}

func (es *EnvironmentStmt) GetType() NodeType {
	return EnvironmentStatement
}

func (es *EnvironmentStmt) GetToken() tokens.Token {
	return es.Env.GetToken()
}

func (es *EnvironmentStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type AssignmentStmt struct {
	Values      []Node
	Identifiers []Node
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

type MacroDeclarationStmt struct {
	Name      tokens.Token
	Params    []tokens.Token
	MacroType MacroType
	Tokens    []tokens.Token
}

func (mds *MacroDeclarationStmt) GetType() NodeType {
	return MacroDeclarationStatement
}

func (mds *MacroDeclarationStmt) GetToken() tokens.Token {
	return mds.Name
}

func (mds *MacroDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type AliasDeclarationStmt struct {
	IsLocal     bool
	Token       tokens.Token
	Alias       tokens.Token
	AliasedType *TypeExpr
}

func (vds *AliasDeclarationStmt) GetType() NodeType {
	return AliasDeclarationStatement
}

func (vds *AliasDeclarationStmt) GetToken() tokens.Token {
	return vds.Token
}

func (vds *AliasDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

// type TypeDeclarationStmt struct {
// 	AliasedType *TypeExpr
// 	Alias tokens.Token
// 	Token tokens.Token
// }

// func (vds *TypeDeclarationStmt) GetType() NodeType {
// 	return TypeDeclarationStatement
// }

// func (vds *TypeDeclarationStmt) GetToken() tokens.Token {
// 	return vds.Token
// }

// func (vds *TypeDeclarationStmt) GetValueType() PrimitiveValueType {
// 	return Invalid
// }

type VariableDeclarationStmt struct {
	Identifiers []tokens.Token
	Type        *TypeExpr
	Values      []Node
	IsLocal     bool
	IsConst     bool
	Token       tokens.Token
}

func (vds *VariableDeclarationStmt) GetType() NodeType {
	return VariableDeclarationStatement
}

func (vds *VariableDeclarationStmt) GetToken() tokens.Token {
	return vds.Token
}

func (vds *VariableDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type ClassDeclarationStmt struct {
	Token       tokens.Token
	Name        tokens.Token
	Fields      []FieldDeclarationStmt
	Constructor *ConstructorStmt
	Methods     []MethodDeclarationStmt
	IsLocal     bool
}

func (sds *ClassDeclarationStmt) GetType() NodeType {
	return ClassDeclarationStatement
}

func (sds *ClassDeclarationStmt) GetToken() tokens.Token {
	return sds.Token
}

func (sds *ClassDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type DestroyStmt struct {
	Identifier Node
	Args       []Node
	Generics   []*TypeExpr
	EntityName string
	EnvName    string
	Token      tokens.Token
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

type EntityDeclarationStmt struct {
	Token     tokens.Token
	Name      tokens.Token
	Fields    []FieldDeclarationStmt
	Spawner   *EntityFunctionDeclarationStmt
	Destroyer *EntityFunctionDeclarationStmt
	Callbacks []*EntityFunctionDeclarationStmt
	Methods   []MethodDeclarationStmt
	IsLocal   bool
}

func (sds *EntityDeclarationStmt) GetType() NodeType {
	return EntityDeclarationStatement
}

func (sds *EntityDeclarationStmt) GetToken() tokens.Token {
	return sds.Token
}

func (sds *EntityDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type EntityFunctionDeclarationStmt struct {
	Type     EntityFunctionType
	Generics []*IdentifierExpr
	Params   []Param
	Returns  []*TypeExpr
	Body     []Node
	Token    tokens.Token
}

func (eds *EntityFunctionDeclarationStmt) GetType() NodeType {
	return EntityFunctionDeclarationStatemet
}

func (eds *EntityFunctionDeclarationStmt) GetToken() tokens.Token {
	return eds.Token
}

func (eds *EntityFunctionDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type EnumDeclarationStmt struct {
	Name    tokens.Token
	Fields  []tokens.Token
	IsLocal bool
}

func (eds *EnumDeclarationStmt) GetType() NodeType {
	return EnumDeclarationStatement
}

func (eds *EnumDeclarationStmt) GetToken() tokens.Token {
	return eds.Name
}

func (eds *EnumDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type ConstructorStmt struct {
	Token    tokens.Token
	Body     []Node
	Return   []*TypeExpr
	Params   []Param
	Generics []*IdentifierExpr
}

func (cs *ConstructorStmt) GetType() NodeType {
	return ConstructorStatement
}

func (cs *ConstructorStmt) GetToken() tokens.Token {
	return cs.Token
}

func (cs *ConstructorStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type FieldDeclarationStmt struct {
	Identifiers []tokens.Token
	Type        *TypeExpr
	Values      []Node
	Token       tokens.Token
}

func (fds *FieldDeclarationStmt) GetType() NodeType {
	return FieldDeclarationStatement
}

func (fds *FieldDeclarationStmt) GetToken() tokens.Token {
	return fds.Token
}

func (fds *FieldDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type Param struct {
	Type *TypeExpr
	Name tokens.Token
}

type FunctionDeclarationStmt struct {
	Name          tokens.Token
	Return        []*TypeExpr
	GenericParams []*IdentifierExpr
	Params        []Param
	IsLocal       bool
	Body          []Node
}

func (fds *FunctionDeclarationStmt) GetType() NodeType {
	return FunctionDeclarationStatement
}

func (fds *FunctionDeclarationStmt) GetToken() tokens.Token {
	return fds.Name
}

func (fds *FunctionDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
}

type MethodDeclarationStmt struct {
	Owner    tokens.Token
	Name     tokens.Token
	Return   []*TypeExpr
	Params   []Param
	Generics []*IdentifierExpr
	IsLocal  bool
	Body     []Node
}

func (mds *MethodDeclarationStmt) GetType() NodeType {
	return FunctionDeclarationStatement
}

func (mds *MethodDeclarationStmt) GetToken() tokens.Token {
	return mds.Name
}

func (mds *MethodDeclarationStmt) GetValueType() PrimitiveValueType {
	return Invalid
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
	return ms.ExprToMatch.GetToken()
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
	Condtion Node
	Body     []Node
	Token    tokens.Token
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
