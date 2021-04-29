package graph

type EdgeKind int

const (
	EdgeKindAsset = iota
	EdgeKindLink
)

type Edge struct {
	Kind EdgeKind
	Node *Node
}

type Node struct {
	Url string
	Out []Edge

	// depth from root in our traversal order
	Depth int

	// how may nodes point at this one?
	Popularity int

	// is it an asset that's never linked to with <a>?
	PureAsset bool
}

// Traverse performs a reverse pre-order traversal on a graph. Each node is
// visited once (even in the presence of cycles).
func Traverse(node *Node, f func(node *Node)) {
	var stack []*Node
	visited := make(map[*Node]bool)

	for {
		if _, haveVisited := visited[node]; !haveVisited {
			visited[node] = true
			f(node)

			for i := range node.Out {
				e := &node.Out[i]
				if _, haveVisited := visited[e.Node]; !haveVisited {
					stack = append(stack, e.Node)
				}
			}
		}

		if len(stack) == 0 {
			break
		}

		// pop stack
		node, stack = stack[len(stack)-1], stack[:len(stack)-1]
	}
}
