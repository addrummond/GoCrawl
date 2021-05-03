package limited_source

import (
	"testing"

	MS "multiverse.io/crawler/crawler/test_helpers/mock_source"
)

func TestLimitedSource(t *testing.T) {
	universe := MS.MockSourceUniverse{
		Links: map[string][]string{
			"/":      []string{"/page1", "/page2", "/page3"},
			"/page1": []string{"/page2"},
			"/page2": []string{"/page3"},
			"/page3": []string{"/page2", "/", "page4"},
			"/page4": []string{"/page5"},
			"/page5": []string{"/page1"},
		},
		Assets: map[string][]string{},
	}

	source := &MS.MockSource{Universe: &universe, Url: "/"}
	limitedSource := MakeSource(
		source,
		10, // maxTotalRequests
		2,  // maxDepth
	)

	withinDepthOuts := limitedSource.GetOuts().Links[0].Source.GetOuts().Links[0].Source.GetOuts()
	if len(withinDepthOuts.Links) == 0 {
		t.Errorf("Unexpectedly got no links from within depth page")
	}

	overDepthOuts := limitedSource.GetOuts().Links[0].Source.GetOuts().Links[0].Source.GetOuts().Links[0].Source.GetOuts()
	if len(overDepthOuts.Links) != 0 {
		t.Errorf("Unexpectedly got links from over depth page")
	}

	// If we go right up to the request limit, we still get links back. We've
	// already made 6 reqs, so 3 more takes us right to the limit of 10.
	limitedSource.GetOuts()
	limitedSource.GetOuts()
	withinMaxReqsOuts := limitedSource.GetOuts()
	if len(withinMaxReqsOuts.Links) == 0 {
		t.Errorf("Unexpectedly got no links from within max reqs page")
	}

	// The next request is over the request limit, so we should get no links back.
	overMaxReqsOuts := limitedSource.GetOuts()
	if len(overMaxReqsOuts.Links) != 0 {
		t.Errorf("Unexpectedly got links from over max reqs page")
	}
}
