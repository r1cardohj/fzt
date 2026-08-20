# fzt key bindings for zsh —— 在 ~/.zshrc 里:
#   eval "$(fzt --zsh)"
#
# Ctrl-F: 打开 fzt 文件树，选中的路径插入到当前命令行光标处

fzt-widget() {
  local selected
  selected="$(command fzt --height=60% 2>/dev/null)" || { zle reset-prompt; return 1; }
  [ -n "$selected" ] || { zle reset-prompt; return 1; } # 防御:退出码 0 但无输出时不插入
  LBUFFER+="${(q)selected}"
  zle reset-prompt
  return 0
}
zle -N fzt-widget
bindkey '^F' fzt-widget
