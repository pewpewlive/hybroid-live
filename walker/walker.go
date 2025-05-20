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
	Name            string
	luaPath         string // dynamic lua path
	hybroidPath     string
	Type            ast.EnvType
	Scope           Scope
	UsedWalkers     []*Walker
	UsedLibraries   map[Library]bool
	UsedBuiltinVars []string
	Structs         map[string]*ClassVal
	Entities        map[string]*EntityVal
	CustomTypes     map[string]*CustomType
	AliasTypes      map[string]*AliasType
	EnvStmt         *ast.EnvironmentDecl
}

func (e *Environment) AddBuiltinVar(name string) {
	if slices.Contains(e.UsedBuiltinVars, name) {
		return
	}

	e.UsedBuiltinVars = append(e.UsedBuiltinVars, name)
}

func NewEnvironment(hybroidPath, luaPath string) *Environment {
	scope := Scope{
		Tag:       &UntaggedTag{},
		Variables: map[string]*VariableVal{},
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
		Structs:     map[string]*ClassVal{},
		Entities:    map[string]*EntityVal{},
		CustomTypes: map[string]*CustomType{},
		AliasTypes:  make(map[string]*AliasType),
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

	currentEnvironment *Environment
	environment        *Environment
	walkers            map[string]*Walker
	nodes              []ast.Node
	context            Context
	Walked             bool
}

func (w *Walker) Alert(alertType alerts.Alert, args ...any) {
	w.Alert_(alertType, args...)
}

func (w *Walker) AlertI(alert alerts.Alert) {
	w.AlertI_(alert)
}

func (w *Walker) AlertSingle(alertType alerts.Alert, token tokens.Token, args ...any) {
	args = append([]any{alerts.NewSingle(token)}, args...)
	w.Alert(alertType, args...)
}

func NewWalker(hybroidPath, luaPath string) *Walker {
	walker := &Walker{
		environment: NewEnvironment(hybroidPath, luaPath),
		nodes:       []ast.Node{},
		context: Context{
			Node:   &ast.Improper{},
			Value:  &Unknown{},
			Value2: &Unknown{},
		},
		Collector: alerts.NewCollector(),
	}
	walker.currentEnvironment = walker.environment
	return walker
}

func (w *Walker) SetProgram(program []ast.Node) {
	w.nodes = program
}

func (w *Walker) Action(wlkrs map[string]*Walker) {
	w.walkers = wlkrs
	nodes := w.nodes

	if len(nodes) == 0 {
		return
	}

	if nodes[0].GetType() != ast.EnvironmentDeclaration {
		w.AlertSingle(&alerts.ExpectedEnvironment{}, nodes[0].GetToken())
		return
	}

	scope := &w.environment.Scope
	for i := range nodes {
		w.WalkNode(&nodes[i], scope)
	}

	// for i := range w.nodes {
	// 	w.WalkNode2(&w.nodes[i], scope)
	// }

	w.Walked = true
}

func (w *Walker) WalkNode(node *ast.Node, scope *Scope) {
	switch node := (*node).(type) {
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
		w.environment.EnvStmt = node
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
		w.VariableDeclarationStmt(newNode, scope)
	case *ast.IfStmt:
		w.IfStmt(newNode, scope)
	case *ast.FunctionDecl:
		w.FunctionDeclarationStmt(newNode, scope, Function)
	case *ast.ReturnStmt:
		w.ReturnStmt(newNode, scope)
	case *ast.YieldStmt:
		w.YieldStmt(newNode, scope)
	case *ast.BreakStmt:
		w.BreakStmt(newNode, scope)
	case *ast.ContinueStmt:
		w.ContinueStmt(newNode, scope)
	case *ast.RepeatStmt:
		w.RepeatStmt(newNode, scope)
	case *ast.WhileStmt:
		w.WhileStmt(newNode, scope)
	case *ast.ForStmt:
		w.ForloopStmt(newNode, scope)
	case *ast.TickStmt:
		w.TickStmt(newNode, scope)
	case *ast.CallExpr:
		val := w.GetNodeValue(&newNode.Caller, scope)
		_, finalNode := w.CallExpr(val, newNode, scope)
		*node = finalNode
	case *ast.MethodCallExpr:
		_, *node = w.MethodCallExpr(newNode, scope)
	case *ast.EnvAccessExpr:
		_, newVersion := w.EnvAccessExpr(newNode)
		if newVersion != nil {
			*node = newVersion
		}
	case *ast.ClassDecl:
		w.ClassDeclarationStmt(newNode, scope)
	case *ast.EnumDecl:
		w.EnumDeclarationStmt(newNode, scope)
	case *ast.MatchStmt:
		w.MatchStmt(newNode, false, scope)
	case *ast.AssignmentStmt:
		w.AssignmentStmt(newNode, scope)
	case *ast.UseStmt:
		w.UseStmt(newNode, scope)
	case *ast.DestroyStmt:
		w.DestroyStmt(newNode, scope)
	case *ast.SpawnExpr:
		w.SpawnExpr(newNode, scope)
	case *ast.NewExpr:
		w.NewExpr(newNode, scope)
	case *ast.AliasDecl:
		w.AliasDeclarationStmt(newNode, scope)
	// case *ast.TypeDeclarationStmt:
	// 	TypeDeclarationStmt(newNode, scope)
	case *ast.Improper:
		// w.Error(newNode.GetToken(), "Improper statement: parser fault")
	case *ast.MacroDecl:
	case *ast.EntityDecl:
		w.EntityDeclarationStmt(newNode, scope)
	default:
		// w.Error(newNode.GetToken(), "Expected statement")
	}
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
		localVal, finalNode := w.CallExpr(callVal, newNode, scope)
		val = localVal
		*node = finalNode
	case *ast.MethodCallExpr:
		val, *node = w.MethodCallExpr(newNode, scope)
	case *ast.MapExpr:
		val = w.MapExpr(newNode, scope)
	case *ast.FunctionExpr:
		val = w.FunctionExpr(newNode, scope)
	case *ast.StructExpr:
		val = w.StructExpr(newNode, scope)
	case *ast.MemberExpr:
		val = w.MemberExpr(newNode, scope)
	case *ast.FieldExpr:
		val = w.FieldExpr(newNode, scope)
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
	// case *ast.CastExpr:
	// 	val = CastExpr(newNode, scope)
	case *ast.UseStmt:
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
			// w.Warn((*body)[i].GetToken(), "unreachable code detected")
			endIndex = i
			break
		}
		w.WalkNode(&(*body)[i], scope)
	}
	if endIndex != -1 {
		*body = (*body)[:endIndex]
	}
}

// func WalkBody(w *Walker, body *[]ast.Node, tag ExitableTag, scope *Scope) {
// 	endIndex := -1
// 	for i := range *body {
// 		if tag.GetIfExits(All) {
// 			w.Warn((*body)[i].GetToken(), "unreachable code detected")
// 			endIndex = i
// 			break
// 		}
// 		WalkNode(&(*body)[i], scope)
// 	}
// 	if endIndex != -1 {
// 		*body = (*body)[:endIndex]
// 	}
// }

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
		params = append(params, w.TypeExpr(param.Type, scope, false))
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
