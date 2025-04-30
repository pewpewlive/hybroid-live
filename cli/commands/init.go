package commands

import (
	"fmt"
	"hybroid/helpers"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

const levelTemplate = `env %s as Level

// Hello, world!
tick with i {
  if i %% 2 == 0 {
	  Pewpew:Print("Foo")
	} else {
	  Pewpew:Print("Bar")
	}
}
`

func Initialize() *cli.Command {
	return &cli.Command{
		Name:    "initialize",
		Aliases: []string{"init", "i"},
		Usage:   "Initializes a new Hybroid project",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "package", Required: true, Usage: "The package name of the project"},
			&cli.StringFlag{Name: "name", Required: true, Usage: "The level name of the project"},
			&cli.StringFlag{Name: "output", Required: true, Usage: "What output directory to use when building"},
		},
		Args:            true,
		SkipFlagParsing: false,
		Action: func(ctx *cli.Context) error {
			if ctx.NumFlags() != 4 {
				return fmt.Errorf("invalid amount of arguments (needed: 4, given: %v)", len(ctx.FlagNames()))
			}
			return initialize(ctx)
		},
	}
}

func initialize(ctx *cli.Context) error {
	pkgName, levelName, output := ctx.String("package"), ctx.String("name"), ctx.String("output")
	if pkgName == "" || levelName == "" || output == "" {
		return fmt.Errorf("invalid arguments, run `hybroid help init` for more information")
	}

	config := helpers.HybroidConfig{
		Level: helpers.LevelManifest{
			Name:         levelName,
			Descriptions: []string{"Change me!"},
			Information:  "Change me!",
			EntryPoint:   "level.hyb",
			IsCasual:     true,
		},
		Project: helpers.ProjectConfig{
			Name:            pkgName,
			OutputDirectory: output,
		},
	}

	configFile, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed generating Hybroid config file: %v", err)
	}

	if err = os.WriteFile("hybconfig.toml", configFile, 0644); err != nil {
		return fmt.Errorf("failed to write the Hybroid config file to disk: %v", err)
	}
	if err = os.WriteFile("level.hyb", []byte(fmt.Sprintf(levelTemplate, levelName)), 0644); err != nil {
		return fmt.Errorf("failed to write a level template to disk: %v", err)
	}

	return nil
}
