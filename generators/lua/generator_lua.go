package lua

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"math"
	"strconv"
	"strings"
)

// func (ge *GenError) generatorError() string {
// 	return fmt.Sprintf("Error: %v, at line: %v (%v)", ge.Message, ge.Token.Location.LineStart, ge.Token.ToString())
// }

func (gen *Generator) error(token lexer.Token, message string) {
	gen.Errors = append(gen.Errors, ast.Error{Token: token, Message: message})
}

type StringBuilder struct {
	strings.Builder
}

func (sb *StringBuilder) Append(chunks ...string) {
	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
}

type Generator struct {
	Errors    []ast.Error
	Src       StringBuilder
	TabsCount int
}

func (gen *Generator) getTabs() string {
	tabs := StringBuilder{}
	for i := 0; i < gen.TabsCount; i++ {
		tabs.Append("\t")
	}

	return tabs.String()
}

func (gen Generator) GetErrors() []ast.Error {
	return gen.Errors
}

func (gen *Generator) GetSrc() string {
	return gen.Src.String()
}

func (gen *Generator) Generate(program []ast.Node) {
	generatedStr := ""
	for _, node := range program {
		generatedStr = gen.GenerateNode(node)
		gen.Src.Append(generatedStr, "\n")
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

func (gen *Generator) GenerateNode(node ast.Node) string {
	switch newNode := node.(type) {
	case ast.VariableDeclarationStmt:
		return gen.variableDeclarationStmt(newNode)
	case ast.IfStmt:
		return gen.ifStmt(newNode)
	case ast.AssignmentStmt:
		return gen.assignmentStmt(newNode)
	case ast.FunctionDeclarationStmt:
		return gen.functionDeclarationStmt(newNode)
	case ast.ReturnStmt:
		return gen.returnStmt(newNode)
	case ast.RepeatStmt:
		return gen.repeatStmt(newNode)
	case ast.TickStmt:
		return gen.tickStmt(newNode)
	case ast.LiteralExpr:
		return gen.literalExpr(newNode)
	case ast.BinaryExpr:
		return gen.binaryExpr(newNode)
	case ast.IdentifierExpr:
		return gen.identifierExpr(newNode)
	case ast.GroupExpr:
		return gen.groupingExpr(newNode)
	case ast.ListExpr:
		return gen.listExpr(newNode)
	case ast.UnaryExpr:
		return gen.unaryExpr(newNode)
	case ast.CallExpr:
		return gen.callExpr(newNode)
	case ast.MapExpr:
		return gen.mapExpr(newNode)
	case ast.MemberExpr:
		return gen.memberExpr(newNode)
	case ast.DirectiveExpr:
		return gen.directiveExpr(newNode)
	}

	return ""
}
