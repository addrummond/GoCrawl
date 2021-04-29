package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"

	G "multiverse.io/crawler/crawler/graph"
)

//go:embed cytoscape.min.js
var cytoscapeSrc string

//go:embed render.js
var renderSrc string

type NodeMetadata struct {
	Depth      int
	Popularity int
	PureAsset  bool
}

type Link struct {
	IsAsset bool
	ToUrl   string
}

type GraphJson struct {
	Links        map[string][]Link
	NodeMetadata map[string]NodeMetadata
}

func GraphToJson(node *G.Node) GraphJson {
	links := make(map[string][]Link)
	nodeMetadata := make(map[string]NodeMetadata)

	G.Traverse(node, func(node *G.Node) {
		links[node.Url] = []Link{}
		for _, out := range node.Out {
			if out.Node == node {
				continue
			}

			links[node.Url] =
				append(links[node.Url],
					Link{
						IsAsset: out.Kind == G.EdgeKindAsset,
						ToUrl:   out.Node.Url,
					})
		}

		nodeMetadata[node.Url] =
			NodeMetadata{
				Depth:      node.Depth,
				Popularity: node.Popularity,
				PureAsset:  node.PureAsset,
			}
	})

	return GraphJson{
		Links:        links,
		NodeMetadata: nodeMetadata,
	}
}

func ExportHtml(node *G.Node) string {
	parsed, _ := url.Parse(node.Url)
	var stripPrefix string = parsed.Scheme + "://" + parsed.Host + "/"
	stripPrefixJson, _ := json.Marshal(stripPrefix)

	graphJson := GraphToJson(node)
	marshaledLinks, _ := json.Marshal(graphJson.Links)
	marshaledNodeMetadata, _ := json.Marshal(graphJson.NodeMetadata)
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Graph</title>
		<style>
		  body {
				width: 100vw;
				height: 100vh;
			}
			svg {
				width: 100vw;
				height: 100vh;
			}
		</style>
		<script>
		%s
		</script>
		<script>
		const STRIP_PREFIX = %s;
		const GRAPH = %s;
		const NODE_METADATA = %s;
		</script>
		<script>
		%s
		</script>
	</head>
	<body>
	</body>
	</html>
	`, cytoscapeSrc, stripPrefixJson, marshaledLinks, marshaledNodeMetadata, renderSrc)
}
