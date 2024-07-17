package api

import "hybroid/ast"

var PewPew struct {
	Methods      map[string]string
	MethodParams map[string][]ast.Param
	WeaponType *VariableVal
}