package main

import (
	"os"
	"path/filepath"
	"strings"
)

// Session state lives in a temp directory (see cmdUI) and is shared with the
// hidden subcommands that fzf bindings call back into:
//
//	<session>/root      absolute path of the scanned root directory
//	<session>/expanded  newline-separated paths of expanded directories
//	<session>/mode      "tree" or "search"

func sessionRoot(dir string) string {
	b, _ := os.ReadFile(filepath.Join(dir, "root"))
	return strings.TrimSpace(string(b))
}

func loadExpanded(dir string) map[string]bool {
	m := map[string]bool{}
	b, _ := os.ReadFile(filepath.Join(dir, "expanded"))
	for _, line := range strings.Split(string(b), "\n") {
		if line != "" {
			m[line] = true
		}
	}
	return m
}

func saveExpanded(dir string, m map[string]bool) {
	lines := []string{}
	for k := range m {
		lines = append(lines, k)
	}
	os.WriteFile(filepath.Join(dir, "expanded"), []byte(strings.Join(lines, "\n")), 0o644)
}

func loadMode(dir string) string {
	b, _ := os.ReadFile(filepath.Join(dir, "mode"))
	return strings.TrimSpace(string(b))
}

func saveMode(dir, mode string) {
	os.WriteFile(filepath.Join(dir, "mode"), []byte(mode), 0o644)
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
