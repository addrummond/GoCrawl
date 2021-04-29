// Package limited_source wraps a Source with limits on the depth of traversal
// from the origin and the total number of requests. If either is exceeded, the
// Source 'lies' and says that the page in question has no assets and no links.
package limited_source

import (
	"sync/atomic"

	S "multiverse.io/crawler/crawler/source"
)

type LimitedSource struct {
	source           S.Source
	maxTotalRequests uint64
	maxDepth         uint64
	totalRequestsPtr *uint64
	depth            uint64
}

func (s *LimitedSource) GetUrl() string {
	return s.source.GetUrl()
}

func (s *LimitedSource) GetOuts() S.Outs {
	// GetOuts() will be called from multiple threads, so we must use an atomic
	// update for the request counter.
	atomic.AddUint64(s.totalRequestsPtr, 1)

	if s.depth > s.maxDepth || *s.totalRequestsPtr > s.maxTotalRequests {
		return S.Outs{}
	}

	origOuts := s.source.GetOuts()

	newLinks := make([]S.Link, len(origOuts.Links))
	for i, l := range origOuts.Links {
		newLinks[i].Url = l.Url
		newLinks[i].Source = &LimitedSource{
			source:           l.Source,
			maxTotalRequests: s.maxTotalRequests,
			maxDepth:         s.maxDepth,
			totalRequestsPtr: s.totalRequestsPtr,
			depth:            s.depth + 1,
		}
	}

	return S.Outs{
		Assets: origOuts.Assets,
		Links:  newLinks,
	}
}

func Get(source S.Source, maxTotalRequests, maxDepth uint64) *LimitedSource {
	var totalRequests uint64

	return &LimitedSource{
		source:           source,
		maxTotalRequests: maxTotalRequests,
		maxDepth:         maxDepth,
		totalRequestsPtr: &totalRequests,
		depth:            0,
	}
}
