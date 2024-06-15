package helpers

import "reflect"

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