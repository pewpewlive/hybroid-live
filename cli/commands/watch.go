package commands

import (
	"fmt"
	"hybroid/core"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

func Watch() *cli.Command {
	return &cli.Command{
		Name:        "watch",
		Aliases:     []string{"w"},
		Usage:       "Starts a watcher process",
		Description: "The Hybroid Live watcher will keep track of the project files and will automatically build them when they are updated, to remove the need for running the transpiler every time",
		Action: func(ctx *cli.Context) error {
			return watch(ctx)
		},
	}
}

func watch(ctx *cli.Context) error {
	cwd, _ := os.Getwd()

	configFile, err := os.ReadFile(cwd + "/hybconfig.toml")
	if err != nil {
		return fmt.Errorf("failed reading Hybroid Live config file: %v", err)
	}
	config := core.HybroidConfig{}
	if err := toml.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed parsing Hybroid Live config file: %v", err)
	}
	outputDirAbs, _ := filepath.Abs(filepath.Join(cwd, config.Project.OutputDirectory))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to start a watcher process: %s", err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		var pending *time.Timer
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				abs, _ := filepath.Abs(event.Name)
				if strings.HasPrefix(abs, outputDirAbs) || strings.Contains(event.Name, ".lua") {
					continue
				}
				if pending != nil {
					pending.Stop()
				}
				pending = time.AfterFunc(150*time.Millisecond, func() {
					Build_()
				})
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(cwd)
	if err != nil {
		return fmt.Errorf("failed to start a watcher process: %s", err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})

	return nil
}
