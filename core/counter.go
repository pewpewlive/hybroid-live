package core

import "fmt"

type Counter struct {
	val  int
	Name string
}

func NewCounter(name string, defaultValue ...int) Counter {
	val := 0

	if len(defaultValue) == 1 {
		val = defaultValue[0]
	}

	return Counter{val: val, Name: name}
}

func (c *Counter) Increment() {
	c.val++
}

func (c *Counter) Decrement() {
	if c.val == 0 {
		panic(fmt.Sprintf("Attempt to decrement Counter (%q) of value 0", c.Name))
	}

	c.val--
}

func (c *Counter) Value() int {
	return c.val
}
