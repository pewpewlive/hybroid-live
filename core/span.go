package core

type Span struct {
	StartByte int
	EndByte   int
	Line      int
	Column    int
	Length    int
}

func NewSpan(startByte, endByte, line, column int) Span {
	return Span{
		StartByte: startByte,
		EndByte:   endByte,
		Line:      line,
		Column:    column,
		Length:    endByte - startByte,
	}
}

func (s *Span) UpdateEnd(endByte int) {
	s.EndByte = endByte
	s.Length = s.EndByte - s.StartByte
}

func MergeSpans(start, end Span) Span {
	return Span{
		StartByte: start.StartByte,
		EndByte:   end.EndByte,
		Line:      start.Line,
		Column:    start.Column,
		Length:    end.EndByte - start.StartByte,
	}
}
