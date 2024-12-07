package commands

import (
	"encoding/json"
	"fmt"
	"hybroid/evaluator"
	"hybroid/generator"
	"hybroid/helpers"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

func Build() *cli.Command {
	return &cli.Command{
		Name:        "build",
		Aliases:     []string{"b"},
		Usage:       "Builds a Hybroid project",
		Description: "This will take the current project in the location the command was ran, and will transpile the project into its destination folder, based on the config file",
		Action: func(ctx *cli.Context) error {
			return build(ctx)
		},
	}
}

func build(ctx *cli.Context, filesToBuild ...helpers.FileInformation) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed getting current working directory: %v", err)
	}
	cwd += "/"

	configFile, err := os.ReadFile(cwd + "hybconfig.toml")
	if err != nil {
		return fmt.Errorf("failed reading Hybroid config file: %v", err)
	}

	config := helpers.HybroidConfig{}
	if err := toml.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed parsing Hybroid config file: %v", err)
	}

	evalError := make(chan error)
	go func(err chan error) {
		outputDir := config.Project.OutputDirectory

		if outputDir != "" {
			os.MkdirAll(cwd+outputDir, 0644)
		}

		manifestConfig := config.Level
		manifestConfig.EntryPoint = "level.lua"
		manifestConfig.IsCasual = !config.Level.IsCasual
		manifest, manifestErr := json.MarshalIndent(manifestConfig, "", "  ")
		if manifestErr != nil {
			err <- fmt.Errorf("failed creating level manifest file: %v", manifestErr)
		}
		os.WriteFile(filepath.Join(cwd, outputDir, "/manifest.json"), manifest, 0644)

		eval := evaluator.NewEvaluator(generator.Generator{
			Scope: generator.GenScope{
				Src: generator.StringBuilder{},
			},
		})

		if len(filesToBuild) == 0 {
			files, filesErr := helpers.CollectFiles(cwd)
			if filesErr != nil {
				err <- filesErr
				return
			}
			filesToBuild = append(filesToBuild, files...)
		}

		for _, file := range filesToBuild {
			eval.AssignFile(file)
		}
		err <- eval.Action(cwd, outputDir)
	}(evalError)
	if err = <-evalError; err != nil {
		return fmt.Errorf("failed evaluation: %v", err)
	}

	return nil
}
