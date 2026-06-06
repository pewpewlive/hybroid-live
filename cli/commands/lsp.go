package commands

import (
	"hybroid/lsp"

	"github.com/urfave/cli/v2"
)

func Lsp() *cli.Command {
	return &cli.Command{
		Name:    "language-server",
		Aliases: []string{"server", "lsp"},
		Usage:   "Starts HybroidLS, an integrated Language Server for Hybroid Live",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable verbose debug logging",
			},
		},
		Action: func(ctx *cli.Context) error {
			return languageServer(ctx)
		},
	}
}

func languageServer(ctx *cli.Context) error {
	lsp.Init(ctx.Bool("debug"))
	return nil
}
