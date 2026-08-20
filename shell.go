package main

import _ "embed"

// Shell key bindings are embedded so the binary is self-contained:
// `eval "$(fzt --bash)"` / `eval "$(fzt --zsh)"` in the shell rc file.
// The .bash/.zsh files in the repo root remain the single source of truth.

//go:embed key-bindings.bash
var bashKeyBindings string

//go:embed key-bindings.zsh
var zshKeyBindings string
