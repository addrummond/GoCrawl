package http_source

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	H "golang.org/x/net/html"
	S "multiverse.io/crawler/crawler/source"
)

type StatusCodeError struct {
	statusCode string
}

func (e *StatusCodeError) Error() string {
	return fmt.Sprintf("Status code %v", e.statusCode)
}

type BadProtocolError struct {
	protocol string
}

func (e *BadProtocolError) Error() string {
	return fmt.Sprintf("Protocol '%v' not supported", e.protocol)
}

type HttpSource struct {
	url          string
	host         string // only load from this domain
	protocol     string // use this protocol by default for relative links
	path         string
	contents     *string
	errorHandler func(url string, err error)
}

func Get(u string, errorHandler func(errorUrl string, err error)) (HttpSource, error) {
	var s HttpSource

	parsed, err := url.Parse(u)
	if err != nil {
		return s, err
	}

	if !supportedProtocol(parsed.Scheme) {
		return s, &BadProtocolError{protocol: parsed.Scheme}
	}

	s.url = u
	s.host = strings.ToLower(parsed.Host)
	s.protocol = strings.ToLower(parsed.Scheme)
	s.path = parsed.Path
	s.errorHandler = errorHandler
	return s, nil
}

func (s *HttpSource) GetUrl() string {
	return s.url
}

func (s *HttpSource) GetOuts() (outs S.Outs) {
	fmt.Fprintf(os.Stderr, "GET %v\n", s.url)

	resp, err := http.Get(s.url)
	if err != nil {
		s.errorHandler(s.url, err)
		return
	}

	if resp.StatusCode != 200 {
		s.errorHandler(s.url, &StatusCodeError{statusCode: resp.Status})
	}

	contentType := resp.Header["Content-Type"]
	if len(contentType) > 0 && strings.HasPrefix(contentType[0], "text/html") {
		// URL may have changed due to redirect. If we got redirected to a new
		// domain, ignore it.
		finalUrl := resp.Request.URL.String()
		parsed, err := url.Parse(finalUrl)
		if err == nil && hostMatches(s.host, parsed.Host) {
			outs = parseHtml(s, parsed.Path, resp.Body)
		}
	}

	return outs
}

func parseHtml(s *HttpSource, newPath string, reader io.Reader) (outs S.Outs) {
	z := H.NewTokenizer(reader)

	existingLinks := make(map[string]bool)
	existingAssets := make(map[string]bool)

	for {
		tt := z.Next()
		if tt == H.ErrorToken && z.Err() == io.EOF {
			break
		}

		for {
			nameBytes, val, more := z.TagAttr()
			name := string(nameBytes)

			if name == "href" || name == "src" {
				url, ok := normalizeUrl(s.protocol, s.host, newPath, string(val))
				if ok {
					newSource, err := Get(url, s.errorHandler)
					if err != nil {
						s.errorHandler(url, err)
					} else {
						tagNameBytes, _ := z.TagName()
						tagName := string(tagNameBytes)

						if _, already := existingLinks[url]; !already && tagName == "a" {
							existingLinks[url] = true
							outs.Links = append(outs.Links, S.Link{Url: url, Source: &newSource})
						} else if _, already := existingAssets[url]; !already {
							existingAssets[url] = true
							outs.Assets = append(outs.Assets, S.Asset{Url: url})
						}
					}
				}
			}

			if !more {
				break
			}
		}
	}

	return outs
}

func normalizeUrl(protocol, host, basePath, httpUrl string) (string, bool) {
	parsed, err := url.Parse(httpUrl)
	if err != nil {
		return "", false
	}

	if parsed.Scheme != "" && !supportedProtocol(parsed.Scheme) {
		return "", false
	}

	if parsed.IsAbs() && strings.ToLower(parsed.Host) != host {
		return "", false
	}

	if !parsed.IsAbs() {
		parsed.Scheme = protocol
		parsed.Host = host
		if len(parsed.Path) > 0 && parsed.Path[0] != '/' {
			parsed.Path = path.Join(path.Dir(basePath), parsed.Path)
		}
	}

	parsed.Fragment = ""

	return parsed.String(), true
}

func supportedProtocol(protocol string) bool {
	return protocol == "http" || protocol == "https"
}

func hostMatches(hostWeAreCrawling, host string) bool {
	hostWeAreCrawlingLc := strings.ToLower(hostWeAreCrawling)
	hostLc := strings.ToLower(host)
	return hostWeAreCrawlingLc == hostLc || strings.HasSuffix(hostLc, "."+hostWeAreCrawlingLc)
}
