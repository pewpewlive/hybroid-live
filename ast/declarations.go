package ast

import (
	"hybroid/tokens"
	"slices"
)

type EnvironmentDecl struct {
	EnvType      *EnvTypeExpr
	Env          *EnvPathExpr
	Requirements []string
}

func (ed *EnvironmentDecl) GetType() NodeType                { return EnvironmentDeclaration }
func (ed *EnvironmentDecl) GetToken() tokens.Token           { return ed.Env.GetToken() }
func (ed *EnvironmentDecl) GetValueType() PrimitiveValueType { return Invalid }

func (ed *EnvironmentDecl) AddRequirement(path string) bool {
	if slices.Contains(ed.Requirements, path) {
		return false
	}
	ed.Requirements = append(ed.Requirements, path)
	return true
}

type MacroDecl struct {
	Name      tokens.Token
	Params    []*IdentifierExpr
	MacroType MacroType
	Tokens    []tokens.Token
}

func (md *MacroDecl) GetType() NodeType                { return MacroDeclaration }
func (md *MacroDecl) GetToken() tokens.Token           { return md.Name }
func (md *MacroDecl) GetValueType() PrimitiveValueType { return Invalid }

type AliasDecl struct {
	IsPub bool
	Token tokens.Token
	Name  tokens.Token
	Type  *TypeExpr
}

func (ad *AliasDecl) GetType() NodeType                { return AliasDeclaration }
func (ad *AliasDecl) GetToken() tokens.Token           { return ad.Token }
func (ad *AliasDecl) GetValueType() PrimitiveValueType { return Invalid }

type EntityDecl struct {
	Token     tokens.Token
	Name      tokens.Token
	Fields    []VariableDecl
	Spawner   *EntityFunctionDecl
	Destroyer *EntityFunctionDecl
	Callbacks []*EntityFunctionDecl
	Methods   []MethodDecl
	IsPub     bool
}

func (ed *EntityDecl) GetType() NodeType                { return EntityDeclaration }
func (ed *EntityDecl) GetToken() tokens.Token           { return ed.Token }
func (ed *EntityDecl) GetValueType() PrimitiveValueType { return Invalid }

type EntityFunctionDecl struct {
	Body

	Type     EntityFunctionType
	Generics []*IdentifierExpr
	Params   []FunctionParam
	Returns  []*TypeExpr
	Token    tokens.Token
}

func (efd *EntityFunctionDecl) GetType() NodeType                { return EntityFunctionDeclaration }
func (efd *EntityFunctionDecl) GetToken() tokens.Token           { return efd.Token }
func (efd *EntityFunctionDecl) GetValueType() PrimitiveValueType { return Invalid }

type EnumDecl struct {
	Token  tokens.Token
	Name   tokens.Token
	Fields []*IdentifierExpr
	IsPub  bool
}

func (ed *EnumDecl) GetType() NodeType                { return EnumDeclaration }
func (ed *EnumDecl) GetToken() tokens.Token           { return ed.Name }
func (ed *EnumDecl) GetValueType() PrimitiveValueType { return Invalid }

type ConstructorDecl struct {
	Body

	Token    tokens.Token
	Params   []FunctionParam
	Generics []*IdentifierExpr
}

func (cd *ConstructorDecl) GetType() NodeType                { return ConstructorDeclaration }
func (cd *ConstructorDecl) GetToken() tokens.Token           { return cd.Token }
func (cd *ConstructorDecl) GetValueType() PrimitiveValueType { return Invalid }

type FunctionParam struct {
	Type *TypeExpr
	Name tokens.Token
}

type FunctionDecl struct {
	Body

	Token    tokens.Token
	Name     tokens.Token
	IsPub    bool
	Generics []*IdentifierExpr
	Params   []FunctionParam
	Returns  []*TypeExpr
}

func (fd *FunctionDecl) GetType() NodeType                { return FunctionDeclaration }
func (fd *FunctionDecl) GetToken() tokens.Token           { return fd.Token }
func (fd *FunctionDecl) GetValueType() PrimitiveValueType { return Invalid }

type MethodDecl struct {
	Body

	Owner    tokens.Token
	Name     tokens.Token
	Returns  []*TypeExpr
	Params   []FunctionParam
	Generics []*IdentifierExpr
	IsPub    bool
}

func (md *MethodDecl) GetType() NodeType                { return MethodDeclaration }
func (md *MethodDecl) GetToken() tokens.Token           { return md.Name }
func (md *MethodDecl) GetValueType() PrimitiveValueType { return Invalid }

type VariableDecl struct {
	Identifiers []*IdentifierExpr
	Type        *TypeExpr
	Expressions []Node
	IsPub       bool
	IsConst     bool
	Token       tokens.Token
}

func (vd *VariableDecl) GetType() NodeType                { return VariableDeclaration }
func (vd *VariableDecl) GetToken() tokens.Token           { return vd.Token }
func (vd *VariableDecl) GetValueType() PrimitiveValueType { return Invalid }

type ClassDecl struct {
	Token       tokens.Token
	Name        tokens.Token
	Constructor *ConstructorDecl
	Fields      []VariableDecl
	Methods     []MethodDecl
	IsPub       bool
}

func (cd *ClassDecl) GetType() NodeType                { return ClassDeclaration }
func (cd *ClassDecl) GetToken() tokens.Token           { return cd.Token }
func (cd *ClassDecl) GetValueType() PrimitiveValueType { return Invalid }
