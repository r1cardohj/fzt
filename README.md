# fzt — fzf-powered file tree explorer

[中文文档](README.zh-CN.md)

A file tree explorer built on fzf's UI: tree navigation by default, `/` for
fuzzy search, selected path printed to stdout.

## Install

Download a prebuilt binary from
[Releases](https://github.com/r1cardohj/fzt/releases) (Linux / macOS / Windows):

```bash
# e.g. Linux amd64
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.1/fzt-v0.0.1-linux-amd64.tar.gz
tar xzf fzt-v0.0.1-linux-amd64.tar.gz
sudo install -m755 fzt-v0.0.1-linux-amd64/fzt /usr/local/bin/
```

Or with Go:

```bash
go install github.com/r1cardohj/fzt@latest
```

## Build from source

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
