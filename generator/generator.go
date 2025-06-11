package generator

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/generator/mapping"
	"strings"
)

const charset = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const hyGotoLabel = "GL"
const hyVar = "H"
const hyClass = "HC"
const hyEntity = "HE"

var envMap = map[string]string{}
var varCounter = 0
var envCounter = 0

func ResetGlobalGeneratorValues() {
	envMap = map[string]string{}
	varCounter = 0
	envCounter = 0
}

func ResolveVarCounter(varname *core.StringBuilder, counter int) {
	if counter > len(charset)-1 {
		newCounter := counter - len(charset)
		varname.WriteByte(charset[len(charset)-1])
		ResolveVarCounter(varname, newCounter)
	} else {
		varname.WriteByte(charset[counter])
	}
}

func GenerateVar(prefix string) string {
	varName := core.StringBuilder{}
	varName.Write(prefix)
	ResolveVarCounter(&varName, varCounter)
	varCounter++
	return varName.String()
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

type YieldContext struct {
	label string
	vars  []string
}

func NewYieldContext(vars []string, label string) YieldContext {
	return YieldContext{vars: vars, label: label}
}

type Generator struct {
	alerts.Collector

	env      ast.Env
	src      core.StringBuilder
	tabCount int

	envName       string
	envPrefixName string

	ContinueLabels core.Stack[string]
	BreakLabels    core.Stack[string]
	YieldContexts  core.Stack[YieldContext]

	buffer        core.StringBuilder
	writeToBuffer bool
}

func (gen *Generator) Twrite(chunks ...string) {
	chunks = append([]string{gen.tabString()}, chunks...)
	if gen.writeToBuffer {
		gen.buffer.Write(chunks...)
	} else {
		gen.src.Write(chunks...)
	}
}

func (gen *Generator) Write(chunks ...string) {
	if gen.writeToBuffer {
		gen.buffer.Write(chunks...)
	} else {
		gen.src.Write(chunks...)
	}
}

func NewGenerator() Generator {
	return Generator{
		Collector:      alerts.NewCollector(),
		ContinueLabels: core.NewStack[string]("ContinueLabels"),
		BreakLabels:    core.NewStack[string]("BreakLabels"),
		YieldContexts:  core.NewStack[YieldContext]("YieldContext"),
	}
}

func (gen *Generator) tabString() string {
	return strings.Repeat("\t", gen.tabCount)
}

func (gen *Generator) SetUniqueEnvName(name string) {
	uniqueName := core.StringBuilder{}
	uniqueName.WriteByte('E')
	ResolveVarCounter(&uniqueName, envCounter)
	envCounter++

	envMap[name] = uniqueName.String()
}

func (gen *Generator) SetEnv(name string, env ast.Env) {
	gen.envName = envMap[name]
	gen.env = env
	gen.envPrefixName = gen.envName
}

func (gen *Generator) WriteVar(name string) string {
	defer func() {
		if gen.envPrefixName != gen.envName {
			gen.envPrefixName = gen.envName
		}
	}()
	return gen.envPrefixName + name
}

func (gen *Generator) WriteVarExtra(name, extra string) string {
	defer func() {
		if gen.envPrefixName != gen.envName {
			gen.envPrefixName = gen.envName
		}
	}()
	return extra + gen.envPrefixName + name
}

func (gen *Generator) GetSrc() string {
	return gen.src.String()
}

func (gen *Generator) Generate(program []ast.Node, builtins []string) {
	for i := range builtins {
		gen.Write(mapping.Functions[builtins[i]])
		gen.Write("\n")
	}
	for _, node := range program {
		gen.GenerateStmt(node)
	}
}

func (gen *Generator) GenerateWithBuiltins(program []ast.Node) {
	gen.Write(mapping.ParseSoundFunction, "\n", mapping.ToStringFunction, "\n")
	for _, node := range program {
		gen.GenerateStmt(node)
	}
}

func (gen *Generator) GenerateBody(body ast.Body) {
	gen.tabCount++
	for _, node := range body {
		gen.GenerateStmt(node)
	}
	gen.tabCount--
}

func (gen *Generator) GenerateStmt(node ast.Node) {
	switch newNode := node.(type) {
	case *ast.EnvironmentDecl:
		gen.envStmt(*newNode)
	case *ast.AssignmentStmt:
		assignStmts := gen.breakDownAssignStmt(*newNode)
		for _, assignStmt := range assignStmts {
			gen.assignmentStmt(assignStmt)
			gen.Write("\n")
		}
		return
	case *ast.BreakStmt:
		gen.breakStmt(*newNode)
	case *ast.ReturnStmt:
		gen.returnStmt(*newNode)
	case *ast.YieldStmt:
		gen.yieldStmt(*newNode)
	case *ast.ContinueStmt:
		gen.continueStmt(*newNode)
	case *ast.MatchStmt:
		gen.matchStmt(*newNode)
	case *ast.IfStmt:
		gen.ifStmt(*newNode)
	case *ast.RepeatStmt:
		gen.repeatStmt(*newNode)
	case *ast.WhileStmt:
		gen.whileStmt(*newNode)
	case *ast.ForStmt:
		gen.forStmt(*newNode)
	case *ast.TickStmt:
		gen.tickStmt(*newNode)
	case *ast.VariableDecl:
		varDecls := gen.breakDownVariableDeclaration(*newNode)
		for _, varDecl := range varDecls {
			gen.variableDeclaration(varDecl)
			gen.Write("\n")
		}
		return
	case *ast.CallExpr:
		val := gen.callExpr(*newNode, true)
		gen.Write(val)
	case *ast.MethodCallExpr:
		val := gen.methodCallExpr(*newNode, true)
		gen.Write(val)
	case *ast.SpawnExpr:
		val := gen.spawnExpr(*newNode, true)
		gen.Write(val)
	case *ast.NewExpr:
		val := gen.newExpr(*newNode, true)
		gen.Write(val)
	case *ast.FunctionDecl:
		gen.functionDeclaration(*newNode)
	case *ast.ClassDecl:
		gen.classDeclaration(*newNode)
	case *ast.EnvAccessExpr:
		val := gen.envAccessExpr(*newNode)
		gen.Write(val)
	case *ast.EntityDecl:
		gen.entityDeclaration(*newNode)
	case *ast.DestroyStmt:
		gen.destroyStmt(*newNode)
	default:
		return
	}
	gen.Write("\n")
}

func (gen *Generator) GenerateExpr(node ast.Node) string {
	switch newNode := node.(type) {
	case *ast.LiteralExpr:
		return gen.literalExpr(*newNode)
	case *ast.EntityEvaluationExpr:
		return gen.entityExpr(*newNode)
	case *ast.BinaryExpr:
		return gen.binaryExpr(*newNode)
	case *ast.IdentifierExpr:
		return gen.identifierExpr(*newNode)
	case *ast.GroupExpr:
		return gen.groupingExpr(*newNode)
	case *ast.ListExpr:
		return gen.listExpr(*newNode)
	case *ast.UnaryExpr:
		return gen.unaryExpr(*newNode)
	case *ast.CallExpr:
		return gen.callExpr(*newNode, false)
	case *ast.MapExpr:
		return gen.mapExpr(*newNode)
	case *ast.AccessExpr:
		return gen.accessExpr(*newNode)
	case *ast.FunctionExpr:
		return gen.functionExpr(*newNode)
	case *ast.StructExpr:
		return gen.structExpr(*newNode)
	case *ast.SelfExpr:
		return gen.selfExpr(*newNode)
	case *ast.NewExpr:
		return gen.newExpr(*newNode, false)
	case *ast.MatchExpr:
		return gen.matchExpr(*newNode)
	case *ast.EnvAccessExpr:
		return gen.envAccessExpr(*newNode)
	case *ast.SpawnExpr:
		return gen.spawnExpr(*newNode, false)
	case *ast.MethodCallExpr:
		return gen.methodCallExpr(*newNode, false)
	case *ast.EntityAccessExpr:
		return gen.entityAccessExpr(*newNode)
	case *ast.FieldExpr:
		return gen.fieldExpr(*newNode)
	case *ast.MemberExpr:
		return gen.memberExpr(*newNode)
	}

	return ""
}
