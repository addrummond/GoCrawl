# GoCrawl

A simple web crawler that outputs an HTML page with an interactive graph.


## Basic usage

Crawl a web page:

```sh
go run main.go http://example.com > example.com.html
```

The program outputs the HTML page to stdout and a request log to stderr.

## Graph format

The root node is shown in red. Nodes that are referenced via `<a>` elements are
shown in blue. Nodes that are referenced only by elements other than `<a>` are
regarded as "assets" and shown in grey. Assets are not retrieved (and hence not
parsed for further pages to crawl).

You can click and drag nodes in the graph to modify the layout.


## Command line options

Option     | Default | Description                                                    |
---------- | ------- | -------------------------------------------------------------- |
-maxdepth  | 30      | The maximum depth of the traversal from the root.              |
-maxreqs   | 200     | The maximum number of HTTP requests to make before halting.    |
-noassets  |         | If this flag is present, assets are not included in the graph. |

**Go's command line parser requires flags to come before the URL.**

Usage:

```sh
go run main.go http://example.com [-maxdepth INT] [-maxreqs INT] [-noassets] <URL>
```

## Development notes

Run Go tests:

```sh
go test ./...
```

Run JS tests:

```sh
cd crawler/render && npm test
```
