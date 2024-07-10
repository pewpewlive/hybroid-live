package cli

import (
	"fmt"
	"hybroid/cli/commands"
	"os"

	"github.com/urfave/cli/v2"
)

func RunApp() {
	app := &cli.App{
		Name:  "hybroid",
		Usage: "The Hybroid transpiler CLI",
		Commands: []*cli.Command{
			commands.Add(),
			commands.Build(),
			commands.Initialize(),
			commands.Watch(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
