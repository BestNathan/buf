package parser

import (
	"reflect"

	"github.com/bufbuild/protocompile/ast"
)

type Nodes []ast.Node

func (ns Nodes) Filter(fn func(ast.Node) bool) Nodes {
	nodes := Nodes{}
	for _, n := range ns {
		if fn(n) {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (ns Nodes) TypeNames() []string {
	ts := []string{}
	for _, n := range ns {
		ts = append(ts, reflect.TypeOf(n).String())
	}
	return ts
}
