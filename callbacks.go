package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Hidden subcommands invoked by fzf reload/transform bindings. fzf exports
// FZF_QUERY to these child processes, which is how they tell search mode
// (non-empty query) apart from tree mode.

// rowsCmd is the shell command that prints the tree-view candidate list.
func rowsCmd(sessionDir string) string {
	exe, _ := os.Executable()
	return shellQuote(exe) + " __rows --session " + shellQuote(sessionDir)
}

// searchRowsCmd is the shell command that prints the flat path list used in
// search mode.
func searchRowsCmd(sessionDir string) string {
	exe, _ := os.Executable()
	return shellQuote(exe) + " __rows --session " + shellQuote(sessionDir) + " --search"
}

// cmdRows prints candidate lines: the indented tree view honoring collapsed
// dirs, or (with --search) the flat full-path list for readable search
// results.
func cmdRows(sessionDir string, search bool) {
	root := buildTree(sessionRoot(sessionDir))
	if search {
		for _, line := range pathRows(root) {
			fmt.Println(line)
		}
		return
	}
	for _, line := range rows(root, loadExpanded(sessionDir)) {
		fmt.Println(line)
	}
}

// treeModeActions is the fzf action sequence that switches the UI back to
// tree-navigation mode.
func treeModeActions(sessionDir string) string {
	return "clear-query+hide-input+disable-search+rebind(j,k)+change-prompt(fzt> )" +
		"+reload-sync(" + rowsCmd(sessionDir) + ")"
}

// cmdEnter is the enter:transform callback. Files -> accept; dirs -> toggle
// (and, when invoked from search mode, drop back into tree mode).
func cmdEnter(sessionDir, line string) {
	root := sessionRoot(sessionDir)
	rel := strings.SplitN(line, "\t", 2)[0]
	p := filepath.Join(root, rel)
	info, err := os.Stat(p)
	if err == nil && info.IsDir() {
		expanded := loadExpanded(sessionDir)
		if expanded[p] {
			delete(expanded, p)
		} else {
			expanded[p] = true
		}
		saveExpanded(sessionDir, expanded)
		if loadMode(sessionDir) == "search" {
			saveMode(sessionDir, "tree")
			fmt.Print(treeModeActions(sessionDir))
		} else {
			fmt.Printf("reload-sync(%s)", rowsCmd(sessionDir))
		}
	} else {
		fmt.Println("accept")
	}
}

// cmdEsc is the esc:transform callback: search mode -> back to the tree;
// tree mode -> quit.
func cmdEsc(sessionDir string) {
	if loadMode(sessionDir) == "search" {
		saveMode(sessionDir, "tree")
		fmt.Print(treeModeActions(sessionDir))
	} else {
		fmt.Println("abort")
	}
}
