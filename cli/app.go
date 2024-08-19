package cli

import (
	"fmt"
	"hybroid/cli/commands"
	"os"

	"github.com/urfave/cli/v2"
)

func RunApp() {
	app := &cli.App{
		Name:      "hybroid",
		Usage:     "The Hybroid transpiler CLI",
		Version:   "0.0.0",
		Copyright: "Copyright (C) Hybroid Team, 2024\nLicensed under Apache-2.0",
		Commands: []*cli.Command{
			commands.Add(),
			commands.Build(),
			commands.Initialize(),
			commands.Watch(),
			// LSP is not yet implemented
			// commands.Lsp(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
