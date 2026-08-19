# fzt — fzf-powered file tree explorer

[中文文档](README.zh-CN.md)

A file tree explorer built on fzf's UI: tree navigation by default, `/` for
fuzzy search, selected path printed to stdout.

<img width="1498" height="850" alt="image" src="https://github.com/user-attachments/assets/64d6bac3-3afd-4a9b-ad65-ee9b60795982" />


## Install

Download a prebuilt binary from
[Releases](https://github.com/r1cardohj/fzt/releases), or use `go install`.

### Linux

```bash
# amd64 (x86_64)
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.2/fzt-v0.0.2-linux-amd64.tar.gz
tar xzf fzt-v0.0.2-linux-amd64.tar.gz
sudo install -m755 fzt-v0.0.2-linux-amd64/fzt /usr/local/bin/

# arm64 (e.g. Raspberry Pi, ARM servers)
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.2/fzt-v0.0.2-linux-arm64.tar.gz
tar xzf fzt-v0.0.2-linux-arm64.tar.gz
sudo install -m755 fzt-v0.0.2-linux-arm64/fzt /usr/local/bin/
```

### macOS

```bash
# Apple Silicon (M1/M2/M3...)
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.2/fzt-v0.0.2-darwin-arm64.tar.gz
tar xzf fzt-v0.0.2-darwin-arm64.tar.gz
sudo install -m755 fzt-v0.0.2-darwin-arm64/fzt /usr/local/bin/

# Intel
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.2/fzt-v0.0.2-darwin-amd64.tar.gz
tar xzf fzt-v0.0.2-darwin-amd64.tar.gz
sudo install -m755 fzt-v0.0.2-darwin-amd64/fzt /usr/local/bin/
```

### With Go

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
