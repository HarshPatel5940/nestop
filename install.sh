#!/usr/bin/env bash
set -euo pipefail

REPO="harshpatel5940/nestop"
BIN_NAME="nestop"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
  Darwin) goos="darwin" ;;
  Linux) goos="linux" ;;
  *)
    echo "Unsupported OS: $os. Download manually from https://github.com/${REPO}/releases" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64) goarch="amd64" ;;
  arm64|aarch64) goarch="arm64" ;;
  *)
    echo "Unsupported arch: $arch. Download manually from https://github.com/${REPO}/releases" >&2
    exit 1
    ;;
esac

echo "Fetching latest release info..."
tag="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep -m1 '"tag_name"' | cut -d '"' -f4)"
if [ -z "$tag" ]; then
  echo "Could not determine latest release tag" >&2
  exit 1
fi
version="${tag#v}"

archive="${BIN_NAME}_${version}_${goos}_${goarch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${tag}/${archive}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

echo "Downloading ${url}..."
curl -fsSL "$url" -o "$tmpdir/$archive"

tar -xzf "$tmpdir/$archive" -C "$tmpdir"

if [ -w "$INSTALL_DIR" ]; then
  mv "$tmpdir/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
else
  sudo mv "$tmpdir/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
fi
chmod +x "$INSTALL_DIR/$BIN_NAME"

echo "Installed ${BIN_NAME} ${tag} to ${INSTALL_DIR}/${BIN_NAME}"
"$INSTALL_DIR/$BIN_NAME" --help >/dev/null 2>&1 || true
