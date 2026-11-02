// Command revise archives the current version of a Hugo page and promotes a new
// version label. It is a writing-time tool: it never runs during `hugo`.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ifeitao/hugo-theme-resound/tools/internal/archive"
	"github.com/ifeitao/hugo-theme-resound/tools/internal/fm"
)

var version = "dev"

const usageText = `revise %s — version history for Hugo content

Usage:
  revise [flags] <content/posts/foo.md | content/posts/foo-bundle>
  revise undo <content/posts/foo.md | content/posts/foo-bundle>

Flags:
  --base-url string   the page's current URL, e.g. /posts/foo/
                      (needed when [permalinks] rewrites the path and Hugo is unavailable)
  --label string      label for the new version (default: today's date)

How it works:
  1. copies the current version to <name>.revisions/<label>/ (a sibling of the
     page, never inside the bundle itself)
  2. publishes it at <base-url>revisions/<label>/ with build.list=never
  3. writes the label -> URL list into front matter (revisions, revisions_schema)
  4. touches lastmod only — date is never rewritten, so permalinks stay stable

Commit before revising: revise undo <path> restores the previous state.
`

func usage() {
	fmt.Fprintf(os.Stderr, usageText+"\n", version)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	switch args[0] {
	case "-h", "--help", "help":
		usage()
		return
	case "-v", "--version", "version":
		fmt.Println("revise " + version)
		return
	case "undo":
		fs := flag.NewFlagSet("undo", flag.ExitOnError)
		fs.Usage = usage
		if err := fs.Parse(args[1:]); err != nil {
			os.Exit(2)
		}
		if fs.NArg() != 1 {
			usage()
			os.Exit(2)
		}
		res, err := archive.Undo(fs.Arg(0))
		if err != nil {
			fail(err)
		}
		fmt.Printf("undone    %s\n", res.Archived)
		fmt.Println("remaining versions:")
		printEntries(res.Entries)
		return
	}

	rest := args
	if args[0] == "revise" {
		rest = args[1:]
	}
	fs := flag.NewFlagSet("revise", flag.ExitOnError)
	fs.Usage = usage
	baseURL := fs.String("base-url", "", "the page's current URL")
	label := fs.String("label", "", "label for the new version")
	if err := fs.Parse(rest); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 1 {
		usage()
		os.Exit(2)
	}

	res, err := archive.Revise(fs.Arg(0), archive.Options{BaseURL: *baseURL, Label: *label})
	if err != nil {
		fail(err)
	}

	fmt.Printf("archived  %s  ->  %s\n", res.Label, res.Archived)
	fmt.Printf("live      %s  ->  %s   (URL from %s)\n", res.NewLabel, res.BaseURL, res.BaseSource)
	fmt.Println("versions:")
	printEntries(res.Entries)
	fmt.Printf("\nnext: git add -A && git commit -m \"archive %s\"\n", res.Label)
}

func printEntries(entries []fm.Revision) {
	for _, e := range entries {
		fmt.Printf("  %-12s  %s\n", e.Label, e.URL)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "revise: "+err.Error())
	os.Exit(1)
}
