package cli

import (
	"fmt"
	"hybroid/cli/commands"
	"os"

	"github.com/urfave/cli/v2"
)

func RunApp() {
	app := &cli.App{
		Name:      "hybroid-live",
		Usage:     "The Hybroid Live transpiler CLI",
		Version:   "0.1.0",
		Copyright: "© Hybroid Team, 2026\nLicensed under Apache-2.0",
		Commands: []*cli.Command{
			commands.Add(),
			commands.Build(),
			commands.Initialize(),
			commands.Watch(),
			commands.Lsp(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
