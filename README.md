# fzt — fzf-powered file tree explorer

[中文文档](README.zh-CN.md)

A file tree explorer built on fzf's UI: tree navigation by default, `/` for
fuzzy search, selected path printed to stdout.

## Build

```bash
make build        # or: go build -o fzt .
make install      # installs to /usr/local/bin (PREFIX=... to override)
```

## Usage

```bash
fzt [directory]   # defaults to the current directory
```

### Tree mode (default)

| Key | Action |
|-----|--------|
| `j`/`k`/arrow keys | Move cursor |
| `Enter` | Directory: collapse/expand · File: select and exit |
| `/` | Open the search box |
| `esc`/`ctrl-c` | Quit (exit code 130) |

### Search mode (after `/`)

- Type to fuzzy-search against full relative paths
- `Enter` on a file → output it; on a directory → jump back to tree mode
  with that directory expanded
- `esc` → back to tree mode

### Output

stdout carries only the selected path, so fzt composes in pipelines:

```bash
vim $(fzt)
fzt | xargs wc -l
```

## Shell integration (like fzf's Ctrl-T)

```bash
# bash: add to ~/.bashrc
source /path/to/fzt/key-bindings.bash

# zsh: add to ~/.zshrc
source /path/to/fzt/key-bindings.zsh
```

Press **Ctrl-F** on the command line to open the tree; the selected path is
inserted at the cursor.

## Testing

```bash
make test
```
