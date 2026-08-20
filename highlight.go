package main

import (
	"bufio"
	"os"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Syntax highlighting for the preview window. The terminal16 formatter maps
// the style's colors onto the terminal's own 16-color palette, so
// highlighting follows the user's terminal theme (light or dark) and stays
// consistent with the fzf UI instead of forcing a fixed RGB theme.
var (
	previewFormatter = formatters.Get("terminal16")
	previewStyle     = styles.Get("monokai")
)

// highlight writes data to stdout with syntax colors, choosing a lexer by
// file name first and falling back to content analysis (shebangs etc.).
// Returns false when no lexer matches or rendering fails, so the caller can
// fall back to the plain dump.
func highlight(name string, data []byte) bool {
	lexer := lexers.Match(name)
	if lexer == nil {
		sniff := data
		if len(sniff) > 8000 {
			sniff = sniff[:8000]
		}
		lexer = lexers.Analyse(string(sniff))
	}
	if lexer == nil {
		return false
	}
	it, err := chroma.Coalesce(lexer).Tokenise(nil, string(data))
	if err != nil {
		return false
	}
	w := bufio.NewWriter(os.Stdout)
	if err := previewFormatter.Format(w, previewStyle, it); err != nil {
		return false
	}
	w.Flush()
	return true
}
