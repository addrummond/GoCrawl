package graph

import "sort"

// For determinism, we order the list of outgoing edges by
//   (i)  assets before links
//   (ii) lexicographic order of URL
func Sort(node *Node) {
	Traverse(node, func(node *Node) {
		sortEdges(node)
	})
}

func sortEdges(root *Node) {
	sort.Sort(EdgeOrder(root.Out))
}

type EdgeOrder []Edge

func (a EdgeOrder) Len() int { return len(a) }
func (a EdgeOrder) Less(i, j int) bool {
	if a[i].Kind == EdgeKindAsset && a[j].Kind == EdgeKindLink {
		return true
	}
	if a[i].Kind == EdgeKindLink && a[j].Kind == EdgeKindAsset {
		return false
	}
	return a[i].Node.Url < a[j].Node.Url
}
func (a EdgeOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
