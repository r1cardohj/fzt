# fzt key bindings for zsh —— 在 ~/.zshrc 里:
#   source /path/to/fzt/key-bindings.zsh
#
# Ctrl-F: 打开 fzt 文件树，选中的路径插入到当前命令行光标处

fzt-widget() {
  local selected
  selected="$(command fzt 2>/dev/null)" || { zle reset-prompt; return 1; }
  LBUFFER+="${(q)selected}"
  zle reset-prompt
  return 0
}
zle -N fzt-widget
bindkey '^F' fzt-widget
