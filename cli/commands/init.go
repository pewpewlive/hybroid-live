package commands

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type HybroidConfig struct {
	ProjectName string
}

func Initialize() error {
	var file HybroidConfig
	s, _ := os.ReadFile("config.toml")	
	toml.Unmarshal([]byte(s), &file)
	fmt.Printf("%v", file)
  return nil
}