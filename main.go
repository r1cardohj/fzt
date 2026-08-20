// fzt - an fzf-powered file tree explorer.
//
// fzt embeds fzf (github.com/junegunn/fzf/src) so that the entire UI is a
// single fzf instance: tree navigation by default, a fuzzy search box on "/",
// and the selected path printed to stdout.
//
// File layout:
//
//	main.go      entry point and subcommand dispatch
//	tree.go      directory scanning and candidate rendering
//	session.go   per-session state (expanded dirs, UI mode)
//	highlight.go syntax highlighting for the preview window
//	callbacks.go hidden subcommands invoked by fzf reload/transform bindings
//	ui.go        the embedded fzf UI
package main

import (
	"fmt"
	"os"
	"strings"

	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/protector"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: fzt [--height=HEIGHT] [directory]")
	os.Exit(fzf.ExitError)
}

func main() {
	protector.Protect()
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "__rows":
			cmdRows(flagValue(args, "--session"), hasFlag(args, "--search"))
			return
		case "__enter":
			line := ""
			for i, a := range args {
				if a == "--" && i+1 < len(args) {
					line = args[i+1]
				}
			}
			cmdEnter(flagValue(args, "--session"), line)
			return
		case "__esc":
			cmdEsc(flagValue(args, "--session"))
			return
		case "__preview":
			line := ""
			for i, a := range args {
				if a == "--" && i+1 < len(args) {
					line = args[i+1]
				}
			}
			cmdPreview(flagValue(args, "--session"), line)
			return
		case "__mode":
			saveMode(flagValue(args, "--session"), args[len(args)-1])
			return
		case "-h", "--help":
			usage()
		}
	}
	root := "."
	var extraOpts []string
	for _, a := range args {
		if strings.HasPrefix(a, "--height") {
			extraOpts = append(extraOpts, a) // passed through to fzf
		} else {
			root = a
		}
	}
	os.Exit(cmdUI(root, extraOpts))
}

func flagValue(args []string, name string) string {
	for i, a := range args {
		if a == name && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func hasFlag(args []string, name string) bool {
	for _, a := range args {
		if a == name {
			return true
		}
	}
	return false
}
