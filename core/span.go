package core

type Span[T any] struct {
	Start T
	End   T
}

func NewSpan[T any](start, end T) Span[T] {
	return Span[T]{start, end}
}

func (s *Span[T]) SetStart(start T) {
	s.Start = start
}

func (s *Span[T]) SetEnd(end T) {
	s.End = end
}
