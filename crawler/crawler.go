package crawler

import (
	"runtime"
	"sync"

	G "multiverse.io/crawler/crawler/graph"
	S "multiverse.io/crawler/crawler/source"
)

type pendingRequest struct {
	node   *G.Node
	source S.Source
}

type pendingGraphUpdate struct {
	node *G.Node
	outs S.Outs
}

type AssetsMode int

const (
	AssetsModeIgnoreAssets = iota
	AssetsModeIncludeAssets
)

// Crawl constructs a graph by crawling a sites links and assets from a root
// source.
func Crawl(source S.Source, assetsMode AssetsMode) *G.Node {
	root := &G.Node{
		Url:        source.GetUrl(),
		Out:        []G.Edge{},
		Depth:      0,
		Popularity: 0,
		PureAsset:  false,
	}

	pendingGraphUpdateChan := make(chan pendingGraphUpdate)
	pendingRequestChan := make(chan pendingRequest)

	var wg sync.WaitGroup

	// Handle pending queries
	for i := 0; i < runtime.NumCPU(); i++ {
		go handleRequests(&wg, pendingGraphUpdateChan, pendingRequestChan)
	}

	// Handle pending graph updates. We only use one worker here because we don't
	// want the graph to be updated by multiple threads at once.
	go handleGraphUpdates(&wg, root, assetsMode, pendingGraphUpdateChan, pendingRequestChan)()

	wg.Add(1)
	pendingRequestChan <- pendingRequest{root, source}

	wg.Wait()

	close(pendingGraphUpdateChan)
	close(pendingRequestChan)

	G.Sort(root)
	return root
}

func handleRequests(wg *sync.WaitGroup, pendingGraphUpdateChan chan<- pendingGraphUpdate, pendingRequestChan <-chan pendingRequest) {
	for nextRequest := range pendingRequestChan {
		outs := nextRequest.source.GetOuts()
		wg.Add(1)
		pendingGraphUpdateChan <- pendingGraphUpdate{nextRequest.node, outs}
		wg.Done()
	}
}

func handleGraphUpdates(wg *sync.WaitGroup, root *G.Node, assetsMode AssetsMode, pendingGraphUpdateChan <-chan pendingGraphUpdate, pendingRequestChan chan<- pendingRequest) func() {
	urlToNode := make(map[string]*G.Node)
	urlToNode[root.Url] = root

	// We can't block when sending updates to pendingGraphUpdateChan because its
	// handler blocks on us, and that could give rise to a deadlock. So we add
	// updates to a queue and then send them to the channel as available.
	var queuedRequests []pendingRequest

	handleUpdate := func(pu pendingGraphUpdate, url string, edgeKind G.EdgeKind) *G.Node {
		var linkNode *G.Node

		if urlNode := urlToNode[url]; urlNode != nil {
			linkNode = urlNode
			linkNode.Popularity++
		} else {
			linkNode = &G.Node{
				Url:        url,
				Depth:      pu.node.Depth + 1,
				Popularity: 1,
				PureAsset:  true,
			}
		}

		pu.node.Out = append(pu.node.Out, G.Edge{Kind: edgeKind, Node: linkNode})

		return linkNode
	}

	return func() {
		for pu := range pendingGraphUpdateChan {
			for _, link := range pu.outs.Links {
				linkNode := handleUpdate(pu, link.Url, G.EdgeKindLink)
				linkNode.PureAsset = false

				if urlToNode[link.Url] == nil {
					wg.Add(1)
					queuedRequests = append(queuedRequests, pendingRequest{linkNode, link.Source})
				}

				urlToNode[link.Url] = linkNode
			}

			if assetsMode == AssetsModeIncludeAssets {
				for _, asset := range pu.outs.Assets {
					linkNode := handleUpdate(pu, asset.Url, G.EdgeKindAsset)
					urlToNode[asset.Url] = linkNode
				}
			}

			queuedRequests = dequeueRequests(pendingRequestChan, queuedRequests)

			wg.Done()
		}
	}
}

func dequeueRequests(pendingRequestChan chan<- pendingRequest, pending []pendingRequest) []pendingRequest {
	for len(pending) > 0 {
		r := pending[0]

		select {
		case pendingRequestChan <- r:
			pending = pending[1:]
		default:
			return pending
		}
	}

	return pending
}
