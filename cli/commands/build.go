package commands

import (
	"encoding/json"
	"fmt"
	"hybroid/core"
	"hybroid/evaluator"
	"os"
	"path/filepath"
	"regexp"

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
			return Build_()
		},
	}
}

func runEvaluator(config core.HybroidConfig, filesToBuild []core.FileInformation, cwd string) error {
	outputDir := config.Project.OutputDirectory

	if outputDir != "" {
		os.MkdirAll(cwd+outputDir, os.ModePerm)
	}

	regex := regexp.MustCompile(`(?:#\{[0-9a-fA-F]+\})|(?:#[0-9a-fA-F]{1,8})`)
	matches := regex.FindAllStringSubmatchIndex(config.Level.Name, -1)
	index := -1
	var err error
	config.Level.Name = string(regex.ReplaceAllFunc([]byte(config.Level.Name), func(b []byte) []byte {
		index++
		var color []byte
		if b[0] == '#' && b[1] == '{' {
			color = append([]byte{'#'}, b[2:len(b)-1]...)
			if len(color) == 5 {
				color = []byte{'#', color[1], color[1], color[2], color[2], color[3], color[3], color[4], color[4]}
			}
		} else {
			color = b
		}
		colorLen := len(color)
		if colorLen != 9 {
			err = fmt.Errorf("color '%s' at around %d-%d in level title must be either 4 characters (braced) or 8 characters (braced/raw)", b, matches[index][0]+9, matches[index][1]+9)
		}
		//fmt.Printf("%s\n", b)
		return color
	}))
	if err != nil {
		return err
	}

	manifestConfig := config.Level
	manifestConfig.EntryPoint = "level.lua"
	manifestConfig.IsCasual = !config.Level.IsCasual

	if len(filesToBuild) == 0 {
		files, filesErr := core.CollectFiles(cwd)
		if filesErr != nil {
			return filesErr

		}
		filesToBuild = append(filesToBuild, files...)
	}

	evaluator := evaluator.NewEvaluator(filesToBuild)
	err = evaluator.Action(cwd, outputDir)
	if err != nil {
		return err
	}

	manifest, manifestErr := json.MarshalIndent(manifestConfig, "", "  ")
	if manifestErr != nil {
		return fmt.Errorf("failed creating level manifest file: %v", manifestErr)

	}
	os.WriteFile(filepath.Join(cwd, outputDir, "/manifest.json"), manifest, os.ModePerm)

	return nil
}

func Build_(filesToBuild ...core.FileInformation) error {
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

	err = runEvaluator(config, filesToBuild, cwd)
	if err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	return nil
}
