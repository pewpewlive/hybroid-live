package commands

import (
	"encoding/json"
	"fmt"
	"hybroid/evaluator"
	"hybroid/generators/lua"
	"hybroid/walker"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

func Build(ctx *cli.Context, files ...FileInformation) error {
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

	evalError := make(chan error)
	go func(err chan error) {
		outputDir := config.Project.OutputDirectory
		entryPoint := config.Level.EntryPoint

		if outputDir != "" {
			os.Mkdir(cwd+outputDir, 0644)
		}

		manifestConfig := config.Level
		manifestConfig.EntryPoint = "level.lua"
		manifestConfig.IsCasual = !config.Level.IsCasual
		manifest, manifestErr := json.MarshalIndent(manifestConfig, "", "  ")
		if manifestErr != nil {
			err <- fmt.Errorf("failed creating level manifest file: %v", manifestErr)
		}
		os.WriteFile(cwd+outputDir+"/manifest.json", manifest, 0644)

		envs := map[string]*walker.Environment{}

		eval := evaluator.NewEvaluator(lua.Generator{Scope: lua.GenScope{Src: lua.StringBuilder{}}}, &envs)
		var evalErr error = nil
		if len(files) == 0 {
			eval.AssignFile(cwd+entryPoint, cwd+outputDir+"/"+strings.Replace(entryPoint, ".hyb", ".lua", -1))
			evalErr = eval.Action()
		} else {
			for _, file := range files {
				sourceFilePath := file.Path()
				outputFilePath := file.NewPath(outputDir, ".lua")
				eval.AssignFile(cwd+sourceFilePath, cwd+outputFilePath)
				evalErr = eval.Action()
			}
		}
		err <- evalErr
	}(evalError)
	if err = <-evalError; err != nil {
		return fmt.Errorf("failed evaluation: %v", err)
	}

	return nil
}
