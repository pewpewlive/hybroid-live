package commands

import (
	"encoding/json"
	"fmt"
	"hybroid/core"
	"hybroid/evaluator"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

func Build() *cli.Command {
	return &cli.Command{
		Name:        "build",
		Aliases:     []string{"b"},
		Usage:       "Builds a Hybroid Live project",
		Description: "This will take the current project in the location the command was ran, and will transpile the project into its destination folder, based on the config file",
		Action: func(ctx *cli.Context) error {
			return Build_(ctx)
		},
	}
}

func runEvaluator(config core.HybroidConfig, filesToBuild []core.FileInformation, cwd string, err chan<- error) {
	outputDir := config.Project.OutputDirectory

	if outputDir != "" {
		os.MkdirAll(cwd+outputDir, os.ModePerm)
	}

	manifestConfig := config.Level
	manifestConfig.EntryPoint = "level.lua"
	manifestConfig.IsCasual = !config.Level.IsCasual
	manifest, manifestErr := json.MarshalIndent(manifestConfig, "", "  ")
	if manifestErr != nil {
		err <- fmt.Errorf("failed creating level manifest file: %v", manifestErr)
		return
	}
	os.WriteFile(filepath.Join(cwd, outputDir, "/manifest.json"), manifest, os.ModePerm)

	if len(filesToBuild) == 0 {
		files, filesErr := core.CollectFiles(cwd)
		if filesErr != nil {
			err <- filesErr
			return
		}
		filesToBuild = append(filesToBuild, files...)
	}

	evaluator := evaluator.NewEvaluator(filesToBuild)
	err <- evaluator.Action(cwd, outputDir)
}

func Build_(ctx *cli.Context, filesToBuild ...core.FileInformation) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed getting current working directory: %v", err)
	}
	cwd += "/"

	configFile, err := os.ReadFile(cwd + "hybconfig.toml")
	if err != nil {
		return fmt.Errorf("failed reading Hybroid Live config file: %v", err)
	}

	config := core.HybroidConfig{}
	if err := toml.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed parsing Hybroid Live config file: %v", err)
	}

	evalError := make(chan error)
	go runEvaluator(config, filesToBuild, cwd, evalError)
	if err = <-evalError; err != nil {
		return fmt.Errorf("failed evaluation: %v", err)
	}

	return nil
}
