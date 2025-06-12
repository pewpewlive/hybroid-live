package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func MapsAreSame[T comparable, E comparable](map1 map[E]T, map2 map[E]T) bool {
	if len(map1) != len(map2) {
		return false
	}

	for k := range map1 {
		if _, found := map2[k]; !found {
			return false
		}
	}

	return true
}

// should be used only for simple lists
//
// if a list has a pointer to a value this wont work
func ListsAreSame[T comparable](list1 []T, list2 []T) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i := range list1 {
		if list1[i] != list2[i] {
			return false
		}
	}

	return true
}

func CollectFiles(dir string) ([]FileInformation, error) {
	files := make([]FileInformation, 0)
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

			files = append(files, FileInformation{
				DirectoryPath: filepath.ToSlash(directoryPath),
				FileName:      strings.Replace(d.Name(), ".hyb", "", -1),
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
