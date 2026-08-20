# fzt — fzf-powered file tree explorer

[中文文档](README.zh-CN.md)

A file tree explorer built on fzf's UI: tree navigation by default, `/` for
fuzzy search, selected path printed to stdout.

![demo](docs/demo.gif)


## Install

```bash
curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | sh
```

Installs the latest release to `~/.local/bin` (no sudo needed; the archive
is verified against the release checksums). To pin a version or install
system-wide:

```bash
curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | sh -s -- v0.0.5
curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | PREFIX=/usr/local sh
```

Prebuilt archives for Linux/macOS/Windows are also on the
[Releases](https://github.com/r1cardohj/fzt/releases) page.

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
| `ctrl-/` | Toggle the preview window |
| `esc`/`ctrl-c` | Quit (exit code 130) |

The preview window (fzf's native `--preview`, right side) shows a directory
listing for directories and file contents for files (binary files are
detected and not dumped). File previews are syntax-highlighted with
[chroma](https://github.com/alecthomas/chroma) mapped onto the terminal's
16-color palette, so colors follow your terminal theme.

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
eval "$(fzt --bash)"

# zsh: add to ~/.zshrc
eval "$(fzt --zsh)"
```

Press **Ctrl-F** on the command line to open the tree; the selected path is
inserted at the cursor.

## Testing

```bash
make test
```

## License

[MIT](LICENSE)
