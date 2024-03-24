package ast

import "hybroid/err"

type Walker struct {
	nodes   []Node
	current int
	Errors  []err.Error
}

func (w *Walker) Walk(nodes []Node) []Node {
	w.nodes = nodes

	newNodes := make([]Node, len(nodes))


	return newNodes
}