package ast

import (
	"hybroid/tokens"
	"strings"
)

type EnvTypeExpr struct {
	Type  Env
	Token tokens.Token
}

func (ete *EnvTypeExpr) GetType() NodeType      { return EnvironmentTypeExpression }
func (ete *EnvTypeExpr) GetToken() tokens.Token { return ete.Token }

type EnvPathExpr struct {
	Path tokens.Token
}

func (epe *EnvPathExpr) GetType() NodeType      { return EnvironmentPathExpression }
func (epe *EnvPathExpr) GetToken() tokens.Token { return epe.Path }

func (epe *EnvPathExpr) Combine(token tokens.Token) {
	epe.Path.Lexeme += ":" + token.Lexeme
	epe.Path.Column.SetEnd(token.Column.End)
}

type EnvAccessExpr struct {
	PathExpr *EnvPathExpr
	Accessed *IdentifierExpr
}

func (eae *EnvAccessExpr) GetType() NodeType      { return EnvironmentAccessExpression }
func (eae *EnvAccessExpr) GetToken() tokens.Token { return eae.Accessed.GetToken() }

// type MacroCallExpr struct {
// 	Caller *CallExpr
// }

// func (mce *MacroCallExpr) GetType() NodeType                { return MacroCallExpression }
// func (mce *MacroCallExpr) GetToken() tokens.Token           { return mce.Caller.GetToken() }

type EntityAccessExpr struct {
	Expr       Node
	EntityName string
	EnvName    string
}

func (pe *EntityAccessExpr) GetType() NodeType      { return EntityAccessExpression }
func (pe *EntityAccessExpr) GetToken() tokens.Token { return pe.Expr.GetToken() }

type LiteralExpr struct {
	Value string
	Token tokens.Token
}

func (le *LiteralExpr) GetType() NodeType      { return LiteralExpression }
func (le *LiteralExpr) GetToken() tokens.Token { return le.Token }

type UnaryExpr struct {
	Value    Node
	Operator tokens.Token
}

func (ue *UnaryExpr) GetType() NodeType      { return UnaryExpression }
func (ue *UnaryExpr) GetToken() tokens.Token { return ue.Value.GetToken() }

type TypeExpr struct {
	WrappedTypes []*TypeExpr
	Name         Node
	Params       []*TypeExpr
	Returns      []*TypeExpr
	Fields       []FunctionParam
	IsVariadic   bool
}

func (te *TypeExpr) GetType() NodeType      { return TypeExpression }
func (te *TypeExpr) GetToken() tokens.Token { return te.Name.GetToken() }

type EntityEvaluationExpr struct {
	Expr               Node
	Type               *TypeExpr
	ConvertedVarName   *tokens.Token
	OfficialEntityType bool
	EntityName         string
	EnvName            string
	Operator           tokens.Token
	Token              tokens.Token
}

func (eee *EntityEvaluationExpr) GetType() NodeType      { return EntityEvaluationExpression }
func (eee *EntityEvaluationExpr) GetToken() tokens.Token { return eee.Token }

type GroupExpr struct {
	Expr  Node
	Token tokens.Token
}

func (ge *GroupExpr) GetType() NodeType      { return GroupExpression }
func (ge *GroupExpr) GetToken() tokens.Token { return ge.Token }

type BinaryExpr struct {
	Left, Right Node
	Operator    tokens.Token
}

func (be *BinaryExpr) GetType() NodeType      { return BinaryExpression }
func (be *BinaryExpr) GetToken() tokens.Token { return be.Operator }

type CallNode interface {
	GetReturnAmount() int
}

type CallExpr struct {
	Caller       Node
	GenericArgs  []*TypeExpr
	Args         []Node
	ReturnAmount int
}

func (ce *CallExpr) GetType() NodeType      { return CallExpression }
func (ce *CallExpr) GetToken() tokens.Token { return ce.Caller.GetToken() }

func (ce *CallExpr) GetGenerics() []*TypeExpr { return ce.GenericArgs }
func (ce *CallExpr) GetArgs() []Node          { return ce.Args }
func (ce *CallExpr) GetCaller() Node          { return ce.Caller }
func (ce *CallExpr) GetReturnAmount() int     { return ce.ReturnAmount }

type MethodCallExpr struct {
	MethodInfo
	Caller       Node
	GenericArgs  []*TypeExpr
	Args         []Node
	ReturnAmount int
}

func (mce *MethodCallExpr) GetType() NodeType      { return MethodCallExpression }
func (mce *MethodCallExpr) GetToken() tokens.Token { return mce.Caller.GetToken() }

func (mce *MethodCallExpr) GetGenerics() []*TypeExpr { return mce.GenericArgs }
func (mce *MethodCallExpr) GetArgs() []Node          { return mce.Args }
func (mce *MethodCallExpr) GetCaller() Node          { return mce.Caller }
func (mce *MethodCallExpr) GetReturnAmount() int     { return mce.ReturnAmount }

type FunctionExpr struct {
	Body

	Token    tokens.Token
	Returns  []*TypeExpr
	Params   []FunctionParam
	Generics []*IdentifierExpr
}

func (fe *FunctionExpr) GetType() NodeType      { return FunctionExpression }
func (fe *FunctionExpr) GetToken() tokens.Token { return fe.Token }

type StructExpr struct {
	Token       tokens.Token
	Fields      []*IdentifierExpr
	Expressions []Node
}

func (ase *StructExpr) GetType() NodeType      { return StructExpression }
func (ase *StructExpr) GetToken() tokens.Token { return ase.Token }

type MatchExpr struct {
	MatchStmt    MatchStmt
	ReturnAmount int
}

func (me *MatchExpr) GetType() NodeType      { return MatchExpression }
func (me *MatchExpr) GetToken() tokens.Token { return me.MatchStmt.GetToken() }

type SelfExpr struct {
	Token      tokens.Token
	EntityName string
	Type       MethodCallType
}

func (se *SelfExpr) GetType() NodeType      { return SelfExpression }
func (se *SelfExpr) GetToken() tokens.Token { return se.Token }

type NewExpr struct {
	Type        *TypeExpr
	GenericArgs []*TypeExpr
	Args        []Node
	Token       tokens.Token
	EnvName     string
}

func (ne *NewExpr) GetType() NodeType      { return NewExpession }
func (ne *NewExpr) GetToken() tokens.Token { return ne.Token }

func (ne *NewExpr) GetGenerics() []*TypeExpr { return ne.GenericArgs }
func (ne *NewExpr) GetCaller() Node          { return ne.Type }
func (ne *NewExpr) GetArgs() []Node          { return ne.Args }

type SpawnExpr struct {
	Type        *TypeExpr
	Args        []Node
	GenericArgs []*TypeExpr
	Token       tokens.Token
	EnvName     string
}

func (ne *SpawnExpr) GetType() NodeType      { return SpawnExpression }
func (ne *SpawnExpr) GetToken() tokens.Token { return ne.Token }

func (ne *SpawnExpr) GetGenerics() []*TypeExpr { return ne.GenericArgs }
func (ne *SpawnExpr) GetCaller() Node          { return ne.Type }
func (ne *SpawnExpr) GetArgs() []Node          { return ne.Args }

type IdentifierType int

const (
	Other IdentifierType = iota
	Raw
)

type AccessExpr struct {
	Start    Node
	Accessed []Node
}

func (ae *AccessExpr) GetType() NodeType      { return ae.Accessed[len(ae.Accessed)-1].GetType() }
func (ae *AccessExpr) GetToken() tokens.Token { return ae.Start.GetToken() }

type FieldExpr struct {
	Field Node
	Index int
}

func (fe *FieldExpr) GetType() NodeType      { return FieldExpression }
func (fe *FieldExpr) GetToken() tokens.Token { return fe.Field.GetToken() }

type MemberExpr struct {
	Member Node
	IsList bool
}

func (me *MemberExpr) GetType() NodeType      { return MemberExpression }
func (me *MemberExpr) GetToken() tokens.Token { return me.Member.GetToken() }

func (me *MemberExpr) GetIdentifier() Node { return me.Member }

type MapExpr struct {
	Token        tokens.Token
	Type         *TypeExpr
	KeyValueList []Property
}

func (me *MapExpr) GetType() NodeType      { return MapExpression }
func (me *MapExpr) GetToken() tokens.Token { return me.Token }

type ListExpr struct {
	List  []Node
	Type  *TypeExpr
	Token tokens.Token
}

func (le *ListExpr) GetType() NodeType      { return ListExpression }
func (le *ListExpr) GetToken() tokens.Token { return le.Token }

type IdentifierExpr struct {
	Type IdentifierType
	Name tokens.Token
}

func (ie *IdentifierExpr) GetType() NodeType      { return Identifier }
func (ie *IdentifierExpr) GetToken() tokens.Token { return ie.Name }

type Improper struct {
	Token tokens.Token
	Type  NodeType
}

func NewImproper(token tokens.Token, nodeType NodeType) *Improper {
	return &Improper{
		Token: token,
		Type:  nodeType,
	}
}

func (i *Improper) GetType() NodeType      { return NA }
func (i *Improper) GetToken() tokens.Token { return i.Token }

func IsImproper(improper Node, nodeType NodeType) bool {
	return improper.GetType() == NA && improper.(*Improper).Type == nodeType
}

func IsImproperNotStatement(improper Node) bool {
	if improper.GetType() != NA {
		return false
	}
	improperType := improper.(*Improper).Type
	str := strings.ToLower(string(improperType))
	return !strings.Contains(str, "statement") &&
		!strings.Contains(str, "declaration") &&
		!strings.Contains(str, "call") &&
		improperType != NewExpession &&
		improperType != SpawnExpression
}
