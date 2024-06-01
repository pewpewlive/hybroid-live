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

func GetValOfInterface[T any, E any](val E) *T {
	value := reflect.ValueOf(val)
	ah := reflect.TypeFor[T]()
	if value.CanConvert(ah) {
		test := value.Convert(ah).Interface()
		tVal := test.(T)
		return &tVal
	}

	return nil
}
