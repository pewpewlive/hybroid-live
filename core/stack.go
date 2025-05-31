package core

import "fmt"

type stackEntry[T any] struct {
	Name string
	Item T
}

type Stack[T any] struct {
	Name  string
	items []stackEntry[T]
}

func NewStack[T any](name string) Stack[T] {
	return Stack[T]{
		Name:  name,
		items: make([]stackEntry[T], 0),
	}
}

func (s *Stack[T]) Push(name string, item T) {
	s.items = append(s.items, stackEntry[T]{Name: name, Item: item})
}

func (s Stack[T]) Top() stackEntry[T] {
	count := s.Count()
	if count == 0 {
		panic(fmt.Sprintf("Attempt to get the top of the Stack (%q) of size: %d", s.Name, count))
	}
	return s.items[count-1]
}

func (s *Stack[T]) Pop(name string) stackEntry[T] {
	count := s.Count()
	if count == 0 {
		panic(fmt.Sprintf("Attempt to pop the Stack (%q) of size: %d", s.Name, count))
	}

	item := s.Top()

	if item.Name != name {
		s.printStack()
		panic(fmt.Sprintf("Attempt to pop the Stack (%q) with an invalid pop name, expected: %q, but got: %q", s.Name, item.Name, name))
	}

	s.items = s.items[:count-1]
	return item
}

func (s Stack[T]) Count() int {
	return len(s.items)
}

func (s Stack[T]) printStack() {
	fmt.Printf("Stack(%q):\n", s.Name)
	for i := s.Count() - 1; i >= 0; i-- {
		item := s.items[i]
		fmt.Printf("%d: | %s (%v) |\n", i, item.Name, item.Item)
	}
}
