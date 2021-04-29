package crawler

import (
	"testing"

	MS "multiverse.io/crawler/crawler/test_helpers/mock_source"
)

func TestCrawlOnSimpleSiteWithCycles(t *testing.T) {
	universe := MS.MockSourceUniverse{
		Links: map[string][]string{
			"/":      []string{"/page1", "/page2", "/page3"},
			"/page1": []string{"/page2"},
			"/page2": []string{"/page3"},
			"/page3": []string{"/page2", "/"},
		},
		Assets: map[string][]string{},
	}

	source := &MS.MockSource{Universe: &universe, Url: "/"}

	root := Crawl(source, AssetsModeIncludeAssets)
	if root == nil {
		t.Errorf("Unexpected nil root")
	}

	if root.Url != "/" {
		t.Errorf("Unexpected root URL (%v instead of /)", root.Url)
	}

	if len(root.Out) != 3 {
		t.Errorf("Unexpected graph structure (wrong number of outward edges for root node)")
	}

	if root.Out[0].Node.Url != "/page1" || root.Out[1].Node.Url != "/page2" || root.Out[2].Node.Url != "/page3" {
		t.Errorf("Unexpected URLS for children of root: %v %v %v", root.Out[0].Node.Url, root.Out[1].Node.Url, root.Out[2].Node.Url)
	}

	if len(root.Out[1].Node.Out) != 1 || root.Out[1].Node.Out[0].Node.Url != "/page3" {
		t.Errorf("Unexpected outward link from page 2; expected /page3, got %v\n", root.Out[1].Node.Out[0].Node.Url)
	}

	if len(root.Out[2].Node.Out) != 2 || root.Out[2].Node.Out[0].Node.Url != "/" || root.Out[2].Node.Out[1].Node.Url != "/page2" {
		t.Errorf("Unexpected outward links from page 3; expected / and /page2, got %v and %v\n", root.Out[2].Node.Out[0].Node.Url, root.Out[2].Node.Out[1].Node.Url)
	}

	page1 := root.Out[0].Node
	page2 := root.Out[1].Node
	page3 := root.Out[2].Node

	if root.Depth != 0 || root.Popularity != 1 || root.PureAsset {
		t.Errorf("Bad fields for root Depth=%v, Popularity=%v, PureAsset=%v\n", root.Depth, root.Popularity, root.PureAsset)
	}

	if page1.Depth != 1 || page1.Popularity != 1 || page1.PureAsset {
		t.Errorf("Bad fields for page1: Depth=%v, Popularity=%v, PureAsset=%v\n", page1.Depth, page1.Popularity, page1.PureAsset)
	}

	if page2.Depth != 1 || page2.Popularity != 3 || page2.PureAsset {
		t.Errorf("Bad fields for page2: Depth=%v, Popularity=%v, PureAsset=%v\n", page2.Depth, page2.Popularity, page2.PureAsset)
	}

	if page3.Depth != 1 || page3.Popularity != 2 || page3.PureAsset {
		t.Errorf("Bad fields for page3: Depth=%v, Popularity=%v, PureAsset=%v\n", page3.Depth, page3.Popularity, page3.PureAsset)
	}
}

func TestCrawlOnSimpleSiteWithCyclesIncludingAssets(t *testing.T) {
	universe := MS.MockSourceUniverse{
		Links: map[string][]string{
			"/":      []string{"/page1", "/page2"},
			"/page1": []string{},
			"/page2": []string{},
		},
		Assets: map[string][]string{
			"/page1": []string{"/asset1", "/asset2"},
			"/page2": []string{"/asset2"},
		},
	}

	source := &MS.MockSource{Universe: &universe, Url: "/"}

	root := Crawl(source, AssetsModeIncludeAssets)
	if root == nil {
		t.Errorf("Unexpected nil root")
	}

	if len(root.Out) != 2 {
		t.Errorf("Unexpected graph structure (wrong number of outward edges for root node)")
	}

	if root.Out[0].Node.Url != "/page1" || root.Out[1].Node.Url != "/page2" {
		t.Errorf("Unexpected URLS for children of root: %v %v", root.Out[0].Node.Url, root.Out[1].Node.Url)
	}

	if len(root.Out[0].Node.Out) != 2 {
		t.Errorf("Unexpected number of outward links from page 1")
	}

	if len(root.Out[1].Node.Out) != 1 {
		t.Errorf("Unexpected number of outward links from page 2")
	}
}

func TestCrawlOnSimpleSiteWithCyclesNotIncludingAssets(t *testing.T) {
	universe := MS.MockSourceUniverse{
		Links: map[string][]string{
			"/":      []string{"/page1", "/page2"},
			"/page1": []string{},
			"/page2": []string{},
		},
		Assets: map[string][]string{
			"/page1": []string{"/asset1", "/asset2"},
			"/page2": []string{"/asset2"},
		},
	}

	source := &MS.MockSource{Universe: &universe, Url: "/"}

	root := Crawl(source, AssetsModeIncludeAssets)
	if root == nil {
		t.Errorf("Unexpected nil root")
	}

	if len(root.Out) != 2 {
		t.Errorf("Unexpected graph structure (wrong number of outward edges for root node)")
	}

	if root.Out[0].Node.Url != "/page1" || root.Out[1].Node.Url != "/page2" {
		t.Errorf("Unexpected URLS for children of root: %v %v", root.Out[0].Node.Url, root.Out[1].Node.Url)
	}

	if len(root.Out[0].Node.Out) != 2 {
		t.Errorf("Unexpected number of outward links from page 1")
	}

	if len(root.Out[1].Node.Out) != 1 {
		t.Errorf("Unexpected number of outward links from page 2")
	}
}
