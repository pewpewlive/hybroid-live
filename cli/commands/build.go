package commands

import (
	"fmt"
	"hybroid/evaluator"
	"hybroid/generators/lua"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

func Build(ctx *cli.Context) error {
	cwd, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("Error getting current working directory: %v", err)
	}

	configFile, err := os.ReadFile(cwd + "/example/hybconfig.toml")

	if err != nil {
		return fmt.Errorf("Error reading Hybroid config file: %v", err)
	}

	config := HybroidConfig{}
	toml.Unmarshal(configFile, &config)

	var eval evaluator.Evaluator
	if config.Target == "ppl" {
		eval = evaluator.New(cwd+"/example", cwd+config.OutputDirectory, lua.Generator{})
		eval.SrcPath += config.EntryPoint
		eval.DstPath += strings.Replace(config.EntryPoint, ".hyb", ".lua", -1)
	} else {
		return fmt.Errorf("Other targets are not implemented yet. Only 'ppl' is allowed.")
	}

	eval.Action()

	return nil
}
