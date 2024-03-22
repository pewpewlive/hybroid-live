package main

import (
	"hybroid/cli"
	"os"
)

func main() {
	cwd, _ := os.Getwd()
	os.Chdir(cwd + "/example")
	cli.RunApp()
}
