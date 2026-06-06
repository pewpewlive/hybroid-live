package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func CollectFiles(dir string) ([]File, error) {
	files := make([]File, 0)
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			ext := filepath.Ext(path)
			if ext != ".hyb" {
				return nil
			}

			directoryPath, err := filepath.Rel(dir, filepath.Dir(dir+"/"+path))
			if err != nil {
				return err
			}

			files = append(files, File{
				DirectoryPath: filepath.ToSlash(directoryPath),
				FileName:      strings.ReplaceAll(d.Name(), ".hyb", ""),
				FileExtension: ext,
			})
		}

		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}
