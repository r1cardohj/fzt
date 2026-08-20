#!/bin/sh
# Overlay key-hint labels onto docs/demo.gif (recorded with docs/demo.tape).
# Run from the repo root. Time ranges mirror the sleeps in docs/demo.tape;
# adjust both together.
set -eu

FONT=/usr/share/fonts/truetype/dejavu/DejaVuSansMono-Bold.ttf
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# label text + start/end seconds, one per line
cat > "$tmp/segments" <<'EOF'
type: nvim|3.1|4.4
Ctrl+F — tree mode|4.4|6.9
j — move|6.9|8.2
Enter — expand dir|8.2|9.7
j j — move|9.7|11.9
Enter — insert path|11.9|13.9
Ctrl+U — clear line|13.9|15.0
type: nvim|15.0|16.3
Ctrl+F — search mode|16.3|18.8
/ — search|18.8|19.7
type: config|19.7|22.1
Enter — insert path|22.1|24.2
EOF

filters=""
i=0
while IFS='|' read -r label start end; do
  f="$tmp/k$i.txt"
  printf '%s' "$label" > "$f"
  filters="${filters}drawtext=fontfile=$FONT:textfile=$f:fontsize=26:fontcolor=white:box=1:boxcolor=black@0.55:boxborderw=10:x=24:y=h-th-24:enable='between(t,$start,$end)',"
  i=$((i + 1))
done < "$tmp/segments"

# two-pass palette keeps the re-encoded GIF crisp
ffmpeg -v error -i docs/demo.gif \
  -vf "${filters%,},split[a][b];[a]palettegen[p];[b][p]paletteuse" \
  -y -f gif docs/demo.gif.tmp && mv docs/demo.gif.tmp docs/demo.gif
