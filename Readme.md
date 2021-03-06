# GoCrawl

A simple web crawler written in Go that outputs an HTML page with an interactive graph.

**This is a toy project written for a coding challenge.**


## Basic usage

Crawl a web page:

```sh
go run main.go http://example.com > example.com.html
```

The program prints the HTML page to stdout and a request log to stderr.

## Domain restriction

Only pages on the original domain (or subdomains thereof) will be crawled. Links
to other domains are not shown in the graph.

## Graph format

The root node is shown in red. Nodes that are referenced via `<a>` elements are
shown in blue. Nodes that are referenced only by elements other than `<a>` are
regarded as "assets" and shown in grey. Assets are not retrieved (and hence not
parsed for further pages to crawl). Only documents sent with a content type of
`text/html` are parsed for links.

You can click and drag nodes in the graph to modify the layout.


## Command line options

**Go's command line parser requires flags to come before the URL.**

Option     | Default | Description                                                    |
---------- | ------- | -------------------------------------------------------------- |
-maxdepth  | 30      | The maximum depth of the traversal from the root.              |
-maxreqs   | 200     | The maximum number of HTTP requests to make before halting.    |
-noassets  |         | If this flag is present, assets are not included in the graph. |

Usage:

```sh
go run main.go [-maxdepth INT] [-maxreqs INT] [-noassets] <URL>
```

## Development notes

Run Go tests (Go >= 1.16):

```sh
go test ./...
```

Run JS tests (Node.js >= 15.0):

```sh
(cd crawler/render && npm test)
```
