package tester

// import (
// 	"hybroid/evaluator"
// 	"hybroid/generators/lua"
// 	"hybroid/walker"
// 	"os"
// 	"strings"
// )

// type Tester struct {
// 	path      string
// 	evaluator evaluator.Evaluator
// }

// func NewTester(path string) Tester {
// 	walkers := map[string]*walker.Walker{}
// 	eval := evaluator.NewEvaluator(lua.Generator{},&walkers)
// 	eval.AssignFile(path, "")
// 	return Tester{
// 		path:path,
// 		evaluator: eval,
// 	}
// }


// func (t *Tester) Run() (string, error) {
// 	genSrc, err := t.evaluator.Action(false)

// 	dirFiles, fileReadErr := os.ReadDir(t.path)
// 	if fileReadErr != nil {
// 		// error
// 	}

// 	walkers := map[string]*walker.Walker{}

// 	for _, file := range dirFiles {
// 		if file.IsDir() {
// 			// add to next read list
// 			continue
// 		}

// 		filename := file.Name()

// 		if !strings.Contains(filename, ".hyb") {
// 			continue
// 		}

// 		eval := evaluator.NewEvaluator(lua.Generator{Scope: lua.GenScope{Src: lua.StringBuilder{}}}, &walkers)
// 		eval.AssignFile(t.path, "")
// 		_, evalErr := eval.Action(false)

// 		if evalErr != nil {
// 			// error
// 		}
// 	}

// 	return genSrc, err
// }

// func (t *Tester) Fail() {
	
// }