package http_source

import (
	"strings"
	"testing"
)

func TestParseHtml(t *testing.T) {
	input := `
	<a href="http://foo.com/link1">Foo!</a>
	<a href="http://otherdomain.com/link1">Foo!</a>
	<p>Blah blah blah</p>
	<link rel="stylesheet" href="foo.css">
	<img src="foo.jpeg">
	<img src="http://otherdomain.com/foo.jpeg">
	`

	outs := parseHtml(
		&HttpSource{
			url:          "http://foo.com",
			host:         "foo.com",
			protocol:     "http",
			errorHandler: func(httpUrl string, err error) { panic("Not expecting error") },
		},
		"/",
		strings.NewReader(input),
	)

	if len(outs.Links) != 1 || len(outs.Assets) != 2 {
		t.Errorf("Unexpected number of links/assets")
	}

	if outs.Links[0].Url != "http://foo.com/link1" || outs.Assets[0].Url != "http://foo.com/foo.css" || outs.Assets[1].Url != "http://foo.com/foo.jpeg" {
		t.Errorf("Unexpected links/assets")
	}

	t.Logf("%+v", outs)
}

// TestParseHtmlErrorHandling checks that the HTML parser extracts links
// correctly even in the presence of various HTML errors.
func TestParseHtmlErrorHandling(t *testing.T) {
	// the html below intentionally contains typos/errors
	input := `
	<a href="http://foo.com/link1">Foo!</aa>
	<a href="http://otherdomain.com/link1">Foo!</a>
	</pp>Blah blah blah</p>
	<linkk rel="stylesheet" href="foo.css">
	<img src="foo.jpeg"
	<img src="http://otherdomain.com/foo.jpeg">
	`

	outs := parseHtml(
		&HttpSource{
			url:          "http://foo.com",
			host:         "foo.com",
			protocol:     "http",
			errorHandler: func(httpUrl string, err error) { panic("Not expecting error") },
		},
		"/",
		strings.NewReader(input),
	)

	if len(outs.Links) != 1 || len(outs.Assets) != 2 {
		t.Errorf("Unexpected number of links/assets")
	}

	if outs.Links[0].Url != "http://foo.com/link1" || outs.Assets[0].Url != "http://foo.com/foo.css" || outs.Assets[1].Url != "http://foo.com/foo.jpeg" {
		t.Errorf("Unexpected links/assets")
	}

	t.Logf("%+v", outs)
}

func TestNormalizeUrl(t *testing.T) {
	type test struct {
		protocol            string
		host                string
		basePath            string
		url                 string
		expectedUrl         string
		expectedShouldCrawl bool
	}

	tests := []test{
		{"http", "foo.com", "", "index.html", "http://foo.com/index.html", true},
		{"http", "foo.com", "", "/bar", "http://foo.com/bar", true},
		{"http", "foo.com", "/", "index.html", "http://foo.com/index.html", true},
		{"http", "foo.com", "/", "/bar", "http://foo.com/bar", true},
		{"http", "foo.com", "/", "https://foo.com/amp", "https://foo.com/amp", true},
		{"http", "foo.com", "/", "https://anotherdomain.com/amp", "", false},
		{"http", "foo.com", "/", "ftp://foo.com/amp", "", false},
	}

	for _, tst := range tests {
		normalized, shouldCrawl := normalizeUrl(tst.protocol, tst.host, tst.basePath, tst.url)
		if shouldCrawl != tst.expectedShouldCrawl {
			t.Errorf("Expected (%v, %v, %v, %v) to have shouldCrawl=%v; got %v", tst.protocol, tst.host, tst.basePath, tst.url, tst.expectedShouldCrawl, shouldCrawl)

		}
		if normalized != tst.expectedUrl {
			t.Errorf("Expected (%v, %v, %v, %v) to normalized to %v; got %v", tst.protocol, tst.host, tst.basePath, tst.url, tst.expectedUrl, normalized)
		}
	}
}
