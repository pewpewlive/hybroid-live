package main

import (
	"hybroid/cli"
	"os"
)

func main() {
	cwd, _ := os.Getwd()
	os.Chdir(cwd + "/examples/" + os.Args[2])
	cli.RunApp()
}
