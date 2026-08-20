# fzt demo recording. Regenerate docs/demo.gif from the repo root with:
#   sh docs/demo-setup.sh && vhs docs/demo.tape
#
# Three scenes, all executed for real:
#   1. cat <Ctrl-F> pick a nested file in tree mode, run it
#   2. rm <Ctrl-F> pick a file in search mode, run it (then ls proves it)
#   3. wc -l $(fzt) and fzt | xargs wc -l

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
Type 'eval "$(fzt --bash)"'
Enter
Type "cd /tmp/fzt-demo && clear"
Enter
Sleep 500ms
Show

# Scene 1: cat a nested file, picked in tree mode
Type "cat "
Sleep 600ms
Ctrl+F
Sleep 2s
Type "j"
Sleep 800ms
Enter
Sleep 1s
Type "j"
Sleep 800ms
Enter
Sleep 1s
Type "j"
Sleep 800ms
Enter
Sleep 1.2s
Enter
Sleep 2s

# Scene 2: rm a file, found via search mode
Type "rm "
Sleep 600ms
Ctrl+F
Sleep 2s
Type "/"
Sleep 600ms
Type "test_m"
Sleep 1.2s
Enter
Sleep 1.2s
Enter
Sleep 1.5s
Type "ls tests"
Enter
Sleep 1.5s

# Scene 3: fzt inside command substitution and pipes
Type "wc -l $(fzt)"
Sleep 600ms
Enter
Sleep 2s
Type "jjj"
Sleep 800ms
Enter
Sleep 2s
Type "fzt | xargs wc -l"
Sleep 600ms
Enter
Sleep 2s
Enter
Sleep 600ms
Type "j"
Sleep 800ms
Enter
Sleep 2s
