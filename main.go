package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	C "multiverse.io/crawler/crawler"
	H "multiverse.io/crawler/crawler/http_source"
	L "multiverse.io/crawler/crawler/limited_source"
	R "multiverse.io/crawler/crawler/render"
)

func main() {
	args, err := getCommandArgs(os.Stderr, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

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

func getCommandArgs(usageOutput io.Writer, argv []string) (args CommandArgs, err error) {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.SetOutput(usageOutput)

	flagSet.Uint64Var(&args.depthLimit, "maxdepth", defaultDepthLimit, "the maximum depth of the traversal from the root")
	flagSet.Uint64Var(&args.nRequestsLimit, "maxreqs", defaultNRequestsLimit, "the maximum number of HTTP requests to make before halting")
	flagSet.BoolVar(&args.noAssets, "noassets", false, "if this flag is present, assets are not included in the graph")

	if err = flagSet.Parse(argv); err != nil {
		return
	}

	if flagSet.NArg() != 1 {
		err = errors.New("You must provide exactly one URL.\n")
		fmt.Fprintf(usageOutput, "%v", err)
		return
	}

	args.url = flagSet.Arg(0)

	return
}
