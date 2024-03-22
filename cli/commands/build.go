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
		return fmt.Errorf("failed getting current working directory: %v", err)
	}
	cwd += "/"

	configFile, err := os.ReadFile(cwd + "hybconfig.toml")
	if err != nil {
		return fmt.Errorf("failed reading Hybroid config file: %v", err)
	}

	config := HybroidConfig{}
	if err := toml.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed parsing Hybroid config file: %v", err)
	}

	if config.Project.Target != "ppl" {
		panic("other targets apart from 'ppl' are not implemented")
	}

	ok := make(chan bool)
	go func(okay chan bool) {
		eval := evaluator.New(cwd+config.Level.EntryPoint, cwd+config.Project.OutputDirectory+"/"+strings.Replace(config.Level.EntryPoint, ".hyb", ".lua", -1), lua.Generator{})
		eval.Action()
		//bar.Add(1)
		okay <- true
	}(ok)
	if !<-ok {
		return fmt.Errorf("failed evaluation")
	}
	//bar.Finish()

	return nil
}
