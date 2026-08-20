package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// previewCmd is the shell command fzf runs to render the preview window for
// the current candidate ({} is fzf's placeholder for the selected line).
func previewCmd(sessionDir string) string {
	exe, _ := os.Executable()
	return shellQuote(exe) + " __preview --session " + shellQuote(sessionDir)
}

// cmdRows prints candidate lines. Tree mode lazily scans only the expanded
// directories; search mode needs every path as a candidate, so it falls
// back to a full fastwalk scan.
func cmdRows(sessionDir string, search bool) {
	root := sessionRoot(sessionDir)
	if search {
		for _, line := range pathRows(buildTree(root)) {
			fmt.Println(line)
		}
		return
	}
	expanded := loadExpanded(sessionDir)
	for _, line := range rows(buildTreeLazy(root, expanded), expanded) {
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

// Preview limits: files are capped in both bytes and lines so opening a huge
// minified file or log doesn't block the UI.
const (
	previewMaxBytes   = 256 * 1024
	previewMaxLines   = 1000
	previewMaxEntries = 200 // directory listings
)

// cmdPreview is the fzf --preview callback. Directories get a listing (dirs
// first, like the tree view); files get their contents, with a binary sniff
// so previewing an image doesn't spew garbage into the terminal.
func cmdPreview(sessionDir, line string) {
	root := sessionRoot(sessionDir)
	rel := strings.SplitN(line, "\t", 2)[0]
	p := filepath.Join(root, rel)
	info, err := os.Stat(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	if info.IsDir() {
		previewDir(p)
		return
	}
	previewFile(p, info)
}

func previewDir(p string) {
	entries, err := os.ReadDir(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})
	count := 0
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") { // same hidden-file rule as the tree
			continue
		}
		if count >= previewMaxEntries {
			fmt.Println("\u2026")
			return
		}
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		fmt.Println(name)
		count++
	}
	if count == 0 {
		fmt.Println("(empty directory)")
	}
}

func previewFile(p string, info os.FileInfo) {
	f, err := os.Open(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	buf := make([]byte, previewMaxBytes)
	n, _ := f.Read(buf)
	data := buf[:n]
	// Binary sniff: NUL in the first chunk is the same heuristic grep uses.
	sniff := data
	if len(sniff) > 8000 {
		sniff = sniff[:8000]
	}
	if bytes.IndexByte(sniff, 0) >= 0 {
		fmt.Printf("binary file (%d bytes)\n", info.Size())
		return
	}
	lines := strings.SplitAfter(string(data), "\n")
	for i, l := range lines {
		if i >= previewMaxLines {
			fmt.Println("\u2026")
			return
		}
		fmt.Print(l)
	}
	if n == previewMaxBytes {
		fmt.Println("\u2026")
	}
}
