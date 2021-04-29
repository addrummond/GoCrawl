package render

import (
	"strings"
	"testing"

	G "multiverse.io/crawler/crawler/graph"
)

func TestGraphToJson(t *testing.T) {
	root := makeTestGraph()

	json := GraphToJson(root)
	if len(json.Links) != 4 {
		t.Errorf("Expected 4 entries for 'links', got %v\n", len(json.Links))
	}
	for _, url := range []string{"A", "B", "C", "D"} {
		outs, ok := json.Links[url]
		if !ok {
			t.Errorf("Couldn't find entry in 'links' for %v\n", url)
		}
		var expectedEdgesOut bool
		switch url {
		case "A":
			expectedEdgesOut = len(outs) == 2 && !outs[0].IsAsset && outs[0].ToUrl == "B" && !outs[1].IsAsset && outs[1].ToUrl == "C"
		case "B":
			expectedEdgesOut = len(outs) == 2 && !outs[0].IsAsset && outs[0].ToUrl == "C" && outs[1].IsAsset && outs[1].ToUrl == "D"
		case "C":
			expectedEdgesOut = len(outs) == 1 && outs[0].IsAsset && outs[0].ToUrl == "D"
		case "D":
			expectedEdgesOut = len(outs) == 1 && !outs[0].IsAsset && outs[0].ToUrl == "A"
		}
		if !expectedEdgesOut {
			t.Errorf("Bad outward edges for %v\n", url)
		}
	}

	hasExpectedMetadata :=
		json.NodeMetadata["A"].Depth == 0 &&
			json.NodeMetadata["A"].Popularity == 1 &&
			!json.NodeMetadata["A"].PureAsset &&
			json.NodeMetadata["B"].Depth == 1 &&
			json.NodeMetadata["B"].Popularity == 1 &&
			!json.NodeMetadata["B"].PureAsset &&
			json.NodeMetadata["C"].Depth == 1 &&
			json.NodeMetadata["C"].Popularity == 2 &&
			!json.NodeMetadata["C"].PureAsset &&
			json.NodeMetadata["D"].Depth == 2 &&
			json.NodeMetadata["D"].Popularity == 2 &&
			json.NodeMetadata["D"].PureAsset

	if !hasExpectedMetadata {
		t.Errorf("Bad node metadata %+v\n", json.NodeMetadata)
	}
}

func TestExportHtml(t *testing.T) {
	root := makeTestGraph()

	html := ExportHtml(root)

	if !(strings.Contains(html, "<!DOCTYPE html>") && strings.Contains(html, "const STRIP_PREFIX =") && strings.Contains(html, "const GRAPH =") && strings.Contains(html, "const NODE_METADATA =")) {
		t.Errorf("Bad html output.\n")
	}
}

func makeTestGraph() *G.Node {
	//       A <---
	//      / \   |
	//     B-->C  |
	//      \ /   |
	//       D____|
	//
	//     (edges point down unless otherwise indicated)
	//
	a := &G.Node{Url: "A", Depth: 0, Popularity: 1, PureAsset: false}
	b := &G.Node{Url: "B", Depth: 1, Popularity: 1, PureAsset: false}
	c := &G.Node{Url: "C", Depth: 1, Popularity: 2, PureAsset: false}
	d := &G.Node{Url: "D", Depth: 2, Popularity: 2, PureAsset: true}
	a.Out = []G.Edge{{G.EdgeKindLink, b}, {G.EdgeKindLink, c}}
	b.Out = []G.Edge{{G.EdgeKindLink, c}, {G.EdgeKindAsset, d}}
	c.Out = []G.Edge{{G.EdgeKindAsset, d}}
	d.Out = []G.Edge{{G.EdgeKindLink, a}}
	return a
}
