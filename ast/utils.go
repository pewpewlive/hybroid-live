package ast

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type NodeDrawer interface {
	DrawNode(str *strings.Builder, depth int) *strings.Builder
}

func DrawNodes(nodes []Node) {
	cwd, erre := os.Getwd()

	if erre != nil {
		panic(erre)
		//error
	}

	out, err := json.MarshalIndent(nodes, " ", " ")

	if err != nil {
		fmt.Print(err.Error())
	}

	writeErr := os.WriteFile(cwd+"/astdebug.json", out, 0644)

	if writeErr != nil {
		fmt.Print(err.Error())
	}
}
