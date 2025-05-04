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
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	out, err := json.MarshalIndent(nodes, "", "  ")

	if err != nil {
		fmt.Println(err.Error())
	}

	writeErr := os.WriteFile(cwd+"/astdebug.json", out, os.ModePerm)

	if writeErr != nil {
		fmt.Println(err.Error())
	}
}
