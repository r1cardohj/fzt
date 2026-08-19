package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fzf "github.com/junegunn/fzf/src"
)

// cmdUI runs the embedded fzf UI and prints the selected path to stdout.
func cmdUI(rootArg string) int {
	abs, err := filepath.Abs(rootArg)
	if err != nil || !isDir(abs) {
		fmt.Fprintf(os.Stderr, "fzt: %s is not a directory\n", rootArg)
		return fzf.ExitError
	}
	sessionDir, err := os.MkdirTemp("", "fzt-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "fzt:", err)
		return fzf.ExitError
	}
	defer os.RemoveAll(sessionDir)
	os.WriteFile(filepath.Join(sessionDir, "root"), []byte(abs), 0o644)
	saveExpanded(sessionDir, map[string]bool{abs: true})
	saveMode(sessionDir, "tree")

	exe, _ := os.Executable()
	searchMode := "execute-silent(" + shellQuote(exe) + " __mode --session " + shellQuote(sessionDir) + " search)" +
		"+show-input+enable-search+unbind(j,k)+change-prompt(search> )"
	opts, err := fzf.ParseOptions(true, []string{
		"--delimiter", "\t",
		"--nth", "1", // search the path field
		"--with-nth", "2", // display the tree/path line
		"--scheme", "path",
		"--no-sort", // keep tree order
		"--track",
		"--ansi",
		"--disabled", // tree mode: pure navigation, no query input
		"--layout", "reverse",
		"--prompt", "fzt> ",
		"--header", abs + "\nj/k: move · enter: toggle dir / select file · /: search · esc: quit",
		"--bind", "start:hide-input",
		"--bind", "j:down",
		"--bind", "k:up",
		"--bind", "/:" + searchMode,
		"--bind", "change:reload-sync(" + rowsCmd(sessionDir) + ")",
		"--bind", "enter:transform(" + shellQuote(exe) + " __enter --session " + shellQuote(sessionDir) + " -- {})",
		"--bind", "esc:transform(" + shellQuote(exe) + " __esc --session " + shellQuote(sessionDir) + ")",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "fzt:", err)
		return fzf.ExitError
	}

	// Initial list via channel; subsequent updates happen through reload.
	initial := rows(buildTree(abs), map[string]bool{abs: true})
	opts.Input = make(chan string, len(initial))
	for _, line := range initial {
		opts.Input <- line
	}
	close(opts.Input)

	accepted := []string{}
	opts.Output = make(chan string, 16)
	go func() {
		for s := range opts.Output {
			accepted = append(accepted, s)
		}
	}()

	code, err := fzf.Run(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fzt:", err)
	}
	if code == fzf.ExitOk && len(accepted) > 0 {
		rel := strings.SplitN(accepted[0], "\t", 2)[0]
		fmt.Println(filepath.Clean(filepath.Join(rootArg, rel)))
	}
	return code
}

func isDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}
