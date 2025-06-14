package walker // THE ACTUAL WALKING

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
	"slices"
)

type Environment struct {
	Name        string
	luaPath     string // dynamic lua path
	hybroidPath string
	Type        ast.Env

	Scope Scope

	importedWalkers []*Walker // the walkers imported through UseStmt
	UsedLibraries   []Library
	UsedBuiltinVars []string

	Classes  map[string]*ClassVal
	Entities map[string]*EntityVal
	Enums    map[string]*EnumVal

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
		Tag:         &UntaggedTag{},
		Variables:   map[string]*VariableVal{},
		AliasTypes:  make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	}
	global := &Environment{
		hybroidPath:   hybroidPath,
		luaPath:       luaPath,
		Type:          ast.InvalidEnv,
		Scope:         scope,
		UsedLibraries: make([]Library, 0),
		Classes:       map[string]*ClassVal{},
		Entities:      map[string]*EntityVal{},
		Enums:         map[string]*EnumVal{},
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

	walkers      map[string]*Walker
	program      []ast.Node
	context      Context
	Walked       bool
	ignoreAlerts bool
}

func (w *Walker) Alert(alertType alerts.Alert, args ...any) {
	w.Alert_(alertType, args...)
}

func (w *Walker) AlertI(alert alerts.Alert) {
	w.AlertI_(alert)
}

func (w *Walker) AlertSingle(alert alerts.Alert, token tokens.Token, args ...any) {
	if w.ignoreAlerts {
		return
	}
	args = append([]any{alerts.NewSingle(token)}, args...)
	w.Alert(alert, args...)
}

func (w *Walker) AlertMulti(alert alerts.Alert, start, end tokens.Token, args ...any) {
	if w.ignoreAlerts {
		return
	}
	args = append([]any{alerts.NewMulti(start, end)}, args...)
	w.Alert(alert, args...)
}

func NewWalker(hybroidPath, luaPath string) *Walker {
	return &Walker{
		environment: NewEnvironment(hybroidPath, luaPath),
		program:     []ast.Node{},
		context: Context{
			EntityCasts: core.NewQueue[EntityCast]("EntityCasts"),
		},
		Collector: alerts.NewCollector(),
	}
}

func (w *Walker) Env() Environment {
	return *w.environment
}

func (w *Walker) SetProgram(program []ast.Node) {
	w.program = program
}

func (w *Walker) Program() []ast.Node {
	return w.program
}

func (w *Walker) PreWalk(walkers map[string]*Walker) {
	if w.walkers == nil && walkers != nil {
		w.walkers = walkers
	}

	if len(w.program) == 0 {
		return
	}

	if w.program[0].GetType() != ast.EnvironmentDeclaration {
		w.AlertSingle(&alerts.ExpectedEnvironment{}, w.program[0].GetToken())
		return
	}

	for _, node := range w.program {
		if environmentDecl, ok := node.(*ast.EnvironmentDecl); ok {
			w.environmentDeclaration(environmentDecl)
		}
	}
}

func (w *Walker) Walk() {
	if len(w.program) == 0 {
		return
	}

	if w.program[0].GetType() != ast.EnvironmentDeclaration {
		return
	}

	scope := &w.environment.Scope
	for i := range w.program {
		w.walkNode(&w.program[i], scope)
	}
	w.CheckUniqueVariables()

	w.Walked = true
}

func (w *Walker) PostWalk() {
	for _, v := range w.environment.Scope.Variables {
		if !v.IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, v.Token, "variable")
		}
	}
	for _, v := range w.environment.Entities {
		if !v.Type.IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, v.Token, "entity type")
		}
	}
	for _, v := range w.environment.Classes {
		if !v.Type.IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, v.Token, "class type")
		}
	}
	for _, v := range w.environment.Enums {
		if !v.Type.IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, v.Token, "enum type")
		}
	}
	for _, v := range w.environment.Scope.AliasTypes {
		if !v.IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, v.Token, "alias type")
		}
	}
}

func (w *Walker) CheckUniqueVariables() {
	if w.environment.Type == ast.MeshEnv {
		variable, ok := w.environment.Scope.Variables["meshes"]
		if !ok {
			w.AlertSingle(&alerts.MissingPewpewVariable{}, w.environment._envStmt.GetToken(), "meshes", "Mesh")
			return
		}
		variable.IsUsed = true
		if !variable.IsPub || !TypeEquals(variable.Value.GetType(), MeshesType) {
			w.AlertSingle(&alerts.InvalidPewpewVariable{}, variable.Token, "meshes", MeshType)
		}
		return
	}
	if w.environment.Type == ast.SoundEnv {
		variable, ok := w.environment.Scope.Variables["sounds"]
		if !ok {
			w.AlertSingle(&alerts.MissingPewpewVariable{}, w.environment._envStmt.GetToken(), "sounds", "Sound")
			return
		}
		variable.IsUsed = true
		if !variable.IsPub || !TypeEquals(variable.Value.GetType(), SoundsType) {
			w.AlertSingle(&alerts.InvalidPewpewVariable{}, variable.Token, "sounds", SoundType)
		}
	}
}

func (w *Walker) walkNode(node *ast.Node, scope *Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentDecl:
	case *ast.VariableDecl:
		w.variableDeclaration(newNode, scope, false)
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
		w.callExpression(val, node, scope)
	case *ast.ClassDecl:
		w.classDeclaration(newNode, scope)
	case *ast.EnumDecl:
		w.enumDeclaration(newNode, scope)
	case *ast.MatchStmt:
		w.matchStatement(newNode, scope)
	case *ast.AssignmentStmt:
		w.assignmentStatement(newNode, scope)
	case *ast.UseStmt:
		w.useStatement(newNode, scope)
	case *ast.DestroyStmt:
		w.destroyStatement(newNode, scope)
	case *ast.SpawnExpr:
		w.spawnExpression(newNode, scope)
	case *ast.NewExpr:
		w.newExpression(newNode, scope)
	case *ast.AliasDecl:
		w.aliasDeclaration(newNode, scope)
	case *ast.Improper:
		// w.Error(newNode.GetToken(), "Improper statement: parser fault")
	case *ast.EntityDecl:
		w.entityDeclaration(newNode, scope)
	default:
		// w.Error(newNode.GetToken(), "Expected statement")
	}
}

func (w *Walker) GetActualNodeValue(node *ast.Node, scope *Scope) Value {
	val := w.GetNodeValue(node, scope)
	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}
	if constVal, ok := val.(*ConstVal); ok {
		val = constVal.Val
	}

	return val
}

func (w *Walker) GetNodeValue(node *ast.Node, scope *Scope) Value {
	var val Value

	switch newNode := (*node).(type) {
	case *ast.LiteralExpr:
		val = w.literalExpression(newNode)
	case *ast.BinaryExpr:
		val = w.binaryExpression(newNode, scope)
	case *ast.IdentifierExpr:
		val = w.identifierExpression(node, scope)
	case *ast.GroupExpr:
		val = w.groupExpression(newNode, scope)
	case *ast.ListExpr:
		val = w.listExpression(newNode, scope)
	case *ast.UnaryExpr:
		val = w.unaryExpression(newNode, scope)
	case *ast.CallExpr:
		callVal := w.GetNodeValue(&newNode.Caller, scope)
		localVal := w.callExpression(callVal, node, scope)
		val = localVal
	case *ast.MapExpr:
		val = w.mapExpression(newNode, scope)
	case *ast.FunctionExpr:
		val = w.functionExpression(newNode, scope)
	case *ast.StructExpr:
		val = w.structExpression(newNode, scope)
	case *ast.AccessExpr:
		val = w.accessExpression(node, scope)
	case *ast.NewExpr:
		val = w.newExpression(newNode, scope)
	case *ast.SelfExpr:
		val = w.selfExpression(newNode, scope)
	case *ast.MatchExpr:
		val = w.matchExpression(newNode, scope)
	case *ast.EntityEvaluationExpr:
		val = w.entityEvaluationExpression(newNode, scope)
	case *ast.EnvAccessExpr:
		val = w.environmentAccessExpression(node)
	case *ast.SpawnExpr:
		val = w.spawnExpression(newNode, scope)
	default:
		// w.Error(newNode.GetToken(), "Expected expression")
		return &Invalid{}
	}

	return val
}
func (w *Walker) walkBody(body *ast.Body, tag ExitableTag, scope *Scope) {
	endIndex := -1
	bodySlice := *body
	for i := range bodySlice {
		if tag.GetIfExits(ControlFlow) {
			w.AlertMulti(&alerts.UnreachableCode{}, bodySlice[i].GetToken(), bodySlice[body.Size()-1].GetToken())
			endIndex = i
			break
		}
		w.walkNode(body.Node(i), scope)
	}
	if endIndex != -1 {
		*body = bodySlice[:endIndex]
	}
	for k := range scope.Variables {
		if !scope.Variables[k].IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, scope.Variables[k].Token, "variable")
		}
	}
	for k := range scope.AliasTypes {
		if !scope.AliasTypes[k].IsUsed {
			w.AlertSingle(&alerts.UnusedElement{}, scope.Variables[k].Token, "alias type")
		}
	}
}

func (w *Walker) walkFuncBody(node ast.Node, body *ast.Body, tag *FuncTag, scope *Scope) {
	w.walkBody(body, tag, scope)

	if !tag.GetIfExits(Return) && len(tag.ReturnTypes) != 0 {
		w.AlertSingle(&alerts.NotAllCodePathsExit{}, node.GetToken(), "return")
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
