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
//	callbacks.go hidden subcommands invoked by fzf reload/transform bindings
//	ui.go        the embedded fzf UI
package main

import (
	"fmt"
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/protector"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: fzt [directory]")
	os.Exit(fzf.ExitError)
}

func main() {
	protector.Protect()
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "__rows":
			cmdRows(flagValue(args, "--session"))
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
		case "__mode":
			saveMode(flagValue(args, "--session"), args[len(args)-1])
			return
		case "-h", "--help":
			usage()
		}
	}
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	os.Exit(cmdUI(root))
}

func flagValue(args []string, name string) string {
	for i, a := range args {
		if a == name && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
