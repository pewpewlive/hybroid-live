package generator

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"math"
	"strconv"
)

const charset = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const hyGotoLabel = "GL"
const hyVar = "H"
const hyClass = "HC"
const hyEntity = "HE"

var envMap = map[string]string{}
var varCounter = 0
var envCounter = 0
var TabsCount int

func ResolveVarCounter(varname *StringBuilder, counter int) {
	if counter > len(charset)-1 {
		newCounter := counter - len(charset)
		varname.WriteByte(charset[len(charset)-1])
		ResolveVarCounter(varname, newCounter)
	} else {
		varname.WriteByte(charset[counter])
	}
}

func GenerateVar(prefix string) string {
	varName := StringBuilder{}
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
	StringBuilder

	envName        string
	env            ast.Env
	ContinueLabels core.Stack[string]
	YieldContexts  core.Stack[YieldContext]
	Future         string
}

func NewGenerator() Generator {
	return Generator{
		Collector: alerts.NewCollector(),
	}
}

func (gen *Generator) SetUniqueEnvName(name string) {
	uniqueName := StringBuilder{}
	uniqueName.WriteByte('E')
	ResolveVarCounter(&uniqueName, envCounter)
	envCounter++

	envMap[name] = uniqueName.String()
}

func (gen *Generator) SetEnv(name string, env ast.Env) {
	gen.envName = envMap[name]
	gen.env = env
}

func (gen *Generator) WriteVar(name string) string {
	return gen.envName + name
}

func (gen *Generator) WriteVarExtra(name, middle string) string {
	return gen.envName + middle + name
}

func getTabs() string {
	tabs := StringBuilder{}
	for range TabsCount {
		tabs.Write("\t")
	}

	return tabs.String()
}

func (gen *Generator) GetSrc() string {
	return gen.String()
}

func (gen *Generator) Generate(program []ast.Node, builtins []string) {
	for i := range builtins {
		gen.WriteString(functions[builtins[i]])
	}
	for _, node := range program {
		gen.GenerateStmt(node)
		gen.WriteString("\n")
	}
}

func (gen *Generator) GenerateWithBuiltins(program []ast.Node) {
	gen.WriteString(ParseSoundFunction)
	gen.WriteString(ToStringFunction)
	for _, node := range program {
		gen.GenerateStmt(node)
		gen.WriteString("\n")
	}
}

func (gen *Generator) GenerateBody(body ast.Body) {
	TabsCount += 1
	if gen.Future != "" {
		gen.WriteTabbed(gen.Future)
		gen.Future = ""
	}
	for _, node := range body {
		gen.GenerateStmt(node)
		gen.Write("\n")
	}
	TabsCount -= 1
}

func fixedToFx(floatstr string) string {
	float, _ := strconv.ParseFloat(floatstr, 64)
	abs_float := math.Abs(float)
	integer := min(math.Floor(abs_float), (2 << 51))
	var sign string
	if float < 0 {
		sign = "-"
	} else {
		sign = ""
	}

	frac := math.Floor((abs_float - integer) * 4096)
	frac_str := ""
	if frac != 0 {
		frac_str = "." + fmt.Sprintf("%v", frac)
	}

	// sign + int + frac_str + "fx"
	return fmt.Sprintf("%s%v%s", sign, integer, frac_str)
}

func degToRad(floatstr string) string {
	float, _ := strconv.ParseFloat(floatstr, 64)
	radians := float * math.Pi / 180
	return fixedToFx(fmt.Sprintf("%v", radians))
}

func (gen *Generator) GenerateStmt(node ast.Node) {
	switch newNode := node.(type) {
	case *ast.EnvironmentDecl:
		gen.envStmt(*newNode)
	case *ast.AssignmentStmt:
		gen.assignmentStmt(*newNode)
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
		}
	case *ast.CallExpr:
		val := gen.callExpr(*newNode, true)
		gen.WriteString(val)
	case *ast.MethodCallExpr:
		val := gen.methodCallExpr(*newNode, true)
		gen.WriteString(val)
	case *ast.SpawnExpr:
		val := gen.spawnExpr(*newNode, true)
		gen.WriteString(val)
	case *ast.NewExpr:
		val := gen.newExpr(*newNode, true)
		gen.WriteString(val)
	case *ast.FunctionDecl:
		gen.functionDeclaration(*newNode)
	case *ast.EnumDecl:
		gen.enumDeclaration(*newNode)
	case *ast.ClassDecl:
		gen.classDeclaration(*newNode)
	case *ast.EnvAccessExpr:
		val := gen.envAccessExpr(*newNode)
		gen.WriteString(val)
	case *ast.EntityDecl:
		gen.entityDeclaration(*newNode)
	case *ast.DestroyStmt:
		gen.destroyStmt(*newNode)
	}
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
		// case *ast.CastExpr:
		// 	return gen.castExpr(*newNode)
	}

	return ""
}
