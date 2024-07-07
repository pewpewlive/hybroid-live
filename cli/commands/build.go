package commands

import (
	"encoding/json"
	"fmt"
	"hybroid/evaluator"
	"hybroid/generators/lua"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

// func GetDirFiles(cwd string, dir string) ([]FileInformation, error) {
// 	dirEntries, err := os.ReadDir(cwd+dir)
// 	files := make([]FileInformation, 0)

// 	if err != nil { return files, err }

// 	for _, dirEntry := range dirEntries {
// 		info, err := dirEntry.Info()
// 		if err != nil { return files, err }

// 		if info.IsDir() {
// 			//fmt.Println(cwd+info.Name())
// 			subfiles, err := GetDirFiles(cwd, dir+info.Name())
// 			if err != nil { return files, err }

// 			files = append(files, subfiles...)
// 		}else {
// 			fmt.Printf("Dir: %s, Info: %s \n", dir, info.Name())
// 			pathFile, err := filepath.Rel(cwd, cwd+"/"+dir+"/"+info.Name())
// 			fmt.Println(pathFile)
// 			if err != nil { return files, err }

// 			ext := filepath.Ext(pathFile)
// 			if ext != ".hyb" {
// 				continue
// 			}
// 			files = append(files, FileInformation{
// 				DirectoryPath: dir,
// 				FileName: strings.Replace(info.Name(),".hyb", "", -1),
// 				FileExtension: ext,
// 			})
// 		}
// 	}

// 	return files, nil
// }

func Accumulate(cwd string, dir string) ([]FileInformation, error) {
	dirEntries, err := os.ReadDir(cwd+dir)
	files := make([]FileInformation, 0)

	if err != nil { return files, err }

	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil { return files, err }

		if info.IsDir() {
			//fmt.Println(cwd+info.Name())
			subfiles, err := GetDirFiles(cwd, dir+info.Name())
			if err != nil { return files, err }

			files = append(files, subfiles...)
		}else {
			fmt.Printf("Dir: %s, Info: %s \n", dir, info.Name())
			pathFile, err := filepath.Rel(cwd, cwd+"/"+dir+"/"+info.Name())
			fmt.Println(pathFile)
			if err != nil { return files, err }

			ext := filepath.Ext(pathFile)
			if ext != ".hyb" {
				continue
			}
			files = append(files, FileInformation{
				DirectoryPath: dir,
				FileName: strings.Replace(info.Name(),".hyb", "", -1),
				FileExtension: ext,
			})
		}
	}

	return files, nil
}

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

		eval := evaluator.NewEvaluator(lua.Generator{Scope: lua.GenScope{Src: lua.StringBuilder{}}})

		if len(files) == 0 {
			var filesErr error
			files, filesErr = GetDirFiles(cwd, "")
			if filesErr != nil {
				err <- filesErr
				return
			}
			fmt.Printf("Files:\n %v", files)
		}

		if len(files) == 0 {			
			eval.AssignFile(cwd+entryPoint, cwd+outputDir+"/"+strings.Replace(entryPoint, ".hyb", ".lua", -1))
			err <- eval.Action()
		} else {
			for _, file := range files {
				sourceFilePath := cwd+file.Path()
				outputFilePath := cwd+file.NewPath(outputDir, ".lua")
				eval.AssignFile(sourceFilePath, outputFilePath)
			}
			err <- eval.Action()
		}
	}(evalError)
	if err = <-evalError; err != nil {
		return fmt.Errorf("failed evaluation: %v", err)
	}

	return nil
}
