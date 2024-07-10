package commands

import "github.com/urfave/cli/v2"

func Add() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Installs packages from the PewPew Marketplace",
		Action: func(ctx *cli.Context) error {
			return add(ctx)
		},
	}
}

func add(ctx *cli.Context) error {
	return nil
}
