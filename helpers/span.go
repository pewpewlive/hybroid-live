package helpers

type Span[T any] struct {
	Start T
	End   T
}

func NewSpan[T any](start, end T) Span[T] {
	return Span[T]{start, end}
}
