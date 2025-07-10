package commands

import (
	"fmt"
	"hybroid/core"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to start a watcher process: %s", err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) && !strings.Contains(event.Name, ".lua") {
					directoryPath, _ := filepath.Rel(cwd, filepath.Dir(event.Name))
					fileName := strings.Split(filepath.Base(event.Name), ".")[0]
					fileExtension := filepath.Ext(event.Name)

					Build_(core.FileInformation{DirectoryPath: directoryPath, FileName: fileName, FileExtension: fileExtension})
				}
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
