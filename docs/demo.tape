# fzt demo recording. Regenerate docs/demo.gif from the repo root with:
#   vhs docs/demo.tape && docs/demo-overlay.sh
# The demo tree under /tmp/fzt-demo mirrors the layout of `make test`.
# Key hints are not recorded; demo-overlay.sh overlays them with ffmpeg.
#
# Two rounds: tree mode first (move/expand/select), then search mode.

Output docs/demo.gif

Set Shell bash
Set FontSize 16
Set Width 1200
Set Height 700
Set Padding 16
Set Framerate 24
Set TypingSpeed 100ms
Set Theme "Catppuccin Mocha"

Hide
Type "export PATH=/home/houjun/project/fzt:$PATH"
Enter
Type "cd /tmp/fzt-demo && clear"
Enter
Sleep 500ms
Show

Type 'eval "$(fzt --bash)"'
Enter
Sleep 1s

# Round 1: tree mode — move, expand, select a file
Type "nvim "
Sleep 800ms
Ctrl+F
Sleep 2.5s
Type "j"
Sleep 1.2s
Enter
Sleep 1.5s
Type "jj"
Sleep 1.2s
Enter
Sleep 2s

# Transition: clear the command line
Ctrl+U
Sleep 1s

# Round 2: search mode — fuzzy match straight into a nested path
Type "nvim "
Sleep 800ms
Ctrl+F
Sleep 2.5s
Type "/"
Sleep 800ms
Type "config"
Sleep 1.8s
Enter
Sleep 2s
