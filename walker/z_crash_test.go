package walker_test

import (
	"hybroid/lsp"
	"hybroid/walker"
	"testing"
)

func TestReproCrash(t *testing.T) {
	code := `env helloworld as Level

use Pewpew
use MyHelper

let width = 500.5f
let height = 500f

let angle1 = -90d
let angle2 = 1r

// myhelper.hyb
Print(Greet("Hello"))
Pewpew:Print(MyHelper:Greet("Hello"))
Pewpew:NewShip(width, height, 0)

// mesh.hyb
Pewpew:SetEntityMesh(Pewpew:NewEntity(0f, 0f), MyMesh, 0

// sound.hyb
Pewpew:PlaySound(MySound, 0, 100f, 100f)
Pewpew:PlaySound(MySound, 0, 100f, 100f)
Pewpew:PlaySound(MySound, 0, 100f, 100f)
Pewpew:PlaySound(MySound, 0, 100f, 100f)
`

	// Test Analyze function directly
	walkerMap := make(map[string]*walker.Walker)
	result := lsp.Analyze("file:///c:/Users/Dominykas/Documents/Development/hybroid projects/hello-world/level.hyb", code, walkerMap, false)

	t.Logf("Analyze returned %d diagnostics", len(result.Diagnostics))
	for _, d := range result.Diagnostics {
		t.Logf("Diagnostic: %s (Line: %d)", d.Message, d.Range.Start.Line)
	}
}
