package commands

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

type LevelManifest struct {
	Name              string         `toml:"name"`
	Description       string         `toml:"description"`
	Information       string         `toml:"information"`
	EntryPoint        string         `toml:"entry_point"`
	IsCasual          bool           `toml:"casual"`
	MedalRequirements map[string]int `toml:"medal_requirements"`
}

type ProjectConfig struct {
	Target          string `toml:"target"` // ppl or else throw an error
	OutputDirectory string `toml:"output_directory"`
}

type HybroidConfig struct {
	Level   LevelManifest `toml:"level"`
	Project ProjectConfig `toml:"project"`
	//Packages        []PackageConfig `toml:"packages"`
}

func Initialize(ctx *cli.Context) error {
	config := HybroidConfig{LevelManifest{Name: ctx.Args().Get(0), Description: "Change me!", Information: "Change me!", EntryPoint: "level.hyb", IsCasual: true}, ProjectConfig{Target: ctx.Args().Get(1), OutputDirectory: ctx.Args().Get(2)}}

	output, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed generating Hybroid config file: %v", err)
	}

	os.WriteFile("hybconfig.toml", output, 0644)
	return nil
}
