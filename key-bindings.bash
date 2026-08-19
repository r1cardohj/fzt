# fzt key bindings for bash —— 在 ~/.bashrc 里:
#   source /path/to/fzt/key-bindings.bash
#
# Ctrl-F: 打开 fzt 文件树，选中的路径插入到当前命令行光标处
# （类似 fzf 官方的 Ctrl-T）

__fzt_select__() {
  local selected
  selected="$(command fzt "$@" 2>/dev/null)" || return
  printf -v selected '%q' "$selected" # shell 转义，防空格/特殊字符
  READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${selected}${READLINE_LINE:$READLINE_POINT}"
  READLINE_POINT=$((READLINE_POINT + ${#selected}))
}

bind -m emacs-standard -x '"\C-f": __fzt_select__'
bind -m vi-command -x '"\C-f": __fzt_select__'
bind -m vi-insert -x '"\C-f": __fzt_select__'
