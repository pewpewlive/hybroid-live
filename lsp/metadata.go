package lsp

import (
	"hybroid/walker"
	"strings"
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
	"let":      "Declares a local variable.",
	"match":    "Starts a pattern-matching block or expression.",
	"new":      "Instantiates a new class instance.",
	"or":       "Logical OR operator.",
	"pub":      "Declares a global variable.",
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
	"number": "An integer number.",
	"fixed":  "A fixed-point number.",
	"text":   "A string of characters.",
	"bool":   "A boolean value.",
	"list":   "A dynamic array-like collection of elements.",
	"map":    "A collection of key-value pairs.",
	"struct": "A user-defined collection of named fields.",
	"entity": "A reference to a game entity.",
}

var namespaceDocs = map[string]string{
	"Pewpew": "The main API for working with PewPew Live. Provides functions for entities, graphics, and game state.",
	"Fmath":  "Fixed-point math library.",
	"Math":   "Floating-point math library.",
	"String": "Utilities for string manipulation and formatting.",
	"Table":  "Utilities for manipulating lists and maps.",
}

var environmentDocs = map[string]string{
	"Level":  "Game level environment. Access to `Pewpew` and `Fmath` libraries. Mandatory for level scripts.",
	"Mesh":   "Mesh generation environment. Used for creating procedurally generated 3D models.",
	"Sound":  "Sound generation environment. Used for creating procedurally generated sound effects.",
	"Shared": "Shared environment. Contains code that can be used by Level, Mesh, or Sound scripts.",
}

var builtinDocs = map[string]string{
	"ToString":   "```hybroid\nToString(value) -> string\n```\nConverts any value to a string.",
	"ParseSound": "```hybroid\nParseSound(string jfxrUrl) -> Sound\n```\nAllows you to parse a sound from a [JFXR](https://pewpew.live/jfxr/index.html) URL. Only available in sound environments.",
}

var aliasDocs = map[string]string{
	"Mesh":     "A struct representing a 3D mesh with `vertexes`, `segments`, and optionally `colors`.",
	"Meshes":   "A list of `Mesh` objects.",
	"Vertex":   "A list of 3 numbers representing a point in 3D space.",
	"Vertexes": "A list of `Vertex` objects.",
	"Segment":  "A list of 2 numbers representing the indices of two vertexes forming a segment.",
	"Segments": "A list of `Segment` objects.",
	"Colors":   "A list of colors, where each color is a number.",
	"Center":   "A struct with `x`, `y`, and `z` fields representing the center of an entity.",
	"Sound":    "A struct representing a sound configuration for procedural generation.",
}

func getSymbolMetadata(w *walker.Walker, walkers map[string]*walker.Walker, label string) (detail string, doc string) {
	if d, ok := environmentDocs[label]; ok {
		return "Environment", d
	}
	if w2, ok := walkers[label]; ok && w2.Env() != nil && w2.Env().Name == label {
		return string(w2.Env().Type), ""
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

	// Handle Namespace:Symbol or Namespace.Symbol
	if strings.Contains(label, ":") || strings.Contains(label, ".") {
		parts := strings.FieldsFunc(label, func(r rune) bool { return r == ':' || r == '.' })
		if len(parts) == 2 {
			ns := parts[0]
			sym := parts[1]

			env := resolveBuiltinEnvByName(ns)

			// Check custom namespaces in walkers map
			if env == nil && walkers != nil {
				if w2, ok := walkers[ns]; ok {
					envVal := w2.Env()
					env = envVal
				}
			}

			// If not a namespace, check if it's an entity/enum/class in the current walker
			if env == nil && w != nil {
				if ev, ok := w.Env().Enums[ns]; ok {
					if field, _, found := ev.ContainsField(sym); found {
						return field.Value.GetType().String(), ""
					}
				}
				if ev, ok := w.Env().Entities[ns]; ok {
					if v, _, found := ev.ContainsField(sym); found {
						return v.Value.GetType().String(), ""
					}
					if v, found := ev.ContainsMethod(sym); found {
						return v.Value.GetType().String(), ""
					}
				}
				if cv, ok := w.Env().Classes[ns]; ok {
					if v, _, found := cv.ContainsField(sym); found {
						return v.Value.GetType().String(), ""
					}
					if v, found := cv.ContainsMethod(sym); found {
						return v.Value.GetType().String(), ""
					}
				}
			}

			isBuiltin := ns == "Pewpew" || ns == "Fmath" || ns == "Math" || ns == "String" || ns == "Table"

			if env != nil {
				// Check for auto-generated API docs
				if d, ok := ApiDocs[ns+":"+sym]; ok {
					// Also need to get the type detail
					if v, ok := env.Scope.Variables[sym]; ok {
						return v.Value.GetType().String(), d
					}
					if ev, ok := env.Enums[sym]; ok {
						return "enum " + ev.Type.Name, d
					}
				}

				// Check variables
				if v, ok := env.Scope.Variables[sym]; ok {
					if isBuiltin || v.IsPub {
						return v.Value.GetType().String(), ""
					}
				}
				// Check enums in this namespace
				if ev, ok := env.Enums[sym]; ok {
					if isBuiltin || ev.IsPub {
						return "enum " + ev.Type.Name, ""
					}
				}
				// Check if ns is an enum
				if ev, ok := env.Enums[ns]; ok {
					if field, _, found := ev.ContainsField(sym); found {
						return field.Value.GetType().String(), ""
					}
				}
			}
		}
	}

	// Check Builtin
	if d, ok := builtinDocs[label]; ok {
		if v, ok := walker.BuiltinEnv.Scope.Variables[label]; ok {
			return v.Value.GetType().String(), d
		}
	}
	if v, ok := walker.BuiltinEnv.Scope.Variables[label]; ok {
		return v.Value.GetType().String(), "Builtin"
	}

	// Check current walker's context (entities, classes, aliases, and imports)
	if w != nil {
		env := w.Env()

		// 1. Current file types
		if ev, ok := env.Enums[label]; ok {
			return "enum " + ev.Type.Name, "Enum"
		}
		if ev, ok := env.Entities[label]; ok {
			return "entity " + ev.Type.Name, "Entity"
		}
		if cv, ok := env.Classes[label]; ok {
			return "class " + cv.Type.Name, "Class"
		}
		if alias, ok := env.Scope.AliasTypes[label]; ok {
			if d, ok := aliasDocs[label]; ok {
				return alias.UnderlyingType.String(), d
			}
			return alias.UnderlyingType.String(), "Alias"
		}

		// 2. Check imported namespaces via 'use'
		for _, imp := range env.Imports() {
			if imp.ThroughUse {
				impEnv := imp.Env()
				if v, ok := impEnv.Scope.Variables[label]; ok && v.IsPub {
					if d, ok := ApiDocs[impEnv.Name+":"+label]; ok {
						return v.Value.GetType().String(), d
					}
					return v.Value.GetType().String(), impEnv.Name
				}
				if ev, ok := impEnv.Enums[label]; ok && ev.IsPub {
					if d, ok := ApiDocs[impEnv.Name+":"+label]; ok {
						return "enum " + ev.Type.Name, d
					}
					return "enum " + ev.Type.Name, impEnv.Name
				}
				if cv, ok := impEnv.Classes[label]; ok && cv.IsPub {
					return "class " + label, impEnv.Name
				}
				if ev, ok := impEnv.Entities[label]; ok && ev.IsPub {
					return "entity " + label, impEnv.Name
				}
			}
		}

		// 3. Check used libraries (Pewpew, Fmath, etc.) - only those explicitly imported via 'use'
		for _, lib := range env.ImportedLibraries {
			libEnv := walker.BuiltinLibraries[lib]
			if libEnv != nil {
				if v, ok := libEnv.Scope.Variables[label]; ok {
					if d, ok := ApiDocs[libEnv.Name+":"+label]; ok {
						return v.Value.GetType().String(), d
					}
					return v.Value.GetType().String(), libEnv.Name
				}
				if ev, ok := libEnv.Enums[label]; ok {
					if d, ok := ApiDocs[libEnv.Name+":"+label]; ok {
						return "enum " + ev.Type.Name, d
					}
					return "enum " + ev.Type.Name, libEnv.Name
				}
			}
		}
	}

	return "", ""
}
