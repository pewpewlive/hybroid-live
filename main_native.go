//go:build !js

package main

import (
	"hybroid/cli"
)

func main() {
	// os.Chdir("examples/test")
	cli.RunApp()
}
