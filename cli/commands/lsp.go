package commands

import (
	"hybroid/lsp"

	"github.com/urfave/cli/v2"
)

func Lsp() *cli.Command {
	return &cli.Command{
		Name:  "lsp",
		Usage: "Starts Internal Hybroid Language Server",
		Action: func(ctx *cli.Context) error {
			return languageServer(ctx)
		},
	}
}

func languageServer(ctx *cli.Context) error {
	lsp.Init()
	return nil
}
