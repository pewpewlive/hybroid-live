package helpers

type Span struct {
	Start int
	End   int
}

func NewSpan(start, end int) Span {
	return Span{start, end}
}
