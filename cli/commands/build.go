package commands

import (
	"encoding/json"
	"fmt"
	"hybroid/evaluator"
	"hybroid/generator"
	"hybroid/helpers"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func collectFiles(cwd string) ([]helpers.FileInformation, error) {
	files := make([]helpers.FileInformation, 0)
	err := fs.WalkDir(os.DirFS(cwd), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			ext := filepath.Ext(path)
			if ext != ".hyb" {
				return nil
			}

			directoryPath, err := filepath.Rel(cwd, filepath.Dir(cwd+"/"+path))
			if err != nil {
				return err
			}

			files = append(files, helpers.FileInformation{
				DirectoryPath: filepath.ToSlash(directoryPath),
				FileName:      strings.Replace(d.Name(), ".hyb", "", -1),
				FileExtension: ext,
			})
		}

		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
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

	if config.Project.Target != "ppl" {
		panic("other targets apart from 'ppl' are not implemented")
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
			files, filesErr := collectFiles(cwd)
			if filesErr != nil {
				err <- filesErr
				return
			}
			filesToBuild = append(filesToBuild, files...)
		}

		for _, file := range filesToBuild {
			fmt.Printf("%v", file)
			eval.AssignFile(file)
		}
		err <- eval.Action(cwd, outputDir)
	}(evalError)
	if err = <-evalError; err != nil {
		return fmt.Errorf("failed evaluation: %v", err)
	}

	return nil
}
