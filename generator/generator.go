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

	LatestSrc *core.StringBuilder
}

func (gen *Generator) Twrite(src *core.StringBuilder, chunks ...string) {
	chunks = append([]string{gen.tabString()}, chunks...)
	src.Write(chunks...)
}

func NewGenerator() Generator {
	gen := Generator{
		Collector:      alerts.NewCollector(),
		ContinueLabels: core.NewStack[string]("ContinueLabels"),
		BreakLabels:    core.NewStack[string]("BreakLabels"),
		YieldContexts:  core.NewStack[YieldContext]("YieldContext"),
	}
	gen.LatestSrc = &gen.src
	return gen
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

func (gen *Generator) GenerateUsedLibaries(usedLibraries []ast.Library) {
	for _, v := range usedLibraries {
		str := v.String()
		gen.src.Write("local ", str, " = ", str, "\n")
	}
}

func (gen *Generator) Generate(program []ast.Node, builtins []string) {
	for i := range builtins {
		gen.src.Write(mapping.Functions[builtins[i]])
		gen.src.Write("\n")
	}
	for _, node := range program {
		gen.src.Write(gen.GenerateStmt(node), "\n")
	}
}

func (gen *Generator) GenerateWithBuiltins(program []ast.Node) {
	gen.src.Write(mapping.ParseSoundFunction, "\n", mapping.ToStringFunction, "\n")
	for _, node := range program {
		gen.src.Write(gen.GenerateStmt(node), "\n")
	}
}

func (gen *Generator) GenerateBody(src *core.StringBuilder, body ast.Body) {
	gen.tabCount++
	prevSrc := gen.LatestSrc
	gen.LatestSrc = src
	for _, node := range body {
		src.Write(gen.GenerateStmt(node), "\n")
	}
	gen.LatestSrc = prevSrc
	gen.tabCount--
}

func (gen *Generator) GenerateStmt(node ast.Node) string {
	stmt := ""
	switch newNode := node.(type) {
	case *ast.EnvironmentDecl:
		stmt = gen.envStmt(*newNode)
	case *ast.AssignmentStmt:
		assignStmts := gen.breakDownAssignStmt(*newNode)
		src := core.StringBuilder{}
		for i, assignStmt := range assignStmts {
			if i != 0 {
				src.Write("\n")
			}
			src.Write(gen.assignmentStmt(assignStmt))
		}
		return src.String()
	case *ast.BreakStmt:
		stmt = gen.breakStmt(*newNode)
	case *ast.ReturnStmt:
		stmt = gen.returnStmt(*newNode)
	case *ast.YieldStmt:
		stmt = gen.yieldStmt(*newNode)
	case *ast.ContinueStmt:
		stmt = gen.continueStmt(*newNode)
	case *ast.MatchStmt:
		stmt = gen.matchStmt(*newNode)
	case *ast.IfStmt:
		stmt = gen.ifStmt(*newNode)
	case *ast.RepeatStmt:
		stmt = gen.repeatStmt(*newNode)
	case *ast.WhileStmt:
		stmt = gen.whileStmt(*newNode)
	case *ast.ForStmt:
		stmt = gen.forStmt(*newNode)
	case *ast.TickStmt:
		stmt = gen.tickStmt(*newNode)
	case *ast.VariableDecl:
		src := core.StringBuilder{}
		varDecls := gen.breakDownVariableDeclaration(*newNode)
		for i, varDecl := range varDecls {
			if i != 0 {
				src.Write("\n")
			}
			src.Write(gen.variableDeclaration(varDecl))
		}
		return src.String()
	case *ast.CallExpr:
		stmt = gen.callExpr(*newNode, true)
	case *ast.MethodCallExpr:
		stmt = gen.methodCallExpr(*newNode, true)
	case *ast.SpawnExpr:
		stmt = gen.spawnExpr(*newNode, true)
	case *ast.NewExpr:
		stmt = gen.newExpr(*newNode, true)
	case *ast.FunctionDecl:
		stmt = gen.functionDeclaration(*newNode)
	case *ast.ClassDecl:
		stmt = gen.classDeclaration(*newNode)
	case *ast.EnvAccessExpr:
		stmt = gen.envAccessExpr(*newNode)
	case *ast.EntityDecl:
		stmt = gen.entityDeclaration(*newNode)
	case *ast.DestroyStmt:
		stmt = gen.destroyStmt(*newNode)
	default:
		return ""
	}
	return stmt
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
	case *ast.MethodExpr:
		return gen.methodExpr(*newNode)
	case *ast.EntityAccessExpr:
		return gen.entityAccessExpr(*newNode)
	case *ast.FieldExpr:
		return gen.fieldExpr(*newNode)
	case *ast.MemberExpr:
		return gen.memberExpr(*newNode)
	}

	return ""
}
