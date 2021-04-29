package mock_source

import (
	S "multiverse.io/crawler/crawler/source"
)

type MockSourceUniverse struct {
	Links  map[string][]string // Which URLs are linked to from each URL?
	Assets map[string][]string // Which assets are referenced from each URL?
}

type MockSource struct {
	Universe *MockSourceUniverse
	Url      string
}

func (s *MockSource) GetUrl() string {
	return s.Url
}

func (s *MockSource) GetOuts() S.Outs {
	var outs S.Outs
	if s.Universe != nil {
		for _, url := range s.Universe.Links[s.Url] {
			outs.Links = append(outs.Links, S.Link{Url: url, Source: &MockSource{s.Universe, url}})
		}
		for _, url := range s.Universe.Assets[s.Url] {
			outs.Assets = append(outs.Assets, S.Asset{Url: url})
		}
	}
	return outs
}
