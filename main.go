package main

import (
	"flag"
	"fmt"
	"os"

	C "multiverse.io/crawler/crawler"
	H "multiverse.io/crawler/crawler/http_source"
	L "multiverse.io/crawler/crawler/limited_source"
	R "multiverse.io/crawler/crawler/render"
)

func main() {
	args := getCommandArgs()

	source, err := H.Get(args.url, handleError)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	limitedSource := L.Get(&source, args.nRequestsLimit, args.depthLimit)

	var assetsMode C.AssetsMode
	if args.noAssets {
		assetsMode = C.AssetsModeIgnoreAssets
	} else {
		assetsMode = C.AssetsModeIncludeAssets
	}

	root := C.Crawl(limitedSource, assetsMode)

	html := R.ExportHtml(root)
	fmt.Printf("%v\n", html)
}

func handleError(httpUrl string, err error) {
	fmt.Fprintf(os.Stderr, "Error for %v: %v\n", httpUrl, err)
}

type CommandArgs struct {
	url            string
	depthLimit     uint64
	nRequestsLimit uint64
	noAssets       bool
}

const defaultDepthLimit = 30
const defaultNRequestsLimit = 200

func getCommandArgs() (args CommandArgs) {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)

	flagSet.Uint64Var(&args.depthLimit, "maxdepth", defaultDepthLimit, "the maximum depth of the traversal from the root")
	flagSet.Uint64Var(&args.nRequestsLimit, "maxreqs", defaultNRequestsLimit, "the maximum number of HTTP requests to make before halting")
	flagSet.BoolVar(&args.noAssets, "noassets", false, "if this flag is present, assets are not included in the graph")

	flagSet.Parse(os.Args[1:])

	args.url = flagSet.Arg(0)
	if args.url == "" {
		fmt.Fprintf(os.Stderr, "You must provide a URL.\n")
		os.Exit(1)
	}

	return
}
