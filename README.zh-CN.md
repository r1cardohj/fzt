# fzt — fzf 驱动的文件树浏览器

[English](README.md)

基于 fzf UI 的文件树浏览器：默认树形导航，`/` 模糊搜索，选中的路径输出到 stdout。

![demo](docs/demo.gif)

## 安装

```bash
curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | sh
```

安装最新版本到 `~/.local/bin`（无需 sudo，并会校验 release 的
checksums）。固定版本或系统级安装：

```bash
curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | sh -s -- v0.0.5
curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | PREFIX=/usr/local sh
```

Linux/macOS/Windows 的预编译归档也可以在
[Releases](https://github.com/r1cardohj/fzt/releases) 页面手动下载。

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
| `ctrl-/` | 开关预览窗口 |
| `esc`/`ctrl-c` | 退出（退出码 130） |

预览窗口使用 fzf 原生的 `--preview`，显示在右侧：目录显示内容列表，
文件显示内容（二进制文件会被识别，不会直接输出乱码）。文件预览通过
[chroma](https://github.com/alecthomas/chroma) 做语法高亮，颜色映射到
终端的 16 色调色板，随终端主题自动变化。

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
eval "$(fzt --bash)"

# zsh: 加到 ~/.zshrc
eval "$(fzt --zsh)"
```

命令行按 **Ctrl-F** 打开文件树，选中的路径插入到光标处。

## 测试

```bash
make test
```

## 许可证

[MIT](LICENSE)
