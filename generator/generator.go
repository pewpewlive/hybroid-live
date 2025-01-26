package generator

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/helpers"
	"math"
	"strconv"
)

const charset = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const hyGTL = "GL"
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

type ReplaceType int

const (
	YieldReplacement ReplaceType = iota
	GotoReplacement
	ContinueReplacement
	VariadicParamReplacement
	VariableReplacement
)

type Replacement struct {
	Tag  ReplaceType
	Span helpers.Span[int]
}

type ReplaceSettings map[ReplaceType]string

type GenScope struct { // 0 3
	StringBuilder

	Parent          *GenScope
	Replacements    []Replacement
	ReplaceSettings ReplaceSettings
}

func (gs *GenScope) AddReplacement(tag ReplaceType, span helpers.Span[int]) {
	gs.Replacements = append(gs.Replacements, Replacement{Tag: tag, Span: span})
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func ResolveReplacement(rType ReplaceType, scope *GenScope) string {
	if r, ok := scope.ReplaceSettings[rType]; ok {
		return r
	}

	if scope.Parent == nil {
		return "SKILL ISSUE"
	}

	return ResolveReplacement(rType, scope.Parent)
}

func (gs *GenScope) ReplaceAll() {
	lengthBefore := gs.Len()

	for i := len(gs.Replacements) - 1; i >= 0; i-- {
		replace := ResolveReplacement(gs.Replacements[i].Tag, gs)

		gs.ReplaceSpan(replace, gs.Replacements[i].Span)
		length := gs.Len() - lengthBefore

		for j := i + 1; j < len(gs.Replacements); j++ {
			gs.Replacements[j].Span.Start += length
			gs.Replacements[j].Span.End += length
		}

		lengthBefore = gs.Len()

		RemoveIndex(gs.Replacements, i)
	}
}

func NewGenScope(scope *GenScope) GenScope {
	new := GenScope{
		Parent:          scope,
		StringBuilder:   StringBuilder{},
		Replacements:    []Replacement{},
		ReplaceSettings: map[ReplaceType]string{},
	}

	return new
}

type Generator struct {
	alerts.AlertHandler // ideally should not be ever triggered here, but if triggered something has gone really wrong

	envName     string
	envType     ast.EnvType
	Scope       GenScope
	Future      string
	Errors      []ast.Error
	libraryVars *map[string]string
}

func (gen *Generator) SetUniqueEnvName(name string) {
	uniqueName := StringBuilder{}
	uniqueName.WriteByte('E')
	ResolveVarCounter(&uniqueName, envCounter)
	envCounter++

	envMap[name] = uniqueName.String()
}

func (gen *Generator) SetEnv(name string, typ ast.EnvType) {
	gen.envName = envMap[name]
	gen.envType = typ
}

func (gen *Generator) WriteVar(name string) string {
	return gen.envName + name
}

func (gen *Generator) WriteVarExtra(name, middle string) string {
	return gen.envName + middle + name
}

func (gen *Generator) Clear() {
	gen.Scope = NewGenScope(nil)
	gen.Errors = make([]ast.Error, 0)
}

func getTabs() string {
	tabs := StringBuilder{}
	for i := 0; i < TabsCount; i++ {
		tabs.Write("\t")
	}

	return tabs.String()
}

func (gen Generator) GetErrors() []ast.Error {
	return gen.Errors
}

func (gen *Generator) GetSrc() string {
	return gen.Scope.String()
}

func (gen *Generator) Generate(program []ast.Node, builtins []string) {
	for i := range builtins {
		gen.Scope.WriteString(functions[builtins[i]])
	}
	for _, node := range program {
		gen.GenerateStmt(node, &gen.Scope)
		gen.Scope.WriteString("\n")
	}
}

func (gen *Generator) GenerateWithBuiltins(program []ast.Node) {
	gen.Scope.WriteString(ParseSoundFunction)
	gen.Scope.WriteString(ToStringFunction)
	for _, node := range program {
		gen.GenerateStmt(node, &gen.Scope)
		gen.Scope.WriteString("\n")
	}
}

func (gen *Generator) GenerateBody(program []ast.Node, scope *GenScope) {
	TabsCount += 1
	if gen.Future != "" {
		scope.WriteTabbed(gen.Future)
		gen.Future = ""
	}
	for _, node := range program {
		gen.GenerateStmt(node, scope)
		scope.Write("\n")
	}
	TabsCount -= 1
}

func fixedToFx(floatstr string) string {
	float, _ := strconv.ParseFloat(floatstr, 64)
	abs_float := math.Abs(float)
	integer := math.Floor(abs_float)
	if integer > (2 << 51) {
		integer = (2 << 51)
	}
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

func (gen *Generator) GenerateStmt(node ast.Node, scope *GenScope) {
	switch newNode := node.(type) {
	case *ast.EnvironmentDecl:
		gen.envStmt(*newNode, scope)
	case *ast.AssignmentStmt:
		gen.assignmentStmt(*newNode, scope)
	case *ast.BreakStmt:
		gen.breakStmt(*newNode, scope)
	case *ast.ReturnStmt:
		gen.returnStmt(*newNode, scope)
	case *ast.YieldStmt:
		gen.yieldStmt(*newNode, scope)
	case *ast.ContinueStmt:
		gen.continueStmt(*newNode, scope)
	case *ast.MatchStmt:
		gen.matchStmt(*newNode, scope)
	case *ast.IfStmt:
		gen.ifStmt(*newNode, scope)
	case *ast.RepeatStmt:
		gen.repeatStmt(*newNode, scope)
	case *ast.WhileStmt:
		gen.whileStmt(*newNode, scope)
	case *ast.ForStmt:
		gen.forStmt(*newNode, scope)
	case *ast.TickStmt:
		gen.tickStmt(*newNode, scope)
	case *ast.VariableDecl:
		gen.variableDeclarationStmt(*newNode, scope)
	case *ast.CallExpr:
		val := gen.callExpr(*newNode, true, scope)
		scope.WriteString(val)
	case *ast.MethodCallExpr:
		val := gen.methodCallExpr(*newNode, true, scope)
		scope.WriteString(val)
	case *ast.SpawnExpr:
		val := gen.spawnExpr(*newNode, true, scope)
		scope.WriteString(val)
	case *ast.NewExpr:
		val := gen.newExpr(*newNode, true, scope)
		scope.WriteString(val)
	case *ast.FunctionDecl:
		gen.functionDeclarationStmt(*newNode, scope)
	case *ast.EnumDecl:
		gen.enumDeclarationStmt(*newNode, scope)
	case *ast.ClassDecl:
		gen.classDeclarationStmt(*newNode, scope)
	case *ast.EnvAccessExpr:
		val := gen.envAccessExpr(*newNode, scope)
		scope.WriteString(val)
	case *ast.EntityDecl:
		gen.entityDeclarationStmt(*newNode, scope)
	case *ast.DestroyStmt:
		gen.destroyStmt(*newNode, scope)
	}
}

func (gen *Generator) GenerateExpr(node ast.Node, scope *GenScope) string {
	switch newNode := node.(type) {
	case *ast.LiteralExpr:
		return gen.literalExpr(*newNode)
	case *ast.EntityExpr:
		return gen.entityExpr(*newNode, scope)
	case *ast.BinaryExpr:
		return gen.binaryExpr(*newNode, scope)
	case *ast.IdentifierExpr:
		return gen.identifierExpr(*newNode, scope)
	case *ast.GroupExpr:
		return gen.groupingExpr(*newNode, scope)
	case *ast.ListExpr:
		return gen.listExpr(*newNode, scope)
	case *ast.UnaryExpr:
		return gen.unaryExpr(*newNode, scope)
	case *ast.CallExpr:
		return gen.callExpr(*newNode, false, scope)
	case *ast.MapExpr:
		return gen.mapExpr(*newNode, scope)
	case *ast.FieldExpr:
		return gen.fieldExpr(*newNode, scope)
	case *ast.MemberExpr:
		return gen.memberExpr(*newNode, scope)
	case *ast.FunctionExpr:
		return gen.functionExpr(*newNode, scope)
	case *ast.StructExpr:
		return gen.structExpr(*newNode, scope)
	case *ast.SelfExpr:
		return gen.selfExpr(*newNode, scope)
	case *ast.NewExpr:
		return gen.newExpr(*newNode, false, scope)
	case *ast.MatchExpr:
		return gen.matchExpr(*newNode, scope)
	case *ast.EnvAccessExpr:
		return gen.envAccessExpr(*newNode, scope)
	case *ast.SpawnExpr:
		return gen.spawnExpr(*newNode, false, scope)
	case *ast.MethodCallExpr:
		return gen.methodCallExpr(*newNode, false, scope)
		// case *ast.CastExpr:
		// 	return gen.castExpr(*newNode, scope)
	}

	return ""
}
