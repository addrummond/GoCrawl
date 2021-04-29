package graph

import (
	"testing"
)

func TestTraverse(t *testing.T) {
	//       A <---
	//      / \   |
	//     B-->C  |
	//      \ /   |
	//       D____|
	//
	//     (edges point down unless otherwise indicated)
	//
	a := &Node{Url: "A"}
	b := &Node{Url: "B"}
	c := &Node{Url: "C"}
	d := &Node{Url: "D"}
	a.Out = []Edge{{EdgeKindLink, b}, {EdgeKindLink, c}}
	b.Out = []Edge{{EdgeKindLink, c}, {EdgeKindLink, d}}
	c.Out = []Edge{{EdgeKindLink, d}}
	d.Out = []Edge{{EdgeKindLink, a}}

	count := 0
	seen := make(map[*Node]bool)

	Traverse(a, func(node *Node) {
		count++
		if _, ok := seen[node]; ok {
			t.Errorf("Unexpectedly visited node %v twice.\n", node.Url)
		}
	})

	if count != 4 {
		t.Errorf("Expected visitor to be called 4 times but it was called %v times\n", count)
	}
}
