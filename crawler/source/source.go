// Package source provides a generic representation of a loadable resource that
// references assets and contains links to other loadable resources.
package source

type Asset struct {
	Url string
}

type Link struct {
	Url    string
	Source Source
}

type Outs struct {
	Links  []Link
	Assets []Asset
}

type Source interface {
	GetUrl() string
	GetOuts() Outs
}
