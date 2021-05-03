package main

import (
	"strings"
	"testing"
)

func TestGetCommandArgs(t *testing.T) {
	{
		var usageOutput strings.Builder
		args, err := getCommandArgs(&usageOutput, []string{"http://foo.com"})
		if err != nil || args.depthLimit != defaultDepthLimit || args.nRequestsLimit != defaultNRequestsLimit || args.noAssets || args.url != "http://foo.com" {
			t.Errorf("Unexpected command line parse result for simple URL argument.\n")
		}
	}

	{
		var usageOutput strings.Builder
		_, err := getCommandArgs(&usageOutput, []string{"http://foo.com", "http://foo.com"})
		if err == nil {
			t.Errorf("Expected error if two URLs given.\n")
		}
	}

	{
		var usageOutput strings.Builder
		args, err := getCommandArgs(&usageOutput, []string{"-maxdepth", "999", "http://foo.com"})
		if err != nil || args.depthLimit != 999 || args.nRequestsLimit != defaultNRequestsLimit || args.noAssets || args.url != "http://foo.com" {
			t.Errorf("Couldn't set -maxdepth 999.\n")
		}
	}

	{
		var usageOutput strings.Builder
		args, err := getCommandArgs(&usageOutput, []string{"-maxreqs", "999", "http://foo.com"})
		if err != nil || args.depthLimit != defaultDepthLimit || args.nRequestsLimit != 999 || args.noAssets || args.url != "http://foo.com" {
			t.Errorf("Couldn't set -maxreqs 999.\n")
		}
	}

	{
		var usageOutput strings.Builder
		args, err := getCommandArgs(&usageOutput, []string{"-noassets", "http://foo.com"})
		if err != nil || args.depthLimit != defaultDepthLimit || args.nRequestsLimit != defaultNRequestsLimit || !args.noAssets || args.url != "http://foo.com" {
			t.Errorf("Couldn't set -noassets.\n")
		}
	}

	{
		var usageOutput strings.Builder
		_, err := getCommandArgs(&usageOutput, []string{"http://foo.com", "-noassets"})
		if err == nil {
			t.Errorf("Expected error if trying to set flag following URL.\n")
		}
	}
}
