package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/charlievieth/fastwalk"
)

// Node is one entry in the directory tree.
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
}

// buildTree scans root concurrently with fastwalk (the same walker fzf uses
// internally), skipping hidden files, and returns the root node. Children are
// sorted: directories first, then files.
func buildTree(root string) *Node {
	rootNode := &Node{Name: filepath.Base(root), Path: root, IsDir: true}
	type entry struct {
		path  string
		isDir bool
	}
	var (
		entries []entry
		mu      sync.Mutex
	)
	conf := fastwalk.Config{Follow: false}
	fastwalk.Walk(&conf, root, func(path string, de os.DirEntry, err error) error {
		if err != nil || path == root { // fastwalk emits the root itself too
			return nil
		}
		if strings.HasPrefix(de.Name(), ".") { // skip hidden files
			if de.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		mu.Lock()
		entries = append(entries, entry{path, de.IsDir()})
		mu.Unlock()
		return nil
	})

	// Assemble the tree: parents before children (sorted by depth).
	sort.Slice(entries, func(i, j int) bool {
		return strings.Count(entries[i].path, string(filepath.Separator)) <
			strings.Count(entries[j].path, string(filepath.Separator))
	})
	nodes := map[string]*Node{root: rootNode}
	for _, e := range entries {
		n := &Node{Name: filepath.Base(e.path), Path: e.path, IsDir: e.isDir}
		nodes[e.path] = n
		if parent, ok := nodes[filepath.Dir(e.path)]; ok {
			parent.Children = append(parent.Children, n)
		}
	}
	for _, n := range nodes {
		sort.Slice(n.Children, func(i, j int) bool {
			a, b := n.Children[i], n.Children[j]
			if a.IsDir != b.IsDir {
				return a.IsDir // directories first
			}
			return a.Name < b.Name
		})
	}
	return rootNode
}

// relpath returns n's path relative to root, falling back to the full path.
func relpath(root *Node, n *Node) string {
	rel, err := filepath.Rel(root.Path, n.Path)
	if err != nil {
		return n.Path
	}
	return rel
}

// rows renders candidate lines "relpath<TAB>tree display" for every node
// except the root. Collapsed dirs hide their children. Directories carry a
// ▾/▸ marker; indentation alone conveys nesting (no ├── connectors).
func rows(root *Node, expanded map[string]bool) []string {
	out := []string{}
	var walk func(n *Node, depth int)
	walk = func(n *Node, depth int) {
		var b strings.Builder
		b.WriteString(strings.Repeat("  ", depth))
		if n.IsDir {
			if expanded[n.Path] {
				b.WriteString("▾ ")
			} else {
				b.WriteString("▸ ")
			}
		} else {
			b.WriteString("  ")
		}
		b.WriteString(n.Name)
		out = append(out, relpath(root, n)+"\t"+b.String())
		if n.IsDir && expanded[n.Path] {
			for _, c := range n.Children {
				walk(c, depth+1)
			}
		}
	}
	for _, c := range root.Children {
		walk(c, 0)
	}
	return out
}

// pathRows renders every node (except root) as a flat, fully-qualified path
// line for search mode: "relpath<TAB>styled relpath" where the directory
// part is dimmed and directories are colored, so nesting stays readable.
func pathRows(root *Node) []string {
	out := []string{}
	var walk func(n *Node)
	walk = func(n *Node) {
		rel := relpath(root, n)
		dir, base := filepath.Split(rel)
		var disp string
		switch {
		case n.IsDir:
			disp = "\x1b[34m" + rel + "/\x1b[0m"
		case dir != "":
			disp = "\x1b[2m" + dir + "\x1b[0m" + base
		default:
			disp = base
		}
		out = append(out, rel+"\t"+disp)
		for _, c := range n.Children {
			walk(c)
		}
	}
	for _, c := range root.Children {
		walk(c)
	}
	return out
}
