package commands

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

type PackageConfig struct {
	PackageName    string
	PackageVersion string
}

type HybroidConfig struct {
	ProjectName     string
	EntryPoint      string
	OutputDirectory string
	Target          string // PPL or else throw an error
	Packages        []PackageConfig
}

func Initialize(ctx *cli.Context) error {
	file := HybroidConfig{"Hybroid Test", "example.hyb", "./out", "PPL", []PackageConfig{{"some-package", "0.1.1"}}}
	output, _ := toml.Marshal(file)
	fmt.Printf("%s", output)
	return nil
}
