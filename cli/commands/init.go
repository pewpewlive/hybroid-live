package commands

import (
	"fmt"
	"os"

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
	Target          string // ppl or else throw an error
	Packages        []PackageConfig
}

func Initialize(ctx *cli.Context) error {
	config := HybroidConfig{"Hybroid Test", "example.hyb", "/out", "ppl", []PackageConfig{{"some-package", "0.1.1"}}}

	output, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed generating Hybroid config file: %v", err)
	}

	os.WriteFile("hybconfig.toml", output, 0644)
	return nil
}
