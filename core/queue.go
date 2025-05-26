package core

import "fmt"

type Queue[T any] struct {
	Name  string
	items []T
}

func NewQueue[T any](name string) Queue[T] {
	return Queue[T]{
		Name:  name,
		items: make([]T, 0),
	}
}

func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() T {
	count := q.Count()
	if count == 0 {
		panic(fmt.Sprintf("Attempt to pop the Queue (%q) of size: %d", q.Name, count))
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item
}

func (q *Queue[T]) Count() int {
	return len(q.items)
}

func (q *Queue[T]) Clear() {
	q.items = nil
}
