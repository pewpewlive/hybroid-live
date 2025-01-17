package helpers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func IsZero[T comparable](v T) bool {
	var z T
	return v == z
}

func Contains[T comparable](list []T, thing T) bool {
	for i := range list {
		if list[i] == thing {
			return true
		}
	}
	return false
}

func GetValOfInterface[T, E any](val E) *T {
	value := reflect.ValueOf(val)
	ah := reflect.TypeFor[T]()
	if value.CanConvert(ah) {
		test := value.Convert(ah).Interface()
		tVal := test.(T)
		return &tVal
	}

	return nil
}

func XORNIL[T any](a, b *T) bool {
	if (a == nil || b == nil) && !(a == nil && b == nil) {
		return true
	}

	return false
}

func HasContents[T any](contents ...[]T) bool {
	sumContents := make([]T, 0)
	for _, v := range contents {
		sumContents = append(sumContents, v...)
	}
	return len(sumContents) != 0
}

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

func CollectFiles(cwd string) ([]FileInformation, error) {
	files := make([]FileInformation, 0)
	err := fs.WalkDir(os.DirFS(cwd), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			ext := filepath.Ext(path)
			if ext != ".hyb" {
				return nil
			}

			directoryPath, err := filepath.Rel(cwd, filepath.Dir(cwd+"/"+path))
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

type StackEntry[T any] struct {
	Name string
	Item T
}

type Stack[T any] struct {
	Name  string
	items []StackEntry[T]
}

func (s *Stack[T]) Push(name string, item T) {
	s.items = append(s.items, StackEntry[T]{Name: name, Item: item})
}

func (s *Stack[T]) Peek() StackEntry[T] {
	return s.items[len(s.items)-1]
}

func (s *Stack[T]) Pop(name string) StackEntry[T] {
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]

	if item.Name == name {
		fmt.Printf("Stack \"%s\" had an invalid pop name, expected: %s got: %s", s.Name, name, item.Name)
	}
	return item
}

func (s *Stack[T]) Count() int {
	return len(s.items)
}
