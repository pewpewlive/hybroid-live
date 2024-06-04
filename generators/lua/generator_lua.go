package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/generators"
	"hybroid/lexer"
	"math"
	"strconv"
)

// func (ge *GenError) generatorError() string {
// 	return fmt.Sprintf("Error: %v, at line: %v (%v)", ge.Message, ge.Token.Location.LineStart, ge.Token.ToString())
// }

func (gen *Generator) error(token lexer.Token, message string) {
	gen.Errors = append(gen.Errors, ast.Error{Token: token, Message: message})
}

type StringBuilder struct {
	generators.BetterBuilder
}

func (sb *StringBuilder) Append(chunks ...string) {
	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
}

func (sb *StringBuilder) AppendTabbed(chunks ...string) {
	sb.WriteString(getTabs())

	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
}

var charset = []byte("_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var charsetLength = len(charset)
var varCounter = 0

func ResolveVarCounter(varname *StringBuilder, counter int) {
	if counter > charsetLength-1 {
		newCounter :=  counter-charsetLength
		varname.WriteByte(charset[charsetLength-1])
		ResolveVarCounter(varname, newCounter)
	} else {
		varname.WriteByte(charset[counter])
	}
}

func GenerateVar() string {
	varName := StringBuilder{}
	varName.WriteByte('H')
	ResolveVarCounter(&varName, varCounter)
	varCounter++
	return varName.String()
}

var TabsCount int

type ReplaceType int

const (
	YieldReplacement ReplaceType = iota
	GotoReplacement
	ContinueReplacement
)

func (sb *StringBuilder) ReplaceRange(str string, _range Range) {
	buf := sb.GetBuf()

	endSlice := make([]byte, len(*buf)-_range.End)

	copy(endSlice, (*buf)[_range.End:len(*buf)])
	*buf = (*buf)[:_range.Start]

	sb.WriteString(str)
	sb.Write(endSlice)
}

type Range struct {
	Start int
	End   int
}

type Replacement struct {
	Tag   ReplaceType
	Range Range
}

func NewRange(start int, end int) Range {
	return Range{Start: start, End: end}
}

type ReplaceSettings map[ReplaceType]string

type GenScope struct { // 0 3
	Parent          *GenScope
	Src             StringBuilder
	Replacements    []Replacement
	ReplaceSettings ReplaceSettings
}

func (gs *GenScope) AddReplacement(tag ReplaceType, _range Range) {
	gs.Replacements = append(gs.Replacements, Replacement{Tag: tag, Range: _range})
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

func (gs *GenScope) ReplaceAll(replacement ReplaceSettings) {
	lengthBefore := gs.Src.Len()

	for i := len(gs.Replacements) - 1; i >= 0; i-- {
		replace := ResolveReplacement(gs.Replacements[i].Tag, gs)

		gs.Src.ReplaceRange(replace, gs.Replacements[i].Range)
		length := gs.Src.Len() - lengthBefore

		for j := i + 1; j < len(gs.Replacements); j++ {
			gs.Replacements[j].Range.Start += length
			gs.Replacements[j].Range.End += length
		}

		lengthBefore = gs.Src.Len()

		RemoveIndex(gs.Replacements, i)
	}
}

func NewGenScope(scope *GenScope) GenScope {
	return GenScope{
		Parent:          scope,
		Src:             StringBuilder{},
		Replacements:    []Replacement{},
		ReplaceSettings: map[ReplaceType]string{},
	}
}

func (gs *GenScope) Write(src StringBuilder) {
	gs.Src.WriteString(src.String())
}

func (gs *GenScope) WriteString(src string) {
	gs.Src.WriteString(src)
}

func (gs *GenScope) Append(strs ...string) {
	gs.Src.Append(strs...)
}

func (gs *GenScope) AppendTabbed(strs ...string) {
	gs.Src.AppendTabbed(strs...)
}

type Generator struct {
	Scope  GenScope
	Errors []ast.Error
}

func getTabs() string {
	tabs := StringBuilder{}
	for i := 0; i < TabsCount; i++ {
		tabs.Append("\t")
	}

	return tabs.String()
}

func (gen Generator) GetErrors() []ast.Error {
	return gen.Errors
}

func (gen *Generator) GetSrc() string {
	return gen.Scope.Src.String()
}

func (gen *Generator) Generate(program []ast.Node) {
	for _, node := range program {
		gen.GenerateStmt(node, &gen.Scope)
		gen.Scope.WriteString("\n")
	}
}

func (gen *Generator) GenerateString(program []ast.Node, scope *GenScope) {
	for _, node := range program {
		gen.GenerateStmt(node, scope)
		scope.Src.WriteString("\n")
	}
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
	case ast.AssignmentStmt:
		gen.assignmentStmt(newNode, scope)
	case ast.BreakStmt:
		gen.breakStmt(newNode, scope)
	case ast.ReturnStmt:
		gen.returnStmt(newNode, scope)
	case ast.YieldStmt:
		gen.yieldStmt(newNode, scope)
	case ast.ContinueStmt:
		gen.continueStmt(newNode, scope)
	case ast.MatchStmt:
		gen.matchStmt(newNode, scope)
	case ast.IfStmt:
		gen.ifStmt(newNode, scope)
	case ast.RepeatStmt:
		gen.repeatStmt(newNode, scope)
	case ast.ForStmt:
		gen.forStmt(newNode, scope)
	case ast.TickStmt:
		gen.tickStmt(newNode, scope)
	case ast.VariableDeclarationStmt:
		gen.variableDeclarationStmt(newNode, scope)
	case ast.UseStmt:
		gen.useStmt(newNode, scope)
	case ast.MethodCallExpr:
		val := gen.methodCallExpr(newNode, scope) // koocing
		scope.WriteString(val)
	case ast.CallExpr:
		val := gen.callExpr(newNode, scope) // koocing
		scope.WriteString(val)
	case ast.FunctionDeclarationStmt:
		gen.functionDeclarationStmt(newNode, scope)
	case ast.StructDeclarationStmt:
		gen.structDeclarationStmt(newNode, scope)
	}
}

func (gen *Generator) GenerateExpr(node ast.Node, scope *GenScope) string {
	switch newNode := node.(type) {
	case ast.LiteralExpr:
		return gen.literalExpr(newNode)
	case ast.BinaryExpr:
		return gen.binaryExpr(newNode, scope)
	case ast.IdentifierExpr:
		return gen.identifierExpr(newNode, scope)
	case ast.GroupExpr:
		return gen.groupingExpr(newNode, scope)
	case ast.ListExpr:
		return gen.listExpr(newNode, scope)
	case ast.UnaryExpr:
		return gen.unaryExpr(newNode, scope)
	case ast.CallExpr:
		return gen.callExpr(newNode, scope)
	case ast.MapExpr:
		return gen.mapExpr(newNode, scope)
	case ast.FieldExpr:
		return gen.fieldExpr(newNode, scope)
	case ast.MemberExpr:
		return gen.memberExpr(newNode, scope)
	case ast.DirectiveExpr:
		return gen.directiveExpr(newNode, scope)
	case ast.AnonFnExpr:
		return gen.anonFnExpr(newNode, scope)
	case ast.SelfExpr:
		return gen.selfExpr(newNode, scope)
	case ast.NewExpr:
		return gen.newExpr(newNode, scope)
	case ast.MatchExpr:
		return gen.matchExpr(newNode, scope)
	case ast.MethodCallExpr:
		return gen.methodCallExpr(newNode, scope)
	}

	return ""
}
