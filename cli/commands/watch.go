package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
)

type FileInformation struct {
	DirectoryPath string // The directory the file is located at (relative)
	FileName      string // The name of the file (without an extension)
	FileExtension string // The extension of the file
}

func (fi *FileInformation) Path() string {
	return fmt.Sprintf("%s/%s%s", fi.DirectoryPath, fi.FileName, fi.FileExtension)
}

func (fi *FileInformation) NewPath(start string, end string) string {
	return fmt.Sprintf("%s/%s/%s%s", start, fi.DirectoryPath, fi.FileName, end)
}

func Watch(ctx *cli.Context) error {
	cwd, _ := os.Getwd()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
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

					Build(ctx, FileInformation{directoryPath, fileName, fileExtension})
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
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})

	return nil
}
