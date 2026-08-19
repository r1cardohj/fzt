# fzt — fzf 驱动的文件树浏览器

[English](README.md)

基于 fzf UI 的文件树浏览器：默认树形导航，`/` 模糊搜索，选中的路径输出到 stdout。

## 安装

从 [Releases](https://github.com/r1cardohj/fzt/releases) 下载预编译二进制，或用 `go install`。

### Linux

```bash
# amd64 (x86_64)
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.1/fzt-v0.0.1-linux-amd64.tar.gz
tar xzf fzt-v0.0.1-linux-amd64.tar.gz
sudo install -m755 fzt-v0.0.1-linux-amd64/fzt /usr/local/bin/

# arm64（树莓派、ARM 服务器等）
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.1/fzt-v0.0.1-linux-arm64.tar.gz
tar xzf fzt-v0.0.1-linux-arm64.tar.gz
sudo install -m755 fzt-v0.0.1-linux-arm64/fzt /usr/local/bin/
```

### macOS

```bash
# Apple Silicon（M1/M2/M3……）
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.1/fzt-v0.0.1-darwin-arm64.tar.gz
tar xzf fzt-v0.0.1-darwin-arm64.tar.gz
sudo install -m755 fzt-v0.0.1-darwin-arm64/fzt /usr/local/bin/

# Intel
curl -LO https://github.com/r1cardohj/fzt/releases/download/v0.0.1/fzt-v0.0.1-darwin-amd64.tar.gz
tar xzf fzt-v0.0.1-darwin-amd64.tar.gz
sudo install -m755 fzt-v0.0.1-darwin-amd64/fzt /usr/local/bin/
```

### 用 Go 安装

```bash
go install github.com/r1cardohj/fzt@latest
```

## 从源码构建

```bash
make build        # 或者: go build -o fzt .
make install      # 安装到 /usr/local/bin（可用 PREFIX=... 覆盖）
```

## 用法

```bash
fzt [目录]   # 默认当前目录
```

### 树模式（默认）

| 按键 | 动作 |
|------|------|
| `j`/`k`/方向键 | 移动光标 |
| `Enter` | 目录：折叠/展开 · 文件：选中并退出 |
| `/` | 唤出搜索框 |
| `esc`/`ctrl-c` | 退出（退出码 130） |

### 搜索模式（按 `/` 后）

- 输入即对完整相对路径做模糊搜索
- `Enter` 文件 → 输出；目录 → 跳回树模式并展开定位
- `esc` → 返回树模式

### 输出

stdout 只包含选中的路径，可管道复用：

```bash
vim $(fzt)
fzt | xargs wc -l
```

## Shell 集成（类似 fzf 的 Ctrl-T）

```bash
# bash: 加到 ~/.bashrc
source /path/to/fzt/key-bindings.bash

# zsh: 加到 ~/.zshrc
source /path/to/fzt/key-bindings.zsh
```

命令行按 **Ctrl-F** 打开文件树，选中的路径插入到光标处。

## 测试

```bash
make test
```
