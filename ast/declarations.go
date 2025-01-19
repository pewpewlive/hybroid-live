package ast

import "hybroid/tokens"

type EnvironmentDecl struct {
	EnvType      *EnvTypeExpr
	Env          *EnvPathExpr
	Requirements Paths
}

func (ed *EnvironmentDecl) AddRequirement(path string) bool {
	for i := range ed.Requirements {
		if ed.Requirements[i] == path {
			return false
		}
	}
	ed.Requirements = append(ed.Requirements, path)
	return true
}

func (ed *EnvironmentDecl) GetType() NodeType {
	return EnvironmentDeclaration
}

func (ed *EnvironmentDecl) GetToken() tokens.Token {
	return ed.Env.GetToken()
}

func (ed *EnvironmentDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type MacroDecl struct {
	Name      tokens.Token
	Params    []tokens.Token
	MacroType MacroType
	Tokens    []tokens.Token
}

func (md *MacroDecl) GetType() NodeType {
	return MacroDeclaration
}

func (md *MacroDecl) GetToken() tokens.Token {
	return md.Name
}

func (md *MacroDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type AliasDecl struct {
	IsLocal     bool
	Token       tokens.Token
	Alias       tokens.Token
	AliasedType *TypeExpr
}

func (ad *AliasDecl) GetType() NodeType {
	return AliasDeclaration
}

func (ad *AliasDecl) GetToken() tokens.Token {
	return ad.Token
}

func (ad *AliasDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type EntityDecl struct {
	Token     tokens.Token
	Name      tokens.Token
	Fields    []FieldDecl
	Spawner   *EntityFunctionDecl
	Destroyer *EntityFunctionDecl
	Callbacks []*EntityFunctionDecl
	Methods   []MethodDecl
	IsLocal   bool
}

func (ed *EntityDecl) GetType() NodeType {
	return EntityDeclaration
}

func (ed *EntityDecl) GetToken() tokens.Token {
	return ed.Token
}

func (ed *EntityDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type EntityFunctionDecl struct {
	Type     EntityFunctionType
	Generics []*IdentifierExpr
	Params   []Param
	Returns  []*TypeExpr
	Body     []Node
	Token    tokens.Token
}

func (efd *EntityFunctionDecl) GetType() NodeType {
	return EntityFunctionDeclaration
}

func (efd *EntityFunctionDecl) GetToken() tokens.Token {
	return efd.Token
}

func (efd *EntityFunctionDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type EnumDecl struct {
	Name    tokens.Token
	Fields  []tokens.Token
	IsLocal bool
}

func (ed *EnumDecl) GetType() NodeType {
	return EnumDeclaration
}

func (ed *EnumDecl) GetToken() tokens.Token {
	return ed.Name
}

func (ed *EnumDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type ConstructorDecl struct {
	Token    tokens.Token
	Body     []Node
	Return   []*TypeExpr
	Params   []Param
	Generics []*IdentifierExpr
}

func (cd *ConstructorDecl) GetType() NodeType {
	return ConstructorStatement
}

func (cd *ConstructorDecl) GetToken() tokens.Token {
	return cd.Token
}

func (cd *ConstructorDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type FieldDecl struct {
	Identifiers []tokens.Token
	Type        *TypeExpr
	Values      []Node
	Token       tokens.Token
}

func (fd *FieldDecl) GetType() NodeType {
	return FieldDeclaration
}

func (fd *FieldDecl) GetToken() tokens.Token {
	return fd.Token
}

func (fd *FieldDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type Param struct {
	Type *TypeExpr
	Name tokens.Token
}

type FunctionDecl struct {
	Name          tokens.Token
	Return        []*TypeExpr
	GenericParams []*IdentifierExpr
	Params        []Param
	IsLocal       bool
	Body          []Node
}

func (fd *FunctionDecl) GetType() NodeType {
	return FunctionDeclaration
}

func (fd *FunctionDecl) GetToken() tokens.Token {
	return fd.Name
}

func (fd *FunctionDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type MethodDecl struct {
	Owner    tokens.Token
	Name     tokens.Token
	Return   []*TypeExpr
	Params   []Param
	Generics []*IdentifierExpr
	IsLocal  bool
	Body     []Node
}

func (md *MethodDecl) GetType() NodeType {
	return FunctionDeclaration
}

func (md *MethodDecl) GetToken() tokens.Token {
	return md.Name
}

func (md *MethodDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type VariableDecl struct {
	Identifiers []tokens.Token
	Type        *TypeExpr
	Values      []Node
	IsLocal     bool
	IsConst     bool
	Token       tokens.Token
}

func (vd *VariableDecl) GetType() NodeType {
	return VariableDeclaration
}

func (vd *VariableDecl) GetToken() tokens.Token {
	return vd.Token
}

func (vd *VariableDecl) GetValueType() PrimitiveValueType {
	return Invalid
}

type ClassDecl struct {
	Token       tokens.Token
	Name        tokens.Token
	Fields      []FieldDecl
	Constructor *ConstructorDecl
	Methods     []MethodDecl
	IsLocal     bool
}

func (cd *ClassDecl) GetType() NodeType {
	return ClassDeclaration
}

func (cd *ClassDecl) GetToken() tokens.Token {
	return cd.Token
}

func (cd *ClassDecl) GetValueType() PrimitiveValueType {
	return Invalid
}
