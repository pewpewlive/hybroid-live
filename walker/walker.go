package walker // THE ACTUAL WALKING

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"slices"
)

var LibraryEnvs = map[Library]*Environment{
	Pewpew: PewpewEnv,
	Fmath:  FmathEnv,
	Math:   MathEnv,
	String: StringEnv,
	Table:  TableEnv,
}

type Environment struct {
	Name        string
	luaPath     string // dynamic lua path
	hybroidPath string
	Type        ast.Env

	Scope Scope

	importedWalkers []*Walker // the walkers imported through UseStmt
	UsedLibraries   map[Library]bool
	UsedBuiltinVars []string

	Classes  map[string]*ClassVal
	Entities map[string]*EntityVal

	_envStmt *ast.EnvironmentDecl
}

func (e *Environment) AddRequirement(path string) bool {
	return e._envStmt.AddRequirement(path)
}

func (e *Environment) Requirements() []string {
	return e._envStmt.Requirements
}

func (e *Environment) AddBuiltinVar(name string) {
	if slices.Contains(e.UsedBuiltinVars, name) {
		return
	}

	e.UsedBuiltinVars = append(e.UsedBuiltinVars, name)
}

func NewEnvironment(hybroidPath, luaPath string) *Environment {
	scope := Scope{
		Tag:        &UntaggedTag{},
		Variables:  map[string]*VariableVal{},
		AliasTypes: make(map[string]*AliasType),
	}
	global := &Environment{
		hybroidPath: hybroidPath,
		luaPath:     luaPath,
		Type:        ast.InvalidEnv,
		Scope:       scope,
		UsedLibraries: map[Library]bool{
			Pewpew: false,
			Table:  false,
			String: false,
			Math:   false,
			Fmath:  false,
		},
		Classes:  map[string]*ClassVal{},
		Entities: map[string]*EntityVal{},
	}

	global.Scope.Environment = global
	return global
}

type Library int

const (
	Pewpew Library = iota
	Fmath
	Math
	String
	Table
)

type Walker struct {
	alerts.Collector

	// ENVIRONMENT SHOULD NEVER CHANGE ONCE INITIALIZED
	environment *Environment

	walkers map[string]*Walker
	program []ast.Node
	context Context
	Walked  bool
}

func (w *Walker) Alert(alertType alerts.Alert, args ...any) {
	w.Alert_(alertType, args...)
}

func (w *Walker) AlertI(alert alerts.Alert) {
	w.AlertI_(alert)
}

func (w *Walker) AlertSingle(alert alerts.Alert, token tokens.Token, args ...any) {
	args = append([]any{alerts.NewSingle(token)}, args...)
	w.Alert(alert, args...)
}

func (w *Walker) AlertMulti(alert alerts.Alert, start, end tokens.Token, args ...any) {
	args = append([]any{alerts.NewMulti(start, end)}, args...)
	w.Alert(alert, args...)
}

func NewWalker(hybroidPath, luaPath string) *Walker {
	return &Walker{
		environment: NewEnvironment(hybroidPath, luaPath),
		program:     []ast.Node{},
		context: Context{
			Value: &Unknown{},
		},
		Collector: alerts.NewCollector(),
	}
}

func (w *Walker) SetProgram(program []ast.Node) {
	w.program = program
}

func (w *Walker) Pass1(wlkrs map[string]*Walker) {
	w.walkers = wlkrs
	nodes := w.program

	if len(nodes) == 0 {
		return
	}

	if nodes[0].GetType() != ast.EnvironmentDeclaration {
		w.AlertSingle(&alerts.ExpectedEnvironment{}, nodes[0].GetToken())
		return
	}

	scope := &w.environment.Scope
	for i := range nodes {
		w.WalkNode(nodes[i], scope)
	}
}

func (w *Walker) Pass2() {
	nodes := w.program

	if len(nodes) == 0 {
		return
	}

	scope := &w.environment.Scope
	for i := range nodes {
		w.WalkNode2(&nodes[i], scope)
	}

	w.Walked = true
}

func (w *Walker) WalkNode(node ast.Node, scope *Scope) {
	switch node := node.(type) {
	case *ast.EnvironmentDecl:
		if w.environment.Name != "" {
			w.AlertSingle(&alerts.EnvironmentRedaclaration{}, node.GetToken())
			return
		}
		switch node.EnvType.Token.Lexeme {
		case "Level":
			node.EnvType.Type = ast.LevelEnv
		case "Mesh":
			node.EnvType.Type = ast.MeshEnv
		case "Sound":
			node.EnvType.Type = ast.SoundEnv
		default:
			w.AlertSingle(&alerts.InvalidEnvironmentType{}, node.EnvType.Token, node.EnvType.Token.Lexeme)
		}
		w.environment.Type = node.EnvType.Type
		w.environment.Name = node.Env.Path.Lexeme
		w.environment._envStmt = node
		if w2, ok := w.walkers[w.environment.Name]; ok {
			w.AlertSingle(&alerts.DuplicateEnvironmentNames{}, node.GetToken(), w.environment.hybroidPath, w2.environment.hybroidPath)
			return
		}

		w.walkers[w.environment.Name] = w
	default:
	}
}

func (w *Walker) WalkNode2(node *ast.Node, scope *Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentDecl:
	case *ast.VariableDecl:
		w.variableDeclaration(newNode, scope)
	case *ast.IfStmt:
		w.ifStatement(newNode, scope)
	case *ast.FunctionDecl:
		w.functionDeclaration(newNode, scope, Function)
	case *ast.ReturnStmt:
		w.returnStatement(newNode, scope)
	case *ast.YieldStmt:
		w.yieldStatement(newNode, scope)
	case *ast.BreakStmt:
		w.breakStatement(newNode, scope)
	case *ast.ContinueStmt:
		w.continueStatement(newNode, scope)
	case *ast.RepeatStmt:
		w.repeatStatement(newNode, scope)
	case *ast.WhileStmt:
		w.whileStatement(newNode, scope)
	case *ast.ForStmt:
		w.forStatement(newNode, scope)
	case *ast.TickStmt:
		w.tickStatement(newNode, scope)
	case *ast.CallExpr:
		val := w.GetNodeValue(&newNode.Caller, scope)
		w.CallExpr(val, node, scope)
	case *ast.EnvAccessExpr:
		_, newVersion := w.EnvAccessExpr(newNode)
		if newVersion != nil {
			*node = newVersion
		}
	case *ast.ClassDecl:
		w.classDeclaration(newNode, scope)
	case *ast.EnumDecl:
		w.enumDeclaration(newNode, scope)
	case *ast.MatchStmt:
		w.matchStatement(newNode, false, scope)
	case *ast.AssignmentStmt:
		w.assignmentStatement(newNode, scope)
	case *ast.UseStmt:
		w.useStatement(newNode, scope)
	case *ast.DestroyStmt:
		w.destroyStatement(newNode, scope)
	case *ast.SpawnExpr:
		w.SpawnExpr(newNode, scope)
	case *ast.NewExpr:
		w.NewExpr(newNode, scope)
	case *ast.AliasDecl:
		w.aliasDeclaration(newNode, scope)
	// case *ast.TypeDeclarationStmt:
	// 	TypeDeclarationStmt(newNode, scope)
	case *ast.Improper:
		// w.Error(newNode.GetToken(), "Improper statement: parser fault")
	case *ast.EntityDecl:
		w.entityDeclaration(newNode, scope)
	default:
		// w.Error(newNode.GetToken(), "Expected statement")
	}
}

func (w *Walker) GetNodeActualValue(node *ast.Node, scope *Scope) Value {
	val := w.GetNodeValue(node, scope)
	if val, ok := val.(*VariableVal); ok {
		return val.Value
	}

	return val
}

func (w *Walker) GetNodeValue(node *ast.Node, scope *Scope) Value {
	var val Value

	switch newNode := (*node).(type) {
	case *ast.LiteralExpr:
		val = w.LiteralExpr(newNode)
	case *ast.BinaryExpr:
		val = w.BinaryExpr(newNode, scope)
	case *ast.IdentifierExpr:
		val = w.IdentifierExpr(node, scope)
	case *ast.GroupExpr:
		val = w.GroupingExpr(newNode, scope)
	case *ast.ListExpr:
		val = w.ListExpr(newNode, scope)
	case *ast.UnaryExpr:
		val = w.UnaryExpr(newNode, scope)
	case *ast.CallExpr:
		callVal := w.GetNodeValue(&newNode.Caller, scope)
		localVal := w.CallExpr(callVal, node, scope)
		val = localVal
	case *ast.MapExpr:
		val = w.MapExpr(newNode, scope)
	case *ast.FunctionExpr:
		val = w.FunctionExpr(newNode, scope)
	case *ast.StructExpr:
		val = w.StructExpr(newNode, scope)
	case *ast.AccessExpr:
		val = w.AccessExpr(newNode, scope)
	case *ast.NewExpr:
		val = w.NewExpr(newNode, scope)
	case *ast.SelfExpr:
		val = w.SelfExpr(newNode, scope)
	case *ast.MatchExpr:
		val = w.MatchExpr(newNode, scope)
	case *ast.EntityExpr:
		val = w.EntityExpr(newNode, scope)
	case *ast.EnvAccessExpr:
		var newVersion ast.Node
		val, newVersion = w.EnvAccessExpr(newNode)
		if newVersion != nil {
			*node = newVersion
		}
	case *ast.SpawnExpr:
		val = w.SpawnExpr(newNode, scope)
	default:
		// w.Error(newNode.GetToken(), "Expected expression")
		return &Invalid{}
	}

	if field, ok := w.context.Node.(*ast.FieldExpr); ok {
		if w.context.Value.GetType().GetType() == Strct {
			field.Index = -1
			return val
		}
		if w.context.PewpewVarFound {
			field.Index = -1
			w.context.PewpewVarFound = false
			return val
		}
		if container, ok := w.context.Value.(FieldContainer); ok {
			_, field.Index, _ = container.ContainsField((*node).GetToken().Lexeme)
		}
	}
	return val
}

func (w *Walker) WalkBody(body *[]ast.Node, tag ExitableTag, scope *Scope) {
	endIndex := -1
	for i := range *body {
		if tag.GetIfExits(All) {
			w.AlertMulti(&alerts.UnreachableCode{}, (*body)[i].GetToken(), (*body)[len(*body)-1].GetToken())
			endIndex = i
			break
		}
		w.WalkNode2(&(*body)[i], scope)
	}
	if endIndex != -1 {
		*body = (*body)[:endIndex]
	}
}

func (w *Walker) TypeifyNodeList(nodes *[]ast.Node, scope *Scope) []Type {
	arguments := make([]Type, 0)
	for i := range *nodes {
		val := w.GetNodeValue(&(*nodes)[i], scope)
		if function, ok := val.(*FunctionVal); ok {
			arguments = append(arguments, function.Returns...)
		} else {
			arguments = append(arguments, val.GetType())
		}
	}
	return arguments
}

func (w *Walker) WalkParams(parameters []ast.FunctionParam, scope *Scope, declare func(name tokens.Token, value Value)) []Type {
	variadicParams := make(map[tokens.Token]int)
	params := make([]Type, 0)
	for i, param := range parameters {
		params = append(params, w.TypeExpr(param.Type, scope))
		if params[i].GetType() == Variadic {
			variadicParams[parameters[i].Name] = i
		}
		value := w.TypeToValue(params[i])
		declare(param.Name, value)
	}

	if len(variadicParams) > 1 {
		// w.Error(parameters[0].Name, "can only have one vartiadic parameter")
	} else if len(variadicParams) != 0 {
		// for k, v := range variadicParams {
		// 	if v != len(parameters)-1 {
		// 		w.Error(k, "variadic parameter should be last")
		// 	}
		// }
	}

	return params
}
