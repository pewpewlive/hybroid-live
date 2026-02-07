package lsp

import (
	"hybroid/walker"
)

var keywordDocs = map[string]string{
// ... (rest of the map remains the same)
	"is":       "Checks if a value is of a certain entity type.",
	"isnt":     "Checks if a value is NOT of a certain entity type.",
	"alias":    "Creates a new name for an existing type.",
	"and":      "Logical AND operator.",
	"as":       "Used in environment declarations or type casting.",
	"break":    "Exits the innermost loop or match case.",
	"by":       "Used in range-based for loops to specify the step.",
	"const":    "Declares a constant value that cannot be reassigned.",
	"continue": "Skips to the next iteration of the innermost loop.",
	"else":     "Executes when the 'if' condition is false.",
	"entity":   "Defines a new game entity type or refers to the generic entity type.",
	"enum":     "Defines a set of named constants.",
	"env":      "Declares the environment (Level, Mesh, Sound, Shared) for the current file.",
	"false":    "Boolean false value.",
	"fn":       "Defines a function or function type.",
	"to":       "Specifies the end of a range in a for loop.",
	"for":      "Starts a loop over a collection or range.",
	"if":       "Starts a conditional block.",
	"in":       "Used in for loops to specify the collection.",
	"let":      "Declares a mutable variable.",
	"match":    "Starts a pattern-matching block or expression.",
	"new":      "Instantiates a new class instance.",
	"or":       "Logical OR operator.",
	"pub":      "Makes a declaration accessible from other files.",
	"repeat":   "Starts a loop that repeats a specific number of times.",
	"return":   "Exits a function and optionally returns values.",
	"self":     "Refers to the current class or entity instance.",
	"spawn":    "Creates a new instance of an entity.",
	"struct":   "Defines a collection of named fields.",
	"class":    "Defines a new class with fields and methods.",
	"tick":     "Starts a block that executes every game tick.",
	"true":     "Boolean true value.",
	"use":      "Imports another environment or library.",
	"from":     "Specifies the start of a range in a for loop.",
	"while":    "Starts a loop that continues while a condition is true.",
	"with":     "Used in certain expressions to provide additional context.",
	"yield":    "Returns a value from a match expression.",
	"destroy":  "Removes an entity from the game.",
	"every":    "Specifies a frequency for tick-based logic.",
}

var typeDocs = map[string]string{
	"number": "A 64-bit floating-point number.",
	"fixed":  "A 32-bit fixed-point number used for deterministic calculations.",
	"text":   "A UTF-8 encoded string of characters.",
	"bool":   "A boolean value: either 'true' or 'false'.",
	"list":   "A dynamic array-like collection of elements.",
	"map":    "A collection of key-value pairs.",
	"struct": "A user-defined collection of named fields.",
	"entity": "A reference to a game entity.",
}

var namespaceDocs = map[string]string{
	"Pewpew": "The main API for interacting with the PewPew Live game engine. Provides functions for entities, graphics, and game state.",
	"Fmath":  "Fixed-point mathematics library. Essential for deterministic gameplay logic across different devices.",
	"Math":   "Standard floating-point mathematics library.",
	"String": "Utilities for string manipulation and formatting.",
	"Table":  "Utilities for manipulating lists and maps.",
}

var environmentDocs = map[string]string{
	"Level":  "Game level environment. Access to `Pewpew` and `Fmath` libraries. Mandatory for level scripts.",
	"Mesh":   "Mesh generation environment. Used for creating procedurally generated 3D models.",
	"Sound":  "Sound generation environment. Used for creating procedurally generated sound effects.",
	"Shared": "Shared environment. Contains code that can be used by Level, Mesh, or Sound scripts.",
}

func getSymbolMetadata(label string) (detail string, doc string) {
	if d, ok := environmentDocs[label]; ok {
		return "Environment", d
	}
	if d, ok := namespaceDocs[label]; ok {
		return "Namespace", d
	}
	if d, ok := typeDocs[label]; ok {
		return "Native Type", d
	}
	if d, ok := keywordDocs[label]; ok {
		return "Keyword", d
	}

	// Check Builtin
	if v, ok := walker.BuiltinEnv.Scope.Variables[label]; ok {
		return "Builtin", v.Value.GetType().String()
	}

	// Check Pewpew
	if v, ok := walker.PewpewAPI.Scope.Variables[label]; ok {
		return "Pewpew", v.Value.GetType().String()
	}

	// Check Fmath
	if v, ok := walker.FmathAPI.Scope.Variables[label]; ok {
		return "Fmath", v.Value.GetType().String()
	}

	return "", ""
}
