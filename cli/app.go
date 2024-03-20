package cli

import (
	"fmt"
	"hybroid/cli/commands"
	"os"

	cli "github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:        "hybroid",
	Description: "The Hybroid transpiler CLI",
	Commands: []*cli.Command{
		{
			Name:        "build",
			Aliases:     []string{"b"},
			Usage:       "Builds a Hybroid project",
			Description: "This will take the current project in the location the command was ran, and will transpile the project into its destination folder, based on the config file",
			Action: func(ctx *cli.Context) error {
				return commands.Build(ctx)
			},
		},
		{
			Name:        "watch",
			Aliases:     []string{"w"},
			Usage:       "Starts a watcher proccess",
			Description: "The Hybroid watcher will keep track of the project files and will automatically build them when they are updated, to remove the need for running the transpiler every time",
			Action: func(ctx *cli.Context) error {
				return commands.Watch(ctx)
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initializes a new Hybroid project",
			Action: func(ctx *cli.Context) error {
				return commands.Initialize(ctx)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Installs packages from the Hybroid Package Repository",
			Action: func(ctx *cli.Context) error {
				return commands.Add(ctx)
			},
		},
	},
}

func RunApp() {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
