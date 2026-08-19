package main

import (
	"os"
	"path/filepath"
	"strings"
)

// Node is one entry in the directory tree.
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
}

// buildTree scans root recursively (skipping hidden files) and returns the
// root node. Children are sorted: directories first, then files.
func buildTree(root string) *Node {
	rootNode := &Node{Name: filepath.Base(root), Path: root, IsDir: true}
	stack := []*Node{rootNode}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		entries, err := os.ReadDir(n.Path) // already sorted by name
		if err != nil {
			continue
		}
		var dirs, files []*Node
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), ".") { // skip hidden files
				continue
			}
			c := &Node{Name: e.Name(), Path: filepath.Join(n.Path, e.Name()), IsDir: e.IsDir()}
			if c.IsDir {
				dirs = append(dirs, c)
			} else {
				files = append(files, c)
			}
		}
		n.Children = append(dirs, files...)
		for i := len(dirs) - 1; i >= 0; i-- {
			stack = append(stack, dirs[i])
		}
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
