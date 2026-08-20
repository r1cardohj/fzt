#!/bin/sh
# fzt installer for Linux/macOS.
#
#   curl -fsSL https://raw.githubusercontent.com/r1cardohj/fzt/main/install.sh | sh
#   curl -fsSL ... | sh -s -- v0.0.5    # pin a version
#   PREFIX=~/.local sh install.sh       # install without sudo
set -eu

REPO="r1cardohj/fzt"
VERSION="${1:-latest}"
PREFIX="${PREFIX:-/usr/local}"

case "$(uname -s)" in
  Linux)  os=linux ;;
  Darwin) os=darwin ;;
  *) echo "fzt: unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64 | amd64)  arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *) echo "fzt: unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

if [ "$VERSION" = latest ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
    sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
  [ -n "$VERSION" ] || { echo "fzt: cannot resolve the latest release" >&2; exit 1; }
fi

name="fzt-$VERSION-$os-$arch"
base="https://github.com/$REPO/releases/download/$VERSION"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "==> Downloading $base/$name.tar.gz"
curl -fsSL "$base/$name.tar.gz" -o "$tmp/$name.tar.gz"
curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt"

echo "==> Verifying checksum"
expected=$(grep " $name\.tar.gz\$" "$tmp/checksums.txt" | awk '{print $1}')
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$tmp/$name.tar.gz" | awk '{print $1}')
else # macOS
  actual=$(shasum -a 256 "$tmp/$name.tar.gz" | awk '{print $1}')
fi
if [ -z "$expected" ] || [ "$expected" != "$actual" ]; then
  echo "fzt: checksum mismatch for $name.tar.gz" >&2
  exit 1
fi

tar xzf "$tmp/$name.tar.gz" -C "$tmp"
sudo=""
[ -w "$PREFIX/bin" ] || sudo="sudo"
$sudo install -m 0755 "$tmp/$name/fzt" "$PREFIX/bin/fzt"
echo "==> Installed fzt $VERSION to $PREFIX/bin/fzt"
