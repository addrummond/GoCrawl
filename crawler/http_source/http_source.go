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
	errorHandler func(url string, err error)
}

func MakeSource(u string, errorHandler func(errorUrl string, err error)) (HttpSource, error) {
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
		return outs
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

	tokenize(z, func(tagName string, attributes map[string]string) {
		url, ok := getUrlFromTagAttributes(attributes)
		if !ok {
			return
		}

		normalizedUrl, ok := normalizeUrl(s.protocol, s.host, newPath, url)
		if !ok {
			return
		}

		newSource, err := MakeSource(normalizedUrl, s.errorHandler)
		if err != nil {
			s.errorHandler(normalizedUrl, err)
			return
		}

		if !existingLinks[normalizedUrl] && tagName == "a" {
			existingLinks[normalizedUrl] = true
			outs.Links = append(outs.Links, S.Link{Url: normalizedUrl, Source: &newSource})
		} else if !existingAssets[normalizedUrl] {
			existingAssets[normalizedUrl] = true
			outs.Assets = append(outs.Assets, S.Asset{Url: normalizedUrl})
		}
	})

	return outs
}

func tokenize(z *H.Tokenizer, f func(tagName string, attributes map[string]string)) {
	for {
		tt := z.Next()
		if tt == H.ErrorToken && z.Err() == io.EOF {
			break
		}

		tagNameBytes, _ := z.TagName()
		tagName := string(tagNameBytes)

		attributes := make(map[string]string)
		for {
			nameBytes, valBytes, more := z.TagAttr()
			nameString := string(nameBytes)

			if nameString == "href" || nameString == "src" {
				attributes[nameString] = string(valBytes)
			}

			if !more {
				break
			}
		}

		f(tagName, attributes)
	}
}

func getUrlFromTagAttributes(attributes map[string]string) (string, bool) {
	if href := attributes["href"]; href != "" {
		return href, true
	}
	if src := attributes["src"]; src != "" {
		return src, true
	}
	return "", false
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
